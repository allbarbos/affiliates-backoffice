package log

import (
	"context"
	"log"

	"affiliates-backoffice-backend/pkg/config"
	ctxpkg "affiliates-backoffice-backend/pkg/context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type LoggerI interface {
	With(ctx context.Context, args ...interface{}) LoggerI
	Desugar() *zap.Logger
	Sync() error

	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

var (
	levels = map[string]zapcore.Level{
		"debug": zap.DebugLevel,
		"":      zap.InfoLevel,
		"info":  zap.InfoLevel,
		"warn":  zap.WarnLevel,
		"error": zap.ErrorLevel,
	}
)

type logger struct {
	*zap.SugaredLogger
}

func New() LoggerI {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.DisableStacktrace = true
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.CallerKey = ""
	cfg.Level.SetLevel(levels[config.Vars.ApiLogLevel])
	logger, err := cfg.Build()
	if err != nil {
		log.Fatal(err.Error())
	}
	return NewWithZap(logger)
}

func NewWithZap(l *zap.Logger) LoggerI {
	return &logger{l.Sugar()}
}

func NewForTest() (LoggerI, *observer.ObservedLogs) {
	core, recorded := observer.New(zapcore.InfoLevel)
	return NewWithZap(zap.New(core)), recorded
}

func (l *logger) With(ctx context.Context, args ...interface{}) LoggerI {
	if id := ctxpkg.GetValue(ctx, ctxpkg.KeyRequestID); id != "" {
		args = append(args, zap.String("requestID", id))
	}
	if id := ctxpkg.GetValue(ctx, ctxpkg.KeyCorrelationID); id != "" {
		args = append(args, zap.String("correlationID", id))
	}
	if len(args) > 0 {
		return &logger{l.SugaredLogger.With(args...)}
	}
	return l
}
