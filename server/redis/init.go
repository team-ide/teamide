package redis

import (
	"fmt"
	"server/config"
	"strings"
)

var (
	RedisService Service
)

func init() {
	var service interface{}
	var err error
	address := config.Config.Redis.Address
	auth := config.Config.Redis.Auth
	cluster := strings.Contains(address, ";")
	fmt.Println("Redis初始化：address:", address)
	if cluster {
		service, err = CreateRedisClusterService(address, auth)
		if err != nil {
			panic(err)
		}
	} else {
		service, err = CreateRedisPoolService(address, auth)
		if err != nil {
			panic(err)
		}
	}
	RedisService = service.(Service)

	_, err = RedisService.Get("_")
	if err != nil {
		panic(err)
	}
	fmt.Println("Redis连接成功！")
}
func Init() {
}
