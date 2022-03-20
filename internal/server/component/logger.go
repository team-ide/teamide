package component

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"teamide/internal/config"
)

var (
	errParam = errors.New("param is error")
	Logger   = initLogger()
)

func LogStr(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	return fmt.Sprint(args...)
}

// newZapLogger creator a new zap logger
// hook {Filename, Maxsize(megabytes), MaxBackups, MaxAge(days)}
// level zap.Level { DebugLevel, InfoLevel, WarnLevel, ErrorLevel, }
func newZapLogger() (*zap.Logger, error) {

	var hook *lumberjack.Logger = &lumberjack.Logger{
		Filename:   config.Config.Log.Filename,
		MaxSize:    config.Config.Log.MaxSize,
		MaxAge:     config.Config.Log.MaxAge,
		MaxBackups: config.Config.Log.MaxBackups,
		Compress:   false,
	}

	if hook == nil {
		return nil, errParam
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	var level zapcore.Level
	switch config.Config.Log.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.DebugLevel
	}
	atomicLevel := zap.NewAtomicLevelAt(level)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(hook), atomicLevel)
	caller := zap.AddCaller()

	return zap.New(core, caller), nil
}

func initLogger() *zap.Logger {
	logger, _ := newZapLogger()
	logger.Info(LogStr("日志初始化完成,filename:", config.Config.Log.Filename, ",level:", config.Config.Log.Level))
	return logger
}
