package service

import (
	"context"
	"fmt"

	"github.com/ydonggwui/blog-api/internal/domain"
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/domain/repository"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
	"github.com/ydonggwui/blog-api/internal/util"
)

type projectService struct {
	projectRepo repository.ProjectRepository
}

func NewProjectService(projectRepo repository.ProjectRepository) domainService.ProjectService {
	return &projectService{projectRepo: projectRepo}
}

// Public API

func (s *projectService) ListProjects(ctx context.Context) ([]entity.Project, error) {
	projects, err := s.projectRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("projectService.ListProjects: %w", err)
	}
	return projects, nil
}

func (s *projectService) ListFeaturedProjects(ctx context.Context) ([]entity.Project, error) {
	projects, err := s.projectRepo.ListFeatured(ctx)
	if err != nil {
		return nil, fmt.Errorf("projectService.ListFeaturedProjects: %w", err)
	}
	return projects, nil
}

func (s *projectService) GetProjectBySlug(ctx context.Context, slug string) (*entity.Project, error) {
	project, err := s.projectRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("projectService.GetProjectBySlug: %w", err)
	}
	return project, nil
}

// Admin API

func (s *projectService) GetProject(ctx context.Context, id int32) (*entity.Project, error) {
	project, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("projectService.GetProject: %w", err)
	}
	return project, nil
}

func (s *projectService) CreateProject(ctx context.Context, cmd domainService.CreateProjectCommand) (*entity.Project, error) {
	// Generate slug
	slug := cmd.Title
	if cmd.Slug != "" {
		slug = cmd.Slug
	}
	slug = util.GenerateSlug(slug)

	// Check if slug exists
	exists, err := s.projectRepo.SlugExists(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("projectService.CreateProject: slug check failed: %w", err)
	}
	if exists {
		return nil, domain.ErrProjectSlugExists
	}

	// Create project entity
	project := &entity.Project{
		Title:       cmd.Title,
		Slug:        slug,
		Description: cmd.Description,
		Content:     cmd.Content,
		TechStack:   cmd.TechStack,
		DemoURL:     cmd.DemoURL,
		GithubURL:   cmd.GithubURL,
		Thumbnail:   cmd.Thumbnail,
		Images:      cmd.Images,
		IsFeatured:  cmd.IsFeatured,
		SortOrder:   cmd.SortOrder,
	}

	created, err := s.projectRepo.Create(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("projectService.CreateProject: create failed: %w", err)
	}
	return created, nil
}

func (s *projectService) UpdateProject(ctx context.Context, id int32, cmd domainService.UpdateProjectCommand) (*entity.Project, error) {
	// Get existing project
	existing, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("projectService.UpdateProject: find project failed: %w", err)
	}

	// Update fields if provided
	if cmd.Title != nil {
		existing.Title = *cmd.Title
	}

	if cmd.Slug != nil && *cmd.Slug != "" {
		newSlug := util.GenerateSlug(*cmd.Slug)
		if newSlug != existing.Slug {
			// Check if new slug exists
			exists, err := s.projectRepo.SlugExistsExcept(ctx, newSlug, id)
			if err != nil {
				return nil, fmt.Errorf("projectService.UpdateProject: slug check failed: %w", err)
			}
			if exists {
				return nil, domain.ErrProjectSlugExists
			}
			existing.Slug = newSlug
		}
	}

	if cmd.Description != nil {
		existing.Description = *cmd.Description
	}
	if cmd.Content != nil {
		existing.Content = *cmd.Content
	}
	if cmd.DemoURL != nil {
		existing.DemoURL = *cmd.DemoURL
	}
	if cmd.GithubURL != nil {
		existing.GithubURL = *cmd.GithubURL
	}
	if cmd.Thumbnail != nil {
		existing.Thumbnail = *cmd.Thumbnail
	}
	if cmd.IsFeatured != nil {
		existing.IsFeatured = *cmd.IsFeatured
	}
	if cmd.SortOrder != nil {
		existing.SortOrder = *cmd.SortOrder
	}

	// Update arrays (replace entirely if provided)
	if cmd.TechStack != nil {
		existing.TechStack = cmd.TechStack
	}
	if cmd.Images != nil {
		existing.Images = cmd.Images
	}

	updated, err := s.projectRepo.Update(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("projectService.UpdateProject: update failed: %w", err)
	}
	return updated, nil
}

func (s *projectService) DeleteProject(ctx context.Context, id int32) error {
	// Check if project exists
	_, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("projectService.DeleteProject: find project failed: %w", err)
	}

	if err := s.projectRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("projectService.DeleteProject: delete failed: %w", err)
	}
	return nil
}

func (s *projectService) ReorderProjects(ctx context.Context, orders []entity.ProjectOrder) error {
	for _, order := range orders {
		if err := s.projectRepo.UpdateOrder(ctx, order.ID, order.SortOrder); err != nil {
			return fmt.Errorf("projectService.ReorderProjects: update order failed for id %d: %w", order.ID, err)
		}
	}
	return nil
}
