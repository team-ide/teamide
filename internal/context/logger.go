package context

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"teamide/internal/config"
)

// newZapLogger creator a new zap logger
// hook {Filename, Maxsize(megabytes), MaxBackups, MaxAge(days)}
// level zap.Level { DebugLevel, InfoLevel, WarnLevel, ErrorLevel, }
func newZapLogger(serverConfig *config.ServerConfig) *zap.Logger {
	var hook = &lumberjack.Logger{
		Filename:   serverConfig.Log.Filename,
		MaxSize:    serverConfig.Log.MaxSize,
		MaxAge:     serverConfig.Log.MaxAge,
		MaxBackups: serverConfig.Log.MaxBackups,
		Compress:   false,
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
	switch serverConfig.Log.Level {
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

	return zap.New(core, caller)
}
