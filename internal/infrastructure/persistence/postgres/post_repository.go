package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	"github.com/ydonggwui/blog-api/internal/domain"
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/domain/repository"
)

type postRepository struct {
	queries *sqlc.Queries
}

func NewPostRepository(queries *sqlc.Queries) repository.PostRepository {
	return &postRepository{queries: queries}
}

// Basic CRUD

func (r *postRepository) FindByID(ctx context.Context, id int32) (*entity.PostWithDetails, error) {
	post, err := r.queries.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPostNotFound
		}
		return nil, fmt.Errorf("postRepository.FindByID: %w", err)
	}

	tags, err := r.getTags(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("postRepository.FindByID: get tags failed: %w", err)
	}

	return toPostWithDetails(post, tags), nil
}

func (r *postRepository) FindBySlug(ctx context.Context, slug string) (*entity.PostWithDetails, error) {
	post, err := r.queries.GetPostBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPostNotFound
		}
		return nil, fmt.Errorf("postRepository.FindBySlug: %w", err)
	}

	tags, err := r.getTags(ctx, post.ID)
	if err != nil {
		return nil, fmt.Errorf("postRepository.FindBySlug: get tags failed: %w", err)
	}

	// Convert GetPostBySlugRow to GetPostByIDRow format
	postWithDetails := &entity.PostWithDetails{
		Post: *toPostEntity(sqlc.Post{
			ID:          post.ID,
			Title:       post.Title,
			Slug:        post.Slug,
			Content:     post.Content,
			Excerpt:     post.Excerpt,
			CategoryID:  post.CategoryID,
			Status:      post.Status,
			ViewCount:   post.ViewCount,
			ReadingTime: post.ReadingTime,
			Thumbnail:   post.Thumbnail,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			PublishedAt: post.PublishedAt,
		}),
		Tags: tags,
	}

	return postWithDetails, nil
}

func (r *postRepository) Create(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	created, err := r.queries.CreatePost(ctx, toCreatePostParams(post))
	if err != nil {
		return nil, fmt.Errorf("postRepository.Create: %w", err)
	}
	return toPostEntity(created), nil
}

func (r *postRepository) Update(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	updated, err := r.queries.UpdatePost(ctx, toUpdatePostParams(post))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPostNotFound
		}
		return nil, fmt.Errorf("postRepository.Update: %w", err)
	}
	return toPostEntity(updated), nil
}

func (r *postRepository) Delete(ctx context.Context, id int32) error {
	err := r.queries.DeletePost(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrPostNotFound
		}
		return fmt.Errorf("postRepository.Delete: %w", err)
	}
	return nil
}

// Slug validation

func (r *postRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	exists, err := r.queries.CheckSlugExists(ctx, slug)
	if err != nil {
		return false, fmt.Errorf("postRepository.SlugExists: %w", err)
	}
	return exists, nil
}

func (r *postRepository) SlugExistsExcept(ctx context.Context, slug string, excludeID int32) (bool, error) {
	exists, err := r.queries.CheckSlugExistsExcept(ctx, sqlc.CheckSlugExistsExceptParams{
		Slug: slug,
		ID:   excludeID,
	})
	if err != nil {
		return false, fmt.Errorf("postRepository.SlugExistsExcept: %w", err)
	}
	return exists, nil
}

// Published posts (public API)

func (r *postRepository) FindPublishedBySlug(ctx context.Context, slug string) (*entity.PostWithDetails, error) {
	post, err := r.queries.GetPublishedPostBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPostNotFound
		}
		return nil, fmt.Errorf("postRepository.FindPublishedBySlug: %w", err)
	}

	tags, err := r.getTags(ctx, post.ID)
	if err != nil {
		return nil, fmt.Errorf("postRepository.FindPublishedBySlug: get tags failed: %w", err)
	}

	return toPostWithDetailsFromSlug(post, tags), nil
}

func (r *postRepository) ListPublished(ctx context.Context, limit, offset int32) ([]entity.PostWithDetails, error) {
	posts, err := r.queries.ListPublishedPosts(ctx, sqlc.ListPublishedPostsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("postRepository.ListPublished: %w", err)
	}

	result := make([]entity.PostWithDetails, len(posts))
	for i, p := range posts {
		tags, _ := r.getTags(ctx, p.ID)
		result[i] = toPostWithDetailsFromList(p, tags)
	}
	return result, nil
}

func (r *postRepository) CountPublished(ctx context.Context) (int64, error) {
	count, err := r.queries.CountPublishedPosts(ctx)
	if err != nil {
		return 0, fmt.Errorf("postRepository.CountPublished: %w", err)
	}
	return count, nil
}

// Filter by category

func (r *postRepository) ListPublishedByCategory(ctx context.Context, categoryID int32, limit, offset int32) ([]entity.PostWithDetails, error) {
	posts, err := r.queries.ListPublishedPostsByCategory(ctx, sqlc.ListPublishedPostsByCategoryParams{
		CategoryID: sql.NullInt32{Int32: categoryID, Valid: true},
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, fmt.Errorf("postRepository.ListPublishedByCategory: %w", err)
	}

	result := make([]entity.PostWithDetails, len(posts))
	for i, p := range posts {
		tags, _ := r.getTags(ctx, p.ID)
		result[i] = toPostWithDetailsFromCategory(p, tags)
	}
	return result, nil
}

func (r *postRepository) CountPublishedByCategory(ctx context.Context, categoryID int32) (int64, error) {
	count, err := r.queries.CountPublishedPostsByCategory(ctx, sql.NullInt32{Int32: categoryID, Valid: true})
	if err != nil {
		return 0, fmt.Errorf("postRepository.CountPublishedByCategory: %w", err)
	}
	return count, nil
}

// Filter by tag

func (r *postRepository) ListPublishedByTag(ctx context.Context, tagID int32, limit, offset int32) ([]entity.PostWithDetails, error) {
	posts, err := r.queries.ListPublishedPostsByTag(ctx, sqlc.ListPublishedPostsByTagParams{
		TagID:  tagID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("postRepository.ListPublishedByTag: %w", err)
	}

	result := make([]entity.PostWithDetails, len(posts))
	for i, p := range posts {
		tags, _ := r.getTags(ctx, p.ID)
		result[i] = toPostWithDetailsFromTag(p, tags)
	}
	return result, nil
}

func (r *postRepository) CountPublishedByTag(ctx context.Context, tagID int32) (int64, error) {
	count, err := r.queries.CountPublishedPostsByTag(ctx, tagID)
	if err != nil {
		return 0, fmt.Errorf("postRepository.CountPublishedByTag: %w", err)
	}
	return count, nil
}

// Search

func (r *postRepository) SearchPublished(ctx context.Context, query string, limit, offset int32) ([]entity.PostWithDetails, error) {
	queryParam := sql.NullString{String: query, Valid: query != ""}
	posts, err := r.queries.SearchPublishedPosts(ctx, sqlc.SearchPublishedPostsParams{
		Column1: queryParam,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("postRepository.SearchPublished: %w", err)
	}

	result := make([]entity.PostWithDetails, len(posts))
	for i, p := range posts {
		tags, _ := r.getTags(ctx, p.ID)
		result[i] = toPostWithDetailsFromSearch(p, tags)
	}
	return result, nil
}

func (r *postRepository) CountSearchPublished(ctx context.Context, query string) (int64, error) {
	queryParam := sql.NullString{String: query, Valid: query != ""}
	count, err := r.queries.CountSearchPublishedPosts(ctx, queryParam)
	if err != nil {
		return 0, fmt.Errorf("postRepository.CountSearchPublished: %w", err)
	}
	return count, nil
}

// Admin operations

func (r *postRepository) ListAll(ctx context.Context, limit, offset int32) ([]entity.PostWithDetails, error) {
	posts, err := r.queries.ListAllPosts(ctx, sqlc.ListAllPostsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("postRepository.ListAll: %w", err)
	}

	result := make([]entity.PostWithDetails, len(posts))
	for i, p := range posts {
		tags, _ := r.getTags(ctx, p.ID)
		result[i] = toPostWithDetailsFromAll(p, tags)
	}
	return result, nil
}

func (r *postRepository) CountAll(ctx context.Context) (int64, error) {
	count, err := r.queries.CountAllPosts(ctx)
	if err != nil {
		return 0, fmt.Errorf("postRepository.CountAll: %w", err)
	}
	return count, nil
}

func (r *postRepository) ListByStatus(ctx context.Context, status entity.PostStatus, limit, offset int32) ([]entity.PostWithDetails, error) {
	posts, err := r.queries.ListPostsByStatus(ctx, sqlc.ListPostsByStatusParams{
		Status: sql.NullString{String: string(status), Valid: true},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("postRepository.ListByStatus: %w", err)
	}

	result := make([]entity.PostWithDetails, len(posts))
	for i, p := range posts {
		tags, _ := r.getTags(ctx, p.ID)
		result[i] = toPostWithDetailsFromStatus(p, tags)
	}
	return result, nil
}

func (r *postRepository) CountByStatus(ctx context.Context, status entity.PostStatus) (int64, error) {
	count, err := r.queries.CountPostsByStatus(ctx, sql.NullString{String: string(status), Valid: true})
	if err != nil {
		return 0, fmt.Errorf("postRepository.CountByStatus: %w", err)
	}
	return count, nil
}

// Publish/Unpublish

func (r *postRepository) Publish(ctx context.Context, id int32) (*entity.Post, error) {
	post, err := r.queries.PublishPost(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPostNotFound
		}
		return nil, fmt.Errorf("postRepository.Publish: %w", err)
	}
	return toPostEntity(post), nil
}

func (r *postRepository) Unpublish(ctx context.Context, id int32) (*entity.Post, error) {
	post, err := r.queries.UnpublishPost(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPostNotFound
		}
		return nil, fmt.Errorf("postRepository.Unpublish: %w", err)
	}
	return toPostEntity(post), nil
}

// Tag management

func (r *postRepository) GetTags(ctx context.Context, postID int32) ([]entity.TagBrief, error) {
	return r.getTags(ctx, postID)
}

func (r *postRepository) SetTags(ctx context.Context, postID int32, tagIDs []int32) error {
	// Remove existing tags
	if err := r.queries.RemoveAllPostTags(ctx, postID); err != nil {
		return fmt.Errorf("postRepository.SetTags: remove existing tags failed: %w", err)
	}

	// Add new tags
	for _, tagID := range tagIDs {
		if err := r.queries.AddPostTag(ctx, sqlc.AddPostTagParams{
			PostID: postID,
			TagID:  tagID,
		}); err != nil {
			return fmt.Errorf("postRepository.SetTags: add tag %d failed: %w", tagID, err)
		}
	}
	return nil
}

func (r *postRepository) RemoveAllTags(ctx context.Context, postID int32) error {
	if err := r.queries.RemoveAllPostTags(ctx, postID); err != nil {
		return fmt.Errorf("postRepository.RemoveAllTags: %w", err)
	}
	return nil
}

// View count

func (r *postRepository) IncrementViewCount(ctx context.Context, id int32) error {
	if err := r.queries.IncrementViewCount(ctx, id); err != nil {
		return fmt.Errorf("postRepository.IncrementViewCount: %w", err)
	}
	return nil
}

// Helper methods

func (r *postRepository) getTags(ctx context.Context, postID int32) ([]entity.TagBrief, error) {
	tags, err := r.queries.GetPostTags(ctx, postID)
	if err != nil {
		return []entity.TagBrief{}, nil
	}
	return toTagBriefs(tags), nil
}
