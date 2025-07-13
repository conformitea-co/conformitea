package database

import (
	"conformitea/server/migrations"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Runs all pending migrations automatically
func RunMigrations() error {
	// Create migration source from embedded files
	source, err := iofs.New(migrations.MigrationFiles, ".")
	if err != nil {
		return err
	}

	// Create database driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// Create migrator
	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return err
	}

	// Run migrations
	if err := m.Up(); !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	// Close the migrator
	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		return sourceErr
	}
	if dbErr != nil {
		return dbErr
	}

	return nil
}
