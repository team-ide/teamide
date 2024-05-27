package maker

import (
	"fmt"
	"github.com/dop251/goja/ast"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"reflect"
)

func (this_ *CompilerMethod) Statements(statements []ast.Statement) (err error) {
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
		err = this_.Error("statement ["+reflect.TypeOf(statement).String()+"] 不支持", statement)
		util.Logger.Debug(this_.GetKey()+" Statement error", zap.Error(err))
		break
	}
	return
}

func (this_ *CompilerMethod) BlockStatement(statement *ast.BlockStatement) (err error) {
	err = this_.Statements(statement.List)
	if err != nil {
		return
	}
	return
}

func (this_ *CompilerMethod) VariableStatement(statement *ast.VariableStatement) (err error) {
	err = this_.Bindings(statement.List)
	if err != nil {
		return
	}
	return
}

func (this_ *CompilerMethod) VariableDeclaration(variableDeclaration *ast.VariableDeclaration) (err error) {
	fmt.Println("TODO VariableDeclaration:", util.GetStringValue(variableDeclaration))
	fmt.Println("TODO VariableDeclaration code:", this_.GetNodeCode(variableDeclaration))
	return
}

func (this_ *CompilerMethod) Bindings(bindings []*ast.Binding) (err error) {
	for _, binding := range bindings {
		err = this_.Binding(binding)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *CompilerMethod) Binding(binding *ast.Binding) (err error) {
	nameScript, err := this_.GetExpressionScript(binding.Target)
	if err != nil {
		return
	}
	if this_.FindType(nameScript) {
		err = this_.Error("变量 ["+nameScript+"] 已定义", binding)
		util.Logger.Error(this_.GetKey()+" Binding error", zap.Error(err))
		return
	}
	this_.BindingScriptCache[binding] = nameScript
	var varTypeStr string
	if binding.Type != nil {
		for i, t := range binding.Type {
			if i > 0 {
				varTypeStr += "."
			}
			varTypeStr += t.Name.String()
		}
	}
	var valueType *ValueType
	if varTypeStr != "" {
		valueType, err = this_.GetValueType(varTypeStr)
		if err != nil {
			return
		}
	}
	if binding.Initializer != nil {
		_, valueType, err = this_.GetExpressionForType(binding.Initializer)
		if err != nil {
			return
		}
	}
	methodVar := this_.addVar(nameScript, valueType)

	util.Logger.Debug("Binding var ["+nameScript+"]", zap.Any("type", valueType))
	err = this_.script.Set(nameScript, methodVar)
	if err != nil {
		return
	}

	this_.BindingCache[binding] = methodVar
	return
}
func (this_ *CompilerMethod) ExpressionStatement(statement *ast.ExpressionStatement) (err error) {
	err = this_.Expression(statement.Expression)
	return
}
func (this_ *CompilerMethod) ThrowStatement(statement *ast.ThrowStatement) (err error) {
	this_.fullImport("error")
	fmt.Println("TODO ThrowStatement:", util.GetStringValue(statement))
	return
}

func (this_ *CompilerMethod) IfStatement(statement *ast.IfStatement) (err error) {
	err = this_.Expression(statement.Test)
	if err != nil {
		return
	}
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
	_, resV, err := this_.GetExpressionForType(statement.Argument)
	if err != nil {
		return
	}
	err = this_.Result.addValueTypes(resV)
	if err != nil {
		return
	}
	return
}
