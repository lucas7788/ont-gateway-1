package misc

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger interface
type Logger interface {
	Infoln(args ...interface{})
	Debugln(args ...interface{})
	Errorln(args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Sync() error
}

type zapLogger struct {
	z *zap.Logger
	s *zap.SugaredLogger
}

func (l *zapLogger) Infoln(args ...interface{}) {
	l.s.Info(args...)
}

func (l *zapLogger) Debugln(args ...interface{}) {
	l.s.Debug(args...)
}

func (l *zapLogger) Errorln(args ...interface{}) {
	l.s.Error(args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.s.Infof(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.s.Errorf(format, args...)
}

func (l *zapLogger) Sync() error {
	return l.z.Sync()
}

// NewZapLogger for zap.Logger
func NewZapLogger(path, level string) Logger {

	if path != "" {
		panic("zap log path not supported yet")
	}

	lv := zapcore.Level(0)
	lv.Set(level)
	alv := zap.NewAtomicLevelAt(lv)
	ec := zap.NewProductionEncoderConfig()
	ec.CallerKey = ""
	ec.TimeKey = "time"
	ec.StacktraceKey = ""
	ec.EncodeTime = zapcore.ISO8601TimeEncoder
	config := zap.Config{
		Encoding:         "json",
		EncoderConfig:    ec,
		Level:            alv,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	z, err := config.Build()

	if err != nil {
		panic(fmt.Sprintf("NewZapLogger:%v", err))
	}

	s := z.Sugar()

	return &zapLogger{z: z, s: s}
}
