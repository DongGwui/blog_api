package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sqlc-dev/pqtype"
	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	"github.com/ydonggwui/blog-api/internal/domain"
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/domain/repository"
)

type projectRepository struct {
	queries *sqlc.Queries
}

func NewProjectRepository(queries *sqlc.Queries) repository.ProjectRepository {
	return &projectRepository{queries: queries}
}

// Basic CRUD

func (r *projectRepository) FindByID(ctx context.Context, id int32) (*entity.Project, error) {
	project, err := r.queries.GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrProjectNotFound
		}
		return nil, fmt.Errorf("projectRepository.FindByID: %w", err)
	}
	return toProjectEntity(&project), nil
}

func (r *projectRepository) FindBySlug(ctx context.Context, slug string) (*entity.Project, error) {
	project, err := r.queries.GetProjectBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrProjectNotFound
		}
		return nil, fmt.Errorf("projectRepository.FindBySlug: %w", err)
	}
	return toProjectEntity(&project), nil
}

func (r *projectRepository) Create(ctx context.Context, project *entity.Project) (*entity.Project, error) {
	params := toCreateProjectParams(project)
	created, err := r.queries.CreateProject(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("projectRepository.Create: %w", err)
	}
	return toProjectEntity(&created), nil
}

func (r *projectRepository) Update(ctx context.Context, project *entity.Project) (*entity.Project, error) {
	params := toUpdateProjectParams(project)
	updated, err := r.queries.UpdateProject(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrProjectNotFound
		}
		return nil, fmt.Errorf("projectRepository.Update: %w", err)
	}
	return toProjectEntity(&updated), nil
}

func (r *projectRepository) Delete(ctx context.Context, id int32) error {
	err := r.queries.DeleteProject(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrProjectNotFound
		}
		return fmt.Errorf("projectRepository.Delete: %w", err)
	}
	return nil
}

// List operations

func (r *projectRepository) ListAll(ctx context.Context) ([]entity.Project, error) {
	projects, err := r.queries.ListProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("projectRepository.ListAll: %w", err)
	}
	return toProjectEntities(projects), nil
}

func (r *projectRepository) ListFeatured(ctx context.Context) ([]entity.Project, error) {
	projects, err := r.queries.ListFeaturedProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("projectRepository.ListFeatured: %w", err)
	}
	return toProjectEntities(projects), nil
}

// Slug validation

func (r *projectRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	_, err := r.queries.GetProjectBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("projectRepository.SlugExists: %w", err)
	}
	return true, nil
}

func (r *projectRepository) SlugExistsExcept(ctx context.Context, slug string, excludeID int32) (bool, error) {
	project, err := r.queries.GetProjectBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("projectRepository.SlugExistsExcept: %w", err)
	}
	return project.ID != excludeID, nil
}

// Reorder

func (r *projectRepository) UpdateOrder(ctx context.Context, id int32, sortOrder int32) error {
	if err := r.queries.UpdateProjectOrder(ctx, sqlc.UpdateProjectOrderParams{
		ID:        id,
		SortOrder: sql.NullInt32{Int32: sortOrder, Valid: true},
	}); err != nil {
		return fmt.Errorf("projectRepository.UpdateOrder: %w", err)
	}
	return nil
}

// Mapper functions

func toProjectEntity(p *sqlc.Project) *entity.Project {
	project := &entity.Project{
		ID:         p.ID,
		Title:      p.Title,
		Slug:       p.Slug,
		TechStack:  []string{},
		Images:     []string{},
		IsFeatured: p.IsFeatured.Bool,
		SortOrder:  p.SortOrder.Int32,
		CreatedAt:  p.CreatedAt.Time,
	}

	if p.Description.Valid {
		project.Description = p.Description.String
	}
	if p.Content.Valid {
		project.Content = p.Content.String
	}
	if p.DemoUrl.Valid {
		project.DemoURL = p.DemoUrl.String
	}
	if p.GithubUrl.Valid {
		project.GithubURL = p.GithubUrl.String
	}
	if p.Thumbnail.Valid {
		project.Thumbnail = p.Thumbnail.String
	}
	if p.UpdatedAt.Valid {
		project.UpdatedAt = &p.UpdatedAt.Time
	}

	// Parse TechStack JSON
	if p.TechStack.Valid && len(p.TechStack.RawMessage) > 0 {
		var techStack []string
		if err := json.Unmarshal(p.TechStack.RawMessage, &techStack); err == nil {
			project.TechStack = techStack
		}
	}

	// Parse Images JSON
	if p.Images.Valid && len(p.Images.RawMessage) > 0 {
		var images []string
		if err := json.Unmarshal(p.Images.RawMessage, &images); err == nil {
			project.Images = images
		}
	}

	return project
}

func toProjectEntities(projects []sqlc.Project) []entity.Project {
	result := make([]entity.Project, len(projects))
	for i, p := range projects {
		result[i] = *toProjectEntity(&p)
	}
	return result
}

func toCreateProjectParams(p *entity.Project) sqlc.CreateProjectParams {
	params := sqlc.CreateProjectParams{
		Title:       p.Title,
		Slug:        p.Slug,
		Description: sql.NullString{String: p.Description, Valid: p.Description != ""},
		Content:     sql.NullString{String: p.Content, Valid: p.Content != ""},
		DemoUrl:     sql.NullString{String: p.DemoURL, Valid: p.DemoURL != ""},
		GithubUrl:   sql.NullString{String: p.GithubURL, Valid: p.GithubURL != ""},
		Thumbnail:   sql.NullString{String: p.Thumbnail, Valid: p.Thumbnail != ""},
		IsFeatured:  sql.NullBool{Bool: p.IsFeatured, Valid: true},
		SortOrder:   sql.NullInt32{Int32: p.SortOrder, Valid: true},
	}

	// Handle TechStack JSON
	if len(p.TechStack) > 0 {
		techStackJSON, err := json.Marshal(p.TechStack)
		if err == nil {
			params.TechStack = pqtype.NullRawMessage{RawMessage: techStackJSON, Valid: true}
		}
	}

	// Handle Images JSON
	if len(p.Images) > 0 {
		imagesJSON, err := json.Marshal(p.Images)
		if err == nil {
			params.Images = pqtype.NullRawMessage{RawMessage: imagesJSON, Valid: true}
		}
	}

	return params
}

func toUpdateProjectParams(p *entity.Project) sqlc.UpdateProjectParams {
	params := sqlc.UpdateProjectParams{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Description: sql.NullString{String: p.Description, Valid: p.Description != ""},
		Content:     sql.NullString{String: p.Content, Valid: p.Content != ""},
		DemoUrl:     sql.NullString{String: p.DemoURL, Valid: p.DemoURL != ""},
		GithubUrl:   sql.NullString{String: p.GithubURL, Valid: p.GithubURL != ""},
		Thumbnail:   sql.NullString{String: p.Thumbnail, Valid: p.Thumbnail != ""},
		IsFeatured:  sql.NullBool{Bool: p.IsFeatured, Valid: true},
		SortOrder:   sql.NullInt32{Int32: p.SortOrder, Valid: true},
	}

	// Handle TechStack JSON
	if len(p.TechStack) > 0 {
		techStackJSON, err := json.Marshal(p.TechStack)
		if err == nil {
			params.TechStack = pqtype.NullRawMessage{RawMessage: techStackJSON, Valid: true}
		}
	} else {
		params.TechStack = pqtype.NullRawMessage{Valid: false}
	}

	// Handle Images JSON
	if len(p.Images) > 0 {
		imagesJSON, err := json.Marshal(p.Images)
		if err == nil {
			params.Images = pqtype.NullRawMessage{RawMessage: imagesJSON, Valid: true}
		}
	} else {
		params.Images = pqtype.NullRawMessage{Valid: false}
	}

	return params
}
