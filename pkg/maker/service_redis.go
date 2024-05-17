package maker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/redis"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
)

func (this_ *Invoker) GetRedisService() (res *ServiceRedis, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get redis service error", zap.Any("error", err))
		}
	}()
	res, err = this_.GetRedisServiceByName("")
	return
}

func (this_ *Invoker) GetRedisServiceByName(name string) (res *ServiceRedis, err error) {
	if name == "" {
		name = "default"
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get redis service by name error", zap.Any("name", name), zap.Any("error", err))
		}
	}()
	this_.redisServiceCacheLock.Lock()
	defer this_.redisServiceCacheLock.Unlock()
	res, find := this_.redisServiceCache[name]

	if find {
		return
	}

	util.Logger.Info("redis service not found,now create service", zap.Any("name", name))

	var config *modelers.ConfigRedisModel
	config = this_.app.GetConfigRedis(name)
	if config == nil {
		err = errors.New("config redis [" + name + "] is not exist")
		util.Logger.Error("create redis service error", zap.Any("error", err))
		return
	}
	ser, err := redis.New(
		&redis.Config{
			Address:  config.Address,
			Username: config.Username,
			Auth:     config.Auth,
			CertPath: config.CertPath,
		})
	if err != nil {
		util.Logger.Error("create redis service error", zap.Any("error", err))
		return
	}
	res = &ServiceRedis{
		ser: ser,
	}
	this_.redisServiceCache[name] = res

	var scriptVar = "redis"
	if name != "default" {
		scriptVar = "redis_" + name
	}
	err = this_.setScriptVar(scriptVar, res)

	return
}

type ServiceRedis struct {
	ser redis.IService
}

func (this_ *ServiceRedis) ShouldMappingFunc() bool {
	return true
}

func (this_ *ServiceRedis) Get(args ...interface{}) (res any, err error) {
	argLen := len(args)
	if argLen == 0 {
		err = errors.New("redis get args error")
		return
	}
	key := util.GetStringValue(args[0])
	util.Logger.Debug("redis get", zap.Any("key", key))

	v, err := this_.ser.Get(key)
	if err != nil {
		return
	}
	if v == "" {
		return
	}
	res = map[string]interface{}{}
	err = json.Unmarshal([]byte(v), &res)
	if err != nil {
		return
	}

	return
}

func (this_ *ServiceRedis) Set(args ...interface{}) (res any, err error) {
	argLen := len(args)
	if argLen < 2 {
		err = errors.New("redis set args error")
		return
	}
	key := util.GetStringValue(args[0])
	value := util.GetStringValue(args[1])
	util.Logger.Debug("redis set", zap.Any("key", key), zap.Any("value", value))

	err = this_.ser.Set(key, value)

	return
}
