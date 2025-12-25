package service

import (
	"context"

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
	return s.categoryRepo.FindAll(ctx)
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id int32) (*entity.Category, error) {
	return s.categoryRepo.FindByID(ctx, id)
}

func (s *categoryService) GetCategoryBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	return s.categoryRepo.FindBySlug(ctx, slug)
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
		return nil, err
	}
	if exists {
		return nil, domain.ErrCategorySlugExists
	}

	category := entity.NewCategory(cmd.Name, slug, cmd.Description, cmd.SortOrder)

	created, err := s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	return s.categoryRepo.FindByID(ctx, created.ID)
}

func (s *categoryService) UpdateCategory(ctx context.Context, id int32, cmd domainService.UpdateCategoryCommand) (*entity.Category, error) {
	// Check if category exists
	existing, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
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
			return nil, err
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
		return nil, err
	}

	return s.categoryRepo.FindByID(ctx, id)
}

func (s *categoryService) DeleteCategory(ctx context.Context, id int32) error {
	// Check if category exists
	_, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if category has posts
	postCount, err := s.categoryRepo.GetPostCount(ctx, id)
	if err != nil {
		return err
	}
	if postCount > 0 {
		return domain.ErrCategoryHasPosts
	}

	return s.categoryRepo.Delete(ctx, id)
}
