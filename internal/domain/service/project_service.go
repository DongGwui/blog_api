package service

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// CreateProjectCommand represents the command for creating a project
type CreateProjectCommand struct {
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
}

// UpdateProjectCommand represents the command for updating a project
type UpdateProjectCommand struct {
	Title       *string
	Slug        *string
	Description *string
	Content     *string
	TechStack   []string
	DemoURL     *string
	GithubURL   *string
	Thumbnail   *string
	Images      []string
	IsFeatured  *bool
	SortOrder   *int32
}

// ProjectService defines the interface for project business logic
type ProjectService interface {
	// Public API
	ListProjects(ctx context.Context) ([]entity.Project, error)
	ListFeaturedProjects(ctx context.Context) ([]entity.Project, error)
	GetProjectBySlug(ctx context.Context, slug string) (*entity.Project, error)

	// Admin API
	GetProject(ctx context.Context, id int32) (*entity.Project, error)
	CreateProject(ctx context.Context, cmd CreateProjectCommand) (*entity.Project, error)
	UpdateProject(ctx context.Context, id int32, cmd UpdateProjectCommand) (*entity.Project, error)
	DeleteProject(ctx context.Context, id int32) error
	ReorderProjects(ctx context.Context, orders []entity.ProjectOrder) error
}
