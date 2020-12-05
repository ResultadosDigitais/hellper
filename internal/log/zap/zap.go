package zap

import (
	"context"
	"hellper/internal/config"

	"hellper/internal/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapWriter interface {
	Write(...zap.Field)
}

type zapLogger interface {
	Check(zapcore.Level, string) zapWriter
	Sync() error
	With(...zap.Field) *zap.Logger
}

type zapLoggerDelegate struct {
	*zap.Logger
}

func NewZapLoggerDelegate(logger *zap.Logger) zapLogger {
	return &zapLoggerDelegate{
		Logger: logger,
	}
}

func (logger *zapLoggerDelegate) Check(level zapcore.Level, msg string) zapWriter {
	return logger.Logger.Check(level, msg)
}

func NewZapLogger(level log.Level, output log.Out) (*zap.Logger, error) {
	var (
		zapLevel  zapcore.Level
		errLevel  = zapLevel.Set(level.String())
		zapOutput = output.String()
		cfg       zap.Config
	)
	if errLevel != nil {
		return nil, errLevel
	}

	if config.Env.Environment == "development" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.Config{
			Level:         zap.NewAtomicLevelAt(zapLevel),
			DisableCaller: true,
			Development:   false,
			Encoding:      "json",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "time",
				LevelKey:       "level",
				NameKey:        "logger",
				MessageKey:     "message",
				StacktraceKey:  "stack",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
			},
			OutputPaths:      []string{zapOutput},
			ErrorOutputPaths: []string{zapOutput},
		}
	}

	return cfg.Build()
}

func NewZapLoggerDefault() *zap.Logger {
	zapLogger, _ := NewZapLogger(log.DEBUG, log.STDOUT)
	return zapLogger
}

type logger struct {
	adapter   zapLogger
	functions []log.ContextFunction
}

func (logger logger) log(ctx context.Context, level log.Level, msg string, values ...log.Value) {
	var zapLevel zapcore.Level
	_ = zapLevel.Set(level.String())
	if writer := logger.adapter.Check(zapLevel, msg); writer != nil {
		values = append(values, log.ResolveContextFunctions(ctx, logger.functions...)...)
		fields := make([]zapcore.Field, len(values))
		for index, logValue := range values {
			fields[index] = zap.Any(logValue.Name, logValue.Value)
		}
		writer.Write(fields...)
	}
}

func (logger logger) Debug(ctx context.Context, msg string, values ...log.Value) {
	logger.log(ctx, log.DEBUG, msg, values...)
}

func (logger logger) Info(ctx context.Context, msg string, values ...log.Value) {
	logger.log(ctx, log.INFO, msg, values...)
}

func (logger logger) Warn(ctx context.Context, msg string, values ...log.Value) {
	logger.log(ctx, log.WARN, msg, values...)
}

func (logger logger) Error(ctx context.Context, msg string, values ...log.Value) {
	logger.log(ctx, log.ERROR, msg, values...)
}

func (logger logger) With(values ...log.Value) log.Logger {
	fields := make([]zapcore.Field, len(values))
	for index, logValue := range values {
		fields[index] = zap.Any(logValue.Name, logValue.Value)
	}

	configuredLogger := logger.adapter.With(fields...)
	zapLogger := NewZapLoggerDelegate(configuredLogger)

	return New(zapLogger)
}

func (logger logger) Close(context.Context) {
	_ = logger.adapter.Sync()
}

func New(adapter zapLogger, functions ...log.ContextFunction) logger {
	return logger{
		adapter:   adapter,
		functions: functions,
	}
}

func NewDefault() logger {
	return New(
		NewZapLoggerDelegate(
			NewZapLoggerDefault(),
		),
	)
}
