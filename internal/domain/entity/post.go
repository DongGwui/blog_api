package entity

import "time"

// PostStatus represents the status of a post
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
)

// Post represents a blog post entity
type Post struct {
	ID          int32
	Title       string
	Slug        string
	Content     string
	Excerpt     string
	CategoryID  *int32
	Status      PostStatus
	ViewCount   int32
	ReadingTime int32
	Thumbnail   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time
}

// IsPublished returns true if the post is published
func (p *Post) IsPublished() bool {
	return p.Status == PostStatusPublished
}

// IsDraft returns true if the post is a draft
func (p *Post) IsDraft() bool {
	return p.Status == PostStatusDraft
}

// TagBrief represents minimal tag information for post responses
type TagBrief struct {
	ID   int32
	Name string
	Slug string
}

// PostWithDetails represents a post with its category and tags
type PostWithDetails struct {
	Post
	CategoryName string
	CategorySlug string
	Tags         []TagBrief
}
