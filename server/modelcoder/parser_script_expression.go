package modelcoder

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/parser"
)

type scriptExpressionParser struct {
	script         string
	application    *Application
	callDotNames   []string
	structDotNames []string
	expression     ast.Expression
}

func (this_ *scriptExpressionParser) parse() error {
	this_.callDotNames = []string{}
	this_.structDotNames = []string{}

	program, err := parser.ParseFile(nil, "", this_.script, 0)
	if err != nil {
		return err
	}
	if len(program.Body) == 0 || len(program.Body) > 1 {
		err = errors.New("please enter the correct value expression")
		return err
	}
	expression, ok := interface{}(program.Body[0]).(*ast.ExpressionStatement)
	if !ok {
		err = errors.New("please enter the correct value expression")
		return err
	}
	this_.expression = expression.Expression
	err = this_.parseExpression(this_.expression)
	if err != nil {
		return err
	}
	return nil
}

func (this_ *scriptExpressionParser) addCallNames(expression *ast.CallExpression) (err error) {
	name := ""

	calleeDotExpression, ok := interface{}(expression.Callee).(*ast.DotExpression)
	if ok {
		if calleeDotExpression.Left != nil && !isIdentifier(calleeDotExpression.Left) {
			return
		}
		name, err = getDotExpressionName(calleeDotExpression)
		if err != nil {
			return
		}
	} else {
		calleeIdentifier, ok := interface{}(expression.Callee).(*ast.Identifier)
		if ok {
			name = calleeIdentifier.Name.String()
		} else {
			err = errors.New(fmt.Sprint("addCallNames type not match:", reflect.TypeOf(expression.Callee).Elem().Name()))
		}
	}

	this_.callDotNames = append(this_.callDotNames, name)
	return
}

func (this_ *scriptExpressionParser) addStructDotNames(expression *ast.DotExpression) (err error) {
	if expression.Left != nil && !isIdentifier(expression.Left) {
		return
	}
	var name string
	name, err = getDotExpressionName(expression)
	if err != nil {
		return
	}
	this_.structDotNames = append(this_.structDotNames, name)
	return
}

func isIdentifier(expression ast.Expression) bool {
	_, ok := interface{}(expression).(*ast.Identifier)
	return ok

}
func getDotExpressionName(expression *ast.DotExpression) (name string, err error) {
	name = ""
	if expression.Left != nil {
		leftIdentifier, ok := interface{}(expression.Left).(*ast.Identifier)
		if ok {
			name = leftIdentifier.Name.String()
			name += "."
		} else {
			err = errors.New(fmt.Sprint("getDotExpressionName type not match:", reflect.TypeOf(expression.Left).Elem().Name()))
		}
	}
	name += expression.Identifier.Name.String()
	return
}

func (this_ *scriptExpressionParser) parseExpressions(expressions []ast.Expression) (err error) {
	for _, one := range expressions {
		err = this_.parseExpression(one)
		if err != nil {
			return err
		}
	}
	return err
}
func (this_ *scriptExpressionParser) parseExpression(expression ast.Expression) (err error) {
	var ok = false

	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.BinaryExpression)
		ok = ok_
		if ok_ {
			err = this_.parseBinaryExpression(expression_)
		}
	}
	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.ConditionalExpression)
		ok = ok_
		if ok_ {
			err = this_.parseConditionalExpression(expression_)
		}
	}
	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.CallExpression)
		ok = ok_
		if ok_ {
			err = this_.parseCallExpression(expression_)
		}
	}
	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.DotExpression)
		ok = ok_
		if ok_ {
			err = this_.parseDotExpression(expression_)
		}
	}
	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.Identifier)
		ok = ok_
		if ok_ {
			err = this_.parseIdentifier(expression_)
		}
	}
	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.StringLiteral)
		ok = ok_
		if ok_ {
			err = this_.parseStringLiteral(expression_)
		}
	}
	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.NumberLiteral)
		ok = ok_
		if ok_ {
			err = this_.parseNumberLiteral(expression_)
		}
	}

	if err != nil {
		return err
	}
	if !ok {
		err = errors.New(fmt.Sprint("parseExpression type not match:", reflect.TypeOf(expression).Elem().Name()))
	}
	return err
}

func (this_ *scriptExpressionParser) parseBinaryExpression(expression *ast.BinaryExpression) (err error) {

	err = this_.parseExpression(expression.Left)
	if err != nil {
		return err
	}
	err = this_.parseExpression(expression.Right)
	if err != nil {
		return err
	}
	return nil
}
func (this_ *scriptExpressionParser) parseConditionalExpression(expression *ast.ConditionalExpression) (err error) {

	err = this_.parseExpression(expression.Test)
	if err != nil {
		return err
	}
	err = this_.parseExpression(expression.Consequent)
	if err != nil {
		return err
	}
	err = this_.parseExpression(expression.Alternate)
	if err != nil {
		return err
	}
	return nil
}

func (this_ *scriptExpressionParser) parseCallExpression(expression *ast.CallExpression) (err error) {
	err = this_.addCallNames(expression)
	if err != nil {
		return
	}
	err = this_.parseExpressions(expression.ArgumentList)
	if err != nil {
		return
	}
	return
}

func (this_ *scriptExpressionParser) parseDotExpression(expression *ast.DotExpression) (err error) {
	err = this_.addStructDotNames(expression)
	if err != nil {
		return
	}
	return
}

func (this_ *scriptExpressionParser) parseIdentifier(expression *ast.Identifier) (err error) {
	return
}

func (this_ *scriptExpressionParser) parseStringLiteral(expression *ast.StringLiteral) (err error) {
	return
}

func (this_ *scriptExpressionParser) parseNumberLiteral(expression *ast.NumberLiteral) (err error) {
	return
}
