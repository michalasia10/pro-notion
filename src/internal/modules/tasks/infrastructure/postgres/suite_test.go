package postgres_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTasksPostgres(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tasks Postgres Repository Suite")
}
