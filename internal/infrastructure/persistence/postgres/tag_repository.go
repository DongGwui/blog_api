package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	"github.com/ydonggwui/blog-api/internal/domain"
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/domain/repository"
)

type tagRepository struct {
	queries *sqlc.Queries
}

// NewTagRepository creates a new PostgreSQL tag repository
func NewTagRepository(queries *sqlc.Queries) repository.TagRepository {
	return &tagRepository{
		queries: queries,
	}
}

func (r *tagRepository) FindAll(ctx context.Context) ([]entity.Tag, error) {
	tags, err := r.queries.ListTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("tagRepository.FindAll: %w", err)
	}

	result := toTagEntities(tags)

	// Load post counts for each tag
	for i := range result {
		postCount, _ := r.queries.GetTagPostCount(ctx, result[i].ID)
		result[i].PostCount = postCount
	}

	return result, nil
}

func (r *tagRepository) FindAllWithPostCount(ctx context.Context) ([]entity.Tag, error) {
	tags, err := r.queries.ListTagsWithPostCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("tagRepository.FindAllWithPostCount: %w", err)
	}

	return toTagEntitiesWithPostCount(tags), nil
}

func (r *tagRepository) FindByID(ctx context.Context, id int32) (*entity.Tag, error) {
	tag, err := r.queries.GetTagByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTagNotFound
		}
		return nil, fmt.Errorf("tagRepository.FindByID: %w", err)
	}

	result := toTagEntity(tag)
	result.PostCount, _ = r.queries.GetTagPostCount(ctx, id)

	return result, nil
}

func (r *tagRepository) FindBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	tag, err := r.queries.GetTagBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTagNotFound
		}
		return nil, fmt.Errorf("tagRepository.FindBySlug: %w", err)
	}

	result := toTagEntity(tag)
	result.PostCount, _ = r.queries.GetTagPostCount(ctx, result.ID)

	return result, nil
}

func (r *tagRepository) Create(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
	params := toCreateTagParams(tag)
	created, err := r.queries.CreateTag(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("tagRepository.Create: %w", err)
	}

	return toTagEntity(created), nil
}

func (r *tagRepository) Update(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
	params := toUpdateTagParams(tag)
	updated, err := r.queries.UpdateTag(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("tagRepository.Update: %w", err)
	}

	result := toTagEntity(updated)
	result.PostCount, _ = r.queries.GetTagPostCount(ctx, result.ID)

	return result, nil
}

func (r *tagRepository) Delete(ctx context.Context, id int32) error {
	if err := r.queries.DeleteTag(ctx, id); err != nil {
		return fmt.Errorf("tagRepository.Delete: %w", err)
	}
	return nil
}

func (r *tagRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	_, err := r.queries.GetTagBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("tagRepository.SlugExists: %w", err)
	}
	return true, nil
}

func (r *tagRepository) SlugExistsExcept(ctx context.Context, slug string, excludeID int32) (bool, error) {
	tag, err := r.queries.GetTagBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("tagRepository.SlugExistsExcept: %w", err)
	}
	// Slug exists but it's the same tag
	if tag.ID == excludeID {
		return false, nil
	}
	return true, nil
}

func (r *tagRepository) GetPostCount(ctx context.Context, id int32) (int64, error) {
	count, err := r.queries.GetTagPostCount(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("tagRepository.GetPostCount: %w", err)
	}
	return count, nil
}
