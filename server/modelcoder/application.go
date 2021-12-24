package modelcoder

import (
	"fmt"
	"reflect"
	"strings"
)

type Application struct {
	context            *applicationContext
	factory            FactoryScript
	factoryScriptCache map[string]interface{}
	OutDebug           func() bool
	OutInfo            func() bool
	Debug              func(args ...interface{})
	Info               func(args ...interface{})
	Warn               func(args ...interface{})
	Error              func(args ...interface{})
}

func NewApplication(applicationModel *ApplicationModel, options ...interface{}) *Application {
	if applicationModel == nil {
		return nil
	}
	res := &Application{}
	res.initOption(&LoggerDefault{})
	res.initOption(&FactoryScriptDefault{})

	if len(options) > 0 {
		for _, option := range options {
			res.initOption(option)
		}
	}

	res.context = newApplicationContext(applicationModel)
	return res
}

type LoggerOption interface {
	OutDebug() bool
	OutInfo() bool
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

func (this_ *Application) initOption(option interface{}) *Application {
	if option == nil {
		return this_
	}
	LoggerOption, LoggerOptionOk := option.(LoggerOption)
	if LoggerOptionOk {
		this_.initLoggerOption(LoggerOption)
	}
	FactoryScript, FactoryScriptOk := option.(FactoryScript)
	if FactoryScriptOk {
		this_.initFactoryOption(FactoryScript)
	}
	return this_
}

func (this_ *Application) initLoggerOption(option LoggerOption) *Application {
	if option == nil {
		return this_
	}
	this_.OutDebug = option.OutDebug
	this_.OutInfo = option.OutInfo
	this_.Debug = option.Debug
	this_.Info = option.Info
	this_.Warn = option.Warn
	this_.Error = option.Error
	return this_
}

func (this_ *Application) initFactoryOption(option FactoryScript) *Application {
	if option == nil {
		return this_
	}

	this_.factoryScriptCache = make(map[string]interface{})
	this_.factory = option

	reflectType := reflect.TypeOf(this_.factory)
	count := reflectType.NumMethod()
	var i = 0
	for i = 0; i < count; i++ {
		method := reflectType.Method(i)
		this_.factoryScriptCache[method.Name] = true
		this_.factoryScriptCache[strings.ToLower(method.Name[0:1])+method.Name[1:]] = true
	}
	return this_
}

func (this_ *Application) NewInvokeVariable(VariableDatas ...*VariableData) *invokeVariable {
	res := &invokeVariable{
		VariableDatas: []*VariableData{},
		Parent:        nil,
		InvokeCache:   map[string]interface{}{},
	}
	if len(VariableDatas) > 0 {
		res.AddVariableData(VariableDatas...)
	}
	return res
}

func (this_ *Application) executeSqlInsert(database string, sql string, sqlParams []interface{}) (err error) {
	if this_.OutDebug() {
		this_.Debug("execute sql insert sql   :", sql)
		this_.Debug("execute sql insert params:", ToJSON(sqlParams))
	}

	return
}

type LoggerDefault struct {
}

func (this_ *LoggerDefault) OutDebug() bool {
	return true
}
func (this_ *LoggerDefault) OutInfo() bool {
	return true
}
func (this_ *LoggerDefault) Debug(args ...interface{}) {
	fmt.Println(args...)
}
func (this_ *LoggerDefault) Info(args ...interface{}) {
	fmt.Println(args...)
}
func (this_ *LoggerDefault) Warn(args ...interface{}) {
	fmt.Println(args...)
}
func (this_ *LoggerDefault) Error(args ...interface{}) {
	fmt.Println(args...)
}
