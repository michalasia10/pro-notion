package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"

	"src/internal/database"
	projectpg "src/internal/modules/projects/infrastructure/postgres"
)

func init() {
	goose.AddMigrationContext(upCreateProjects, downCreateProjects)
}

func upCreateProjects(ctx context.Context, _ *sql.Tx) error {
	m := database.Migrator()
	return m.AutoMigrate(&projectpg.ProjectRecord{})
}

func downCreateProjects(ctx context.Context, _ *sql.Tx) error {
	m := database.Migrator()
	return m.DropTable(&projectpg.ProjectRecord{})
}
