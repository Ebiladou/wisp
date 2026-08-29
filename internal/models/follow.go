package models

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	FollowerID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_follower_following"`
	FollowingID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_follower_following"`
	CreatedAt   time.Time
}
