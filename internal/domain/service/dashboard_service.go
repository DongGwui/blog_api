package service

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// DashboardService defines the interface for dashboard operations
type DashboardService interface {
	// GetStats returns dashboard statistics
	GetStats(ctx context.Context) (*entity.DashboardStats, error)
}
