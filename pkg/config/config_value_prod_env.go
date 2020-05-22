// +build prod

package config

import "time"

var config = Value{
	Prod: true,
	RestConfig: RestConfig{
		PublicAddr:      ":2020",
		IntraAddr:       ":2021",
		GracefulUpgrade: false,
		PIDFile:         "/tmp/ont-gateway.pid",
		ReadTimeout:     time.Second * 3,
	},
	LoggerConfig: LoggerConfig{
		LogLevel: "info",
	},
	RedisCacheConfig: RedisConfig{
		Addr: "master.sagamarket-prod.pbcbnm.apne1.cache.amazonaws.com:6379",
	},
	MongoConfig: MongoConfig{
		ConnectionString: "mongodb://sagamarket:S+n7eL1+rSZKtxzC@sagemarket-prod-mongo.cluster-cburg362yfls.ap-northeast-1.docdb.amazonaws.com:27017/ont?authSource=ont&maxidletimems=3000",
		Timeout:          time.Second * 3,
	},
	CICDConfig: CICDConfig{
		AddonDeployAPI: AkSkURL{
			Host: "a0d771952588111ea89590659513bb5d-1585432770.ap-southeast-1.elb.amazonaws.com:8000",
			URI:  "/api/v1/ss",
		},
	},
}
