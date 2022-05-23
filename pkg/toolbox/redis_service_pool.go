package toolbox

import (
	"context"
	"sort"
	"time"

	redis "github.com/go-redis/redis/v8"
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

func (this_ *RedisPoolService) GetClient(ctx context.Context) *redis.Client {
	defer func() {
		this_.lastUseTime = GetNowTime()
	}()
	return this_.client
}

func (this_ *RedisPoolService) SelectDatabase(ctx context.Context, client *redis.Client, database int) (err error) {
	conn := client.Conn(ctx)
	cmdSelect := conn.Select(context.TODO(), database)
	_, err = cmdSelect.Result()
	return
}

func (this_ *RedisPoolService) Keys(ctx context.Context, database int, cursor uint64, pattern string, size int64) (resCursor uint64, keys []string, err error) {

	client := this_.GetClient(ctx)

	err = this_.SelectDatabase(ctx, client, database)
	if err != nil {
		return
	}

	cmdKeys := client.Scan(ctx, cursor, pattern, size)
	var list []string
	list, resCursor, err = cmdKeys.Result()
	if err != nil {
		return
	}

	sor := sort.StringSlice(list)
	sor.Sort()

	return
}

func (this_ *RedisPoolService) KeyType(ctx context.Context, database int, key string) (keyType string, err error) {

	client := this_.GetClient(ctx)

	err = this_.SelectDatabase(ctx, client, database)
	if err != nil {
		return
	}

	cmdType := client.Type(ctx, key)
	keyType, err = cmdType.Result()
	return
}

func (this_ *RedisPoolService) Get(ctx context.Context, database int, key string, valueStart, valueSize int64, cursor uint64) (valueInfo RedisValueInfo, err error) {
	var keyType string
	keyType, err = this_.KeyType(ctx, database, key)
	if err != nil {
		return
	}
	var value interface{}
	client := this_.GetClient(ctx)
	defer client.Close()

	err = this_.SelectDatabase(client, database)
	if err != nil {
		return
	}

	if keyType == "none" {

	} else if keyType == "string" {
		var reply interface{}
		reply, err = client.Do("get", key)
		if err != nil {
			return
		}
		if reply != nil {
			value, err = redis.String(reply, err)
		}
	} else if keyType == "list" {
		var reply interface{}
		reply, err = client.Do("llen", key)
		if err != nil {
			return
		}
		valueInfo.ValueCount, err = redis.Int64(reply, err)
		if err != nil {
			return
		}

		valueInfo.ValueStart = valueStart
		valueInfo.ValueEnd = valueInfo.ValueStart + valueSize

		reply, err = client.Do("lrange", key, valueInfo.ValueStart, valueInfo.ValueEnd)
		if err != nil {
			return
		}
		value, err = redis.Strings(reply, err)
	} else if keyType == "set" {

		var reply interface{}
		reply, err = client.Do("scard", key)
		if err != nil {
			return
		}
		valueInfo.ValueCount, err = redis.Int64(reply, err)
		if err != nil {
			return
		}

		reply, err = client.Do("sscan", key, cursor, "match", "*", "count", valueSize)
		if err != nil {
			return
		}
		value, err = redis.Strings(reply, err)
	} else if keyType == "hash" {

		var reply interface{}
		reply, err = client.Do("hlen", key)
		if err != nil {
			return
		}
		valueInfo.ValueCount, err = redis.Int64(reply, err)
		if err != nil {
			return
		}

		reply, err = client.Do("hscan", key, cursor, "match", "*", "count", valueSize)
		if err != nil {
			return
		}
		value, err = redis.StringMap(reply, err)
	} else {
		println(keyType)
	}
	valueInfo.Type = keyType
	valueInfo.Value = value

	return
}

func (this_ *RedisPoolService) Set(ctx context.Context, database int, key string, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()

	err = this_.SelectDatabase(client, database)
	if err != nil {
		return
	}

	_, err = client.Do("set", key, value)
	return
}

func (this_ *RedisPoolService) Sadd(ctx context.Context, database int, key string, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()

	err = this_.SelectDatabase(client, database)
	if err != nil {
		return
	}

	_, err = client.Do("sadd", key, value)
	return
}

func (this_ *RedisPoolService) Srem(ctx context.Context, database int, key string, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()

	err = this_.SelectDatabase(client, database)
	if err != nil {
		return
	}

	_, err = client.Do("srem", key, value)
	return
}

func (this_ *RedisPoolService) Lpush(ctx context.Context, database int, key string, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()

	err = this_.SelectDatabase(client, database)
	if err != nil {
		return
	}

	_, err = client.Do("lpush", key, value)
	return
}

func (this_ *RedisPoolService) Rpush(ctx context.Context, database int, key string, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()

	err = this_.SelectDatabase(client, database)
	if err != nil {
		return
	}

	_, err = client.Do("rpush", key, value)
	return
}

func (this_ *RedisPoolService) Lset(ctx context.Context, database int, key string, index int64, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()

	err = this_.SelectDatabase(client, database)
	if err != nil {
		return
	}

	_, err = client.Do("lset", key, index, value)
	return
}

func (this_ *RedisPoolService) Lrem(ctx context.Context, database int, key string, count int64, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()

	err = this_.SelectDatabase(client, database)
	if err != nil {
		return
	}

	_, err = client.Do("lrem", key, count, value)
	return
}

func (this_ *RedisPoolService) Hset(ctx context.Context, database int, key string, field string, value string) (err error) {

	client := this_.GetClient()
	defer client.Close()

	err = this_.SelectDatabase(client, database)
	if err != nil {
		return
	}

	_, err = client.Do("hset", key, field, value)
	return
}

func (this_ *RedisPoolService) Hdel(ctx context.Context, database int, key string, field string) (err error) {

	client := this_.GetClient()
	defer client.Close()

	err = this_.SelectDatabase(client, database)
	if err != nil {
		return
	}

	_, err = client.Do("hdel", key, field)
	return
}

func (this_ *RedisPoolService) Del(ctx context.Context, database int, key string) (count int, err error) {
	count = 0
	client := this_.GetClient()
	defer client.Close()

	err = this_.SelectDatabase(client, database)
	if err != nil {
		return
	}

	_, err = client.Do("del", key)
	if err == nil {
		count++
	}
	return
}

func (this_ *RedisPoolService) DelPattern(ctx context.Context, database int, pattern string) (count int, err error) {
	count = 0
	var list []string
	_, list, err = this_.Keys(database, pattern, 0)

	client := this_.GetClient()
	defer client.Close()

	err = this_.SelectDatabase(client, database)
	if err != nil {
		return
	}

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
