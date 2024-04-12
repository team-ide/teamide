package invokers

import (
	"context"
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/redis"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
	"time"
)

func (this_ *Invoker) invokeRedisStep(from string, step *modelers.StepRedisModel, invokeData *InvokeData) (err error) {
	funcInvoke := invokeStart(from+" redis", invokeData)

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(funcInvoke.name + " error:" + fmt.Sprint(e))
			util.Logger.Error(funcInvoke.name+" error", zap.Any("error", err))
		}
		funcInvoke.end(err)
		util.Logger.Debug(funcInvoke.name+" end", zap.Any("use", funcInvoke.use()))
	}()
	util.Logger.Debug(funcInvoke.name+" start", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))

	redisService, err := this_.GetRedisServiceByName(step.Config)
	if err != nil {
		util.Logger.Error(funcInvoke.name+" get redis service error", zap.Any("config", step.Config), zap.Any("error", err))
		return
	}
	param := &redis.Param{
		Ctx: context.Background(),
	}

	var varValue interface{}

	switch step.GetType() {
	case modelers.RedisGet:
		varValue, err = this_.invokeRedisGet(funcInvoke, step, invokeData, redisService, param)
		if err != nil {
			return
		}
		break
	case modelers.RedisSet:
		varValue, err = this_.invokeRedisSet(funcInvoke, step, invokeData, redisService, param)
		if err != nil {
			return
		}
		break
	default:
		err = errors.New(funcInvoke.name + " [" + step.Redis + "] can not be support")
		util.Logger.Error(funcInvoke.name+" error", zap.Any("error", err))
		return
	}
	if step.SetVar != "" {
		util.Logger.Debug(funcInvoke.name+" set var", zap.Any("setVar", step.SetVar), zap.Any("setVarType", step.SetVarType), zap.Any("varValue", varValue))
		err = invokeData.AddVar(step.SetVar, varValue, step.SetVarType)
		if err != nil {
			util.Logger.Error(funcInvoke.name+" set var error", zap.Any("setVar", step.SetVar), zap.Any("varValue", varValue), zap.Any("error", err))
			return
		}
	}

	return
}

func (this_ *Invoker) invokeRedisGet(funcInvoke *FuncInvoke, step *modelers.StepRedisModel, invokeData *InvokeData, redisService redis.IService, param *redis.Param) (res interface{}, err error) {

	funcInvoke.name = funcInvoke.name + " get"

	var key string
	key, err = invokeData.InvokeByStringRule(step.Key)
	if err != nil {
		util.Logger.Error(funcInvoke.name+" get key name error", zap.Any("key", step.Key), zap.Any("error", err))
		return
	}

	if key == "" {
		err = errors.New(funcInvoke.name + " key is empty")
		util.Logger.Error(funcInvoke.name+" error", zap.Any("error", err))
		return
	}
	startTime := time.Now()
	res, err = redisService.Get(key, param)
	if err != nil {
		util.Logger.Error(funcInvoke.name+" error", zap.Any("error", err))
		return
	}
	endTime := time.Now()
	util.Logger.Debug(funcInvoke.name + " use " + GetDurationFormatByMillisecond(endTime.UnixMilli()-startTime.UnixMilli()))

	return
}

func (this_ *Invoker) invokeRedisSet(funcInvoke *FuncInvoke, step *modelers.StepRedisModel, invokeData *InvokeData, redisService redis.IService, param *redis.Param) (res interface{}, err error) {

	funcInvoke.name = funcInvoke.name + " set"

	var key string
	key, err = invokeData.InvokeByStringRule(step.Key)
	if err != nil {
		util.Logger.Error(funcInvoke.name+" get key name error", zap.Any("key", step.Key), zap.Any("error", err))
		return
	}

	if key == "" {
		err = errors.New(funcInvoke.name + " key is empty")
		util.Logger.Error(funcInvoke.name+" error", zap.Any("error", err))
		return
	}
	var value interface{}
	value, err = invokeData.InvokeScript(step.Value)
	if err != nil {
		util.Logger.Error(funcInvoke.name+" value error", zap.Any("value", step.Value), zap.Any("error", err))
		return
	}
	var setValue string
	setValue = util.GetStringValue(value)
	err = redisService.Set(key, setValue, param)
	if err != nil {
		util.Logger.Error(funcInvoke.name+" error", zap.Any("error", err))
		return
	}

	return
}
