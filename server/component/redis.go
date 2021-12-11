package component

import (
	"server/base"
	"server/config"
	"strings"
)

var (
	Redis RedisService
)

func init() {
	var service interface{}
	var err error
	address := config.Config.Redis.Address
	auth := config.Config.Redis.Auth
	cluster := strings.Contains(address, ";")
	base.Logger.Info(base.LogStr("Redis初始化:address:", address))
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
	Redis = service.(RedisService)

	_, err = Redis.GetString("_")
	if err != nil {
		panic(err)
	}
	base.Logger.Info(base.LogStr("Redis连接成功!"))
}

type RedisService interface {
	Keys(pattern string, size int) (count int, keys []string, err error)
	Exists(key string) (value bool, err error)
	GetString(key string) (value string, err error)
	GetInt64(key string) (value int64, err error)
	IncrBy(key string, num int64) (value int64, err error)
	Set(key string, value string) (err error)
	SetInt64(key string, value int64) (err error)
	Sadd(key string, value string) (err error)
	Srem(key string, value string) (err error)
	Lpush(key string, value string) (err error)
	Rpush(key string, value string) (err error)
	Lset(key string, index int64, value string) (err error)
	Lrem(key string, count int64, value string) (err error)
	Hset(key string, field string, value string) (err error)
	Hdel(key string, field string) (err error)
	Del(key string) (count int, err error)
	DelPattern(key string) (count int, err error)
	Lock(key string, expire int, timeout int64) (unlock func() (err error), err error)
}

type RedisLock interface {
	Unlock() (err error)
}