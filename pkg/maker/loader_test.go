package maker

import (
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"teamide/pkg/maker/modelers"
	"teamide/pkg/util"
	"testing"
)

func LoadDemoApp() (app *modelers.Application, err error) {
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
	dir := rootDir + "/demo"
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
