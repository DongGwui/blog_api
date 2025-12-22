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
	ErrTagNotFound   = errors.New("tag not found")
	ErrTagSlugExists = errors.New("tag slug already exists")
	ErrTagHasPosts   = errors.New("tag has posts")
)

type TagService struct {
	queries *sqlc.Queries
}

func NewTagService(queries *sqlc.Queries) *TagService {
	return &TagService{
		queries: queries,
	}
}

// ListTags returns all tags
func (s *TagService) ListTags(ctx context.Context) ([]model.TagResponse, error) {
	tags, err := s.queries.ListTags(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]model.TagResponse, len(tags))
	for i, t := range tags {
		postCount, _ := s.queries.GetTagPostCount(ctx, t.ID)
		result[i] = model.TagResponse{
			ID:        t.ID,
			Name:      t.Name,
			Slug:      t.Slug,
			PostCount: &postCount,
			CreatedAt: t.CreatedAt.Time,
		}
	}

	return result, nil
}

// ListTagsWithPostCount returns all tags with post count
func (s *TagService) ListTagsWithPostCount(ctx context.Context) ([]model.TagResponse, error) {
	tags, err := s.queries.ListTagsWithPostCount(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]model.TagResponse, len(tags))
	for i, t := range tags {
		postCount := t.PostCount
		result[i] = model.TagResponse{
			ID:        t.ID,
			Name:      t.Name,
			Slug:      t.Slug,
			PostCount: &postCount,
			CreatedAt: t.CreatedAt.Time,
		}
	}

	return result, nil
}

// GetTagByID returns a tag by ID
func (s *TagService) GetTagByID(ctx context.Context, id int32) (*model.TagResponse, error) {
	tag, err := s.queries.GetTagByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}

	postCount, _ := s.queries.GetTagPostCount(ctx, tag.ID)
	return &model.TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		PostCount: &postCount,
		CreatedAt: tag.CreatedAt.Time,
	}, nil
}

// GetTagBySlug returns a tag by slug
func (s *TagService) GetTagBySlug(ctx context.Context, slug string) (*model.TagResponse, error) {
	tag, err := s.queries.GetTagBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}

	postCount, _ := s.queries.GetTagPostCount(ctx, tag.ID)
	return &model.TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		PostCount: &postCount,
		CreatedAt: tag.CreatedAt.Time,
	}, nil
}

// CreateTag creates a new tag
func (s *TagService) CreateTag(ctx context.Context, req *model.CreateTagRequest) (*model.TagResponse, error) {
	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = util.GenerateSlug(req.Name)
	}

	// Check if slug exists
	_, err := s.queries.GetTagBySlug(ctx, slug)
	if err == nil {
		return nil, ErrTagSlugExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	tag, err := s.queries.CreateTag(ctx, sqlc.CreateTagParams{
		Name: req.Name,
		Slug: slug,
	})
	if err != nil {
		return nil, err
	}

	return s.GetTagByID(ctx, tag.ID)
}

// UpdateTag updates a tag
func (s *TagService) UpdateTag(ctx context.Context, id int32, req *model.UpdateTagRequest) (*model.TagResponse, error) {
	// Check if tag exists
	existing, err := s.queries.GetTagByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}

	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = util.GenerateSlug(req.Name)
	}

	// Check if slug exists (excluding current tag)
	if slug != existing.Slug {
		_, err := s.queries.GetTagBySlug(ctx, slug)
		if err == nil {
			return nil, ErrTagSlugExists
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	_, err = s.queries.UpdateTag(ctx, sqlc.UpdateTagParams{
		ID:   id,
		Name: req.Name,
		Slug: slug,
	})
	if err != nil {
		return nil, err
	}

	return s.GetTagByID(ctx, id)
}

// DeleteTag deletes a tag
func (s *TagService) DeleteTag(ctx context.Context, id int32) error {
	// Check if tag exists
	_, err := s.queries.GetTagByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTagNotFound
		}
		return err
	}

	return s.queries.DeleteTag(ctx, id)
}
