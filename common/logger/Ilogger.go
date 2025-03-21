package logger

import "go.uber.org/zap"

type (
	Logger interface {
		Debug(msg string, fields ...zap.Field)
		Info(msg string, fields ...zap.Field)
		Warn(msg string, fields ...zap.Field)
		Error(err error, fields ...zap.Field)
	}
)
