package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/domain/repository"
)

type dashboardRepository struct {
	queries *sqlc.Queries
}

func NewDashboardRepository(queries *sqlc.Queries) repository.DashboardRepository {
	return &dashboardRepository{queries: queries}
}

func (r *dashboardRepository) GetPostStats(ctx context.Context) (*entity.PostStats, error) {
	total, err := r.queries.CountAllPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("dashboardRepository.GetPostStats: count total failed: %w", err)
	}

	published, err := r.queries.CountPostsByStatus(ctx, sql.NullString{String: "published", Valid: true})
	if err != nil {
		return nil, fmt.Errorf("dashboardRepository.GetPostStats: count published failed: %w", err)
	}

	draft, err := r.queries.CountPostsByStatus(ctx, sql.NullString{String: "draft", Valid: true})
	if err != nil {
		return nil, fmt.Errorf("dashboardRepository.GetPostStats: count draft failed: %w", err)
	}

	return &entity.PostStats{
		Total:     total,
		Published: published,
		Draft:     draft,
	}, nil
}

func (r *dashboardRepository) GetCategoryStats(ctx context.Context) ([]entity.CategoryStats, error) {
	categories, err := r.queries.ListCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("dashboardRepository.GetCategoryStats: %w", err)
	}

	result := make([]entity.CategoryStats, len(categories))
	for i, c := range categories {
		postCount, _ := r.queries.GetCategoryPostCount(ctx, sql.NullInt32{Int32: c.ID, Valid: true})
		result[i] = entity.CategoryStats{
			ID:        c.ID,
			Name:      c.Name,
			Slug:      c.Slug,
			PostCount: postCount,
		}
	}

	return result, nil
}

func (r *dashboardRepository) GetRecentPosts(ctx context.Context, limit int32) ([]entity.RecentPost, error) {
	posts, err := r.queries.ListAllPosts(ctx, sqlc.ListAllPostsParams{
		Limit:  limit,
		Offset: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("dashboardRepository.GetRecentPosts: %w", err)
	}

	result := make([]entity.RecentPost, len(posts))
	for i, p := range posts {
		result[i] = entity.RecentPost{
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
