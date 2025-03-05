package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/sh1ro/todo-api/pkg/config"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
    
    // Sleep for 10 seconds
    time.Sleep(10 * time.Second)

	// Construct database URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	// Create a new migrate instance
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("Migration failed to initialize: %v", err)
	}

	// Add logging
	m.Log = &MigrateLogger{}

	// Check command line arguments
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/migrate/main.go [up|down|version|create]")
	}

	// Execute the appropriate command
	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
		log.Println("Migrations applied successfully")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		log.Println("Migrations rolled back successfully")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		log.Printf("Current migration version: %d, Dirty: %v", version, dirty)

	case "create":
		if len(os.Args) < 3 {
			log.Fatal("Usage: go run cmd/migrate/main.go create <migration_name>")
		}
		name := os.Args[2]
		timestamp := time.Now().Unix()
		upFile := fmt.Sprintf("migrations/%d_%s.up.sql", timestamp, name)
		downFile := fmt.Sprintf("migrations/%d_%s.down.sql", timestamp, name)

		// Create up migration file
		if err := createFile(upFile, "-- Migration Up\n\n"); err != nil {
			log.Fatalf("Failed to create up migration file: %v", err)
		}

		// Create down migration file
		if err := createFile(downFile, "-- Migration Down\n\n"); err != nil {
			log.Fatalf("Failed to create down migration file: %v", err)
		}

		log.Printf("Created migration files: %s, %s", upFile, downFile)

	default:
		log.Fatal("Usage: go run cmd/migrate/main.go [up|down|version|create]")
	}
}

// MigrateLogger implements migrate.Logger interface
type MigrateLogger struct{}

// Printf implements migrate.Logger interface
func (l *MigrateLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Verbose implements migrate.Logger interface
func (l *MigrateLogger) Verbose() bool {
	return true
}

// Create a new file with the given content
func createFile(filename, content string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}
