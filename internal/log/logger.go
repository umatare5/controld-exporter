// Package log provides a simple logging interface.
package log

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// SetLogLevel sets the log level based on the debug flag
func SetLogLevel(logLevel string) {
	if logLevel == "warn" {
		logger.SetLevel(logrus.WarnLevel)
		return
	}
	if logLevel == "error" {
		logger.SetLevel(logrus.ErrorLevel)
		return
	}
	if logLevel == "debug" {
		logger.SetLevel(logrus.DebugLevel)
		return
	}
	logger.SetLevel(logrus.InfoLevel)
}

// Info logs a message at level Info.
func Info(args ...any) {
	logger.Info(args...)
}

// Infof logs a message at level Info.
func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

// Warnf logs a message at level Warn.
func Warnf(format string, args ...any) {
	logger.Warnf(format, args...)
}

// Errorf logs a message at level Error.
func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

// Debugf logs a message at level Debug.
func Debugf(format string, args ...any) {
	logger.Debugf(format, args...)
}

// Fatal logs a message at level Fatal.
func Fatal(args ...any) {
	logger.Fatal(args...)
}
