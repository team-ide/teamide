package maker

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/javascript"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"reflect"
	"teamide/pkg/maker/modelers"
)

func NewCompiler(app *Application) (compiler *Compiler, err error) {
	compiler = &Compiler{
		Application: app,
	}

	err = compiler.init()

	return
}

type Compiler struct {
	*Application

	script *Script
}

func (this_ *Compiler) setScriptVar(name string, value interface{}) (err error) {
	err = this_.script.Set(name, value)
	if err != nil {
		return
	}
	return
}

func (this_ *Compiler) initScript() (err error) {
	return
}

func (this_ *Compiler) initVar() (err error) {
	return
}

func (this_ *Compiler) init() (err error) {
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
			err = this_.setScriptVar(o.Name, o)
			if err != nil {
				util.Logger.Error("compiler data init set constant value error", zap.Any("name", o.Name), zap.Any("error", err))
				return
			}
			this_.constantContext[o.Name] = o
		}
	}

	err = this_.setScriptVar("error", this_.errorContext)
	if err != nil {
		return
	}
	for _, one := range this_.GetErrorList() {
		for _, o := range one.Options {
			err = this_.setScriptVar(o.Name, o)
			if err != nil {
				util.Logger.Error("compiler data init set error value error", zap.Any("name", o.Name), zap.Any("error", err))
				return
			}
			this_.errorContext[o.Name] = o
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
			component = NewRedisCompiler(one).ToContext()
			return
		})
		if err != nil {
			return
		}
	}
	for _, one := range this_.GetConfigDbList() {
		err = this_.BindComponent("db", one.Name, func() (component interface{}, err error) {
			component = NewDbCompiler(one).ToContext()
			return
		})
		if err != nil {
			return
		}
	}
	for _, one := range this_.GetConfigZkList() {
		err = this_.BindComponent("zk", one.Name, func() (component interface{}, err error) {
			component = NewZkCompiler(one).ToContext()
			return
		})
		if err != nil {
			return
		}
	}
	for _, one := range this_.GetConfigElasticsearchList() {
		err = this_.BindComponent("es", one.Name, func() (component interface{}, err error) {
			component = NewEsCompiler(one).ToContext()
			return
		})
		if err != nil {
			return
		}
	}
	for _, one := range this_.GetConfigKafkaList() {
		err = this_.BindComponent("kafka", one.Name, func() (component interface{}, err error) {
			component = NewKafkaCompiler(one).ToContext()
			return
		})
		if err != nil {
			return
		}
	}
	for _, one := range this_.GetConfigMongodbList() {
		err = this_.BindComponent("mongodb", one.Name, func() (component interface{}, err error) {
			component = NewMongodbCompiler(one).ToContext()
			return
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

func (this_ *Compiler) BindFunc(f *modelers.FuncModel) (err error) {
	this_.funcProgram[f.Name], err = this_.script.CompileScript(f.Func)
	if err != nil {
		util.Logger.Error("compiler bind func compile script error", zap.Any("name", f.Name), zap.Any("error", err))
		return
	}
	SetBySlash(this_.funcContext, f.Name, f)
	return
}

func (this_ *Compiler) BindComponent(componentType, name string, create func() (component interface{}, err error)) (err error) {
	if name == "" {
		name = "default"
	}
	logStr := "bind [" + componentType + "] component by name [" + name + "] "
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error(logStr+"error", zap.Any("error", err))
		}
	}()

	component, err := create()
	if err != nil {
		util.Logger.Error(logStr+"create error", zap.Any("error", err))
		return
	}

	var scriptVar = componentType
	if name != "default" {
		scriptVar = componentType + "_" + name
	}
	err = this_.setScriptVar(scriptVar, component)

	return
}

func (this_ *Compiler) BindDao(dao *modelers.DaoModel) (err error) {
	this_.daoProgram[dao.Name], err = this_.script.CompileScript(dao.Func)
	if err != nil {
		util.Logger.Error("compiler bind dao compile script error", zap.Any("name", dao.Name), zap.Any("error", err))
		return
	}
	SetBySlash(this_.daoContext, dao.Name, dao)
	return
}

func (this_ *Compiler) BindService(service *modelers.ServiceModel) (err error) {
	this_.serviceProgram[service.Name], err = this_.script.CompileScript(service.Func)
	if err != nil {
		util.Logger.Error("compiler bind service compile script error", zap.Any("name", service.Name), zap.Any("error", err))
		return
	}
	SetBySlash(this_.serviceContext, service.Name, service)
	return
}

func (this_ *Compiler) CompileFunc(f *modelers.FuncModel) (res *CompileResult, err error) {
	if f == nil {
		err = errors.New("compile func error, func is null")
		return
	}
	funcInvoke := invokeStart("compile func "+f.Name, nil)
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(funcInvoke.name + " error:" + fmt.Sprint(e))
			util.Logger.Error("compile func error", zap.Any("error", err))
		}
		funcInvoke.end(err)
		util.Logger.Debug(funcInvoke.name+" end", zap.Any("use", funcInvoke.use()))
	}()

	p := this_.funcProgram[f.Name]
	if p == nil {
		err = errors.New("compile func [" + f.Name + "] error, func program is null")
		return
	}

	util.Logger.Debug(funcInvoke.name + " start")

	script, err := this_.script.NewScriptByArgs(f.Args)
	if err != nil {
		return
	}

	res, err = this_.CompileProgram(funcInvoke.name, script, p)
	if err != nil {
		return
	}
	return
}

func (this_ *Compiler) CompileDaoByName(name string) (res *CompileResult, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("compile dao by name [" + name + "] error:" + fmt.Sprint(e))
			util.Logger.Error("compile dao by name error", zap.Any("error", err))
		}
	}()

	dao := this_.GetDao(name)
	if dao == nil {
		err = errors.New("dao [" + name + "] is not exist")
		util.Logger.Error("compile dao by name error", zap.Any("error", err))
		return
	}
	res, err = this_.CompileDao(dao)
	return
}

func (this_ *Compiler) CompileDao(dao *modelers.DaoModel) (res *CompileResult, err error) {
	if dao == nil {
		err = errors.New("compile dao error,dao is null")
		return
	}
	funcInvoke := invokeStart("compile dao "+dao.Name, nil)
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(funcInvoke.name + " error:" + fmt.Sprint(e))
			util.Logger.Error("compile dao error", zap.Any("error", err))
		}
		funcInvoke.end(err)
		util.Logger.Debug(funcInvoke.name+" end", zap.Any("use", funcInvoke.use()))
	}()

	p := this_.daoProgram[dao.Name]
	if p == nil {
		err = errors.New("compile dao [" + dao.Name + "] error, dao program is null")
		return
	}

	util.Logger.Debug(funcInvoke.name + " start")

	script, err := this_.script.NewScriptByArgs(dao.Args)
	if err != nil {
		return
	}

	res, err = this_.CompileProgram(funcInvoke.name, script, p)
	if err != nil {
		return
	}
	return
}

func (this_ *Compiler) CompileServiceByName(name string) (res *CompileResult, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("compile service by name [" + name + "] error:" + fmt.Sprint(e))
			util.Logger.Error("compile service by name error", zap.Any("error", err))
		}
	}()

	service := this_.GetService(name)
	if service == nil {
		err = errors.New("service [" + name + "] is not exist")
		util.Logger.Error("compile service by name error", zap.Any("error", err))
		return
	}
	res, err = this_.CompileService(service)
	return
}

func (this_ *Compiler) CompileService(service *modelers.ServiceModel) (res *CompileResult, err error) {
	if service == nil {
		err = errors.New("compile service error,service is null")
		return
	}
	funcInvoke := invokeStart("compile service "+service.Name, nil)
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(funcInvoke.name + " error:" + fmt.Sprint(e))
			util.Logger.Error("compile service error", zap.Any("error", err))
		}
		funcInvoke.end(err)
		util.Logger.Debug(funcInvoke.name+" end", zap.Any("use", funcInvoke.use()))
	}()

	p := this_.serviceProgram[service.Name]
	if p == nil {
		err = errors.New("compile service [" + service.Name + "] error, service program is null")
		return
	}
	util.Logger.Debug(funcInvoke.name + " start")

	script, err := this_.script.NewScriptByArgs(service.Args)
	if err != nil {
		return
	}

	res, err = this_.CompileProgram(funcInvoke.name, script, p)
	if err != nil {
		return
	}
	return
}

func (this_ *Compiler) CompileProgram(from string, script *Script, p *CompileProgram) (res *CompileResult, err error) {

	fmt.Println(util.GetStringValue(p.program.Program))
	codes := p.program.GetCode()
	for _, c := range codes {
		fmt.Println(reflect.TypeOf(c))
		fmt.Println(util.GetStringValue(c))
	}
	return
}
