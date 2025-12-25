package service

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/domain/repository"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
)

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
}

func NewDashboardService(dashboardRepo repository.DashboardRepository) domainService.DashboardService {
	return &dashboardService{
		dashboardRepo: dashboardRepo,
	}
}

func (s *dashboardService) GetStats(ctx context.Context) (*entity.DashboardStats, error) {
	// Get post stats
	postStats, err := s.dashboardRepo.GetPostStats(ctx)
	if err != nil {
		return nil, err
	}

	// Get category stats
	categoryStats, err := s.dashboardRepo.GetCategoryStats(ctx)
	if err != nil {
		return nil, err
	}

	// Get recent posts (limit 5)
	recentPosts, err := s.dashboardRepo.GetRecentPosts(ctx, 5)
	if err != nil {
		return nil, err
	}

	return &entity.DashboardStats{
		Posts:       *postStats,
		Categories:  categoryStats,
		RecentPosts: recentPosts,
	}, nil
}
