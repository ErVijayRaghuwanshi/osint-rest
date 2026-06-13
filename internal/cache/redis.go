package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache is an implementation of Cache using Redis.
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache instantiates a new Redis cache client.
func NewRedisCache(addr, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisCache{
		client: client,
		ctx:    context.Background(),
	}
}

// Get fetches a key from Redis. Returns false on cache miss or error.
func (r *RedisCache) Get(key string) ([]byte, bool) {
	val, err := r.client.Get(r.ctx, key).Bytes()
	if err == redis.Nil {
		return nil, false
	} else if err != nil {
		return nil, false
	}
	return val, true
}

// Set stores a key in Redis with a TTL.
func (r *RedisCache) Set(key string, value []byte, ttl time.Duration) {
	_ = r.client.Set(r.ctx, key, value, ttl).Err()
}

// Delete removes a key from Redis.
func (r *RedisCache) Delete(key string) {
	_ = r.client.Del(r.ctx, key).Err()
}

// Flush removes all keys in the current Redis database.
func (r *RedisCache) Flush() {
	_ = r.client.FlushDB(r.ctx).Err()
}

// Close closes the underlying Redis connection pool.
func (r *RedisCache) Close() error {
	return r.client.Close()
}
