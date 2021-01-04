package logs

import (
	"context"

	"github.com/sirupsen/logrus"
)

var stdLogger = &logrusLogger{logger, nil, ""}

func WithPackage(pkg string) Logger {
	return stdLogger.WithPackage(pkg)
}

func WithContext(ctx context.Context) Logger {
	return stdLogger.WithContext(ctx)
}

func EnableDebug() {
	logger.Info("logger is entering debug mode")
	logger.SetLevel(logrus.DebugLevel)
}

func DisableDebug() {
	logger.Info("logger is leaving debug mode")
	logger.SetLevel(logrus.InfoLevel)
}
