package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sort"
	"teamide/pkg/util"
	"time"
)

func CreateClusterService(servers []string, auth string) (service *ClusterService, err error) {
	service = &ClusterService{
		servers: servers,
		auth:    auth,
	}
	err = service.init()
	return
}

type ClusterService struct {
	servers      []string
	auth         string
	redisCluster *redis.ClusterClient
	lastUseTime  int64
}

func (this_ *ClusterService) init() (err error) {
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

func (this_ *ClusterService) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *ClusterService) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *ClusterService) Stop() {
	_ = this_.redisCluster.Close()
}

func (this_ *ClusterService) GetClient(ctx context.Context, database int) (client redis.Cmdable, err error) {
	defer func() {
		this_.lastUseTime = util.GetNowTime()
	}()
	client = this_.redisCluster
	if ctx != nil && database >= 0 {
		return
	}
	return
}

func (this_ *ClusterService) Info(ctx context.Context) (res string, err error) {

	client, err := this_.GetClient(ctx, 0)
	if err != nil {
		return
	}

	return Info(ctx, client)
}

func (this_ *ClusterService) Keys(ctx context.Context, database int, pattern string, size int64) (count int, keys []string, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return ClusterKeys(ctx, client.(*redis.ClusterClient), pattern, size)
}

func (this_ *ClusterService) Exists(ctx context.Context, database int, key string) (res int64, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return Exists(ctx, client, key)
}

func (this_ *ClusterService) Get(ctx context.Context, database int, key string, valueStart, valueSize int64) (valueInfo *ValueInfo, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return Get(ctx, client, key, valueStart, valueSize)
}

func (this_ *ClusterService) Set(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return Set(ctx, client, key, value)
}

func (this_ *ClusterService) SAdd(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return SAdd(ctx, client, key, value)
}

func (this_ *ClusterService) SRem(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return SRem(ctx, client, key, value)
}

func (this_ *ClusterService) LPush(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return LPush(ctx, client, key, value)
}

func (this_ *ClusterService) RPush(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return RPush(ctx, client, key, value)
}

func (this_ *ClusterService) LSet(ctx context.Context, database int, key string, index int64, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return LSet(ctx, client, key, index, value)
}

func (this_ *ClusterService) LRem(ctx context.Context, database int, key string, count int64, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return LRem(ctx, client, key, count, value)
}

func (this_ *ClusterService) HSet(ctx context.Context, database int, key string, field string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return HSet(ctx, client, key, field, value)
}

func (this_ *ClusterService) HDel(ctx context.Context, database int, key string, field string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return HDel(ctx, client, key, field)
}

func (this_ *ClusterService) Del(ctx context.Context, database int, key string) (count int, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return Del(ctx, client, key)
}

func (this_ *ClusterService) DelPattern(ctx context.Context, database int, pattern string) (count int, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	var keys []string
	_, keys, err = ClusterKeys(ctx, client.(*redis.ClusterClient), pattern, -1)
	if err != nil {
		return
	}

	count = 0
	for _, key := range keys {
		_, err = Del(ctx, client, key)
		if err == nil {
			count++
		}
	}
	return
}

func ClusterKeys(ctx context.Context, client *redis.ClusterClient, pattern string, size int64) (count int, keys []string, err error) {

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
