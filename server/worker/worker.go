package worker

import (
	"base"
	"sync"
	"time"
)

type Worker struct {
	Name    string
	Text    string
	Icon    string
	Comment string
	Configs interface{}
	WorkMap map[string]func(interface{}) (interface{}, error)
}

var (
	WorkerCache            map[string]*Worker
	AutomaticShutdownCache map[string]*AutomaticShutdown
	lock                   sync.Mutex
)

func init() {
	WorkerCache = map[string]*Worker{}
	AutomaticShutdownCache = map[string]*AutomaticShutdown{}
	go startToolboxAutomaticShutdownTimer()
}

func AddWorker(worker *Worker) {
	name := worker.Name
	WorkerCache[name] = worker
}

func GetAutomaticShutdown(key string, create func(*AutomaticShutdown) error) (res *AutomaticShutdown, err error) {
	lock.Lock()
	defer lock.Unlock()
	var ok bool
	res, ok = AutomaticShutdownCache[key]
	if !ok {
		res = &AutomaticShutdown{}
		err = create(res)
		if err != nil {
			if res.Stop != nil {
				res.Stop()
			}
			return
		}
		AutomaticShutdownCache[key] = res
	}
	return
}

type AutomaticShutdown struct {
	AutomaticShutdown int64
	LastUseTimestamp  int64
	Stop              func()
	Service           interface{}
}

func startToolboxAutomaticShutdownTimer() {
	time.Sleep(1 * time.Second)
	cache := AutomaticShutdownCache
	nowTime := base.GetNowTime()
	for key, one := range cache {
		if one.AutomaticShutdown <= 0 {
			continue
		}
		lastTime := nowTime - one.LastUseTimestamp
		if lastTime >= one.AutomaticShutdown {
			delete(AutomaticShutdownCache, key)
			(*one).Stop()
		}

	}
}
