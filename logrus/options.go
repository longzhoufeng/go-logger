package logrus

import (
	"github.com/longzhoufeng/go-logger"
	"github.com/sirupsen/logrus"
)

type Options struct {
	go_logger.Options
	Formatter logrus.Formatter
	Hooks     logrus.LevelHooks
	// Flag for whether to log caller info (off by default)
	ReportCaller bool
	// Exit Function to call when FatalLevel log
	ExitFunc func(int)
}

type formatterKey struct{}

func WithTextTextFormatter(formatter *logrus.TextFormatter) go_logger.Option {
	return go_logger.SetOption(formatterKey{}, formatter)
}
func WithJSONFormatter(formatter *logrus.JSONFormatter) go_logger.Option {
	return go_logger.SetOption(formatterKey{}, formatter)
}

type hooksKey struct{}

func WithLevelHooks(hooks logrus.LevelHooks) go_logger.Option {
	return go_logger.SetOption(hooksKey{}, hooks)
}

type reportCallerKey struct{}

// warning to use this option. because logrus doest not open CallerDepth option
// this will only print this package
func ReportCaller() go_logger.Option {
	return go_logger.SetOption(reportCallerKey{}, true)
}

type exitKey struct{}

func WithExitFunc(exit func(int)) go_logger.Option {
	return go_logger.SetOption(exitKey{}, exit)
}

type logrusLoggerKey struct{}

func WithLogger(l logrus.StdLogger) go_logger.Option {
	return go_logger.SetOption(logrusLoggerKey{}, l)
}
