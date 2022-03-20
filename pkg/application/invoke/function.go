package invoke

import (
	"sync"
	"teamide/pkg/application/base"
	common2 "teamide/pkg/application/common"

	"github.com/dop251/goja/ast"
)

func ExecuteFunctionScript(app common2.IApplication, invokeNamespace *common2.InvokeNamespace, script string) (res interface{}, err error) {
	if app.GetScript().IsEmpty(script) {
		res = nil
	} else {
		var parser *FunctionParser = getFunctionParser(script)
		invokeInfo := &InvokeInfo{App: app, InvokeNamespace: invokeNamespace}
		res, err = parser.invoke(invokeInfo)
	}
	return
}

var (
	FunctionParserCache = make(map[string]*FunctionParser)
)

var (
	FunctionParserMutex sync.Mutex
)

func getFunctionParser(script string) (res *FunctionParser) {
	if base.IsEmpty(script) {
		return
	}
	FunctionParserMutex.Lock()
	defer FunctionParserMutex.Unlock()
	var ok bool
	res, ok = FunctionParserCache[script]
	if !ok {
		res = &FunctionParser{
			script: script,
		}
		FunctionParserCache[script] = res
	}
	return
}

type FunctionParser struct {
	script   string
	function *ast.FunctionLiteral
	Name     string
}

func NewFunctionParser(script string) *FunctionParser {
	return &FunctionParser{script: script}
}
