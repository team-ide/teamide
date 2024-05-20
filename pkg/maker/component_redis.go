package maker

import (
	"encoding/json"
	"errors"
	"github.com/team-ide/go-tool/redis"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
)

func NewRedisCompiler(config *modelers.ConfigRedisModel) *Component {
	component := &Component{}
	return component
}

func NewComponentRedis(config *modelers.ConfigRedisModel) (res *ComponentRedis, err error) {
	ser, err := redis.New(&redis.Config{
		Address:  config.Address,
		Username: config.Username,
		Auth:     config.Auth,
		CertPath: config.CertPath,
	})
	if err != nil {
		return
	}
	res = &ComponentRedis{
		ser: ser,
	}
	return
}

type ComponentRedis struct {
	ser redis.IService
}

func (this_ *ComponentRedis) ShouldMappingFunc() bool {
	return true
}

func (this_ *ComponentRedis) Get(args ...interface{}) (res any, err error) {
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

func (this_ *ComponentRedis) Set(args ...interface{}) (res any, err error) {
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
