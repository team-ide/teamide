package maker

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/team-ide/go-tool/elasticsearch"
	"github.com/team-ide/go-tool/javascript"
	"github.com/team-ide/go-tool/kafka"
	"github.com/team-ide/go-tool/mongodb"
	"github.com/team-ide/go-tool/util"
	"github.com/team-ide/go-tool/zookeeper"
	"go.uber.org/zap"
	"strings"
	"sync"
	"teamide/pkg/maker/modelers"
	"time"
)

func NewInvoker(app *Application) (runner *Invoker, err error) {
	runner = &Invoker{
		app:                     app,
		redisServiceCache:       make(map[string]*ServiceRedis),
		redisServiceCacheLock:   &sync.Mutex{},
		esServiceCache:          make(map[string]elasticsearch.IService),
		esServiceCacheLock:      &sync.Mutex{},
		zkServiceCache:          make(map[string]zookeeper.IService),
		zkServiceCacheLock:      &sync.Mutex{},
		dbServiceCache:          make(map[string]*ServiceDb),
		dbServiceCacheLock:      &sync.Mutex{},
		mongodbServiceCache:     make(map[string]mongodb.IService),
		mongodbServiceCacheLock: &sync.Mutex{},
		kafkaServiceCache:       make(map[string]kafka.IService),
		kafkaServiceCacheLock:   &sync.Mutex{},

		constantContext: make(map[string]interface{}),
		errorContext:    make(map[string]*Error),
		serviceContext:  make(map[string]interface{}),
		daoContext:      make(map[string]interface{}),

		daoProgram:     make(map[string]*goja.Program),
		serviceProgram: make(map[string]*goja.Program),
	}

	err = runner.init()

	return
}

type Invoker struct {
	app                     *Application
	redisServiceCache       map[string]*ServiceRedis
	redisServiceCacheLock   sync.Locker
	esServiceCache          map[string]elasticsearch.IService
	esServiceCacheLock      sync.Locker
	zkServiceCache          map[string]zookeeper.IService
	zkServiceCacheLock      sync.Locker
	dbServiceCache          map[string]*ServiceDb
	dbServiceCacheLock      sync.Locker
	mongodbServiceCache     map[string]mongodb.IService
	mongodbServiceCacheLock sync.Locker
	kafkaServiceCache       map[string]kafka.IService
	kafkaServiceCacheLock   sync.Locker

	constantContext map[string]interface{}
	errorContext    map[string]*Error
	serviceContext  map[string]interface{}
	daoContext      map[string]interface{}
	daoProgram      map[string]*goja.Program
	serviceProgram  map[string]*goja.Program

	script *Script
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
		return
	}
	return
}

func (this_ *Invoker) initScript() (err error) {

	this_.script, err = NewScript()
	scriptContext := javascript.NewContext()
	for key, value := range scriptContext {
		err = this_.setScriptVar(key, value)
		if err != nil {
			return
		}
	}

	return
}
func (this_ *Invoker) init() (err error) {
	err = this_.initScript()
	if err != nil {
		return
	}

	err = this_.setScriptVar("constant", this_.constantContext)
	if err != nil {
		return
	}
	// 将 常量 error func 填充 至 script 变量域中
	for _, one := range this_.app.GetConstantList() {
		for _, o := range one.Options {
			err = this_.setScriptVar(o.Name, o.Value)
			if err != nil {
				util.Logger.Error("invoke data init set constant value error", zap.Any("name", o.Name), zap.Any("value", o.Value), zap.Any("error", err))
				return
			}
			this_.constantContext[o.Name] = o.Value
		}
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

	for _, one := range this_.app.GetFuncList() {
		err = this_.setScriptVar(one.Name, func(args ...interface{}) {
			util.Logger.Debug("func "+one.Name+" run start", zap.Any("func", one))
		})
		if err != nil {
			util.Logger.Error("invoke data init set func value error", zap.Any("name", one.Name), zap.Any("func", one), zap.Any("error", err))
			return
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
	for _, one := range this_.app.GetConfigElasticsearchList() {
		_, err = this_.GetEsServiceByName(one.Name)
		if err != nil {
			util.Logger.Error("invoker init get es service error", zap.Any("name", one.Name), zap.Any("error", err))
			return
		}
	}
	for _, one := range this_.app.GetConfigKafkaList() {
		_, err = this_.GetKafkaServiceByName(one.Name)
		if err != nil {
			util.Logger.Error("invoker init get kafka service error", zap.Any("name", one.Name), zap.Any("error", err))
			return
		}
	}
	for _, one := range this_.app.GetConfigMongodbList() {
		_, err = this_.GetMongodbServiceByName(one.Name)
		if err != nil {
			util.Logger.Error("invoker init get mongodb service error", zap.Any("name", one.Name), zap.Any("error", err))
			return
		}
	}

	err = this_.setScriptVar("dao", this_.daoContext)
	if err != nil {
		return
	}
	for _, one := range this_.app.GetDaoList() {
		err = this_.BindDao(one)
		if err != nil {
			return
		}
	}

	err = this_.setScriptVar("service", this_.serviceContext)
	if err != nil {
		return
	}
	for _, one := range this_.app.GetServiceList() {
		err = this_.BindService(one)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *Invoker) BindDao(dao *modelers.DaoModel) (err error) {
	this_.daoProgram[dao.Name], err = this_.script.CompileScript(dao.Func)
	if err != nil {
		util.Logger.Error("invoker bind dao compile script error", zap.Any("name", dao.Name), zap.Any("error", err))
		return
	}
	var run = func(args ...interface{}) (res any, err error) {
		data, err := this_.NewInvokeDataByArgs(dao.Args, args)
		if err != nil {
			return
		}
		res, err = this_.InvokeDao(dao, data)
		return
	}
	SetBySlash(this_.daoContext, dao.Name, run)
	return
}

func (this_ *Invoker) BindService(service *modelers.ServiceModel) (err error) {
	this_.serviceProgram[service.Name], err = this_.script.CompileScript(service.Func)
	if err != nil {
		util.Logger.Error("invoker bind service compile script error", zap.Any("name", service.Name), zap.Any("error", err))
		return
	}
	var run = func(args ...interface{}) (res any, err error) {
		data, err := this_.NewInvokeDataByArgs(service.Args, args)
		if err != nil {
			return
		}
		res, err = this_.InvokeService(service, data)
		return
	}
	SetBySlash(this_.serviceContext, service.Name, run)
	return
}

func SetBySlash(data map[string]interface{}, name string, value any) {
	//fmt.Println("SetBySlash:", name)
	index := strings.Index(name, "/")
	if index < 0 {
		data[name] = value
		return
	}
	pName := name[:index]
	cName := name[index+1:]
	//fmt.Println("SetBySlash pName:", pName, "cName:", cName)
	parent := data[pName]
	if parent == nil {
		parent = map[string]interface{}{}
		data[pName] = parent
	}
	SetBySlash(parent.(map[string]interface{}), cName, value)
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
		invokeData, err = this_.NewInvokeData()
		if err != nil {
			return
		}
	}
	if invokeData.app == nil {
		invokeData.app = this_.app
	}

	p := this_.serviceProgram[service.Name]
	if p == nil {
		err = errors.New("invoke service [" + service.Name + "] error, service program is null")
		return
	}
	util.Logger.Debug(funcInvoke.name + " start")

	res, err = this_.InvokeFunc(funcInvoke.name, p, invokeData)
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
		err = errors.New("invoke dao error,dao is null")
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
		invokeData, err = this_.NewInvokeData()
		if err != nil {
			return
		}
	}
	if invokeData.app == nil {
		invokeData.app = this_.app
	}

	p := this_.daoProgram[dao.Name]
	if p == nil {
		err = errors.New("invoke dao [" + dao.Name + "] error, dao program is null")
		return
	}

	util.Logger.Debug(funcInvoke.name + " start")

	res, err = this_.InvokeFunc(funcInvoke.name, p, invokeData)
	if err != nil {
		return
	}
	return
}

func (this_ *Invoker) InvokeFunc(from string, p *goja.Program, invokeData *InvokeData) (res interface{}, err error) {
	funcInvoke := invokeStart(from+" run func program", invokeData)
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
	v, err := invokeData.script.vm.RunProgram(p)
	if err != nil {
		util.Logger.Error(funcInvoke.name+" invoke error", zap.Any("error", err))
		return
	}
	res = v.Export()

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
