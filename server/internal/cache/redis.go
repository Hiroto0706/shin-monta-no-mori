package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"shin-monta-no-mori/pkg/util"

	"github.com/redis/go-redis/v9"
)

// RedisClient インターフェースは、Redis クライアントの操作を定義します。
type RedisClient interface {
	Get(ctx context.Context, key string, i interface{}) error
	Set(ctx context.Context, key string, i interface{}, expiration time.Duration) error
	Del(ctx context.Context, key []string) error
}

// RedisContext は、Redis クライアントを管理する構造体です。
type RedisContext struct {
	client *redis.Client
}

// NewRedisClient は、新しい Redis クライアントを作成します。
func NewRedisClient(config util.Config) RedisClient {
	rds := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		DB:       config.RedisDbNumber,
		Password: config.RedisPassword,
	})

	return &RedisContext{
		client: rds,
	}
}

// Set は、データを Redis に保存します。
func (r *RedisContext) Set(ctx context.Context, key string, i interface{}, expiration time.Duration) error {
	data, err := json.Marshal(i)

	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := r.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set data in Redis: %w", err)
	}

	return nil
}

// Get は、Redis からデータを取得します。
func (r *RedisContext) Get(ctx context.Context, key string, i interface{}) error {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("key %s does not exist in Redis: %w", key, err)
		}

		return fmt.Errorf("failed to get data from Redis: %w", err)
	}

	if err := json.Unmarshal([]byte(data), i); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}

// Del は、指定されたキーパターンを含むすべてのキーを Redis から削除します。
func (r *RedisContext) Del(ctx context.Context, patterns []string) error {
	var errs []string
	for _, pattern := range patterns {
		iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			key := iter.Val()
			if err := r.client.Del(ctx, key).Err(); err != nil {
				errs = append(errs, fmt.Sprintf("failed to delete key %s: %v", key, err))
			}
		}

		if err := iter.Err(); err != nil {
			errs = append(errs, fmt.Sprintf("failed to iterate keys with pattern %s: %v", pattern, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, ", "))
	}

	return nil
}
