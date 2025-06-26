package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fehepe/pet-store/backend/internal/config"
	"github.com/go-redis/redis/v8"
)

// CacheInterface defines the interface for cache operations
type CacheInterface interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error
	Delete(ctx context.Context, keys ...string) error

	InvalidatePattern(ctx context.Context, pattern string) error

	Close() error
	Ping() error
}

// Ensure Cache implements CacheInterface
var _ CacheInterface = (*Cache)(nil)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(cfg *config.Config) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Cache{
		client: client,
		ttl:    5 * time.Minute, // Default TTL
	}, nil
}

func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key not found")
	} else if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	expiration := c.ttl
	if len(ttl) > 0 {
		expiration = ttl[0]
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}

func (c *Cache) InvalidatePattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var keys []string

	for {
		var err error
		var result []string
		result, cursor, err = c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		keys = append(keys, result...)

		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		return c.Delete(ctx, keys...)
	}

	return nil
}

// Close closes the cache connection
func (c *Cache) Close() error {
	return c.client.Close()
}

// Ping checks if the cache connection is alive
func (c *Cache) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.client.Ping(ctx).Err()
}

// Cache key helpers
func PetCacheKey(storeID, petID string) string {
	return fmt.Sprintf("pet:%s:%s", storeID, petID)
}

func StoreCacheKey(storeID string) string {
	return fmt.Sprintf("store:%s", storeID)
}
