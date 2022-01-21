package invoke

import (
	"errors"
	"fmt"
	"reflect"
	"teamide/application/base"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/parser"
)

func (this_ *StatementParser) Parse(parseInfo *ParseInfo) (err error) {

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
	err = this_.parseStatements(this_.statements, parseInfo)
	return err
}

func (this_ *StatementParser) parseStatements(statements []ast.Statement, parseInfo *ParseInfo) (err error) {
	for _, one := range statements {
		err = this_.parseStatement(one, parseInfo)
		if err != nil {
			return err
		}
	}
	return err
}

func (this_ *StatementParser) parseStatement(statement ast.Statement, parseInfo *ParseInfo) (err error) {
	if statement == nil {
		return
	}
	switch statement_ := statement.(type) {
	case *ast.BlockStatement:
		err = this_.parseBlockStatement(statement_, parseInfo)
	case *ast.VariableStatement:
		err = this_.parseVariableStatement(statement_, parseInfo)
	case *ast.ForStatement:
		err = this_.parseForStatement(statement_, parseInfo)
	case *ast.IfStatement:
		err = this_.parseIfStatement(statement_, parseInfo)
	case *ast.ThrowStatement:
		err = this_.parseThrowStatement(statement_, parseInfo)
	case *ast.ExpressionStatement:
		err = this_.parseExpressionStatement(statement_, parseInfo)
	case *ast.ReturnStatement:
		err = this_.parseReturnStatement(statement_, parseInfo)
	default:
		err = errors.New(fmt.Sprint("parseStatement type not match:", reflect.TypeOf(statement).Elem().Name()))
	}

	if err != nil {
		return err
	}
	return err
}

func (this_ *StatementParser) parseBlockStatement(statement *ast.BlockStatement, parseInfo *ParseInfo) (err error) {
	for _, one := range statement.List {
		err = this_.parseStatement(one, parseInfo)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *StatementParser) parseReturnStatement(statement *ast.ReturnStatement, parseInfo *ParseInfo) (err error) {
	expressionParser := &ExpressionParser{
		expression: statement.Argument,
	}

	err = expressionParser.Parse(parseInfo)
	if err != nil {
		return
	}
	return
}

func (this_ *StatementParser) parseExpressionStatement(statement *ast.ExpressionStatement, parseInfo *ParseInfo) (err error) {
	expressionParser := &ExpressionParser{
		expression: statement.Expression,
	}

	err = expressionParser.Parse(parseInfo)
	if err != nil {
		return
	}
	return
}

func (this_ *StatementParser) parseVariableStatement(statement *ast.VariableStatement, parseInfo *ParseInfo) (err error) {
	for _, one := range statement.List {
		fmt.Println("parseVariableStatement Target:", base.ToJSON(one.Target))
		fmt.Println("parseVariableStatement Initializer:", base.ToJSON(one.Initializer))
	}
	return
}

func (this_ *StatementParser) parseForStatement(statement *ast.ForStatement, parseInfo *ParseInfo) (err error) {
	fmt.Println("parseForStatement For:", base.ToJSON(statement.For))
	fmt.Println("parseForStatement Initializer:", base.ToJSON(statement.Initializer))
	fmt.Println("parseForStatement Test:", base.ToJSON(statement.Test))
	fmt.Println("parseForStatement Update:", base.ToJSON(statement.Update))
	fmt.Println("parseForStatement Body:", base.ToJSON(statement.Body))
	return
}

func (this_ *StatementParser) parseIfStatement(statement *ast.IfStatement, parseInfo *ParseInfo) (err error) {
	expressionParser := &ExpressionParser{
		expression: statement.Test,
	}

	err = expressionParser.Parse(parseInfo)
	if err != nil {
		return
	}
	err = this_.parseStatement(statement.Alternate, parseInfo)
	if err != nil {
		return
	}
	err = this_.parseStatement(statement.Consequent, parseInfo)
	if err != nil {
		return
	}

	return
}

func (this_ *StatementParser) parseThrowStatement(statement *ast.ThrowStatement, parseInfo *ParseInfo) (err error) {
	expressionParser := &ExpressionParser{
		expression: statement.Argument,
	}
	err = expressionParser.Parse(parseInfo)
	if err != nil {
		return
	}
	return
}
