package http_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProjectsRouter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Projects Router Suite")
}
