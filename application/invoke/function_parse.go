package invoke

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/parser"
)

func (this_ *FunctionParser) Parse(parseInfo *ParseInfo) (err error) {
	if this_.function == nil {
		var program *ast.Program
		program, err = parser.ParseFile(nil, "", this_.script, 0)
		if err != nil {
			fmt.Println("error script:")
			fmt.Println(this_.script)
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
	err = this_.parseFunction(this_.function, parseInfo)
	if err != nil {
		return
	}
	return
}

func (this_ *FunctionParser) parseFunction(function *ast.FunctionLiteral, parseInfo *ParseInfo) (err error) {
	this_.Name = function.Name.Name.String()

	err = this_.parseParameterList(function.ParameterList, parseInfo)
	if err != nil {
		return
	}
	statementParser := &StatementParser{
		statements: []ast.Statement{function.Body},
	}
	err = statementParser.Parse(parseInfo)
	if err != nil {
		return
	}
	return
}

func (this_ *FunctionParser) parseParameterList(parameterList *ast.ParameterList, parseInfo *ParseInfo) (err error) {
	for _, one := range parameterList.List {
		identifier, ok := one.Target.(*ast.Identifier)
		if !ok {
			err = errors.New(fmt.Sprint("parseParameterList type not match:", reflect.TypeOf(one.Target).Elem().Name()))
			return
		}
		parseInfo.ParameterList = append(parseInfo.ParameterList, identifier.Name.String())
	}
	return
}
