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
	Size       int64  `json:"size"`
	DoType     string `json:"doType"`
	Index      int64  `json:"index"`
	Count      int64  `json:"count"`
	Field      string `json:"field"`
	TaskKey    string `json:"taskKey,omitempty"`
}

type RedisConfig struct {
	Address string `json:"address"`
	Auth    string `json:"auth"`
}

func redisWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {

	var redisConfig RedisConfig
	var configBS []byte
	configBS, err = json.Marshal(config)
	if err != nil {
		return
	}
	err = json.Unmarshal(configBS, &redisConfig)
	if err != nil {
		return
	}

	var service RedisService
	service, err = getRedisService(redisConfig)
	if err != nil {
		return
	}

	var dataBS []byte
	dataBS, err = json.Marshal(data)
	if err != nil {
		return
	}
	request := &RedisBaseRequest{}
	err = json.Unmarshal(dataBS, request)
	if err != nil {
		return
	}

	ctx := context.TODO()
	res = map[string]interface{}{}
	switch work {
	case "get":
		var valueInfo RedisValueInfo
		valueInfo, err = service.Get(ctx, request.Database, request.Key, request.ValueStart, request.ValueSize)
		res["database"] = request.Database
		res["key"] = request.Key
		res["valueType"] = valueInfo.ValueType
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
		switch request.DoType {
		case "set":
			err = service.Set(ctx, request.Database, request.Key, request.Value)
		case "SAdd":
			err = service.SAdd(ctx, request.Database, request.Key, request.Value)
		case "SRem":
			err = service.SRem(ctx, request.Database, request.Key, request.Value)
		case "LPush":
			err = service.LPush(ctx, request.Database, request.Key, request.Value)
		case "RPush":
			err = service.RPush(ctx, request.Database, request.Key, request.Value)
		case "LSet":
			err = service.LSet(ctx, request.Database, request.Key, request.Index, request.Value)
		case "LRem":
			err = service.LRem(ctx, request.Database, request.Key, request.Count, request.Value)
		case "HSet":
			err = service.HSet(ctx, request.Database, request.Key, request.Field, request.Value)
		case "HDel":
			err = service.HDel(ctx, request.Database, request.Key, request.Field)

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

	case "import":

		taskKey := util.GenerateUUID()

		var redisImportTask = &redisImportTask{}
		err = json.Unmarshal(dataBS, redisImportTask)
		if err != nil {
			return
		}

		redisImportTask.request = request
		redisImportTask.Key = taskKey
		redisImportTask.service = service

		addRedisImportTask(redisImportTask)

		res["taskKey"] = taskKey
	case "importStatus":
		redisImportTask := redisImportTaskCache[request.TaskKey]
		res["task"] = redisImportTask
	case "importStop":
		redisImportTask := redisImportTaskCache[request.TaskKey]
		if redisImportTask != nil {
			redisImportTask.Stop()
		}
	case "importClean":
		delete(redisImportTaskCache, request.TaskKey)
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

		ctx := context.TODO()
		_, err = s.Get(ctx, 0, "_", 0, 0)
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
	ValueType  string      `json:"valueType"`
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
	Keys(ctx context.Context, database int, pattern string, size int64) (count int, keys []string, err error)
	Get(ctx context.Context, database int, key string, valueStart, valueSize int64) (valueInfo RedisValueInfo, err error)
	Set(ctx context.Context, database int, key string, value string) (err error)
	SAdd(ctx context.Context, database int, key string, value string) (err error)
	SRem(ctx context.Context, database int, key string, value string) (err error)
	LPush(ctx context.Context, database int, key string, value string) (err error)
	RPush(ctx context.Context, database int, key string, value string) (err error)
	LSet(ctx context.Context, database int, key string, index int64, value string) (err error)
	LRem(ctx context.Context, database int, key string, count int64, value string) (err error)
	HSet(ctx context.Context, database int, key string, field string, value string) (err error)
	HDel(ctx context.Context, database int, key string, field string) (err error)
	Del(ctx context.Context, database int, key string) (count int, err error)
	DelPattern(ctx context.Context, database int, key string) (count int, err error)
}
