package models

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	FollowerID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_follower_following"`
	FollowingID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_follower_following"`
	CreatedAt   time.Time
}
