package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeText  MessageType = "TEXT"
	MessageTypeImage MessageType = "IMAGE"
	MessageTypeVideo MessageType = "VIDEO"
	MessageTypeAudio MessageType = "AUDIO"
)

// every message will first check for an existing chat id, if there is none between the sender and the intended receiver, then a chat will be created and the id assigned predefining the relationship and then the chat is published to the receiver. consequently, receiver id is irrelevant when there is a chat id that determines the relationship.

// access check (block) will sit here too. not sure yet how to verify the access, but I suppose we simply just call the access method after quering chat from chat id, and passing both users id as the argument. issue is this operation might be expensive since every message calls the block policy. not sure (for now).

// users can in the settings define a universal expiry time for their messages (from 1 min to 30 days). and on individual messages to be sent, they can specify the lifespan of that message, overriding the general messages settings (hopefully, this does not mutate the general expiry state).

type Message struct {
	ID        uuid.UUID   `gorm:"type:uuid;primaryKey"`
	ChatID    uuid.UUID   `gorm:"type:uuid;not null;index"`
	SenderID  uuid.UUID   `gorm:"type:uuid;not null;index"`
	Type      MessageType `gorm:"type:varchar(20);not null"`
	Content   *string     `gorm:"type:text"`
	MediaID   *uuid.UUID  `gorm:"type:uuid;index"`
	Media     *Media      `gorm:"foreignKey:MediaID"`
	CreatedAt time.Time
	ExpiresAt time.Time `gorm:"not null;index"`
	Deleted   bool      `gorm:"not null;default:false"`
}
