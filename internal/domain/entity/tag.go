package entity

import "time"

// Tag represents a blog post tag
type Tag struct {
	ID        int32
	Name      string
	Slug      string
	PostCount int64
	CreatedAt time.Time
}

// NewTag creates a new Tag entity
func NewTag(name, slug string) *Tag {
	return &Tag{
		Name:      name,
		Slug:      slug,
		CreatedAt: time.Now(),
	}
}
