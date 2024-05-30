package golang

import (
	"sort"
	"strings"
)

var (
	loggerCode = `
var (
	Logger = newConsoleLogger()
)

func init() {
	util.Logger = Logger
}

func Init(c *config.Log) {
	Logger = NewZapLogger(c)
	util.Logger = Logger
	return
}

// NewZapLogger creator a new zap logger
// hook {Filename, Maxsize(megabytes), MaxBackups, MaxAge(days)}
// level zap.Level { DebugLevel, InfoLevel, WarnLevel, ErrorLevel, }
func NewZapLogger(c *config.Log) *zap.Logger {
	var writer io.Writer
	var cLevel string
	if c != nil {
		cLevel = c.Level
	}
	if c != nil && !c.Console {
		writer = &lumberjack.Logger{
			Filename:   c.Filename,
			MaxSize:    c.MaxSize,
			MaxAge:     c.MaxAge,
			MaxBackups: c.MaxBackups,
			Compress:   true,
		}
	} else {
		writer = os.Stdout
	}
	var level zapcore.Level
	switch cLevel {
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
		zapcore.AddSync(writer),
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
`
)

func (this_ *Generator) GenLogger() (err error) {
	dir := this_.golang.GetLoggerDir(this_.Dir)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "logger.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	var imports []string

	imports = append(imports, this_.golang.GetConfigImport())
	pack := this_.golang.GetLoggerPack()

	builder.AppendTabLine("package " + pack)
	builder.NewLine()

	builder.AppendTabLine("import(")
	builder.Tab()

	ss := strings.Split(`
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strings"
	"time"
`, "\n")
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		s = strings.TrimPrefix(s, `"`)
		s = strings.TrimSuffix(s, `"`)
		imports = append(imports, s)
	}

	sort.Strings(imports)
	for _, im := range imports {
		builder.AppendTabLine("\"" + im + "\"")
	}
	builder.Indent()
	builder.AppendTabLine(")")
	builder.NewLine()

	code := strings.ReplaceAll(loggerCode, "{pack}", pack)

	builder.AppendCode(code)
	return
}
