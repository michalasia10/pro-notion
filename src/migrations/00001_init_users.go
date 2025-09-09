package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"

	"src/internal/database"
	userpg "src/internal/modules/users/infrastructure/postgres"
)

func init() {
	goose.AddMigrationContext(upInitUsers, downInitUsers)
}

func upInitUsers(ctx context.Context, _ *sql.Tx) error {
	m := database.Migrator()
	return m.AutoMigrate(&userpg.UserRecord{})
}

func downInitUsers(ctx context.Context, _ *sql.Tx) error {
	m := database.Migrator()
	return m.DropTable(&userpg.UserRecord{})
}
