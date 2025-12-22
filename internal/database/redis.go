package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ydonggwui/blog-api/internal/config"
)

func NewRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
		DB:       0,
	})

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return client, nil
}
