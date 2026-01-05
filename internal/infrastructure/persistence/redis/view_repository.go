package redis

import (
	"context"
	"fmt"
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
	result, err := r.client.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		return false, fmt.Errorf("viewRepository.SetViewIfNotExists: %w", err)
	}
	return result, nil
}

func (r *viewRepository) HasView(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("viewRepository.HasView: %w", err)
	}
	return exists > 0, nil
}
