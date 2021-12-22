package modelcoder

import (
	"fmt"
	"testing"
)

func TestScriptParser(t *testing.T) {
	initApplication()
	script := "user.userId + '1' + $factory.GetID() +'aaa'.length()"
	scriptParser := &scriptParser{
		script:             script,
		factory:            application.factory,
		factoryScriptCache: application.factoryScriptCache,
	}

	err := scriptParser.parse()

	if err != nil {
		panic(err)
	}
}

func TestApplication(t *testing.T) {
	initApplication()

	variable := application.NewInvokeVariable(&ParamData{
		Name: "user",
		Data: map[string]interface{}{
			"name": "张三",
			"age":  16,
		},
	})

	res, err := application.InvokeServiceByName("user/insert", variable)
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

	initStructModel()
	initDaoModel()
	initServiceModel()
}
func initStructModel() {

	applicationModel.AppendStruct(&StructModel{
		Name: "user",
		Fields: []*StructFieldModel{
			{Name: "name"},
			{Name: "arg"},
		},
	})
	applicationModel.AppendStruct(&StructModel{
		Name: "user",
		Fields: []*StructFieldModel{
			{Name: "name"},
			{Name: "arg"},
		},
	})
}
func initDaoModel() {

	applicationModel.AppendDao(&DaoSqlInsertModel{
		Name:  "user/insert",
		Table: "IM_USER",
		Columns: []*DaoSqlInsertColumn{
			{Name: "name", ValueScript: "user.name", Required: true},
			{Name: "age", ValueScript: "user.age", Required: true},
		},
		Params: []*ParamModel{
			{Name: "user", DataType: "user"},
		},
		Result: &ParamModel{
			Name: "user", DataType: "user",
		},
	})
}
func initServiceModel() {

	applicationModel.AppendService(&ServiceFlowModel{
		Name: "user/insert",
		Params: []*ParamModel{
			{Name: "user", DataType: "user"},
		},
		Steps: []ServiceFlowStepModel{
			&ServiceFlowStepStartModel{Name: "start", Next: "insert"},
			&ServiceFlowStepDaoModel{Name: "insert", DaoName: "user/insert"},
		},
	})
}
