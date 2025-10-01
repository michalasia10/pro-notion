package domain

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines persistence operations for the Task aggregate.
type Repository interface {
	// Save persists a new task.
	Save(ctx context.Context, task *Task) error

	// Update persists changes to an existing task.
	Update(ctx context.Context, task *Task) error

	// Upsert inserts or updates a task based on Notion identifiers.
	Upsert(ctx context.Context, task *Task) error

	// FindByID retrieves a task by canonical ID.
	FindByID(ctx context.Context, id uuid.UUID) (*Task, error)

	// FindByPublicID retrieves a task by its public identifier.
	FindByPublicID(ctx context.Context, publicID string) (*Task, error)

	// FindByNotionIDs retrieves a task by Notion database/page IDs.
	FindByNotionIDs(ctx context.Context, notionDatabaseID, notionPageID string) (*Task, error)

	// WithTransaction executes the provided function within a transaction boundary
	// and provides a repository instance bound to that transaction.
	WithTransaction(ctx context.Context, fn func(ctx context.Context, repo Repository) error) error
}
