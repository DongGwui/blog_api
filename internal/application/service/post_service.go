package service

import (
	"context"

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
	return s.postRepo.FindPublishedBySlug(ctx, slug)
}

func (s *postService) ListPublishedPosts(ctx context.Context, limit, offset int32) ([]entity.PostWithDetails, int64, error) {
	posts, err := s.postRepo.ListPublished(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.postRepo.CountPublished(ctx)
	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

func (s *postService) ListPublishedPostsByCategory(ctx context.Context, categoryID int32, limit, offset int32) ([]entity.PostWithDetails, int64, error) {
	posts, err := s.postRepo.ListPublishedByCategory(ctx, categoryID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.postRepo.CountPublishedByCategory(ctx, categoryID)
	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

func (s *postService) ListPublishedPostsByTag(ctx context.Context, tagID int32, limit, offset int32) ([]entity.PostWithDetails, int64, error) {
	posts, err := s.postRepo.ListPublishedByTag(ctx, tagID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.postRepo.CountPublishedByTag(ctx, tagID)
	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

func (s *postService) SearchPublishedPosts(ctx context.Context, query string, limit, offset int32) ([]entity.PostWithDetails, int64, error) {
	posts, err := s.postRepo.SearchPublished(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.postRepo.CountSearchPublished(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

// Admin API

func (s *postService) GetPost(ctx context.Context, id int32) (*entity.PostWithDetails, error) {
	return s.postRepo.FindByID(ctx, id)
}

func (s *postService) ListAllPosts(ctx context.Context, limit, offset int32) ([]entity.PostWithDetails, int64, error) {
	posts, err := s.postRepo.ListAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.postRepo.CountAll(ctx)
	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

func (s *postService) ListPostsByStatus(ctx context.Context, status entity.PostStatus, limit, offset int32) ([]entity.PostWithDetails, int64, error) {
	posts, err := s.postRepo.ListByStatus(ctx, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.postRepo.CountByStatus(ctx, status)
	if err != nil {
		return nil, 0, err
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
		return nil, err
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
		return nil, err
	}

	// Add tags
	if len(cmd.TagIDs) > 0 {
		if err := s.postRepo.SetTags(ctx, created.ID, cmd.TagIDs); err != nil {
			return nil, err
		}
	}

	// Return full post with details
	return s.postRepo.FindByID(ctx, created.ID)
}

func (s *postService) UpdatePost(ctx context.Context, id int32, cmd domainService.UpdatePostCommand) (*entity.PostWithDetails, error) {
	// Check if post exists
	existing, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Generate slug if not provided
	slug := cmd.Slug
	if slug == "" {
		slug = util.GenerateSlug(cmd.Title)
	}

	// Check if slug exists (excluding current post)
	exists, err := s.postRepo.SlugExistsExcept(ctx, slug, id)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// Update tags
	if err := s.postRepo.SetTags(ctx, id, cmd.TagIDs); err != nil {
		return nil, err
	}

	// Return full post with details
	return s.postRepo.FindByID(ctx, id)
}

func (s *postService) DeletePost(ctx context.Context, id int32) error {
	// Check if post exists
	_, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Remove tags first
	if err := s.postRepo.RemoveAllTags(ctx, id); err != nil {
		return err
	}

	return s.postRepo.Delete(ctx, id)
}

func (s *postService) PublishPost(ctx context.Context, id int32, publish bool) (*entity.PostWithDetails, error) {
	var err error
	if publish {
		_, err = s.postRepo.Publish(ctx, id)
	} else {
		_, err = s.postRepo.Unpublish(ctx, id)
	}
	if err != nil {
		return nil, err
	}

	return s.postRepo.FindByID(ctx, id)
}

// View tracking

func (s *postService) IncrementViewCount(ctx context.Context, id int32) error {
	return s.postRepo.IncrementViewCount(ctx, id)
}

func (s *postService) GetPostIDBySlug(ctx context.Context, slug string) (int32, error) {
	post, err := s.postRepo.FindBySlug(ctx, slug)
	if err != nil {
		return 0, err
	}
	return post.ID, nil
}
