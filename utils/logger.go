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
	_DevLogWriter  io.Writer
	_HTTPLogWriter io.Writer
	_SQLLogWriter  io.Writer
	_ErrLogWriter  io.Writer
	HTTPLogger     *logrus.Logger
	DevLogger      *logrus.Logger
	ErrLogger      *logrus.Logger
	SQLLogger      logger.Interface
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

	DevLogFile, err := os.OpenFile("./log/DevLog.conf", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	ErrLogFile, err := os.OpenFile("./log/ErrLog.conf", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	if IsDebug {
		_HTTPLogWriter = io.MultiWriter(HTTPLogFile, os.Stdout)
		_SQLLogWriter = io.MultiWriter(SQLLogFile, os.Stdout)
	} else {
		_HTTPLogWriter = HTTPLogFile
		_SQLLogWriter = SQLLogFile
	}
	_DevLogWriter = io.MultiWriter(DevLogFile, os.Stdout)
	_ErrLogWriter = io.MultiWriter(ErrLogFile, os.Stdout)

	HTTPLogger = logrus.New()
	HTTPLogger.SetOutput(_HTTPLogWriter)

	DevLogger = logrus.New()
	DevLogger.SetOutput(_DevLogWriter)

	ErrLogger = logrus.New()
	ErrLogger.SetOutput(_ErrLogWriter)

	SQLLogger = logger.New(
		log.New(_SQLLogWriter, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  false,
		},
	)
}
