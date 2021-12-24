package modelcoder

func invokeScriptStatement(application *Application, variable *invokeVariable, script string) (err error) {
	if IsEmpty(script) {
	} else {
		var parser *scriptStatementParser
		parser, err = getScriptStatementParser(application, script)
		if err != nil {
			return
		}
		if parser != nil {
			err = parser.invoke(variable)
		}
	}
	return
}

var (
	scriptStatementParserCache = make(map[string]*scriptStatementParser)
)

func getScriptStatementParser(application *Application, script string) (res *scriptStatementParser, err error) {
	if IsEmpty(script) {
		return
	}
	var ok bool
	res, ok = scriptStatementParserCache[script]
	if !ok {
		res = &scriptStatementParser{
			script:      script,
			application: application,
		}
		err = res.parse()
		if err != nil {
			return nil, err
		}
		scriptStatementParserCache[script] = res
	}
	return
}

func (this_ *scriptStatementParser) check(variable *invokeVariable) (err error) {

	return
}

func (this_ *scriptStatementParser) invoke(variable *invokeVariable) (err error) {
	err = this_.check(variable)
	if err != nil {
		return
	}
	return
}
