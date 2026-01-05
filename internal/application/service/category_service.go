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

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repository.CategoryRepository) domainService.CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) ListCategories(ctx context.Context) ([]entity.Category, error) {
	categories, err := s.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("categoryService.ListCategories: %w", err)
	}
	return categories, nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id int32) (*entity.Category, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("categoryService.GetCategoryByID: %w", err)
	}
	return category, nil
}

func (s *categoryService) GetCategoryBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	category, err := s.categoryRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("categoryService.GetCategoryBySlug: %w", err)
	}
	return category, nil
}

func (s *categoryService) CreateCategory(ctx context.Context, cmd domainService.CreateCategoryCommand) (*entity.Category, error) {
	// Generate slug if not provided
	slug := cmd.Slug
	if slug == "" {
		slug = util.GenerateSlug(cmd.Name)
	}

	// Check if slug exists
	exists, err := s.categoryRepo.SlugExists(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("categoryService.CreateCategory: slug check failed: %w", err)
	}
	if exists {
		return nil, domain.ErrCategorySlugExists
	}

	category := entity.NewCategory(cmd.Name, slug, cmd.Description, cmd.SortOrder)

	created, err := s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("categoryService.CreateCategory: create failed: %w", err)
	}

	result, err := s.categoryRepo.FindByID(ctx, created.ID)
	if err != nil {
		return nil, fmt.Errorf("categoryService.CreateCategory: fetch result failed: %w", err)
	}
	return result, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id int32, cmd domainService.UpdateCategoryCommand) (*entity.Category, error) {
	// Check if category exists
	existing, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("categoryService.UpdateCategory: find category failed: %w", err)
	}

	// Generate slug if not provided
	slug := cmd.Slug
	if slug == "" {
		slug = util.GenerateSlug(cmd.Name)
	}

	// Check if slug exists (excluding current category)
	if slug != existing.Slug {
		exists, err := s.categoryRepo.SlugExistsExcept(ctx, slug, id)
		if err != nil {
			return nil, fmt.Errorf("categoryService.UpdateCategory: slug check failed: %w", err)
		}
		if exists {
			return nil, domain.ErrCategorySlugExists
		}
	}

	category := &entity.Category{
		ID:          id,
		Name:        cmd.Name,
		Slug:        slug,
		Description: cmd.Description,
		SortOrder:   cmd.SortOrder,
	}

	_, err = s.categoryRepo.Update(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("categoryService.UpdateCategory: update failed: %w", err)
	}

	result, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("categoryService.UpdateCategory: fetch result failed: %w", err)
	}
	return result, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id int32) error {
	// Check if category exists
	_, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("categoryService.DeleteCategory: find category failed: %w", err)
	}

	// Check if category has posts
	postCount, err := s.categoryRepo.GetPostCount(ctx, id)
	if err != nil {
		return fmt.Errorf("categoryService.DeleteCategory: get post count failed: %w", err)
	}
	if postCount > 0 {
		return domain.ErrCategoryHasPosts
	}

	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("categoryService.DeleteCategory: delete failed: %w", err)
	}
	return nil
}
