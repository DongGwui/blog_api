package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ydonggwui/blog-api/internal/domain/repository"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
)

const (
	viewKeyPrefix = "post:view:"
	viewTTL       = 24 * time.Hour
)

type viewService struct {
	viewRepo         repository.ViewRepository
	viewCountUpdater repository.ViewCountUpdater
}

func NewViewService(viewRepo repository.ViewRepository, viewCountUpdater repository.ViewCountUpdater) domainService.ViewService {
	return &viewService{
		viewRepo:         viewRepo,
		viewCountUpdater: viewCountUpdater,
	}
}

func (s *viewService) RecordView(ctx context.Context, postID int32, clientIP string) (bool, error) {
	// Create a unique key for this post + IP combination
	key := s.createViewKey(postID, clientIP)

	// Try to set the key with NX (only if not exists) and TTL
	isNew, err := s.viewRepo.SetViewIfNotExists(ctx, key, viewTTL)
	if err != nil {
		return false, fmt.Errorf("viewService.RecordView: set view failed: %w", err)
	}

	// If the key was set (new view), increment the view count in DB
	if isNew {
		if err := s.viewCountUpdater.IncrementViewCount(ctx, postID); err != nil {
			// Log error but don't fail the request
			// The view was already recorded in Redis
			return true, nil
		}
	}

	return isNew, nil
}

func (s *viewService) HasViewed(ctx context.Context, postID int32, clientIP string) (bool, error) {
	key := s.createViewKey(postID, clientIP)
	hasViewed, err := s.viewRepo.HasView(ctx, key)
	if err != nil {
		return false, fmt.Errorf("viewService.HasViewed: %w", err)
	}
	return hasViewed, nil
}

// createViewKey creates a Redis key for view tracking
// Uses SHA256 hash of IP to avoid storing raw IPs
func (s *viewService) createViewKey(postID int32, clientIP string) string {
	hash := sha256.Sum256([]byte(clientIP))
	ipHash := hex.EncodeToString(hash[:8]) // Use first 8 bytes for shorter key
	return fmt.Sprintf("%s%d:%s", viewKeyPrefix, postID, ipHash)
}
