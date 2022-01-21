package invoke

import (
	"errors"
	"fmt"
	"reflect"
	"teamide/application/common"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/parser"
)

func (this_ *FunctionParser) check(invokeInfo *InvokeInfo) (err error) {

	return
}
func (this_ *FunctionParser) invoke(invokeInfo *InvokeInfo) (res interface{}, err error) {

	if this_.function == nil {
		var program *ast.Program
		program, err = parser.ParseFile(nil, "", this_.script, 0)
		if err != nil {
			return
		}
		if len(program.Body) == 0 || len(program.Body) > 1 {
			err = errors.New("please enter the correct value expression")
			return
		}
		functionDeclaration, ok := interface{}(program.Body[0]).(*ast.FunctionDeclaration)
		if !ok {
			err = errors.New("please enter the correct value expression")
			return
		}
		this_.function = functionDeclaration.Function
	}
	err = this_.check(invokeInfo)
	if err != nil {
		return
	}
	res, err = this_.invokeFunction(this_.function, invokeInfo)
	if err != nil {
		return
	}
	return
}

func (this_ *FunctionParser) invokeFunction(function *ast.FunctionLiteral, invokeInfo *InvokeInfo) (res interface{}, err error) {
	err = this_.invokeParameterList(function.ParameterList, invokeInfo)
	if err != nil {
		return
	}
	statementParser := &StatementParser{
		statements: []ast.Statement{function.Body},
	}
	res, err = statementParser.invoke(invokeInfo)
	if err != nil {
		return
	}
	return
}

func (this_ *FunctionParser) invokeParameterList(parameterList *ast.ParameterList, invokeInfo *InvokeInfo) (err error) {
	for _, one := range parameterList.List {
		identifier, ok := one.Target.(*ast.Identifier)
		if !ok {
			err = errors.New(fmt.Sprint("invokeParameterList type not match:", reflect.TypeOf(one.Target).Elem().Name()))
			return
		}
		name := identifier.Name.String()
		var data *common.InvokeData
		data, err = invokeInfo.InvokeNamespace.GetData(name)
		if err != nil {
			return
		}
		if data == nil {
			// fmt.Println(fmt.Sprint("invokeParameterList parameter [", name, "] not defind"))
			// fmt.Println(base.ToJSON(invokeInfo.InvokeNamespace))
			err = errors.New(fmt.Sprint("invokeParameterList parameter [", name, "] not defind"))
			return
		}
	}
	return
}
