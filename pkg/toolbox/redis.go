package toolbox

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"teamide/pkg/redis"
	"teamide/pkg/util"
)

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
	Expire     int64  `json:"expire"`
}

func RedisWork(work string, config *redis.Config, data map[string]interface{}) (res map[string]interface{}, err error) {

	var service redis.Service
	service, err = getRedisService(*config)
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
	case "info":
		var info string
		info, err = service.Info(ctx)
		if err != nil {
			return
		}
		res["info"] = info
	case "get":
		var valueInfo *redis.ValueInfo
		valueInfo, err = service.Get(ctx, request.Database, request.Key, request.ValueStart, request.ValueSize)
		res["database"] = request.Database
		res["key"] = request.Key
		res["valueType"] = valueInfo.ValueType
		res["value"] = valueInfo.Value
		res["memoryUsage"] = valueInfo.MemoryUsage
		res["valueStart"] = valueInfo.ValueStart
		res["valueEnd"] = valueInfo.ValueEnd
		res["valueCount"] = valueInfo.ValueCount
		res["ttl"] = valueInfo.TTL
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
			if err != nil {
				return
			}
			if request.Expire > 0 {
				_, err = service.Expire(ctx, request.Database, request.Key, request.Expire)
			}
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
	case "expire":
		var success bool
		success, err = service.Expire(ctx, request.Database, request.Key, request.Expire)
		res["success"] = success
	case "ttl":
		var ttl int64
		ttl, err = service.TTL(ctx, request.Database, request.Key)
		res["ttl"] = ttl
	case "persist":
		var success bool
		success, err = service.Persist(ctx, request.Database, request.Key)
		res["success"] = success

	case "import":

		taskKey := util.UUID()

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
	if redisConfig.Username != "" {
		key += "-" + util.GetMd5String(key+redisConfig.Username)
	}
	if redisConfig.Auth != "" {
		key += "-" + util.GetMd5String(key+redisConfig.Auth)
	}
	if redisConfig.CertPath != "" {
		key += "-" + util.GetMd5String(key+redisConfig.CertPath)
	}
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s redis.Service
		s, err = redis.CreateService(redisConfig)
		if err != nil {
			util.Logger.Error("getRedisService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Stop()
			}
			return
		}

		ctx := context.TODO()
		_, err = s.Exists(ctx, 0, "_")
		if err != nil {
			util.Logger.Error("getRedisService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Stop()
			}
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
