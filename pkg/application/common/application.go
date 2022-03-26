package common

import (
	"reflect"
	"teamide/pkg/application/model"

	"github.com/gin-gonic/gin"
)

type IApplication interface {
	GetContext() *model.ModelContext
	GetScript() IScript
	ScriptExist(name string) bool
	GetScriptMethod(name string) reflect.Method
	GetSqlExecutor(name string) (ISqlExecutor, error)
	GetRedisExecutor(name string) (IRedisExecutor, error)
	GetKafkaExecutor(name string) (IKafkaExecutor, error)
	GetZookeeperExecutor(name string) (IZookeeperExecutor, error)
	GetLogger() ILogger
	InvokeActionByName(name string, invokeNamespace *InvokeNamespace) (res interface{}, err error)
	InvokeAction(action *model.ActionModel, invokeNamespace *InvokeNamespace) (res interface{}, err error)
	InvokeTestByName(name string) (res *TestResult, err error)
	InvokeTest(test *model.TestModel) (res *TestResult, err error)
	StartServers() (err error)
	StartServerWeb(serverWeb *model.ServerWebModel) (err error)
	BindServerWebApis(serverWebToken *model.ServerWebToken, gouterGroup *gin.RouterGroup) (err error)
}
