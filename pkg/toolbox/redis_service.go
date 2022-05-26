package toolbox

import (
	"context"
	"go.uber.org/zap"
	"sort"
	"teamide/pkg/util"
	"time"

	"github.com/go-redis/redis/v8"
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
	client      *redis.Client
	lastUseTime int64
}

func (this_ *RedisPoolService) init() (err error) {
	client := redis.NewClient(&redis.Options{
		Addr:         this_.address,
		DialTimeout:  100 * time.Second,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
		Password:     this_.auth,
	})
	this_.client = client
	return
}

func (this_ *RedisPoolService) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *RedisPoolService) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *RedisPoolService) Stop() {
	_ = this_.client.Close()
}

func (this_ *RedisPoolService) GetClient(ctx context.Context, database int) (client *redis.Client, err error) {
	defer func() {
		this_.lastUseTime = GetNowTime()
	}()
	cmd := this_.client.Do(ctx, "select", database)
	_, err = cmd.Result()
	if err != nil {
		return
	}
	client = this_.client
	return
}

func (this_ *RedisPoolService) Keys(ctx context.Context, database int, pattern string, size int64) (count int, keys []string, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	var list []string
	cmdKeys := client.Keys(ctx, pattern)
	list, err = cmdKeys.Result()
	if err != nil {
		return
	}
	count = len(list)

	sor := sort.StringSlice(list)
	sor.Sort()

	if int64(count) <= size || size < 0 {
		keys = list
	} else {
		keys = list[0:size]
	}
	return
}

func (this_ *RedisPoolService) ValueType(ctx context.Context, database int, key string) (ValueType string, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmdType := client.Type(ctx, key)
	ValueType, err = cmdType.Result()
	return
}

func (this_ *RedisPoolService) Get(ctx context.Context, database int, key string, valueStart, valueSize int64) (valueInfo RedisValueInfo, err error) {
	var ValueType string
	ValueType, err = this_.ValueType(ctx, database, key)
	if err != nil {
		return
	}
	var value interface{}

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	if ValueType == "none" {

	} else if ValueType == "string" {
		cmd := client.Get(ctx, key)
		value, err = cmd.Result()
		if err != nil {
			util.Logger.Error("Get Error", zap.Any("key", key), zap.Error(err))
			return
		}
	} else if ValueType == "list" {

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

	} else if ValueType == "set" {

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
	} else if ValueType == "hash" {

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
		println(ValueType)
	}
	valueInfo.ValueType = ValueType
	valueInfo.Value = value

	return
}

func (this_ *RedisPoolService) Set(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.Set(ctx, key, value, time.Duration(0))
	_, err = cmd.Result()
	return
}

func (this_ *RedisPoolService) SAdd(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.SAdd(ctx, key, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisPoolService) SRem(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.SRem(ctx, key, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisPoolService) LPush(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.LPush(ctx, key, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisPoolService) RPush(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.RPush(ctx, key, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisPoolService) LSet(ctx context.Context, database int, key string, index int64, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.LSet(ctx, key, index, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisPoolService) LRem(ctx context.Context, database int, key string, count int64, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.LRem(ctx, key, count, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisPoolService) HSet(ctx context.Context, database int, key string, field string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.HSet(ctx, key, field, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisPoolService) HDel(ctx context.Context, database int, key string, field string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.HDel(ctx, key, field)
	_, err = cmd.Result()
	return
}

func (this_ *RedisPoolService) Del(ctx context.Context, database int, key string) (count int, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	cmd := client.Del(ctx, key)
	_, err = cmd.Result()
	if err == nil {
		count++
	}
	return
}

func (this_ *RedisPoolService) DelPattern(ctx context.Context, database int, pattern string) (count int, err error) {
	count = 0
	var keys []string
	_, keys, err = this_.Keys(ctx, database, pattern, -1)

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	for _, key := range keys {
		cmd := client.Del(ctx, key)
		_, err = cmd.Result()
		if err == nil {
			count++
		} else {
			return
		}
	}
	return
}
