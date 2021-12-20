package main

import (
	"fmt"
	"server/modelcoder"
	"testing"
)

func TestApplication(t *testing.T) {
	initApplication()

	variable := application.NewInvokeVariable(nil)

	res, err := application.InvokeDaoByName("user/queryPage", variable)
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

	applicationModel.AppendDao(&modelcoder.DaoSqlSelectOne{
		Name: "user/queryPage",
	})
}
