package postgres_test

import (
	"context"
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
	"src/internal/modules/projects/domain"
	projectRepo "src/internal/modules/projects/infrastructure/postgres"
)

var (
	pgContainer testcontainers.Container
	db          *gorm.DB
	repo        *projectRepo.ProjectRepository
)

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	var (
		dbName = "test_database"
		dbPwd  = "test_password"
		dbUser = "test_user"
	)

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	testConfig := &config.Config{
		Port: 8080, // Default for tests
	}
	testConfig.Database.Host = dbHost
	testConfig.Database.Port = dbPort.Port()
	testConfig.Database.Username = dbUser
	testConfig.Database.Password = dbPwd
	testConfig.Database.Database = dbName
	testConfig.Database.Schema = "public"

	config.SetForTests(testConfig)

	return dbContainer.Terminate, err
}

var teardown func(context.Context, ...testcontainers.TerminateOption) error

var _ = BeforeSuite(func() {
	var err error
	teardown, err = mustStartPostgresContainer()
	Expect(err).ToNot(HaveOccurred())

	db = database.GormDB()

	// Run migrations using GORM AutoMigrate for tests
	migrator := database.Migrator()
	if err := migrator.AutoMigrate(&projectRepo.ProjectRecord{}); err != nil {
		Fail("Failed to run AutoMigrate: " + err.Error())
	}

	repo = projectRepo.NewProjectRepository(db)
})

var _ = AfterSuite(func() {
	if teardown != nil {
		if err := teardown(context.Background()); err != nil {
			log.Fatalf("could not teardown postgres container: %v", err)
		}
	}
})

var _ = Describe("ProjectRepository", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		// Clean up database before each test
		db.Exec("TRUNCATE TABLE projects CASCADE")
	})

	Describe("Save and FindByID", func() {
		It("should save and retrieve a project successfully", func() {
			// Create a project using domain constructor
			userID := uuid.New()
			notionDatabaseID := "database_123"
			webhookSecret := "secret_123"

			project, err := domain.NewProject(userID, notionDatabaseID, webhookSecret, &mockIDGenerator{}, &mockClock{})
			Expect(err).ToNot(HaveOccurred())

			// Save project
			err = repo.Save(ctx, &project)
			Expect(err).ToNot(HaveOccurred())

			// Retrieve project
			found, err := repo.FindByID(ctx, project.ID)
			Expect(err).ToNot(HaveOccurred())
			Expect(found).ToNot(BeNil())
			Expect(found.ID).To(Equal(project.ID))
			Expect(found.UserID).To(Equal(userID))
			Expect(found.NotionDatabaseID).To(Equal(notionDatabaseID))
			Expect(found.NotionWebhookSecret).To(Equal(webhookSecret))
		})
	})

	Describe("FindByID", func() {
		It("should return error when project does not exist", func() {
			_, err := repo.FindByID(ctx, uuid.New())
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(domain.ErrProjectNotFound))
		})
	})

	Describe("FindByUserID", func() {
		It("should find all projects for a user", func() {
			userID1 := uuid.New()
			userID2 := uuid.New()

			// Use separate ID generators to avoid conflicts
			idGen1 := &mockIDGenerator{counter: 0}
			idGen2 := &mockIDGenerator{counter: 10}
			idGen3 := &mockIDGenerator{counter: 20}

			// Create projects for user1
			project1, _ := domain.NewProject(userID1, "db1", "secret1", idGen1, &mockClock{})
			project2, _ := domain.NewProject(userID1, "db2", "secret2", idGen2, &mockClock{})
			// Create project for user2
			project3, _ := domain.NewProject(userID2, "db3", "secret3", idGen3, &mockClock{})

			// Save all projects
			err := repo.Save(ctx, &project1)
			Expect(err).ToNot(HaveOccurred())
			err = repo.Save(ctx, &project2)
			Expect(err).ToNot(HaveOccurred())
			err = repo.Save(ctx, &project3)
			Expect(err).ToNot(HaveOccurred())

			// Find projects for user1
			projects, err := repo.FindByUserID(ctx, userID1)
			Expect(err).ToNot(HaveOccurred())
			Expect(projects).To(HaveLen(2))

			// Check that both projects belong to user1
			for _, p := range projects {
				Expect(p.UserID).To(Equal(userID1))
			}
		})
	})

	Describe("FindByNotionDatabaseID", func() {
		It("should find project by Notion database ID", func() {
			userID := uuid.New()
			notionDBID := "notion_db_123"

			project, _ := domain.NewProject(userID, notionDBID, "secret", &mockIDGenerator{}, &mockClock{})
			repo.Save(ctx, &project)

			found, err := repo.FindByNotionDatabaseID(ctx, notionDBID)
			Expect(err).ToNot(HaveOccurred())
			Expect(found).ToNot(BeNil())
			Expect(found.ID).To(Equal(project.ID))
			Expect(found.NotionDatabaseID).To(Equal(notionDBID))
		})

		It("should return error when database ID does not exist", func() {
			_, err := repo.FindByNotionDatabaseID(ctx, "nonexistent")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(domain.ErrProjectNotFound))
		})
	})

	Describe("Update", func() {
		It("should update an existing project", func() {
			userID := uuid.New()
			project, _ := domain.NewProject(userID, "db1", "secret1", &mockIDGenerator{}, &mockClock{})
			repo.Save(ctx, &project)

			// Update project
			project.NotionWebhookSecret = "updated_secret"
			err := repo.Update(ctx, &project)
			Expect(err).ToNot(HaveOccurred())

			// Verify update
			found, err := repo.FindByID(ctx, project.ID)
			Expect(err).ToNot(HaveOccurred())
			Expect(found.NotionWebhookSecret).To(Equal("updated_secret"))
		})
	})

	Describe("Delete", func() {
		It("should delete a project", func() {
			userID := uuid.New()
			project, _ := domain.NewProject(userID, "db1", "secret1", &mockIDGenerator{}, &mockClock{})
			repo.Save(ctx, &project)

			// Delete project
			err := repo.Delete(ctx, project.ID)
			Expect(err).ToNot(HaveOccurred())

			// Verify deletion
			_, err = repo.FindByID(ctx, project.ID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(domain.ErrProjectNotFound))
		})
	})
})

// Mock implementations for tests
type mockClock struct{}

func (m *mockClock) Now() time.Time {
	return time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
}

type mockIDGenerator struct {
	counter int
}

func (m *mockIDGenerator) NewID(prefix string) string {
	m.counter++
	return prefix + "_test_" + string(rune(m.counter+'0'))
}
