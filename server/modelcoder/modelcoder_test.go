package modelcoder

import (
	"fmt"
	"testing"
)

func TestScriptParser(t *testing.T) {
	initApplication()
	script := "'张三'+(user.age>=18 || user.name=='张三'?'成年':'未成年')"

	variable := application.NewInvokeVariable(&VariableData{
		Name: "user",
		Data: map[string]interface{}{
			"name": "张三",
			"age":  16,
		},
		DataStruct: application.context.GetStruct("user"),
	})

	value, err := getScriptValue(application, variable, script)

	if err != nil {
		panic(err)
	}

	println("script:", script)
	println("value:", value)
	println("value:", fmt.Sprint(value))
	println("value:", value == 12)
}

func TestApplication(t *testing.T) {
	initApplication()

	variable := application.NewInvokeVariable(&VariableData{
		Name: "user",
		Data: map[string]interface{}{
			"name": "张三",
			"age":  16,
		},
		DataStruct: application.context.GetStruct("user"),
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
			{Name: "age"},
		},
	})
	applicationModel.AppendStruct(&StructModel{
		Name: "user",
		Fields: []*StructFieldModel{
			{Name: "name"},
			{Name: "age"},
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
		Parameters: []*VariableModel{
			{Name: "user", DataType: "user"},
		},
		Result: &VariableModel{
			Name: "user", DataType: "user",
		},
	})
}
func initServiceModel() {

	applicationModel.AppendService(&ServiceFlowModel{
		Name: "user/insert",
		Parameters: []*VariableModel{
			{Name: "user", DataType: "user"},
		},
		Steps: []ServiceFlowStepModel{
			&ServiceFlowStepStartModel{Name: "start", Next: "insert"},
			&ServiceFlowStepDaoModel{Name: "insert", DaoName: "user/insert"},
		},
	})
}
