package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/budhilaw/personal-website-backend/config"
	"github.com/budhilaw/personal-website-backend/db"
	"github.com/budhilaw/personal-website-backend/internal/logger"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

// handleDBCommand handles database migration commands
func handleDBCommand() {
	// Check if command is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/api/main.go [db:migrate|db:create|db:rollback|db:reset]")
		os.Exit(1)
	}

	command := os.Args[1]

	// Process command
	switch command {
	case "db:migrate":
		runMigrations()
	case "db:create":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run cmd/api/main.go db:create <migration_name>")
			os.Exit(1)
		}
		createMigration(os.Args[2], true) // true for SQL migration
	case "db:rollback":
		rollbackMigration()
	case "db:reset":
		resetDatabase()
	default:
		// If not a db command, return to continue with normal app flow
		return
	}

	// Exit after handling command
	os.Exit(0)
}

// runMigrations runs all database migrations
func runMigrations() {
	cfg := config.InitConfig()

	// Initialize logger
	_ = logger.InitLogger(cfg.IsProduction())

	database, err := db.InitDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	if err := db.RunMigrations(database); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	logger.Info("Migrations completed successfully")
}

// createMigration creates a new migration file
func createMigration(name string, sql bool) {
	// Initialize logger with development mode
	_ = logger.InitLogger(false)

	name = strings.ReplaceAll(name, " ", "_")
	err := db.CreateMigration(name, sql)
	if err != nil {
		logger.Fatal("Failed to create migration", zap.Error(err), zap.String("name", name))
	}

	logger.Info("Migration created", zap.String("name", name))
}

// rollbackMigration rolls back the most recent migration
func rollbackMigration() {
	cfg := config.InitConfig()

	// Initialize logger
	_ = logger.InitLogger(cfg.IsProduction())

	database, err := db.InitDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	if err := rollback(database); err != nil {
		logger.Fatal("Failed to rollback migration", zap.Error(err))
	}

	logger.Info("Rollback completed successfully")
}

// resetDatabase drops all tables and reruns migrations
func resetDatabase() {
	cfg := config.InitConfig()

	// Initialize logger
	_ = logger.InitLogger(cfg.IsProduction())

	// Connect to database
	database, err := db.InitDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	// First down all migrations
	logger.Info("Reverting all migrations")
	if err := goose.Reset(database, "db/migration"); err != nil {
		logger.Fatal("Failed to reset migrations", zap.Error(err))
	}

	// Then run migrations again
	logger.Info("Applying migrations")
	if err := db.RunMigrations(database); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	logger.Info("Database reset completed successfully")
}

// rollback rolls back the most recent migration
func rollback(db *sql.DB) error {
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Down(db, "db/migration"); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}
