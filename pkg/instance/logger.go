package instance

import (
	"sync"

	"github.com/zhiqiangxu/ont-gateway/pkg/logger"
	"go.uber.org/zap"
)

var (
	loggerInstance *zap.Logger
	loggerOnce     sync.Once
)

// Logger is singleton for zap.Logger
func Logger() *zap.Logger {
	loggerOnce.Do(func() {
		loggerInstance = logger.New()
	})

	return loggerInstance
}
