package logger

import (
	"goredis/common/config"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zaplogger struct {
	zl *zap.Logger
}

var _ Logger = (*zaplogger)(nil)

func NewLogger(cfg *config.Config) (*zaplogger, error) {
	return &zaplogger{
		zl: buildZapLogger(cfg),
	}, nil
}

func (l *zaplogger) Debug(msg string, fields ...zap.Field) {
	l.zl.Debug(msg, fields...)
}

func (l *zaplogger) Warn(msg string, fields ...zap.Field) {
	l.zl.Warn(msg, fields...)
}

func (l *zaplogger) Info(msg string, fields ...zap.Field) {
	l.zl.Info(msg, fields...)
}

func (l *zaplogger) Error(err error, fields ...zap.Field) {
	l.zl.Error(err.Error(), fields...)
}

func buildZapLogger(cfg *config.Config) *zap.Logger {

	encodeConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      zapcore.OmitKey,
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(parseZapLevel(cfg.Loginfo.Level)),
		Encoding:         "json",
		EncoderConfig:    encodeConfig,
		Sampling:         nil,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	logger, _ := config.Build()
	return logger
}

func parseZapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "info":
		return zap.InfoLevel
	case "debug":
		return zap.DebugLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	case "dpanic":
		return zap.DPanicLevel
	default:
		return zap.InfoLevel
	}
}
