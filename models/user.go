package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name" binding:"notnumeric"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Phone     string         `gorm:"type:varchar(15);uniqueIndex" json:"phone"` // New column
	Password  string         `gorm:"not null" json:"-"`
	RoleID    uuid.UUID      `gorm:"type:uuid;not null" json:"role_id"` // FK to roles table
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
