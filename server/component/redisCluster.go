package component

import (
	"context"
	"errors"
	"server/base"
	"strings"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
)

type RedisClusterService struct {
	address      string
	auth         string
	redisCluster *redisv8.ClusterClient
}

func CreateRedisClusterService(address string, auth string) (service *RedisClusterService, err error) {
	service = &RedisClusterService{
		address: address,
		auth:    auth,
	}
	err = service.init()
	return
}

func (service *RedisClusterService) init() (err error) {
	redisCluster := redisv8.NewClusterClient(&redisv8.ClusterOptions{
		Addrs:        strings.Split(service.address, ";"),
		DialTimeout:  100 * time.Second,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
		Password:     service.auth,
	})
	service.redisCluster = redisCluster
	return
}

func (service *RedisClusterService) Keys(pattern string, size int) (count int, keys []string, err error) {
	cmd := service.redisCluster.Keys(context.TODO(), pattern)
	var list []string
	list, err = cmd.Result()
	if err != nil {
		return
	}
	count = len(list)
	if count <= size {
		keys = list
	} else {
		for index, one := range list {
			if index < size {
				keys = append(keys, one)
			} else {
				break
			}
		}
	}
	return
}

func (service *RedisClusterService) KeyType(key string) (keyType string, err error) {
	cmd := service.redisCluster.Type(context.TODO(), key)
	keyType, err = cmd.Result()
	return
}

func (service *RedisClusterService) Exists(key string) (value bool, err error) {
	cmd := service.redisCluster.Exists(context.TODO(), key)
	var res int64
	res, err = cmd.Result()
	value = res == 1
	return
}
func (service *RedisClusterService) GetString(key string) (value string, err error) {
	cmd := service.redisCluster.Get(context.TODO(), key)
	value, err = cmd.Result()
	return
}

func (service *RedisClusterService) GetInt64(key string) (value int64, err error) {
	cmd := service.redisCluster.Get(context.TODO(), key)
	value, err = cmd.Int64()
	return
}

func (service *RedisClusterService) IncrBy(key string, num int64) (value int64, err error) {
	cmd := service.redisCluster.IncrBy(context.TODO(), key, num)
	value, err = cmd.Result()
	return
}

func (service *RedisClusterService) Set(key string, value string) (err error) {
	cmd := service.redisCluster.Set(context.TODO(), key, value, time.Duration(0))
	_, err = cmd.Result()
	return
}

func (service *RedisClusterService) SetInt64(key string, value int64) (err error) {
	cmd := service.redisCluster.Set(context.TODO(), key, value, time.Duration(0))
	_, err = cmd.Result()
	return
}

func (service *RedisClusterService) Sadd(key string, value string) (err error) {
	cmd := service.redisCluster.SAdd(context.TODO(), key, value)
	_, err = cmd.Result()
	return
}

func (service *RedisClusterService) Srem(key string, value string) (err error) {
	cmd := service.redisCluster.SRem(context.TODO(), key, value)
	_, err = cmd.Result()
	return
}

func (service *RedisClusterService) Lpush(key string, value string) (err error) {
	cmd := service.redisCluster.LPush(context.TODO(), key, value)
	_, err = cmd.Result()
	return
}

func (service *RedisClusterService) Rpush(key string, value string) (err error) {
	cmd := service.redisCluster.RPush(context.TODO(), key, value)
	_, err = cmd.Result()
	return
}

func (service *RedisClusterService) Lset(key string, index int64, value string) (err error) {
	cmd := service.redisCluster.LSet(context.TODO(), key, index, value)
	_, err = cmd.Result()
	return
}

func (service *RedisClusterService) Lrem(key string, count int64, value string) (err error) {
	cmd := service.redisCluster.LRem(context.TODO(), key, count, value)
	_, err = cmd.Result()
	return
}

func (service *RedisClusterService) Hset(key string, field string, value string) (err error) {
	cmd := service.redisCluster.HSet(context.TODO(), key, field, value)
	_, err = cmd.Result()
	return
}

func (service *RedisClusterService) Hdel(key string, field string) (err error) {
	cmd := service.redisCluster.HDel(context.TODO(), key, field)
	_, err = cmd.Result()
	return
}

func (service *RedisClusterService) Del(key string) (count int, err error) {
	count = 0
	cmd := service.redisCluster.Del(context.TODO(), key)
	_, err = cmd.Result()
	if err == nil {
		count++
	}
	return
}

func (service *RedisClusterService) DelPattern(pattern string) (count int, err error) {
	count = 0
	cmd := service.redisCluster.Keys(context.TODO(), pattern)
	var list []string
	list, err = cmd.Result()
	if err != nil {
		return
	}
	for _, key := range list {
		var num int
		num, err = service.Del(key)
		if err != nil {
			return
		}
		count += num
	}
	return
}

func (service *RedisClusterService) Lock(key string, expire int, timeout int64) (unlock func() (err error), err error) {
	value := base.GenerateUUID()
	start := base.GetNowTime()
	var wait int = 5
	for {
		var res int
		cmd := clusterUpdateLockExpireUidScript.Run(context.TODO(), service.redisCluster, []string{key}, value, expire)

		res, err = cmd.Int()
		if err != nil {
			break
		}
		if res == 1 {
			break
		}
		end := base.GetNowTime()
		if (end - start) >= int64(timeout-5) {
			err = errors.New("Lock timeout")
			break
		}
		time.Sleep(time.Duration(wait) * time.Millisecond)
	}
	if err != nil {
		return
	}
	redisLock := &clusterRedisLock{
		Key:     key,
		Value:   value,
		Locked:  true,
		service: service,
	}
	unlock = redisLock.Unlock
	return
}

type clusterRedisLock struct {
	Key     string
	Value   string
	Locked  bool
	service *RedisClusterService
}

func (redisLock *clusterRedisLock) Unlock() (err error) {
	if !redisLock.Locked {
		return
	}
	redisLock.Locked = false
	cmd := clusterUpdateLockExpireUidScript.Run(context.TODO(), redisLock.service.redisCluster, []string{redisLock.Key}, redisLock.Value)
	_ = cmd.String()
	if err != nil {
		return
	}
	return
}

var (
	clusterUpdateLockExpireUidScript = redisv8.NewScript(`
		local res = redis.call("SETNX", KEYS[1], ARGV[1]) 
		if res == 1 then
			return redis.call("EXPIRE", KEYS[1], ARGV[2])
		end
		return res
	`)
	// clusterDeleteLockByUidScript = redisv8.NewScript(`
	// 	local res = redis.call("GET", KEYS[1])
	// 	if res == ARGV[1] then
	// 		return redis.call("DEL", KEYS[1])
	// 	end
	// 	return res
	// `)
)
