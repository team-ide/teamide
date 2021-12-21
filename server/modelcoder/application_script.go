package modelcoder

import (
	js "github.com/dop251/goja"
)

func getColumnValue(variable *invokeVariable, name string, valueScript string) (res interface{}, err error) {
	if IsEmpty(valueScript) {
		paramData := variable.GetParamData(name)
		if paramData != nil {
			res = paramData.Data
		}
	} else {
		res, err = getScriptValue(variable, valueScript)
	}
	return
}

func ifScriptValue(variable *invokeVariable, ifScript string) (res bool, err error) {
	if IsEmpty(ifScript) {
		res = true
		return
	}
	var value interface{}
	value, err = getScriptValue(variable, ifScript)
	if err != nil {

		return
	}
	if value == nil {
		return
	}
	if value == true || value == "1" || value == "true" {
		res = true
		return
	}
	return
}

func getScriptValue(variable *invokeVariable, script string) (res interface{}, err error) {
	if IsEmpty(script) {
		res = nil
	} else {
		processor := getScriptValueProcessor(script)
		if processor != nil {
			res, err = processor.process(variable)
		}
	}
	return
}

var (
	scriptValueProcessorCache = make(map[string]*scriptValueProcessor)
)

func getScriptValueProcessor(script string) *scriptValueProcessor {
	if IsEmpty(script) {
		return nil
	}
	processor, ok := scriptValueProcessorCache[script]
	if !ok {
		processor = &scriptValueProcessor{
			script: script,
		}
		processor.init()
		scriptValueProcessorCache[script] = processor
	}
	return processor
}

type scriptValueProcessor struct {
	script string
}

func (this_ *scriptValueProcessor) init() {
}

func (this_ *scriptValueProcessor) process(variable *invokeVariable) (res interface{}, err error) {

	vm := js.New() // 创建engine实例
	for _, paramData := range variable.ParamDatas {
		vm.Set(paramData.Name, paramData.Data)
	}
	r, err := vm.RunString(this_.script) // 执行javascript代码
	res = r.Export()                     // 将执行的结果转换为Golang对应的类型
	return
}
