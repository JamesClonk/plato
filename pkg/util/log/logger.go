package log

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/spf13/viper"
)

var (
	logger   *slog.Logger
	disabled bool
)

func Initialize() {
	if !disabled {
		logger = newLogger(os.Stdout)
	}
}

func newLogger(writer *os.File) *slog.Logger {
	var logLevel slog.Level
	switch viper.GetString("plato.log_level") {
	case "info":
		logLevel = slog.LevelInfo
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelDebug
	}

	opts := &tint.Options{
		NoColor:    !isatty.IsTerminal(writer.Fd()),
		Level:      logLevel,
		TimeFormat: time.RFC3339,
	}
	logger = slog.New(tint.NewHandler(writer, opts))
	//Debugf("logger created, with log_level [%v]", logLevel)

	return logger
}

func Disable() {
	disabled = true
	logger = slog.New(slog.DiscardHandler)
}

func Info(format string, args ...interface{}) {
	logger.Info(format, args...)
}
func Infof(format string, args ...interface{}) {
	logger.Info(fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	logger.Warn(format, args...)
}
func Warnf(format string, args ...interface{}) {
	logger.Warn(fmt.Sprintf(format, args...))
}

func Debug(format string, args ...interface{}) {
	logger.Debug(format, args...)
}
func Debugf(format string, args ...interface{}) {
	logger.Debug(fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	logger.Error(format, args...)
}
func Errorf(format string, args ...interface{}) {
	logger.Error(fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...interface{}) {
	logger.Error(format, args...)
	os.Exit(1)
}
func Fatalf(format string, args ...interface{}) {
	logger.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}
