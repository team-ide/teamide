package invoke

import (
	"teamide/pkg/application/common"
	"teamide/pkg/application/model"
)

func StartServerWeb(app common.IApplication, serverWeb *model.ServerWebModel) (err error) {
	if app.GetLogger() != nil && app.GetLogger().OutDebug() {
		app.GetLogger().Debug("start server web [", serverWeb.Name, "] start")
	}

	err = startServerWeb(app, serverWeb)
	if err != nil {
		if app.GetLogger() != nil {
			app.GetLogger().Error("start server web [", serverWeb.Name, "] error:", err)
		}
		return
	}
	return
}
