package util

import "go.uber.org/zap"

var (
	Logger *zap.Logger
)

func init() {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.Development = false
	Logger, _ = loggerConfig.Build()
}
