package main

import (
	"os"
	"os/signal"
	"server/base"
	"server/cache"
	"server/config"
	"server/db"
	"server/install"
	"server/redis"
	"server/service"
	"server/web"
	"server/worker"
	"server/zookeeper"
	"syscall"
)

func init() {
	Init()
}
func Init() {
	base.Init()
	config.Init()
	worker.Init()
	db.Init()
	install.Init()
	redis.Init()
	service.Init()
	cache.Init()
	zookeeper.Init()
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
