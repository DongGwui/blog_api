package service

import (
	"context"
	"io"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// UploadMediaCommand represents the input for uploading a media file
type UploadMediaCommand struct {
	File         io.Reader
	OriginalName string
	MimeType     string
	Size         int64
}

// MediaService defines the interface for media operations
type MediaService interface {
	// ListMedia returns a paginated list of media files
	ListMedia(ctx context.Context, limit, offset int32) ([]entity.Media, int64, error)

	// GetMediaByID returns a media file by ID
	GetMediaByID(ctx context.Context, id int32) (*entity.Media, error)

	// UploadMedia uploads a file and saves metadata
	UploadMedia(ctx context.Context, cmd UploadMediaCommand) (*entity.UploadedFile, error)

	// DeleteMedia removes a media file
	DeleteMedia(ctx context.Context, id int32) error
}
