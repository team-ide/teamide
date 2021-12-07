package base

import (
	"config"
	"fmt"
	"server"
)

var (
	REDIS_PREFIX     string = config.Config.Redis.Prefix
	ZOOKEEPER_PREFIX string = config.Config.Zookeeper.Namespace
)

type IDType int8

// ID类型
const (
	ID_TYPE_USER          IDType = 1
	ID_TYPE_USER_METADATA IDType = 2
)

func GetRedisKey() (key string) {
	key = fmt.Sprint(REDIS_PREFIX, ":", "server", ":", server.GetServerId())
	return
}

func GetIDRedisKey(idType IDType) (key string) {
	key = fmt.Sprint(GetRedisKey(), ":", "id", ":", idType)
	return
}

func GetUserInsertLockRedisKey(accountKey string) (key string) {
	key = fmt.Sprint(GetRedisKey(), ":", "user:", accountKey, ":insert:lock")
	return
}

func GetUserMetadataLockRedisKey(userId int64) (key string) {
	key = fmt.Sprint(GetRedisKey(), ":", "user:", userId, ":metadata:lock")
	return
}
