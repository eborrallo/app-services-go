package cache

import (
	"app-services-go/configs"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// RedisCache Cache implements the Cache interface using github.com/go-redis/redis/v9.
type RedisCache struct {
	client   redis.UniversalClient
	duration time.Duration
}

func RedisConnection(cfg configs.RedisConfig) redis.UniversalClient {
	return redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)},
	})
}

// NewRedisCache creates a new instance of RedisCache.
func NewRedisCache(client redis.UniversalClient, duration time.Duration) *RedisCache {
	return &RedisCache{
		client:   client,
		duration: duration,
	}
}

// Get retrieves the value for a given key from the cache.
func (r *RedisCache) Get(key Key) ([]byte, error) {
	result, err := r.client.Get(context.Background(), string(key)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, errors.New("not found") // Cache miss, return nil without an error
	} else if err != nil {
		return nil, err
	}

	return []byte(result), nil
}

// Set stores a value for a given key in the cache.
func (r *RedisCache) Set(key Key, value Value) {
	jsonValue, _ := json.Marshal(value)
	r.client.Set(context.Background(), string(key), jsonValue, r.duration)
}
