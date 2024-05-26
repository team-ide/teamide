package main

import (
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker/coder"
	"teamide/pkg/maker/coder/golang"
	"testing"
)

func TestXxxSx(t *testing.T) {
	var bs = []string{"1", "2", "3"}
	bs2 := bs[0:2]
	bs[0] = "--"
	fmt.Println(bs)
	fmt.Println(bs2)
}

func TestCoder(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("TestCoder error", zap.Any("error", e))
		}
	}()

	util.Logger.Debug("TestCoder start")

	compiler := LoadDemoCompiler()

	options := &coder.Options{
		Dir: compiler.GetDir() + "gen-golang",
	}

	coder_, err := coder.NewCoder(compiler, options)
	if err != nil {
		panic(err)
	}

	err = golang.FullGenerator(coder_)
	if err != nil {
		panic(err)
	}

	err = coder_.Gen()

	if err != nil {
		panic(err)
	}
	util.Logger.Debug("TestCoder end")

}
