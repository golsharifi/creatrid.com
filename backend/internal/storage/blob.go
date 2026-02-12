package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
)

type BlobStorage struct {
	client    *azblob.Client
	container string
	baseURL   string
}

func NewBlobStorage(accountName, accountKey, container string) (*BlobStorage, error) {
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("invalid storage credentials: %w", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net", accountName)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob client: %w", err)
	}

	return &BlobStorage{
		client:    client,
		container: container,
		baseURL:   fmt.Sprintf("%s/%s", serviceURL, container),
	}, nil
}

func (b *BlobStorage) Upload(ctx context.Context, blobName string, reader io.Reader, contentType string) (string, error) {
	_, err := b.client.UploadStream(ctx, b.container, blobName, reader, &blockblob.UploadStreamOptions{
		HTTPHeaders: &blob.HTTPHeaders{BlobContentType: &contentType},
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload blob: %w", err)
	}

	return fmt.Sprintf("%s/%s", b.baseURL, blobName), nil
}

func (b *BlobStorage) Delete(ctx context.Context, blobURL string) error {
	blobName := strings.TrimPrefix(blobURL, b.baseURL+"/")
	if blobName == blobURL {
		return nil
	}
	_, err := b.client.DeleteBlob(ctx, b.container, blobName, nil)
	return err
}
