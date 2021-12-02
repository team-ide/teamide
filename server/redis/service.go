package redis

import (
	"base"
	"fmt"
	"strings"
	"worker"
)

type Automatic struct {
	address           string
	auth              string
	automaticShutdown *worker.AutomaticShutdown
}

func (automatic *Automatic) CreateAutomaticShutdown(automaticShutdown *worker.AutomaticShutdown) error {
	service, err := CreateRedisPoolService(automatic.address, automatic.auth)
	if err != nil {
		return err
	}

	_, err = service.Get("test_key")
	if err != nil {
		service.pool.Close()
		return err
	}
	// 默认10分钟自动关闭
	automaticShutdown.AutomaticShutdown = 10 * 60
	automaticShutdown.Service = service
	automaticShutdown.Stop = func() {
		service.pool.Close()
	}
	automatic.automaticShutdown = automaticShutdown

	return err
}

func (automatic *Automatic) CreateClusterAutomaticShutdown(automaticShutdown *worker.AutomaticShutdown) error {
	service, err := CreateRedisClusterService(automatic.address, automatic.auth)
	if err != nil {
		return err
	}
	_, err = service.Get("test_key")
	if err != nil {
		service.redisCluster.Close()
		return err
	}
	// 默认10分钟自动关闭
	automaticShutdown.AutomaticShutdown = 10 * 60
	automaticShutdown.Service = service
	automaticShutdown.Stop = func() {
		service.redisCluster.Close()
	}
	automatic.automaticShutdown = automaticShutdown

	return err
}

func getService(address string, auth string) (service Service, err error) {
	automatic := &Automatic{
		address: address,
		auth:    auth,
	}
	cluster := strings.Contains(address, ";")
	if cluster {
		key := "redis-cluster-" + address + "-" + auth + "-" + fmt.Sprint(cluster)
		var automaticShutdown *worker.AutomaticShutdown
		automaticShutdown, err = worker.GetAutomaticShutdown(key, automatic.CreateClusterAutomaticShutdown)
		if err != nil {
			return
		}
		automaticShutdown.LastUseTimestamp = base.GetNowTime()
		service = automaticShutdown.Service.(*RedisClusterService)

	} else {
		key := "redis-" + address + "-" + auth
		var automaticShutdown *worker.AutomaticShutdown
		automaticShutdown, err = worker.GetAutomaticShutdown(key, automatic.CreateAutomaticShutdown)
		if err != nil {
			return
		}
		automaticShutdown.LastUseTimestamp = base.GetNowTime()
		service = automaticShutdown.Service.(*RedisPoolService)
	}
	return
}

type Service interface {
	Keys(pattern string, size int) (count int, keys []string, err error)
	Exists(key string) (value bool, err error)
	Get(key string) (valueInfo ValueInfo, err error)
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
