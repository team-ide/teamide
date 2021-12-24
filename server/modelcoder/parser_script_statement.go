package modelcoder

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/parser"
)

type scriptStatementParser struct {
	script      string
	application *Application
}

func (this_ *scriptStatementParser) parse() error {

	program, err := parser.ParseFile(nil, "", this_.script, 0)
	if err != nil {
		return err
	}
	if len(program.Body) == 0 {
		err = errors.New("please enter the correct script")
		return err
	}
	err = this_.parseStatements(program.Body)
	return err
}

func (this_ *scriptStatementParser) parseStatements(statements []ast.Statement) (err error) {
	for _, one := range statements {
		err = this_.parseStatement(one)
		if err != nil {
			return err
		}
	}
	return err
}

func (this_ *scriptStatementParser) parseStatement(statement ast.Statement) (err error) {
	var ok = false

	if !ok {
		statement_, ok_ := interface{}(statement).(*ast.VariableStatement)
		ok = ok_
		if ok_ {
			err = this_.parseVariableStatement(statement_)
		}
	}
	if !ok {
		statement_, ok_ := interface{}(statement).(*ast.ForStatement)
		ok = ok_
		if ok_ {
			err = this_.parseForStatement(statement_)
		}
	}

	if err != nil {
		return err
	}
	if !ok {
		err = errors.New(fmt.Sprint("parseStatement type not match:", reflect.TypeOf(statement).Elem().Name()))
	}
	return err
}

func (this_ *scriptStatementParser) parseVariableStatement(statement *ast.VariableStatement) (err error) {
	for _, one := range statement.List {
		fmt.Println("parseVariableStatement Target:", ToJSON(one.Target))
		fmt.Println("parseVariableStatement Initializer:", ToJSON(one.Initializer))
	}
	return
}

func (this_ *scriptStatementParser) parseForStatement(statement *ast.ForStatement) (err error) {
	fmt.Println("parseForStatement For:", ToJSON(statement.For))
	fmt.Println("parseForStatement Initializer:", ToJSON(statement.Initializer))
	fmt.Println("parseForStatement Test:", ToJSON(statement.Test))
	fmt.Println("parseForStatement Update:", ToJSON(statement.Update))
	fmt.Println("parseForStatement Body:", ToJSON(statement.Body))
	return
}
