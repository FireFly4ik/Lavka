package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	Price       float32   `gorm:"no null"`
	Discount    float32   `gorm:"no null" default:"0"`
	Image_url   *string   `gorm:"type:varchar(255);"`

	Market_product   []Market_product   `gorm:"foreignKey:ProductID"`
	Category_product []Category_product `gorm:"foreignKey:ProductID"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
