package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"teamide/pkg/util"
	"time"
)

func CreateRedisService(address string, auth string) (service *V8Service, err error) {
	service = &V8Service{
		address: address,
		auth:    auth,
	}
	err = service.init()
	return
}

type ValueInfo struct {
	ValueType   string      `json:"valueType"`
	Value       interface{} `json:"value"`
	ValueCount  int64       `json:"valueCount"`
	ValueStart  int64       `json:"valueStart"`
	ValueEnd    int64       `json:"valueEnd"`
	Cursor      uint64      `json:"cursor"`
	MemoryUsage int64       `json:"memoryUsage"`
}

type V8Service struct {
	address     string
	auth        string
	client      *redis.Client
	lastUseTime int64
}

func (this_ *V8Service) init() (err error) {
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

func (this_ *V8Service) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *V8Service) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *V8Service) Stop() {
	_ = this_.client.Close()
}

func (this_ *V8Service) GetClient(ctx context.Context, database int) (client redis.Cmdable, err error) {
	defer func() {
		this_.lastUseTime = util.GetNowTime()
	}()
	cmd := this_.client.Do(ctx, "select", database)
	_, err = cmd.Result()
	if err != nil {
		return
	}
	client = this_.client
	return
}

func (this_ *V8Service) Keys(ctx context.Context, database int, pattern string, size int64) (count int, keys []string, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return Keys(ctx, client, pattern, size)
}

func (this_ *V8Service) Exists(ctx context.Context, database int, key string) (res int64, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return Exists(ctx, client, key)
}

func (this_ *V8Service) Get(ctx context.Context, database int, key string, valueStart, valueSize int64) (valueInfo *ValueInfo, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return Get(ctx, client, key, valueStart, valueSize)
}

func (this_ *V8Service) Set(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return Set(ctx, client, key, value)
}

func (this_ *V8Service) SAdd(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return SAdd(ctx, client, key, value)
}

func (this_ *V8Service) SRem(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return SRem(ctx, client, key, value)
}

func (this_ *V8Service) LPush(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return LPush(ctx, client, key, value)
}

func (this_ *V8Service) RPush(ctx context.Context, database int, key string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return RPush(ctx, client, key, value)
}

func (this_ *V8Service) LSet(ctx context.Context, database int, key string, index int64, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return LSet(ctx, client, key, index, value)
}

func (this_ *V8Service) LRem(ctx context.Context, database int, key string, count int64, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return LRem(ctx, client, key, count, value)
}

func (this_ *V8Service) HSet(ctx context.Context, database int, key string, field string, value string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return HSet(ctx, client, key, field, value)
}

func (this_ *V8Service) HDel(ctx context.Context, database int, key string, field string) (err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return HDel(ctx, client, key, field)
}

func (this_ *V8Service) Del(ctx context.Context, database int, key string) (count int, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return Del(ctx, client, key)
}

func (this_ *V8Service) DelPattern(ctx context.Context, database int, pattern string) (count int, err error) {

	client, err := this_.GetClient(ctx, database)
	if err != nil {
		return
	}

	return DelPattern(ctx, client, pattern)
}
