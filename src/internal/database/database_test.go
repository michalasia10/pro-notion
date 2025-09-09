package database

import (
	"context"
	"log"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"src/internal/config"
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

	testConfig.Redis.Host = "localhost"
	testConfig.Redis.Port = "6379"
	testConfig.Redis.Password = ""

	// Set test config
	config.SetForTests(testConfig)

	return dbContainer.Terminate, err
}

var _ = BeforeSuite(func() {
	var err error
	teardown, err = mustStartPostgresContainer()
	Expect(err).ToNot(HaveOccurred())
})

var teardown func(context.Context, ...testcontainers.TerminateOption) error

var _ = AfterSuite(func() {
	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container")
	}
})

var _ = Describe("Database", func() {
	BeforeEach(func() {
		// Reset singleton for each test
		dbInstance = nil
	})

	It("New returns non-nil", func() {
		srv := New()
		Expect(srv).ToNot(BeNil())
	})

	It("Health returns up/healthy", func() {
		srv := New()
		stats := srv.Health()
		Expect(stats["status"]).To(Equal("up"))
		_, hasErr := stats["error"]
		Expect(hasErr).To(BeFalse())
		Expect(stats["message"]).To(Equal("It's healthy"))
	})

	It("Close returns nil", func() {
		srv := New()
		Expect(srv.Close()).To(BeNil())
	})
})
