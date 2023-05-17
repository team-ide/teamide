package base

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	serviceCache     = map[string]*ServiceInfo{}
	serviceCacheLock = &sync.Mutex{}
	lockCacheLock    = &sync.Mutex{}
	lockCache        = map[string]sync.Locker{}
)

func init() {
	go startServiceTimer()
}
func getLockCache(key string) sync.Locker {
	lockCacheLock.Lock()
	defer lockCacheLock.Unlock()
	res, ok := lockCache[key]
	if !ok {
		res = &sync.Mutex{}
		lockCache[key] = res
	}
	return res
}
func removeLockCache(key string) {
	lockCacheLock.Lock()
	defer lockCacheLock.Unlock()
	delete(lockCache, key)
}

func GetService(key string, create func() (serviceInfo *ServiceInfo, err error)) (*ServiceInfo, error) {
	res, ok := FindService(key)
	if ok {
		return res, nil
	}
	res, err := createService(key, create)
	return res, err
}

func createService(key string, create func() (serviceInfo *ServiceInfo, err error)) (*ServiceInfo, error) {
	lock := getLockCache(key)
	lock.Lock()
	defer removeLockCache(key)
	defer lock.Unlock()

	res, ok := FindService(key)
	if ok {
		return res, nil
	}
	util.Logger.Info("缓存暂无该服务，创建服务", zap.Any("Key", key))
	res, err := create()
	if err != nil {
		if res != nil {
			res.Stop()
		}
		return nil, err
	}
	setService(key, res)
	return res, err
}

func FindService(key string) (*ServiceInfo, bool) {
	serviceCacheLock.Lock()
	defer serviceCacheLock.Unlock()

	res, ok := serviceCache[key]
	if ok {
		return res, true
	}
	return nil, false
}

func setService(key string, ser *ServiceInfo) {
	serviceCacheLock.Lock()
	defer serviceCacheLock.Unlock()

	serviceCache[key] = ser
}

type ServiceInfo struct {
	WaitTime    int64
	LastUseTime int64
	Service     interface{}
	Stop        func()
}

func (this_ *ServiceInfo) SetLastUseTime() {
	this_.LastUseTime = util.GetNowMilli()
}

func startServiceTimer() {
	for {
		time.Sleep(1 * time.Minute)
		cleanCache()
	}
}

func cleanCache() {
	serviceCacheLock.Lock()
	defer serviceCacheLock.Unlock()
	nowTime := util.GetNowMilli()
	for key, one := range serviceCache {
		if one.WaitTime <= 0 {
			continue
		}
		t := nowTime - one.LastUseTime
		if t >= one.WaitTime {
			delete(serviceCache, key)
			if one.Stop != nil {
				go one.Stop()
			}
			util.Logger.Info("缓存服务回收", zap.Any("Key", key), zap.Any("WaitTime", one.WaitTime), zap.Any("NowTime", nowTime), zap.Any("LastUseTime", one.LastUseTime))
		}
	}
}
