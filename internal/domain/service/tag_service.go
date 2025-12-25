package service

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// CreateTagCommand represents the command to create a tag
type CreateTagCommand struct {
	Name string
	Slug string
}

// UpdateTagCommand represents the command to update a tag
type UpdateTagCommand struct {
	Name string
	Slug string
}

// TagService defines the interface for tag business logic
type TagService interface {
	// ListTags returns all tags
	ListTags(ctx context.Context) ([]entity.Tag, error)

	// ListTagsWithPostCount returns all tags with post count (optimized)
	ListTagsWithPostCount(ctx context.Context) ([]entity.Tag, error)

	// GetTagByID returns a tag by ID
	GetTagByID(ctx context.Context, id int32) (*entity.Tag, error)

	// GetTagBySlug returns a tag by slug
	GetTagBySlug(ctx context.Context, slug string) (*entity.Tag, error)

	// CreateTag creates a new tag
	CreateTag(ctx context.Context, cmd CreateTagCommand) (*entity.Tag, error)

	// UpdateTag updates an existing tag
	UpdateTag(ctx context.Context, id int32, cmd UpdateTagCommand) (*entity.Tag, error)

	// DeleteTag deletes a tag
	DeleteTag(ctx context.Context, id int32) error
}
