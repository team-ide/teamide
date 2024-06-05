package maker

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/javascript"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
)

func NewCompiler(app *Application) (compiler *Compiler, err error) {
	compiler = &Compiler{
		Application: app,
		spaceCache:  make(map[string]*CompilerSpace),
	}

	err = compiler.init()

	return
}

type Compiler struct {
	*Application

	SpaceList  []*CompilerSpace
	spaceCache map[string]*CompilerSpace
	script     *Script
}

func (this_ *Compiler) setScriptVar(name string, value interface{}) (err error) {
	err = this_.script.Set(name, value)
	if err != nil {
		return
	}
	return
}

func (this_ *Compiler) initScript() (err error) {
	util.Logger.Debug("init script start")
	this_.script, err = this_.NewScript()
	scriptContext := javascript.NewContext()
	for key, value := range scriptContext {
		err = this_.setScriptVar(key, value)
		if err != nil {
			return
		}
	}
	util.Logger.Debug("init script end")
	return
}

func (this_ *Compiler) init() (err error) {
	util.Logger.Debug("init start")
	err = this_.initScript()
	if err != nil {
		return
	}

	util.Logger.Debug("init constant start")
	err = this_.setScriptVar("constant", this_.constantContext)
	space := this_.GetOrCreateSpace("constant")
	// 将 常量 error func 填充 至 script 变量域中
	for _, one := range this_.GetConstantList() {
		_, class := space.GetClass(one.Name, true)
		class.Constant = one
		for _, o := range one.Options {
			var valueType *ValueType
			valueType, err = this_.GetValueType(o.Type)
			if err != nil {
				util.Logger.Error("compiler init set constant value error", zap.Any("name", one.Name), zap.Any("error", err))
				return
			}
			field := class.addField(o.Name, valueType)
			field.ConstantOption = o

			err = this_.setScriptVar(o.Name, field)
			if err != nil {
				util.Logger.Error("compiler init set constant value error", zap.Any("name", o.Name), zap.Any("error", err))
				return
			}
			this_.constantContext[o.Name] = field
		}
	}
	util.Logger.Debug("init constant end")

	util.Logger.Debug("init error start")
	err = this_.setScriptVar("error", this_.errorContext)
	if err != nil {
		return
	}
	space = this_.GetOrCreateSpace("error")
	for _, one := range this_.GetErrorList() {
		_, class := space.GetClass(one.Name, true)
		class.Error = one
		for _, o := range one.Options {
			err = this_.setScriptVar(o.Name, o)
			if err != nil {
				util.Logger.Error("compiler init set error value error", zap.Any("name", o.Name), zap.Any("error", err))
				return
			}
			field := class.addField(o.Name, ValueTypeError)
			field.ErrorOption = o
			this_.errorContext[o.Name] = o
		}
	}

	util.Logger.Debug("init error end")

	util.Logger.Debug("init struct start")
	err = this_.setScriptVar("struct", this_.strictContext)
	if err != nil {
		return
	}
	space = this_.GetOrCreateSpace("struct")
	for _, one := range this_.GetStructList() {
		_, class := space.GetClass(one.Name, true)
		class.Struct = one
		var valueType *ValueType
		valueType, err = this_.GetValueType(one.Name)
		if err != nil {
			util.Logger.Error("compiler init set error strict error", zap.Any("name", one.Name), zap.Any("error", err))
			return
		}
		for _, f := range one.Fields {
			field := class.addField(f.Name, valueType.FieldTypes[f.Name])
			field.StructField = f
		}
		this_.strictContext[one.Name] = valueType
	}
	util.Logger.Debug("init struct end")

	util.Logger.Debug("init common start")

	// 初始化服务
	err = this_.BindComponent("common", "", func() (component interface{}, err error) {
		component = NewCommonCompiler(this_.GetApp()).ToContext()
		return
	})
	if err != nil {
		return
	}
	util.Logger.Debug("init common end")

	util.Logger.Debug("init redis start")

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
	util.Logger.Debug("init redis end")

	util.Logger.Debug("init db start")
	for _, one := range this_.GetConfigDbList() {
		err = this_.BindComponent("db", one.Name, func() (component interface{}, err error) {
			component = NewDbCompiler(one).ToContext()
			return
		})
		if err != nil {
			return
		}
	}
	util.Logger.Debug("init db end")

	util.Logger.Debug("init zk start")
	for _, one := range this_.GetConfigZkList() {
		err = this_.BindComponent("zk", one.Name, func() (component interface{}, err error) {
			component = NewZkCompiler(one).ToContext()
			return
		})
		if err != nil {
			return
		}
	}
	util.Logger.Debug("init zk end")

	util.Logger.Debug("init es start")
	for _, one := range this_.GetConfigEsList() {
		err = this_.BindComponent("es", one.Name, func() (component interface{}, err error) {
			component = NewEsCompiler(one).ToContext()
			return
		})
		if err != nil {
			return
		}
	}
	util.Logger.Debug("init es end")

	util.Logger.Debug("init kafka start")
	for _, one := range this_.GetConfigKafkaList() {
		err = this_.BindComponent("kafka", one.Name, func() (component interface{}, err error) {
			component = NewKafkaCompiler(one).ToContext()
			return
		})
		if err != nil {
			return
		}
	}
	util.Logger.Debug("init kafka end")

	util.Logger.Debug("init mongodb start")
	for _, one := range this_.GetConfigMongodbList() {
		err = this_.BindComponent("mongodb", one.Name, func() (component interface{}, err error) {
			component = NewMongodbCompiler(one).ToContext()
			return
		})
		if err != nil {
			return
		}
	}
	util.Logger.Debug("init mongodb end")

	util.Logger.Debug("init func start")
	err = this_.setScriptVar("func", this_.funcContext)
	if err != nil {
		return
	}
	space = this_.GetOrCreateSpace("func")
	var method *CompilerMethod
	for _, one := range this_.GetFuncList() {
		methodName, class := space.GetClass(one.Name, false)
		method, err = class.CreateMethod(methodName, one.Args)
		err = method.BindCode(one.Func)
		if err != nil {
			return
		}
		SetBySlash(this_.funcContext, one.Name, method)
	}
	util.Logger.Debug("init func end")

	util.Logger.Debug("init dao start")
	err = this_.setScriptVar("dao", this_.daoContext)
	if err != nil {
		return
	}
	space = this_.GetOrCreateSpace("dao")
	for _, one := range this_.GetDaoList() {
		methodName, class := space.GetClass(one.Name, false)
		var args = []*modelers.ArgModel{{Name: "ctx", Type: "context"}}
		args = append(args, one.Args...)
		method, err = class.CreateMethod(methodName, args)
		err = method.BindCode(one.Func)
		if err != nil {
			return
		}
		SetBySlash(this_.daoContext, one.Name, method)
	}
	util.Logger.Debug("init dao end")

	util.Logger.Debug("init service start")

	err = this_.setScriptVar("service", this_.serviceContext)
	if err != nil {
		return
	}
	space = this_.GetOrCreateSpace("service")
	for _, one := range this_.GetServiceList() {
		methodName, class := space.GetClass(one.Name, false)
		var args = []*modelers.ArgModel{{Name: "ctx", Type: "context"}}
		args = append(args, one.Args...)
		method, err = class.CreateMethod(methodName, args)
		err = method.BindCode(one.Func)
		if err != nil {
			return
		}
		SetBySlash(this_.serviceContext, one.Name, method)
	}
	util.Logger.Debug("init service end")

	util.Logger.Debug("init end")
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

func (this_ *Compiler) ToValueByValueType(originalValue any, valueType *ValueType) (targetValue any, err error) {
	return
}

func (this_ *Compiler) Compile(hasErrorContinue bool) (compileErrors []*CompileError) {
	util.Logger.Debug("compile start")
	for _, space := range this_.SpaceList {
		util.Logger.Debug("compile " + space.GetKey() + " start")
		for _, pack := range space.PackList {
			util.Logger.Debug("compile " + pack.GetKey() + " start")
			for _, class := range pack.ClassList {
				util.Logger.Debug("compile " + class.GetKey() + " start")
				for _, method := range class.MethodList {
					_, err := method.Compile()
					if err != nil {
						compileErrors = append(compileErrors, &CompileError{
							Err:    err,
							Method: method,
						})
						if !hasErrorContinue {
							return
						}
					}
				}
				util.Logger.Debug("compile " + class.GetKey() + " end")
			}
			util.Logger.Debug("compile " + pack.GetKey() + " end")
		}
		util.Logger.Debug("compile " + space.GetKey() + " end")
	}

	util.Logger.Debug("compile end")
	return
}

type CompileError struct {
	Err    error
	Method *CompilerMethod
}
