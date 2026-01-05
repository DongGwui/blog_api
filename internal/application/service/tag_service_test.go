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

func TestTagService_ListTags(t *testing.T) {
	mockRepo := &mocks.MockTagRepository{
		FindAllFunc: func(ctx context.Context) ([]entity.Tag, error) {
			return []entity.Tag{
				{ID: 1, Name: "Go", Slug: "go", CreatedAt: time.Now()},
				{ID: 2, Name: "Rust", Slug: "rust", CreatedAt: time.Now()},
			}, nil
		},
	}

	svc := NewTagService(mockRepo)
	tags, err := svc.ListTags(context.Background())

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}
}

func TestTagService_ListTagsWithPostCount(t *testing.T) {
	mockRepo := &mocks.MockTagRepository{
		FindAllWithPostCountFunc: func(ctx context.Context) ([]entity.Tag, error) {
			return []entity.Tag{
				{ID: 1, Name: "Go", Slug: "go", PostCount: 10, CreatedAt: time.Now()},
				{ID: 2, Name: "Rust", Slug: "rust", PostCount: 5, CreatedAt: time.Now()},
			}, nil
		},
	}

	svc := NewTagService(mockRepo)
	tags, err := svc.ListTagsWithPostCount(context.Background())

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}
	if tags[0].PostCount != 10 {
		t.Errorf("expected post count 10, got %d", tags[0].PostCount)
	}
}

func TestTagService_GetTagByID(t *testing.T) {
	expected := &entity.Tag{
		ID:        1,
		Name:      "Go",
		Slug:      "go",
		CreatedAt: time.Now(),
	}

	mockRepo := &mocks.MockTagRepository{
		FindByIDFunc: func(ctx context.Context, id int32) (*entity.Tag, error) {
			if id == 1 {
				return expected, nil
			}
			return nil, domain.ErrTagNotFound
		},
	}

	svc := NewTagService(mockRepo)

	t.Run("existing tag", func(t *testing.T) {
		tag, err := svc.GetTagByID(context.Background(), 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if tag.ID != expected.ID {
			t.Errorf("expected ID %d, got %d", expected.ID, tag.ID)
		}
	})

	t.Run("non-existing tag", func(t *testing.T) {
		_, err := svc.GetTagByID(context.Background(), 999)
		if !errors.Is(err, domain.ErrTagNotFound) {
			t.Errorf("expected ErrTagNotFound, got %v", err)
		}
	})
}

func TestTagService_CreateTag(t *testing.T) {
	createdTag := &entity.Tag{
		ID:        1,
		Name:      "Go",
		Slug:      "go",
		CreatedAt: time.Now(),
	}

	mockRepo := &mocks.MockTagRepository{
		SlugExistsFunc: func(ctx context.Context, slug string) (bool, error) {
			return false, nil
		},
		CreateFunc: func(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
			tag.ID = 1
			return tag, nil
		},
		FindByIDFunc: func(ctx context.Context, id int32) (*entity.Tag, error) {
			return createdTag, nil
		},
	}

	svc := NewTagService(mockRepo)

	t.Run("successful creation", func(t *testing.T) {
		cmd := domainService.CreateTagCommand{
			Name: "Go",
			Slug: "go",
		}
		tag, err := svc.CreateTag(context.Background(), cmd)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if tag.Name != "Go" {
			t.Errorf("expected name Go, got %s", tag.Name)
		}
	})

	t.Run("slug already exists", func(t *testing.T) {
		mockRepo.SlugExistsFunc = func(ctx context.Context, slug string) (bool, error) {
			return true, nil
		}

		cmd := domainService.CreateTagCommand{
			Name: "Go",
			Slug: "go",
		}
		_, err := svc.CreateTag(context.Background(), cmd)
		if !errors.Is(err, domain.ErrTagSlugExists) {
			t.Errorf("expected ErrTagSlugExists, got %v", err)
		}
	})
}

func TestTagService_UpdateTag(t *testing.T) {
	existingTag := &entity.Tag{
		ID:        1,
		Name:      "Go",
		Slug:      "go",
		CreatedAt: time.Now(),
	}

	mockRepo := &mocks.MockTagRepository{
		FindByIDFunc: func(ctx context.Context, id int32) (*entity.Tag, error) {
			if id == 1 {
				return existingTag, nil
			}
			return nil, domain.ErrTagNotFound
		},
		SlugExistsExceptFunc: func(ctx context.Context, slug string, excludeID int32) (bool, error) {
			return false, nil
		},
		UpdateFunc: func(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
			return tag, nil
		},
	}

	svc := NewTagService(mockRepo)

	t.Run("successful update", func(t *testing.T) {
		cmd := domainService.UpdateTagCommand{
			Name: "Golang",
			Slug: "golang",
		}
		tag, err := svc.UpdateTag(context.Background(), 1, cmd)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if tag == nil {
			t.Error("expected tag, got nil")
		}
	})

	t.Run("tag not found", func(t *testing.T) {
		cmd := domainService.UpdateTagCommand{
			Name: "Golang",
			Slug: "golang",
		}
		_, err := svc.UpdateTag(context.Background(), 999, cmd)
		if !errors.Is(err, domain.ErrTagNotFound) {
			t.Errorf("expected ErrTagNotFound, got %v", err)
		}
	})
}

func TestTagService_DeleteTag(t *testing.T) {
	existingTag := &entity.Tag{
		ID:        1,
		Name:      "Go",
		Slug:      "go",
		CreatedAt: time.Now(),
	}

	mockRepo := &mocks.MockTagRepository{
		FindByIDFunc: func(ctx context.Context, id int32) (*entity.Tag, error) {
			if id == 1 {
				return existingTag, nil
			}
			return nil, domain.ErrTagNotFound
		},
		DeleteFunc: func(ctx context.Context, id int32) error {
			return nil
		},
	}

	svc := NewTagService(mockRepo)

	t.Run("successful delete", func(t *testing.T) {
		err := svc.DeleteTag(context.Background(), 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("tag not found", func(t *testing.T) {
		err := svc.DeleteTag(context.Background(), 999)
		if !errors.Is(err, domain.ErrTagNotFound) {
			t.Errorf("expected ErrTagNotFound, got %v", err)
		}
	})
}
