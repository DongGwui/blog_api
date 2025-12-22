package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/sqlc-dev/pqtype"
	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	"github.com/ydonggwui/blog-api/internal/model"
	"github.com/ydonggwui/blog-api/internal/util"
)

var (
	ErrProjectNotFound   = errors.New("project not found")
	ErrProjectSlugExists = errors.New("project slug already exists")
)

type ProjectService struct {
	queries *sqlc.Queries
}

func NewProjectService(queries *sqlc.Queries) *ProjectService {
	return &ProjectService{
		queries: queries,
	}
}

// ListProjects returns all projects ordered by sort_order
func (s *ProjectService) ListProjects(ctx context.Context) ([]model.ProjectListResponse, error) {
	projects, err := s.queries.ListProjects(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]model.ProjectListResponse, len(projects))
	for i, p := range projects {
		result[i] = s.toProjectListResponse(&p)
	}

	return result, nil
}

// ListFeaturedProjects returns featured projects
func (s *ProjectService) ListFeaturedProjects(ctx context.Context) ([]model.ProjectListResponse, error) {
	projects, err := s.queries.ListFeaturedProjects(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]model.ProjectListResponse, len(projects))
	for i, p := range projects {
		result[i] = s.toProjectListResponse(&p)
	}

	return result, nil
}

// GetProjectByID returns a project by ID
func (s *ProjectService) GetProjectByID(ctx context.Context, id int32) (*model.ProjectResponse, error) {
	project, err := s.queries.GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	return s.toProjectResponse(&project), nil
}

// GetProjectBySlug returns a project by slug
func (s *ProjectService) GetProjectBySlug(ctx context.Context, slug string) (*model.ProjectResponse, error) {
	project, err := s.queries.GetProjectBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	return s.toProjectResponse(&project), nil
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, req *model.CreateProjectRequest) (*model.ProjectResponse, error) {
	// Generate slug if not provided
	slug := req.Title
	if req.Slug != nil && *req.Slug != "" {
		slug = *req.Slug
	}
	slug = util.GenerateSlug(slug)

	// Check if slug exists
	_, err := s.queries.GetProjectBySlug(ctx, slug)
	if err == nil {
		return nil, ErrProjectSlugExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Prepare params
	params := sqlc.CreateProjectParams{
		Title: req.Title,
		Slug:  slug,
	}

	if req.Description != nil {
		params.Description = sql.NullString{String: *req.Description, Valid: true}
	}
	if req.Content != nil {
		params.Content = sql.NullString{String: *req.Content, Valid: true}
	}
	if req.DemoURL != nil {
		params.DemoUrl = sql.NullString{String: *req.DemoURL, Valid: true}
	}
	if req.GithubURL != nil {
		params.GithubUrl = sql.NullString{String: *req.GithubURL, Valid: true}
	}
	if req.Thumbnail != nil {
		params.Thumbnail = sql.NullString{String: *req.Thumbnail, Valid: true}
	}
	if req.IsFeatured != nil {
		params.IsFeatured = sql.NullBool{Bool: *req.IsFeatured, Valid: true}
	}
	if req.SortOrder != nil {
		params.SortOrder = sql.NullInt32{Int32: *req.SortOrder, Valid: true}
	}

	// Handle TechStack JSON
	if len(req.TechStack) > 0 {
		techStackJSON, err := json.Marshal(req.TechStack)
		if err != nil {
			return nil, err
		}
		params.TechStack = pqtype.NullRawMessage{RawMessage: techStackJSON, Valid: true}
	}

	// Handle Images JSON
	if len(req.Images) > 0 {
		imagesJSON, err := json.Marshal(req.Images)
		if err != nil {
			return nil, err
		}
		params.Images = pqtype.NullRawMessage{RawMessage: imagesJSON, Valid: true}
	}

	project, err := s.queries.CreateProject(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.toProjectResponse(&project), nil
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(ctx context.Context, id int32, req *model.UpdateProjectRequest) (*model.ProjectResponse, error) {
	// Check if project exists
	existing, err := s.queries.GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	// Prepare update params with existing values
	params := sqlc.UpdateProjectParams{
		ID:          id,
		Title:       existing.Title,
		Slug:        existing.Slug,
		Description: existing.Description,
		Content:     existing.Content,
		TechStack:   existing.TechStack,
		DemoUrl:     existing.DemoUrl,
		GithubUrl:   existing.GithubUrl,
		Thumbnail:   existing.Thumbnail,
		Images:      existing.Images,
		IsFeatured:  existing.IsFeatured,
		SortOrder:   existing.SortOrder,
	}

	// Update fields if provided
	if req.Title != nil {
		params.Title = *req.Title
	}
	if req.Slug != nil && *req.Slug != "" {
		newSlug := util.GenerateSlug(*req.Slug)
		if newSlug != existing.Slug {
			// Check if new slug exists
			_, err := s.queries.GetProjectBySlug(ctx, newSlug)
			if err == nil {
				return nil, ErrProjectSlugExists
			}
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
			params.Slug = newSlug
		}
	}
	if req.Description != nil {
		params.Description = sql.NullString{String: *req.Description, Valid: true}
	}
	if req.Content != nil {
		params.Content = sql.NullString{String: *req.Content, Valid: true}
	}
	if req.DemoURL != nil {
		params.DemoUrl = sql.NullString{String: *req.DemoURL, Valid: true}
	}
	if req.GithubURL != nil {
		params.GithubUrl = sql.NullString{String: *req.GithubURL, Valid: true}
	}
	if req.Thumbnail != nil {
		params.Thumbnail = sql.NullString{String: *req.Thumbnail, Valid: true}
	}
	if req.IsFeatured != nil {
		params.IsFeatured = sql.NullBool{Bool: *req.IsFeatured, Valid: true}
	}
	if req.SortOrder != nil {
		params.SortOrder = sql.NullInt32{Int32: *req.SortOrder, Valid: true}
	}

	// Handle TechStack JSON
	if req.TechStack != nil {
		if len(req.TechStack) > 0 {
			techStackJSON, err := json.Marshal(req.TechStack)
			if err != nil {
				return nil, err
			}
			params.TechStack = pqtype.NullRawMessage{RawMessage: techStackJSON, Valid: true}
		} else {
			params.TechStack = pqtype.NullRawMessage{Valid: false}
		}
	}

	// Handle Images JSON
	if req.Images != nil {
		if len(req.Images) > 0 {
			imagesJSON, err := json.Marshal(req.Images)
			if err != nil {
				return nil, err
			}
			params.Images = pqtype.NullRawMessage{RawMessage: imagesJSON, Valid: true}
		} else {
			params.Images = pqtype.NullRawMessage{Valid: false}
		}
	}

	project, err := s.queries.UpdateProject(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.toProjectResponse(&project), nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(ctx context.Context, id int32) error {
	// Check if project exists
	_, err := s.queries.GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrProjectNotFound
		}
		return err
	}

	return s.queries.DeleteProject(ctx, id)
}

// ReorderProjects updates the sort order of multiple projects
func (s *ProjectService) ReorderProjects(ctx context.Context, orders []model.ProjectOrder) error {
	for _, order := range orders {
		err := s.queries.UpdateProjectOrder(ctx, sqlc.UpdateProjectOrderParams{
			ID:        order.ID,
			SortOrder: sql.NullInt32{Int32: order.SortOrder, Valid: true},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// toProjectResponse converts a sqlc.Project to model.ProjectResponse
func (s *ProjectService) toProjectResponse(p *sqlc.Project) *model.ProjectResponse {
	resp := &model.ProjectResponse{
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
		resp.Description = &p.Description.String
	}
	if p.Content.Valid {
		resp.Content = &p.Content.String
	}
	if p.DemoUrl.Valid {
		resp.DemoURL = &p.DemoUrl.String
	}
	if p.GithubUrl.Valid {
		resp.GithubURL = &p.GithubUrl.String
	}
	if p.Thumbnail.Valid {
		resp.Thumbnail = &p.Thumbnail.String
	}
	if p.UpdatedAt.Valid {
		resp.UpdatedAt = &p.UpdatedAt.Time
	}

	// Parse TechStack JSON
	if p.TechStack.Valid && len(p.TechStack.RawMessage) > 0 {
		var techStack []string
		if err := json.Unmarshal(p.TechStack.RawMessage, &techStack); err == nil {
			resp.TechStack = techStack
		}
	}

	// Parse Images JSON
	if p.Images.Valid && len(p.Images.RawMessage) > 0 {
		var images []string
		if err := json.Unmarshal(p.Images.RawMessage, &images); err == nil {
			resp.Images = images
		}
	}

	return resp
}

// toProjectListResponse converts a sqlc.Project to model.ProjectListResponse
func (s *ProjectService) toProjectListResponse(p *sqlc.Project) model.ProjectListResponse {
	resp := model.ProjectListResponse{
		ID:         p.ID,
		Title:      p.Title,
		Slug:       p.Slug,
		TechStack:  []string{},
		IsFeatured: p.IsFeatured.Bool,
		SortOrder:  p.SortOrder.Int32,
	}

	if p.Description.Valid {
		resp.Description = &p.Description.String
	}
	if p.Thumbnail.Valid {
		resp.Thumbnail = &p.Thumbnail.String
	}

	// Parse TechStack JSON
	if p.TechStack.Valid && len(p.TechStack.RawMessage) > 0 {
		var techStack []string
		if err := json.Unmarshal(p.TechStack.RawMessage, &techStack); err == nil {
			resp.TechStack = techStack
		}
	}

	return resp
}
