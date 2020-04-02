package instance

import (
	"sync"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/storage"
)

var (
	instanceRedisCache *storage.RedisCache
	lockRedisCache     sync.Mutex
)

// RedisCache is singleton for storage.RedisCache
func RedisCache() *storage.RedisCache {
	if instanceRedisCache != nil {
		return instanceRedisCache
	}

	lockRedisCache.Lock()
	defer lockRedisCache.Unlock()
	if instanceRedisCache != nil {
		return instanceRedisCache
	}

	config := config.Load().RedisCacheConfig
	instanceRedisCache = storage.NewRedisCache(&config)

	return instanceRedisCache
}
