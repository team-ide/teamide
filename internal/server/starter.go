package server

import (
	"teamide/internal/config"
	"teamide/internal/server/web"
)

func Start() (serverUrl string, err error) {
	webServer, err := web.NewWebServer(config.Config)
	if err != nil {
		return
	}
	return webServer.Start()
}
