package maker

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja/ast"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
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

type CompileInfo struct {
	from           string
	script         *Script
	expressionInfo *ExpressionInfo
	returnList     []*ValueType
	paramTypes     map[string][]*ValueType
	varTypes       map[string][]*ValueType
}

func (this_ *CompileInfo) findType(name string) (find bool) {
	_, find = this_.varTypes[name]
	if !find {
		_, find = this_.paramTypes[name]
	}
	return
}

func (this_ *CompileInfo) getType(name string) (types []*ValueType) {
	types = append(types, this_.varTypes[name]...)
	if len(types) == 0 {
		types = append(types, this_.paramTypes[name]...)
	}
	return
}

func (this_ *CompileInfo) addParamType(name string, types ...*ValueType) {
	for _, t := range types {
		var find bool
		for _, t_ := range this_.paramTypes[name] {
			if t_ == t {
				find = true
				break
			}
		}
		if !find {
			this_.varTypes[name] = append(this_.paramTypes[name], t)
			if t.FieldTypes != nil {
				for n, nT := range t.FieldTypes {
					this_.addParamType(name+"."+n, nT)
				}
			}
		}
	}

	return
}

func (this_ *CompileInfo) addVarType(name string, types ...*ValueType) {
	for _, t := range types {
		var find bool
		for _, t_ := range this_.varTypes[name] {
			if t_ == t {
				find = true
				break
			}
		}
		if !find {
			this_.varTypes[name] = append(this_.varTypes[name], t)
			if t.FieldTypes != nil {
				for n, nT := range t.FieldTypes {
					this_.addVarType(name+"."+n, nT)
				}
			}
		}
	}

	return
}

func (this_ *CompileProgram) Compile(from string, args []*modelers.ArgModel) (info *CompileInfo, err error) {

	util.Logger.Debug("compile start", zap.Any("from", from))

	script, err := this_.script.NewScript()
	if err != nil {
		return
	}

	info = &CompileInfo{
		from:       from,
		script:     script,
		paramTypes: make(map[string][]*ValueType),
		varTypes:   make(map[string][]*ValueType),
	}

	for _, arg := range args {
		var t *ValueType
		t, err = script.compiler.GetValueType(arg.Type)
		if err != nil {
			return
		}
		info.addParamType(arg.Name, t)
		err = script.Set(arg.Name, t)
		if err != nil {
			return
		}
	}

	err = this_.VariableDeclarations(info, this_.astProgram.DeclarationList)
	if err != nil {
		return
	}
	err = this_.Statements(info, this_.astProgram.Body)
	if err != nil {
		return
	}
	util.Logger.Debug("compile end", zap.Any("from", from), zap.Any("returnList", info.returnList))

	return
}

func (this_ *CompileProgram) VariableDeclarations(info *CompileInfo, variableDeclarations []*ast.VariableDeclaration) (err error) {

	fmt.Println("TODO VariableDeclarations:", util.GetStringValue(variableDeclarations))
	for _, variableDeclaration := range variableDeclarations {
		err = this_.VariableDeclaration(info, variableDeclaration)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *CompileProgram) VariableDeclaration(info *CompileInfo, variableDeclaration *ast.VariableDeclaration) (err error) {
	fmt.Println("TODO VariableDeclaration:", util.GetStringValue(variableDeclaration))
	return
}

func (this_ *CompileProgram) FunctionLiteral(info *CompileInfo, expression *ast.FunctionLiteral) (err error) {
	err = this_.ParameterList(info, expression.ParameterList)
	if err != nil {
		return
	}
	err = this_.VariableDeclarations(info, expression.DeclarationList)
	if err != nil {
		return
	}
	err = this_.BlockStatement(info, expression.Body)
	if err != nil {
		return
	}
	return
}

func (this_ *CompileProgram) ParameterList(info *CompileInfo, parameterList *ast.ParameterList) (err error) {
	fmt.Println("TODO ParameterList:", util.GetStringValue(parameterList))
	return
}
