package maker

import (
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"teamide/pkg/util"
	"testing"
)

func TestLoader(t *testing.T) {
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
	app := Load(dir)

	bs, err := json.Marshal(app)
	if err != nil {
		util.Logger.Error("app to json error", zap.Error(err))
		return
	}
	println(string(bs))
}
