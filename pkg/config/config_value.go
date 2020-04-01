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
}
