package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strings"
)

type Config struct {
	Address string `json:"address"`
	Auth    string `json:"auth"`
}

func CreateService(config Config) (service Service, err error) {
	if !strings.Contains(config.Address, ",") && !strings.Contains(config.Address, ";") {
		service, err = CreateRedisService(config.Address, config.Auth)
	} else {
		var servers []string
		if strings.Contains(config.Address, ",") {
			servers = strings.Split(config.Address, ",")
		} else if strings.Contains(config.Address, ";") {
			servers = strings.Split(config.Address, ";")
		} else {
			servers = []string{config.Address}
		}
		service, err = CreateClusterService(servers, config.Auth)
	}
	return
}

type Service interface {
	GetWaitTime() int64
	GetLastUseTime() int64
	Stop()
	GetClient(ctx context.Context, database int) (client redis.Cmdable, err error)
	Exists(ctx context.Context, database int, key string) (res int64, err error)
	Keys(ctx context.Context, database int, pattern string, size int64) (count int, keys []string, err error)
	Get(ctx context.Context, database int, key string, valueStart, valueSize int64) (valueInfo *ValueInfo, err error)
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
