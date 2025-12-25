package repository

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// DashboardRepository defines the interface for dashboard statistics operations
type DashboardRepository interface {
	// GetPostStats returns post statistics (total, published, draft counts)
	GetPostStats(ctx context.Context) (*entity.PostStats, error)

	// GetCategoryStats returns all categories with their post counts
	GetCategoryStats(ctx context.Context) ([]entity.CategoryStats, error)

	// GetRecentPosts returns the most recent posts
	GetRecentPosts(ctx context.Context, limit int32) ([]entity.RecentPost, error)
}
