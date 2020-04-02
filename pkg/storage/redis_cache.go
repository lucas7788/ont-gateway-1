package storage

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
)

// Cache defines interface for cache
type Cache interface {
	Set(k, v string, expire time.Duration) error
	Get(k string) (string, error)
}

var _ Cache = (*RedisCache)(nil)

// NewRedisCache for redis cache
func NewRedisCache(conf *config.RedisConfig) *RedisCache {
	return &RedisCache{r: NewRedis(conf)}
}

// RedisCache is cached based on redis
type RedisCache struct {
	r *redis.Client
}

// Set value
func (cache *RedisCache) Set(k, v string, expire time.Duration) (err error) {
	err = cache.r.Set(k, v, expire).Err()
	return
}

// Get value
func (cache *RedisCache) Get(k string) (v string, err error) {
	v, err = cache.r.Get(k).Result()
	return
}
