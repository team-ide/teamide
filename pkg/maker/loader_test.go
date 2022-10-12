package maker

import (
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"teamide/pkg/maker/code/javascript"
	"teamide/pkg/maker/model"
	"teamide/pkg/util"
	"testing"
)

func LoadDemoApp() (app *model.Application, err error) {
	rootDir, err := os.Getwd()
	if err != nil {
		util.Logger.Error("os get wd error", zap.Error(err))
		return
	}

	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		util.Logger.Error("filepath abs error", zap.Error(err))
		return
	}
	dir := rootDir + "/model/demo"
	app = model.Load(dir)
	return
}

func TestLoader(t *testing.T) {
	app, err := LoadDemoApp()
	if err != nil {
		util.Logger.Error("load demo app error", zap.Error(err))
		return
	}

	bs, err := json.Marshal(app)
	if err != nil {
		util.Logger.Error("app to json error", zap.Error(err))
		return
	}
	println(string(bs))

	println("service list start")
	for _, one := range app.ServiceList {
		println("service [" + one.Name + "] start")
		js, err := javascript.GetServiceJavascript(app, one)
		if err != nil {
			util.Logger.Error("get service javascript error", zap.Error(err))
			return
		}
		println(js)
		println("service [" + one.Name + "] end")
	}
	println("service list end")
}
