package repository

import (
	"context"
	"time"
)

// ViewRepository defines the interface for view tracking operations (Redis-based)
type ViewRepository interface {
	// SetViewIfNotExists sets a view record if it doesn't exist, returns true if set (new view)
	SetViewIfNotExists(ctx context.Context, key string, ttl time.Duration) (bool, error)

	// HasView checks if a view record exists
	HasView(ctx context.Context, key string) (bool, error)
}

// ViewCountUpdater defines the interface for updating view counts in the database
type ViewCountUpdater interface {
	// IncrementViewCount increments the view count for a post
	IncrementViewCount(ctx context.Context, postID int32) error
}
