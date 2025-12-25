package repository

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// TagRepository defines the interface for tag data access
type TagRepository interface {
	// FindAll returns all tags
	FindAll(ctx context.Context) ([]entity.Tag, error)

	// FindAllWithPostCount returns all tags with post count (optimized query)
	FindAllWithPostCount(ctx context.Context) ([]entity.Tag, error)

	// FindByID returns a tag by ID
	FindByID(ctx context.Context, id int32) (*entity.Tag, error)

	// FindBySlug returns a tag by slug
	FindBySlug(ctx context.Context, slug string) (*entity.Tag, error)

	// Create creates a new tag
	Create(ctx context.Context, tag *entity.Tag) (*entity.Tag, error)

	// Update updates an existing tag
	Update(ctx context.Context, tag *entity.Tag) (*entity.Tag, error)

	// Delete deletes a tag by ID
	Delete(ctx context.Context, id int32) error

	// SlugExists checks if a slug already exists
	SlugExists(ctx context.Context, slug string) (bool, error)

	// SlugExistsExcept checks if a slug exists excluding a specific ID
	SlugExistsExcept(ctx context.Context, slug string, excludeID int32) (bool, error)

	// GetPostCount returns the number of posts with this tag
	GetPostCount(ctx context.Context, id int32) (int64, error)
}
