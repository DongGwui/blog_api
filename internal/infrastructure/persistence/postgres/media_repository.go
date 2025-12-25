package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ydonggwui/blog-api/internal/database/sqlc"
	"github.com/ydonggwui/blog-api/internal/domain"
	"github.com/ydonggwui/blog-api/internal/domain/entity"
	"github.com/ydonggwui/blog-api/internal/domain/repository"
)

type mediaRepository struct {
	queries *sqlc.Queries
}

func NewMediaRepository(queries *sqlc.Queries) repository.MediaRepository {
	return &mediaRepository{queries: queries}
}

func (r *mediaRepository) FindByID(ctx context.Context, id int32) (*entity.Media, error) {
	media, err := r.queries.GetMediaByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrMediaNotFound
		}
		return nil, err
	}
	return toMediaEntity(media), nil
}

func (r *mediaRepository) Create(ctx context.Context, media *entity.Media) (*entity.Media, error) {
	created, err := r.queries.CreateMedia(ctx, toCreateMediaParams(media))
	if err != nil {
		return nil, err
	}
	return toMediaEntity(created), nil
}

func (r *mediaRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteMedia(ctx, id)
}

func (r *mediaRepository) List(ctx context.Context, limit, offset int32) ([]entity.Media, error) {
	media, err := r.queries.ListMedia(ctx, sqlc.ListMediaParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return toMediaEntities(media), nil
}

func (r *mediaRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountMedia(ctx)
}
