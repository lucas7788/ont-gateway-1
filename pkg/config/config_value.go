// +build !test,!prod

package config

import "time"

var config = Value{
	RestConfig: RestConfig{
		PublicAddr:      ":8080",
		IntraAddr:       ":8081",
		GracefulUpgrade: false,
		PIDFile:         "/tmp/ont-gateway.pid",
		ReadTimeout:     time.Second * 3,
	},
	LoggerConfig: LoggerConfig{
		LogLevel: "debug",
	},
	RedisCacheConfig: RedisConfig{
		Addr: "172.168.3.46:6379",
	},
	MongoConfig: MongoConfig{
		ConnectionString: "mongodb://172.168.3.46:27017/ont",
		Timeout:          time.Second * 3,
	},
	CICDConfig: CICDConfig{
		AddonDeployAPI: "",
	},
}
