package db

import (
	"context"
	"fmt"
	"time"

	"github.com/budhilaw/personal-website-backend/config"
	"github.com/budhilaw/personal-website-backend/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib" // Use pgx driver
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

var (
	// DBPool is the global database connection pool
	DBPool *sqlx.DB
)

// InitDB initializes the database connection pool
func InitDB(cfg config.Config) (*sqlx.DB, error) {
	if DBPool != nil {
		return DBPool, nil
	}

	// Create a connection pool
	db, err := sqlx.Open("pgx", cfg.GetPostgresConnString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pooling
	db.SetMaxOpenConns(25)                  // Maximum number of open connections to the database
	db.SetMaxIdleConns(10)                  // Maximum number of connections in the idle connection pool
	db.SetConnMaxLifetime(time.Hour)        // Maximum time a connection may be reused
	db.SetConnMaxIdleTime(30 * time.Minute) // Maximum time a connection may be idle

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection pool established",
		zap.String("host", cfg.PostgresHost),
		zap.String("database", cfg.PostgresDB),
		zap.Int("max_open_conns", 25),
		zap.Int("max_idle_conns", 10))

	// Set the global pool
	DBPool = db

	return db, nil
}

// GetDB returns the database connection pool
func GetDB() *sqlx.DB {
	return DBPool
}

// CloseDB closes the database connection pool
func CloseDB() error {
	if DBPool != nil {
		logger.Info("Closing database connection pool")
		return DBPool.Close()
	}
	return nil
}

// RunMigrations runs the database migrations
func RunMigrations(db *sqlx.DB) error {
	// Set up goose
	goose.SetBaseFS(nil)

	// Set the migration directory
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Run migrations
	if err := goose.Up(db.DB, "db/migration"); err != nil {
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
