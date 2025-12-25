package mocks

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// MockCategoryRepository is a mock implementation of CategoryRepository
type MockCategoryRepository struct {
	FindAllFunc           func(ctx context.Context) ([]entity.Category, error)
	FindByIDFunc          func(ctx context.Context, id int32) (*entity.Category, error)
	FindBySlugFunc        func(ctx context.Context, slug string) (*entity.Category, error)
	CreateFunc            func(ctx context.Context, category *entity.Category) (*entity.Category, error)
	UpdateFunc            func(ctx context.Context, category *entity.Category) (*entity.Category, error)
	DeleteFunc            func(ctx context.Context, id int32) error
	SlugExistsFunc        func(ctx context.Context, slug string) (bool, error)
	SlugExistsExceptFunc  func(ctx context.Context, slug string, excludeID int32) (bool, error)
	GetPostCountFunc      func(ctx context.Context, id int32) (int64, error)
}

func (m *MockCategoryRepository) FindAll(ctx context.Context) ([]entity.Category, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockCategoryRepository) FindByID(ctx context.Context, id int32) (*entity.Category, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockCategoryRepository) FindBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	if m.FindBySlugFunc != nil {
		return m.FindBySlugFunc(ctx, slug)
	}
	return nil, nil
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, category)
	}
	return nil, nil
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, category)
	}
	return nil, nil
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id int32) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockCategoryRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	if m.SlugExistsFunc != nil {
		return m.SlugExistsFunc(ctx, slug)
	}
	return false, nil
}

func (m *MockCategoryRepository) SlugExistsExcept(ctx context.Context, slug string, excludeID int32) (bool, error) {
	if m.SlugExistsExceptFunc != nil {
		return m.SlugExistsExceptFunc(ctx, slug, excludeID)
	}
	return false, nil
}

func (m *MockCategoryRepository) GetPostCount(ctx context.Context, id int32) (int64, error) {
	if m.GetPostCountFunc != nil {
		return m.GetPostCountFunc(ctx, id)
	}
	return 0, nil
}
