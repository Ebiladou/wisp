package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenType string

const (
	TokenTypeEmailVerification TokenType = "email_verification"
	TokenTypePasswordReset     TokenType = "password_reset"
)

type Token struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"not null;uniqueIndex"`
	TokenType TokenType `gorm:"type:varchar(50);not null"`
	ExpiresAt time.Time `gorm:"not null"`
	UsedAt    *time.Time
	CreatedAt time.Time
}

func (token *Token) BeforeCreate(tx *gorm.DB) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	return nil
}
