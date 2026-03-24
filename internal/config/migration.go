package config

import (
	"errors"
	"log"

	"go-e-commerce/db"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"
)

// RunMigrations executes database migrations automatically using the embedded SQL files
func RunMigrations(gormDB *gorm.DB) {
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("failed to get sql db for migrations: %v", err)
	}

	sourceDriver, err := iofs.New(db.MigrationFS, "migrations")
	if err != nil {
		log.Fatalf("failed to create migration source driver: %v", err)
	}

	dbDriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create database driver: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("failed to run migrate up: %v", err)
	}

	log.Println("Database migrations successfully checked and applied")
}
