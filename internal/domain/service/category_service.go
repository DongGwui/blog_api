package service

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// CreateCategoryCommand represents the command to create a category
type CreateCategoryCommand struct {
	Name        string
	Slug        string
	Description string
	SortOrder   int32
}

// UpdateCategoryCommand represents the command to update a category
type UpdateCategoryCommand struct {
	Name        string
	Slug        string
	Description string
	SortOrder   int32
}

// CategoryService defines the interface for category business logic
type CategoryService interface {
	// ListCategories returns all categories
	ListCategories(ctx context.Context) ([]entity.Category, error)

	// GetCategoryByID returns a category by ID
	GetCategoryByID(ctx context.Context, id int32) (*entity.Category, error)

	// GetCategoryBySlug returns a category by slug
	GetCategoryBySlug(ctx context.Context, slug string) (*entity.Category, error)

	// CreateCategory creates a new category
	CreateCategory(ctx context.Context, cmd CreateCategoryCommand) (*entity.Category, error)

	// UpdateCategory updates an existing category
	UpdateCategory(ctx context.Context, id int32, cmd UpdateCategoryCommand) (*entity.Category, error)

	// DeleteCategory deletes a category
	DeleteCategory(ctx context.Context, id int32) error
}
