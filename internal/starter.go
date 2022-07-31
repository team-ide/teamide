package internal

import (
	"teamide/internal/context"
	"teamide/internal/web"
)

func Start(serverContext *context.ServerContext) (serverUrl string, err error) {
	webServer := web.NewWebServer(serverContext)
	return webServer.Start()
}
