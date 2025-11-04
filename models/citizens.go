package models

import (
	"time"

	"github.com/google/uuid"
)

type Citizens struct {
	UUID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	NationalID int       `gorm:"uniqueIndex;not null"`
	FirstName  string    `gorm:"not null"`
	LastName   string    `gorm:"not null"`
	Phone      string    `gorm:"not null"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
