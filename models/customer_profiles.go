package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CustomerProfile struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;unique" json:"user_id"`
	Addresses datatypes.JSON `gorm:"type:json" json:"addresses"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
