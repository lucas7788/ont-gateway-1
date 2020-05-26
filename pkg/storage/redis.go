package storage

import (
	"crypto/tls"

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

	if conf.TLS {
		o.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	if conf.PoolSize > 0 {
		o.PoolSize = conf.PoolSize
	}
	client := redis.NewClient(o)

	return client
}
