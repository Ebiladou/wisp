package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserOneID uuid.UUID `gorm:"type:uuid;not null;index"`
	UserTwoID uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
