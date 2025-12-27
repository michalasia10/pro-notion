package http_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAuthRouter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Router Suite")
}
