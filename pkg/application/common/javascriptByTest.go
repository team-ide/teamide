package common

import (
	"teamide/pkg/application/base"
	"teamide/pkg/application/model"
)

func GetTestJavascriptByTestStep(app IApplication, test *model.TestModel) (javascript string, err error) {
	methodName := GetJavascriptMethodName(test.Name)
	javascript += ""
	javascript += "function test_" + methodName + "() {\n"

	var stepsJavascript string
	stepsJavascript, err = GetJavascriptBySteps(app, test.Steps, 1)
	if err != nil {
		return
	}
	if base.IsNotEmpty(stepsJavascript) {
		javascript += stepsJavascript
	}

	javascript += "}"
	return
}
