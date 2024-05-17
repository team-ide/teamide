package maker

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/elasticsearch"
	"github.com/team-ide/go-tool/kafka"
	"github.com/team-ide/go-tool/mongodb"
	"github.com/team-ide/go-tool/util"
	"github.com/team-ide/go-tool/zookeeper"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
)

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
	config = this_.app.GetConfigZk(name)
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

	this_.zkServiceCache[name] = res
	var scriptVar = "zk"
	if name != "default" {
		scriptVar = "zk_" + name
	}
	err = this_.setScriptVar(scriptVar, res)
	return
}

func (this_ *Invoker) GetEsService() (res elasticsearch.IService, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get es service error", zap.Any("error", err))
		}
	}()
	res, err = this_.GetEsServiceByName("")
	return
}

func (this_ *Invoker) GetEsServiceByName(name string) (res elasticsearch.IService, err error) {
	if name == "" {
		name = "default"
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get es service by name error", zap.Any("name", name), zap.Any("error", err))
		}
	}()
	this_.esServiceCacheLock.Lock()
	defer this_.esServiceCacheLock.Unlock()
	res, find := this_.esServiceCache[name]

	if find {
		return
	}

	util.Logger.Info("es service not found,now create service", zap.Any("name", name))

	var config *modelers.ConfigEsModel
	config = this_.app.GetConfigElasticsearch(name)
	if config == nil {
		err = errors.New("config es [" + name + "] is not exist")
		util.Logger.Error("create es service error", zap.Any("error", err))
		return
	}
	res, err = elasticsearch.New(&elasticsearch.Config{
		Url:      config.Url,
		Username: config.Username,
		Password: config.Password,
	})
	if err != nil {
		util.Logger.Error("create es service error", zap.Any("error", err))
		return
	}

	this_.esServiceCache[name] = res
	var scriptVar = "es"
	if name != "default" {
		scriptVar = "es_" + name
	}
	err = this_.setScriptVar(scriptVar, res)
	return
}

func (this_ *Invoker) GetKafkaService() (res kafka.IService, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get kafka service error", zap.Any("error", err))
		}
	}()
	res, err = this_.GetKafkaServiceByName("")
	return
}

func (this_ *Invoker) GetKafkaServiceByName(name string) (res kafka.IService, err error) {
	if name == "" {
		name = "default"
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get kafka service by name error", zap.Any("name", name), zap.Any("error", err))
		}
	}()
	this_.kafkaServiceCacheLock.Lock()
	defer this_.kafkaServiceCacheLock.Unlock()
	res, find := this_.kafkaServiceCache[name]

	if find {
		return
	}

	util.Logger.Info("kafka service not found,now create service", zap.Any("name", name))

	var config *modelers.ConfigKafkaModel
	config = this_.app.GetConfigKafka(name)
	if config == nil {
		err = errors.New("config kafka [" + name + "] is not exist")
		util.Logger.Error("create kafka service error", zap.Any("error", err))
		return
	}
	res, err = kafka.New(&kafka.Config{
		Address:  config.Address,
		Username: config.Username,
		Password: config.Password,
	})
	if err != nil {
		util.Logger.Error("create kafka service error", zap.Any("error", err))
		return
	}

	this_.kafkaServiceCache[name] = res
	var scriptVar = "kafka"
	if name != "default" {
		scriptVar = "kafka_" + name
	}
	err = this_.setScriptVar(scriptVar, res)
	return
}

func (this_ *Invoker) GetMongodbService() (res mongodb.IService, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get mongodb service error", zap.Any("error", err))
		}
	}()
	res, err = this_.GetMongodbServiceByName("")
	return
}

func (this_ *Invoker) GetMongodbServiceByName(name string) (res mongodb.IService, err error) {
	if name == "" {
		name = "default"
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get mongodb service by name error", zap.Any("name", name), zap.Any("error", err))
		}
	}()
	this_.mongodbServiceCacheLock.Lock()
	defer this_.mongodbServiceCacheLock.Unlock()
	res, find := this_.mongodbServiceCache[name]

	if find {
		return
	}

	util.Logger.Info("mongodb service not found,now create service", zap.Any("name", name))

	var config *modelers.ConfigKafkaModel
	config = this_.app.GetConfigKafka(name)
	if config == nil {
		err = errors.New("config mongodb [" + name + "] is not exist")
		util.Logger.Error("create mongodb service error", zap.Any("error", err))
		return
	}
	res, err = mongodb.New(&mongodb.Config{
		Address:  config.Address,
		Username: config.Username,
		Password: config.Password,
	})
	if err != nil {
		util.Logger.Error("create mongodb service error", zap.Any("error", err))
		return
	}

	this_.mongodbServiceCache[name] = res
	var scriptVar = "mongodb"
	if name != "default" {
		scriptVar = "mongodb_" + name
	}
	err = this_.setScriptVar(scriptVar, res)
	return
}
