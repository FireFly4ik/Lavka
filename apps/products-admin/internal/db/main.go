package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(url string) error {
	var err error
	DB, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %w", err)
	}

	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return fmt.Errorf("Failed to create extension uuid-ossp: %w", err)
	}

	if err := DB.AutoMigrate(
	//TODO вставить моделм
	); err != nil {
		return fmt.Errorf("Failed to run migrations: %w", err)
	}

	return nil
}
