package internal

import (
	"teamide/internal/context"
	"teamide/internal/module"
	"teamide/internal/web"
)

func Start(serverContext *context.ServerContext) (serverUrl string, err error) {
	var api *module.Api
	api, err = module.NewApi(serverContext)
	if err != nil {
		return
	}
	webServer := web.NewWebServer(serverContext, api)
	return webServer.Start()
}
