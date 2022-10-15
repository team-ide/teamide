package maker

import (
	"go.uber.org/zap"
	"teamide/pkg/maker/coders/javascript"
	"teamide/pkg/util"
	"testing"
)

func TestCoder(t *testing.T) {
	app, err := LoadDemoApp()
	if err != nil {
		util.Logger.Error("load demo app error", zap.Error(err))
		return
	}

	//coderFactory := javascript.NewFactory(app)
	//applicationCoder := coders.NewApplicationCoder(app, coderFactory)

	println("dao list start")
	for _, one := range app.DaoList {
		println("dao [" + one.Name + "] start")
		js, err := javascript.GetDaoJavascript(app, one)
		if err != nil {
			util.Logger.Error("get dao javascript error", zap.Error(err))
			return
		}
		println(js)
		println("dao [" + one.Name + "] end")
	}
	println("dao list end")

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
