package application_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUsersApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Users Application Suite")
}
