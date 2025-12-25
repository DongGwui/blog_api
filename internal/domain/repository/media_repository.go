package repository

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// MediaRepository defines the interface for media metadata storage
type MediaRepository interface {
	// CRUD operations
	FindByID(ctx context.Context, id int32) (*entity.Media, error)
	Create(ctx context.Context, media *entity.Media) (*entity.Media, error)
	Delete(ctx context.Context, id int32) error

	// List operations
	List(ctx context.Context, limit, offset int32) ([]entity.Media, error)
	Count(ctx context.Context) (int64, error)
}
