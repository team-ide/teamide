package maker

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

func NewInvoker(app *Application) (runner *Invoker, err error) {
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
		errorContext:          make(map[string]*Error),
	}

	err = runner.init()

	return
}

type Invoker struct {
	app                   *Application
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
	errorContext          map[string]*Error
}

type Error struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func (this_ *Error) Error() string {
	return fmt.Sprintf("code:%s,msg:%s", this_.Code, this_.Msg)
}

func (this_ *Invoker) setScriptVar(name string, value interface{}) (err error) {
	err = this_.script.Set(name, value)
	if err != nil {
		util.Logger.Error("invoker set script var error", zap.Any("name", name), zap.Any("error", err))
		return
	}
	return
}

func (this_ *Invoker) init() (err error) {
	this_.script, err = javascript.NewScript()
	if err != nil {
		util.Logger.Error("invoker init new script error", zap.Any("error", err))
		return
	}
	err = this_.setScriptVar("error", this_.errorContext)
	if err != nil {
		return
	}
	for _, one := range this_.app.GetErrorList() {
		for _, o := range one.Options {
			this_.errorContext[o.Name] = &Error{
				Code: o.Code,
				Msg:  o.Msg,
			}
		}
	}

	// 初始化服务
	for _, one := range this_.app.GetConfigRedisList() {
		_, err = this_.GetRedisServiceByName(one.Name)
		if err != nil {
			util.Logger.Error("invoker init get redis service error", zap.Any("name", one.Name), zap.Any("error", err))
			return
		}
	}
	for _, one := range this_.app.GetConfigDbList() {
		_, err = this_.GetDbServiceByName(one.Name)
		if err != nil {
			util.Logger.Error("invoker init get db service error", zap.Any("name", one.Name), zap.Any("error", err))
			return
		}
	}
	for _, one := range this_.app.GetConfigZkList() {
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

	res, err = this_.InvokeFunc(funcInvoke.name, service.Func, invokeData)
	if err != nil {
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

	res, err = this_.InvokeFunc(funcInvoke.name, dao.Func, invokeData)
	if err != nil {
		return
	}
	return
}

func (this_ *Invoker) InvokeFunc(from string, code string, invokeData *InvokeData) (res interface{}, err error) {
	funcInvoke := invokeStart(from+" func code", invokeData)
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
	res, err = invokeData.InvokeScript(code)
	if err != nil {
		util.Logger.Error(funcInvoke.name+" invoke error", zap.Any("error", err))
		return
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
	this_.err = err
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
