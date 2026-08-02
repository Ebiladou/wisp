package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name              string    `gorm:"unique; not null"`
	Username          string    `gorm:"unique; not null"`
	Email             string    `gorm:"unique; not null"`
	Password          string    `gorm:"not null"`
	ProfilePicture    string
	DateOfBirth       *time.Time
	Active            bool
	DeletionRequested bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
