package util

import (
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	serviceCache  = map[string]Service{}
	lockCacheLock sync.Mutex
	lockCache     = map[string]sync.Mutex{}
)

func init() {
	go startServiceTimer()
}
func getLockCache(key string) sync.Mutex {
	lockCacheLock.Lock()
	defer lockCacheLock.Unlock()
	res, ok := lockCache[key]
	if !ok {
		lockCache[key] = sync.Mutex{}
	}
	return res
}
func removeLockCache(key string) {
	lockCacheLock.Lock()
	defer lockCacheLock.Unlock()
	delete(lockCache, key)
}

func GetService(key string, create func() (Service, error)) (Service, error) {
	lock := getLockCache(key)
	lock.Lock()
	defer removeLockCache(key)
	defer lock.Unlock()

	res, ok := serviceCache[key]
	if ok {
		return res, nil
	}
	Logger.Info("缓存暂无该服务，创建服务", zap.Any("Key", key))
	res, err := create()
	if err != nil {
		if res != nil {
			res.Stop()
		}
		return nil, err
	}
	serviceCache[key] = res
	return res, err
}

type Service interface {
	GetWaitTime() int64
	GetLastUseTime() int64
	Stop()
}

func startServiceTimer() {
	for {
		time.Sleep(1 * time.Second)
		nowTime := GetNowTime()
		for key, one := range serviceCache {
			if one.GetWaitTime() <= 0 {
				continue
			}
			t := nowTime - one.GetLastUseTime()
			if t >= one.GetWaitTime() {
				delete(serviceCache, key)
				one.Stop()
				Logger.Info("缓存服务回收", zap.Any("Key", key), zap.Any("WaitTime", one.GetWaitTime()), zap.Any("NowTime", nowTime), zap.Any("LastUseTime", one.GetLastUseTime()))
			}
		}
	}
}
