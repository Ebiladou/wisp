package models

import (
	"time"

	"github.com/google/uuid"
)

type ViewingSession struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	MessageID  uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	OpenedAt   time.Time `gorm:"not null"`
	ExpiresAt  time.Time `gorm:"not null;index"`
	ConsumedAt *time.Time
}
