package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"src/internal/modules/projects/domain"
)

// Mock implementations for testing
type mockClock struct {
	now time.Time
}

func (m *mockClock) Now() time.Time {
	return m.now
}

type mockIDGenerator struct {
	id string
}

func (m *mockIDGenerator) NewID(prefix string) string {
	return prefix + "_" + m.id
}

var _ = Describe("Project", func() {
	Describe("NewProject", func() {
		var (
			userID              uuid.UUID
			notionDatabaseID    string
			notionWebhookSecret string
			clock               domain.Clock
			idGen               domain.IDGenerator
		)

		BeforeEach(func() {
			userID = uuid.New()
			notionDatabaseID = "database_123"
			notionWebhookSecret = "secret_123"
			clock = &mockClock{now: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)}
			idGen = &mockIDGenerator{id: "test_id"}
		})

		It("should create a project successfully", func() {
			project, err := domain.NewProject(userID, notionDatabaseID, notionWebhookSecret, idGen, clock)

			Expect(err).ToNot(HaveOccurred())
			Expect(project.ID).ToNot(BeNil())
			Expect(project.PublicID).To(Equal("project_test_id"))
			Expect(project.UserID).To(Equal(userID))
			Expect(project.NotionDatabaseID).To(Equal(notionDatabaseID))
			Expect(project.NotionWebhookSecret).To(Equal(notionWebhookSecret))
			Expect(project.CreatedAt).To(Equal(clock.Now()))
			Expect(project.UpdatedAt).To(Equal(clock.Now()))
		})

		It("should return error for empty user ID", func() {
			_, err := domain.NewProject(uuid.Nil, notionDatabaseID, notionWebhookSecret, idGen, clock)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid user ID"))
		})

		It("should return error for empty notion database ID", func() {
			_, err := domain.NewProject(userID, "", notionWebhookSecret, idGen, clock)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("notion database ID cannot be empty"))
		})

		It("should return error for empty webhook secret", func() {
			_, err := domain.NewProject(userID, notionDatabaseID, "", idGen, clock)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("notion webhook secret cannot be empty"))
		})
	})
})

func TestProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Project Domain Suite")
}
