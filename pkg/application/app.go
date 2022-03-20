package application

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"teamide/pkg/application/base"
	common2 "teamide/pkg/application/common"
	invoke2 "teamide/pkg/application/invoke"
	model2 "teamide/pkg/application/model"

	"github.com/gin-gonic/gin"
)

func NewApplication(context *model2.ModelContext, options ...interface{}) common2.IApplication {
	if context == nil {
		return nil
	}
	context.Init()
	res := &Application{
		context: context,
	}
	res.initLoggerOption(&common2.LoggerDefault{OutDebug_: true, OutInfo_: true})
	res.initScriptOption(&common2.ScriptDefault{})

	if len(options) > 0 {
		for _, option := range options {
			res.initOption(option)
		}
	}

	return res
}

type Application struct {
	context                 *model2.ModelContext
	script                  common2.IScript
	javascriptExecutor      common2.IJavascriptExecutor
	scriptCache             map[string]reflect.Method
	logger                  common2.ILogger
	sqlExecutorCacheMutex   sync.Mutex
	sqlExecutorCache        map[string]common2.ISqlExecutor
	redisExecutorCacheMutex sync.Mutex
	redisExecutorCache      map[string]common2.IRedisExecutor
}

func (this_ *Application) initOption(option interface{}) *Application {
	if option == nil {
		return this_
	}
	logger, loggerOk := option.(common2.ILogger)
	if loggerOk {
		this_.initLoggerOption(logger)
	}
	script, scriptOk := option.(common2.IScript)
	if scriptOk {
		this_.initScriptOption(script)
	}
	return this_
}

func (this_ *Application) initLoggerOption(option common2.ILogger) *Application {
	if option == nil {
		return this_
	}
	this_.logger = option
	return this_
}

func (this_ *Application) initScriptOption(option common2.IScript) *Application {
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

func (this_ *Application) GetContext() *model2.ModelContext {
	return this_.context
}

func (this_ *Application) GetLogger() common2.ILogger {
	return this_.logger
}

func (this_ *Application) GetScript() common2.IScript {
	return this_.script
}
func (this_ *Application) GetJavascriptExecutor() common2.IJavascriptExecutor {
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

func (this_ *Application) GetSqlExecutor(name string) (executor common2.ISqlExecutor, err error) {
	datasource := this_.context.GetDatasourceDatabase(name)
	if datasource == nil {
		err = errors.New(fmt.Sprint("database datasource [", name, "] is not defind"))
		return
	}
	key := fmt.Sprint("database:", datasource.Host, ":", datasource.Port, ":", datasource.Username, ":", datasource.Password)
	this_.sqlExecutorCacheMutex.Lock()
	defer this_.sqlExecutorCacheMutex.Unlock()
	if this_.sqlExecutorCache == nil {
		this_.sqlExecutorCache = make(map[string]common2.ISqlExecutor)
	}
	executor = this_.sqlExecutorCache[key]

	if executor == nil {
		executor, err = common2.CreateSqlExecutor(datasource)
		if err != nil {
			return
		}
		this_.sqlExecutorCache[key] = executor
	}
	return
}

func (this_ *Application) GetRedisExecutor(name string) (executor common2.IRedisExecutor, err error) {
	datasource := this_.context.GetDatasourceRedis(name)
	if datasource == nil {
		err = errors.New(fmt.Sprint("redis datasource [", name, "] is not defind"))
		return
	}
	key := fmt.Sprint("redis:", datasource.Address, ":", datasource.Auth)
	this_.redisExecutorCacheMutex.Lock()
	defer this_.redisExecutorCacheMutex.Unlock()
	if this_.redisExecutorCache == nil {
		this_.redisExecutorCache = make(map[string]common2.IRedisExecutor)
	}
	executor = this_.redisExecutorCache[key]

	if executor == nil {
		executor, err = common2.CreateRedisExecutor(datasource)
		if err != nil {
			return
		}
		this_.redisExecutorCache[key] = executor
	}
	return
}

func (this_ *Application) GetKafkaExecutor(name string) (executor common2.IKafkaExecutor, err error) {
	datasource := this_.context.GetDatasourceKafka(name)
	if datasource == nil {
		err = errors.New(fmt.Sprint("kafka datasource [", name, "] is not defind"))
		return
	}
	return
}

func (this_ *Application) GetZookeeperExecutor(name string) (executor common2.IZookeeperExecutor, err error) {
	datasource := this_.context.GetDatasourceZookeeper(name)
	if datasource == nil {
		err = errors.New(fmt.Sprint("zookeeper datasource [", name, "] is not defind"))
		return
	}
	return
}

func (this_ *Application) InvokeActionByName(name string, invokeNamespace *common2.InvokeNamespace) (res interface{}, err error) {
	action := this_.context.GetAction(name)
	if action == nil {
		err = base.NewErrorActionIsNull("invoke action model [", name, "] is null")
		return
	}
	res, err = this_.InvokeAction(action, invokeNamespace)
	return
}
func (this_ *Application) InvokeAction(action *model2.ActionModel, invokeNamespace *common2.InvokeNamespace) (res interface{}, err error) {
	if invokeNamespace == nil {
		err = base.NewErrorVariableIsNull("invoke action ", action.Name, " invokeNamespace is null")
		return
	}
	res, err = invoke2.InvokeAction(this_, invokeNamespace, action)
	if err != nil {
		return
	}
	return
}

func (this_ *Application) InvokeTestByName(name string) (res *common2.TestResult, err error) {
	test := this_.context.GetTest(name)
	if test == nil {
		err = base.NewErrorActionIsNull("test model [", name, "] is null")
		return
	}
	res, err = this_.InvokeTest(test)
	return
}

func (this_ *Application) InvokeTest(test *model2.TestModel) (res *common2.TestResult, err error) {
	res, err = invoke2.InvokeTest(this_, test)
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

func (this_ *Application) StartServerWeb(server *model2.ServerWebModel) (err error) {

	err = invoke2.StartServerWeb(this_, server)

	if err != nil {
		return
	}

	return
}

func (this_ *Application) BindServerWebApis(serverWebToken *model2.ServerWebToken, gouterGroup *gin.RouterGroup) (err error) {

	err = invoke2.ServerWebBindApis(this_, serverWebToken, gouterGroup)
	if err != nil {
		return
	}
	return
}
