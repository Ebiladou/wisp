package dto

import (
	"time"

	"github.com/Ebiladou/wisp/internal/models"
	"github.com/google/uuid"
)

type MessageCreate struct {
	RecipientID uuid.UUID          `json:"recipient_id" binding:"required"`
	Type        models.MessageType `json:"type" binding:"required"`
	Content     *string            `json:"content"`
	MediaID     *uuid.UUID         `json:"media_id"`
}

type MessageResponse struct {
	ID          uuid.UUID          `json:"id"`
	ChatID      uuid.UUID          `json:"chat_id"`
	SenderID    uuid.UUID          `json:"sender_id"`
	RecipientID uuid.UUID          `json:"recipient_id"`
	Type        models.MessageType `json:"type"`
	Content     *string            `json:"content"`
	MediaID     *uuid.UUID         `json:"media_id"`
	CreatedAt   time.Time          `json:"created_at"`
	ExpiresAt   time.Time          `json:"expires_at"`
}
