package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddUsersPublicID, downAddUsersPublicID)
}

func upAddUsersPublicID(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS pgcrypto;`); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		ALTER TABLE users
		ADD COLUMN IF NOT EXISTS public_id varchar(255);
	`); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE users
		SET public_id = 'user_' || gen_random_uuid()
		WHERE public_id IS NULL;
	`); err != nil {
		return err
	}

	_, err := tx.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_public_id
		ON users (public_id);
	`)
	return err
}

func downAddUsersPublicID(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `DROP INDEX IF EXISTS idx_users_public_id;`); err != nil {
		return err
	}

	_, err := tx.ExecContext(ctx, `
		ALTER TABLE users
		DROP COLUMN IF EXISTS public_id;
	`)
	return err
}
