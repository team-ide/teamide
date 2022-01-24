package common

import (
	"fmt"
	"strings"
	"teamide/application/base"
	"teamide/application/model"
)

func GetWebApiJavascriptByAction(app IApplication, action *model.ActionModel, shouldValidataToken bool) (javascript string, err error) {
	methodName := GetJavascriptMethodName(action.Name)
	javascript += ""
	javascript += "function web_api_" + methodName + "() {\n"

	if shouldValidataToken {
		base.AppendLine(&javascript, "validataRequestToken()", 1)
	}

	var variablesJavascript string
	variablesJavascript, err = getWebApiJavascriptByVariables(app, action, action.InVariables, 1)
	if err != nil {
		return
	}
	if base.IsNotEmpty(variablesJavascript) {
		javascript += variablesJavascript
	}
	callArgs := `"` + action.Name + `", `
	for _, one := range action.InVariables {
		callArgs += one.Name + ", "
	}
	callArgs = strings.TrimSuffix(callArgs, ", ")
	base.AppendLine(&javascript, "// 调用服务方法", 1)

	base.AppendLine(&javascript, "return action("+callArgs+")", 1)

	javascript += "}"
	return
}

func getWebApiJavascriptByVariables(app IApplication, action *model.ActionModel, variables []*model.VariableModel, tab int) (javascript string, err error) {
	for _, one := range variables {
		if base.IsNotEmpty(one.Comment) {
			base.AppendLine(&javascript, "// "+one.Comment, tab)
		}
		base.AppendLine(&javascript, `addDataInfo("`+one.Name+`", "`+one.DataType+`", "`+one.Comment+`", "`+one.Value+`", `+fmt.Sprint(one.IsList)+`)`, tab)
		base.AppendLine(&javascript, one.Name+" = "+`getRequestData("`+one.DataPlace+`", "`+one.Name+`", "`+one.Value+`")`, tab)
	}

	return
}
