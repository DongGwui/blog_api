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

type adminRepository struct {
	queries *sqlc.Queries
}

func NewAdminRepository(queries *sqlc.Queries) repository.AdminRepository {
	return &adminRepository{queries: queries}
}

func (r *adminRepository) FindByID(ctx context.Context, id int32) (*entity.Admin, error) {
	admin, err := r.queries.GetAdminByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAdminNotFound
		}
		return nil, err
	}
	return toAdminEntity(admin), nil
}

func (r *adminRepository) FindByUsername(ctx context.Context, username string) (*entity.Admin, error) {
	admin, err := r.queries.GetAdminByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAdminNotFound
		}
		return nil, err
	}
	return toAdminEntity(admin), nil
}

func (r *adminRepository) Create(ctx context.Context, admin *entity.Admin) (*entity.Admin, error) {
	created, err := r.queries.CreateAdmin(ctx, sqlc.CreateAdminParams{
		Username: admin.Username,
		Password: admin.Password,
	})
	if err != nil {
		return nil, err
	}
	return toAdminEntity(created), nil
}

func (r *adminRepository) UpdatePassword(ctx context.Context, id int32, hashedPassword string) error {
	return r.queries.UpdateAdminPassword(ctx, sqlc.UpdateAdminPasswordParams{
		ID:       id,
		Password: hashedPassword,
	})
}
