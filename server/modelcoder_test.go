package main

import (
	"fmt"
	"server/modelcoder"
	"testing"
)

func TestApplication(t *testing.T) {
	initApplication()

	variable := application.NewInvokeVariable(nil)

	res, err := application.InvokeServiceByName("user/queryPage", variable)
	if err != nil {
		panic(err)
	}
	fmt.Println(ToJSON(res))
}

var (
	application      *modelcoder.Application
	applicationModel *modelcoder.ApplicationModel
)

type Logger struct {
}

func (this_ *Logger) OutDebug() bool {
	return true
}
func (this_ *Logger) Debug(args ...interface{}) {
	fmt.Println(args...)
}
func (this_ *Logger) Info(args ...interface{}) {
	fmt.Println(args...)
}
func (this_ *Logger) Warn(args ...interface{}) {
	fmt.Println(args...)
}
func (this_ *Logger) Error(args ...interface{}) {
	fmt.Println(args...)
}
func initApplication() {
	initApplicationModel()
	application = modelcoder.NewApplication(applicationModel, &Logger{})
}
func initApplicationModel() {
	applicationModel = &modelcoder.ApplicationModel{}

	initDaoModel()
	initServiceModel()
}
func initDaoModel() {

	applicationModel.AppendDao(&modelcoder.DaoSqlSelectOne{
		Name: "user/queryPage",
	})
}
func initServiceModel() {

	applicationModel.AppendService(&modelcoder.ServiceFlow{
		Name: "user/queryPage",
	})
}
