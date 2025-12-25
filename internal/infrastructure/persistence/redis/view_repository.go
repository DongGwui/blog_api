package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ydonggwui/blog-api/internal/domain/repository"
)

type viewRepository struct {
	client *redis.Client
}

func NewViewRepository(client *redis.Client) repository.ViewRepository {
	return &viewRepository{client: client}
}

func (r *viewRepository) SetViewIfNotExists(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, "1", ttl).Result()
}

func (r *viewRepository) HasView(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
