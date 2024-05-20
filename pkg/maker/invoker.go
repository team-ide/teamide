package maker

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/javascript"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
)

func NewInvoker(compiler *Compiler) (runner *Invoker, err error) {
	runner = &Invoker{
		Compiler: compiler,
	}

	err = runner.init()

	return
}

type Invoker struct {
	*Compiler
}

type Error struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func (this_ *Error) Error() string {
	return fmt.Sprintf("code:%s,msg:%s", this_.Code, this_.Msg)
}

func (this_ *Invoker) init() (err error) {
	this_.script, err = this_.NewScript()
	scriptContext := javascript.NewContext()
	for key, value := range scriptContext {
		err = this_.setScriptVar(key, value)
		if err != nil {
			return
		}
	}

	err = this_.setScriptVar("constant", this_.constantContext)
	if err != nil {
		return
	}
	// 将 常量 error func 填充 至 script 变量域中
	for _, one := range this_.GetConstantList() {
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
	for _, one := range this_.GetErrorList() {
		for _, o := range one.Options {
			this_.errorContext[o.Name] = &Error{
				Code: o.Code,
				Msg:  o.Msg,
			}
		}
	}

	for _, one := range this_.GetFuncList() {
		err = this_.BindFunc(one)
		if err != nil {
			return
		}
	}

	// 初始化服务
	for _, one := range this_.GetConfigRedisList() {
		err = this_.BindComponent("redis", one.Name, func() (component interface{}, err error) {
			return NewComponentRedis(one)
		})
		if err != nil {
			return
		}
	}
	for _, one := range this_.GetConfigDbList() {
		err = this_.BindComponent("db", one.Name, func() (component interface{}, err error) {
			return NewComponentDb(one)
		})
		if err != nil {
			return
		}
	}
	for _, one := range this_.GetConfigZkList() {
		err = this_.BindComponent("zk", one.Name, func() (component interface{}, err error) {
			return NewComponentZk(one)
		})
		if err != nil {
			return
		}
	}
	for _, one := range this_.GetConfigElasticsearchList() {
		err = this_.BindComponent("es", one.Name, func() (component interface{}, err error) {
			return NewComponentEs(one)
		})
		if err != nil {
			return
		}
	}
	for _, one := range this_.GetConfigKafkaList() {
		err = this_.BindComponent("kafka", one.Name, func() (component interface{}, err error) {
			return NewComponentKafka(one)
		})
		if err != nil {
			return
		}
	}
	for _, one := range this_.GetConfigMongodbList() {
		err = this_.BindComponent("mongodb", one.Name, func() (component interface{}, err error) {
			return NewComponentMongodb(one)
		})
		if err != nil {
			return
		}
	}

	err = this_.setScriptVar("dao", this_.daoContext)
	if err != nil {
		return
	}
	for _, one := range this_.GetDaoList() {
		err = this_.BindDao(one)
		if err != nil {
			return
		}
	}

	err = this_.setScriptVar("service", this_.serviceContext)
	if err != nil {
		return
	}
	for _, one := range this_.GetServiceList() {
		err = this_.BindService(one)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *Invoker) BindFunc(f *modelers.FuncModel) (err error) {
	this_.funcProgram[f.Name], err = this_.script.CompileScript(f.Func)
	if err != nil {
		util.Logger.Error("invoker bind func compile script error", zap.Any("name", f.Name), zap.Any("error", err))
		return
	}
	var run = func(args ...interface{}) (res any, err error) {
		data, err := this_.NewInvokeDataByArgs(f.Args, args)
		if err != nil {
			return
		}
		res, err = this_.InvokeFunc(f, data)
		return
	}
	SetBySlash(this_.funcContext, f.Name, run)
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

func (this_ *Invoker) InvokeServiceByName(name string, invokeData *InvokeData) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("invoke service by name [" + name + "] error:" + fmt.Sprint(e))
			util.Logger.Error("invoke service by name error", zap.Any("error", err))
		}
	}()

	service := this_.GetService(name)
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

	p := this_.serviceProgram[service.Name]
	if p == nil {
		err = errors.New("invoke service [" + service.Name + "] error, service program is null")
		return
	}
	util.Logger.Debug(funcInvoke.name + " start")

	res, err = this_.InvokeProgram(funcInvoke.name, p, invokeData)
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

	dao := this_.GetDao(name)
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

	p := this_.daoProgram[dao.Name]
	if p == nil {
		err = errors.New("invoke dao [" + dao.Name + "] error, dao program is null")
		return
	}

	util.Logger.Debug(funcInvoke.name + " start")

	res, err = this_.InvokeProgram(funcInvoke.name, p, invokeData)
	if err != nil {
		return
	}
	return
}

func (this_ *Invoker) InvokeFunc(f *modelers.FuncModel, invokeData *InvokeData) (res interface{}, err error) {
	if f == nil {
		err = errors.New("invoke func error, func is null")
		return
	}
	funcInvoke := invokeStart("func "+f.Name, invokeData)
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(funcInvoke.name + " error:" + fmt.Sprint(e))
			util.Logger.Error("invoke func error", zap.Any("error", err))
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

	p := this_.funcProgram[f.Name]
	if p == nil {
		err = errors.New("invoke func [" + f.Name + "] error, func program is null")
		return
	}

	util.Logger.Debug(funcInvoke.name + " start")

	res, err = this_.InvokeProgram(funcInvoke.name, p, invokeData)
	if err != nil {
		return
	}
	return
}

func (this_ *Invoker) InvokeProgram(from string, p *CompileProgram, invokeData *InvokeData) (res interface{}, err error) {

	//var res interface{}
	v, err := invokeData.script.vm.RunProgram(p.program)
	if err != nil {
		util.Logger.Error(from+" invoke error", zap.Any("error", err))
		return
	}
	res = v.Export()

	return
}
