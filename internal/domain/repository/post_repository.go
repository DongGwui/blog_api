package repository

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// PostRepository defines the interface for post data access
type PostRepository interface {
	// Basic CRUD
	FindByID(ctx context.Context, id int32) (*entity.PostWithDetails, error)
	FindBySlug(ctx context.Context, slug string) (*entity.PostWithDetails, error)
	Create(ctx context.Context, post *entity.Post) (*entity.Post, error)
	Update(ctx context.Context, post *entity.Post) (*entity.Post, error)
	Delete(ctx context.Context, id int32) error

	// Slug validation
	SlugExists(ctx context.Context, slug string) (bool, error)
	SlugExistsExcept(ctx context.Context, slug string, excludeID int32) (bool, error)

	// Published posts (public API)
	FindPublishedBySlug(ctx context.Context, slug string) (*entity.PostWithDetails, error)
	ListPublished(ctx context.Context, limit, offset int32) ([]entity.PostWithDetails, error)
	CountPublished(ctx context.Context) (int64, error)

	// Filter by category
	ListPublishedByCategory(ctx context.Context, categoryID int32, limit, offset int32) ([]entity.PostWithDetails, error)
	CountPublishedByCategory(ctx context.Context, categoryID int32) (int64, error)

	// Filter by tag
	ListPublishedByTag(ctx context.Context, tagID int32, limit, offset int32) ([]entity.PostWithDetails, error)
	CountPublishedByTag(ctx context.Context, tagID int32) (int64, error)

	// Search
	SearchPublished(ctx context.Context, query string, limit, offset int32) ([]entity.PostWithDetails, error)
	CountSearchPublished(ctx context.Context, query string) (int64, error)

	// Admin operations
	ListAll(ctx context.Context, limit, offset int32) ([]entity.PostWithDetails, error)
	CountAll(ctx context.Context) (int64, error)
	ListByStatus(ctx context.Context, status entity.PostStatus, limit, offset int32) ([]entity.PostWithDetails, error)
	CountByStatus(ctx context.Context, status entity.PostStatus) (int64, error)

	// Publish/Unpublish
	Publish(ctx context.Context, id int32) (*entity.Post, error)
	Unpublish(ctx context.Context, id int32) (*entity.Post, error)

	// Tag management
	GetTags(ctx context.Context, postID int32) ([]entity.TagBrief, error)
	SetTags(ctx context.Context, postID int32, tagIDs []int32) error
	RemoveAllTags(ctx context.Context, postID int32) error

	// View count
	IncrementViewCount(ctx context.Context, id int32) error
}
