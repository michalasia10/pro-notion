package database

import (
	"database/sql"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	gormOnce sync.Once
	gormDB   *gorm.DB
)

// ensureSQL ensures that the sql.DB has been initialized via New().
func ensureSQL() *sql.DB {
	if dbInstance == nil {
		_ = New() // initialize singleton
	}
	return dbInstance.db
}

// GormDB returns a singleton *gorm.DB connected to the same database.
func GormDB() *gorm.DB {
	gormOnce.Do(func() {
		var err error
		// Reuse the same *sql.DB pool to avoid duplicate connections
		sqlDB := ensureSQL()
		gormDB, err = gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to open gorm connection: %v", err)
		}
	})
	return gormDB
}

// SQLDB exposes the underlying *sql.DB for libraries like goose.
func SQLDB() *sql.DB {
	return ensureSQL()
}

// Migrator exposes GORM's migrator for migrations.
func Migrator() gorm.Migrator {
	return GormDB().Migrator()
}
