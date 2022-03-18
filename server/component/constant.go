package component

import (
	"fmt"
	"teamide/server/config"
)

var (
	REDIS_PREFIX     string
	ZOOKEEPER_PREFIX string
)

func init() {
	if config.Config.Redis != nil {
		REDIS_PREFIX = config.Config.Redis.Prefix
	}
	if config.Config.Zookeeper != nil {
		ZOOKEEPER_PREFIX = config.Config.Zookeeper.Namespace
	}
}

const (
	HTTP_AES_KEY string = "Q56hFAauWk18Gy2i"
)

type IDType int8

// ID类型
const (
	ID_TYPE_USER          IDType = 1
	ID_TYPE_USER_METADATA IDType = 2
)

func GetRedisKey() (key string) {
	key = fmt.Sprint(REDIS_PREFIX)
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

func GeAccessTokenKey(accessToken string) (key string) {
	key = fmt.Sprint(GetRedisKey(), ":", "access:token:", accessToken)
	return
}
