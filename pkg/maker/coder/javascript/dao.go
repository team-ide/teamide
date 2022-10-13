package javascript

import (
	"go.uber.org/zap"
	"strings"
	"teamide/pkg/maker/model"
	"teamide/pkg/util"
)

func GetDaoMethodName(name string) (methodName string) {
	methodName = GetFormatMethodName(name)
	methodName = "dao" + util.Capitalize(methodName)
	return
}

func GetDaoJavascript(app *model.Application, dao *model.DaoModel) (javascript string, err error) {
	serviceMethodName := GetDaoMethodName(dao.Name)
	javascript += ""
	javascript += "function " + serviceMethodName + "("
	for _, arg := range dao.Args {
		javascript += arg.Name + ", "
	}
	javascript = strings.TrimSuffix(javascript, ", ")

	javascript += ") {"

	javascript += "\n"

	var stepsJavascript string
	stepsJavascript, err = GetJavascriptBySteps(app, dao.Steps, 1)
	if err != nil {
		util.Logger.Error("GetDaoJavascript GetJavascriptBySteps error", zap.Any("dao", dao), zap.Error(err))
		return
	}
	if util.IsNotEmpty(stepsJavascript) {
		javascript += stepsJavascript
	}

	if util.IsNotEmpty(dao.Return) {
		if dao.Return != "-" {
			util.AppendLine(&javascript, "return "+dao.Return, 1)
		} else {
			util.AppendLine(&javascript, "return", 1)
		}
	}

	javascript += "}"
	// fmt.Println(javascript)
	return
}
