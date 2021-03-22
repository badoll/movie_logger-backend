package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"time"
)

var loggerMap map[string]*logrus.Logger
var mLog *logrus.Logger

const logPath = "/home/ubuntu/movie_logger-backend_log"

func init() {
	var err error
	mLog, err = newLogger("mlog")
	if err != nil {
		panic(err)
	}
	loggerMap = make(map[string]*logrus.Logger)
	loggerMap["mlog"] = mLog
}

func GetLogger(name string) *logrus.Logger {
	if l, ok := loggerMap[name]; ok {
		return l
	}
	logger, err := newLogger(name)
	if err != nil {
		mLog.Errorf("GetLogger error, name: %s, err: %v", name, err)
		return mLog
	}
	loggerMap[name] = logger
	return logger
}

func GetDefaultLogger() *logrus.Logger {
	return mLog
}

func newLogger(name string) (*logrus.Logger, error) {
	logfname := fmt.Sprintf("%s/%s", logPath, name)
	/* 日志轮转相关函数
	`WithLinkName` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	 WithMaxAge 和 WithRotationCount二者只能设置一个
	`WithMaxAge` 设置文件清理前的最长保存时间
	`WithRotationCount` 设置文件清理前最多保存的个数
	*/
	writer, err := rotatelogs.New(
		logfname+".%Y%m%d",
		rotatelogs.WithLinkName(logfname),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil, err
	}
	logger := logrus.New()
	logger.SetOutput(writer)
	logger.SetLevel(logrus.DebugLevel)
	formatter := &logrus.TextFormatter{}
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	logger.SetFormatter(formatter)
	return logger, nil
}
