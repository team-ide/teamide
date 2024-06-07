package main

import (
	"encoding/json"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"os"
	"strings"
	"teamide/pkg/maker"
	"teamide/pkg/maker/modelers"
	"testing"
)

var (
	localDir string
)

func LoadDemoApp() (app *maker.Application, err error) {
	rootDir, err := os.Getwd()
	if err != nil {
		util.Logger.Error("os get wd error", zap.Error(err))
		return
	}

	rootDir = util.FormatPath(rootDir)
	localDir = rootDir
	if !strings.HasSuffix(rootDir, "/demo") {
		localDir = rootDir + "/demo"
	}
	exist, _ := util.PathExists(localDir)
	if !exist {
		localDir = rootDir + "/pkg/maker/demo"
	}

	fmt.Println("demo app load start dir:", localDir)
	app = maker.Load(localDir)
	fmt.Println("demo app load success")
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

	res, _, err := app.Save(modelers.TypeStorage, "user/getAllUsers", map[string]interface{}{
		"comment": "获取所有用户",
	}, false, false)
	if err != nil {
		util.Logger.Error("app save model error", zap.Error(err))
		return
	}
	println(util.GetStringValue(res))
	bs, err = json.MarshalIndent(app, "", "  ")
	if err != nil {
		util.Logger.Error("app to json error", zap.Error(err))
		return
	}
	println(string(bs))
}
