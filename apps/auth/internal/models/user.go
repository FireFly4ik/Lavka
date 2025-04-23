package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	CUSTOMER Role = "customer"
	COURIER  Role = "courier"
	ADMIN    Role = "admin"
)

// User представляет пользователя системы
type User struct {
	gorm.Model
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Avatar          *string   `gorm:"type:varchar(255);"`
	Email           string    `gorm:"uniqueIndex;not null"`
	Phone           *string
	Password        string `gorm:"not null"`
	Username        string `gorm:"uniqueIndex;not null"`
	IsEmailVerified bool   `gorm:"default:false"`
	Role            Role   `gorm:"default:'customer';not null"`

	Resets        []Reset        `gorm:"foreignKey:UserID"`
	Verifications []Verification `gorm:"foreignKey:UserID"`
	Sessions      []Session      `gorm:"foreignKey:UserID"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func formatPhoneNumber(phone string) (string, error) {
	re := regexp.MustCompile(`[^0-9+]`)
	normalizedPhone := re.ReplaceAllString(phone, "")

	if !regexp.MustCompile(`^\+?[1-9]\d{1,14}$`).MatchString(normalizedPhone) {
		return "", fmt.Errorf("invalid phone number format")
	}

	if !strings.HasPrefix(normalizedPhone, "+") {
		normalizedPhone = "+" + normalizedPhone
	}

	return normalizedPhone, nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Phone != nil {
		formattedPhone, err := formatPhoneNumber(*u.Phone)
		if err != nil {
			return err
		}
		u.Phone = &formattedPhone
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if u.Phone != nil {
		formattedPhone, err := formatPhoneNumber(*u.Phone)
		if err != nil {
			return err
		}
		u.Phone = &formattedPhone
	}
	return nil
}
