package javascript

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"strings"
	"teamide/pkg/maker/base"
	"teamide/pkg/maker/coders/common"
	"teamide/pkg/maker/modelers"
)

type serviceCoder struct {
	*appCoder
}

func (this_ *serviceCoder) Gen(code *common.Code, model *modelers.ServiceModel) (err error) {

	return
}

func GetFormatMethodName(name string) (methodName string) {
	methodName = strings.ReplaceAll(name, ".", "_")
	chars := strings.Split(methodName, "")
	needToUp := false
	methodName = ""
	for _, char := range chars {
		if char == "/" {
			needToUp = true
		} else {
			if char == "." {
				char = "_"
			}
			if needToUp {
				char = strings.ToUpper(char)
				needToUp = false
			}
			methodName += char
		}
	}
	return
}

func GetServiceMethodName(name string) (methodName string) {
	methodName = GetFormatMethodName(name)
	methodName = "service" + util.FirstToUpper(methodName)
	return
}

func GetServiceJavascript(app *modelers.Application, service *modelers.ServiceModel) (javascript string, err error) {
	serviceMethodName := GetServiceMethodName(service.Name)
	javascript += ""
	javascript += "function " + serviceMethodName + "("
	for _, arg := range service.Args {
		javascript += arg.Name + ", "
	}
	javascript = strings.TrimSuffix(javascript, ", ")

	javascript += ") {"

	javascript += "\n"

	var stepsJavascript string
	stepsJavascript, err = GetJavascriptBySteps(app, service.Steps, 1)
	if err != nil {
		util.Logger.Error("GetServiceJavascript GetJavascriptBySteps error", zap.Any("service", service), zap.Error(err))
		return
	}
	if util.IsNotEmpty(stepsJavascript) {
		javascript += stepsJavascript
	}

	if util.IsNotEmpty(service.Return) {
		if service.Return != "-" {
			base.AppendLine(&javascript, "return "+service.Return, 1)
		} else {
			base.AppendLine(&javascript, "return", 1)
		}
	}

	javascript += "}"
	// fmt.Println(javascript)
	return
}
