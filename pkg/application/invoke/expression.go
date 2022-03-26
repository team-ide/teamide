package invoke

import (
	"sync"
	"teamide/pkg/application/base"
	"teamide/pkg/application/common"

	"github.com/dop251/goja/ast"
)

func ExecuteExpressionScript(app common.IApplication, invokeNamespace *common.InvokeNamespace, script string) (res interface{}, err error) {
	if app.GetScript().IsEmpty(script) {
		res = nil
	} else {
		var parser *ExpressionParser = getExpressionParser(script)
		invokeInfo := &InvokeInfo{App: app, InvokeNamespace: invokeNamespace}
		res, err = parser.invoke(invokeInfo)
	}
	return
}

var (
	ExpressionParserCache = make(map[string]*ExpressionParser)
)

var (
	ExpressionParserMutex sync.Mutex
)

func getExpressionParser(script string) (res *ExpressionParser) {
	if base.IsEmpty(script) {
		return
	}
	ExpressionParserMutex.Lock()
	defer ExpressionParserMutex.Unlock()
	var ok bool
	res, ok = ExpressionParserCache[script]
	if !ok {
		res = &ExpressionParser{
			script: script,
		}
		ExpressionParserCache[script] = res
	}
	return
}

type ExpressionParser struct {
	script     string
	expression ast.Expression
}
