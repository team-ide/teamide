package toolbox

import (
	"context"
	"sort"
	"time"

	redis "github.com/go-redis/redis/v8"
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
	this_.redisCluster.Close()
}

func (this_ *RedisClusterService) GetClient() {
	defer func() {
		this_.lastUseTime = GetNowTime()
	}()

}

func (this_ *RedisClusterService) Keys(ctx context.Context, database int, pattern string, size int) (count int, keys []string, err error) {
	this_.GetClient()
	var list []string
	this_.redisCluster.ForEachMaster(context.TODO(), func(ctx context.Context, client *redis.Client) (err error) {
		cmd := client.Keys(ctx, pattern)
		var ls []string
		ls, err = cmd.Result()
		if err != nil {
			return
		}
		list = append(list, ls...)
		return
	})
	sor := sort.StringSlice(list)
	sor.Sort()
	count = len(list)
	if count <= size {
		keys = list
	} else {
		for index, one := range list {
			if index < size {
				keys = append(keys, one)
			} else {
				break
			}
		}
	}
	return
}

func (this_ *RedisClusterService) KeyType(ctx context.Context, database int, key string) (keyType string, err error) {
	this_.GetClient()
	cmd := this_.redisCluster.Type(context.TODO(), key)
	keyType, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) Get(ctx context.Context, database int, key string, valueStart, valueSize int64, cursor uint64) (valueInfo RedisValueInfo, err error) {
	this_.GetClient()
	var keyType string
	keyType, err = this_.KeyType(ctx, database, key)
	if err != nil {
		return
	}
	var value interface{}

	if keyType == "none" {

	} else if keyType == "string" {
		cmd := this_.redisCluster.Get(context.TODO(), key)
		value, err = cmd.Result()
	} else if keyType == "list" {

		cmd := this_.redisCluster.LLen(context.TODO(), key)

		valueInfo.ValueCount, err = cmd.Result()
		if err != nil {
			return
		}
		valueInfo.ValueStart = valueStart
		valueInfo.ValueEnd = valueInfo.ValueStart + valueSize

		cmdRange := this_.redisCluster.LRange(context.TODO(), key, valueInfo.ValueStart, valueInfo.ValueEnd)
		value, err = cmdRange.Result()
	} else if keyType == "set" {

		cmdSCard := this_.redisCluster.SCard(context.TODO(), key)
		valueInfo.ValueCount, err = cmdSCard.Result()

		cmd := this_.redisCluster.SScan(context.TODO(), key, cursor, "*", valueSize)
		value, valueInfo.Cursor, err = cmd.Result()
	} else if keyType == "hash" {

		cmdHLen := this_.redisCluster.HLen(context.TODO(), key)
		valueInfo.ValueCount, err = cmdHLen.Result()

		cmd := this_.redisCluster.HScan(context.TODO(), key, cursor, "*", valueSize)
		value, valueInfo.Cursor, err = cmd.Result()
	} else {
		println(keyType)
	}
	valueInfo.Type = keyType
	valueInfo.Value = value
	return
}

func (this_ *RedisClusterService) Set(ctx context.Context, database int, key string, value string) (err error) {
	this_.GetClient()
	cmd := this_.redisCluster.Set(context.TODO(), key, value, time.Duration(0))
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) Sadd(ctx context.Context, database int, key string, value string) (err error) {
	this_.GetClient()
	cmd := this_.redisCluster.SAdd(context.TODO(), key, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) Srem(ctx context.Context, database int, key string, value string) (err error) {
	this_.GetClient()
	cmd := this_.redisCluster.SRem(context.TODO(), key, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) Lpush(ctx context.Context, database int, key string, value string) (err error) {
	this_.GetClient()
	cmd := this_.redisCluster.LPush(context.TODO(), key, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) Rpush(ctx context.Context, database int, key string, value string) (err error) {
	this_.GetClient()
	cmd := this_.redisCluster.RPush(context.TODO(), key, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) Lset(ctx context.Context, database int, key string, index int64, value string) (err error) {
	this_.GetClient()
	cmd := this_.redisCluster.LSet(context.TODO(), key, index, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) Lrem(ctx context.Context, database int, key string, count int64, value string) (err error) {
	this_.GetClient()
	cmd := this_.redisCluster.LRem(context.TODO(), key, count, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) Hset(ctx context.Context, database int, key string, field string, value string) (err error) {
	this_.GetClient()
	cmd := this_.redisCluster.HSet(context.TODO(), key, field, value)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) Hdel(ctx context.Context, database int, key string, field string) (err error) {
	this_.GetClient()
	cmd := this_.redisCluster.HDel(context.TODO(), key, field)
	_, err = cmd.Result()
	return
}

func (this_ *RedisClusterService) Del(ctx context.Context, database int, key string) (count int, err error) {
	this_.GetClient()
	count = 0
	cmd := this_.redisCluster.Del(context.TODO(), key)
	_, err = cmd.Result()
	if err == nil {
		count++
	}
	return
}

func (this_ *RedisClusterService) DelPattern(ctx context.Context, database int, pattern string) (count int, err error) {
	this_.GetClient()
	count = 0
	cmd := this_.redisCluster.Keys(context.TODO(), pattern)
	var list []string
	list, err = cmd.Result()
	if err != nil {
		return
	}
	for _, key := range list {
		var num int
		num, err = this_.Del(ctx, database, key)
		if err != nil {
			return
		}
		count += num
	}
	return
}
