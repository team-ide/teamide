package application

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"teamide/pkg/application/base"
	"teamide/pkg/application/common"
	"teamide/pkg/application/invoke"
	"teamide/pkg/application/model"

	"github.com/gin-gonic/gin"
)

func NewApplication(context *model.ModelContext, options ...interface{}) common.IApplication {
	if context == nil {
		return nil
	}
	context.Init()
	res := &Application{
		context: context,
	}
	res.initLoggerOption(&common.LoggerDefault{OutDebug_: true, OutInfo_: true})
	res.initScriptOption(&common.ScriptDefault{})

	if len(options) > 0 {
		for _, option := range options {
			res.initOption(option)
		}
	}

	return res
}

type Application struct {
	context                 *model.ModelContext
	script                  common.IScript
	javascriptExecutor      common.IJavascriptExecutor
	scriptCache             map[string]reflect.Method
	logger                  common.ILogger
	sqlExecutorCacheMutex   sync.Mutex
	sqlExecutorCache        map[string]common.ISqlExecutor
	redisExecutorCacheMutex sync.Mutex
	redisExecutorCache      map[string]common.IRedisExecutor
}

func (this_ *Application) initOption(option interface{}) *Application {
	if option == nil {
		return this_
	}
	logger, loggerOk := option.(common.ILogger)
	if loggerOk {
		this_.initLoggerOption(logger)
	}
	script, scriptOk := option.(common.IScript)
	if scriptOk {
		this_.initScriptOption(script)
	}
	return this_
}

func (this_ *Application) initLoggerOption(option common.ILogger) *Application {
	if option == nil {
		return this_
	}
	this_.logger = option
	return this_
}

func (this_ *Application) initScriptOption(option common.IScript) *Application {
	if option == nil {
		return this_
	}

	this_.scriptCache = make(map[string]reflect.Method)
	this_.script = option

	reflectType := reflect.TypeOf(this_.script)
	count := reflectType.NumMethod()
	var i = 0
	for i = 0; i < count; i++ {
		method := reflectType.Method(i)
		this_.scriptCache[method.Name] = method
		this_.scriptCache[strings.ToLower(method.Name[0:1])+method.Name[1:]] = method
	}
	return this_
}

func (this_ *Application) GetContext() *model.ModelContext {
	return this_.context
}

func (this_ *Application) GetLogger() common.ILogger {
	return this_.logger
}

func (this_ *Application) GetScript() common.IScript {
	return this_.script
}
func (this_ *Application) GetJavascriptExecutor() common.IJavascriptExecutor {
	return this_.javascriptExecutor
}

func (this_ *Application) ScriptExist(name string) bool {
	_, ok := this_.scriptCache[name]
	return ok
}

func (this_ *Application) GetScriptMethod(name string) reflect.Method {
	method := this_.scriptCache[name]
	return method
}

func (this_ *Application) GetSqlExecutor(name string) (executor common.ISqlExecutor, err error) {
	datasource := this_.context.GetDatasourceDatabase(name)
	if datasource == nil {
		err = errors.New(fmt.Sprint("database datasource [", name, "] is not defind"))
		return
	}
	key := fmt.Sprint("database:", datasource.Host, ":", datasource.Port, ":", datasource.Username, ":", datasource.Password)
	this_.sqlExecutorCacheMutex.Lock()
	defer this_.sqlExecutorCacheMutex.Unlock()
	if this_.sqlExecutorCache == nil {
		this_.sqlExecutorCache = make(map[string]common.ISqlExecutor)
	}
	executor = this_.sqlExecutorCache[key]

	if executor == nil {
		executor, err = common.CreateSqlExecutor(datasource)
		if err != nil {
			return
		}
		this_.sqlExecutorCache[key] = executor
	}
	return
}

func (this_ *Application) GetRedisExecutor(name string) (executor common.IRedisExecutor, err error) {
	datasource := this_.context.GetDatasourceRedis(name)
	if datasource == nil {
		err = errors.New(fmt.Sprint("redis datasource [", name, "] is not defind"))
		return
	}
	key := fmt.Sprint("redis:", datasource.Address, ":", datasource.Auth)
	this_.redisExecutorCacheMutex.Lock()
	defer this_.redisExecutorCacheMutex.Unlock()
	if this_.redisExecutorCache == nil {
		this_.redisExecutorCache = make(map[string]common.IRedisExecutor)
	}
	executor = this_.redisExecutorCache[key]

	if executor == nil {
		executor, err = common.CreateRedisExecutor(datasource)
		if err != nil {
			return
		}
		this_.redisExecutorCache[key] = executor
	}
	return
}

func (this_ *Application) GetKafkaExecutor(name string) (executor common.IKafkaExecutor, err error) {
	datasource := this_.context.GetDatasourceKafka(name)
	if datasource == nil {
		err = errors.New(fmt.Sprint("kafka datasource [", name, "] is not defind"))
		return
	}
	return
}

func (this_ *Application) GetZookeeperExecutor(name string) (executor common.IZookeeperExecutor, err error) {
	datasource := this_.context.GetDatasourceZookeeper(name)
	if datasource == nil {
		err = errors.New(fmt.Sprint("zookeeper datasource [", name, "] is not defind"))
		return
	}
	return
}

func (this_ *Application) InvokeActionByName(name string, invokeNamespace *common.InvokeNamespace) (res interface{}, err error) {
	action := this_.context.GetAction(name)
	if action == nil {
		err = base.NewErrorActionIsNull("invoke action model [", name, "] is null")
		return
	}
	res, err = this_.InvokeAction(action, invokeNamespace)
	return
}
func (this_ *Application) InvokeAction(action *model.ActionModel, invokeNamespace *common.InvokeNamespace) (res interface{}, err error) {
	if invokeNamespace == nil {
		err = base.NewErrorVariableIsNull("invoke action ", action.Name, " invokeNamespace is null")
		return
	}
	res, err = invoke.InvokeAction(this_, invokeNamespace, action)
	if err != nil {
		return
	}
	return
}

func (this_ *Application) InvokeTestByName(name string) (res *common.TestResult, err error) {
	test := this_.context.GetTest(name)
	if test == nil {
		err = base.NewErrorActionIsNull("test model [", name, "] is null")
		return
	}
	res, err = this_.InvokeTest(test)
	return
}

func (this_ *Application) InvokeTest(test *model.TestModel) (res *common.TestResult, err error) {
	res, err = invoke.InvokeTest(this_, test)
	if err != nil {
		return
	}
	return
}

func (this_ *Application) StartServers() (err error) {
	if len(this_.GetContext().ServerWebs) > 0 {
		for _, one := range this_.GetContext().ServerWebs {
			err = this_.StartServerWeb(one)
			if err != nil {
				return
			}
		}
	}

	return
}

func (this_ *Application) StartServerWeb(server *model.ServerWebModel) (err error) {

	err = invoke.StartServerWeb(this_, server)

	if err != nil {
		return
	}

	return
}

func (this_ *Application) BindServerWebApis(serverWebToken *model.ServerWebToken, gouterGroup *gin.RouterGroup) (err error) {

	err = invoke.ServerWebBindApis(this_, serverWebToken, gouterGroup)
	if err != nil {
		return
	}
	return
}
