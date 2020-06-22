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
		Timeout          time.Duration
	}

	// RedisConfig for redis
	RedisConfig struct {
		Addr     string
		Password string
		TLS      bool
		PoolSize int
	}

	// AkSkURL is protected with aksk
	AkSkURL struct {
		Host string
		URI  string
	}
	// CICDConfig for cicd
	CICDConfig struct {
		AddonDeployAPI AkSkURL
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

	// DocuConfig for docu
	DocuConfig struct {
		IntegratorKey string
		PrivateKey    string
		KeyPairID     string
		APIAccountID  string
		APIUserID     string
	}

	// MySQLConfig for mysql
	MySQLConfig struct {
		ConnectionString string
		ConnMaxLifetime  int
		MaxOpenConn      int
		MaxIdleConn      int
		ShowSQL          bool
		Slaves           []MySQLConfig
	}

	// EsignConfig for esign
	EsignConfig struct {
		DocuConfig DocuConfig
	}
	// Value is combined config info
	Value struct {
		Prod             bool
		RestConfig       RestConfig
		LoggerConfig     LoggerConfig
		RedisCacheConfig RedisConfig
		MongoConfig      MongoConfig
		CICDConfig       CICDConfig
		EsignConfig      EsignConfig
		MySQLConfig      MySQLConfig
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
