package entity

import "time"

// Admin represents an admin user entity
type Admin struct {
	ID        int32
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// TokenInfo represents JWT token information
type TokenInfo struct {
	Token     string
	ExpiresAt time.Time
}

// Claims represents JWT claims
type Claims struct {
	UserID   int32
	Username string
}
