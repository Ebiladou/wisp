package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserOneID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_chat_users"`
	UserTwoID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_chat_users"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
