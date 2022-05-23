package toolbox

import (
	"context"
	"encoding/json"
	"strings"
	"teamide/pkg/form"
	"teamide/pkg/util"
)

func init() {
	worker_ := &Worker{
		Name: "redis",
		Text: "Redis",
		Work: redisWork,
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{Label: "连接地址（127.0.0.1:6379）", Name: "address", DefaultValue: "127.0.0.1:6379",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
				{Label: "密码", Name: "auth", Type: "password"},
			},
		},
	}

	AddWorker(worker_)
}

type RedisBaseRequest struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	ValueSize  int64  `json:"valueSize"`
	ValueStart int64  `json:"valueStart"`
	Pattern    string `json:"pattern"`
	Database   int    `json:"database"`
	Size       int    `json:"size"`
	Type       string `json:"type"`
	Index      int64  `json:"index"`
	Count      int64  `json:"count"`
	Field      string `json:"field"`
	Cursor     uint64 `json:"cursor"`
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

	ctx := context.TODO()
	res = map[string]interface{}{}
	switch work {
	case "get":
		var valueInfo RedisValueInfo
		valueInfo, err = service.Get(ctx, request.Database, request.Key, request.ValueStart, request.ValueSize, request.Cursor)
		res["database"] = request.Database
		res["key"] = request.Key
		res["type"] = valueInfo.Type
		res["value"] = valueInfo.Value
	case "keys":
		var count int
		var keys []string
		count, keys, err = service.Keys(ctx, request.Database, request.Pattern, request.Size)
		if err != nil {
			return
		}

		var dataList []map[string]interface{}
		for _, key := range keys {
			var one = map[string]interface{}{}
			one["key"] = key
			one["database"] = request.Database
			dataList = append(dataList, one)
		}
		res["count"] = count
		res["dataList"] = dataList
	case "do":
		if request.Type == "set" {
			err = service.Set(ctx, request.Database, request.Key, request.Value)
		} else if request.Type == "sadd" {
			err = service.Sadd(ctx, request.Database, request.Key, request.Value)
		} else if request.Type == "srem" {
			err = service.Srem(ctx, request.Database, request.Key, request.Value)
		} else if request.Type == "lpush" {
			err = service.Lpush(ctx, request.Database, request.Key, request.Value)
		} else if request.Type == "rpush" {
			err = service.Rpush(ctx, request.Database, request.Key, request.Value)
		} else if request.Type == "lset" {
			err = service.Lset(ctx, request.Database, request.Key, request.Index, request.Value)
		} else if request.Type == "lrem" {
			err = service.Lrem(ctx, request.Database, request.Key, request.Count, request.Value)
		} else if request.Type == "hset" {
			err = service.Hset(ctx, request.Database, request.Key, request.Field, request.Value)
		} else if request.Type == "hdel" {
			err = service.Hdel(ctx, request.Database, request.Key, request.Field)
		}
		if err != nil {
			return
		}
	case "delete":
		var count int
		count, err = service.Del(ctx, request.Database, request.Key)
		res["count"] = count
	case "deletePattern":
		var count int
		count, err = service.DelPattern(ctx, request.Database, request.Pattern)
		res["count"] = count
	}
	return
}

func getRedisService(redisConfig RedisConfig) (res RedisService, err error) {
	key := "redis-" + redisConfig.Address
	if redisConfig.Auth != "" {
		key += "-" + util.GetMd5String(key+redisConfig.Auth)
	}
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s RedisService
		s, err = CreateRedisService(redisConfig)
		if err != nil {
			return
		}
		_, err = s.Get(context.TODO(), 0, "_", 0, 0, 0)
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
	Type       string      `json:"type"`
	Value      interface{} `json:"value"`
	ValueCount int64       `json:"valueCount"`
	ValueStart int64       `json:"valueStart"`
	ValueEnd   int64       `json:"valueEnd"`
	Cursor     uint64      `json:"cursor"`
}

type RedisService interface {
	GetWaitTime() int64
	GetLastUseTime() int64
	Stop()
	Keys(ctx context.Context, database int, pattern string, size int) (count int, keys []string, err error)
	Get(ctx context.Context, database int, key string, valueStart, valueSize int64, cursor uint64) (valueInfo RedisValueInfo, err error)
	Set(ctx context.Context, database int, key string, value string) (err error)
	Sadd(ctx context.Context, database int, key string, value string) (err error)
	Srem(ctx context.Context, database int, key string, value string) (err error)
	Lpush(ctx context.Context, database int, key string, value string) (err error)
	Rpush(ctx context.Context, database int, key string, value string) (err error)
	Lset(ctx context.Context, database int, key string, index int64, value string) (err error)
	Lrem(ctx context.Context, database int, key string, count int64, value string) (err error)
	Hset(ctx context.Context, database int, key string, field string, value string) (err error)
	Hdel(ctx context.Context, database int, key string, field string) (err error)
	Del(ctx context.Context, database int, key string) (count int, err error)
	DelPattern(ctx context.Context, database int, key string) (count int, err error)
}
