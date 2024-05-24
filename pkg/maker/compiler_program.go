package maker

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja/ast"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
)

func (this_ *Script) CompileScript(script string) (compileProgram *CompileProgram, err error) {
	runScript := `(function (){
` + script + `
})()`
	astProgram, err := goja.Parse("", runScript)
	if err != nil {
		util.Logger.Error("compile script parse error", zap.Any("error", err))
		return
	}
	program, err := goja.CompileAST(astProgram, false)
	if err != nil {
		return
	}

	compileProgram = &CompileProgram{
		code:       runScript,
		program:    program,
		script:     this_,
		astProgram: astProgram,
	}

	return
}

type CompileProgram struct {
	script     *Script
	code       string
	program    *goja.Program
	astProgram *ast.Program
}

func (this_ *CompilerMethod) CompileMethod(method *CompilerMethod) (res *CompilerMethodResult, err error) {
	key := method.GetKey()
	res = method.result
	if res != nil {
		util.Logger.Debug("compile method ["+key+"] already compiled", zap.Any("result", res))
		return
	}

	res, err = method.Compile()
	if err != nil {
		return
	}
	return
}

func (this_ *CompilerMethod) Compile() (res *CompilerMethodResult, err error) {
	key := this_.GetKey()
	res = &CompilerMethodResult{}
	this_.result = res

	util.Logger.Debug("compile method [" + key + "] start")
	this_.script, err = this_.Compiler.script.NewScript()
	if err != nil {
		util.Logger.Error("compile method new script error", zap.Error(err))
		return
	}

	for _, param := range this_.paramList {
		util.Logger.Debug("compile method [" + key + "] param [" + param.name + "] init")
		err = this_.script.Set(param.name, param)
		if err != nil {
			util.Logger.Error("compile method ["+key+"] param ["+param.name+"] init error", zap.Error(err))
			return
		}
	}

	err = this_.VariableDeclarations(this_.program.DeclarationList)
	if err != nil {
		util.Logger.Error("compile method ["+key+"] variable declarations error", zap.Error(err))
		return
	}
	err = this_.Statements(this_.program.Body)
	if err != nil {
		util.Logger.Error("compile method ["+key+"] statements error", zap.Error(err))
		return
	}
	util.Logger.Debug("compile method ["+key+"] end", zap.Any("result", res))

	return
}

func (this_ *CompilerMethod) VariableDeclarations(variableDeclarations []*ast.VariableDeclaration) (err error) {

	fmt.Println("TODO VariableDeclarations:", util.GetStringValue(variableDeclarations))
	for _, variableDeclaration := range variableDeclarations {
		err = this_.VariableDeclaration(variableDeclaration)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *CompilerMethod) VariableDeclaration(variableDeclaration *ast.VariableDeclaration) (err error) {
	fmt.Println("TODO VariableDeclaration:", util.GetStringValue(variableDeclaration))
	return
}

func (this_ *CompilerMethod) FunctionLiteral(expression *ast.FunctionLiteral) (err error) {
	err = this_.ParameterList(expression.ParameterList)
	if err != nil {
		return
	}
	err = this_.VariableDeclarations(expression.DeclarationList)
	if err != nil {
		return
	}
	err = this_.BlockStatement(expression.Body)
	if err != nil {
		return
	}
	return
}

func (this_ *CompilerMethod) ParameterList(parameterList *ast.ParameterList) (err error) {
	fmt.Println("TODO ParameterList:", util.GetStringValue(parameterList))
	return
}
