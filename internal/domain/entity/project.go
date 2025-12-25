package entity

import "time"

// Project represents a portfolio project entity
type Project struct {
	ID          int32
	Title       string
	Slug        string
	Description string
	Content     string
	TechStack   []string
	DemoURL     string
	GithubURL   string
	Thumbnail   string
	Images      []string
	IsFeatured  bool
	SortOrder   int32
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

// ProjectOrder represents a project reorder item
type ProjectOrder struct {
	ID        int32
	SortOrder int32
}
