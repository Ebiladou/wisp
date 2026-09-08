package models

import (
	"time"

	"github.com/google/uuid"
)

type MediaType string

const (
	MediaTypeImage MediaType = "IMAGE"
	MediaTypeVideo MediaType = "VIDEO"
	MediaTypeAudio MediaType = "AUDIO"
)

type Media struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	StorageKey string    `gorm:"type:text;not null;uniqueIndex"`
	Type       MediaType `gorm:"type:varchar(20);not null"`
	MimeType   string    `gorm:"type:varchar(100);not null"`
	Size       int64     `gorm:"not null"`
	CreatedAt  time.Time
}
