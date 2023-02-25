package invokers

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
	"teamide/pkg/util"
)

func (this_ *Invoker) invokeRedisStep(step *modelers.StepRedisModel, invokeData *InvokeData) (ok bool, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("invoke redis step error", zap.Any("error", err))
		}
		util.Logger.Debug("invoke redis step end", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))
	}()
	util.Logger.Debug("invoke redis step start", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))

	var key string
	key, err = this_.GetNameByRule(step.Key, invokeData)
	if err != nil {
		util.Logger.Error("invoke redis step get key name error", zap.Any("key", step.Key), zap.Any("error", err))
		return
	}

	redisService, err := this_.GetRedisServiceByName(step.Datasource)
	if err != nil {
		util.Logger.Error("invoke redis step get redis service error", zap.Any("datasource", step.Datasource), zap.Any("error", err))
		return
	}
	var ctx = context.Background()

	var res interface{}

	switch step.GetType() {
	case modelers.RedisGet:

		if key == "" {
			err = errors.New("redis get key is empty")
			util.Logger.Error("invoke redis step error", zap.Any("error", err))
			return
		}
		util.Logger.Debug("invoke redis get start", zap.Any("key", key))
		res, err = redisService.Get(ctx, 0, key, 0, 0)
		if err != nil {
			util.Logger.Error("redis get error", zap.Any("datasource", step.Datasource), zap.Any("error", err))
			return
		}
		util.Logger.Debug("invoke redis get end", zap.Any("key", key), zap.Any("res", res))
		break
	case modelers.RedisSet:
		break
	}

	ok = true
	return
}
