package toolbox

import (
	"go.uber.org/zap"
	"time"
)

//GetNowTime 获取当前时间戳
func GetNowTime() int64 {
	return time.Now().UnixNano() / 1e6
}

var (
	Logger *zap.Logger
)

func init() {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.Development = false
	Logger, _ = loggerConfig.Build()
}
