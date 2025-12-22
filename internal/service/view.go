package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	viewKeyPrefix = "post:view:"
	viewTTL       = 24 * time.Hour
)

type ViewService struct {
	redis       *redis.Client
	postService *PostService
}

func NewViewService(redis *redis.Client, postService *PostService) *ViewService {
	return &ViewService{
		redis:       redis,
		postService: postService,
	}
}

// RecordView records a view for a post, returns true if it's a new view
func (s *ViewService) RecordView(ctx context.Context, postID int32, clientIP string) (bool, error) {
	// Create a unique key for this post + IP combination
	key := s.createViewKey(postID, clientIP)

	// Try to set the key with NX (only if not exists) and TTL
	result, err := s.redis.SetNX(ctx, key, "1", viewTTL).Result()
	if err != nil {
		return false, err
	}

	// If the key was set (new view), increment the view count in DB
	if result {
		if err := s.postService.IncrementViewCount(ctx, postID); err != nil {
			// Log error but don't fail the request
			// The view was already recorded in Redis
			return true, nil
		}
	}

	return result, nil
}

// createViewKey creates a Redis key for view tracking
// Uses SHA256 hash of IP to avoid storing raw IPs
func (s *ViewService) createViewKey(postID int32, clientIP string) string {
	hash := sha256.Sum256([]byte(clientIP))
	ipHash := hex.EncodeToString(hash[:8]) // Use first 8 bytes for shorter key
	return fmt.Sprintf("%s%d:%s", viewKeyPrefix, postID, ipHash)
}

// HasViewed checks if the client has already viewed the post
func (s *ViewService) HasViewed(ctx context.Context, postID int32, clientIP string) (bool, error) {
	key := s.createViewKey(postID, clientIP)
	exists, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
