package service

import (
	"context"
)

// ViewService defines the interface for view tracking operations
type ViewService interface {
	// RecordView records a view for a post, returns true if it's a new view
	RecordView(ctx context.Context, postID int32, clientIP string) (bool, error)

	// HasViewed checks if the client has already viewed the post
	HasViewed(ctx context.Context, postID int32, clientIP string) (bool, error)
}
