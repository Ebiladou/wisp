package models

import (
	"time"

	"github.com/google/uuid"
)

type Block struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	BlockerID uuid.UUID `gorm:"type:uuid;not null;index"`
	BlockedID uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time

	Blocker User `gorm:"foreignKey:BlockerID;constraint:OnDelete:CASCADE"`
	Blocked User `gorm:"foreignKey:BlockedID;constraint:OnDelete:CASCADE"`
}
