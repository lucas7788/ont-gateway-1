// +build prod

package config

import "time"

var config = Value{
	Prod: true,
	RestConfig: RestConfig{
		PublicAddr:      "0.0.0.0:2020",
		IntraAddr:       "0.0.0.0:2021",
		GracefulUpgrade: false,
		PIDFile:         "/tmp/ont-gateway.pid",
		ReadTimeout:     time.Second * 3,
	},
	LoggerConfig: LoggerConfig{
		LogLevel: "info",
	},
	RedisCacheConfig: RedisConfig{
		Addr:     "master.sagamarket-prod.pbcbnm.apne1.cache.amazonaws.com:6379",
		TLS:      true,
		Password: "E062AB4QL+g9eo2RNrajw6LjcDc=",
	},
	MongoConfig: MongoConfig{
		ConnectionString: "mongodb://sagamarket:QvJ4ikuhNcZperMhDd8v@sagemarket-prod-mongo.cluster-cburg362yfls.ap-northeast-1.docdb.amazonaws.com:27017/ont?authSource=admin&maxidletimems=3000",
		Timeout:          time.Second * 3,
	},
	CICDConfig: CICDConfig{
		AddonDeployAPI: AkSkURL{
			Host: "af9a305da83ae11ea809406ad1589d0c-1615907691.ap-northeast-1.elb.amazonaws.com:8000",
			URI:  "/api/v1/ss",
		},
	},
}
