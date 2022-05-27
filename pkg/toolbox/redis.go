package toolbox

import (
	"context"
	"encoding/json"
	"teamide/pkg/form"
	"teamide/pkg/redis"
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

func redisWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {

	var redisConfig redis.Config
	var configBS []byte
	configBS, err = json.Marshal(config)
	if err != nil {
		return
	}
	err = json.Unmarshal(configBS, &redisConfig)
	if err != nil {
		return
	}

	var service redis.Service
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
		var valueInfo *redis.ValueInfo
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

		var importTask = &redis.ImportTask{}
		err = json.Unmarshal(dataBS, importTask)
		if err != nil {
			return
		}

		importTask.Key = taskKey
		importTask.Service = service

		redis.StartImportTask(importTask)

		res["taskKey"] = taskKey
	case "importStatus":
		task := redis.GetImportTask(request.TaskKey)
		res["task"] = task
	case "importStop":
		redis.StopImportTask(request.TaskKey)
	case "importClean":
		redis.CleanImportTask(request.TaskKey)
	}
	return
}

func getRedisService(redisConfig redis.Config) (res redis.Service, err error) {
	key := "redis-" + redisConfig.Address
	if redisConfig.Auth != "" {
		key += "-" + util.GetMd5String(key+redisConfig.Auth)
	}
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s redis.Service
		s, err = redis.CreateService(redisConfig)
		if err != nil {
			return
		}

		ctx := context.TODO()
		_, err = s.Exists(ctx, 0, "_")
		if err != nil {
			return
		}
		res = s
		return
	})
	if err != nil {
		return
	}
	res = service.(redis.Service)
	return
}
