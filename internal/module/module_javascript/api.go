package module_javascript

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/javascript"
	"github.com/team-ide/go-tool/javascript/context_map"
	_ "github.com/team-ide/go-tool/javascript/context_service"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"strings"
	"sync"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
	"time"
)

type api struct {
	toolboxService *module_toolbox.ToolboxService
}

func NewApi(toolboxService *module_toolbox.ToolboxService) *api {
	return &api{
		toolboxService: toolboxService,
	}
}

var (
	Power      = base.AppendPower(&base.PowerAction{Action: "javascript", Text: "Javascript", ShouldLogin: true, StandAlone: true})
	getModules = base.AppendPower(&base.PowerAction{Action: "getModules", Text: "getModules", ShouldLogin: true, StandAlone: true, Parent: Power})
	run        = base.AppendPower(&base.PowerAction{Action: "run", Text: "Run", ShouldLogin: true, StandAlone: true, Parent: Power})
	load       = base.AppendPower(&base.PowerAction{Action: "load", Text: "Load", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: getModules, Do: this_.getModules})
	apis = append(apis, &base.ApiWorker{Power: run, Do: this_.run})
	apis = append(apis, &base.ApiWorker{Power: load, Do: this_.load})

	return
}

type BaseRequest struct {
	Javascript string `json:"javascript"`
	WorkerId   string `json:"workerId,omitempty"`
}

func (this_ *api) getModules(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	res = context_map.ModuleList
	return
}

func (this_ *api) run(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	run := &JavascriptRun{
		BaseRequest: request,
	}
	setJavascriptRun(request.WorkerId, run)

	go func() {
		run.run()
	}()

	return
}
func (this_ *api) load(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	find := getJavascriptRun(request.WorkerId)
	if find != nil {
		find.Log = find.buffer.String()
	}
	res = find
	return
}

type JavascriptRun struct {
	*BaseRequest
	Key       string      `json:"key"`
	IsEnd     bool        `json:"isEnd"`
	StartTime time.Time   `json:"startTime"`
	Start     int64       `json:"start"`
	EndTime   time.Time   `json:"endTime"`
	End       int64       `json:"end"`
	Error     string      `json:"error"`
	Result    interface{} `json:"result"`
	Log       string      `json:"log"`
	buffer    *bytes.Buffer
}

func (this_ *JavascriptRun) run() {
	this_.buffer = &bytes.Buffer{}
	logger := newLogger(this_.buffer)
	logger.Debug("javascript run start")
	var err error
	this_.StartTime = time.Now()
	this_.Start = this_.StartTime.UnixMilli()
	defer func() {
		logger.Debug("javascript run end")
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			logger.Error("javascript run panic", zap.Error(err))
		}
		if err != nil {
			this_.Error = err.Error()
		}
		this_.EndTime = time.Now()
		this_.End = this_.EndTime.UnixMilli()
		this_.IsEnd = true
	}()

	scriptContext := javascript.NewContext()

	scriptContext["logger"] = map[string]interface{}{
		"debug": func(args ...interface{}) {
			msg, fields := formatZapArgs(args...)
			logger.Debug(fmt.Sprint(msg...), fields...)
		},
		"info": func(args ...interface{}) {
			msg, fields := formatZapArgs(args...)
			logger.Info(fmt.Sprint(msg...), fields...)
		},
		"warn": func(args ...interface{}) {
			msg, fields := formatZapArgs(args...)
			logger.Warn(fmt.Sprint(msg...), fields...)
		},
		"error": func(args ...interface{}) {
			msg, fields := formatZapArgs(args...)
			logger.Error(fmt.Sprint(msg...), fields...)
		},
		"any": func(key, value interface{}) zap.Field {
			return zap.Any(util.GetStringValue(key), value)
		},
	}

	this_.Result, err = javascript.RunScript(this_.Javascript, scriptContext)
	if err != nil {
		logger.Error("javascript run error", zap.Error(err))
		return
	}

}

var (
	javascriptRunCache     = map[string]*JavascriptRun{}
	javascriptRunCacheLock = &sync.Mutex{}
)

func getJavascriptRun(key string) (res *JavascriptRun) {
	javascriptRunCacheLock.Lock()
	defer javascriptRunCacheLock.Unlock()

	res = javascriptRunCache[key]

	return
}

func setJavascriptRun(key string, value *JavascriptRun) {
	javascriptRunCacheLock.Lock()
	defer javascriptRunCacheLock.Unlock()

	javascriptRunCache[key] = value

	return
}

func removeJavascriptRun(key string) {
	javascriptRunCacheLock.Lock()
	defer javascriptRunCacheLock.Unlock()

	delete(javascriptRunCache, key)

	return
}

func newLogger(w io.Writer) *zap.Logger {

	// 默认 日志 输出在 控制台
	// 日志格式如下
	// [2023-05-17 11:17:33.063] [DEBUG] [util/logger.go:55] [logger test debug] [{"arg1": 1, "arg2": "2"}]
	// [2023-05-17 11:17:33.067] [INFO ] [util/logger.go:56] [logger test info] [{"arg1": 1, "arg2": "2"}]
	// [2023-05-17 11:17:33.067] [WARN ] [util/logger.go:57] [logger test warn] [{"arg1": 1, "arg2": "2"}]
	// [2023-05-17 11:17:33.067] [ERROR] [util/logger.go:58] [logger test error] [{"arg1": 1, "arg2": "2"}]

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
			enc.AppendString(caller.TrimmedPath())
		},
		EncodeName: func(s string, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(strings.ToUpper(s))
		},
	}
	level := zapcore.DebugLevel
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(w),
		zap.NewAtomicLevelAt(level),
	)
	logger := zap.New(
		core,
		// 表示 输出 文件名 以及 行号
		zap.AddCaller(),
		// 表示 输出 堆栈跟踪 传入 level 表示 在哪个级别下输出
		zap.AddStacktrace(zapcore.ErrorLevel),
		//zap.AddCallerSkip(0),
	)
	//Logger.Debug("logger test debug", zap.Any("arg1", 1), zap.Any("arg2", "2"))
	//Logger.Info("logger test info", zap.Any("arg1", 1), zap.Any("arg2", "2"))
	//Logger.Warn("logger test warn", zap.Any("arg1", 1), zap.Any("arg2", "2"))
	//Logger.Error("logger test error", zap.Any("arg1", 1), zap.Any("arg2", "2"))

	return logger
}

func formatZapArgs(args ...interface{}) (msg []interface{}, fields []zap.Field) {
	for _, arg := range args {
		switch tV := arg.(type) {
		case zap.Field:
			fields = append(fields, tV)
		default:
			msg = append(msg, tV)
		}
	}

	return
}
