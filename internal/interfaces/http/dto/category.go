package dto

import "time"

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	SortOrder   int32  `json:"sort_order,omitempty"`
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	SortOrder   int32  `json:"sort_order,omitempty"`
}

// CategoryResponse represents a category in API responses
type CategoryResponse struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`
	SortOrder   int32     `json:"sort_order"`
	PostCount   *int64    `json:"post_count,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
