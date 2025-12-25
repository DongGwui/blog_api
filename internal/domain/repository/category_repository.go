package repository

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	// FindAll returns all categories
	FindAll(ctx context.Context) ([]entity.Category, error)

	// FindByID returns a category by ID
	FindByID(ctx context.Context, id int32) (*entity.Category, error)

	// FindBySlug returns a category by slug
	FindBySlug(ctx context.Context, slug string) (*entity.Category, error)

	// Create creates a new category
	Create(ctx context.Context, category *entity.Category) (*entity.Category, error)

	// Update updates an existing category
	Update(ctx context.Context, category *entity.Category) (*entity.Category, error)

	// Delete deletes a category by ID
	Delete(ctx context.Context, id int32) error

	// SlugExists checks if a slug already exists
	SlugExists(ctx context.Context, slug string) (bool, error)

	// SlugExistsExcept checks if a slug exists excluding a specific ID
	SlugExistsExcept(ctx context.Context, slug string, excludeID int32) (bool, error)

	// GetPostCount returns the number of posts in a category
	GetPostCount(ctx context.Context, id int32) (int64, error)
}
