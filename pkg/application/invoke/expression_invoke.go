package invoke

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	base2 "teamide/pkg/application/base"
	"teamide/pkg/application/common"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/parser"
)

func (this_ *ExpressionParser) check(invokeInfo *InvokeInfo) (err error) {

	return
}

func (this_ *ExpressionParser) invoke(invokeInfo *InvokeInfo) (res interface{}, err error) {
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
	err = this_.check(invokeInfo)
	if err != nil {
		return
	}
	res, err = this_.invokeExpression(this_.expression, invokeInfo)
	if err != nil {
		return
	}
	return
}

func (this_ *ExpressionParser) invokeExpressions(expressions []ast.Expression, invokeInfo *InvokeInfo) (ress []interface{}, err error) {
	for _, one := range expressions {
		var res interface{}
		res, err = this_.invokeExpression(one, invokeInfo)
		if err != nil {
			return
		}
		ress = append(ress, res)
	}
	return
}

func (this_ *ExpressionParser) invokeExpression(expression ast.Expression, invokeInfo *InvokeInfo) (res interface{}, err error) {
	if expression == nil {
		return
	}
	switch expression_ := expression.(type) {
	case *ast.BinaryExpression:
		res, err = this_.invokeBinaryExpression(expression_, invokeInfo)
	case *ast.ConditionalExpression:
		res, err = this_.invokeConditionalExpression(expression_, invokeInfo)
	case *ast.CallExpression:
		res, err = this_.invokeCallExpression(expression_, invokeInfo)
	case *ast.DotExpression:
		res, err = this_.invokeDotExpression(expression_, invokeInfo)
	case *ast.Identifier:
		res, err = this_.invokeIdentifier(expression_, invokeInfo)
	case *ast.StringLiteral:
		res, err = this_.invokeStringLiteral(expression_, invokeInfo)
	case *ast.NumberLiteral:
		res, err = this_.invokeNumberLiteral(expression_, invokeInfo)
	case *ast.NullLiteral:
		res, err = this_.invokeNullLiteral(expression_, invokeInfo)
	case *ast.TemplateLiteral:
		res, err = this_.invokeTemplateLiteral(expression_, invokeInfo)
	case *ast.AssignExpression:
		res, err = this_.invokeAssignExpression(expression_, invokeInfo)
	case *ast.ObjectLiteral:
		res, err = this_.invokeObjectLiteral(expression_, invokeInfo)
	case *ast.ArrayLiteral:
		res, err = this_.invokeArrayLiteral(expression_, invokeInfo)
	case *ast.BooleanLiteral:
		res, err = this_.invokeBooleanLiteral(expression_, invokeInfo)
	case *ast.BracketExpression:
		res, err = this_.invokeBracketExpression(expression_, invokeInfo)
	default:
		err = errors.New(fmt.Sprint("invokeExpression type not match:", reflect.TypeOf(expression).Elem().Name()))
	}
	if err != nil {
		return
	}

	return
}

func (this_ *ExpressionParser) invokeBooleanLiteral(expression *ast.BooleanLiteral, invokeInfo *InvokeInfo) (res interface{}, err error) {
	res = expression.Value
	return
}
func (this_ *ExpressionParser) invokeObjectLiteral(expression *ast.ObjectLiteral, invokeInfo *InvokeInfo) (res interface{}, err error) {
	res = map[string]interface{}{}
	return
}

func (this_ *ExpressionParser) invokeArrayLiteral(expression *ast.ArrayLiteral, invokeInfo *InvokeInfo) (res interface{}, err error) {
	res = []interface{}{}
	return
}

func (this_ *ExpressionParser) invokeTemplateLiteral(expression *ast.TemplateLiteral, invokeInfo *InvokeInfo) (res interface{}, err error) {
	str := ""
	for _, one := range expression.Elements {
		str += one.Literal
	}
	res = str
	return
}
func (this_ *ExpressionParser) invokeNullLiteral(expression *ast.NullLiteral, invokeInfo *InvokeInfo) (res interface{}, err error) {
	res = nil
	return
}
func (this_ *ExpressionParser) invokeBinaryExpression(expression *ast.BinaryExpression, invokeInfo *InvokeInfo) (res interface{}, err error) {

	var leftValue interface{}
	leftValue, err = this_.invokeExpression(expression.Left, invokeInfo)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("invoke binary left error:", err)
		}
		return
	}
	var rightValue interface{}
	rightValue, err = this_.invokeExpression(expression.Right, invokeInfo)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("invoke binary right error:", err)
		}
		return
	}

	leftInt64, leftInt64Ok := base2.ToInt64(leftValue)
	leftFloat64, leftFloat64Ok := base2.ToFloat64(leftValue)

	rightInt64, rightInt64Ok := base2.ToInt64(rightValue)
	rightFloat64, rightFloat64Ok := base2.ToFloat64(rightValue)

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
		res = invokeInfo.App.GetScript().IsTrue(leftValue) && invokeInfo.App.GetScript().IsTrue(rightValue)
	case "||":
		res = invokeInfo.App.GetScript().IsTrue(leftValue) || invokeInfo.App.GetScript().IsTrue(rightValue)
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

func (this_ *ExpressionParser) invokeConditionalExpression(expression *ast.ConditionalExpression, invokeInfo *InvokeInfo) (res interface{}, err error) {
	var test interface{}
	test, err = this_.invokeExpression(expression.Test, invokeInfo)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("invoke conditional test error:", err)
		}
		return
	}
	if invokeInfo.App.GetScript().IsTrue(test) {
		res, err = this_.invokeExpression(expression.Consequent, invokeInfo)
	} else {
		res, err = this_.invokeExpression(expression.Alternate, invokeInfo)
	}
	return
}

func getName(expression ast.Expression) (name string, err error) {
	switch expression_ := expression.(type) {
	case *ast.DotExpression:
		var leftName string
		leftName, err = getName(expression_.Left)
		if err != nil {
			return
		}
		var rightName string
		rightName, err = getName(&expression_.Identifier)
		if err != nil {
			return
		}
		name = leftName + "." + rightName
	case *ast.Identifier:
		name = expression_.Name.String()
	case *ast.BracketExpression:
		var leftName string
		leftName, err = getName(expression_.Left)
		if err != nil {
			return
		}
		var rightName string
		rightName, err = getName(expression_.Member)
		if err != nil {
			return
		}

		name = leftName + "[" + rightName + "]"

	case *ast.NumberLiteral:
		name = expression_.Literal
	case *ast.StringLiteral:
		name = expression_.Literal
	default:
		err = errors.New(fmt.Sprint("getName type not match:", reflect.TypeOf(expression).Elem().Name()))
		return
	}
	return
}

func (this_ *ExpressionParser) invokeCallExpression(expression *ast.CallExpression, invokeInfo *InvokeInfo) (res interface{}, err error) {

	var funcName string
	funcName, err = getName(expression.Callee)
	if err != nil {
		return
	}

	var args []interface{}
	args, err = this_.invokeExpressions(expression.ArgumentList, invokeInfo)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("call func [", funcName, "] build args error:", err)
		}
		return
	}
	for callKey, callFun := range invokeCallMap {
		if strings.HasSuffix(funcName, callKey) {
			prefixName := ""
			if strings.Contains(funcName, ".") {
				prefixName = funcName[:strings.LastIndex(funcName, ".")]
			}
			if prefixName == "_" {
				prefixName = ""
			}
			// if invokeInfo.App.GetLogger() != nil && invokeInfo.App.GetLogger().OutDebug() {
			// 	invokeInfo.App.GetLogger().Debug("call func [", funcName, "] prefixName [", prefixName, "] args:", base.ToJSON(args))
			// }
			res, err = callFun(invokeInfo, prefixName, args)
			if err != nil {
				// if invokeInfo.App.GetLogger() != nil {
				// 	invokeInfo.App.GetLogger().Error("call func [", funcName, "] error:", err)
				// }
				return
			}
			return
		}
	}

	funcName = strings.TrimPrefix(funcName, "$factory.")

	method := invokeInfo.App.GetScriptMethod(funcName)
	if method.Name == "" {
		err = errors.New(fmt.Sprint("func [", funcName, "] not defind"))
		return
	}
	var callArgs []reflect.Value = []reflect.Value{}
	callArgs = append(callArgs, reflect.ValueOf(invokeInfo.App.GetScript()))
	if len(args) > 0 {

		// fmt.Println("callName:", method.Name)
		for index, arg := range args {
			argType := method.Type.In(index + 1)
			switch argType.Name() {
			case "int":
				if arg == nil {
					arg = 0
				} else {
					arg, err = strconv.Atoi(fmt.Sprint(arg))
					if err != nil {
						return
					}
				}
			case "string":
				if arg == nil {
					arg = ""
				}
			}
			callArgs = append(callArgs, reflect.ValueOf(arg))
		}
	}
	values := method.Func.Call(callArgs)
	if len(values) == 1 {
		res = values[0].Interface()
	} else if len(values) == 2 {
		res = values[0].Interface()
		if values[1].Interface() != nil {
			err = values[1].Interface().(error)
		}
	}
	if res != nil {
		_, isError := res.(error)
		if isError {
			err = res.(error)
			return
		}
	}
	return
}

func (this_ *ExpressionParser) invokeAssignExpression(expression *ast.AssignExpression, invokeInfo *InvokeInfo) (res interface{}, err error) {
	var name string
	name, err = getName(expression.Left)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("invoke assign left error:", err)
		}
		return
	}
	var value interface{}
	value, err = this_.invokeExpression(expression.Right, invokeInfo)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("invoke assign right error:", err)
		}
		return
	}

	err = invokeInfo.InvokeNamespace.SetData(name, value, nil)

	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("assign invoke data [", name, "] value [", value, "] error:", err)
		}
	}
	return
}

func (this_ *ExpressionParser) invokeBracketExpression(expression *ast.BracketExpression, invokeInfo *InvokeInfo) (res interface{}, err error) {
	var name string
	name, err = getName(expression)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("invoke bracket name error:", err)
		}
		return
	}

	var data *common.InvokeData
	data, err = invokeInfo.InvokeNamespace.GetData(name)
	if err != nil {
		return
	}
	if data == nil {
		err = base2.NewError("", "invoke bracket data [", name, "] not defind")
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("invoke bracket [", name, "] data error:", err)
		}
		return
	}
	res = data.Value

	return
}

func (this_ *ExpressionParser) invokeDotExpression(expression *ast.DotExpression, invokeInfo *InvokeInfo) (res interface{}, err error) {
	var name string
	name, err = getName(expression)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("invoke dot name error:", err)
		}
		return
	}

	var data *common.InvokeData
	data, err = invokeInfo.InvokeNamespace.GetData(name)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("invoke dot get [", name, "] data error:", err)
		}
		return
	}
	if data == nil {
		err = base2.NewError("", "invoke dot data [", name, "] not defind")
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("invoke dot [", name, "] data error:", err)
		}
		return
	}
	res = data.Value

	return
}

func (this_ *ExpressionParser) invokeIdentifier(expression *ast.Identifier, invokeInfo *InvokeInfo) (res interface{}, err error) {
	name := expression.Name.String()

	var data *common.InvokeData
	data, err = invokeInfo.InvokeNamespace.GetData(name)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("invoke identifier get [", name, "] data error:", err)
		}
		return
	}
	if data == nil {
		err = base2.NewError("", "data [", name, "] not defind")
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("invoke identifier [", name, "] data error:", err)
		}
		return
	}
	res = data.Value
	return
}

func (this_ *ExpressionParser) invokeStringLiteral(expression *ast.StringLiteral, invokeInfo *InvokeInfo) (res interface{}, err error) {
	res = expression.Value.String()
	return
}

func (this_ *ExpressionParser) invokeNumberLiteral(expression *ast.NumberLiteral, invokeInfo *InvokeInfo) (res interface{}, err error) {
	reT := base2.GetRefType(expression.Value)
	reV := base2.GetRefValue(expression.Value)
	res = base2.GetFieldTypeValue(reT, reV)
	return
}
