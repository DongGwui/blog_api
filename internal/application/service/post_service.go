package service

import (
	"context"
	"fmt"

	"github.com/ydonggwui/blog-api/internal/domain"
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/domain/repository"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
	"github.com/ydonggwui/blog-api/internal/util"
)

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) domainService.PostService {
	return &postService{postRepo: postRepo}
}

// Public API

func (s *postService) GetPublishedPost(ctx context.Context, slug string) (*entity.PostWithDetails, error) {
	post, err := s.postRepo.FindPublishedBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("postService.GetPublishedPost: %w", err)
	}
	return post, nil
}

func (s *postService) ListPublishedPosts(ctx context.Context, limit, offset int32) ([]entity.PostWithDetails, int64, error) {
	posts, err := s.postRepo.ListPublished(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("postService.ListPublishedPosts: list failed: %w", err)
	}

	count, err := s.postRepo.CountPublished(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("postService.ListPublishedPosts: count failed: %w", err)
	}

	return posts, count, nil
}

func (s *postService) ListPublishedPostsByCategory(ctx context.Context, categoryID int32, limit, offset int32) ([]entity.PostWithDetails, int64, error) {
	posts, err := s.postRepo.ListPublishedByCategory(ctx, categoryID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("postService.ListPublishedPostsByCategory: list failed: %w", err)
	}

	count, err := s.postRepo.CountPublishedByCategory(ctx, categoryID)
	if err != nil {
		return nil, 0, fmt.Errorf("postService.ListPublishedPostsByCategory: count failed: %w", err)
	}

	return posts, count, nil
}

func (s *postService) ListPublishedPostsByTag(ctx context.Context, tagID int32, limit, offset int32) ([]entity.PostWithDetails, int64, error) {
	posts, err := s.postRepo.ListPublishedByTag(ctx, tagID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("postService.ListPublishedPostsByTag: list failed: %w", err)
	}

	count, err := s.postRepo.CountPublishedByTag(ctx, tagID)
	if err != nil {
		return nil, 0, fmt.Errorf("postService.ListPublishedPostsByTag: count failed: %w", err)
	}

	return posts, count, nil
}

func (s *postService) SearchPublishedPosts(ctx context.Context, query string, limit, offset int32) ([]entity.PostWithDetails, int64, error) {
	posts, err := s.postRepo.SearchPublished(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("postService.SearchPublishedPosts: search failed: %w", err)
	}

	count, err := s.postRepo.CountSearchPublished(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("postService.SearchPublishedPosts: count failed: %w", err)
	}

	return posts, count, nil
}

// Admin API

func (s *postService) GetPost(ctx context.Context, id int32) (*entity.PostWithDetails, error) {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("postService.GetPost: %w", err)
	}
	return post, nil
}

func (s *postService) ListAllPosts(ctx context.Context, limit, offset int32) ([]entity.PostWithDetails, int64, error) {
	posts, err := s.postRepo.ListAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("postService.ListAllPosts: list failed: %w", err)
	}

	count, err := s.postRepo.CountAll(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("postService.ListAllPosts: count failed: %w", err)
	}

	return posts, count, nil
}

func (s *postService) ListPostsByStatus(ctx context.Context, status entity.PostStatus, limit, offset int32) ([]entity.PostWithDetails, int64, error) {
	posts, err := s.postRepo.ListByStatus(ctx, status, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("postService.ListPostsByStatus: list failed: %w", err)
	}

	count, err := s.postRepo.CountByStatus(ctx, status)
	if err != nil {
		return nil, 0, fmt.Errorf("postService.ListPostsByStatus: count failed: %w", err)
	}

	return posts, count, nil
}

func (s *postService) CreatePost(ctx context.Context, cmd domainService.CreatePostCommand) (*entity.PostWithDetails, error) {
	// Generate slug if not provided
	slug := cmd.Slug
	if slug == "" {
		slug = util.GenerateSlug(cmd.Title)
	}

	// Check if slug exists
	exists, err := s.postRepo.SlugExists(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("postService.CreatePost: slug check failed: %w", err)
	}
	if exists {
		return nil, domain.ErrSlugExists
	}

	// Calculate reading time
	readingTime := util.CalculateReadingTime(cmd.Content)

	// Determine status
	status := cmd.Status
	if status == "" {
		status = entity.PostStatusDraft
	}

	// Create post entity
	post := &entity.Post{
		Title:       cmd.Title,
		Slug:        slug,
		Content:     cmd.Content,
		Excerpt:     cmd.Excerpt,
		CategoryID:  cmd.CategoryID,
		Status:      status,
		ReadingTime: int32(readingTime),
		Thumbnail:   cmd.Thumbnail,
	}

	created, err := s.postRepo.Create(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("postService.CreatePost: create failed: %w", err)
	}

	// Add tags
	if len(cmd.TagIDs) > 0 {
		if err := s.postRepo.SetTags(ctx, created.ID, cmd.TagIDs); err != nil {
			return nil, fmt.Errorf("postService.CreatePost: set tags failed: %w", err)
		}
	}

	// Return full post with details
	result, err := s.postRepo.FindByID(ctx, created.ID)
	if err != nil {
		return nil, fmt.Errorf("postService.CreatePost: fetch result failed: %w", err)
	}
	return result, nil
}

func (s *postService) UpdatePost(ctx context.Context, id int32, cmd domainService.UpdatePostCommand) (*entity.PostWithDetails, error) {
	// Check if post exists
	existing, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("postService.UpdatePost: find post failed: %w", err)
	}

	// Generate slug if not provided
	slug := cmd.Slug
	if slug == "" {
		slug = util.GenerateSlug(cmd.Title)
	}

	// Check if slug exists (excluding current post)
	exists, err := s.postRepo.SlugExistsExcept(ctx, slug, id)
	if err != nil {
		return nil, fmt.Errorf("postService.UpdatePost: slug check failed: %w", err)
	}
	if exists {
		return nil, domain.ErrSlugExists
	}

	// Calculate reading time
	readingTime := util.CalculateReadingTime(cmd.Content)

	// Update post entity
	post := &entity.Post{
		ID:          id,
		Title:       cmd.Title,
		Slug:        slug,
		Content:     cmd.Content,
		Excerpt:     cmd.Excerpt,
		CategoryID:  cmd.CategoryID,
		Status:      existing.Status, // Preserve existing status
		ReadingTime: int32(readingTime),
		Thumbnail:   cmd.Thumbnail,
	}

	_, err = s.postRepo.Update(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("postService.UpdatePost: update failed: %w", err)
	}

	// Update tags
	if err := s.postRepo.SetTags(ctx, id, cmd.TagIDs); err != nil {
		return nil, fmt.Errorf("postService.UpdatePost: set tags failed: %w", err)
	}

	// Return full post with details
	result, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("postService.UpdatePost: fetch result failed: %w", err)
	}
	return result, nil
}

func (s *postService) DeletePost(ctx context.Context, id int32) error {
	// Check if post exists
	_, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("postService.DeletePost: find post failed: %w", err)
	}

	// Remove tags first
	if err := s.postRepo.RemoveAllTags(ctx, id); err != nil {
		return fmt.Errorf("postService.DeletePost: remove tags failed: %w", err)
	}

	if err := s.postRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("postService.DeletePost: delete failed: %w", err)
	}
	return nil
}

func (s *postService) PublishPost(ctx context.Context, id int32, publish bool) (*entity.PostWithDetails, error) {
	var err error
	if publish {
		_, err = s.postRepo.Publish(ctx, id)
	} else {
		_, err = s.postRepo.Unpublish(ctx, id)
	}
	if err != nil {
		return nil, fmt.Errorf("postService.PublishPost: status change failed: %w", err)
	}

	result, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("postService.PublishPost: fetch result failed: %w", err)
	}
	return result, nil
}

// View tracking

func (s *postService) IncrementViewCount(ctx context.Context, id int32) error {
	if err := s.postRepo.IncrementViewCount(ctx, id); err != nil {
		return fmt.Errorf("postService.IncrementViewCount: %w", err)
	}
	return nil
}

func (s *postService) GetPostIDBySlug(ctx context.Context, slug string) (int32, error) {
	post, err := s.postRepo.FindBySlug(ctx, slug)
	if err != nil {
		return 0, fmt.Errorf("postService.GetPostIDBySlug: %w", err)
	}
	return post.ID, nil
}
