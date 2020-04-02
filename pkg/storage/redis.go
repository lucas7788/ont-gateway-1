package storage

import (
	"github.com/go-redis/redis"
	"github.com/zhiqiangxu/ont-gateway/pkg/config"
)

// NewRedis creates a redis client
func NewRedis(conf *config.RedisConfig) *redis.Client {

	o := &redis.Options{
		Addr: conf.Addr,
	}
	if conf.Password != "" {
		o.Password = conf.Password
	}
	if conf.PoolSize > 0 {
		o.PoolSize = conf.PoolSize
	}
	client := redis.NewClient(o)

	return client
}
