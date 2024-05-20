package maker

import (
	"github.com/dop251/goja"
)

func (this_ *Script) CompileScript(script string) (*CompileProgram, error) {
	runScript := `(function (){` + script + `})()
`
	p, err := goja.Compile("", runScript, false)
	if err != nil {
		return nil, err
	}

	compileProgram := &CompileProgram{
		code:    runScript,
		program: p,
		script:  this_,
	}

	return compileProgram, nil
}

type CompileProgram struct {
	script  *Script
	code    string
	program *goja.Program
}

type CompileResult struct {
}
