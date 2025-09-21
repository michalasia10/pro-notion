package domain

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for project data access
type Repository interface {
	// Save persists a project
	Save(ctx context.Context, project *Project) error

	// FindByID retrieves a project by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*Project, error)

	// FindByUserID retrieves all projects for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*Project, error)

	// FindByNotionDatabaseID retrieves a project by Notion database ID
	FindByNotionDatabaseID(ctx context.Context, notionDatabaseID string) (*Project, error)

	// Update updates an existing project
	Update(ctx context.Context, project *Project) error

	// Delete removes a project
	Delete(ctx context.Context, id uuid.UUID) error
}
