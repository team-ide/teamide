package maker

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja/ast"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"reflect"
)

func (this_ *CompilerMethod) Expression(expression ast.Expression) (err error) {
	if expression == nil {
		return
	}
	//fmt.Println("TODO Expression:", util.GetStringValue(expression))
	switch e := expression.(type) {
	case *ast.CallExpression:
		_, err = this_.CallExpression(e)
		break
	case *ast.FunctionLiteral:
		err = this_.FunctionLiteral(e)
		break
	case *ast.AssignExpression:
		err = this_.AssignExpression(e)
		break
	case *ast.TemplateLiteral:
		_, err = this_.TemplateLiteral(e)
		break
	case *ast.BinaryExpression:
		err = this_.BinaryExpression(e)
		break
	case *ast.Identifier:
		err = this_.Identifier(e)
		break
	case *ast.StringLiteral:
		err = this_.StringLiteral(e)
		break
	case *ast.NumberLiteral:
		err = this_.NumberLiteral(e)
		break
	case *ast.NullLiteral:
		err = this_.NullLiteral(e)
		break
	case *ast.DotExpression:
		err = this_.DotExpression(e)
		break
	default:
		err = this_.Error("expression ["+reflect.TypeOf(expression).String()+"] 不支持", expression)
		util.Logger.Debug(this_.GetKey()+" Expression error", zap.Error(err))
		break

	}
	return
}

func (this_ *CompilerMethod) ArgumentList(argumentList []ast.Expression) (values []interface{}, err error) {
	//fmt.Println("TODO ArgumentList:", util.GetStringValue(argumentList))
	var v *ValueType
	var nameScript string
	for _, arg := range argumentList {
		nameScript, v, err = this_.GetExpressionForType(arg)
		if err != nil {
			return
		}
		values = append(values, v)

		if nameScript != "" {
			this_.fullImport(nameScript)
		}
	}
	return
}

func (this_ *CompilerMethod) CallExpression(expression *ast.CallExpression) (res any, err error) {
	//fmt.Println("TODO CallExpression ArgumentList:", util.GetStringValue(expression.ArgumentList))
	args, err := this_.ArgumentList(expression.ArgumentList)
	if err != nil {
		return
	}
	//fmt.Println("TODO CallExpression Callee:", util.GetStringValue(expression.Callee))

	nameScript, v, err := this_.GetExpressionForValue(expression.Callee)
	if err != nil {
		return
	}
	if v == nil {
		if nameScript == "" {
			return
		}
		err = errors.New("call expression method [" + nameScript + "] is null")
		return
	}
	//fmt.Println("CallExpression GetExpressionForValue res:", util.GetStringValue(v))

	if nameScript != "" {
		this_.fullImport(nameScript)
	}
	this_.CallCache[expression] = v
	this_.CallScriptCache[expression] = nameScript
	switch m := v.(type) {
	case *ComponentMethod:
		if m.GetReturnTypes == nil {
			err = errors.New("call expression method [" + nameScript + "] getReturnTypes not set")
			return
		}
		fmt.Println("TODO CallExpression call start:", nameScript, ",args:", util.GetStringValue(args))
		defer func() {
			if e := recover(); e != nil {
				err = errors.New("call method [" + nameScript + "] error")
				util.Logger.Error("CallExpression error", zap.Any("error", err))
			}
		}()
		res = m.GetReturnTypes(args)
		fmt.Println("TODO CallExpression call end:", nameScript, ",res:", util.GetStringValue(res))
		break
	case *CompilerMethod:
		var re *CompilerMethodResult
		re, err = m.CompileMethod(m)
		if err != nil {
			return
		}
		res = re
		break
	case *CompilerMethodResult:
		res = m.valueType
		break
	case *ValueType:
		res = m
		break
	default:
		if reflect.TypeOf(v).Kind() == reflect.Func {
			f := reflect.TypeOf(v)
			n := f.NumOut()
			for i := 0; i < n; i++ {
				out := f.Out(i)
				fmt.Println(out.Kind())
				switch out.Kind() {
				case reflect.String:
					res = ValueTypeString
					break
				case reflect.Int8:
					res = ValueTypeInt8
					break
				case reflect.Int16:
					res = ValueTypeInt16
					break
				case reflect.Int:
					res = ValueTypeInt
					break
				case reflect.Int32:
					res = ValueTypeInt32
					break
				case reflect.Int64:
					res = ValueTypeInt64
					break
				case reflect.Float32:
					res = ValueTypeFloat32
					break
				case reflect.Float64:
					res = ValueTypeFloat64
					break
				case reflect.Bool:
					res = ValueTypeBool
					break
				default:
					err = errors.New("call expression func [" + reflect.TypeOf(v).String() + "] not support result type [" + out.Kind().String() + "]")
					return
				}
				if res != nil {
					break
				}
			}
			break
		}
		err = errors.New("call expression method [" + reflect.TypeOf(v).String() + "] not support")
		break
	}
	return
}

func (this_ *CompilerMethod) AssignExpression(expression *ast.AssignExpression) (err error) {

	//fmt.Println("TODO AssignExpression Left:", util.GetStringValue(expression.Left))
	nameScript, err := this_.GetExpressionScript(expression.Left)
	if err != nil {
		return
	}
	this_.AssignExpressionScriptCache[expression] = nameScript

	methodVar := this_.getVar(nameScript)
	methodParam := this_.getParam(nameScript)
	if methodVar == nil && methodParam == nil {
		err = this_.Error("变量["+nameScript+"]未定义", expression)
		return
	}
	//fmt.Println("AssignExpression Right:", util.GetStringValue(expression.Right))
	vName, v, err := this_.GetExpressionForType(expression.Right)
	if err != nil {
		return
	}

	if vName != "" {
		this_.fullImport(vName)
	}

	if methodVar != nil {
		err = methodVar.addValueTypes(v)
		this_.AssignExpressionScriptTypeCache[expression] = methodVar
	} else {
		err = methodParam.addValueTypes(v)
		this_.AssignExpressionScriptTypeCache[expression] = methodParam
	}
	if err != nil {
		err = this_.Error(err.Error(), expression.Left)
		return
	}
	util.Logger.Debug("AssignExpression var set", zap.Any("name", nameScript), zap.Any("type", v))
	if err != nil {
		return
	}
	return
}

func (this_ *CompilerMethod) TemplateLiteral(expression *ast.TemplateLiteral) (res any, err error) {
	fmt.Println("TODO TemplateLiteral:", util.GetStringValue(expression))
	this_.fullImport("fmt")
	res = ValueTypeString
	return
}
func (this_ *CompilerMethod) Identifier(expression *ast.Identifier) (err error) {
	fmt.Println("TODO Identifier:", util.GetStringValue(expression))
	return
}
func (this_ *CompilerMethod) NullLiteral(expression *ast.NullLiteral) (err error) {
	fmt.Println("TODO NullLiteral:", util.GetStringValue(expression))
	return
}
func (this_ *CompilerMethod) StringLiteral(expression *ast.StringLiteral) (err error) {
	fmt.Println("TODO StringLiteral:", util.GetStringValue(expression))
	return
}
func (this_ *CompilerMethod) NumberLiteral(expression *ast.NumberLiteral) (err error) {
	fmt.Println("TODO NumberLiteral:", util.GetStringValue(expression))
	return
}
func (this_ *CompilerMethod) DotExpression(expression *ast.DotExpression) (err error) {
	fmt.Println("TODO DotExpression:", util.GetStringValue(expression))
	err = this_.Expression(expression.Left)
	if err != nil {
		return
	}
	err = this_.Identifier(&expression.Identifier)
	if err != nil {
		return
	}
	return
}

func (this_ *CompilerMethod) BinaryExpression(expression *ast.BinaryExpression) (err error) {

	err = this_.Expression(expression.Left)
	if err != nil {
		return
	}
	err = this_.Expression(expression.Right)
	if err != nil {
		return
	}
	return
}

func (this_ *CompilerMethod) GetExpressionForValue(expression ast.Expression) (nameScript string, res any, err error) {

	//fmt.Println("TODO GetExpressionValue:", util.GetStringValue(expression))
	switch s := expression.(type) {
	case *ast.TemplateLiteral:
		res, err = this_.TemplateLiteral(s)
		return
	case *ast.CallExpression:
		res, err = this_.CallExpression(s)
		return
	case *ast.StringLiteral:
		//res = s.Value.String()
		res = ValueTypeString
		return
	case *ast.NullLiteral:
		//res = s.Value.String()
		res = ValueTypeNull
		return
	}
	nameScript, err = this_.GetExpressionScript(expression)
	if err != nil {
		return
	}
	if nameScript == "" {
		return
	}
	v, err := this_.script.vm.RunString(nameScript)
	if err != nil {
		return
	}
	res = v.Export()
	util.Logger.Debug(this_.GetKey()+" GetExpressionValue nameScript ["+nameScript+"] ", zap.Any("value", res))

	return
}

func (this_ *CompilerMethod) GetExpressionForType(expression ast.Expression) (nameScript string, res *ValueType, err error) {

	//fmt.Println(fmt.Sprintf("TODO GetExpressionType code:%s,expression:%s", this_.GetNodeCode(expression), util.GetStringValue(expression)))
	switch s := expression.(type) {
	case *ast.TemplateLiteral:
		res = ValueTypeString
		this_.fullImport("fmt")
		return
	case *ast.StringLiteral:
		res = ValueTypeString
		return
	case *ast.NullLiteral:
		res = ValueTypeNull
		return
	case *ast.BinaryExpression:
		// TODO 表达式
		res = ValueTypeAny
		return
	case *ast.NumberLiteral:
		if _, ok := s.Value.(float64); ok {
			res = ValueTypeFloat64
		} else if _, ok := s.Value.(float32); ok {
			res = ValueTypeFloat32
		} else if _, ok := s.Value.(int64); ok {
			res = ValueTypeInt64
		} else if _, ok := s.Value.(int8); ok {
			res = ValueTypeInt8
		} else if _, ok := s.Value.(int16); ok {
			res = ValueTypeInt16
		} else if _, ok := s.Value.(int32); ok {
			res = ValueTypeInt32
		} else {
			res = ValueTypeInt
		}
		return
	case *ast.CallExpression:
		var v any
		v, err = this_.CallExpression(s)
		if err != nil {
			return
		}
		fmt.Println("GetExpressionForType CallExpression res:", util.GetStringValue(v))
		if v != nil {
			if vT, ok := v.(*ValueType); ok {
				res = vT
			} else if vT, ok := v.(*CompilerMethodResult); ok {
				res = vT.valueType
			} else if vT, ok := v.(*CompilerField); ok {
				res = vT.valueType
			} else {
				err = errors.New("GetExpressionForType CallExpression value [" + reflect.TypeOf(v).String() + "] not is ValueType")
				return
			}
		}
		return
	}
	nameScript, err = this_.GetExpressionScript(expression)
	if err != nil {
		return
	}
	if nameScript == "" {
		return
	}

	methodVar := this_.getParam(nameScript)
	if methodVar != nil {
		res = methodVar.valueType
		return
	}

	methodParam := this_.getParam(nameScript)
	if methodParam != nil {
		res = methodParam.valueType
		return
	}

	var v goja.Value
	v, err = this_.script.vm.RunString(nameScript)
	if err != nil {
		return
	}
	if v == goja.Undefined() {
		return
	}
	fmt.Println("TODO GetExpressionType key:", nameScript, ",v:", v)
	vv := v.Export()
	if vT, ok := vv.(*ValueType); ok {
		res = vT
	} else if vT, ok := vv.(*CompilerMethodVar); ok {
		res = vT.valueType
	} else if vT, ok := vv.(*CompilerField); ok {
		res = vT.valueType
	} else {
		err = errors.New("GetExpressionForType nameScript [" + nameScript + "] value [" + reflect.TypeOf(vv).String() + "] not is ValueType")
		return
	}
	fmt.Println("TODO GetExpressionType key:", nameScript, ",type:", res)

	return
}

func (this_ *CompilerMethod) GetExpressionScript(expression ast.Expression) (script string, err error) {

	//fmt.Println("TODO GetExpressionScript:", util.GetStringValue(expression))
	switch s := expression.(type) {
	case *ast.FunctionLiteral:
		err = this_.FunctionLiteral(s)
		break
	case *ast.Identifier:
		script = s.Name.String()
		break
	case *ast.StringLiteral:
		script = "\"" + s.Value.String() + "\""
		break
	case *ast.BracketExpression:
		var leftName string
		leftName, err = this_.GetExpressionScript(s.Left)
		if err != nil {
			return
		}
		var rightName string
		rightName, err = this_.GetExpressionScript(s.Member)
		if err != nil {
			return
		}
		script = leftName + "[" + rightName + "]"
		break
	case *ast.DotExpression:
		var leftName string
		leftName, err = this_.GetExpressionScript(s.Left)
		if err != nil {
			return
		}
		var rightName = s.Identifier.Name.String()
		script = leftName + "." + rightName
		break
	default:
		err = this_.Error("GetExpressionScript ["+reflect.TypeOf(s).String()+"] 不支持", expression)
		util.Logger.Debug(this_.GetKey()+" GetExpressionScript error", zap.Error(err))
		break
	}
	util.Logger.Debug(this_.GetKey()+" GetExpressionScript script ["+script+"] ", zap.Any("code", this_.GetNodeCode(expression)))

	return
}
