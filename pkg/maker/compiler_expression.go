package maker

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja/ast"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"reflect"
	"teamide/pkg/maker/modelers"
)

type ExpressionInfo struct {
	Value any
}

func (this_ *CompileProgram) Expression(info *CompileInfo, expression ast.Expression) (err error) {
	if expression == nil {
		return
	}
	fmt.Println("TODO Expression:", util.GetStringValue(expression))
	info.expressionInfo = &ExpressionInfo{}
	switch e := expression.(type) {
	case *ast.CallExpression:
		info.expressionInfo.Value, err = this_.CallExpression(info, e)
		break
	case *ast.FunctionLiteral:
		err = this_.FunctionLiteral(info, e)
		break
	case *ast.AssignExpression:
		err = this_.AssignExpression(info, e)
		break
	case *ast.TemplateLiteral:
		info.expressionInfo.Value, err = this_.TemplateLiteral(info, e)
		break
	default:
		err = errors.New("expression [" + reflect.TypeOf(expression).String() + "] not support")
		break

	}
	return
}

func (this_ *CompileProgram) ArgumentList(info *CompileInfo, argumentList []ast.Expression) (values []interface{}, err error) {
	fmt.Println("TODO ArgumentList:", util.GetStringValue(argumentList))
	var v []*modelers.ValueType
	for _, arg := range argumentList {
		_, v, err = this_.GetExpressionForType(info, arg)
		if err != nil {
			return
		}
		if len(v) == 0 {
			values = append(values, nil)
		} else if len(v) == 1 {
			values = append(values, v[0])
		} else {
			values = append(values, v)
		}
	}
	return
}

func (this_ *CompileProgram) CallExpression(info *CompileInfo, expression *ast.CallExpression) (res any, err error) {
	fmt.Println("TODO CallExpression ArgumentList:", util.GetStringValue(expression.ArgumentList))
	args, err := this_.ArgumentList(info, expression.ArgumentList)
	if err != nil {
		return
	}
	fmt.Println("TODO CallExpression Callee:", util.GetStringValue(expression.Callee))

	nameScript, v, err := this_.GetExpressionForValue(info, expression.Callee)
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
	fmt.Println("CallExpression GetExpressionForValue res:", util.GetStringValue(v))

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
	case *modelers.DaoModel:
		var in *CompileInfo
		in, err = this_.script.compiler.CompileDao(m)
		if err != nil {
			return
		}
		res = in.returnList
		break
	case *modelers.ServiceModel:
		var in *CompileInfo
		in, err = this_.script.compiler.CompileService(m)
		if err != nil {
			return
		}
		res = in.returnList
		break
	case *modelers.ValueType:
		res = m
		break
	case []*modelers.ValueType:
		res = m
		break
	default:
		err = errors.New("call expression method [" + reflect.TypeOf(v).String() + "] not support")
		break
	}
	return
}

func (this_ *CompileProgram) AssignExpression(info *CompileInfo, expression *ast.AssignExpression) (err error) {

	fmt.Println("TODO AssignExpression Left:", util.GetStringValue(expression.Left))
	nameScript, err := this_.GetExpressionScript(info, expression.Left)
	if err != nil {
		return
	}
	fmt.Println("AssignExpression Right:", util.GetStringValue(expression.Right))
	_, v, err := this_.GetExpressionForType(info, expression.Right)
	if err != nil {
		return
	}

	info.addVarType(nameScript, v...)
	err = info.script.Set(nameScript, v)
	util.Logger.Debug("AssignExpression var set", zap.Any("name", nameScript), zap.Any("type", v))
	if err != nil {
		return
	}
	return
}

func (this_ *CompileProgram) TemplateLiteral(info *CompileInfo, expression *ast.TemplateLiteral) (res any, err error) {
	fmt.Println("TODO TemplateLiteral:", util.GetStringValue(expression))
	res = modelers.ValueTypeString
	return
}

func (this_ *CompileProgram) GetExpressionForValue(info *CompileInfo, expression ast.Expression) (nameScript string, res any, err error) {

	fmt.Println("TODO GetExpressionValue:", util.GetStringValue(expression))
	switch s := expression.(type) {
	case *ast.TemplateLiteral:
		res, err = this_.TemplateLiteral(info, s)
		return
	case *ast.CallExpression:
		res, err = this_.CallExpression(info, s)
		return
	case *ast.StringLiteral:
		//res = s.Value.String()
		res = modelers.ValueTypeString
		return
	case *ast.NullLiteral:
		//res = s.Value.String()
		res = modelers.ValueTypeNull
		return
	}
	nameScript, err = this_.GetExpressionScript(info, expression)
	if err != nil {
		return
	}
	if nameScript == "" {
		return
	}
	v, err := info.script.vm.RunString(nameScript)
	if err != nil {
		return
	}
	res = v.Export()
	fmt.Println("TODO GetExpressionValue key:", nameScript, ",value:", res)

	return
}

func (this_ *CompileProgram) GetExpressionForType(info *CompileInfo, expression ast.Expression) (nameScript string, res []*modelers.ValueType, err error) {

	fmt.Println("TODO GetExpressionType:", util.GetStringValue(expression))
	switch s := expression.(type) {
	case *ast.TemplateLiteral:
		res = append(res, modelers.ValueTypeString)
		return
	case *ast.StringLiteral:
		res = append(res, modelers.ValueTypeString)
		return
	case *ast.NullLiteral:
		res = append(res, modelers.ValueTypeNull)
		return
	case *ast.CallExpression:
		var v any
		v, err = this_.CallExpression(info, s)
		if err != nil {
			return
		}
		fmt.Println("GetExpressionForType CallExpression res:", util.GetStringValue(v))
		if v != nil {
			if vT, ok := v.(*modelers.ValueType); ok {
				res = append(res, vT)
			} else if vTs, ok := v.([]*modelers.ValueType); ok {
				res = append(res, vTs...)
			} else {
				err = errors.New("GetExpressionForType CallExpression value [" + reflect.TypeOf(v).String() + "] not is ValueType")
				return
			}
		}
		return
	}
	nameScript, err = this_.GetExpressionScript(info, expression)
	if err != nil {
		return
	}
	if nameScript == "" {
		return
	}
	res = info.varTypes[nameScript]
	if res == nil {
		var v goja.Value
		v, err = info.script.vm.RunString(nameScript)
		if err != nil {
			return
		}
		vv := v.Export()
		if vT, ok := vv.(*modelers.ValueType); ok {
			res = append(res, vT)
		} else if vTs, ok := vv.([]*modelers.ValueType); ok {
			res = append(res, vTs...)
		} else {
			err = errors.New("GetExpressionForType value [" + reflect.TypeOf(v).String() + "] not is ValueType")
			return
		}
	}
	fmt.Println("TODO GetExpressionType key:", nameScript, ",type:", res)

	return
}

func (this_ *CompileProgram) GetExpressionScript(info *CompileInfo, expression ast.Expression) (script string, err error) {

	switch s := expression.(type) {
	case *ast.FunctionLiteral:
		err = this_.FunctionLiteral(info, s)
		break
	case *ast.Identifier:
		script = s.Name.String()
		break
	case *ast.DotExpression:
		var leftName string
		leftName, err = this_.GetExpressionScript(info, s.Left)
		if err != nil {
			return
		}
		var rightName = s.Identifier.Name.String()
		script = leftName + "." + rightName
		break
	default:
		err = errors.New("GetExpressionScript [" + reflect.TypeOf(s).String() + "] not support")
		break
	}

	return
}
