package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/FireFly4ik/Lavka-products-customer/internal/models"
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
		&models.Product{},
		&models.Category{},
		&models.Market_product{},
		&models.Category_product{},
	); err != nil {
		return fmt.Errorf("Failed to run migrations: %w", err)
	}

	return nil
}
