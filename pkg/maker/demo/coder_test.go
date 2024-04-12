package main

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker/coders"
	"teamide/pkg/maker/coders/javascript"
	"testing"
)

func TestCoder(t *testing.T) {
	app, err := LoadDemoApp()
	if err != nil {
		util.Logger.Error("load demo app error", zap.Error(err))
		return
	}

	coderFactory := javascript.NewFactory(app)
	applicationCoder := coders.NewApplicationCoder(app, coderFactory)

	println("---------constant list start---------")
	for _, one := range app.ConstantList {
		println("---------constant [" + one.Name + "] start---------")
		code, err := applicationCoder.GenConstant(one)
		if err != nil {
			util.Logger.Error("gen constant code error", zap.Error(err))
			return
		}
		println("dir:" + code.Dir)
		println("name:" + code.Name)
		println("namespace:" + code.Namespace)
		println("codeType:" + string(code.CodeType))
		println(code.ToContent())
		println("---------constant [" + one.Name + "] end---------")
	}
	println("---------constant list end---------")

	//println("dao list start")
	//for _, one := range app.DaoList {
	//	println("dao [" + one.Name + "] start")
	//	js, err := javascript.GetDaoJavascript(app, one)
	//	if err != nil {
	//		util.Logger.Error("get dao javascript error", zap.Error(err))
	//		return
	//	}
	//	println(js)
	//	println("dao [" + one.Name + "] end")
	//}
	//println("dao list end")
	//
	//println("service list start")
	//for _, one := range app.ServiceList {
	//	println("service [" + one.Name + "] start")
	//	js, err := javascript.GetServiceJavascript(app, one)
	//	if err != nil {
	//		util.Logger.Error("get service javascript error", zap.Error(err))
	//		return
	//	}
	//	println(js)
	//	println("service [" + one.Name + "] end")
	//}
	//println("service list end")
}
