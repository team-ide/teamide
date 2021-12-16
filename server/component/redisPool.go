package component

import (
	"errors"
	"fmt"
	"server/base"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

type RedisPoolService struct {
	address string
	auth    string
	pool    *redigo.Pool
}

func CreateRedisPoolService(address string, auth string) (service *RedisPoolService, err error) {
	service = &RedisPoolService{
		address: address,
		auth:    auth,
	}
	err = service.init()
	return
}

func (service *RedisPoolService) init() (err error) {
	pool := &redigo.Pool{
		MaxIdle:     10, //空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   100, //最大数
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", service.address)
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
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("ping")
			return err
		},
		Wait: true,
	}
	service.pool = pool
	return
}

func (service *RedisPoolService) Keys(pattern string, size int) (count int, keys []string, err error) {

	client := service.pool.Get()
	defer client.Close()
	var reply interface{}
	reply, err = client.Do("keys", pattern)
	if err != nil {
		return
	}
	if reply != nil {
		var list []string
		list, err = redigo.Strings(reply, err)
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

func (service *RedisPoolService) KeyType(key string) (keyType string, err error) {

	client := service.pool.Get()
	defer client.Close()
	var reply interface{}
	reply, err = client.Do("type", key)
	if err != nil {
		return
	}
	if reply != nil {
		keyType, err = redigo.String(reply, err)
	}
	return
}

func (service *RedisPoolService) Exists(key string) (value bool, err error) {
	client := service.pool.Get()
	defer client.Close()
	var reply interface{}
	reply, err = client.Do("exists", key)
	if err != nil {
		return
	}
	if reply != nil {
		var count int
		count, err = redigo.Int(reply, err)
		value = count == 1
	}

	return
}

func (service *RedisPoolService) GetString(key string) (value string, err error) {
	client := service.pool.Get()
	defer client.Close()
	var reply interface{}
	reply, err = client.Do("get", key)
	if err != nil {
		return
	}
	if reply != nil {
		value, err = redigo.String(reply, err)
	}

	return
}

func (service *RedisPoolService) GetInt64(key string) (value int64, err error) {
	client := service.pool.Get()
	defer client.Close()
	var reply interface{}
	reply, err = client.Do("get", key)
	if err != nil {
		return
	}
	if reply != nil {
		value, err = redigo.Int64(reply, err)
	}

	return
}

func (service *RedisPoolService) IncrBy(key string, num int64) (value int64, err error) {
	client := service.pool.Get()
	defer client.Close()
	var reply interface{}
	reply, err = client.Do("incrby", key, num)
	if err != nil {
		return
	}
	if reply != nil {
		value, err = redigo.Int64(reply, err)
	}

	return
}

func (service *RedisPoolService) Expire(key string, second int64) (value int, err error) {
	client := service.pool.Get()
	defer client.Close()
	var reply interface{}
	reply, err = client.Do("expire", key, second*1000)
	if err != nil {
		return
	}
	if reply != nil {
		value, err = redigo.Int(reply, err)
	}

	return
}

func (service *RedisPoolService) Set(key string, value string) (err error) {

	client := service.pool.Get()
	defer client.Close()
	_, err = client.Do("set", key, value)
	return
}

func (service *RedisPoolService) SetInt64(key string, value int64) (err error) {

	client := service.pool.Get()
	defer client.Close()
	_, err = client.Do("set", key, fmt.Sprint(value))
	return
}

func (service *RedisPoolService) Sadd(key string, value string) (err error) {

	client := service.pool.Get()
	defer client.Close()
	_, err = client.Do("sadd", key, value)
	return
}

func (service *RedisPoolService) Srem(key string, value string) (err error) {

	client := service.pool.Get()
	defer client.Close()
	_, err = client.Do("srem", key, value)
	return
}

func (service *RedisPoolService) Lpush(key string, value string) (err error) {

	client := service.pool.Get()
	defer client.Close()
	_, err = client.Do("lpush", key, value)
	return
}

func (service *RedisPoolService) Rpush(key string, value string) (err error) {

	client := service.pool.Get()
	defer client.Close()
	_, err = client.Do("rpush", key, value)
	return
}

func (service *RedisPoolService) Lset(key string, index int64, value string) (err error) {

	client := service.pool.Get()
	defer client.Close()
	_, err = client.Do("lset", key, index, value)
	return
}

func (service *RedisPoolService) Lrem(key string, count int64, value string) (err error) {

	client := service.pool.Get()
	defer client.Close()
	_, err = client.Do("lrem", key, count, value)
	return
}

func (service *RedisPoolService) Hset(key string, field string, value string) (err error) {

	client := service.pool.Get()
	defer client.Close()
	_, err = client.Do("hset", key, field, value)
	return
}

func (service *RedisPoolService) Hdel(key string, field string) (err error) {

	client := service.pool.Get()
	defer client.Close()
	_, err = client.Do("hdel", key, field)
	return
}

func (service *RedisPoolService) Del(key string) (count int, err error) {
	count = 0
	client := service.pool.Get()
	defer client.Close()
	_, err = client.Do("del", key)
	if err == nil {
		count++
	}
	return
}

func (service *RedisPoolService) DelPattern(pattern string) (count int, err error) {
	count = 0
	var list []string
	_, list, err = service.Keys(pattern, 0)

	client := service.pool.Get()
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

func (service *RedisPoolService) Lock(key string, expire int, timeout int64) (unlock func() (err error), err error) {
	value := base.GenerateUUID()
	client := service.pool.Get()
	defer client.Close()
	start := base.GetNowTime()
	var wait int = 5
	for {
		var res interface{}
		res, err = redigoUpdateLockExpireUidScript.Do(client, key, value, expire)
		if err != nil {
			break
		}
		if res.(int64) == 1 {
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
	redisLock := &RedigoRedisLock{
		Key:     key,
		Value:   value,
		Locked:  true,
		service: service,
	}
	unlock = redisLock.Unlock
	return
}

type RedigoRedisLock struct {
	Key     string
	Value   string
	Locked  bool
	service *RedisPoolService
}

func (redisLock *RedigoRedisLock) Unlock() (err error) {
	if !redisLock.Locked {
		return
	}
	redisLock.Locked = false
	client := redisLock.service.pool.Get()
	defer client.Close()
	_, err = redigoDeleteLockByUidScript.Do(client, redisLock.Key, redisLock.Value)
	if err != nil {
		return
	}
	return
}

var (
	redigoUpdateLockExpireUidScript = redigo.NewScript(1, `
		local res = redis.call("SETNX", KEYS[1], ARGV[1]) 
		if res == 1 then
			return redis.call("EXPIRE", KEYS[1], ARGV[2])
		end
		return res
	`)
	redigoDeleteLockByUidScript = redigo.NewScript(1, `
		local res = redis.call("GET", KEYS[1]) 
		if res == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		end
		return res 
	`)
)
