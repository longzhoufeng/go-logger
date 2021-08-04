package zap

import (
	"github.com/longzhoufeng/go-logger"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	go_logger.Options
}

type callerSkipKey struct{}

func WithCallerSkip(i int) go_logger.Option {
	return go_logger.SetOption(callerSkipKey{}, i)
}

type configKey struct{}

// WithConfig pass zap.Config to logger
func WithConfig(c zap.Config) go_logger.Option {
	return go_logger.SetOption(configKey{}, c)
}

type encoderConfigKey struct{}

// WithEncoderConfig pass zapcore.EncoderConfig to logger
func WithEncoderConfig(c zapcore.EncoderConfig) go_logger.Option {
	return go_logger.SetOption(encoderConfigKey{}, c)
}

type namespaceKey struct{}

func WithNamespace(namespace string) go_logger.Option {
	return go_logger.SetOption(namespaceKey{}, namespace)
}

type writerKey struct{}

func WithOutput(out io.Writer) go_logger.Option {
	return go_logger.SetOption(writerKey{}, out)
}
