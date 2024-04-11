package context

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"teamide/internal/config"
	"teamide/pkg/node"
	"time"
)

var (
	Logger = newConsoleLogger()
)

func init() {
	util.Logger = Logger
	node.Logger = Logger
}

// newZapLogger creator a new zap logger
// hook {Filename, Maxsize(megabytes), MaxBackups, MaxAge(days)}
// level zap.Level { DebugLevel, InfoLevel, WarnLevel, ErrorLevel, }
func newZapLogger(serverConfig *config.ServerConfig) *zap.Logger {
	var hook = &lumberjack.Logger{
		Filename:   serverConfig.Log.Filename,
		MaxSize:    serverConfig.Log.MaxSize,
		MaxAge:     serverConfig.Log.MaxAge,
		MaxBackups: serverConfig.Log.MaxBackups,
		Compress:   true,
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
		zapcore.NewConsoleWrapColEncoder(NewEncoderConfig(), "[", "]"),
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

func newConsoleLogger() *zap.Logger {
	var level = zapcore.DebugLevel
	core := zapcore.NewCore(
		zapcore.NewConsoleWrapColEncoder(NewEncoderConfig(), "[", "]"),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(level),
	)
	caller := zap.AddCaller()

	return zap.New(core, caller,
		// 表示 输出 堆栈跟踪 传入 level 表示 在哪个级别下输出
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "S",
		//TraceKey:         "trackId",
		ConsoleSeparator: " ",
		LineEnding:       "\n",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		},
		EncodeLevel: func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			if l == zapcore.InfoLevel || l == zapcore.WarnLevel {
				enc.AppendString(l.CapitalString() + " ")
			} else {
				enc.AppendString(l.CapitalString())
			}
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
	}
}
