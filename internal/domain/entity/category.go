package entity

import "time"

// Category represents a blog post category
type Category struct {
	ID          int32
	Name        string
	Slug        string
	Description string
	SortOrder   int32
	PostCount   int64
	CreatedAt   time.Time
}

// NewCategory creates a new Category entity
func NewCategory(name, slug, description string, sortOrder int32) *Category {
	return &Category{
		Name:        name,
		Slug:        slug,
		Description: description,
		SortOrder:   sortOrder,
		CreatedAt:   time.Now(),
	}
}

// HasDescription returns true if the category has a description
func (c *Category) HasDescription() bool {
	return c.Description != ""
}
