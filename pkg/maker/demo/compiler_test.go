package main

import (
	"encoding/json"
	_ "github.com/team-ide/go-tool/db/db_type_mysql"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker"
	"teamide/pkg/maker/modelers"
	"testing"
)

func TestCompileUserGet(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("TestInvokerUserGet error", zap.Any("error", e))
		}
	}()

	util.Logger.Debug("TestInvokerUserGet start")

	app, err := LoadDemoApp()
	if err != nil {
		util.Logger.Error("load demo app error", zap.Error(err))
		return
	}

	invoker, err := maker.NewInvoker(app)
	if err != nil {
		util.Logger.Error("NewInvoker error", zap.Error(err))
		return
	}

	invokeData, err := invoker.NewInvokeData()
	if err != nil {
		util.Logger.Error("NewInvokeData error", zap.Error(err))
		return
	}

	err = invokeData.AddArg("userId", 1, modelers.ValueTypeInt64)
	if err != nil {
		util.Logger.Error("invoke data add arg error", zap.Error(err))
		return
	}

	serviceName := "user/get"
	startTime := util.GetNowMilli()
	res, err := invoker.InvokeServiceByName(serviceName, invokeData)
	if err != nil {
		util.Logger.Error("service invoke error", zap.Any("serviceName", serviceName), zap.Error(err))
		return
	}
	endTime := util.GetNowMilli()
	bs, err := json.Marshal(res)
	if err != nil {
		util.Logger.Error("res to json error", zap.Error(err))
		return
	}
	println("service ["+serviceName+"] run success,use", endTime-startTime, "ms")
	println(string(bs))
}
