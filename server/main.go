package main

import (
	"base"
	"cache"
	"config"
	"db"
	"install"
	"redis"
	"service"
	"web"
	"worker"
	"zookeeper"
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
