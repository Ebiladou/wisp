package storage

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/images"
	"github.com/cloudflare/cloudflare-go/v7/option"
)

type CloudflareImages struct {
	client    *cloudflare.Client
	accountID string
}

func NewCloudflareImages(
	accountID string,
	apiToken string,
) *CloudflareImages {
	client := cloudflare.NewClient(
		option.WithAPIToken(apiToken),
	)

	return &CloudflareImages{
		client:    client,
		accountID: accountID,
	}
}

func (storage *CloudflareImages) CreateUploadURL(ctx context.Context, userID string) (*UploadResult, error) {
	directUpload, err := storage.client.Images.V2.DirectUploads.New(
		ctx,
		images.V2DirectUploadNewParams{
			AccountID: cloudflare.F(storage.accountID),
			Creator:   cloudflare.F(userID),
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create cloudflare image upload URL: %w",
			err,
		)
	}

	if directUpload.ID == "" || directUpload.UploadURL == "" {
		return nil, fmt.Errorf(
			"cloudflare returned incomplete upload information",
		)
	}

	return &UploadResult{
		ID:        directUpload.ID,
		UploadURL: directUpload.UploadURL,
	}, nil
}

func (storage *CloudflareImages) GetImage(ctx context.Context, imageID string) (*ImageInfo, error) {
	image, err := storage.client.Images.V1.Get(
		ctx,
		imageID,
		images.V1GetParams{
			AccountID: cloudflare.F(storage.accountID),
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get cloudflare image: %w",
			err,
		)
	}

	return &ImageInfo{
		ID:       image.ID,
		Creator:  image.Creator,
		Uploaded: !image.Uploaded.IsZero(),
	}, nil
}

func (storage *CloudflareImages) DeleteImage(ctx context.Context, imageID string) error {
	_, err := storage.client.Images.V1.Delete(
		ctx,
		imageID,
		images.V1DeleteParams{
			AccountID: cloudflare.F(storage.accountID),
		},
	)
	if err != nil {
		return fmt.Errorf(
			"delete cloudflare image: %w",
			err,
		)
	}

	return nil
}
