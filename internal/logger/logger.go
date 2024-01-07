package logger

import (
	"fmt"
	"log"
	"sync"
)

type Logger interface {
	Print(v ...any)
}

type LogLevel uint

const (
	LogLevelOff   LogLevel = iota
	LogLevelError LogLevel = 1
	LogLevelInfo  LogLevel = 2
	LogLevelDebug LogLevel = 3
)

type impl struct {
	logLevel LogLevel
	logger   Logger
}

var (
	initLogger     sync.Once
	loggerInstance impl
)

func Init(logLevel LogLevel, logger Logger) {
	if logger == nil {
		logger = log.Default()
	}

	initLogger.Do(func() {
		loggerInstance = impl{
			logLevel: logLevel,
			logger:   logger,
		}
	})
}

func (i *impl) log(l LogLevel, format string, args ...any) { // notest
	if i.logLevel < l {
		return
	}

	if i.logger == nil {
		return
	}

	i.logger.Print(fmt.Sprintf(format, args...))
}

func Debug(format string, args ...any) {
	loggerInstance.log(LogLevelDebug, format, args...)
}

func Info(format string, args ...any) {
	loggerInstance.log(LogLevelInfo, format, args...)
}

func Error(format string, args ...any) {
	loggerInstance.log(LogLevelError, format, args...)
}
