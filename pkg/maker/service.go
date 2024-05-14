package maker

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/db"
	"github.com/team-ide/go-tool/redis"
	"github.com/team-ide/go-tool/util"
	"github.com/team-ide/go-tool/zookeeper"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
)

func (this_ *Invoker) GetRedisService() (res redis.IService, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get redis service error", zap.Any("error", err))
		}
	}()
	res, err = this_.GetRedisServiceByName("")
	return
}

func (this_ *Invoker) GetRedisServiceByName(name string) (res redis.IService, err error) {
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
	res, err = redis.New(
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
	_, err = res.Get("_")
	if err != nil {
		util.Logger.Error("create redis service error", zap.Any("error", err))
		return
	}
	this_.redisServiceCache[name] = res

	var scriptVar = "redis"
	if name != "default" {
		scriptVar = "redis_" + name
	}
	err = this_.setScriptVar(scriptVar, res)

	return
}

func (this_ *Invoker) GetDbService() (res db.IService, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get db service error", zap.Any("error", err))
		}
	}()
	res, err = this_.GetDbServiceByName("")
	return
}

func (this_ *Invoker) GetDbServiceByName(name string) (res db.IService, err error) {
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
	res, err = db.New(&db.Config{
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
	this_.dbServiceCache[name] = res

	var scriptVar = "db"
	if name != "default" {
		scriptVar = "db_" + name
	}
	err = this_.setScriptVar(scriptVar, res)

	return
}

func (this_ *Invoker) GetZkService() (res zookeeper.IService, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get zk service error", zap.Any("error", err))
		}
	}()
	res, err = this_.GetZkServiceByName("")
	return
}

func (this_ *Invoker) GetZkServiceByName(name string) (res zookeeper.IService, err error) {
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
	res, err = zookeeper.New(&zookeeper.Config{
		Address:  config.Address,
		Username: config.Username,
		Password: config.Password,
	})
	if err != nil {
		util.Logger.Error("create zk service error", zap.Any("error", err))
		return
	}
	_, err = res.Get("/")
	if err != nil {
		util.Logger.Error("create zk service error", zap.Any("error", err))
		return
	}

	this_.zkServiceCache[name] = res
	var scriptVar = "zk"
	if name != "default" {
		scriptVar = "zk_" + name
	}
	err = this_.setScriptVar(scriptVar, res)
	return
}
