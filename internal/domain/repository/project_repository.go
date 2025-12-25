package repository

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// ProjectRepository defines the interface for project data access
type ProjectRepository interface {
	// Basic CRUD
	FindByID(ctx context.Context, id int32) (*entity.Project, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Project, error)
	Create(ctx context.Context, project *entity.Project) (*entity.Project, error)
	Update(ctx context.Context, project *entity.Project) (*entity.Project, error)
	Delete(ctx context.Context, id int32) error

	// List operations
	ListAll(ctx context.Context) ([]entity.Project, error)
	ListFeatured(ctx context.Context) ([]entity.Project, error)

	// Slug validation
	SlugExists(ctx context.Context, slug string) (bool, error)
	SlugExistsExcept(ctx context.Context, slug string, excludeID int32) (bool, error)

	// Reorder
	UpdateOrder(ctx context.Context, id int32, sortOrder int32) error
}
