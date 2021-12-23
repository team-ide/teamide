package modelcoder

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/dop251/goja/ast"
)

func getColumnValue(application *Application, variable *invokeVariable, name string, valueScript string) (res interface{}, err error) {
	if IsEmpty(valueScript) {
		variableData := variable.GetVariableData(name)
		if variableData != nil {
			res = variableData.Data
		}
	} else {
		res, err = getScriptValue(application, variable, valueScript)
	}
	return
}

func ifScriptValue(application *Application, variable *invokeVariable, ifScript string) (res bool, err error) {
	if IsEmpty(ifScript) {
		res = true
		return
	}
	var value interface{}
	value, err = getScriptValue(application, variable, ifScript)
	if err != nil {
		return
	}
	if value == nil {
		return
	}
	if value == true || value == "1" || value == "true" {
		res = true
		return
	}
	return
}

func getScriptValue(application *Application, variable *invokeVariable, script string) (res interface{}, err error) {
	if IsEmpty(script) {
		res = nil
	} else {
		var parser *scriptValueParser
		parser, err = getScriptValueParser(application, script)
		if err != nil {
			return
		}
		if parser != nil {
			res, err = parser.invoke(variable)
		}
	}
	return
}

func getScriptValueParser(application *Application, script string) (res *scriptValueParser, err error) {
	if IsEmpty(script) {
		return
	}
	var ok bool
	res, ok = application.scriptValueParserCache[script]
	if !ok {
		res = &scriptValueParser{
			script:      script,
			application: application,
		}
		err = res.parse()
		if err != nil {
			return nil, err
		}
		application.scriptValueParserCache[script] = res
	}
	return
}

func (this_ *scriptValueParser) check(variable *invokeVariable) (err error) {

	for _, callName := range this_.callDotNames {
		key := strings.TrimPrefix(callName, "$factory.")
		_, find := this_.application.factoryScriptCache[key]
		if !find {
			err = errors.New(fmt.Sprint("call script [", callName, "] not defind"))
			return
		}
	}
	for _, structDotName := range this_.structDotNames {
		dotIndex := strings.Index(structDotName, ".")
		if dotIndex > 0 {
			variableName := structDotName[0:dotIndex]
			variableData := variable.GetVariableData(variableName)
			if variableData == nil {
				err = errors.New(fmt.Sprint("variable [", variableName, "] not defind"))
				return
			}
			if variableData.DataStruct == nil {
				err = errors.New(fmt.Sprint("variable [", variableName, "] struct not defind"))
				return
			}
			fieldName := structDotName[dotIndex+1:]
			field := variableData.DataStruct.GetField(fieldName)
			if field == nil {
				err = errors.New(fmt.Sprint("variable [", variableName, "] struct field [", fieldName, "] not defind"))
				return
			}
		}
	}
	return
}

func (this_ *scriptValueParser) invoke(variable *invokeVariable) (res interface{}, err error) {
	err = this_.check(variable)
	if err != nil {
		return
	}
	res, err = this_.invokeExpression(this_.expression, variable)
	return
}

func (this_ *scriptValueParser) invokeExpressions(expressions []ast.Expression, variable *invokeVariable) (ress []interface{}, err error) {
	for _, one := range expressions {
		var res interface{}
		res, err = this_.invokeExpression(one, variable)
		if err != nil {
			return
		}
		ress = append(ress, res)
	}
	return
}

func (this_ *scriptValueParser) invokeExpression(expression ast.Expression, variable *invokeVariable) (res interface{}, err error) {
	var ok = false

	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.BinaryExpression)
		ok = ok_
		if ok_ {
			res, err = this_.invokeBinaryExpression(expression_, variable)
		}
	}
	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.ConditionalExpression)
		ok = ok_
		if ok_ {
			res, err = this_.invokeConditionalExpression(expression_, variable)
		}
	}
	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.CallExpression)
		ok = ok_
		if ok_ {
			res, err = this_.invokeCallExpression(expression_, variable)
		}
	}
	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.DotExpression)
		ok = ok_
		if ok_ {
			res, err = this_.invokeDotExpression(expression_, variable)
		}
	}
	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.Identifier)
		ok = ok_
		if ok_ {
			res, err = this_.invokeIdentifier(expression_, variable)
		}
	}
	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.StringLiteral)
		ok = ok_
		if ok_ {
			res, err = this_.invokeStringLiteral(expression_, variable)
		}
	}
	if !ok {
		expression_, ok_ := interface{}(expression).(*ast.NumberLiteral)
		ok = ok_
		if ok_ {
			res, err = this_.invokeNumberLiteral(expression_, variable)
		}
	}

	if err != nil {
		return
	}
	if !ok {
		err = errors.New(fmt.Sprint("invokeExpression type not match:", reflect.TypeOf(expression).Elem().Name()))
	}
	return
}

func ToInt64(value interface{}) (res int64, ok bool) {
	if value == nil {
		return
	}
	refType := GetRefType(value)
	refValue := GetRefValue(value)
	switch refType.Name() {
	case "int", "int8", "int16", "int32", "int64":
		ok = true
		res = refValue.Int()
	}
	return
}

func ToFloat64(value interface{}) (res float64, ok bool) {
	if value == nil {
		return
	}
	refType := GetRefType(value)
	refValue := GetRefValue(value)
	switch refType.Name() {
	case "float32", "float64":
		ok = true
		res = refValue.Float()
	}
	return
}

func (this_ *scriptValueParser) invokeBinaryExpression(expression *ast.BinaryExpression, variable *invokeVariable) (res interface{}, err error) {

	var leftValue interface{}
	leftValue, err = this_.invokeExpression(expression.Left, variable)
	if err != nil {
		return
	}
	var rightValue interface{}
	rightValue, err = this_.invokeExpression(expression.Right, variable)
	if err != nil {
		return
	}

	leftInt64, leftInt64Ok := ToInt64(leftValue)
	leftFloat64, leftFloat64Ok := ToFloat64(leftValue)

	rightInt64, rightInt64Ok := ToInt64(rightValue)
	rightFloat64, rightFloat64Ok := ToFloat64(rightValue)

	switch expression.Operator.String() {
	case "+":
		if leftInt64Ok && rightInt64Ok {
			res = leftInt64 + rightInt64
		} else if (leftInt64Ok && rightFloat64Ok) || (leftFloat64Ok && rightInt64Ok) || (leftFloat64Ok && rightFloat64Ok) {
			if !leftFloat64Ok {
				leftFloat64 = float64(leftInt64)
			}
			if !rightFloat64Ok {
				rightFloat64 = float64(rightInt64)
			}
			res = leftFloat64 + rightFloat64
		} else {
			res = fmt.Sprint(leftValue, rightValue)
		}
	case "-":
		if leftInt64Ok && rightInt64Ok {
			res = leftInt64 - rightInt64
		} else if (leftInt64Ok && rightFloat64Ok) || (leftFloat64Ok && rightInt64Ok) || (leftFloat64Ok && rightFloat64Ok) {
			if !leftFloat64Ok {
				leftFloat64 = float64(leftInt64)
			}
			if !rightFloat64Ok {
				rightFloat64 = float64(rightInt64)
			}
			res = leftFloat64 - rightFloat64
		} else {
			err = errors.New(fmt.Sprint("value [", leftValue, "] - value [", rightValue, "] error"))
		}
	case "*":
		if leftInt64Ok && rightInt64Ok {
			res = leftInt64 * rightInt64
		} else if (leftInt64Ok && rightFloat64Ok) || (leftFloat64Ok && rightInt64Ok) || (leftFloat64Ok && rightFloat64Ok) {
			if !leftFloat64Ok {
				leftFloat64 = float64(leftInt64)
			}
			if !rightFloat64Ok {
				rightFloat64 = float64(rightInt64)
			}
			res = leftFloat64 * rightFloat64
		} else {
			err = errors.New(fmt.Sprint("value [", leftValue, "] * value [", rightValue, "] error"))
		}
	case "/":
		if leftInt64Ok && rightInt64Ok {
			res = leftInt64 / rightInt64
		} else if (leftInt64Ok && rightFloat64Ok) || (leftFloat64Ok && rightInt64Ok) || (leftFloat64Ok && rightFloat64Ok) {
			if !leftFloat64Ok {
				leftFloat64 = float64(leftInt64)
			}
			if !rightFloat64Ok {
				rightFloat64 = float64(rightInt64)
			}
			res = leftFloat64 / rightFloat64
		} else {
			err = errors.New(fmt.Sprint("value [", leftValue, "] / value [", rightValue, "] error"))
		}
	case "%":
		if leftInt64Ok && rightInt64Ok {
			res = leftInt64 % rightInt64
		} else {
			err = errors.New(fmt.Sprint("value [", leftValue, "] % value [", rightValue, "] error"))
		}
	case "&&":
		res = this_.application.factory.IsTrue(leftValue) && this_.application.factory.IsTrue(rightValue)
	case "||":
		res = this_.application.factory.IsTrue(leftValue) || this_.application.factory.IsTrue(rightValue)
	case "==":
		res = leftValue == rightValue
	case "===":
		res = leftValue == rightValue
	case "!=":
		res = leftValue != rightValue
	case "!==":
		res = leftValue != rightValue
	case "<":
		if leftInt64Ok && rightInt64Ok {
			res = leftInt64 < rightInt64
		} else if (leftInt64Ok && rightFloat64Ok) || (leftFloat64Ok && rightInt64Ok) || (leftFloat64Ok && rightFloat64Ok) {
			if !leftFloat64Ok {
				leftFloat64 = float64(leftInt64)
			}
			if !rightFloat64Ok {
				rightFloat64 = float64(rightInt64)
			}
			res = leftFloat64 < rightFloat64
		} else {
			err = errors.New(fmt.Sprint("value [", leftValue, "] < value [", rightValue, "] error"))
		}
	case "<=":
		if leftInt64Ok && rightInt64Ok {
			res = leftInt64 <= rightInt64
		} else if (leftInt64Ok && rightFloat64Ok) || (leftFloat64Ok && rightInt64Ok) || (leftFloat64Ok && rightFloat64Ok) {
			if !leftFloat64Ok {
				leftFloat64 = float64(leftInt64)
			}
			if !rightFloat64Ok {
				rightFloat64 = float64(rightInt64)
			}
			res = leftFloat64 <= rightFloat64
		} else {
			err = errors.New(fmt.Sprint("value [", leftValue, "] <= value [", rightValue, "] error"))
		}
	case ">":
		if leftInt64Ok && rightInt64Ok {
			res = leftInt64 > rightInt64
		} else if (leftInt64Ok && rightFloat64Ok) || (leftFloat64Ok && rightInt64Ok) || (leftFloat64Ok && rightFloat64Ok) {
			if !leftFloat64Ok {
				leftFloat64 = float64(leftInt64)
			}
			if !rightFloat64Ok {
				rightFloat64 = float64(rightInt64)
			}
			res = leftFloat64 > rightFloat64
		} else {
			err = errors.New(fmt.Sprint("value [", leftValue, "] > value [", rightValue, "] error"))
		}
	case ">=":
		fmt.Println("leftInt64Ok:", leftInt64Ok)
		fmt.Println("leftFloat64Ok:", leftFloat64Ok)
		fmt.Println("rightInt64Ok:", rightInt64Ok)
		fmt.Println("rightFloat64Ok:", rightFloat64Ok)
		if leftInt64Ok && rightInt64Ok {
			res = leftInt64 >= rightInt64
		} else if (leftInt64Ok && rightFloat64Ok) || (leftFloat64Ok && rightInt64Ok) || (leftFloat64Ok && rightFloat64Ok) {
			if !leftFloat64Ok {
				leftFloat64 = float64(leftInt64)
			}
			if !rightFloat64Ok {
				rightFloat64 = float64(rightInt64)
			}
			res = leftFloat64 >= rightFloat64
		} else {
			err = errors.New(fmt.Sprint("value [", leftValue, "] >= value [", rightValue, "] error"))
		}
	default:
		err = errors.New(fmt.Sprint("expression operator [", expression.Operator, "] not defind"))
	}
	return
}

func (this_ *scriptValueParser) invokeConditionalExpression(expression *ast.ConditionalExpression, variable *invokeVariable) (res interface{}, err error) {
	var test interface{}
	test, err = this_.invokeExpression(expression.Test, variable)
	if err != nil {
		return
	}
	if this_.application.factory.IsTrue(test) {
		res, err = this_.invokeExpression(expression.Consequent, variable)
	} else {
		res, err = this_.invokeExpression(expression.Alternate, variable)
	}
	return
}
func (this_ *scriptValueParser) invokeCallExpression(expression *ast.CallExpression, variable *invokeVariable) (res interface{}, err error) {
	println("invokeCallExpression:", ToJSON(expression))
	var args []interface{}
	args, err = this_.invokeExpressions(expression.ArgumentList, variable)
	if err != nil {
		return
	}
	println("invokeCallExpression args:", ToJSON(args))
	return
}

func (this_ *scriptValueParser) invokeDotExpression(expression *ast.DotExpression, variable *invokeVariable) (res interface{}, err error) {
	identifier, ok := interface{}(expression.Left).(*ast.Identifier)
	if ok {
		key := identifier.Name.String()
		variableData := variable.GetVariableData(key)
		if variableData == nil || variableData.Data == nil {
			err = errors.New(fmt.Sprint("variable [", key, "] not defind"))
			return
		}
		var dataMap map[string]interface{}
		dataMap, err = StructToMap(variableData.Data, "", "")
		if err != nil {
			return
		}
		key = expression.Identifier.Name.String()
		res = dataMap[key]
	} else {
		err = errors.New(fmt.Sprint("invokeDotExpression left type not match:", reflect.TypeOf(expression.Left).Elem().Name()))
	}
	return
}

func (this_ *scriptValueParser) invokeIdentifier(expression *ast.Identifier, variable *invokeVariable) (res interface{}, err error) {
	key := expression.Name.String()
	variableData := variable.GetVariableData(key)
	if variableData != nil {
		res = variableData.Data
	}
	return
}

func (this_ *scriptValueParser) invokeStringLiteral(expression *ast.StringLiteral, variable *invokeVariable) (res interface{}, err error) {
	res = expression.Value.String()
	return
}

func (this_ *scriptValueParser) invokeNumberLiteral(expression *ast.NumberLiteral, variable *invokeVariable) (res interface{}, err error) {
	reT := GetRefType(expression.Value)
	reV := GetRefValue(expression.Value)
	res = GetFieldTypeValue(reT, reV)
	return
}
