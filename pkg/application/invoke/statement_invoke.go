package invoke

import (
	"errors"
	"fmt"
	"reflect"
	"teamide/pkg/application/base"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/parser"
)

func (this_ *StatementParser) check(invokeInfo *InvokeInfo) (err error) {

	return
}

func (this_ *StatementParser) invoke(invokeInfo *InvokeInfo) (res interface{}, err error) {
	if this_.statements == nil {
		var program *ast.Program
		program, err = parser.ParseFile(nil, "", this_.script, 0)
		if err != nil {
			return
		}
		if len(program.Body) == 0 {
			err = errors.New("please enter the correct script")
			return
		}
		this_.statements = program.Body
	}
	err = this_.check(invokeInfo)
	if err != nil {
		return
	}
	res, err = this_.invokeStatements(this_.statements, invokeInfo)
	if err != nil {
		return
	}
	return
}

func (this_ *StatementParser) invokeStatements(statements []ast.Statement, invokeInfo *InvokeInfo) (res interface{}, err error) {
	for _, one := range statements {
		res, err = this_.invokeStatement(one, invokeInfo)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *StatementParser) invokeStatement(statement ast.Statement, invokeInfo *InvokeInfo) (res interface{}, err error) {
	if statement == nil {
		return
	}
	switch statement_ := statement.(type) {
	case *ast.BlockStatement:
		res, err = this_.invokeBlockStatement(statement_, invokeInfo)
	case *ast.VariableStatement:
		res, err = this_.invokeVariableStatement(statement_, invokeInfo)
	case *ast.ForStatement:
		res, err = this_.invokeForStatement(statement_, invokeInfo)
	case *ast.IfStatement:
		res, err = this_.invokeIfStatement(statement_, invokeInfo)
	case *ast.ThrowStatement:
		res, err = this_.invokeThrowStatement(statement_, invokeInfo)
	case *ast.ExpressionStatement:
		res, err = this_.invokeExpressionStatement(statement_, invokeInfo)
	case *ast.ReturnStatement:
		res, err = this_.invokeReturnStatement(statement_, invokeInfo)
	default:
		err = errors.New(fmt.Sprint("invokeStatement type not match:", reflect.TypeOf(statement).Elem().Name()))
	}

	if err != nil {
		return
	}
	return
}

func (this_ *StatementParser) invokeBlockStatement(statement *ast.BlockStatement, invokeInfo *InvokeInfo) (res interface{}, err error) {
	for _, one := range statement.List {
		res, err = this_.invokeStatement(one, invokeInfo)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *StatementParser) invokeReturnStatement(statement *ast.ReturnStatement, invokeInfo *InvokeInfo) (res interface{}, err error) {
	expressionParser := &ExpressionParser{
		expression: statement.Argument,
	}

	res, err = expressionParser.invoke(invokeInfo)
	if err != nil {
		return
	}
	return
}

func (this_ *StatementParser) invokeExpressionStatement(statement *ast.ExpressionStatement, invokeInfo *InvokeInfo) (res interface{}, err error) {
	expressionParser := &ExpressionParser{
		expression: statement.Expression,
	}

	res, err = expressionParser.invoke(invokeInfo)
	if err != nil {
		return
	}
	return
}

func (this_ *StatementParser) invokeVariableStatement(statement *ast.VariableStatement, invokeInfo *InvokeInfo) (res interface{}, err error) {
	for _, one := range statement.List {
		fmt.Println("invokeVariableStatement Target:", base.ToJSON(one.Target))
		fmt.Println("invokeVariableStatement Initializer:", base.ToJSON(one.Initializer))
	}
	return
}

func (this_ *StatementParser) invokeForStatement(statement *ast.ForStatement, invokeInfo *InvokeInfo) (res interface{}, err error) {
	fmt.Println("invokeForStatement For:", base.ToJSON(statement.For))
	fmt.Println("invokeForStatement Initializer:", base.ToJSON(statement.Initializer))
	fmt.Println("invokeForStatement Test:", base.ToJSON(statement.Test))
	fmt.Println("invokeForStatement Update:", base.ToJSON(statement.Update))
	fmt.Println("invokeForStatement Body:", base.ToJSON(statement.Body))
	return
}

func (this_ *StatementParser) invokeIfStatement(statement *ast.IfStatement, invokeInfo *InvokeInfo) (res interface{}, err error) {
	expressionParser := &ExpressionParser{
		expression: statement.Test,
	}
	var testRes interface{}
	testRes, err = expressionParser.invoke(invokeInfo)
	if err != nil {
		return
	}

	if invokeInfo.App.GetScript().IsTrue(testRes) {
		res, err = this_.invokeStatement(statement.Consequent, invokeInfo)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *StatementParser) invokeThrowStatement(statement *ast.ThrowStatement, invokeInfo *InvokeInfo) (res interface{}, err error) {
	expressionParser := &ExpressionParser{
		expression: statement.Argument,
	}
	res, err = expressionParser.invoke(invokeInfo)
	if err != nil {
		return
	}
	return
}
