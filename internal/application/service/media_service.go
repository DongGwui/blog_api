package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ydonggwui/blog-api/internal/domain"
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/domain/repository"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
	imageutil "github.com/ydonggwui/blog-api/internal/util/image"
)

// Allowed MIME types for upload
var allowedMimeTypes = map[string]bool{
	"image/jpeg":    true,
	"image/png":     true,
	"image/gif":     true,
	"image/webp":    true,
	"image/svg+xml": true,
}

// MIME types that should not be processed (keep original)
var skipProcessingTypes = map[string]bool{
	"image/gif":     true, // Preserve animation
	"image/svg+xml": true, // Vector format
}

// Maximum file size (10MB)
const maxFileSize = 10 * 1024 * 1024

// Image processing settings
const (
	compressionQuality = 85  // WebP quality (0-100)
	thumbnailSmallSize = 150 // Small thumbnail width
	thumbnailMediumSize = 400 // Medium thumbnail width
)

type mediaService struct {
	mediaRepo      repository.MediaRepository
	storageRepo    repository.StorageRepository
	imageProcessor *imageutil.Processor
}

func NewMediaService(mediaRepo repository.MediaRepository, storageRepo repository.StorageRepository) domainService.MediaService {
	return &mediaService{
		mediaRepo:      mediaRepo,
		storageRepo:    storageRepo,
		imageProcessor: imageutil.NewProcessor(compressionQuality),
	}
}

func (s *mediaService) ListMedia(ctx context.Context, limit, offset int32) ([]entity.Media, int64, error) {
	media, err := s.mediaRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("mediaService.ListMedia: list failed: %w", err)
	}

	total, err := s.mediaRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("mediaService.ListMedia: count failed: %w", err)
	}

	return media, total, nil
}

func (s *mediaService) GetMediaByID(ctx context.Context, id int32) (*entity.Media, error) {
	media, err := s.mediaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("mediaService.GetMediaByID: %w", err)
	}
	return media, nil
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

	// Generate unique base filename with UUID
	baseFilename := uuid.New().String()

	// Generate path prefix: year/month/
	now := time.Now()
	pathPrefix := fmt.Sprintf("%d/%02d/", now.Year(), now.Month())

	// Check if we should skip image processing
	if skipProcessingTypes[cmd.MimeType] {
		result, err := s.uploadOriginal(ctx, cmd, baseFilename, pathPrefix)
		if err != nil {
			return nil, fmt.Errorf("mediaService.UploadMedia: %w", err)
		}
		return result, nil
	}

	// Process image (compress and generate thumbnails)
	result, err := s.uploadProcessed(ctx, cmd, baseFilename, pathPrefix)
	if err != nil {
		return nil, fmt.Errorf("mediaService.UploadMedia: %w", err)
	}
	return result, nil
}

// uploadOriginal uploads the file without any processing (for GIF, SVG)
func (s *mediaService) uploadOriginal(ctx context.Context, cmd domainService.UploadMediaCommand, baseFilename, pathPrefix string) (*entity.UploadedFile, error) {
	ext := getExtensionFromMimeType(cmd.MimeType)
	filename := baseFilename + ext
	path := pathPrefix + filename

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
		_ = s.storageRepo.Delete(ctx, path)
		return nil, fmt.Errorf("uploadOriginal: create media record failed: %w", err)
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

// uploadProcessed processes the image (compress to JPEG, generate thumbnails)
func (s *mediaService) uploadProcessed(ctx context.Context, cmd domainService.UploadMediaCommand, baseFilename, pathPrefix string) (*entity.UploadedFile, error) {
	// Read file content into buffer
	fileData, err := io.ReadAll(cmd.File)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read file: %v", domain.ErrUploadFailed, err)
	}

	// Decode image
	img, err := s.imageProcessor.DecodeImage(bytes.NewReader(fileData))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to decode image: %v", domain.ErrUploadFailed, err)
	}

	// Get original dimensions
	width, height := s.imageProcessor.GetDimensions(img)

	// Encode original to JPEG
	jpegData, err := s.imageProcessor.EncodeToJPEG(img)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to encode to jpeg: %v", domain.ErrUploadFailed, err)
	}

	// Generate thumbnails
	thumbnails, err := s.imageProcessor.GenerateThumbnails(img, map[string]int{
		"_sm": thumbnailSmallSize,
		"_md": thumbnailMediumSize,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: failed to generate thumbnails: %v", domain.ErrUploadFailed, err)
	}

	// Prepare file paths
	filename := baseFilename + ".jpg"
	mainPath := pathPrefix + filename
	smPath := pathPrefix + baseFilename + "_sm.jpg"
	mdPath := pathPrefix + baseFilename + "_md.jpg"

	// Upload main image
	err = s.storageRepo.Upload(ctx, mainPath, bytes.NewReader(jpegData), int64(len(jpegData)), "image/jpeg")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to upload main image: %v", domain.ErrUploadFailed, err)
	}

	// Upload thumbnails
	var uploadedPaths []string
	uploadedPaths = append(uploadedPaths, mainPath)

	smData := thumbnails["_sm"].Data
	err = s.storageRepo.Upload(ctx, smPath, bytes.NewReader(smData), int64(len(smData)), "image/jpeg")
	if err != nil {
		s.cleanupFiles(ctx, uploadedPaths)
		return nil, fmt.Errorf("%w: failed to upload small thumbnail: %v", domain.ErrUploadFailed, err)
	}
	uploadedPaths = append(uploadedPaths, smPath)

	mdData := thumbnails["_md"].Data
	err = s.storageRepo.Upload(ctx, mdPath, bytes.NewReader(mdData), int64(len(mdData)), "image/jpeg")
	if err != nil {
		s.cleanupFiles(ctx, uploadedPaths)
		return nil, fmt.Errorf("%w: failed to upload medium thumbnail: %v", domain.ErrUploadFailed, err)
	}
	uploadedPaths = append(uploadedPaths, mdPath)

	// Generate URLs
	mainURL := s.storageRepo.GenerateURL(mainPath)
	smURL := s.storageRepo.GenerateURL(smPath)
	mdURL := s.storageRepo.GenerateURL(mdPath)

	// Save to database
	media := &entity.Media{
		Filename:     filename,
		OriginalName: cmd.OriginalName,
		Path:         mainPath,
		URL:          mainURL,
		MimeType:     "image/jpeg",
		Size:         int64(len(jpegData)),
		Width:        int32(width),
		Height:       int32(height),
		ThumbnailSM:  smURL,
		ThumbnailMD:  mdURL,
	}

	created, err := s.mediaRepo.Create(ctx, media)
	if err != nil {
		s.cleanupFiles(ctx, uploadedPaths)
		return nil, fmt.Errorf("uploadProcessed: create media record failed: %w", err)
	}

	return &entity.UploadedFile{
		ID:           created.ID,
		Filename:     created.Filename,
		OriginalName: created.OriginalName,
		URL:          created.URL,
		MimeType:     created.MimeType,
		Size:         created.Size,
		ThumbnailSM:  created.ThumbnailSM,
		ThumbnailMD:  created.ThumbnailMD,
	}, nil
}

// cleanupFiles deletes uploaded files on error
func (s *mediaService) cleanupFiles(ctx context.Context, paths []string) {
	for _, path := range paths {
		_ = s.storageRepo.Delete(ctx, path)
	}
}

func (s *mediaService) DeleteMedia(ctx context.Context, id int32) error {
	// Get media info
	media, err := s.mediaRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("mediaService.DeleteMedia: find media failed: %w", err)
	}

	// Delete main file from storage
	err = s.storageRepo.Delete(ctx, media.Path)
	if err != nil {
		return fmt.Errorf("mediaService.DeleteMedia: delete file from storage failed: %w", err)
	}

	// Delete thumbnails if they exist
	if media.ThumbnailSM != "" {
		smPath := extractPathFromURL(media.ThumbnailSM)
		if smPath != "" {
			_ = s.storageRepo.Delete(ctx, smPath)
		}
	}
	if media.ThumbnailMD != "" {
		mdPath := extractPathFromURL(media.ThumbnailMD)
		if mdPath != "" {
			_ = s.storageRepo.Delete(ctx, mdPath)
		}
	}

	// Delete from database
	if err := s.mediaRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("mediaService.DeleteMedia: delete record failed: %w", err)
	}
	return nil
}

// extractPathFromURL extracts the storage path from a full URL
// URL format: {publicURL}/{bucket}/{path}
func extractPathFromURL(url string) string {
	// Find the path part after the bucket name
	// Example: http://localhost:9000/blog/2024/01/uuid_sm.webp -> 2024/01/uuid_sm.webp
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return ""
	}
	// Skip protocol, host, and bucket, join the rest
	return strings.Join(parts[4:], "/")
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
