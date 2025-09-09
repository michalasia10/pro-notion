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
)

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
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

	database = dbName
	password = dbPwd
	username = dbUser

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	host = dbHost
	port = dbPort.Port()

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
