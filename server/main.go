package main

import (
	"os"
	"os/signal"
	"server/base"
	"server/component"
	"server/config"
	installService "server/service/install"
	"server/web"
	"syscall"
)

func init() {
	Init()
}
func Init() {
	base.Init()
	config.Init()
	component.Init()
	installService.Init()
	web.Init()
}

var (
	done = make(chan os.Signal, 1)
)

func main() {
	web.StartServer()
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
