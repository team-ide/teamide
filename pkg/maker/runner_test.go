package maker

import (
	"encoding/json"
	"go.uber.org/zap"
	"teamide/pkg/util"
	"testing"
)

func TestRunner(t *testing.T) {
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
}
