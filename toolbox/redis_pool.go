package toolbox

import (
	"sort"
	"time"

	"github.com/gomodule/redigo/redis"
)

func CreateRedisPoolService(address string, auth string) (service *RedisPoolService, err error) {
	service = &RedisPoolService{
		address: address,
		auth:    auth,
	}
	err = service.init()
	return
}

type RedisPoolService struct {
	address     string
	auth        string
	pool        *redis.Pool
	lastUseTime int64
}

func (service *RedisPoolService) init() (err error) {
	pool := &redis.Pool{
		MaxIdle:     2, //空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   3, //最大数
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", service.address)
			if err != nil {
				return nil, err
			}
			if service.auth != "" {
				if _, err := c.Do("auth", service.auth); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("ping")
			return err
		},
	}
	service.pool = pool
	return
}

func (this_ *RedisPoolService) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *RedisPoolService) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *RedisPoolService) Stop() {
	this_.pool.Close()
}

func (this_ *RedisPoolService) GetClient() redis.Conn {
	defer func() {
		this_.lastUseTime = GetNowTime()
	}()
	return this_.pool.Get()
}

func (this_ *RedisPoolService) Keys(pattern string, size int) (count int, keys []string, err error) {

	client := this_.GetClient()
	defer client.Close()
	var reply interface{}
	reply, err = client.Do("keys", pattern)
	if err != nil {
		return
	}
	if reply != nil {
		var list []string
		list, err = redis.Strings(reply, err)

		sor := sort.StringSlice(list)
		sor.Sort()

		count = len(list)
		if count <= size || size <= 0 {
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
	}
	return
}

func (this_ *RedisPoolService) KeyType(key string) (keyType string, err error) {

	client := this_.GetClient()
	defer client.Close()
	var reply interface{}
	reply, err = client.Do("type", key)
	if err != nil {
		return
	}
	if reply != nil {
		keyType, err = redis.String(reply, err)
	}
	return
}

func (this_ *RedisPoolService) Get(key string) (valueInfo RedisValueInfo, err error) {
	var keyType string
	keyType, err = this_.KeyType(key)
	if err != nil {
		return
	}
	var value interface{}
	client := this_.GetClient()
	defer client.Close()
	if keyType == "none" {

	} else if keyType == "string" {
		var reply interface{}
		reply, err = client.Do("get", key)
		if err != nil {
			return
		}
		if reply != nil {
			value, err = redis.String(reply, err)
		}
	} else if keyType == "list" {
		var reply interface{}
		reply, err = client.Do("llen", key)
		if err != nil {
			return
		}
		var len int64
		len, err = redis.Int64(reply, err)
		if err != nil {
			return
		}
		reply, err = client.Do("lrange", key, 0, len)
		if err != nil {
			return
		}
		value, err = redis.Strings(reply, err)
	} else if keyType == "set" {
		var reply interface{}
		reply, err = client.Do("smembers", key)
		if err != nil {
			return
		}
		value, err = redis.Strings(reply, err)
	} else if keyType == "hash" {
		var reply interface{}
		reply, err = client.Do("hgetall", key)
		if err != nil {
			return
		}
		value, err = redis.StringMap(reply, err)
	} else {
		println(keyType)
	}
	valueInfo.Type = keyType
	valueInfo.Value = value

	return
}

func (this_ *RedisPoolService) Set(key string, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()
	_, err = client.Do("set", key, value)
	return
}

func (this_ *RedisPoolService) Sadd(key string, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()
	_, err = client.Do("sadd", key, value)
	return
}

func (this_ *RedisPoolService) Srem(key string, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()
	_, err = client.Do("srem", key, value)
	return
}

func (this_ *RedisPoolService) Lpush(key string, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()
	_, err = client.Do("lpush", key, value)
	return
}

func (this_ *RedisPoolService) Rpush(key string, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()
	_, err = client.Do("rpush", key, value)
	return
}

func (this_ *RedisPoolService) Lset(key string, index int64, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()
	_, err = client.Do("lset", key, index, value)
	return
}

func (this_ *RedisPoolService) Lrem(key string, count int64, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()
	_, err = client.Do("lrem", key, count, value)
	return
}

func (this_ *RedisPoolService) Hset(key string, field string, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()
	_, err = client.Do("hset", key, field, value)
	return
}

func (this_ *RedisPoolService) Hdel(key string, field string) (err error) {

	client := this_.GetClient()
	defer client.Close()
	_, err = client.Do("hdel", key, field)
	return
}

func (this_ *RedisPoolService) Del(key string) (count int, err error) {
	count = 0
	client := this_.GetClient()
	defer client.Close()
	_, err = client.Do("del", key)
	if err == nil {
		count++
	}
	return
}

func (this_ *RedisPoolService) DelPattern(pattern string) (count int, err error) {
	count = 0
	var list []string
	_, list, err = this_.Keys(pattern, 0)

	client := this_.GetClient()
	defer client.Close()

	for _, key := range list {
		_, err = client.Do("del", key)
		if err == nil {
			count++
		} else {
			return
		}
	}
	return
}
