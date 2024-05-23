package maker

import (
	"errors"
	"fmt"
	"github.com/dop251/goja/ast"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"reflect"
)

func (this_ *CompileProgram) Statements(info *CompileInfo, statements []ast.Statement) (err error) {
	fmt.Println("TODO Statements:", util.GetStringValue(statements))
	for _, statement := range statements {
		err = this_.Statement(info, statement)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *CompileProgram) Statement(info *CompileInfo, statement ast.Statement) (err error) {
	if statement == nil {
		return
	}
	fmt.Println("TODO Statement:", util.GetStringValue(statement))

	switch s := statement.(type) {
	case *ast.ExpressionStatement:
		err = this_.ExpressionStatement(info, s)
		break
	case *ast.IfStatement:
		err = this_.IfStatement(info, s)
		break
	case *ast.VariableStatement:
		err = this_.VariableStatement(info, s)
		break
	case *ast.BlockStatement:
		err = this_.BlockStatement(info, s)
		break
	case *ast.ThrowStatement:
		err = this_.ThrowStatement(info, s)
		break
	case *ast.ReturnStatement:
		err = this_.ReturnStatement(info, s)
		break
	default:
		err = errors.New("statement [" + reflect.TypeOf(statement).String() + "] not support")
		break

	}
	return
}

func (this_ *CompileProgram) BlockStatement(info *CompileInfo, statement *ast.BlockStatement) (err error) {
	fmt.Println("TODO BlockStatement:", util.GetStringValue(statement))
	err = this_.Statements(info, statement.List)
	if err != nil {
		return
	}
	return
}

func (this_ *CompileProgram) VariableStatement(info *CompileInfo, statement *ast.VariableStatement) (err error) {
	fmt.Println("TODO VariableStatement:", util.GetStringValue(statement))
	err = this_.Bindings(info, statement.List)
	if err != nil {
		return
	}
	return
}

func (this_ *CompileProgram) Bindings(info *CompileInfo, bindings []*ast.Binding) (err error) {
	fmt.Println("TODO Bindings:", util.GetStringValue(bindings))
	for _, binding := range bindings {
		err = this_.Binding(info, binding)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *CompileProgram) Binding(info *CompileInfo, binding *ast.Binding) (err error) {
	fmt.Println("TODO Binding:", util.GetStringValue(binding))
	nameScript, err := this_.GetExpressionScript(info, binding.Target)
	if err != nil {
		return
	}
	if info.findType(nameScript) {
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
	if varTypeStr != "" {
		var varType *ValueType
		varType, err = info.script.compiler.GetValueType(varTypeStr)
		if err != nil {
			return
		}
		util.Logger.Debug("Binding var set type", zap.Any("name", nameScript), zap.Any("type", varType))
		err = info.script.Set(nameScript, varType)
		if err != nil {
			return
		}
	} else {
		err = info.script.Set(nameScript, nil)
		if err != nil {
			return
		}
	}
	if binding.Initializer != nil {
		var v []*ValueType
		_, v, err = this_.GetExpressionForType(info, binding.Initializer)
		if err != nil {
			return
		}
		util.Logger.Debug("Binding var set initializer type", zap.Any("name", nameScript), zap.Any("type", v))
		info.addVarType(nameScript, v...)
		err = info.script.Set(nameScript, v)
		if err != nil {
			return
		}
	}
	return
}
func (this_ *CompileProgram) ExpressionStatement(info *CompileInfo, statement *ast.ExpressionStatement) (err error) {
	fmt.Println("TODO ExpressionStatement:", util.GetStringValue(statement))
	err = this_.Expression(info, statement.Expression)
	return
}
func (this_ *CompileProgram) ThrowStatement(info *CompileInfo, statement *ast.ThrowStatement) (err error) {
	fmt.Println("TODO ThrowStatement:", util.GetStringValue(statement))
	return
}

func (this_ *CompileProgram) IfStatement(info *CompileInfo, statement *ast.IfStatement) (err error) {
	fmt.Println("TODO IfStatement:", util.GetStringValue(statement))
	err = this_.Statement(info, statement.Consequent)
	if err != nil {
		return
	}

	err = this_.Statement(info, statement.Alternate)
	if err != nil {
		return
	}
	return
}

func (this_ *CompileProgram) ReturnStatement(info *CompileInfo, statement *ast.ReturnStatement) (err error) {
	fmt.Println("TODO ReturnStatement:", util.GetStringValue(statement))
	_, resV, err := this_.GetExpressionForType(info, statement.Argument)
	if err != nil {
		return
	}
	info.returnList = append(info.returnList, resV...)
	return
}
