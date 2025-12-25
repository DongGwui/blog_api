package mapper

import (
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/dto"
)

// ToDashboardStatsResponse converts entity.DashboardStats to dto.DashboardStatsResponse
func ToDashboardStatsResponse(s *entity.DashboardStats) dto.DashboardStatsResponse {
	return dto.DashboardStatsResponse{
		Posts:       toPostStatsResponse(s.Posts),
		Categories:  toCategoryStatsResponses(s.Categories),
		RecentPosts: toRecentPostResponses(s.RecentPosts),
	}
}

func toPostStatsResponse(p entity.PostStats) dto.PostStatsResponse {
	return dto.PostStatsResponse{
		Total:     p.Total,
		Published: p.Published,
		Draft:     p.Draft,
	}
}

func toCategoryStatsResponses(categories []entity.CategoryStats) []dto.CategoryStatsResponse {
	result := make([]dto.CategoryStatsResponse, len(categories))
	for i, c := range categories {
		result[i] = dto.CategoryStatsResponse{
			ID:        c.ID,
			Name:      c.Name,
			Slug:      c.Slug,
			PostCount: c.PostCount,
		}
	}
	return result
}

func toRecentPostResponses(posts []entity.RecentPost) []dto.RecentPostResponse {
	result := make([]dto.RecentPostResponse, len(posts))
	for i, p := range posts {
		result[i] = dto.RecentPostResponse{
			ID:          p.ID,
			Title:       p.Title,
			Slug:        p.Slug,
			Status:      p.Status,
			ViewCount:   p.ViewCount,
			CreatedAt:   p.CreatedAt,
			PublishedAt: p.PublishedAt,
		}
	}
	return result
}
