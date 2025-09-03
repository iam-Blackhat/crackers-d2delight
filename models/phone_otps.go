package models

import (
	"time"

	"gorm.io/datatypes"
)

type PhoneOTP struct {
	ID          uint           `gorm:"primaryKey"`
	Phone       string         `gorm:"index;not null"`
	Purpose     string         `gorm:"not null"` // "register", "login"
	CodeHash    string         `gorm:"not null"`
	Meta        datatypes.JSON `gorm:"type:jsonb"`
	Attempts    int
	MaxAttempts int
	ExpiresAt   time.Time
	CreatedAt   time.Time
}
