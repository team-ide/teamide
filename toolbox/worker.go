package toolbox

import (
	"sync"
	"time"
)

type Worker struct {
	Name    string
	Text    string
	Icon    string
	Comment string
	WorkMap map[string]func(map[string]interface{}) (map[string]interface{}, error)
}

var (
	WorkerCache  map[string]*Worker
	serviceCache map[string]Service
	lock         sync.Mutex
)

func init() {
	WorkerCache = map[string]*Worker{}
	serviceCache = map[string]Service{}
	go startServiceTimer()
}

func AddWorker(worker *Worker) {
	name := worker.Name
	WorkerCache[name] = worker
}

func GetService(key string, create func() (Service, error)) (Service, error) {
	lock.Lock()
	defer lock.Unlock()
	res, ok := serviceCache[key]
	if ok {
		return res, nil
	}
	res, err := create()
	if err != nil {
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
			if one.GetWaitTime() > 0 && nowTime-one.GetLastUseTime() >= one.GetWaitTime() {
				delete(serviceCache, key)
				one.Stop()
			}
		}
	}
}
