package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-zookeeper/zk"
	_ "github.com/team-ide/go-tool/db/db_type_mysql"
	"github.com/team-ide/go-tool/util"
	"github.com/team-ide/go-tool/zookeeper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"teamide/pkg/maker"
	"teamide/pkg/maker/modelers"
	"testing"
	"time"
)

type defaultLogger struct{}

func (*defaultLogger) Printf(format string, args ...interface{}) {
	util.Logger.Info(fmt.Sprintf("zookeeper log:"+format, args...))
}
func TestInvokerUserGet(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("TestInvokerUserGet error", zap.Any("error", e))
		}
	}()
	zookeeper.ZKLogger = &defaultLogger{}
	zk.DefaultLogger = zookeeper.ZKLogger
	config := zap.NewDevelopmentConfig()
	//config.Encoding = "json"
	//config.EncoderConfig = zap.NewProductionEncoderConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(util.DefaultTimeFormatLayout + ".000")
	config.DisableStacktrace = true
	//config.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	util.Logger, _ = config.Build()
	util.Logger.Debug("NewDevelopmentConfig json Logger")

	util.Logger.Debug("TestInvokerUserGet start")

	compiler := LoadDemoCompiler()

	invoker, err := maker.NewInvoker(compiler)
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

func TestInvokerZk(t *testing.T) {
	compiler := LoadDemoCompiler()

	invoker, err := maker.NewInvoker(compiler)
	if err != nil {
		util.Logger.Error("NewInvoker error", zap.Error(err))
		return
	}

	invokeData, err := invoker.NewInvokeData()
	if err != nil {
		util.Logger.Error("NewInvokeData error", zap.Error(err))
		return
	}

	serviceName := "task/zk"
	res, err := invoker.InvokeServiceByName(serviceName, invokeData)
	if err != nil {
		util.Logger.Error("service invoke error", zap.Any("serviceName", serviceName), zap.Error(err))
		return
	}
	bs, err := json.Marshal(res)
	if err != nil {
		util.Logger.Error("res to json error", zap.Error(err))
		return
	}
	println("service [" + serviceName + "] run success")
	println(string(bs))

	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		time.Sleep(time.Second * 5)
		wait.Done()
	}()
	wait.Wait()
}
