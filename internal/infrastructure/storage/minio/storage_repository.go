package minio

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/ydonggwui/blog-api/internal/config"
	"github.com/ydonggwui/blog-api/internal/domain/repository"
)

type storageRepository struct {
	client *minio.Client
	cfg    *config.MinIOConfig
}

func NewStorageRepository(client *minio.Client, cfg *config.MinIOConfig) repository.StorageRepository {
	return &storageRepository{
		client: client,
		cfg:    cfg,
	}
}

func (r *storageRepository) Upload(ctx context.Context, path string, file io.Reader, size int64, contentType string) error {
	_, err := r.client.PutObject(ctx, r.cfg.Bucket, path, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (r *storageRepository) Delete(ctx context.Context, path string) error {
	return r.client.RemoveObject(ctx, r.cfg.Bucket, path, minio.RemoveObjectOptions{})
}

func (r *storageRepository) GenerateURL(path string) string {
	publicURL := strings.TrimRight(r.cfg.PublicURL, "/")
	return fmt.Sprintf("%s/%s/%s", publicURL, r.cfg.Bucket, path)
}
