package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	"github.com/ydonggwui/blog-api/internal/model"
	"github.com/ydonggwui/blog-api/internal/util"
)

var (
	ErrCategoryNotFound  = errors.New("category not found")
	ErrCategorySlugExists = errors.New("category slug already exists")
	ErrCategoryHasPosts  = errors.New("category has posts")
)

type CategoryService struct {
	queries *sqlc.Queries
}

func NewCategoryService(queries *sqlc.Queries) *CategoryService {
	return &CategoryService{
		queries: queries,
	}
}

// ListCategories returns all categories
func (s *CategoryService) ListCategories(ctx context.Context) ([]model.CategoryResponse, error) {
	categories, err := s.queries.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]model.CategoryResponse, len(categories))
	for i, c := range categories {
		postCount, _ := s.queries.GetCategoryPostCount(ctx, sql.NullInt32{Int32: c.ID, Valid: true})
		result[i] = model.CategoryResponse{
			ID:        c.ID,
			Name:      c.Name,
			Slug:      c.Slug,
			SortOrder: c.SortOrder.Int32,
			PostCount: &postCount,
			CreatedAt: c.CreatedAt.Time,
		}
		if c.Description.Valid {
			result[i].Description = &c.Description.String
		}
	}

	return result, nil
}

// GetCategoryByID returns a category by ID
func (s *CategoryService) GetCategoryByID(ctx context.Context, id int32) (*model.CategoryResponse, error) {
	category, err := s.queries.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	postCount, _ := s.queries.GetCategoryPostCount(ctx, sql.NullInt32{Int32: category.ID, Valid: true})
	resp := &model.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		SortOrder: category.SortOrder.Int32,
		PostCount: &postCount,
		CreatedAt: category.CreatedAt.Time,
	}
	if category.Description.Valid {
		resp.Description = &category.Description.String
	}

	return resp, nil
}

// GetCategoryBySlug returns a category by slug
func (s *CategoryService) GetCategoryBySlug(ctx context.Context, slug string) (*model.CategoryResponse, error) {
	category, err := s.queries.GetCategoryBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	postCount, _ := s.queries.GetCategoryPostCount(ctx, sql.NullInt32{Int32: category.ID, Valid: true})
	resp := &model.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		SortOrder: category.SortOrder.Int32,
		PostCount: &postCount,
		CreatedAt: category.CreatedAt.Time,
	}
	if category.Description.Valid {
		resp.Description = &category.Description.String
	}

	return resp, nil
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(ctx context.Context, req *model.CreateCategoryRequest) (*model.CategoryResponse, error) {
	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = util.GenerateSlug(req.Name)
	}

	// Check if slug exists
	_, err := s.queries.GetCategoryBySlug(ctx, slug)
	if err == nil {
		return nil, ErrCategorySlugExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	category, err := s.queries.CreateCategory(ctx, sqlc.CreateCategoryParams{
		Name:        req.Name,
		Slug:        slug,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		SortOrder:   sql.NullInt32{Int32: req.SortOrder, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return s.GetCategoryByID(ctx, category.ID)
}

// UpdateCategory updates a category
func (s *CategoryService) UpdateCategory(ctx context.Context, id int32, req *model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
	// Check if category exists
	existing, err := s.queries.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = util.GenerateSlug(req.Name)
	}

	// Check if slug exists (excluding current category)
	if slug != existing.Slug {
		_, err := s.queries.GetCategoryBySlug(ctx, slug)
		if err == nil {
			return nil, ErrCategorySlugExists
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	_, err = s.queries.UpdateCategory(ctx, sqlc.UpdateCategoryParams{
		ID:          id,
		Name:        req.Name,
		Slug:        slug,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		SortOrder:   sql.NullInt32{Int32: req.SortOrder, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return s.GetCategoryByID(ctx, id)
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(ctx context.Context, id int32) error {
	// Check if category exists
	_, err := s.queries.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCategoryNotFound
		}
		return err
	}

	// Check if category has posts
	postCount, err := s.queries.GetCategoryPostCount(ctx, sql.NullInt32{Int32: id, Valid: true})
	if err != nil {
		return err
	}
	if postCount > 0 {
		return ErrCategoryHasPosts
	}

	return s.queries.DeleteCategory(ctx, id)
}
