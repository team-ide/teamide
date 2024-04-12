package main

import (
	"encoding/json"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"os"
	"strings"
	"teamide/pkg/maker/modelers"
	"testing"
)

func LoadDemoApp() (app *modelers.Application, err error) {
	rootDir, err := os.Getwd()
	if err != nil {
		util.Logger.Error("os get wd error", zap.Error(err))
		return
	}

	rootDir = util.FormatPath(rootDir)
	dir := rootDir
	if !strings.HasSuffix(rootDir, "/demo") {
		dir = rootDir + "/demo"
	}
	exist, _ := util.PathExists(dir)
	if !exist {
		dir = rootDir + "/pkg/maker/demo"
	}

	fmt.Println("demo app dir:", dir)
	app = modelers.Load(dir)
	return
}

func TestLoader(t *testing.T) {
	app, err := LoadDemoApp()
	if err != nil {
		util.Logger.Error("load demo app error", zap.Error(err))
		return
	}

	bs, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		util.Logger.Error("app to json error", zap.Error(err))
		return
	}
	println(string(bs))
}
