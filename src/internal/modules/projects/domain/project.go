package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProjectNotFound = errors.New("project not found")
)

// Project represents a Notion database that is being synchronized
type Project struct {
	ID                  uuid.UUID // Internal UUID for DB relations and ordering
	PublicID            string    // Public ID with prefix for API
	UserID              uuid.UUID
	NotionDatabaseID    string
	NotionWebhookSecret string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// NewProject creates a new project with validation
func NewProject(userID uuid.UUID, notionDatabaseID, notionWebhookSecret string, idGen IDGenerator, clock Clock) (Project, error) {
	if userID == uuid.Nil {
		return Project{}, errors.New("invalid user ID")
	}
	if notionDatabaseID == "" {
		return Project{}, errors.New("notion database ID cannot be empty")
	}
	if notionWebhookSecret == "" {
		return Project{}, errors.New("notion webhook secret cannot be empty")
	}

	now := clock.Now()

	return Project{
		ID:                  uuid.New(),
		PublicID:            idGen.NewID("project"),
		UserID:              userID,
		NotionDatabaseID:    notionDatabaseID,
		NotionWebhookSecret: notionWebhookSecret,
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}

// Clock interface for dependency injection
type Clock interface {
	Now() time.Time
}

// IDGenerator provides unique ID generation for domain entities
type IDGenerator interface {
	NewID(prefix string) string
}
