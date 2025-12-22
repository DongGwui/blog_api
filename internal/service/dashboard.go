package service

import (
	"context"
	"database/sql"

	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	"github.com/ydonggwui/blog-api/internal/model"
)

type DashboardService struct {
	queries *sqlc.Queries
}

func NewDashboardService(queries *sqlc.Queries) *DashboardService {
	return &DashboardService{
		queries: queries,
	}
}

// GetStats returns dashboard statistics
func (s *DashboardService) GetStats(ctx context.Context) (*model.DashboardStats, error) {
	// Get post stats
	postStats, err := s.getPostStats(ctx)
	if err != nil {
		return nil, err
	}

	// Get category stats
	categoryStats, err := s.getCategoryStats(ctx)
	if err != nil {
		return nil, err
	}

	// Get recent posts
	recentPosts, err := s.getRecentPosts(ctx, 5)
	if err != nil {
		return nil, err
	}

	return &model.DashboardStats{
		Posts:       *postStats,
		Categories:  categoryStats,
		RecentPosts: recentPosts,
	}, nil
}

// getPostStats returns post statistics
func (s *DashboardService) getPostStats(ctx context.Context) (*model.PostStats, error) {
	total, err := s.queries.CountAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	published, err := s.queries.CountPostsByStatus(ctx, sql.NullString{String: "published", Valid: true})
	if err != nil {
		return nil, err
	}

	draft, err := s.queries.CountPostsByStatus(ctx, sql.NullString{String: "draft", Valid: true})
	if err != nil {
		return nil, err
	}

	return &model.PostStats{
		Total:     total,
		Published: published,
		Draft:     draft,
	}, nil
}

// getCategoryStats returns category statistics with post counts
func (s *DashboardService) getCategoryStats(ctx context.Context) ([]model.CategoryStats, error) {
	categories, err := s.queries.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]model.CategoryStats, len(categories))
	for i, c := range categories {
		postCount, _ := s.queries.GetCategoryPostCount(ctx, sql.NullInt32{Int32: c.ID, Valid: true})
		result[i] = model.CategoryStats{
			ID:        c.ID,
			Name:      c.Name,
			Slug:      c.Slug,
			PostCount: postCount,
		}
	}

	return result, nil
}

// getRecentPosts returns recent posts
func (s *DashboardService) getRecentPosts(ctx context.Context, limit int32) ([]model.RecentPost, error) {
	posts, err := s.queries.ListAllPosts(ctx, sqlc.ListAllPostsParams{
		Limit:  limit,
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.RecentPost, len(posts))
	for i, p := range posts {
		result[i] = model.RecentPost{
			ID:        p.ID,
			Title:     p.Title,
			Slug:      p.Slug,
			Status:    p.Status.String,
			ViewCount: p.ViewCount.Int32,
			CreatedAt: p.CreatedAt.Time,
		}
		if p.PublishedAt.Valid {
			result[i].PublishedAt = &p.PublishedAt.Time
		}
	}

	return result, nil
}
