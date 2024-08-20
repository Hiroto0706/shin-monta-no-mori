package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"shin-monta-no-mori/pkg/util"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	GET(ctx context.Context, key string, i interface{}) error
	SET(ctx context.Context, key string, i interface{}, expiration time.Duration) error
}

type RedisContext struct {
	client *redis.Client
}

func NewRedisClient(config util.Config) RedisClient {
	rds := redis.NewClient(&redis.Options{
		Addr: config.RedisAddress,
		DB:   config.RedisDB,
	})

	return &RedisContext{
		client: rds,
	}
}

func (r *RedisContext) SET(ctx context.Context, key string, i interface{}, expiration time.Duration) error {
	data, err := json.Marshal(i)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := r.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set data in Redis: %w", err)
	}

	return nil
}

func (r *RedisContext) GET(ctx context.Context, key string, i interface{}) error {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to get data from Redis: %w", err)
	}

	if err := json.Unmarshal([]byte(data), i); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}
