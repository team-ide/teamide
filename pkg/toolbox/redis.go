package toolbox

import (
	"encoding/json"
	"strings"
)

func GetRedisWorker() *Worker {
	worker_ := &Worker{
		Name:    "redis",
		Text:    "Redis",
		WorkMap: map[string]func(map[string]interface{}) (map[string]interface{}, error){},
	}
	worker_.WorkMap["get"] = func(m map[string]interface{}) (map[string]interface{}, error) {
		return redisWork("get", m["config"].(map[string]interface{}), m["data"].(map[string]interface{}))
	}

	return worker_
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

func redisWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {
	var service RedisService
	var address string = config["address"].(string)
	var auth string
	if config["auth"] != nil {
		auth = config["auth"].(string)
	}
	service, err = getRedisService(address, auth)
	if err != nil {
		return
	}

	var bs []byte
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
		res["value"] = valueInfo
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

func getRedisService(address string, auth string) (res RedisService, err error) {
	key := "redis-" + address + "-" + auth
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s RedisService
		s, err = CreateRedisService(address, auth)
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

func CreateRedisService(address string, auth string) (service RedisService, err error) {
	if !strings.Contains(address, ",") && !strings.Contains(address, ";") {
		service, err = CreateRedisPoolService(address, auth)
	} else {
		var servers []string
		if strings.Contains(address, ",") {
			servers = strings.Split(address, ",")
		} else if strings.Contains(address, ";") {
			servers = strings.Split(address, ";")
		} else {
			servers = []string{address}
		}
		service, err = CreateRedisClusterService(servers, auth)
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
