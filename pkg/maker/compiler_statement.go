package maker

import (
	"errors"
	"fmt"
	"github.com/dop251/goja/ast"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"reflect"
)

func (this_ *CompilerMethod) Statements(statements []ast.Statement) (err error) {
	fmt.Println("TODO Statements:", util.GetStringValue(statements))
	for _, statement := range statements {
		err = this_.Statement(statement)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *CompilerMethod) Statement(statement ast.Statement) (err error) {
	if statement == nil {
		return
	}
	fmt.Println("TODO Statement:", util.GetStringValue(statement))

	switch s := statement.(type) {
	case *ast.ExpressionStatement:
		err = this_.ExpressionStatement(s)
		break
	case *ast.IfStatement:
		err = this_.IfStatement(s)
		break
	case *ast.VariableStatement:
		err = this_.VariableStatement(s)
		break
	case *ast.BlockStatement:
		err = this_.BlockStatement(s)
		break
	case *ast.ThrowStatement:
		err = this_.ThrowStatement(s)
		break
	case *ast.ReturnStatement:
		err = this_.ReturnStatement(s)
		break
	default:
		err = errors.New("statement [" + reflect.TypeOf(statement).String() + "] not support")
		break
	}
	return
}

func (this_ *CompilerMethod) BlockStatement(statement *ast.BlockStatement) (err error) {
	fmt.Println("TODO BlockStatement:", util.GetStringValue(statement))
	err = this_.Statements(statement.List)
	if err != nil {
		return
	}
	return
}

func (this_ *CompilerMethod) VariableStatement(statement *ast.VariableStatement) (err error) {
	fmt.Println("TODO VariableStatement:", util.GetStringValue(statement))
	err = this_.Bindings(statement.List)
	if err != nil {
		return
	}
	return
}

func (this_ *CompilerMethod) Bindings(bindings []*ast.Binding) (err error) {
	fmt.Println("TODO Bindings:", util.GetStringValue(bindings))
	for _, binding := range bindings {
		err = this_.Binding(binding)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *CompilerMethod) Binding(binding *ast.Binding) (err error) {
	fmt.Println("TODO Binding:", util.GetStringValue(binding))
	nameScript, err := this_.GetExpressionScript(binding.Target)
	if err != nil {
		return
	}
	if this_.findType(nameScript) {
		err = errors.New("变量[" + nameScript + "]已定义")
		return
	}
	var varTypeStr string
	if binding.Type != nil {
		for i, t := range binding.Type {
			if i > 0 {
				varTypeStr += "."
			}
			varTypeStr += t.Name.String()
		}
	}
	methodVar := this_.addVar(nameScript)
	if varTypeStr != "" {
		var valueType *ValueType
		valueType, err = this_.GetValueType(varTypeStr)
		if err != nil {
			return
		}
		methodVar.addValueType(valueType)

		util.Logger.Debug("Binding var set type", zap.Any("name", nameScript), zap.Any("type", valueType))
	}
	err = this_.script.Set(nameScript, methodVar)
	if err != nil {
		return
	}
	if binding.Initializer != nil {
		var v []*ValueType
		_, v, err = this_.GetExpressionForType(binding.Initializer)
		if err != nil {
			return
		}
		util.Logger.Debug("Binding var set initializer type", zap.Any("name", nameScript), zap.Any("type", v))
		methodVar.addValueType(v...)
	}
	return
}
func (this_ *CompilerMethod) ExpressionStatement(statement *ast.ExpressionStatement) (err error) {
	fmt.Println("TODO ExpressionStatement:", util.GetStringValue(statement))
	err = this_.Expression(statement.Expression)
	return
}
func (this_ *CompilerMethod) ThrowStatement(statement *ast.ThrowStatement) (err error) {
	fmt.Println("TODO ThrowStatement:", util.GetStringValue(statement))
	return
}

func (this_ *CompilerMethod) IfStatement(statement *ast.IfStatement) (err error) {
	fmt.Println("TODO IfStatement:", util.GetStringValue(statement))
	err = this_.Statement(statement.Consequent)
	if err != nil {
		return
	}

	err = this_.Statement(statement.Alternate)
	if err != nil {
		return
	}
	return
}

func (this_ *CompilerMethod) ReturnStatement(statement *ast.ReturnStatement) (err error) {
	fmt.Println("TODO ReturnStatement:", util.GetStringValue(statement))
	_, resV, err := this_.GetExpressionForType(statement.Argument)
	if err != nil {
		return
	}
	this_.result.addValueType(resV...)
	return
}
