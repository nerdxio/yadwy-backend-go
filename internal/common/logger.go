package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = "stacktrace"

	env := "development"
	if env != "production" {
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return config.Build()
}
