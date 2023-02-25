package invokers

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"sync"
	"teamide/pkg/db"
	"teamide/pkg/elasticsearch"
	"teamide/pkg/filework"
	"teamide/pkg/javascript"
	"teamide/pkg/kafka"
	"teamide/pkg/maker/modelers"
	"teamide/pkg/redis"
	"teamide/pkg/task"
	"teamide/pkg/util"
	"teamide/pkg/zookeeper"
)

func NewInvoker(app *modelers.Application) (runner *Invoker) {
	runner = &Invoker{
		app:                        app,
		redisServiceCache:          make(map[string]redis.Service),
		redisServiceCacheLock:      &sync.Mutex{},
		esServiceCache:             make(map[string]elasticsearch.Service),
		esServiceCacheLock:         &sync.Mutex{},
		zkServiceCache:             make(map[string]zookeeper.Service),
		zkServiceCacheLock:         &sync.Mutex{},
		dbServiceCache:             make(map[string]*db.Service),
		dbServiceCacheLock:         &sync.Mutex{},
		kafkaServiceCache:          make(map[string]kafka.Service),
		kafkaServiceCacheLock:      &sync.Mutex{},
		javascriptServiceCache:     make(map[string]javascript.Service),
		javascriptServiceCacheLock: &sync.Mutex{},
		fileServiceCache:           make(map[string]filework.Service),
		fileServiceCacheLock:       &sync.Mutex{},
		taskServiceCache:           make(map[string]task.Service),
		taskServiceCacheLock:       &sync.Mutex{},
	}
	return
}

type Invoker struct {
	app                        *modelers.Application
	redisServiceCache          map[string]redis.Service
	redisServiceCacheLock      sync.Locker
	esServiceCache             map[string]elasticsearch.Service
	esServiceCacheLock         sync.Locker
	zkServiceCache             map[string]zookeeper.Service
	zkServiceCacheLock         sync.Locker
	dbServiceCache             map[string]*db.Service
	dbServiceCacheLock         sync.Locker
	kafkaServiceCache          map[string]kafka.Service
	kafkaServiceCacheLock      sync.Locker
	javascriptServiceCache     map[string]javascript.Service
	javascriptServiceCacheLock sync.Locker
	fileServiceCache           map[string]filework.Service
	fileServiceCacheLock       sync.Locker
	taskServiceCache           map[string]task.Service
	taskServiceCacheLock       sync.Locker
}

func (this_ *Invoker) InvokeServiceByName(name string, invokeData *InvokeData) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("invoke service by name error", zap.Any("error", err))
		}
	}()

	service := this_.app.GetService(name)
	if service == nil {
		err = errors.New("service [" + name + "] is not exist")
		util.Logger.Error("invoke service by name error", zap.Any("error", err))
		return
	}
	res, err = this_.InvokeService(service, invokeData)
	return
}

func (this_ *Invoker) InvokeService(service *modelers.ServiceModel, invokeData *InvokeData) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("invoke service error", zap.Any("error", err))
		}
	}()
	if service == nil {
		err = errors.New("invoke service error,service is null")
		return
	}
	defer func() {
		util.Logger.Info("invoke service end", zap.Any("name", service.Name), zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))
	}()
	if invokeData == nil {
		invokeData = NewInvokeData(this_.app)
	}
	if invokeData.app == nil {
		invokeData.app = this_.app
	}
	util.Logger.Info("invoke service start", zap.Any("name", service.Name), zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))

	err = this_.InvokeSteps(service.Steps, invokeData)
	if err != nil {
		return
	}

	return
}

func (this_ *Invoker) InvokeDaoByName(name string, invokeData *InvokeData) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("invoke dao by name error", zap.Any("error", err))
		}
	}()

	dao := this_.app.GetDao(name)
	if dao == nil {
		err = errors.New("dao [" + name + "] is not exist")
		util.Logger.Error("invoke dao by name error", zap.Any("error", err))
		return
	}
	res, err = this_.InvokeDao(dao, invokeData)
	return
}

func (this_ *Invoker) InvokeDao(dao *modelers.DaoModel, invokeData *InvokeData) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("invoke dao error", zap.Any("error", err))
		}
	}()
	if dao == nil {
		err = errors.New("invoke dao error,service is null")
		return
	}
	defer func() {
		util.Logger.Info("invoke dao end", zap.Any("name", dao.Name), zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))
	}()
	if invokeData == nil {
		invokeData = NewInvokeData(this_.app)
	}
	if invokeData.app == nil {
		invokeData.app = this_.app
	}
	util.Logger.Info("invoke dao start", zap.Any("name", dao.Name), zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))

	err = this_.InvokeSteps(dao.Steps, invokeData)
	if err != nil {
		return
	}

	return
}

func (this_ *Invoker) InvokeSteps(steps []interface{}, invokeData *InvokeData) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("invoke steps error", zap.Any("error", err))
		}
	}()

	for _, step := range steps {
		err = this_.InvokeStep(step, invokeData)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *Invoker) InvokeScript(script string, invokeData *InvokeData) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("invoke script error", zap.Any("script", script), zap.Any("error", err))
		}
	}()

	if script == "" {
		return
	}
	return
}

func (this_ *Invoker) GetNameByRule(rule string, invokeData *InvokeData) (res string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get name by rule error", zap.Any("rule", rule), zap.Any("error", err))
		}
	}()

	if rule == "" {
		return
	}

	res = rule
	return
}
