package postgres

import (
	"time"

	"github.com/google/uuid"

	taskdomain "src/internal/modules/tasks/domain"
)

// TaskRecord represents the tasks table schema for GORM.
type TaskRecord struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	PublicID         string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	ProjectID        uuid.UUID `gorm:"type:uuid;index;not null"`
	NotionPageID     string    `gorm:"type:varchar(255);index:idx_tasks_notion_unique,unique,not null"`
	NotionDatabaseID string    `gorm:"type:varchar(255);index:idx_tasks_notion_unique,unique,not null"`
	PropertiesHash   string    `gorm:"type:text"`
	SyncStatus       string    `gorm:"type:varchar(32);not null"`
	LastSyncedAt     *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `gorm:"index"`
}

// TableName returns the table name for tasks.
func (TaskRecord) TableName() string {
	return "tasks"
}

// ToDomain converts a TaskRecord into a domain Task.
func (r TaskRecord) ToDomain() (*taskdomain.Task, error) {
	return &taskdomain.Task{
		ID:               r.ID,
		PublicID:         r.PublicID,
		ProjectID:        r.ProjectID,
		NotionDatabaseID: r.NotionDatabaseID,
		NotionPageID:     r.NotionPageID,
		PropertiesHash:   r.PropertiesHash,
		SyncStatus:       taskdomain.SyncStatus(r.SyncStatus),
		LastSyncedAt:     r.LastSyncedAt,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
		DeletedAt:        r.DeletedAt,
	}, nil
}

// FromDomain converts a domain Task into a TaskRecord.
func FromDomain(task *taskdomain.Task) TaskRecord {
	return TaskRecord{
		ID:               task.ID,
		PublicID:         task.PublicID,
		ProjectID:        task.ProjectID,
		NotionDatabaseID: task.NotionDatabaseID,
		NotionPageID:     task.NotionPageID,
		PropertiesHash:   task.PropertiesHash,
		SyncStatus:       string(task.SyncStatus),
		LastSyncedAt:     task.LastSyncedAt,
		CreatedAt:        task.CreatedAt,
		UpdatedAt:        task.UpdatedAt,
		DeletedAt:        task.DeletedAt,
	}
}
