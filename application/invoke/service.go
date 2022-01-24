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

func InvokeAction(app common.IApplication, invokeNamespace *common.InvokeNamespace, action *model.ActionModel) (res interface{}, err error) {
	if app.GetLogger() != nil && app.GetLogger().OutDebug() {
		app.GetLogger().Debug("invoke action [", action.Name, "] start")
		// app.GetLogger().Debug("invoke action [", action.Name, "] invokeNamespace:", app.GetScript().DataToJSON(invokeNamespace))
	}

	startTime := base.GetNowTime()
	defer func() {
		endTime := base.GetNowTime()
		if app.GetLogger() != nil && app.GetLogger().OutDebug() {
			app.GetLogger().Debug("invoke action [", action.Name, "] end, use:", (endTime - startTime), "ms")
		}
	}()

	if base.IsEmpty(action.ActionJavascript) {
		action.ActionJavascript, err = common.GetActionJavascriptByAction(app, action)
		if err != nil {
			return
		}
	}
	res, err = invokeJavascript(app, invokeNamespace, action.ActionJavascript)
	if err != nil {
		if app.GetLogger() != nil {
			app.GetLogger().Error("invoke action [", action.Name, "] error:", err)
		}
		return
	}
	return

}
