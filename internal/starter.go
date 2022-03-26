package internal

import (
	"teamide/internal/context"
	"teamide/internal/web"
)

func Start(ServerContext *context.ServerContext) (serverUrl string, err error) {
	webServer := web.NewWebServer(ServerContext)
	return webServer.Start()
}
