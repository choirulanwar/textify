package queue

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func NewLogger() *Logger {
	logger := logrus.New()
	return &Logger{logger: logger}
}

func (logger *Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	logger.logger.Debugf(format, v...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.logger.Debug(args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.logger.Info(args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.logger.Warn(args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.logger.Error(args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.logger.Fatal(args...)
}
