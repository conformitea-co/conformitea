package database

import (
	"database/sql"
	"fmt"
	"time"

	"conformitea/server/internal/config"
	"conformitea/server/internal/logger"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

// Initializes the database connection pool
func Initialize() error {
	cfg := config.GetConfig().Database

	if _db, err := sql.Open("postgres", cfg.URL); err != nil {
		return err
	} else {
		db = _db
	}

	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return err
	}
	logger.GetLogger().Info("database connection established successfully")

	if err := RunMigrations(); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}
	logger.GetLogger().Info("database migrations completed successfully")

	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
