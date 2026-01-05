package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ydonggwui/blog-api/internal/domain"
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/domain/repository/mocks"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
)

func TestCategoryService_ListCategories(t *testing.T) {
	mockRepo := &mocks.MockCategoryRepository{
		FindAllFunc: func(ctx context.Context) ([]entity.Category, error) {
			return []entity.Category{
				{ID: 1, Name: "Tech", Slug: "tech", CreatedAt: time.Now()},
				{ID: 2, Name: "Life", Slug: "life", CreatedAt: time.Now()},
			}, nil
		},
	}

	svc := NewCategoryService(mockRepo)
	categories, err := svc.ListCategories(context.Background())

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(categories))
	}
}

func TestCategoryService_GetCategoryByID(t *testing.T) {
	expected := &entity.Category{
		ID:        1,
		Name:      "Tech",
		Slug:      "tech",
		CreatedAt: time.Now(),
	}

	mockRepo := &mocks.MockCategoryRepository{
		FindByIDFunc: func(ctx context.Context, id int32) (*entity.Category, error) {
			if id == 1 {
				return expected, nil
			}
			return nil, domain.ErrCategoryNotFound
		},
	}

	svc := NewCategoryService(mockRepo)

	t.Run("existing category", func(t *testing.T) {
		category, err := svc.GetCategoryByID(context.Background(), 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if category.ID != expected.ID {
			t.Errorf("expected ID %d, got %d", expected.ID, category.ID)
		}
	})

	t.Run("non-existing category", func(t *testing.T) {
		_, err := svc.GetCategoryByID(context.Background(), 999)
		if !errors.Is(err, domain.ErrCategoryNotFound) {
			t.Errorf("expected ErrCategoryNotFound, got %v", err)
		}
	})
}

func TestCategoryService_CreateCategory(t *testing.T) {
	createdCategory := &entity.Category{
		ID:        1,
		Name:      "Tech",
		Slug:      "tech",
		CreatedAt: time.Now(),
	}

	mockRepo := &mocks.MockCategoryRepository{
		SlugExistsFunc: func(ctx context.Context, slug string) (bool, error) {
			return false, nil
		},
		CreateFunc: func(ctx context.Context, category *entity.Category) (*entity.Category, error) {
			category.ID = 1
			return category, nil
		},
		FindByIDFunc: func(ctx context.Context, id int32) (*entity.Category, error) {
			return createdCategory, nil
		},
	}

	svc := NewCategoryService(mockRepo)

	t.Run("successful creation", func(t *testing.T) {
		cmd := domainService.CreateCategoryCommand{
			Name: "Tech",
			Slug: "tech",
		}
		category, err := svc.CreateCategory(context.Background(), cmd)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if category.Name != "Tech" {
			t.Errorf("expected name Tech, got %s", category.Name)
		}
	})

	t.Run("slug already exists", func(t *testing.T) {
		mockRepo.SlugExistsFunc = func(ctx context.Context, slug string) (bool, error) {
			return true, nil
		}

		cmd := domainService.CreateCategoryCommand{
			Name: "Tech",
			Slug: "tech",
		}
		_, err := svc.CreateCategory(context.Background(), cmd)
		if err != domain.ErrCategorySlugExists {
			t.Errorf("expected ErrCategorySlugExists, got %v", err)
		}
	})
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	existingCategory := &entity.Category{
		ID:        1,
		Name:      "Tech",
		Slug:      "tech",
		CreatedAt: time.Now(),
	}

	mockRepo := &mocks.MockCategoryRepository{
		FindByIDFunc: func(ctx context.Context, id int32) (*entity.Category, error) {
			if id == 1 {
				return existingCategory, nil
			}
			return nil, domain.ErrCategoryNotFound
		},
		SlugExistsExceptFunc: func(ctx context.Context, slug string, excludeID int32) (bool, error) {
			return false, nil
		},
		UpdateFunc: func(ctx context.Context, category *entity.Category) (*entity.Category, error) {
			return category, nil
		},
	}

	svc := NewCategoryService(mockRepo)

	t.Run("successful update", func(t *testing.T) {
		cmd := domainService.UpdateCategoryCommand{
			Name: "Technology",
			Slug: "technology",
		}
		category, err := svc.UpdateCategory(context.Background(), 1, cmd)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if category == nil {
			t.Error("expected category, got nil")
		}
	})

	t.Run("category not found", func(t *testing.T) {
		cmd := domainService.UpdateCategoryCommand{
			Name: "Technology",
			Slug: "technology",
		}
		_, err := svc.UpdateCategory(context.Background(), 999, cmd)
		if !errors.Is(err, domain.ErrCategoryNotFound) {
			t.Errorf("expected ErrCategoryNotFound, got %v", err)
		}
	})
}

func TestCategoryService_DeleteCategory(t *testing.T) {
	existingCategory := &entity.Category{
		ID:        1,
		Name:      "Tech",
		Slug:      "tech",
		CreatedAt: time.Now(),
	}

	mockRepo := &mocks.MockCategoryRepository{
		FindByIDFunc: func(ctx context.Context, id int32) (*entity.Category, error) {
			if id == 1 {
				return existingCategory, nil
			}
			return nil, domain.ErrCategoryNotFound
		},
		GetPostCountFunc: func(ctx context.Context, id int32) (int64, error) {
			return 0, nil
		},
		DeleteFunc: func(ctx context.Context, id int32) error {
			return nil
		},
	}

	svc := NewCategoryService(mockRepo)

	t.Run("successful delete", func(t *testing.T) {
		err := svc.DeleteCategory(context.Background(), 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("category not found", func(t *testing.T) {
		err := svc.DeleteCategory(context.Background(), 999)
		if !errors.Is(err, domain.ErrCategoryNotFound) {
			t.Errorf("expected ErrCategoryNotFound, got %v", err)
		}
	})

	t.Run("category has posts", func(t *testing.T) {
		mockRepo.GetPostCountFunc = func(ctx context.Context, id int32) (int64, error) {
			return 5, nil
		}
		err := svc.DeleteCategory(context.Background(), 1)
		if !errors.Is(err, domain.ErrCategoryHasPosts) {
			t.Errorf("expected ErrCategoryHasPosts, got %v", err)
		}
	})
}
