package service

import (
	"context"

	"github.com/ydonggwui/blog-api/internal/domain"
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/domain/repository"
	domainService "github.com/ydonggwui/blog-api/internal/domain/service"
	"github.com/ydonggwui/blog-api/internal/util"
)

type tagService struct {
	tagRepo repository.TagRepository
}

// NewTagService creates a new tag service
func NewTagService(tagRepo repository.TagRepository) domainService.TagService {
	return &tagService{
		tagRepo: tagRepo,
	}
}

func (s *tagService) ListTags(ctx context.Context) ([]entity.Tag, error) {
	return s.tagRepo.FindAll(ctx)
}

func (s *tagService) ListTagsWithPostCount(ctx context.Context) ([]entity.Tag, error) {
	return s.tagRepo.FindAllWithPostCount(ctx)
}

func (s *tagService) GetTagByID(ctx context.Context, id int32) (*entity.Tag, error) {
	return s.tagRepo.FindByID(ctx, id)
}

func (s *tagService) GetTagBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	return s.tagRepo.FindBySlug(ctx, slug)
}

func (s *tagService) CreateTag(ctx context.Context, cmd domainService.CreateTagCommand) (*entity.Tag, error) {
	// Generate slug if not provided
	slug := cmd.Slug
	if slug == "" {
		slug = util.GenerateSlug(cmd.Name)
	}

	// Check if slug exists
	exists, err := s.tagRepo.SlugExists(ctx, slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrTagSlugExists
	}

	tag := entity.NewTag(cmd.Name, slug)

	created, err := s.tagRepo.Create(ctx, tag)
	if err != nil {
		return nil, err
	}

	return s.tagRepo.FindByID(ctx, created.ID)
}

func (s *tagService) UpdateTag(ctx context.Context, id int32, cmd domainService.UpdateTagCommand) (*entity.Tag, error) {
	// Check if tag exists
	existing, err := s.tagRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Generate slug if not provided
	slug := cmd.Slug
	if slug == "" {
		slug = util.GenerateSlug(cmd.Name)
	}

	// Check if slug exists (excluding current tag)
	if slug != existing.Slug {
		exists, err := s.tagRepo.SlugExistsExcept(ctx, slug, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domain.ErrTagSlugExists
		}
	}

	tag := &entity.Tag{
		ID:   id,
		Name: cmd.Name,
		Slug: slug,
	}

	_, err = s.tagRepo.Update(ctx, tag)
	if err != nil {
		return nil, err
	}

	return s.tagRepo.FindByID(ctx, id)
}

func (s *tagService) DeleteTag(ctx context.Context, id int32) error {
	// Check if tag exists
	_, err := s.tagRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.tagRepo.Delete(ctx, id)
}
