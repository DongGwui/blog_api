package dto

import "time"

// CreateTagRequest represents the request body for creating a tag
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,max=50"`
	Slug string `json:"slug,omitempty"`
}

// UpdateTagRequest represents the request body for updating a tag
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required,max=50"`
	Slug string `json:"slug,omitempty"`
}

// TagResponse represents a tag in API responses
type TagResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	PostCount *int64    `json:"post_count,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
