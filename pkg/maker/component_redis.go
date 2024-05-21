package maker

import (
	"errors"
	"github.com/team-ide/go-tool/redis"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
)

func NewRedisCompiler(config *modelers.ConfigRedisModel) *Component {
	component := &Component{
		Methods: []*ComponentMethod{
			{
				Name: "Get", GetReturnTypes: func(args []interface{}) (returnTypes []*modelers.ValueType) {
					if len(args) == 2 {
						returnTypes = append(returnTypes, args[1].(*modelers.ValueType))
					} else {
						returnTypes = append(returnTypes, modelers.ValueTypeString)
					}
					return
				},
			},
			{
				Name: "Set", GetReturnTypes: func(args []interface{}) (returnTypes []*modelers.ValueType) {
					return
				},
			},
		},
	}
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
	*Compiler
	ser redis.IService
}

func (this_ *ComponentRedis) ShouldMappingFunc() bool {
	return true
}

func (this_ *ComponentRedis) Get(key string, valueType *modelers.ValueType) (res any, err error) {
	workInfo := "redis get "
	if key == "" {
		err = errors.New(workInfo + "key is empty")
		return
	}
	util.Logger.Debug(workInfo, zap.Any("key", key))

	v, err := this_.ser.Get(key)
	if err != nil {
		return
	}
	if v == "" {
		return
	}
	res, err = this_.ToValueByValueType(v, valueType)
	if err != nil {
		return
	}

	return
}

func (this_ *ComponentRedis) Set(key string, value interface{}) (err error) {
	workInfo := "redis set "
	if key == "" {
		err = errors.New(workInfo + "key is empty")
		return
	}
	valueStr := util.GetStringValue(value)
	util.Logger.Debug(workInfo+"", zap.Any("key", key), zap.Any("value", valueStr))

	err = this_.ser.Set(key, valueStr)

	return
}
