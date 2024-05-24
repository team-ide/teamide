package main

import (
	"fmt"
	_ "github.com/team-ide/go-tool/db/db_type_mysql"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker"
	"testing"
)

func LoadDemoCompiler() *maker.Compiler {
	app, err := LoadDemoApp()
	fmt.Println("load demo compiler start")
	if err != nil {
		util.Logger.Error("load demo app error", zap.Error(err))
		panic(err)
	}

	compiler, err := maker.NewCompiler(app)
	if err != nil {
		util.Logger.Error("NewCompiler error", zap.Error(err))
		panic(err)
	}
	fmt.Println("load demo compiler success")
	return compiler
}

func TestCompile(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("TestCompile error", zap.Any("error", e))
		}
	}()

	util.Logger.Debug("TestCompile start")

	compiler := LoadDemoCompiler()

	//fmt.Println(util.GetStringValue(compiler))
	compileErrs := compiler.Compile(false)
	for _, compileErr := range compileErrs {
		fmt.Println("method:", compileErr.Method.GetKey())
		fmt.Println("error:", compileErr.Err.Error())
	}

	util.Logger.Debug("TestCompile end")

}
