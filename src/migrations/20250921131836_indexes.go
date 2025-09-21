package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upIndexes, downIndexes)
}

func upIndexes(ctx context.Context, tx *sql.Tx) error {
	// Performance composite indexes for better query performance

	// Projects: Index for user projects ordered by creation date (most common query)
	_, err := tx.ExecContext(ctx, `
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_projects_user_created 
		ON projects (user_id, created_at DESC);
	`)
	if err != nil {
		return err
	}

	// Projects: Index for active projects (excluding soft-deleted)
	_, err = tx.ExecContext(ctx, `
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_projects_active 
		ON projects (user_id, created_at DESC) 
		WHERE deleted_at IS NULL;
	`)
	if err != nil {
		return err
	}

	// Users: Index for notion workspace queries with email lookup
	_, err = tx.ExecContext(ctx, `
		CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_notion_workspace 
		ON users (notion_workspace_id, email);
	`)
	if err != nil {
		return err
	}

	return nil
}

func downIndexes(ctx context.Context, tx *sql.Tx) error {
	// Drop the composite indexes created in upIndexes

	_, err := tx.ExecContext(ctx, `DROP INDEX CONCURRENTLY IF EXISTS idx_projects_user_created;`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `DROP INDEX CONCURRENTLY IF EXISTS idx_projects_active;`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `DROP INDEX CONCURRENTLY IF EXISTS idx_users_notion_workspace;`)
	if err != nil {
		return err
	}

	return nil
}
