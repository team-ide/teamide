package common

import (
	"reflect"
	model2 "teamide/pkg/application/model"

	"github.com/gin-gonic/gin"
)

type IApplication interface {
	GetContext() *model2.ModelContext
	GetScript() IScript
	ScriptExist(name string) bool
	GetScriptMethod(name string) reflect.Method
	GetSqlExecutor(name string) (ISqlExecutor, error)
	GetRedisExecutor(name string) (IRedisExecutor, error)
	GetKafkaExecutor(name string) (IKafkaExecutor, error)
	GetZookeeperExecutor(name string) (IZookeeperExecutor, error)
	GetLogger() ILogger
	InvokeActionByName(name string, invokeNamespace *InvokeNamespace) (res interface{}, err error)
	InvokeAction(action *model2.ActionModel, invokeNamespace *InvokeNamespace) (res interface{}, err error)
	InvokeTestByName(name string) (res *TestResult, err error)
	InvokeTest(test *model2.TestModel) (res *TestResult, err error)
	StartServers() (err error)
	StartServerWeb(serverWeb *model2.ServerWebModel) (err error)
	BindServerWebApis(serverWebToken *model2.ServerWebToken, gouterGroup *gin.RouterGroup) (err error)
}
