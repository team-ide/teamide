package context

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
	"teamide/internal/config"
	"time"
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
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "S",
		//FunctionKey:      "F",
		ConsoleSeparator: "] [",
		LineEnding:       "]\n",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000"))
		},
		EncodeLevel: func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(util.StrPadRight(strings.ToUpper(l.String()), 5, " "))
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendFloat64(float64(d) / float64(time.Second))
		},
		EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			str := caller.TrimmedPath()
			method := caller.Function
			dot := strings.LastIndex(method, ".")
			if dot > 0 {
				method = method[dot+1:]
				index := strings.LastIndex(str, ":")
				if index > 0 {
					str = str[0:index] + ":" + method + str[index:]
				} else {
					str += ":" + method
				}
			}
			enc.AppendString(str)
		},
		EncodeName: func(s string, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(strings.ToUpper(s))
		},
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
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(hook),
		zap.NewAtomicLevelAt(level),
	)
	res := zap.New(
		core,
		// 表示 输出 文件名 以及 行号
		zap.AddCaller(),
		// 表示 输出 堆栈跟踪 传入 level 表示 在哪个级别下输出
		zap.AddStacktrace(zapcore.ErrorLevel),
		//zap.AddCallerSkip(0),
	)

	return res
}
