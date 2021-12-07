package main

import (
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
func main() {
	web.StartServer()
}
