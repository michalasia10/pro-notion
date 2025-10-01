package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	shared "src/internal/modules/shared/domain"
	tasks "src/internal/modules/tasks/domain"
)

func TestTaskDomain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tasks Domain Suite")
}

var _ = Describe("Task", func() {
	var (
		projectID uuid.UUID
		now       time.Time
		idGen     shared.IDGenerator
	)

	BeforeEach(func() {
		projectID = uuid.New()
		now = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
		idGen = shared.NewUUIDGenerator()
	})

	Describe("NewTask", func() {
		It("creates a task with initial pending sync status", func() {
			task, err := tasks.NewTask(projectID, "db123", "page456", now, idGen)
			Expect(err).NotTo(HaveOccurred())
			Expect(task.ID).NotTo(Equal(uuid.Nil))
			Expect(task.PublicID).To(ContainSubstring("task_"))
			Expect(task.ProjectID).To(Equal(projectID))
			Expect(task.NotionDatabaseID).To(Equal("db123"))
			Expect(task.NotionPageID).To(Equal("page456"))
			Expect(task.SyncStatus).To(Equal(tasks.SyncStatusPending))
			Expect(task.CreatedAt).To(Equal(now))
			Expect(task.UpdatedAt).To(Equal(now))
			Expect(task.LastSyncedAt).To(BeNil())
		})

		It("returns error when projectID is empty", func() {
			task, err := tasks.NewTask(uuid.Nil, "db123", "page456", now, idGen)
			Expect(err).To(MatchError(tasks.ErrInvalidProjectID))
			Expect(task).To(BeNil())
		})

		It("returns error when notion database ID missing", func() {
			task, err := tasks.NewTask(projectID, "", "page456", now, idGen)
			Expect(err).To(MatchError(tasks.ErrInvalidNotionDatabaseID))
			Expect(task).To(BeNil())
		})

		It("returns error when notion page ID missing", func() {
			task, err := tasks.NewTask(projectID, "db123", "", now, idGen)
			Expect(err).To(MatchError(tasks.ErrInvalidNotionPageID))
			Expect(task).To(BeNil())
		})

		It("returns error when id generator missing", func() {
			task, err := tasks.NewTask(projectID, "db123", "page456", now, nil)
			Expect(err).To(MatchError(tasks.ErrInvalidIDGenerator))
			Expect(task).To(BeNil())
		})
	})

	Describe("MarkSynced", func() {
		It("updates sync metadata", func() {
			task, _ := tasks.NewTask(projectID, "db123", "page456", now, idGen)
			syncTime := now.Add(10 * time.Minute)
			task.MarkSynced(syncTime, "hash123")
			Expect(task.LastSyncedAt).NotTo(BeNil())
			Expect(*task.LastSyncedAt).To(Equal(syncTime))
			Expect(task.PropertiesHash).To(Equal("hash123"))
			Expect(task.SyncStatus).To(Equal(tasks.SyncStatusCompleted))
			Expect(task.UpdatedAt).To(Equal(syncTime))
		})
	})

	Describe("MarkSyncFailed", func() {
		It("marks sync as failed and updates timestamp", func() {
			task, _ := tasks.NewTask(projectID, "db123", "page456", now, idGen)
			failTime := now.Add(5 * time.Minute)
			task.MarkSyncFailed(failTime)
			Expect(task.SyncStatus).To(Equal(tasks.SyncStatusFailed))
			Expect(task.UpdatedAt).To(Equal(failTime))
		})
	})

	Describe("UpdateNotionIdentifiers", func() {
		It("updates both database and page IDs", func() {
			task, _ := tasks.NewTask(projectID, "db123", "page456", now, idGen)
			updateTime := now.Add(2 * time.Hour)
			task.UpdateNotionIdentifiers("db999", "page999", updateTime)
			Expect(task.NotionDatabaseID).To(Equal("db999"))
			Expect(task.NotionPageID).To(Equal("page999"))
			Expect(task.UpdatedAt).To(Equal(updateTime))
		})

		It("keeps existing values when empty strings provided", func() {
			task, _ := tasks.NewTask(projectID, "db123", "page456", now, idGen)
			task.UpdateNotionIdentifiers("", "", now.Add(time.Minute))
			Expect(task.NotionDatabaseID).To(Equal("db123"))
			Expect(task.NotionPageID).To(Equal("page456"))
		})
	})
})
