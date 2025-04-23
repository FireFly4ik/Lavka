package models

import (
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

// User представляет пользователя системы
type Category_product struct {
	gorm.Model
	CategoryID uuid.UUID `gorm:"type:uuid;not null"`
	ProductID  uuid.UUID `gorm:"type:uuid;not null"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
