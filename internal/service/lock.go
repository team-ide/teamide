package service

import "sync"

// TMLock 自定义锁
type TMLock interface {
	Lock()
	Unlock()
}

var (
	LockCache = make(map[string]TMLock)
	cacheLock = &sync.RWMutex{}
)

// GetLock 根据Key 获取锁
func GetLock(key string) TMLock {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	lock, find := LockCache[key]
	if !find {
		lock = &sync.RWMutex{}
		LockCache[key] = lock
	}
	return lock
}
