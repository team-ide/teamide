package toolbox

import (
	"encoding/json"
	"strings"
)

func init() {
	worker_ := &Worker{
		Name: "redis",
		Text: "Redis",
		Work: redisWork,
	}

	AddWorker(worker_)
}

type RedisBaseRequest struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Pattern string `json:"pattern"`
	Size    int    `json:"size"`
	Type    string `json:"type"`
	Index   int64  `json:"index"`
	Count   int64  `json:"count"`
	Field   string `json:"field"`
}

type RedisConfig struct {
	Address string `json:"address"`
	Auth    string `json:"auth"`
}

func redisWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {

	var redisConfig RedisConfig
	var bs []byte
	bs, err = json.Marshal(config)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, &redisConfig)
	if err != nil {
		return
	}

	var service RedisService
	service, err = getRedisService(redisConfig)
	if err != nil {
		return
	}

	bs, err = json.Marshal(data)
	if err != nil {
		return
	}
	request := &RedisBaseRequest{}
	err = json.Unmarshal(bs, request)
	if err != nil {
		return
	}

	res = map[string]interface{}{}
	switch work {
	case "get":
		var valueInfo RedisValueInfo
		valueInfo, err = service.Get(request.Key)
		res["type"] = valueInfo.Type
		res["value"] = valueInfo.Value
	case "keys":
		var count int
		var keys []string
		count, keys, err = service.Keys(request.Pattern, request.Size)
		if err != nil {
			return
		}
		res["count"] = count
		res["keys"] = keys
	case "do":
		if request.Type == "set" {
			err = service.Set(request.Key, request.Value)
		} else if request.Type == "sadd" {
			err = service.Sadd(request.Key, request.Value)
		} else if request.Type == "srem" {
			err = service.Srem(request.Key, request.Value)
		} else if request.Type == "lpush" {
			err = service.Lpush(request.Key, request.Value)
		} else if request.Type == "rpush" {
			err = service.Rpush(request.Key, request.Value)
		} else if request.Type == "lset" {
			err = service.Lset(request.Key, request.Index, request.Value)
		} else if request.Type == "lrem" {
			err = service.Lrem(request.Key, request.Count, request.Value)
		} else if request.Type == "hset" {
			err = service.Hset(request.Key, request.Field, request.Value)
		} else if request.Type == "hdel" {
			err = service.Hdel(request.Key, request.Field)
		}
		if err != nil {
			return
		}
	case "delete":
		var count int
		count, err = service.Del(request.Key)
		res["count"] = count
	case "deletePattern":
		var count int
		count, err = service.DelPattern(request.Pattern)
		res["count"] = count
	}
	return
}

func getRedisService(redisConfig RedisConfig) (res RedisService, err error) {
	key := "redis-" + redisConfig.Address + "-" + redisConfig.Auth
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s RedisService
		s, err = CreateRedisService(redisConfig)
		if err != nil {
			return
		}
		_, err = s.Get("_")
		if err != nil {
			return
		}
		res = s
		return
	})
	if err != nil {
		return
	}
	res = service.(RedisService)
	return
}

func CreateRedisService(redisConfig RedisConfig) (service RedisService, err error) {
	if !strings.Contains(redisConfig.Address, ",") && !strings.Contains(redisConfig.Address, ";") {
		service, err = CreateRedisPoolService(redisConfig.Address, redisConfig.Auth)
	} else {
		var servers []string
		if strings.Contains(redisConfig.Address, ",") {
			servers = strings.Split(redisConfig.Address, ",")
		} else if strings.Contains(redisConfig.Address, ";") {
			servers = strings.Split(redisConfig.Address, ";")
		} else {
			servers = []string{redisConfig.Address}
		}
		service, err = CreateRedisClusterService(servers, redisConfig.Auth)
	}
	return
}

type RedisValueInfo struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type RedisService interface {
	GetWaitTime() int64
	GetLastUseTime() int64
	Stop()
	Keys(pattern string, size int) (count int, keys []string, err error)
	Get(key string) (valueInfo RedisValueInfo, err error)
	Set(key string, value string) (err error)
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
}
