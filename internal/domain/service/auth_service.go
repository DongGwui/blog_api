package service

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// LoginCommand represents the input for login
type LoginCommand struct {
	Username string
	Password string
}

// AuthService defines the interface for authentication operations
type AuthService interface {
	// Login authenticates an admin and returns token info
	Login(ctx context.Context, cmd LoginCommand) (*entity.TokenInfo, error)

	// GetAdminByID returns an admin by ID
	GetAdminByID(ctx context.Context, id int32) (*entity.Admin, error)

	// ValidateToken validates a JWT token and returns claims
	ValidateToken(tokenString string) (*entity.Claims, error)

	// EnsureAdminExists creates the initial admin if it doesn't exist
	EnsureAdminExists(ctx context.Context, username, password string) error
}
