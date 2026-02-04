package database

import (
	"context"
	"fmt"

	"github.com/fffeng99999/hcp-server/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Ping to verify connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return rdb, nil
}
