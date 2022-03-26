package invoke

import (
	"sync"
	"teamide/pkg/application/base"
	"teamide/pkg/application/common"

	"github.com/dop251/goja/ast"
)

func ExecuteStatementScript(app common.IApplication, invokeNamespace *common.InvokeNamespace, script string) (res interface{}, err error) {
	if app.GetScript().IsEmpty(script) {
	} else {
		var parser *StatementParser = getStatementParser(script)
		invokeInfo := &InvokeInfo{App: app, InvokeNamespace: invokeNamespace}
		res, err = parser.invoke(invokeInfo)
	}
	return
}

var (
	StatementParserCache = make(map[string]*StatementParser)
)

var (
	StatementParserMutex sync.Mutex
)

func getStatementParser(script string) (res *StatementParser) {
	if base.IsEmpty(script) {
		return
	}
	StatementParserMutex.Lock()
	defer StatementParserMutex.Unlock()
	var ok bool
	res, ok = StatementParserCache[script]
	if !ok {
		res = &StatementParser{
			script: script,
		}
		StatementParserCache[script] = res
	}
	return
}

type StatementParser struct {
	script     string
	statements []ast.Statement
}
