package instance

import (
	"sync"

	"github.com/zhiqiangxu/ont-gateway/pkg/logger"
	"go.uber.org/zap"
)

var (
	loggerInstance *zap.Logger
	loggerMu       sync.Mutex
)

// Logger is singleton for zap.Logger
func Logger() *zap.Logger {
	if loggerInstance != nil {
		return loggerInstance
	}

	loggerMu.Lock()
	defer loggerMu.Unlock()

	if loggerInstance != nil {
		return loggerInstance
	}

	loggerInstance = logger.New()
	return loggerInstance

}
