package common

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

type ILogger interface {
	OutDebug() bool
	OutInfo() bool
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type LoggerDefault struct {
	OutDebug_ bool
	OutInfo_  bool
}

func (this_ *LoggerDefault) OutDebug() bool {
	return this_.OutDebug_
}
func (this_ *LoggerDefault) OutInfo() bool {
	return this_.OutInfo_
}
func (this_ *LoggerDefault) Debug(args ...interface{}) {
	info := GetParentInfo(2)
	info["level"] = "debug"
	info["msg"] = fmt.Sprint(args...)
	fmt.Println(getLogLine(info))
}
func (this_ *LoggerDefault) Info(args ...interface{}) {
	info := GetParentInfo(2)
	info["level"] = "info"
	info["msg"] = fmt.Sprint(args...)
	fmt.Println(getLogLine(info))
}
func (this_ *LoggerDefault) Warn(args ...interface{}) {
	info := GetParentInfo(2)
	info["level"] = "warn"
	info["msg"] = fmt.Sprint(args...)
	fmt.Println(getLogLine(info))
}
func (this_ *LoggerDefault) Error(args ...interface{}) {
	info := GetParentInfo(2)
	info["level"] = "error"
	info["msg"] = fmt.Sprint(args...)
	fmt.Println(getLogLine(info))
}

func getLogLine(info map[string]interface{}) string {
	if len(info["level"].(string)) == 4 {
		info["level"] = info["level"].(string) + " "
	}

	return fmt.Sprint("[:", info["time"], ":]", " [:", info["level"], ":]", " [:", info["funcName"], " ", info["line"], ":]", " [:", info["msg"], ":]")
}

func GetParentInfo(index int) (res map[string]interface{}) {
	res = map[string]interface{}{}
	var pc uintptr
	var file string
	var line int
	var ok bool
	pc, file, line, ok = runtime.Caller(index)
	runtime.FuncForPC(pc).Entry()
	funcName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(file) // Base函数返回路径的最后一个元素
	funcName = strings.TrimPrefix(funcName, "application.")
	funcName = strings.TrimPrefix(funcName, "teamide/application/")
	res["funcName"] = funcName
	res["file"] = file
	res["fileName"] = fileName
	res["line"] = line
	res["line"] = line
	res["ok"] = ok
	res["time"] = time.Now().Format("2006-01-02 15:04:05.999999999")[0:23]
	var fun *runtime.Func = runtime.FuncForPC(pc)
	res["method"] = fun.Name()
	return
}
