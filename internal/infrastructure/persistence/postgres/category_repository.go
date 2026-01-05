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

type categoryRepository struct {
	queries *sqlc.Queries
}

// NewCategoryRepository creates a new PostgreSQL category repository
func NewCategoryRepository(queries *sqlc.Queries) repository.CategoryRepository {
	return &categoryRepository{
		queries: queries,
	}
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]entity.Category, error) {
	categories, err := r.queries.ListCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("categoryRepository.FindAll: %w", err)
	}

	result := toCategoryEntities(categories)

	// Load post counts for each category
	for i := range result {
		postCount, _ := r.queries.GetCategoryPostCount(ctx, sql.NullInt32{Int32: result[i].ID, Valid: true})
		result[i].PostCount = postCount
	}

	return result, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id int32) (*entity.Category, error) {
	category, err := r.queries.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("categoryRepository.FindByID: %w", err)
	}

	result := toCategoryEntity(category)
	result.PostCount, _ = r.queries.GetCategoryPostCount(ctx, sql.NullInt32{Int32: id, Valid: true})

	return result, nil
}

func (r *categoryRepository) FindBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	category, err := r.queries.GetCategoryBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("categoryRepository.FindBySlug: %w", err)
	}

	result := toCategoryEntity(category)
	result.PostCount, _ = r.queries.GetCategoryPostCount(ctx, sql.NullInt32{Int32: result.ID, Valid: true})

	return result, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	params := toCreateCategoryParams(category)
	created, err := r.queries.CreateCategory(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("categoryRepository.Create: %w", err)
	}

	return toCategoryEntity(created), nil
}

func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	params := toUpdateCategoryParams(category)
	updated, err := r.queries.UpdateCategory(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("categoryRepository.Update: %w", err)
	}

	result := toCategoryEntity(updated)
	result.PostCount, _ = r.queries.GetCategoryPostCount(ctx, sql.NullInt32{Int32: result.ID, Valid: true})

	return result, nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int32) error {
	if err := r.queries.DeleteCategory(ctx, id); err != nil {
		return fmt.Errorf("categoryRepository.Delete: %w", err)
	}
	return nil
}

func (r *categoryRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	_, err := r.queries.GetCategoryBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("categoryRepository.SlugExists: %w", err)
	}
	return true, nil
}

func (r *categoryRepository) SlugExistsExcept(ctx context.Context, slug string, excludeID int32) (bool, error) {
	category, err := r.queries.GetCategoryBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("categoryRepository.SlugExistsExcept: %w", err)
	}
	// Slug exists but it's the same category
	if category.ID == excludeID {
		return false, nil
	}
	return true, nil
}

func (r *categoryRepository) GetPostCount(ctx context.Context, id int32) (int64, error) {
	count, err := r.queries.GetCategoryPostCount(ctx, sql.NullInt32{Int32: id, Valid: true})
	if err != nil {
		return 0, fmt.Errorf("categoryRepository.GetPostCount: %w", err)
	}
	return count, nil
}
