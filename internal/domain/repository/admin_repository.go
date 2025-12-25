package repository

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// AdminRepository defines the interface for admin storage operations
type AdminRepository interface {
	// FindByID returns an admin by ID
	FindByID(ctx context.Context, id int32) (*entity.Admin, error)

	// FindByUsername returns an admin by username
	FindByUsername(ctx context.Context, username string) (*entity.Admin, error)

	// Create creates a new admin
	Create(ctx context.Context, admin *entity.Admin) (*entity.Admin, error)

	// UpdatePassword updates an admin's password
	UpdatePassword(ctx context.Context, id int32, hashedPassword string) error
}
