package invoke

import (
	"teamide/application/base"
	"teamide/application/common"
	"teamide/application/model"
)

func invokeJavascript(app common.IApplication, invokeNamespace *common.InvokeNamespace, javascript string) (res interface{}, err error) {
	res, err = ExecuteFunctionScript(app, invokeNamespace, javascript)
	if err != nil {
		return
	}
	return
}

func InvokeService(app common.IApplication, invokeNamespace *common.InvokeNamespace, service *model.ServiceModel) (res interface{}, err error) {
	if app.GetLogger() != nil && app.GetLogger().OutDebug() {
		app.GetLogger().Debug("invoke service [", service.Name, "] start")
		// app.GetLogger().Debug("invoke service [", service.Name, "] invokeNamespace:", app.GetScript().DataToJSON(invokeNamespace))
	}

	startTime := base.GetNowTime()
	defer func() {
		endTime := base.GetNowTime()
		if app.GetLogger() != nil && app.GetLogger().OutDebug() {
			app.GetLogger().Debug("invoke service [", service.Name, "] end, use:", (endTime - startTime), "ms")
		}
	}()

	if base.IsEmpty(service.ServiceJavascript) {
		service.ServiceJavascript, err = common.GetServiceJavascriptByService(app, service)
		if err != nil {
			return
		}
	}
	res, err = invokeJavascript(app, invokeNamespace, service.ServiceJavascript)
	if err != nil {
		if app.GetLogger() != nil {
			app.GetLogger().Error("invoke service [", service.Name, "] error:", err)
		}
		return
	}
	return

}
