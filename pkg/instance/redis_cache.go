package instance

import (
	"sync"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/storage"
)

var (
	instanceRedisCache *storage.RedisCache
	redisOnce          sync.Once
)

// RedisCache is singleton for storage.RedisCache
func RedisCache() *storage.RedisCache {
	redisOnce.Do(func() {
		config := config.Load().RedisCacheConfig
		instanceRedisCache = storage.NewRedisCache(&config)
	})

	return instanceRedisCache
}
