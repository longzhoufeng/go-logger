package logrus

import (
	"context"
	"fmt"
	"github.com/longzhoufeng/go-logger"
	"os"

	"github.com/sirupsen/logrus"
)

type entryLogger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
	WithError(err error) *logrus.Entry

	Log(level logrus.Level, args ...interface{})
	Logf(level logrus.Level, format string, args ...interface{})
}

type logrusLogger struct {
	Logger entryLogger
	opts   Options
}

func (l *logrusLogger) Init(opts ...go_logger.Option) error {
	for _, o := range opts {
		o(&l.opts.Options)
	}

	if formatter, ok := l.opts.Context.Value(formatterKey{}).(logrus.Formatter); ok {
		l.opts.Formatter = formatter
	}
	if hs, ok := l.opts.Context.Value(hooksKey{}).(logrus.LevelHooks); ok {
		l.opts.Hooks = hs
	}
	if caller, ok := l.opts.Context.Value(reportCallerKey{}).(bool); ok && caller {
		l.opts.ReportCaller = caller
	}
	if exitFunction, ok := l.opts.Context.Value(exitKey{}).(func(int)); ok {
		l.opts.ExitFunc = exitFunction
	}

	switch ll := l.opts.Context.Value(logrusLoggerKey{}).(type) {
	case *logrus.Logger:
		// overwrite default options
		l.opts.Level = logrusToLoggerLevel(ll.GetLevel())
		l.opts.Out = ll.Out
		l.opts.Formatter = ll.Formatter
		l.opts.Hooks = ll.Hooks
		l.opts.ReportCaller = ll.ReportCaller
		l.opts.ExitFunc = ll.ExitFunc
		l.Logger = ll
	case *logrus.Entry:
		// overwrite default options
		el := ll.Logger
		l.opts.Level = logrusToLoggerLevel(el.GetLevel())
		l.opts.Out = el.Out
		l.opts.Formatter = el.Formatter
		l.opts.Hooks = el.Hooks
		l.opts.ReportCaller = el.ReportCaller
		l.opts.ExitFunc = el.ExitFunc
		l.Logger = ll
	case nil:
		log := logrus.New() // defaults
		log.SetLevel(loggerToLogrusLevel(l.opts.Level))
		log.SetOutput(l.opts.Out)
		log.SetFormatter(l.opts.Formatter)
		log.ReplaceHooks(l.opts.Hooks)
		log.SetReportCaller(l.opts.ReportCaller)
		log.ExitFunc = l.opts.ExitFunc
		l.Logger = log
	default:
		return fmt.Errorf("invalid logrus type: %T", ll)
	}

	return nil
}

func (l *logrusLogger) String() string {
	return "logrus"
}

func (l *logrusLogger) Fields(fields map[string]interface{}) go_logger.Logger {
	return &logrusLogger{l.Logger.WithFields(fields), l.opts}
}

func (l *logrusLogger) Log(level go_logger.Level, args ...interface{}) {
	l.Logger.Log(loggerToLogrusLevel(level), args...)
}

func (l *logrusLogger) Logf(level go_logger.Level, format string, args ...interface{}) {
	l.Logger.Logf(loggerToLogrusLevel(level), format, args...)
}

func (l *logrusLogger) Options() go_logger.Options {
	// FIXME: How to return full opts?
	return l.opts.Options
}

// New builds a new logger based on options
func NewLogger(opts ...go_logger.Option) go_logger.Logger {
	// Default options
	options := Options{
		Options: go_logger.Options{
			Level:   go_logger.InfoLevel,
			Fields:  make(map[string]interface{}),
			Out:     os.Stderr,
			Context: context.Background(),
		},
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		ReportCaller: false,
		ExitFunc:     os.Exit,
	}
	l := &logrusLogger{opts: options}
	_ = l.Init(opts...)
	return l
}

func loggerToLogrusLevel(level go_logger.Level) logrus.Level {
	switch level {
	case go_logger.TraceLevel:
		return logrus.TraceLevel
	case go_logger.DebugLevel:
		return logrus.DebugLevel
	case go_logger.InfoLevel:
		return logrus.InfoLevel
	case go_logger.WarnLevel:
		return logrus.WarnLevel
	case go_logger.ErrorLevel:
		return logrus.ErrorLevel
	case go_logger.FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func logrusToLoggerLevel(level logrus.Level) go_logger.Level {
	switch level {
	case logrus.TraceLevel:
		return go_logger.TraceLevel
	case logrus.DebugLevel:
		return go_logger.DebugLevel
	case logrus.InfoLevel:
		return go_logger.InfoLevel
	case logrus.WarnLevel:
		return go_logger.WarnLevel
	case logrus.ErrorLevel:
		return go_logger.ErrorLevel
	case logrus.FatalLevel:
		return go_logger.FatalLevel
	default:
		return go_logger.InfoLevel
	}
}
