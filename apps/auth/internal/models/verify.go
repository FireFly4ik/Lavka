package models

import (
	"time"

	"github.com/google/uuid"
)

type Verification struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Token string    `gorm:"uniqueIndex;not null"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`
	User   User

	ExpiresAt time.Time
	CreatedAt time.Time
}
