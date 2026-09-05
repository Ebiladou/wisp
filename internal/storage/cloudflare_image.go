package storage

import (
	"github.com/cloudflare/cloudflare-go/v7"
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
