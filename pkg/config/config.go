package config

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	// MongoConfig for mongo
	MongoConfig struct {
		ConnectionString string
	}

	// RedisConfig for redis
	RedisConfig struct {
		Addr     string
		Password string
		PoolSize int
	}

	// RestConfig for rest
	RestConfig struct {
		PublicAddr      string
		IntraAddr       string
		GracefulUpgrade bool
		PIDFile         string
		ReadTimeout     time.Duration
	}

	// LoggerConfig for logger
	LoggerConfig struct {
		LogLevel string
	}
	// Value is combined config info
	Value struct {
		Prod             bool
		RestConfig       RestConfig
		LoggerConfig     LoggerConfig
		RedisCacheConfig RedisConfig
		MongoConfig      MongoConfig
	}
)

// Load config
func Load() *Value {
	return &config
}

func init() {

	bytes, _ := json.Marshal(config)
	fmt.Println("ont-gateway conf", string(bytes))

}
