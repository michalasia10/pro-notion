package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"

	"src/internal/database"
	taskpg "src/internal/modules/tasks/infrastructure/postgres"
)

func init() {
	goose.AddMigrationContext(upCreateTasks, downCreateTasks)
}

func upCreateTasks(ctx context.Context, _ *sql.Tx) error {
	m := database.Migrator()
	return m.AutoMigrate(&taskpg.TaskRecord{})
}

func downCreateTasks(ctx context.Context, _ *sql.Tx) error {
	m := database.Migrator()
	return m.DropTable(&taskpg.TaskRecord{})
}
