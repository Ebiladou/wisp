package dto

import "time"

type CreateUser struct {
	Name        string     `json:"name" binding:"required"`
	Username    string     `json:"username" binding:"required"`
	Email       string     `json:"email" binding:"required,email"`
	Password    string     `json:"password" binding:"required"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}

type UpdateUserRequest struct {
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

type UserResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profile_picture"`
	Active         bool   `json:"active"`
}

type PublicUserResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}

type ProfilePictureUploadResponse struct {
	ImageID   string `json:"image_id"`
	UploadURL string `json:"upload_url"`
}

type ConfirmProfilePictureRequest struct {
	ImageID string `json:"image_id" binding:"required"`
}
