package repository

import (
	"context"
	"io"
)

// StorageRepository defines the interface for file storage operations
type StorageRepository interface {
	// Upload uploads a file and returns the path
	Upload(ctx context.Context, path string, file io.Reader, size int64, contentType string) error

	// Delete removes a file from storage
	Delete(ctx context.Context, path string) error

	// GenerateURL generates a public URL for the file
	GenerateURL(path string) string
}
