package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/ydonggwui/blog-api/internal/config"
	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	"github.com/ydonggwui/blog-api/internal/model"
)

var (
	ErrMediaNotFound     = errors.New("media not found")
	ErrInvalidFileType   = errors.New("invalid file type")
	ErrFileTooLarge      = errors.New("file too large")
	ErrUploadFailed      = errors.New("upload failed")
)

// Allowed MIME types for upload
var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
	"image/svg+xml": true,
}

// Maximum file size (10MB)
const maxFileSize = 10 * 1024 * 1024

type MediaService struct {
	queries   *sqlc.Queries
	minio     *minio.Client
	minioCfg  *config.MinIOConfig
}

func NewMediaService(queries *sqlc.Queries, minioClient *minio.Client, minioCfg *config.MinIOConfig) *MediaService {
	return &MediaService{
		queries:   queries,
		minio:     minioClient,
		minioCfg:  minioCfg,
	}
}

// ListMedia returns a paginated list of media files
func (s *MediaService) ListMedia(ctx context.Context, limit, offset int32) (*model.MediaListResponse, error) {
	media, err := s.queries.ListMedia(ctx, sqlc.ListMediaParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	total, err := s.queries.CountMedia(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]model.MediaResponse, len(media))
	for i, m := range media {
		items[i] = s.toMediaResponse(&m)
	}

	return &model.MediaListResponse{
		Items: items,
		Total: total,
	}, nil
}

// GetMediaByID returns a media file by ID
func (s *MediaService) GetMediaByID(ctx context.Context, id int32) (*model.MediaResponse, error) {
	media, err := s.queries.GetMediaByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMediaNotFound
		}
		return nil, err
	}

	resp := s.toMediaResponse(&media)
	return &resp, nil
}

// UploadFile uploads a file to MinIO and saves metadata to database
func (s *MediaService) UploadFile(ctx context.Context, file io.Reader, originalName string, mimeType string, size int64) (*model.UploadResponse, error) {
	// Validate file type
	if !allowedMimeTypes[mimeType] {
		return nil, ErrInvalidFileType
	}

	// Validate file size
	if size > maxFileSize {
		return nil, ErrFileTooLarge
	}

	// Generate unique filename with UUID
	ext := filepath.Ext(originalName)
	if ext == "" {
		ext = s.getExtensionFromMimeType(mimeType)
	}
	filename := uuid.New().String() + ext

	// Generate path: year/month/filename
	now := time.Now()
	path := fmt.Sprintf("%d/%02d/%s", now.Year(), now.Month(), filename)

	// Upload to MinIO
	_, err := s.minio.PutObject(ctx, s.minioCfg.Bucket, path, file, size, minio.PutObjectOptions{
		ContentType: mimeType,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	// Generate public URL
	url := s.generatePublicURL(path)

	// Save to database
	media, err := s.queries.CreateMedia(ctx, sqlc.CreateMediaParams{
		Filename:     filename,
		OriginalName: originalName,
		Path:         path,
		Url:          url,
		MimeType:     sql.NullString{String: mimeType, Valid: true},
		Size:         sql.NullInt64{Int64: size, Valid: true},
	})
	if err != nil {
		// Try to delete the uploaded file if database insert fails
		_ = s.minio.RemoveObject(ctx, s.minioCfg.Bucket, path, minio.RemoveObjectOptions{})
		return nil, err
	}

	return &model.UploadResponse{
		ID:           media.ID,
		Filename:     media.Filename,
		OriginalName: media.OriginalName,
		URL:          media.Url,
		MimeType:     stringPtr(mimeType),
		Size:         int64Ptr(size),
	}, nil
}

// DeleteMedia deletes a media file from MinIO and database
func (s *MediaService) DeleteMedia(ctx context.Context, id int32) error {
	// Get media info
	media, err := s.queries.GetMediaByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrMediaNotFound
		}
		return err
	}

	// Delete from MinIO
	err = s.minio.RemoveObject(ctx, s.minioCfg.Bucket, media.Path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Delete from database
	return s.queries.DeleteMedia(ctx, id)
}

// generatePublicURL generates a public URL for the file
func (s *MediaService) generatePublicURL(path string) string {
	publicURL := strings.TrimRight(s.minioCfg.PublicURL, "/")
	return fmt.Sprintf("%s/%s/%s", publicURL, s.minioCfg.Bucket, path)
}

// getExtensionFromMimeType returns the file extension for a MIME type
func (s *MediaService) getExtensionFromMimeType(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	default:
		return ""
	}
}

// toMediaResponse converts a sqlc.Medium to model.MediaResponse
func (s *MediaService) toMediaResponse(m *sqlc.Medium) model.MediaResponse {
	resp := model.MediaResponse{
		ID:           m.ID,
		Filename:     m.Filename,
		OriginalName: m.OriginalName,
		Path:         m.Path,
		URL:          m.Url,
		CreatedAt:    m.CreatedAt.Time,
	}

	if m.MimeType.Valid {
		resp.MimeType = &m.MimeType.String
	}
	if m.Size.Valid {
		resp.Size = &m.Size.Int64
	}
	if m.Width.Valid {
		resp.Width = &m.Width.Int32
	}
	if m.Height.Valid {
		resp.Height = &m.Height.Int32
	}

	return resp
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}
