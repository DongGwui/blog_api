package model

import "time"

// DashboardStats represents the dashboard statistics response
type DashboardStats struct {
	Posts      PostStats       `json:"posts"`
	Categories []CategoryStats `json:"categories"`
	RecentPosts []RecentPost   `json:"recent_posts"`
}

// PostStats represents post statistics
type PostStats struct {
	Total     int64 `json:"total"`
	Published int64 `json:"published"`
	Draft     int64 `json:"draft"`
}

// CategoryStats represents category with post count
type CategoryStats struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	PostCount int64  `json:"post_count"`
}

// RecentPost represents a recent post for dashboard
type RecentPost struct {
	ID          int32      `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Status      string     `json:"status"`
	ViewCount   int32      `json:"view_count"`
	CreatedAt   time.Time  `json:"created_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}
