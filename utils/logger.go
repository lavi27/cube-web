package utils

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"
)

var (
	HTTPLogWriter io.Writer
	SQLLogWriter  io.Writer
	HTTPLogger    *logrus.Logger
	SQLLogger     logger.Interface
)

func init() {
	HTTPLogFile := &lumberjack.Logger{
		Filename:   "./log/HTTPLog.conf",
		MaxSize:    20,
		MaxBackups: 5,
		MaxAge:     28,
		Compress:   true,
	}

	SQLLogFile, err := os.OpenFile("./log/SQLLog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	if IsDebug {
		HTTPLogWriter = io.MultiWriter(HTTPLogFile, os.Stdout)
		SQLLogWriter = io.MultiWriter(SQLLogFile, os.Stdout)
	} else {
		HTTPLogWriter = HTTPLogFile
		SQLLogWriter = SQLLogFile
	}

	HTTPLogger = logrus.New()
	HTTPLogger.SetOutput(HTTPLogWriter)

	SQLLogger = logger.New(
		log.New(SQLLogWriter, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  false,
		},
	)
}
