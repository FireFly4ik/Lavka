package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Market_product struct {
	gorm.Model
	MarketID  uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Stock     int       `gorm:"not null"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
