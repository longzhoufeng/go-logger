package zap

import (
	"context"
	"fmt"
	"github.com/longzhoufeng/go-logger"
	"io"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zaplog struct {
	cfg  zap.Config
	zap  *zap.Logger
	opts go_logger.Options
	sync.RWMutex
	fields map[string]interface{}
}

func (l *zaplog) Init(opts ...go_logger.Option) error {
	//var err error

	for _, o := range opts {
		o(&l.opts)
	}

	zapConfig := zap.NewProductionConfig()
	if zconfig, ok := l.opts.Context.Value(configKey{}).(zap.Config); ok {
		zapConfig = zconfig
	}

	if zcconfig, ok := l.opts.Context.Value(encoderConfigKey{}).(zapcore.EncoderConfig); ok {
		zapConfig.EncoderConfig = zcconfig
	}

	writer, ok := l.opts.Context.Value(writerKey{}).(io.Writer)
	if !ok {
		writer = os.Stdout
	}

	skip, ok := l.opts.Context.Value(callerSkipKey{}).(int)
	if !ok || skip < 1 {
		skip = 1
	}

	// Set log Level if not default
	zapConfig.Level = zap.NewAtomicLevel()
	if l.opts.Level != go_logger.InfoLevel {
		zapConfig.Level.SetLevel(loggerToZapLevel(l.opts.Level))
	}
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapConfig.EncoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(writer)),
		zapConfig.Level)

	log := zap.New(logCore, zap.AddCaller(), zap.AddCallerSkip(skip), zap.AddStacktrace(zap.DPanicLevel))

	// Adding seed fields if exist
	if l.opts.Fields != nil {
		data := []zap.Field{}
		for k, v := range l.opts.Fields {
			data = append(data, zap.Any(k, v))
		}
		log = log.With(data...)
	}

	// Adding namespace
	if namespace, ok := l.opts.Context.Value(namespaceKey{}).(string); ok {
		log = log.With(zap.Namespace(namespace))
	}

	// defer log.Sync() ??

	l.cfg = zapConfig
	l.zap = log
	l.fields = make(map[string]interface{})

	return nil
}

func (l *zaplog) Fields(fields map[string]interface{}) go_logger.Logger {
	l.Lock()
	nfields := make(map[string]interface{}, len(l.fields))
	for k, v := range l.fields {
		nfields[k] = v
	}
	l.Unlock()
	for k, v := range fields {
		nfields[k] = v
	}

	data := make([]zap.Field, 0, len(nfields))
	for k, v := range fields {
		data = append(data, zap.Any(k, v))
	}

	zl := &zaplog{
		cfg:    l.cfg,
		zap:    l.zap,
		opts:   l.opts,
		fields: nfields,
	}

	return zl
}

func (l *zaplog) Error(err error) go_logger.Logger {
	return l.Fields(map[string]interface{}{"error": err})
}

func (l *zaplog) Log(level go_logger.Level, args ...interface{}) {
	l.RLock()
	data := make([]zap.Field, 0, len(l.fields))
	for k, v := range l.fields {
		data = append(data, zap.Any(k, v))
	}
	l.RUnlock()

	lvl := loggerToZapLevel(level)
	msg := fmt.Sprint(args...)
	switch lvl {
	case zap.DebugLevel:
		l.zap.Debug(msg, data...)
	case zap.InfoLevel:
		l.zap.Info(msg, data...)
	case zap.WarnLevel:
		l.zap.Warn(msg, data...)
	case zap.ErrorLevel:
		l.zap.Error(msg, data...)
	case zap.FatalLevel:
		l.zap.Fatal(msg, data...)
	}
}

func (l *zaplog) Logf(level go_logger.Level, format string, args ...interface{}) {
	l.RLock()
	data := make([]zap.Field, 0, len(l.fields))
	for k, v := range l.fields {
		data = append(data, zap.Any(k, v))
	}
	l.RUnlock()

	lvl := loggerToZapLevel(level)
	msg := fmt.Sprintf(format, args...)
	switch lvl {
	case zap.DebugLevel:
		l.zap.Debug(msg, data...)
	case zap.InfoLevel:
		l.zap.Info(msg, data...)
	case zap.WarnLevel:
		l.zap.Warn(msg, data...)
	case zap.ErrorLevel:
		l.zap.Error(msg, data...)
	case zap.FatalLevel:
		l.zap.Fatal(msg, data...)
	}
}

func (l *zaplog) String() string {
	return "zap"
}

func (l *zaplog) Options() go_logger.Options {
	return l.opts
}

// New builds a new logger based on options
func NewLogger(opts ...go_logger.Option) (go_logger.Logger, error) {
	// Default options
	options := go_logger.Options{
		Level:   go_logger.InfoLevel,
		Fields:  make(map[string]interface{}),
		Out:     os.Stderr,
		Context: context.Background(),
	}

	l := &zaplog{opts: options}
	if err := l.Init(opts...); err != nil {
		return nil, err
	}

	return l, nil
}

func loggerToZapLevel(level go_logger.Level) zapcore.Level {
	switch level {
	case go_logger.TraceLevel, go_logger.DebugLevel:
		return zap.DebugLevel
	case go_logger.InfoLevel:
		return zap.InfoLevel
	case go_logger.WarnLevel:
		return zap.WarnLevel
	case go_logger.ErrorLevel:
		return zap.ErrorLevel
	case go_logger.FatalLevel:
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func zapToLoggerLevel(level zapcore.Level) go_logger.Level {
	switch level {
	case zap.DebugLevel:
		return go_logger.DebugLevel
	case zap.InfoLevel:
		return go_logger.InfoLevel
	case zap.WarnLevel:
		return go_logger.WarnLevel
	case zap.ErrorLevel:
		return go_logger.ErrorLevel
	case zap.FatalLevel:
		return go_logger.FatalLevel
	default:
		return go_logger.InfoLevel
	}
}
