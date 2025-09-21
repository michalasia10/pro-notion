package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"src/internal/config"
	"src/internal/database"
	_ "src/migrations" // Import migrations to register them with goose

	"github.com/pressly/goose/v3"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "migrations", "directory with migration files")
)

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[0]

	// Load configuration
	config.Load()

	// Initialize database connection
	db := database.SQLDB()

	// Set goose dialect
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("goose dialect error: %v", err)
	}

	// Note: arguments not used in current implementation but kept for future extensions
	_ = args[1:] // Future: could be used for command-specific arguments

	switch command {
	case "up":
		if err := goose.Up(db, *dir); err != nil {
			log.Fatalf("goose up: %v", err)
		}
		fmt.Println("Migration up completed successfully")

	case "down":
		if err := goose.Down(db, *dir); err != nil {
			log.Fatalf("goose down: %v", err)
		}
		fmt.Println("Migration down completed successfully")

	case "status":
		if err := goose.Status(db, *dir); err != nil {
			log.Fatalf("goose status: %v", err)
		}

	case "version":
		version, err := goose.GetDBVersion(db)
		if err != nil {
			log.Fatalf("goose version: %v", err)
		}
		fmt.Printf("Current database version: %d\n", version)

	case "reset":
		if err := goose.Reset(db, *dir); err != nil {
			log.Fatalf("goose reset: %v", err)
		}
		fmt.Println("Migration reset completed successfully")

	case "redo":
		if err := goose.Redo(db, *dir); err != nil {
			log.Fatalf("goose redo: %v", err)
		}
		fmt.Println("Migration redo completed successfully")

	default:
		log.Printf("%q: no such command", command)
		flags.Usage()
		return
	}
}

func usage() {
	fmt.Println("Usage: migrate COMMAND")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("    up                   Migrate the DB to the most recent version available")
	fmt.Println("    down                 Roll back the version by 1")
	fmt.Println("    status               Dump the migration status for the current DB")
	fmt.Println("    version              Print the current version of the database")
	fmt.Println("    reset                Roll back all migrations")
	fmt.Println("    redo                 Re-run the latest migration")
	fmt.Println()
	fmt.Println("Flags:")
	flags.PrintDefaults()
}
