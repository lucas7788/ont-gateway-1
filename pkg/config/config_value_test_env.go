// +build test

package config

import "time"

var config = Value{
	RestConfig: RestConfig{
		PublicAddr:      ":2020",
		IntraAddr:       ":2021",
		GracefulUpgrade: false,
		PIDFile:         "/tmp/ont-gateway.pid",
		ReadTimeout:     time.Second * 3,
	},
	LoggerConfig: LoggerConfig{
		LogLevel: "debug",
	},
	RedisCacheConfig: RedisConfig{
		Addr: "localhost:6379",
	},
	MongoConfig: MongoConfig{
		ConnectionString: "mongodb://localhost:27017/ont",
		Timeout:          time.Second * 3,
	},
	CICDConfig: CICDConfig{
		AddonDeployAPI: AkSkURL{
			Host: "a0d771952588111ea89590659513bb5d-1585432770.ap-southeast-1.elb.amazonaws.com:8000",
			URI:  "/api/v1/ss",
		},
	},
}
