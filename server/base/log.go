package base

import (
	"errors"
	"fmt"
	"server/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	errParam = errors.New("param is error")
	Logger   *zap.Logger
)

func LogStr(args ...interface{}) string {
	if args == nil || len(args) == 0 {
		return ""
	}
	return fmt.Sprint(args...)
}

// NewZapLogger creator a new zap logger
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
		EncodeTime:     zapcore.ISO8601TimeEncoder,
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

func init() {
	Logger, _ = newZapLogger()
	Logger.Info(LogStr("vrv job ", "log inited ", "logger config:", config.Config.Log))
}
