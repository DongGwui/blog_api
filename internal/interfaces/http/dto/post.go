package dto

import "time"

// CreatePostRequest represents the request for creating a post
type CreatePostRequest struct {
	Title      string  `json:"title" binding:"required,max=200"`
	Slug       string  `json:"slug,omitempty"`
	Content    string  `json:"content" binding:"required"`
	Excerpt    string  `json:"excerpt,omitempty" binding:"max=500"`
	CategoryID *int32  `json:"category_id,omitempty"`
	TagIDs     []int32 `json:"tag_ids,omitempty"`
	Status     string  `json:"status,omitempty"`
	Thumbnail  string  `json:"thumbnail,omitempty"`
}

// UpdatePostRequest represents the request for updating a post
type UpdatePostRequest struct {
	Title      string  `json:"title" binding:"required,max=200"`
	Slug       string  `json:"slug,omitempty"`
	Content    string  `json:"content" binding:"required"`
	Excerpt    string  `json:"excerpt,omitempty" binding:"max=500"`
	CategoryID *int32  `json:"category_id,omitempty"`
	TagIDs     []int32 `json:"tag_ids,omitempty"`
	Thumbnail  string  `json:"thumbnail,omitempty"`
}

// PublishRequest represents the request for publishing/unpublishing a post
type PublishRequest struct {
	Publish bool `json:"publish"`
}

// PostResponse represents a post in API responses
type PostResponse struct {
	ID           int32            `json:"id"`
	Title        string           `json:"title"`
	Slug         string           `json:"slug"`
	Content      string           `json:"content,omitempty"`
	Excerpt      string           `json:"excerpt,omitempty"`
	CategoryID   *int32           `json:"category_id,omitempty"`
	CategoryName string           `json:"category_name,omitempty"`
	CategorySlug string           `json:"category_slug,omitempty"`
	Status       string           `json:"status"`
	ViewCount    int32            `json:"view_count"`
	ReadingTime  int32            `json:"reading_time,omitempty"`
	Thumbnail    string           `json:"thumbnail,omitempty"`
	Tags         []TagBriefInPost `json:"tags,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	PublishedAt  *time.Time       `json:"published_at,omitempty"`
}

// PostListResponse represents a post in list responses (without content)
type PostListResponse struct {
	ID           int32            `json:"id"`
	Title        string           `json:"title"`
	Slug         string           `json:"slug"`
	Excerpt      string           `json:"excerpt,omitempty"`
	CategoryID   *int32           `json:"category_id,omitempty"`
	CategoryName string           `json:"category_name,omitempty"`
	CategorySlug string           `json:"category_slug,omitempty"`
	Status       string           `json:"status"`
	ViewCount    int32            `json:"view_count"`
	ReadingTime  int32            `json:"reading_time,omitempty"`
	Thumbnail    string           `json:"thumbnail,omitempty"`
	Tags         []TagBriefInPost `json:"tags,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	PublishedAt  *time.Time       `json:"published_at,omitempty"`
}

// TagBriefInPost represents a tag in post responses
type TagBriefInPost struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
