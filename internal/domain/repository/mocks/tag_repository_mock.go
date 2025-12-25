package mocks

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain/entity"
)

// MockTagRepository is a mock implementation of TagRepository
type MockTagRepository struct {
	FindAllFunc              func(ctx context.Context) ([]entity.Tag, error)
	FindAllWithPostCountFunc func(ctx context.Context) ([]entity.Tag, error)
	FindByIDFunc             func(ctx context.Context, id int32) (*entity.Tag, error)
	FindBySlugFunc           func(ctx context.Context, slug string) (*entity.Tag, error)
	CreateFunc               func(ctx context.Context, tag *entity.Tag) (*entity.Tag, error)
	UpdateFunc               func(ctx context.Context, tag *entity.Tag) (*entity.Tag, error)
	DeleteFunc               func(ctx context.Context, id int32) error
	SlugExistsFunc           func(ctx context.Context, slug string) (bool, error)
	SlugExistsExceptFunc     func(ctx context.Context, slug string, excludeID int32) (bool, error)
	GetPostCountFunc         func(ctx context.Context, id int32) (int64, error)
}

func (m *MockTagRepository) FindAll(ctx context.Context) ([]entity.Tag, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockTagRepository) FindAllWithPostCount(ctx context.Context) ([]entity.Tag, error) {
	if m.FindAllWithPostCountFunc != nil {
		return m.FindAllWithPostCountFunc(ctx)
	}
	return nil, nil
}

func (m *MockTagRepository) FindByID(ctx context.Context, id int32) (*entity.Tag, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockTagRepository) FindBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	if m.FindBySlugFunc != nil {
		return m.FindBySlugFunc(ctx, slug)
	}
	return nil, nil
}

func (m *MockTagRepository) Create(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, tag)
	}
	return nil, nil
}

func (m *MockTagRepository) Update(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, tag)
	}
	return nil, nil
}

func (m *MockTagRepository) Delete(ctx context.Context, id int32) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockTagRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	if m.SlugExistsFunc != nil {
		return m.SlugExistsFunc(ctx, slug)
	}
	return false, nil
}

func (m *MockTagRepository) SlugExistsExcept(ctx context.Context, slug string, excludeID int32) (bool, error) {
	if m.SlugExistsExceptFunc != nil {
		return m.SlugExistsExceptFunc(ctx, slug, excludeID)
	}
	return false, nil
}

func (m *MockTagRepository) GetPostCount(ctx context.Context, id int32) (int64, error) {
	if m.GetPostCountFunc != nil {
		return m.GetPostCountFunc(ctx, id)
	}
	return 0, nil
}
