package db

import (
	"database/sql"
	"fmt"

	"github.com/budhilaw/personal-website-backend/config"
	"github.com/budhilaw/personal-website-backend/internal/logger"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

// InitDB initializes the database connection
func InitDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.GetPostgresConnString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established",
		zap.String("host", cfg.PostgresHost),
		zap.String("database", cfg.PostgresDB))

	return db, nil
}

// RunMigrations runs the database migrations
func RunMigrations(db *sql.DB) error {
	// Set up goose
	goose.SetBaseFS(nil)

	// Set the migration directory
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Run migrations
	if err := goose.Up(db, "db/migration"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database migrations completed")
	return nil
}

// CreateMigration creates a new migration file
func CreateMigration(name string, sql bool) error {
	var ext string
	if sql {
		ext = "sql"
	} else {
		ext = "go"
	}

	if err := goose.Create(nil, "db/migration", name, ext); err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	return nil
}
