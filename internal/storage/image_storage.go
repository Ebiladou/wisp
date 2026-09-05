package storage

import "context"

type UploadResult struct {
	ID        string
	UploadURL string
}

type ImageStorage interface {
	CreateUploadURL(ctx context.Context, userID string) (*UploadResult, error)
	IsUploaded(ctx context.Context, imageID string) (bool, error)
	DeleteImage(ctx context.Context, imageID string) error
}
