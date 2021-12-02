package main

import (
	"base"
	"cache"
	"config"
	"db"
	"redis"
	"service"
	"sync"
	"version"
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
	version.Init()
	redis.Init()
	service.Init()
	cache.Init()
	zookeeper.Init()
	web.Init()
}
func main() {
	service.TestTotalBatchInsert(50, 100000)
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
