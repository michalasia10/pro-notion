package postgres

import (
	"time"

	"src/internal/modules/projects/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProjectRecord represents the projects table structure in PostgreSQL
type ProjectRecord struct {
	ID                  uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid();index"` // Internal UUID for DB relations and ordering
	PublicID            string         `gorm:"uniqueIndex;type:varchar(255);index"`                  // Public ID with prefix for API
	UserID              uuid.UUID      `gorm:"not null;type:uuid;index"`
	NotionDatabaseID    string         `gorm:"not null;type:varchar(255);uniqueIndex"`
	NotionWebhookSecret string         `gorm:"not null;type:varchar(255)"`
	CreatedAt           time.Time      `gorm:"not null;index"`
	UpdatedAt           time.Time      `gorm:"not null"`
	DeletedAt           gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (ProjectRecord) TableName() string {
	return "projects"
}

// toDomainProject converts a ProjectRecord to a domain Project
func toDomainProject(record ProjectRecord) domain.Project {
	return domain.Project{
		ID:                  record.ID,       // Internal UUID for DB relations and ordering
		PublicID:            record.PublicID, // Public ID with prefix for API
		UserID:              record.UserID,
		NotionDatabaseID:    record.NotionDatabaseID,
		NotionWebhookSecret: record.NotionWebhookSecret,
		CreatedAt:           record.CreatedAt,
		UpdatedAt:           record.UpdatedAt,
	}
}

// toProjectRecord converts a domain Project to a ProjectRecord
func toProjectRecord(project domain.Project) ProjectRecord {
	return ProjectRecord{
		ID:                  project.ID,       // Internal UUID for database relations
		PublicID:            project.PublicID, // Public ID with prefix
		UserID:              project.UserID,
		NotionDatabaseID:    project.NotionDatabaseID,
		NotionWebhookSecret: project.NotionWebhookSecret,
		CreatedAt:           project.CreatedAt,
		UpdatedAt:           project.UpdatedAt,
	}
}
