package config

import (
	"encoding/json"
	"fmt"
	"time"
)

type (

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
		Prod         bool
		RestConfig   RestConfig
		LoggerConfig LoggerConfig
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
