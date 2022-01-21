package invoke

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"teamide/application/base"
	"teamide/application/common"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/parser"
	"github.com/wxnacy/wgo/arrays"
)

func (this_ *ExpressionParser) Parse(parseInfo *ParseInfo) (err error) {

	if this_.expression == nil {
		var program *ast.Program
		program, err = parser.ParseFile(nil, "", this_.script, 0)
		if err != nil {
			return
		}
		if len(program.Body) == 0 || len(program.Body) > 1 {
			err = errors.New("please enter the correct value expression")
			return
		}
		expression, ok := interface{}(program.Body[0]).(*ast.ExpressionStatement)
		if !ok {
			err = errors.New("please enter the correct value expression")
			return
		}
		this_.expression = expression.Expression
	}
	err = this_.parseExpression(this_.expression, parseInfo)
	if err != nil {
		return err
	}
	return nil
}

func (this_ *ExpressionParser) parseExpressions(expressions []ast.Expression, parseInfo *ParseInfo) (err error) {
	for _, one := range expressions {
		err = this_.parseExpression(one, parseInfo)
		if err != nil {
			return err
		}
	}
	return err
}
func (this_ *ExpressionParser) parseExpression(expression ast.Expression, parseInfo *ParseInfo) (err error) {
	if expression == nil {
		return
	}
	switch expression_ := expression.(type) {
	case *ast.BinaryExpression:
		err = this_.parseBinaryExpression(expression_, parseInfo)
	case *ast.ConditionalExpression:
		err = this_.parseConditionalExpression(expression_, parseInfo)
	case *ast.CallExpression:
		err = this_.parseCallExpression(expression_, parseInfo)
	case *ast.DotExpression:
		err = this_.parseDotExpression(expression_, parseInfo)
	case *ast.Identifier:
		err = this_.parseIdentifier(expression_, parseInfo)
	case *ast.StringLiteral:
		err = this_.parseStringLiteral(expression_, parseInfo)
	case *ast.NumberLiteral:
		err = this_.parseNumberLiteral(expression_, parseInfo)
	case *ast.NullLiteral:
		err = this_.parseNullLiteral(expression_, parseInfo)
	case *ast.NewExpression:
		err = this_.parseNewExpression(expression_, parseInfo)
	case *ast.AssignExpression:
		err = this_.parseAssignExpression(expression_, parseInfo)
	case *ast.ObjectLiteral:
		err = this_.parseObjectLiteral(expression_, parseInfo)
	case *ast.ArrayLiteral:
		err = this_.parseArrayLiteral(expression_, parseInfo)
	case *ast.TemplateLiteral:
		err = this_.parseTemplateLiteral(expression_, parseInfo)
	case *ast.BooleanLiteral:
		err = this_.parseBooleanLiteral(expression_, parseInfo)
	case *ast.BracketExpression:
		err = this_.parseBracketExpression(expression_, parseInfo)
	default:
		err = errors.New(fmt.Sprint("parseExpression type not match:", reflect.TypeOf(expression).Elem().Name()))
	}
	if err != nil {
		return err
	}
	return err
}

func (this_ *ExpressionParser) parseBooleanLiteral(expression *ast.BooleanLiteral, parseInfo *ParseInfo) (err error) {

	return
}
func (this_ *ExpressionParser) parseTemplateLiteral(expression *ast.TemplateLiteral, parseInfo *ParseInfo) (err error) {
	// fmt.Println("parseTemplateLiteral Elements:", base.ToJSON(expression.Elements))
	// fmt.Println("parseTemplateLiteral Expressions:", base.ToJSON(expression.Expressions))
	// fmt.Println("parseTemplateLiteral Tag:", base.ToJSON(expression.Tag))
	return
}
func (this_ *ExpressionParser) parseObjectLiteral(expression *ast.ObjectLiteral, parseInfo *ParseInfo) (err error) {
	for _, one := range expression.Value {
		switch one_ := one.(type) {
		case *ast.PropertyKeyed:
			switch key := one_.Key.(type) {
			case *ast.Identifier:
				_ = key.Name.String()
			case *ast.StringLiteral:
				_ = key.Value.String()
			default:
				err = errors.New(fmt.Sprint("parseObjectLiteral value key type not match:", reflect.TypeOf(one_.Key).Elem().Name()))
				return
			}
			err = this_.parseExpression(one_.Value, parseInfo)
			if err != nil {
				return
			}
		default:
			err = errors.New(fmt.Sprint("parseObjectLiteral value type not match:", reflect.TypeOf(one).Elem().Name()))
			return
		}
	}
	return
}

func (this_ *ExpressionParser) parseArrayLiteral(expression *ast.ArrayLiteral, parseInfo *ParseInfo) (err error) {
	for _, one := range expression.Value {
		switch one_ := one.(type) {
		case *ast.PropertyKeyed:
			switch key := one_.Key.(type) {
			case *ast.Identifier:
				_ = key.Name.String()
			case *ast.StringLiteral:
				_ = key.Value.String()
			default:
				err = errors.New(fmt.Sprint("parseArrayLiteral value key type not match:", reflect.TypeOf(one_.Key).Elem().Name()))
				return
			}
			err = this_.parseExpression(one_.Value, parseInfo)
			if err != nil {
				return
			}
		default:
			err = errors.New(fmt.Sprint("parseArrayLiteral value type not match:", reflect.TypeOf(one).Elem().Name()))
			return
		}
	}
	return
}

func (this_ *ExpressionParser) parseNullLiteral(expression *ast.NullLiteral, parseInfo *ParseInfo) (err error) {

	return
}
func (this_ *ExpressionParser) parseNewExpression(expression *ast.NewExpression, parseInfo *ParseInfo) (err error) {
	err = this_.parseExpressions(expression.ArgumentList, parseInfo)
	if err != nil {
		return
	}
	err = this_.parseExpression(expression.Callee, parseInfo)
	if err != nil {
		return
	}
	return
}

func (this_ *ExpressionParser) parseBinaryExpression(expression *ast.BinaryExpression, parseInfo *ParseInfo) (err error) {
	err = this_.parseExpression(expression.Left, parseInfo)
	if err != nil {
		return
	}
	err = this_.parseExpression(expression.Right, parseInfo)
	if err != nil {
		return
	}
	return
}
func (this_ *ExpressionParser) parseConditionalExpression(expression *ast.ConditionalExpression, parseInfo *ParseInfo) (err error) {

	err = this_.parseExpression(expression.Test, parseInfo)
	if err != nil {
		return
	}
	err = this_.parseExpression(expression.Consequent, parseInfo)
	if err != nil {
		return
	}
	err = this_.parseExpression(expression.Alternate, parseInfo)
	if err != nil {
		return
	}
	return
}

func getParseCallArgsForValue(argumentList []ast.Expression, size int) (args []interface{}, err error) {
	for index, one := range argumentList {
		if index >= size {
			return
		}
		switch argument := one.(type) {
		case *ast.StringLiteral:
			args = append(args, argument.Value.String())
		case *ast.BooleanLiteral:
			args = append(args, argument.Value)
		case *ast.NumberLiteral:
			args = append(args, argument.Value)
		case *ast.NullLiteral:
			args = append(args, nil)
		default:
			err = errors.New(fmt.Sprint("getArgsForValue value type not match:", reflect.TypeOf(one).Elem().Name()))
			return
		}
	}
	return
}

func getParseCallArgsForScript(argumentList []ast.Expression) (argScripts []interface{}, err error) {
	for _, one := range argumentList {
		var argScript string
		switch argument := one.(type) {
		case *ast.StringLiteral:
			argScript = fmt.Sprint(argument.Value.String())
		case *ast.BooleanLiteral:
			argScript = fmt.Sprint(argument.Value)
		case *ast.NumberLiteral:
			argScript = fmt.Sprint(argument.Value)
		default:
			argScript, err = getName(one)
			if err != nil {
				return
			}
		}
		argScripts = append(argScripts, argScript)
	}
	return
}
func (this_ *ExpressionParser) parseCallExpression(expression *ast.CallExpression, parseInfo *ParseInfo) (err error) {
	var funcName string
	funcName, err = getName(expression.Callee)
	if err != nil {
		return
	}
	if arrays.ContainsString(parseInfo.UseFunctions, funcName) == -1 {
		parseInfo.UseFunctions = append(parseInfo.UseFunctions, funcName)
	}
	err = this_.parseExpressions(expression.ArgumentList, parseInfo)
	if err != nil {
		return
	}
	for callKey, callFun := range parseCallMap {
		if strings.HasSuffix(funcName, callKey) {
			var args []interface{}
			if callKey == "service" {
				args, err = getParseCallArgsForScript(expression.ArgumentList)
			} else {
				args, err = getParseCallArgsForValue(expression.ArgumentList, 100)
			}
			if err != nil {
				if parseInfo.App.GetLogger() != nil {
					parseInfo.App.GetLogger().Error("parseCall [", funcName, "] getParseCallArgsForValue error argumentList:", base.ToJSON(expression.ArgumentList))
					parseInfo.App.GetLogger().Error("parseCall [", funcName, "] getParseCallArgsForValue error:", err)
				}
				return
			}
			prefixName := ""
			if strings.Contains(funcName, ".") {
				prefixName = funcName[:strings.LastIndex(funcName, ".")]
			}
			if prefixName == "_" {
				prefixName = ""
			}
			err = callFun(parseInfo, prefixName, args)
			if err != nil {
				return
			}
			return
		}
	}

	return
}

func (this_ *ExpressionParser) parseBracketExpression(expression *ast.BracketExpression, parseInfo *ParseInfo) (err error) {
	var name string
	name, err = getName(expression)
	if err != nil {
		if parseInfo.App.GetLogger() != nil {
			parseInfo.App.GetLogger().Error("parse bracket name error:", err)
		}
		return
	}
	var dataInfo *common.InvokeDataInfo
	dataInfo, err = parseInfo.InvokeNamespace.GetDataInfo(name)
	if err != nil {
		return
	}
	if dataInfo != nil {
		dataInfo.IsUse = true
	}

	return
}

func (this_ *ExpressionParser) parseAssignExpression(expression *ast.AssignExpression, parseInfo *ParseInfo) (err error) {
	var name string
	name, err = getName(expression.Left)
	if err != nil {
		return
	}
	var dataInfo *common.InvokeDataInfo
	dataInfo, err = parseInfo.InvokeNamespace.GetDataInfo(name)
	if err != nil {
		return
	}
	if dataInfo == nil {
		err = base.NewError("", "parse assign data info [", name, "] not defind")
		return
	}
	dataInfo.IsSetValue = true
	err = this_.parseExpression(expression.Right, parseInfo)
	if err != nil {
		return
	}
	return
}
func (this_ *ExpressionParser) parseDotExpression(expression *ast.DotExpression, parseInfo *ParseInfo) (err error) {
	var name string
	name, err = getName(expression)
	if err != nil {
		if parseInfo.App.GetLogger() != nil {
			parseInfo.App.GetLogger().Error("parse dot name error:", err)
		}
		return
	}
	var dataInfo *common.InvokeDataInfo
	dataInfo, err = parseInfo.InvokeNamespace.GetDataInfo(name)
	if err != nil {
		return
	}
	if dataInfo == nil {
		if parseInfo.App.GetLogger() != nil {
			parseInfo.App.GetLogger().Error("parse dot name [", name, "] not defind")
		}
		return
	}
	dataInfo.IsUse = true
	return
}

func (this_ *ExpressionParser) parseIdentifier(expression *ast.Identifier, parseInfo *ParseInfo) (err error) {
	name := expression.Name.String()

	var dataInfo *common.InvokeDataInfo
	dataInfo, err = parseInfo.InvokeNamespace.GetDataInfo(name)
	if err != nil {
		return
	}
	dataInfo.IsUse = true
	return
}

func (this_ *ExpressionParser) parseStringLiteral(expression *ast.StringLiteral, parseInfo *ParseInfo) (err error) {
	return
}

func (this_ *ExpressionParser) parseNumberLiteral(expression *ast.NumberLiteral, parseInfo *ParseInfo) (err error) {
	return
}
