package dto

import "time"

// CreateProjectRequest represents the request for creating a project
type CreateProjectRequest struct {
	Title       string   `json:"title" binding:"required,min=1,max=200"`
	Slug        *string  `json:"slug"`
	Description *string  `json:"description"`
	Content     *string  `json:"content"`
	TechStack   []string `json:"tech_stack"`
	DemoURL     *string  `json:"demo_url"`
	GithubURL   *string  `json:"github_url"`
	Thumbnail   *string  `json:"thumbnail"`
	Images      []string `json:"images"`
	IsFeatured  *bool    `json:"is_featured"`
	SortOrder   *int32   `json:"sort_order"`
}

// UpdateProjectRequest represents the request for updating a project
type UpdateProjectRequest struct {
	Title       *string  `json:"title" binding:"omitempty,min=1,max=200"`
	Slug        *string  `json:"slug"`
	Description *string  `json:"description"`
	Content     *string  `json:"content"`
	TechStack   []string `json:"tech_stack"`
	DemoURL     *string  `json:"demo_url"`
	GithubURL   *string  `json:"github_url"`
	Thumbnail   *string  `json:"thumbnail"`
	Images      []string `json:"images"`
	IsFeatured  *bool    `json:"is_featured"`
	SortOrder   *int32   `json:"sort_order"`
}

// ReorderProjectsRequest represents the request for reordering projects
type ReorderProjectsRequest struct {
	Orders []ProjectOrderItem `json:"orders" binding:"required,dive"`
}

// ProjectOrderItem represents a single project order item
type ProjectOrderItem struct {
	ID        int32 `json:"id" binding:"required"`
	SortOrder int32 `json:"sort_order" binding:"required"`
}

// ProjectResponse represents the response for a project
type ProjectResponse struct {
	ID          int32      `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Description string     `json:"description,omitempty"`
	Content     string     `json:"content,omitempty"`
	TechStack   []string   `json:"tech_stack"`
	DemoURL     string     `json:"demo_url,omitempty"`
	GithubURL   string     `json:"github_url,omitempty"`
	Thumbnail   string     `json:"thumbnail,omitempty"`
	Images      []string   `json:"images"`
	IsFeatured  bool       `json:"is_featured"`
	SortOrder   int32      `json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// ProjectListResponse represents the response for a project in list view
type ProjectListResponse struct {
	ID          int32    `json:"id"`
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Description string   `json:"description,omitempty"`
	TechStack   []string `json:"tech_stack"`
	Thumbnail   string   `json:"thumbnail,omitempty"`
	IsFeatured  bool     `json:"is_featured"`
	SortOrder   int32    `json:"sort_order"`
}
