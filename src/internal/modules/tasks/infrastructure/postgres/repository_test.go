package postgres_test

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	"src/internal/config"
	"src/internal/database"
	shared "src/internal/modules/shared/domain"
	taskdomain "src/internal/modules/tasks/domain"
	taskpostgres "src/internal/modules/tasks/infrastructure/postgres"
)

var (
	pgContainer testcontainers.Container
	db          *gorm.DB
	repo        *taskpostgres.Repository
)

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	const (
		dbName = "tasks_test"
		dbUser = "tasks_user"
		dbPass = "tasks_pass"
	)

	container, err := postgres.Run(
		context.Background(),
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, err
	}

	host, err := container.Host(context.Background())
	if err != nil {
		return container.Terminate, err
	}

	port, err := container.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return container.Terminate, err
	}

	testCfg := &config.Config{}
	testCfg.Database.Host = host
	testCfg.Database.Port = port.Port()
	testCfg.Database.Username = dbUser
	testCfg.Database.Password = dbPass
	testCfg.Database.Database = dbName
	testCfg.Database.Schema = "public"
	testCfg.Port = 8080

	config.SetForTests(testCfg)

	return container.Terminate, nil
}

var teardown func(context.Context, ...testcontainers.TerminateOption) error

var _ = BeforeSuite(func() {
	var err error
	teardown, err = mustStartPostgresContainer()
	Expect(err).NotTo(HaveOccurred())

	db = database.GormDB()

	// Apply migration for tasks table
	migrator := database.Migrator()
	Expect(migrator.AutoMigrate(&taskpostgres.TaskRecord{})).To(Succeed())

	repo = taskpostgres.NewRepository()
})

var _ = AfterSuite(func() {
	if teardown != nil {
		if err := teardown(context.Background()); err != nil {
			log.Printf("could not teardown postgres container: %v", err)
		}
	}
})

var _ = Describe("Tasks Repository", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		db.Exec("TRUNCATE TABLE tasks RESTART IDENTITY CASCADE")
	})

	Describe("Save and retrieval", func() {
		It("saves and fetches a task by ID", func() {
			task := buildTask()
			Expect(repo.Save(ctx, task)).To(Succeed())

			found, err := repo.FindByID(ctx, task.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found.ID).To(Equal(task.ID))
			Expect(found.NotionDatabaseID).To(Equal(task.NotionDatabaseID))
			Expect(found.NotionPageID).To(Equal(task.NotionPageID))
		})

		It("finds by public ID", func() {
			task := buildTask()
			Expect(repo.Save(ctx, task)).To(Succeed())

			found, err := repo.FindByPublicID(ctx, task.PublicID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found.ID).To(Equal(task.ID))
		})

		It("finds by Notion identifiers", func() {
			task := buildTask()
			Expect(repo.Save(ctx, task)).To(Succeed())

			found, err := repo.FindByNotionIDs(ctx, task.NotionDatabaseID, task.NotionPageID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found.PublicID).To(Equal(task.PublicID))
		})

		It("returns ErrTaskNotFound when missing", func() {
			_, err := repo.FindByNotionIDs(ctx, "missing", "missing")
			Expect(err).To(MatchError(taskdomain.ErrTaskNotFound))
		})
	})

	Describe("Upsert", func() {
		It("inserts when no record exists", func() {
			task := buildTask()
			Expect(repo.Upsert(ctx, task)).To(Succeed())

			found, err := repo.FindByNotionIDs(ctx, task.NotionDatabaseID, task.NotionPageID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found.ID).To(Equal(task.ID))
		})

		It("updates existing record on conflict", func() {
			task := buildTask()
			Expect(repo.Save(ctx, task)).To(Succeed())

			task.PropertiesHash = "updated"
			task.SyncStatus = taskdomain.SyncStatusFailed
			updateTime := time.Now().Add(time.Minute)
			task.UpdatedAt = updateTime

			tExpect := repo.Upsert(ctx, task)
			Expect(tExpect).To(Succeed())

			found, err := repo.FindByNotionIDs(ctx, task.NotionDatabaseID, task.NotionPageID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found.PropertiesHash).To(Equal("updated"))
			Expect(found.SyncStatus).To(Equal(taskdomain.SyncStatusFailed))
		})
	})

	Describe("WithTransaction", func() {
		It("commits when callback succeeds", func() {
			task := buildTask()
			Expect(repo.WithTransaction(ctx, func(ctx context.Context, txRepo taskdomain.Repository) error {
				return txRepo.Save(ctx, task)
			})).To(Succeed())

			found, err := repo.FindByID(ctx, task.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(found.ID).To(Equal(task.ID))
		})

		It("rolls back when callback fails", func() {
			task := buildTask()
			errTest := repo.WithTransaction(ctx, func(ctx context.Context, txRepo taskdomain.Repository) error {
				Expect(txRepo.Save(ctx, task)).To(Succeed())
				return assertErr
			})
			Expect(errTest).To(HaveOccurred())

			_, err := repo.FindByID(ctx, task.ID)
			Expect(err).To(MatchError(taskdomain.ErrTaskNotFound))
		})
	})
})

var assertErr = errors.New("force rollback")

func buildTask() *taskdomain.Task {
	idGen := shared.NewUUIDGenerator()
	now := time.Now().UTC().Truncate(time.Second)
	task, err := taskdomain.NewTask(uuid.New(), "db"+uuid.New().String()[:6], "page"+uuid.New().String()[:6], now, idGen)
	if err != nil {
		panic(err)
	}
	return task
}
