package dto

import "time"

// DashboardStatsResponse represents the dashboard statistics response
type DashboardStatsResponse struct {
	Posts       PostStatsResponse       `json:"posts"`
	Categories  []CategoryStatsResponse `json:"categories"`
	RecentPosts []RecentPostResponse    `json:"recent_posts"`
}

// PostStatsResponse represents post statistics
type PostStatsResponse struct {
	Total     int64 `json:"total"`
	Published int64 `json:"published"`
	Draft     int64 `json:"draft"`
}

// CategoryStatsResponse represents category with post count
type CategoryStatsResponse struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	PostCount int64  `json:"post_count"`
}

// RecentPostResponse represents a recent post for dashboard
type RecentPostResponse struct {
	ID          int32      `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Status      string     `json:"status"`
	ViewCount   int32      `json:"view_count"`
	CreatedAt   time.Time  `json:"created_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}
