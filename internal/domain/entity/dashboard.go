package entity

import "time"

// DashboardStats represents the dashboard statistics
type DashboardStats struct {
	Posts       PostStats
	Categories  []CategoryStats
	RecentPosts []RecentPost
}

// PostStats represents post statistics
type PostStats struct {
	Total     int64
	Published int64
	Draft     int64
}

// CategoryStats represents category with post count
type CategoryStats struct {
	ID        int32
	Name      string
	Slug      string
	PostCount int64
}

// RecentPost represents a recent post for dashboard
type RecentPost struct {
	ID          int32
	Title       string
	Slug        string
	Status      string
	ViewCount   int32
	CreatedAt   time.Time
	PublishedAt *time.Time
}
