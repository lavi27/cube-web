package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var (
	Logger = logrus.New()
)

// implements logger.Interface
type GormLogger struct{}

func (GormLogger) LogMode(logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}
func (GormLogger) Info(context.Context, string, ...interface{}) {
	if l.LogLevel >= Info {
		l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}
func (GormLogger) Warn(context.Context, string, ...interface{}) {
	if l.LogLevel >= Warn {
		l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}
func (GormLogger) Error(context.Context, string, ...interface{}) {
	if l.LogLevel >= Error {
		l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}
func (GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= Error && (!errors.Is(err, ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == Info:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func init() {
	// logrus.SetLevel(logrus.InfoLevel)
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	// logrus.SetReportCaller(true)

	logFile := &lumberjack.Logger{
		Filename:   "./log/log.txt",
		MaxSize:    20,
		MaxBackups: 5,
		MaxAge:     28,
		Compress:   true,
	}

	var logWriter io.Writer
	if IsDebug {
		logWriter = io.MultiWriter(logFile, os.Stdout)
	} else {
		logWriter = logFile
	}
	Logger.SetOutput(logWriter)
}
