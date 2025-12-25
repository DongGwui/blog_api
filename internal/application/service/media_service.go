package service

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/ydonggwui/blog-api/internal/domain"
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/domain/repository"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
)

// Allowed MIME types for upload
var allowedMimeTypes = map[string]bool{
	"image/jpeg":    true,
	"image/png":     true,
	"image/gif":     true,
	"image/webp":    true,
	"image/svg+xml": true,
}

// Maximum file size (10MB)
const maxFileSize = 10 * 1024 * 1024

type mediaService struct {
	mediaRepo   repository.MediaRepository
	storageRepo repository.StorageRepository
}

func NewMediaService(mediaRepo repository.MediaRepository, storageRepo repository.StorageRepository) domainService.MediaService {
	return &mediaService{
		mediaRepo:   mediaRepo,
		storageRepo: storageRepo,
	}
}

func (s *mediaService) ListMedia(ctx context.Context, limit, offset int32) ([]entity.Media, int64, error) {
	media, err := s.mediaRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.mediaRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return media, total, nil
}

func (s *mediaService) GetMediaByID(ctx context.Context, id int32) (*entity.Media, error) {
	return s.mediaRepo.FindByID(ctx, id)
}

func (s *mediaService) UploadMedia(ctx context.Context, cmd domainService.UploadMediaCommand) (*entity.UploadedFile, error) {
	// Validate file type
	if !allowedMimeTypes[cmd.MimeType] {
		return nil, domain.ErrInvalidFileType
	}

	// Validate file size
	if cmd.Size > maxFileSize {
		return nil, domain.ErrFileTooLarge
	}

	// Generate unique filename with UUID
	ext := filepath.Ext(cmd.OriginalName)
	if ext == "" {
		ext = getExtensionFromMimeType(cmd.MimeType)
	}
	filename := uuid.New().String() + ext

	// Generate path: year/month/filename
	now := time.Now()
	path := fmt.Sprintf("%d/%02d/%s", now.Year(), now.Month(), filename)

	// Upload to storage
	err := s.storageRepo.Upload(ctx, path, cmd.File, cmd.Size, cmd.MimeType)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrUploadFailed, err)
	}

	// Generate public URL
	url := s.storageRepo.GenerateURL(path)

	// Save to database
	media := &entity.Media{
		Filename:     filename,
		OriginalName: cmd.OriginalName,
		Path:         path,
		URL:          url,
		MimeType:     cmd.MimeType,
		Size:         cmd.Size,
	}

	created, err := s.mediaRepo.Create(ctx, media)
	if err != nil {
		// Try to delete the uploaded file if database insert fails
		_ = s.storageRepo.Delete(ctx, path)
		return nil, err
	}

	return &entity.UploadedFile{
		ID:           created.ID,
		Filename:     created.Filename,
		OriginalName: created.OriginalName,
		URL:          created.URL,
		MimeType:     created.MimeType,
		Size:         created.Size,
	}, nil
}

func (s *mediaService) DeleteMedia(ctx context.Context, id int32) error {
	// Get media info
	media, err := s.mediaRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete from storage
	err = s.storageRepo.Delete(ctx, media.Path)
	if err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Delete from database
	return s.mediaRepo.Delete(ctx, id)
}

// getExtensionFromMimeType returns the file extension for a MIME type
func getExtensionFromMimeType(mimeType string) string {
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
