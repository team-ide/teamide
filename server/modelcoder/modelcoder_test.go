package modelcoder

import (
	"fmt"
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
	application      *Application
	applicationModel *ApplicationModel
)

type Logger struct {
}

func (this_ *Logger) OutDebug() bool {
	return true
}
func (this_ *Logger) OutInfo() bool {
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
	application = NewApplication(applicationModel, &Logger{})
}
func initApplicationModel() {
	applicationModel = &ApplicationModel{}

	initDaoModel()
	initServiceModel()
}
func initDaoModel() {

	applicationModel.AppendDao(&DaoSqlSelectOneModel{
		Name: "user/queryPage",
	})
}
func initServiceModel() {

	applicationModel.AppendService(&ServiceFlowModel{
		Name: "user/queryPage",
		Steps: []ServiceFlowStepModel{
			&ServiceFlowStepStartModel{Name: "start", Next: "queryPage"},
			&ServiceFlowStepDaoModel{Name: "queryPage"},
		},
	})
}
