package service

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// CreatePostCommand represents the command for creating a post
type CreatePostCommand struct {
	Title      string
	Slug       string
	Content    string
	Excerpt    string
	CategoryID *int32
	TagIDs     []int32
	Status     entity.PostStatus
	Thumbnail  string
}

// UpdatePostCommand represents the command for updating a post
type UpdatePostCommand struct {
	Title      string
	Slug       string
	Content    string
	Excerpt    string
	CategoryID *int32
	TagIDs     []int32
	Thumbnail  string
}

// PostService defines the interface for post business logic
type PostService interface {
	// Public API
	GetPublishedPost(ctx context.Context, slug string) (*entity.PostWithDetails, error)
	ListPublishedPosts(ctx context.Context, limit, offset int32) ([]entity.PostWithDetails, int64, error)
	ListPublishedPostsByCategory(ctx context.Context, categoryID int32, limit, offset int32) ([]entity.PostWithDetails, int64, error)
	ListPublishedPostsByTag(ctx context.Context, tagID int32, limit, offset int32) ([]entity.PostWithDetails, int64, error)
	SearchPublishedPosts(ctx context.Context, query string, limit, offset int32) ([]entity.PostWithDetails, int64, error)

	// Admin API
	GetPost(ctx context.Context, id int32) (*entity.PostWithDetails, error)
	ListAllPosts(ctx context.Context, limit, offset int32) ([]entity.PostWithDetails, int64, error)
	ListPostsByStatus(ctx context.Context, status entity.PostStatus, limit, offset int32) ([]entity.PostWithDetails, int64, error)
	CreatePost(ctx context.Context, cmd CreatePostCommand) (*entity.PostWithDetails, error)
	UpdatePost(ctx context.Context, id int32, cmd UpdatePostCommand) (*entity.PostWithDetails, error)
	DeletePost(ctx context.Context, id int32) error
	PublishPost(ctx context.Context, id int32, publish bool) (*entity.PostWithDetails, error)

	// View tracking
	IncrementViewCount(ctx context.Context, id int32) error
	GetPostIDBySlug(ctx context.Context, slug string) (int32, error)
}
