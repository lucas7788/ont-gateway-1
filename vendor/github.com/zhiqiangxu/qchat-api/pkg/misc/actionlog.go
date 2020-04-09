package misc

import (
	"io/ioutil"
	"os"
	"time"

	"encoding/json"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/zhiqiangxu/qrpc"
)

// ActionLog for actionlog
type ActionLog struct {
	logger *logrus.Logger
}

// NewActionLog is ctor for ActionLog
func NewActionLog(path string) (al *ActionLog) {
	log := logrus.New()
	log.SetOutput(ioutil.Discard)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	path += "/log"
	infoW, err := rotatelogs.New(
		path+"."+time.Now().Format("20060102"),
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(86400)*time.Second),
	)

	if err != nil {
		panic(err)
	}

	hook := lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel: infoW,
		},
		&ActionLogFormatter{},
	)
	log.AddHook(hook)

	al = &ActionLog{logger: log}

	return
}

// Log writes an actionlog
func (al *ActionLog) Log(app int, action string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["action"] = action
	data["time"] = time.Now().UnixNano()
	data["app"] = app

	bytes, _ := json.Marshal(data)
	al.logger.Infoln(qrpc.String(bytes) + "\r\n")
}

// ActionLogFormatter is formatter for actionlog
type ActionLogFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

// Format implements ActionLogFormatter
func (f *ActionLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}
