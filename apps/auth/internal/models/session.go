package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID   uuid.UUID `gorm:"type:uuid;index;not null"`
	Agent    string
	IP       string
	IsActive bool `gorm:"default:true"`

	IssuedAt  time.Time `gorm:"autoCreateTime"`
	ExpiresAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User User `gorm:"foreignKey:UserID"`
}
