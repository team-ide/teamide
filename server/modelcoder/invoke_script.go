package modelcoder

func getColumnValue(application *Application, variable *invokeVariable, name string, valueScript string) (res interface{}, err error) {
	if IsEmpty(valueScript) {
		paramData := variable.GetParamData(name)
		if paramData != nil {
			res = paramData.Data
		}
	} else {
		res, err = getScriptValue(application, variable, valueScript)
	}
	return
}

func ifScriptValue(application *Application, variable *invokeVariable, ifScript string) (res bool, err error) {
	if IsEmpty(ifScript) {
		res = true
		return
	}
	var value interface{}
	value, err = getScriptValue(application, variable, ifScript)
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

func getScriptValue(application *Application, variable *invokeVariable, script string) (res interface{}, err error) {
	if IsEmpty(script) {
		res = nil
	} else {
		scriptParser := getScriptParser(application, script)
		if scriptParser != nil {
			res, err = scriptParser.invoke(variable)
		}
	}
	return
}

func getScriptParser(application *Application, script string) *scriptParser {
	if IsEmpty(script) {
		return nil
	}
	res, ok := application.scriptParserCache[script]
	if !ok {
		res = &scriptParser{
			script:             script,
			factory:            application.factory,
			factoryScriptCache: application.factoryScriptCache,
		}
		res.parse()
		application.scriptParserCache[script] = res
	}
	return res
}

func (this_ *scriptParser) invoke(variable *invokeVariable) (res interface{}, err error) {
	return
}
