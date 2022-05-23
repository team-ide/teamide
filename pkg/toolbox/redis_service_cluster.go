package toolbox

import (
	"context"
	"go.uber.org/zap"
	"sort"
	"teamide/pkg/util"
	"time"

	"github.com/go-redis/redis/v8"
)

func CreateRedisClusterService(servers []string, auth string) (service *RedisClusterService, err error) {
	service = &RedisClusterService{
		servers: servers,
		auth:    auth,
	}
	err = service.init()
	return
}

type RedisClusterService struct {
	servers      []string
	auth         string
	redisCluster *redis.ClusterClient
	lastUseTime  int64
}

func (this_ *RedisClusterService) init() (err error) {
	redisCluster := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        this_.servers,
		DialTimeout:  100 * time.Second,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
		Password:     this_.auth,
	})
	this_.redisCluster = redisCluster
	return
}

func (this_ *RedisClusterService) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *RedisClusterService) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *RedisClusterService) Stop() {
	_ = this_.redisCluster.Close()
}

func (this_ *RedisClusterService) GetClient(ctx context.Context, database int) (redisCluster *redis.ClusterClient, err error) {
	defer func() {
		this_.lastUseTime = GetNowTime()
	}()
	redisCluster = this_.redisCluster
	if ctx != nil && database >= 0 {
		return
	}
	return
}

func (this_ *RedisClusterService) SelectDatabase(ctx context.Context, database int) (err error) {
	if ctx != nil && database >= 0 {
		return
	}
	return
}

func (this_ *RedisClusterService) Keys(ctx context.Context, database int, pattern string, size int64) (count int, keys []string, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	var list []string
	err = client.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) (err error) {

		var ls []string
		cmd := client.Keys(ctx, pattern)
		ls, err = cmd.Result()
		if err != nil {
			return
		}
		count += len(ls)
		list = append(list, ls...)
		return
	})
	sor := sort.StringSlice(list)
	sor.Sort()
	listCount := len(list)
	if int64(listCount) <= size || size < 0 {
		keys = list
	} else {
		keys = list[0:size]
	}
	return
}

func (this_ *RedisClusterService) KeyType(ctx context.Context, database int, key string) (keyType string, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.Type(ctx, key)
	keyType, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) Get(ctx context.Context, database int, key string, valueStart, valueSize int64) (valueInfo RedisValueInfo, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	var keyType string
	keyType, err = this_.KeyType(ctx, database, key)
	if err != nil {
		return
	}
	var value interface{}

	if keyType == "none" {

	} else if keyType == "string" {
		cmd := client.Get(ctx, key)
		value, err = cmd.Result()
		if err != nil {
			util.Logger.Error("Get Error", zap.Any("key", key), zap.Error(err))
			return
		}
	} else if keyType == "list" {

		cmd := client.LLen(ctx, key)

		valueInfo.ValueCount, err = cmd.Result()
		if err != nil {
			util.Logger.Error("LLen Error", zap.Any("key", key), zap.Error(err))
			return
		}
		valueInfo.ValueStart = valueStart
		valueInfo.ValueEnd = valueInfo.ValueStart + valueSize

		var list []string
		cmdRange := client.LRange(ctx, key, valueInfo.ValueStart, valueInfo.ValueEnd)
		list, err = cmdRange.Result()
		if err != nil {
			util.Logger.Error("LRange Error", zap.Any("key", key), zap.Error(err))
			return
		}

		if int64(len(list)) <= valueSize || valueSize < 0 {
			value = list
		} else {
			value = list[0:valueSize]
		}

	} else if keyType == "set" {

		cmdSCard := client.SCard(ctx, key)
		valueInfo.ValueCount, err = cmdSCard.Result()
		if err != nil {
			util.Logger.Error("SCard Error", zap.Any("key", key), zap.Error(err))
			return
		}
		valueInfo.ValueStart = valueStart
		valueInfo.ValueEnd = valueInfo.ValueStart + valueSize

		var list []string
		cmd := client.SScan(ctx, key, uint64(valueInfo.ValueStart), "*", valueInfo.ValueEnd)
		list, valueInfo.Cursor, err = cmd.Result()
		if err != nil {
			util.Logger.Error("SScan Error", zap.Any("key", key), zap.Error(err))
			return
		}

		if int64(len(list)) <= valueSize || valueSize < 0 {
			value = list
		} else {
			value = list[0:valueSize]
		}
	} else if keyType == "hash" {

		cmdHLen := client.HLen(ctx, key)
		valueInfo.ValueCount, err = cmdHLen.Result()
		if err != nil {
			util.Logger.Error("HLen Error", zap.Any("key", key), zap.Error(err))
			return
		}
		valueInfo.ValueStart = valueStart * 2
		valueInfo.ValueEnd = valueInfo.ValueStart + valueSize*2

		cmdHScan := client.HScan(ctx, key, uint64(valueInfo.ValueStart), "*", valueInfo.ValueEnd)

		var keyValueList []string
		keyValueList, valueInfo.Cursor, err = cmdHScan.Result()
		if err != nil {
			util.Logger.Error("HScan Error", zap.Any("key", key), zap.Error(err))
			return
		}
		var keyValueListSize = int64(len(keyValueList))
		var keyValue = map[string]string{}
		for i := int64(0); i < valueSize*2; i++ {
			if i >= keyValueListSize {
				break
			}
			filed := keyValueList[i]
			filedValue := ""
			if i+1 < keyValueListSize {
				filedValue = keyValueList[i+1]
			}
			keyValue[filed] = filedValue
			i++
		}

		value = keyValue
	} else {
		println(keyType)
	}
	valueInfo.Type = keyType
	valueInfo.Value = value
	return
}

func (this_ *RedisClusterService) Set(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.Set(ctx, key, value, time.Duration(0))
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) SAdd(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.SAdd(ctx, key, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) SRem(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.SRem(ctx, key, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) LPush(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.LPush(ctx, key, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) RPush(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.RPush(ctx, key, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) LSet(ctx context.Context, database int, key string, index int64, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.LSet(ctx, key, index, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) LRem(ctx context.Context, database int, key string, count int64, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.LRem(ctx, key, count, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) HSet(ctx context.Context, database int, key string, field string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.HSet(ctx, key, field, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) HDel(ctx context.Context, database int, key string, field string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.HDel(ctx, key, field)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) Del(ctx context.Context, database int, key string) (count int, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	count = 0
	cmd := client.Del(ctx, key)
	_, err = cmd.Result()
	if err == nil {
		count++
	}
	return
}

func (this_ *RedisClusterService) DelPattern(ctx context.Context, database int, pattern string) (count int, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	var keys []string
	_, keys, err = this_.Keys(ctx, database, pattern, -1)
	if err != nil {
		return
	}

	count = 0
	for _, key := range keys {
		cmd := client.Del(ctx, key)
		_, err = cmd.Result()
		if err == nil {
			count++
		}
	}
	return
}
