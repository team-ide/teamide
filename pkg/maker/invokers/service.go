package invokers

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"teamide/pkg/db"
	"teamide/pkg/maker/modelers"
	"teamide/pkg/redis"
	"teamide/pkg/util"
	"teamide/pkg/zookeeper"
)

func (this_ *Invoker) GetRedisService() (res redis.Service, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get redis service error", zap.Any("error", err))
		}
	}()
	res, err = this_.GetRedisServiceByName("")
	return
}

func (this_ *Invoker) GetRedisServiceByName(name string) (res redis.Service, err error) {
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
	res, err = redis.CreateRedisService(config.Address, config.Username, config.Auth, config.CertPath)
	if err != nil {
		util.Logger.Error("create redis service error", zap.Any("error", err))
		return
	}
	return
}

func (this_ *Invoker) GetDbService() (res *db.Service, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get db service error", zap.Any("error", err))
		}
	}()
	res, err = this_.GetDbServiceByName("")
	return
}

func (this_ *Invoker) GetDbServiceByName(name string) (res *db.Service, err error) {
	if name == "" {
		name = "default"
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get db service by name error", zap.Any("name", name), zap.Any("error", err))
		}
	}()
	this_.dbServiceCacheLock.Lock()
	defer this_.dbServiceCacheLock.Unlock()
	res, find := this_.dbServiceCache[name]

	if find {
		return
	}

	util.Logger.Info("db service not found,now create service", zap.Any("name", name))

	var config *modelers.ConfigDbModel
	config = this_.app.GetConfigDb("default")
	if config == nil {
		err = errors.New("config db [" + name + "] is not exist")
		util.Logger.Error("create db service error", zap.Any("error", err))
		return
	}
	res, err = db.CreateService(&db.DatabaseConfig{
		Type:     config.Type,
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		Database: config.Database,
	})
	if err != nil {
		util.Logger.Error("create db service error", zap.Any("error", err))
		return
	}
	return
}

func (this_ *Invoker) GetZkService() (res zookeeper.Service, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get zk service error", zap.Any("error", err))
		}
	}()
	res, err = this_.GetZkServiceByName("")
	return
}

func (this_ *Invoker) GetZkServiceByName(name string) (res zookeeper.Service, err error) {
	if name == "" {
		name = "default"
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get zk service by name error", zap.Any("name", name), zap.Any("error", err))
		}
	}()
	this_.zkServiceCacheLock.Lock()
	defer this_.zkServiceCacheLock.Unlock()
	res, find := this_.zkServiceCache[name]

	if find {
		return
	}

	util.Logger.Info("zk service not found,now create service", zap.Any("name", name))

	var config *modelers.ConfigZkModel
	config = this_.app.GetConfigZk("default")
	if config == nil {
		err = errors.New("config zk [" + name + "] is not exist")
		util.Logger.Error("create zk service error", zap.Any("error", err))
		return
	}
	res, err = zookeeper.CreateZKService(zookeeper.Config{
		Address:  config.Address,
		Username: config.Username,
		Password: config.Password,
	})
	if err != nil {
		util.Logger.Error("create db service error", zap.Any("error", err))
		return
	}
	return
}
