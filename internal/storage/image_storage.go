package storage

import "context"

type UploadResult struct {
	ID        string
	UploadURL string
}

type ImageInfo struct {
	ID       string
	Creator  string
	Uploaded bool
}

type ImageStorage interface {
	CreateUploadURL(ctx context.Context, userID string) (*UploadResult, error)
	GetImage(ctx context.Context, imageID string) (*ImageInfo, error)
	DeleteImage(ctx context.Context, imageID string) error
}
