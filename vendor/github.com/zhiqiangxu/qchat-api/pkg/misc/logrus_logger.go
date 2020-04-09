package misc

import (
	"fmt"
	"io"
	"time"

	"io/ioutil"
	"os"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	r *logrus.Logger
}

func (l *logrusLogger) Infoln(args ...interface{}) {
	l.r.Infoln(args...)
}
func (l *logrusLogger) Debugln(args ...interface{}) {
	l.r.Debugln(args...)
}
func (l *logrusLogger) Errorln(args ...interface{}) {
	l.r.Errorln(args...)
}
func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.r.Infof(format, args...)
}
func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.r.Errorf(format, args...)
}
func (l *logrusLogger) Sync() error {
	return nil
}

// NewLogrusLogger creates a logrus.Logger
func NewLogrusLogger(path, level string) Logger {

	log := logrus.New()
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		panic(fmt.Sprintln("logrus.ParseLevel fail", err))
	}
	log.SetLevel(lv)
	log.SetOutput(ioutil.Discard)

	var (
		infoW  io.Writer
		errorW io.Writer
	)
	if path == "" {
		infoW = os.Stdout
		errorW = os.Stderr
	} else {
		infoPath := path + ".info.log"
		errorPath := path + ".error.log"
		infoW, err = rotatelogs.New(
			infoPath+"."+time.Now().Format("20060102"),
			rotatelogs.WithLinkName(infoPath),
			rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
			rotatelogs.WithRotationTime(time.Duration(86400)*time.Second),
		)
		errorW, err = rotatelogs.New(
			errorPath+"."+time.Now().Format("20060102"),
			rotatelogs.WithLinkName(errorPath),
			rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
		)
		if err != nil {
			panic(err)
		}
	}

	hook := lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel: infoW,
		},
		&logrus.JSONFormatter{},
	)
	hook.SetDefaultWriter(errorW)
	log.AddHook(hook)
	log.SetFormatter(&logrus.JSONFormatter{})
	return &logrusLogger{r: log}

}
