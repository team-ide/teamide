package invokers

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/db"
	"github.com/team-ide/go-tool/elasticsearch"
	"github.com/team-ide/go-tool/javascript"
	"github.com/team-ide/go-tool/kafka"
	"github.com/team-ide/go-tool/redis"
	"github.com/team-ide/go-tool/util"
	"github.com/team-ide/go-tool/zookeeper"
	"go.uber.org/zap"
	"sync"
	"teamide/pkg/maker/modelers"
	"time"
)

func NewInvoker(app *modelers.Application) (runner *Invoker, err error) {
	runner = &Invoker{
		app:                   app,
		redisServiceCache:     make(map[string]redis.IService),
		redisServiceCacheLock: &sync.Mutex{},
		esServiceCache:        make(map[string]elasticsearch.IService),
		esServiceCacheLock:    &sync.Mutex{},
		zkServiceCache:        make(map[string]zookeeper.IService),
		zkServiceCacheLock:    &sync.Mutex{},
		dbServiceCache:        make(map[string]db.IService),
		dbServiceCacheLock:    &sync.Mutex{},
		kafkaServiceCache:     make(map[string]kafka.IService),
		kafkaServiceCacheLock: &sync.Mutex{},
	}

	err = runner.init()

	return
}

type Invoker struct {
	app                   *modelers.Application
	redisServiceCache     map[string]redis.IService
	redisServiceCacheLock sync.Locker
	esServiceCache        map[string]elasticsearch.IService
	esServiceCacheLock    sync.Locker
	zkServiceCache        map[string]zookeeper.IService
	zkServiceCacheLock    sync.Locker
	dbServiceCache        map[string]db.IService
	dbServiceCacheLock    sync.Locker
	kafkaServiceCache     map[string]kafka.IService
	kafkaServiceCacheLock sync.Locker
	script                *javascript.Script
}

func (this_ *Invoker) init() (err error) {
	// 初始化服务
	for _, one := range this_.app.ConfigRedisList {
		_, err = this_.GetRedisServiceByName(one.Name)
		if err != nil {
			util.Logger.Error("invoker init get redis service error", zap.Any("name", one.Name), zap.Any("error", err))
			return
		}
	}
	for _, one := range this_.app.ConfigDbList {
		_, err = this_.GetDbServiceByName(one.Name)
		if err != nil {
			util.Logger.Error("invoker init get db service error", zap.Any("name", one.Name), zap.Any("error", err))
			return
		}
	}
	for _, one := range this_.app.ConfigZkList {
		_, err = this_.GetZkServiceByName(one.Name)
		if err != nil {
			util.Logger.Error("invoker init get zk service error", zap.Any("name", one.Name), zap.Any("error", err))
			return
		}
	}

	return
}

func (this_ *Invoker) InvokeServiceByName(name string, invokeData *InvokeData) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("invoke service by name [" + name + "] error:" + fmt.Sprint(e))
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
	if service == nil {
		err = errors.New("invoke service error,service is null")
		return
	}
	funcInvoke := invokeStart("service "+service.Name, invokeData)
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(funcInvoke.name + " error:" + fmt.Sprint(e))
			util.Logger.Error("invoke service error", zap.Any("error", err))
		}
		funcInvoke.end(err)
		util.Logger.Debug(funcInvoke.name+" end", zap.Any("use", funcInvoke.use()))
	}()
	if invokeData == nil {
		invokeData, err = NewInvokeData(this_.app)
		if err != nil {
			return
		}
	}
	if invokeData.app == nil {
		invokeData.app = this_.app
	}
	util.Logger.Debug(funcInvoke.name + " start")

	err = this_.InvokeSteps(funcInvoke.name, service.Steps, invokeData)
	if err != nil {
		return
	}

	if service.Return != "" {
		res, err = invokeData.InvokeScript(service.Return)
		if err != nil {
			util.Logger.Error(funcInvoke.name+" get return value error", zap.Any("return", service.Return), zap.Any("error", err))
			return
		}
		return
	}

	return
}

func (this_ *Invoker) InvokeDaoByName(name string, invokeData *InvokeData) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("invoke dao by name [" + name + "] error:" + fmt.Sprint(e))
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
	if dao == nil {
		err = errors.New("invoke dao error,service is null")
		return
	}
	funcInvoke := invokeStart("dao "+dao.Name, invokeData)
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(funcInvoke.name + " error:" + fmt.Sprint(e))
			util.Logger.Error("invoke dao error", zap.Any("error", err))
		}
		funcInvoke.end(err)
		util.Logger.Debug(funcInvoke.name+" end", zap.Any("use", funcInvoke.use()))
	}()
	if invokeData == nil {
		invokeData, err = NewInvokeData(this_.app)
		if err != nil {
			return
		}
	}
	if invokeData.app == nil {
		invokeData.app = this_.app
	}
	util.Logger.Debug(funcInvoke.name + " start")

	err = this_.InvokeSteps(funcInvoke.name, dao.Steps, invokeData)
	if err != nil {
		return
	}

	if dao.Return != "" {
		res, err = invokeData.InvokeScript(dao.Return)
		if err != nil {
			util.Logger.Error(funcInvoke.name+" get return value error", zap.Any("return", dao.Return), zap.Any("error", err))
			return
		}
		return
	}

	return
}

func (this_ *Invoker) InvokeSteps(from string, steps []interface{}, invokeData *InvokeData) (err error) {
	funcInvoke := invokeStart(from+" steps", invokeData)
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(funcInvoke.name + " error:" + fmt.Sprint(e))
			util.Logger.Error(funcInvoke.name+" error", zap.Any("error", err))
		}
		funcInvoke.end(err)
		util.Logger.Debug(funcInvoke.name+" end", zap.Any("use", funcInvoke.use()))
	}()
	util.Logger.Debug(funcInvoke.name + " start")
	//var res interface{}
	var isReturn bool
	for _, step := range steps {
		_, isReturn, err = this_.InvokeStep(funcInvoke.name, step, invokeData)
		if err != nil {
			return
		}
		if isReturn {
			break
		}
	}

	return
}

func invokeStart(name string, invokeData *InvokeData) (funcInvoke *FuncInvoke) {
	funcInvoke = &FuncInvoke{
		name:       name,
		invokeData: invokeData,
	}
	funcInvoke.start()
	return
}

type FuncInvoke struct {
	name       string
	startTime  time.Time
	endTime    time.Time
	err        error
	invokeData *InvokeData
}

func (this_ *FuncInvoke) start() {
	this_.startTime = time.Now()
}

func (this_ *FuncInvoke) end(err error) {
	this_.endTime = time.Now()
}

func (this_ *FuncInvoke) use() string {
	return GetDurationFormatByMillisecond(this_.endTime.UnixMilli() - this_.startTime.UnixMilli())
}

func GetDurationFormatByMillisecond(millisecond int64) (formatString string) {
	if millisecond == 0 {
		return fmt.Sprintf("%d毫秒", 0)
	}

	duration := time.Duration(millisecond) * time.Millisecond
	h := int(duration.Hours())
	m := int(duration.Minutes()) % 60
	s := int(duration.Seconds()) % 60
	ms := int(duration.Milliseconds()) % 1000
	if h > 0 {
		formatString = fmt.Sprintf("%d小时", h)
	}
	if m > 0 {
		formatString += fmt.Sprintf("%d分钟", m)
	}
	if s > 0 {
		formatString += fmt.Sprintf("%d秒", s)
	}
	if ms > 0 {
		formatString += fmt.Sprintf("%d毫秒", ms)
	}
	return
}
