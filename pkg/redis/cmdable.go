package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"sort"
	"teamide/pkg/util"
	"time"
)

func Keys(ctx context.Context, client redis.Cmdable, database int, pattern string, size int64) (keysResult *KeysResult, err error) {
	keysResult = &KeysResult{}
	var list []string
	cmdKeys := client.Keys(ctx, pattern)
	list, err = cmdKeys.Result()
	if err != nil {
		return
	}
	keysResult.Count = len(list)

	sor := sort.StringSlice(list)
	sor.Sort()

	var keys []string
	if int64(keysResult.Count) <= size || size < 0 {
		keys = list
	} else {
		keys = list[0:size]
	}
	for _, key := range keys {
		info := &KeyInfo{
			Key:      key,
			Database: database,
		}
		keysResult.KeyList = append(keysResult.KeyList, info)
	}
	return
}

func Exists(ctx context.Context, client redis.Cmdable, key string) (res int64, err error) {

	cmdExists := client.Exists(ctx, key)
	res, err = cmdExists.Result()
	return
}

func ValueType(ctx context.Context, client redis.Cmdable, key string) (ValueType string, err error) {

	cmdType := client.Type(ctx, key)
	ValueType, err = cmdType.Result()
	if err == redis.Nil {
		err = nil
		return
	}
	return
}

// Expire 让给定键在指定的秒数之后过期
func Expire(ctx context.Context, client redis.Cmdable, key string, expire int64) (res bool, err error) {
	cmd := client.Expire(ctx, key, time.Duration(expire)*time.Second)
	res, err = cmd.Result()
	if err == redis.Nil {
		err = nil
		return
	}
	return
}

// Persist 移除键的过期时间
func Persist(ctx context.Context, client redis.Cmdable, key string) (res bool, err error) {
	cmd := client.Persist(ctx, key)
	res, err = cmd.Result()
	if err == redis.Nil {
		err = nil
		return
	}
	return
}

// TTL 查看给定键距离过期还有多少秒
func TTL(ctx context.Context, client redis.Cmdable, key string) (res int64, err error) {
	cmd := client.TTL(ctx, key)
	r, err := cmd.Result()
	if err == redis.Nil {
		err = nil
		return
	}
	if err != nil {
		return
	}
	if r > 0 {
		res = int64(r / time.Second)
	}
	return
}

func Get(ctx context.Context, client redis.Cmdable, database int, key string, valueStart, valueSize int64) (valueInfo *ValueInfo, err error) {
	var valueType string
	valueType, err = ValueType(ctx, client, key)
	if err != nil {
		return
	}
	valueInfo = &ValueInfo{
		Key:      key,
		Database: database,
	}

	valueInfo.MemoryUsage, _ = MemoryUsage(ctx, client, key)
	valueInfo.TTL, _ = TTL(ctx, client, key)
	var value interface{}

	if valueType == "none" {

	} else if valueType == "string" {
		cmd := client.Get(ctx, key)
		value, err = cmd.Result()
		if err != nil {
			util.Logger.Error("Get Error", zap.Any("key", key), zap.Error(err))
			return
		}
	} else if valueType == "list" {

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

	} else if valueType == "set" {

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
	} else if valueType == "hash" {

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
		util.Logger.Warn("valueType not support", zap.Any("valueType", valueType), zap.Any("key", key))
	}
	valueInfo.ValueType = valueType
	valueInfo.Value = value

	return
}

func Set(ctx context.Context, client redis.Cmdable, key string, value string) (err error) {

	cmd := client.Set(ctx, key, value, time.Duration(0))
	_, err = cmd.Result()
	return
}

func SAdd(ctx context.Context, client redis.Cmdable, key string, value string) (err error) {

	cmd := client.SAdd(ctx, key, value)
	_, err = cmd.Result()
	return
}

func SRem(ctx context.Context, client redis.Cmdable, key string, value string) (err error) {

	cmd := client.SRem(ctx, key, value)
	_, err = cmd.Result()
	return
}

func LPush(ctx context.Context, client redis.Cmdable, key string, value string) (err error) {

	cmd := client.LPush(ctx, key, value)
	_, err = cmd.Result()
	return
}

func RPush(ctx context.Context, client redis.Cmdable, key string, value string) (err error) {

	cmd := client.RPush(ctx, key, value)
	_, err = cmd.Result()
	return
}

func LSet(ctx context.Context, client redis.Cmdable, key string, index int64, value string) (err error) {

	cmd := client.LSet(ctx, key, index, value)
	_, err = cmd.Result()
	return
}

func LRem(ctx context.Context, client redis.Cmdable, key string, count int64, value string) (err error) {

	cmd := client.LRem(ctx, key, count, value)
	_, err = cmd.Result()
	return
}

func HSet(ctx context.Context, client redis.Cmdable, key string, field string, value string) (err error) {

	cmd := client.HSet(ctx, key, field, value)
	_, err = cmd.Result()
	return
}

func HDel(ctx context.Context, client redis.Cmdable, key string, field string) (err error) {

	cmd := client.HDel(ctx, key, field)
	_, err = cmd.Result()
	return
}

func HGet(ctx context.Context, client redis.Cmdable, key string, field string) (value string, err error) {

	cmd := client.HGet(ctx, key, field)
	value, err = cmd.Result()
	if err == redis.Nil {
		err = nil
		return
	}
	return
}

func SetBit(ctx context.Context, client redis.Cmdable, key string, offset int64, value int) (err error) {

	cmd := client.SetBit(ctx, key, offset, value)
	err = cmd.Err()
	return
}

func BitCount(ctx context.Context, client redis.Cmdable, key string) (count int64, err error) {
	cmd := client.BitCount(ctx, key, nil)
	count, err = cmd.Result()
	return
}

func Info(ctx context.Context, client redis.Cmdable) (res string, err error) {
	cmd := client.Info(ctx)
	res, err = cmd.Result()
	return
}

func MemoryUsage(ctx context.Context, client redis.Cmdable, key string) (size int64, err error) {
	cmd := client.MemoryUsage(ctx, key)
	size, err = cmd.Result()
	return
}

func Del(ctx context.Context, client redis.Cmdable, key string) (count int, err error) {

	cmd := client.Del(ctx, key)
	_, err = cmd.Result()
	if err == nil {
		count++
	}
	return
}

func DelPattern(ctx context.Context, client redis.Cmdable, database int, pattern string) (count int, err error) {
	count = 0
	keysResult, err := Keys(ctx, client, database, pattern, -1)
	if err != nil {
		return
	}
	for _, keyInfo := range keysResult.KeyList {
		cmd := client.Del(ctx, keyInfo.Key)
		_, err = cmd.Result()
		if err == nil {
			count++
		} else {
			return
		}
	}
	return
}
