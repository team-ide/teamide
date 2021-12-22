package modelcoder

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/parser"
)

type scriptParser struct {
	script             string
	factory            FactoryScript
	factoryScriptCache map[string]interface{}
	callDotNames       []string
	varDotNames        []string
}

func (this_ *scriptParser) parse() error {
	this_.callDotNames = []string{}
	this_.varDotNames = []string{}

	program, err := parser.ParseFile(nil, "", this_.script, 0)
	if err != nil {
		return err
	}
	for _, one := range program.Body {
		println(reflect.TypeOf(one).Elem().Name())
		expression, ok := interface{}(one).(*ast.ExpressionStatement)
		if ok {
			err = this_.parseExpression(expression.Expression)
			if err != nil {
				return err
			}
		}
	}
	err = this_.checkCall()
	if err != nil {
		return err
	}
	return nil
}

func (this_ *scriptParser) checkCall() (err error) {
	for _, callName := range this_.callDotNames {
		key := strings.TrimPrefix(callName, "$factory.")
		_, find := this_.factoryScriptCache[key]
		if !find {
			err = errors.New(fmt.Sprint("call script [", callName, "] not defind"))
			return
		}
	}
	return
}

func (this_ *scriptParser) addCallNames(expression *ast.CallExpression) (err error) {
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

func (this_ *scriptParser) addDotNames(expression *ast.DotExpression) (err error) {
	if expression.Left != nil && !isIdentifier(expression.Left) {
		return
	}
	var name string
	name, err = getDotExpressionName(expression)
	if err != nil {
		return
	}
	this_.varDotNames = append(this_.varDotNames, name)
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

func (this_ *scriptParser) parseExpressions(expressions []ast.Expression) (err error) {
	for _, one := range expressions {
		err = this_.parseExpression(one)
		if err != nil {
			return err
		}
	}
	return err
}
func (this_ *scriptParser) parseExpression(expression ast.Expression) (err error) {
	var ok = false

	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.BinaryExpression)
		ok = ok_
		if ok_ {
			err = this_.parseBinaryExpression(expression_)
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

func (this_ *scriptParser) parseBinaryExpression(binaryExpression *ast.BinaryExpression) (err error) {
	println("parseBinaryExpression:", ToJSON(binaryExpression))

	err = this_.parseExpression(binaryExpression.Left)
	if err != nil {
		return err
	}
	err = this_.parseExpression(binaryExpression.Right)
	if err != nil {
		return err
	}
	return nil
}

func (this_ *scriptParser) parseCallExpression(callExpression *ast.CallExpression) (err error) {
	println("parseCallExpression:", ToJSON(callExpression))
	err = this_.addCallNames(callExpression)
	if err != nil {
		return
	}
	err = this_.parseExpressions(callExpression.ArgumentList)
	if err != nil {
		return
	}
	return
}
func (this_ *scriptParser) parseDotExpression(dotExpression *ast.DotExpression) (err error) {
	println("parseDotExpression:", ToJSON(dotExpression))
	err = this_.addDotNames(dotExpression)
	if err != nil {
		return
	}
	return
}
func (this_ *scriptParser) parseIdentifier(identifier *ast.Identifier) (err error) {
	println("parseIdentifier:", ToJSON(identifier))
	return
}
func (this_ *scriptParser) parseStringLiteral(stringLiteral *ast.StringLiteral) (err error) {
	println("parseStringLiteral:", ToJSON(stringLiteral))
	return
}
func (this_ *scriptParser) parseNumberLiteral(numberLiteral *ast.NumberLiteral) (err error) {
	println("parseNumberLiteral:", ToJSON(numberLiteral))
	return
}
