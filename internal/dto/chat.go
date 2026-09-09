package dto

import "time"

type ChatResponse struct {
	ID        string    `json:"id"`
	UserOneID string    `json:"user_one_id"`
	UserTwoID string    `json:"user_two_id"`
	CreatedAt time.Time `json:"created_at"`
}
