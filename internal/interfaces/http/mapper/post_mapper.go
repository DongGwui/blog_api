package mapper

import (
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
	"github.com/ydonggwui/blog-api/internal/interfaces/http/dto"
)

// ToCreatePostCommand converts CreatePostRequest to CreatePostCommand
func ToCreatePostCommand(req *dto.CreatePostRequest) domainService.CreatePostCommand {
	status := entity.PostStatus(req.Status)
	if req.Status == "" {
		status = entity.PostStatusDraft
	}

	return domainService.CreatePostCommand{
		Title:      req.Title,
		Slug:       req.Slug,
		Content:    req.Content,
		Excerpt:    req.Excerpt,
		CategoryID: req.CategoryID,
		TagIDs:     req.TagIDs,
		Status:     status,
		Thumbnail:  req.Thumbnail,
	}
}

// ToUpdatePostCommand converts UpdatePostRequest to UpdatePostCommand
func ToUpdatePostCommand(req *dto.UpdatePostRequest) domainService.UpdatePostCommand {
	return domainService.UpdatePostCommand{
		Title:      req.Title,
		Slug:       req.Slug,
		Content:    req.Content,
		Excerpt:    req.Excerpt,
		CategoryID: req.CategoryID,
		TagIDs:     req.TagIDs,
		Thumbnail:  req.Thumbnail,
	}
}

// ToPostResponse converts PostWithDetails entity to PostResponse DTO
func ToPostResponse(p *entity.PostWithDetails) *dto.PostResponse {
	if p == nil {
		return nil
	}

	return &dto.PostResponse{
		ID:           p.ID,
		Title:        p.Title,
		Slug:         p.Slug,
		Content:      p.Content,
		Excerpt:      p.Excerpt,
		CategoryID:   p.CategoryID,
		CategoryName: p.CategoryName,
		CategorySlug: p.CategorySlug,
		Status:       string(p.Status),
		ViewCount:    p.ViewCount,
		ReadingTime:  p.ReadingTime,
		Thumbnail:    p.Thumbnail,
		Tags:         toTagBriefsInPost(p.Tags),
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
		PublishedAt:  p.PublishedAt,
	}
}

// ToPostListResponse converts PostWithDetails entity to PostListResponse DTO
func ToPostListResponse(p entity.PostWithDetails) dto.PostListResponse {
	return dto.PostListResponse{
		ID:           p.ID,
		Title:        p.Title,
		Slug:         p.Slug,
		Excerpt:      p.Excerpt,
		CategoryID:   p.CategoryID,
		CategoryName: p.CategoryName,
		CategorySlug: p.CategorySlug,
		Status:       string(p.Status),
		ViewCount:    p.ViewCount,
		ReadingTime:  p.ReadingTime,
		Thumbnail:    p.Thumbnail,
		Tags:         toTagBriefsInPost(p.Tags),
		CreatedAt:    p.CreatedAt,
		PublishedAt:  p.PublishedAt,
	}
}

// ToPostListResponses converts a slice of PostWithDetails to PostListResponse DTOs
func ToPostListResponses(posts []entity.PostWithDetails) []dto.PostListResponse {
	result := make([]dto.PostListResponse, len(posts))
	for i, p := range posts {
		result[i] = ToPostListResponse(p)
	}
	return result
}

// toTagBriefsInPost converts entity TagBrief slice to DTO TagBriefInPost slice
func toTagBriefsInPost(tags []entity.TagBrief) []dto.TagBriefInPost {
	result := make([]dto.TagBriefInPost, len(tags))
	for i, t := range tags {
		result[i] = dto.TagBriefInPost{
			ID:   t.ID,
			Name: t.Name,
			Slug: t.Slug,
		}
	}
	return result
}
