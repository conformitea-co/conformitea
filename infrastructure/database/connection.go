package database

import (
	"fmt"
	"time"

	"conformitea/infrastructure/config"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize(dbConfigValues config.DatabaseConfig, logger *zap.Logger) (*gorm.DB, error) {
	if err := dbConfigValues.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	db, err := gorm.Open(postgres.Open(dbConfigValues.URL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxOpenConns(dbConfigValues.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(dbConfigValues.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	logger.Info("database connection established successfully")

	if err := RunMigrations(sqlDB); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}
	logger.Info("database migrations completed successfully")

	return db, nil
}
