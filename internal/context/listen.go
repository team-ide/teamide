package context

import (
	"sync"
	"teamide/pkg/util"
	"time"
)

var (
	listenerCache                []*ClientTabListener
	listenerCacheForUserId       = map[int64][]*ClientTabListener{}
	listenerCacheForClientKey    = map[string][]*ClientTabListener{}
	listenerCacheForClientTabKey = map[string]*ClientTabListener{}
	listenerCacheLock            = &sync.Mutex{}
)

func GetListeners() (listeners []*ClientTabListener) {
	listenerCacheLock.Lock()
	defer listenerCacheLock.Unlock()
	listeners = listenerCache
	return
}

func GetListener(clientTabKey string) (listener *ClientTabListener) {
	listenerCacheLock.Lock()
	defer listenerCacheLock.Unlock()
	listener = listenerCacheForClientTabKey[clientTabKey]
	return
}

func GetClientListeners(clientKey string) (listeners []*ClientTabListener) {
	listenerCacheLock.Lock()
	defer listenerCacheLock.Unlock()
	listeners = listenerCacheForClientKey[clientKey]
	return
}

func GetUserListeners(userId int64) (listeners []*ClientTabListener) {
	listenerCacheLock.Lock()
	defer listenerCacheLock.Unlock()
	listeners = listenerCacheForUserId[userId]
	return
}

func ChangeListenerUserId(listener *ClientTabListener, userId int64) {
	if listener.UserId == userId {
		return
	}
	listenerCacheLock.Lock()
	defer listenerCacheLock.Unlock()

	listener.UserId = userId
	if listener.UserId > 0 {
		var newList []*ClientTabListener
		list := listenerCacheForUserId[listener.UserId]
		for _, one := range list {
			if one.ClientTabKey != listener.ClientTabKey {
				newList = append(newList, one)
			}
		}
		listenerCacheForUserId[listener.UserId] = newList
	}
	if userId > 0 {
		var newList []*ClientTabListener
		list := listenerCacheForUserId[userId]
		var find bool
		for _, one := range list {
			if one.ClientTabKey != listener.ClientTabKey {
				newList = append(newList, one)
			} else {
				find = true
			}
		}
		if !find {
			newList = append(newList, listener)
		}
		listenerCacheForUserId[userId] = newList
	}

	listener.AddEvent(&ListenEvent{})
}

func ChangeListenerClientKey(listener *ClientTabListener, clientKey string) {
	if listener.ClientKey == clientKey {
		return
	}
	listenerCacheLock.Lock()
	defer listenerCacheLock.Unlock()

	listener.ClientKey = clientKey
	if listener.ClientKey != "" {
		var newList []*ClientTabListener
		list := listenerCacheForClientKey[listener.ClientKey]
		for _, one := range list {
			if one.ClientTabKey != listener.ClientTabKey {
				newList = append(newList, one)
			}
		}
		listenerCacheForClientKey[listener.ClientKey] = newList
	}
	if clientKey != "" {
		var newList []*ClientTabListener
		list := listenerCacheForClientKey[listener.ClientKey]
		var find bool
		for _, one := range list {
			if one.ClientTabKey != listener.ClientTabKey {
				newList = append(newList, one)
			} else {
				find = true
			}
		}
		if !find {
			newList = append(newList, listener)
		}
		listenerCacheForClientKey[listener.ClientKey] = newList
	}
	listener.AddEvent(&ListenEvent{})
}

func AddListener(listener *ClientTabListener) {
	if listener == nil {
		return
	}
	listenerCacheLock.Lock()
	defer listenerCacheLock.Unlock()

	var newList []*ClientTabListener
	list := listenerCache
	var find bool
	for _, one := range list {
		if one.ClientTabKey != listener.ClientTabKey {
			newList = append(newList, one)
		} else {
			find = true
		}
	}
	if !find {
		newList = append(newList, listener)
	}
	listenerCache = newList

	if listener.ClientTabKey != "" {
		listenerCacheForClientTabKey[listener.ClientTabKey] = listener
	}
	if listener.ClientKey != "" {
		var newList []*ClientTabListener
		list := listenerCacheForClientKey[listener.ClientKey]
		var find bool
		for _, one := range list {
			if one.ClientTabKey != listener.ClientTabKey {
				newList = append(newList, one)
			} else {
				find = true
			}
		}
		if !find {
			newList = append(newList, listener)
		}
		listenerCacheForClientKey[listener.ClientTabKey] = newList
	}
	if listener.UserId != 0 {
		var newList []*ClientTabListener
		list := listenerCacheForUserId[listener.UserId]
		var find bool
		for _, one := range list {
			if one.ClientTabKey != listener.ClientTabKey {
				newList = append(newList, one)
			} else {
				find = true
			}
		}
		if !find {
			newList = append(newList, listener)
		}
		listenerCacheForUserId[listener.UserId] = newList
	}

	return
}

func RemoveListener(listener *ClientTabListener) {
	if listener == nil {
		return
	}
	listenerCacheLock.Lock()
	defer listenerCacheLock.Unlock()

	listener.AddEvent(&ListenEvent{})

	var newList []*ClientTabListener
	list := listenerCache
	for _, one := range list {
		if one.ClientTabKey != listener.ClientTabKey {
			newList = append(newList, one)
		}
	}
	listenerCache = newList

	if listener.ClientTabKey != "" {
		delete(listenerCacheForClientTabKey, listener.ClientTabKey)
	}
	if listener.ClientKey != "" {
		var newList []*ClientTabListener
		list := listenerCacheForClientKey[listener.ClientKey]
		for _, one := range list {
			if one.ClientTabKey != listener.ClientTabKey {
				newList = append(newList, one)
			}
		}
		listenerCacheForClientKey[listener.ClientTabKey] = newList
	}
	if listener.UserId != 0 {
		var newList []*ClientTabListener
		list := listenerCacheForUserId[listener.UserId]
		for _, one := range list {
			if one.ClientTabKey != listener.ClientTabKey {
				newList = append(newList, one)
			}
		}
		listenerCacheForUserId[listener.UserId] = newList
	}

	return
}

func CallUserEvent(userId int64, event *ListenEvent) {
	list := GetUserListeners(userId)
	if len(list) == 0 {
		return
	}
	for _, one := range list {
		go one.AddEvent(event)
	}
}

func CallClientKeyEvent(clientKey string, event *ListenEvent) {
	list := GetClientListeners(clientKey)
	if len(list) == 0 {
		return
	}
	for _, one := range list {
		go one.AddEvent(event)
	}
}

func CallClientTabKeyEvent(clientTabKey string, event *ListenEvent) {
	listener := GetListener(clientTabKey)
	if listener == nil {
		return
	}
	go listener.AddEvent(event)
}

func listenerInit() {

	go checkListener()
}

func checkListener() {

	defer func() {
		time.Sleep(time.Second * 60)
		checkListener()
	}()

	list := GetListeners()
	for _, one := range list {
		nowTime := util.GetNowTime()

		// 最后监听时间 在此 之前的 都为超时
		outTime := nowTime - int64(listenerLastListenTimeout)
		lastListenTime := util.GetTimeTime(one.lastListenTime)
		if lastListenTime > outTime {
			continue
		}
		RemoveListener(one)
	}

}

type ListenEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func NewListenEvent(event string, data interface{}) *ListenEvent {
	return &ListenEvent{
		Event: event,
		Data:  data,
	}
}

func NewClientTabListener(clientKey string, clientTabKey string, userId int64) (listener *ClientTabListener) {

	listener = &ClientTabListener{
		ClientTabKey: clientTabKey,
		ClientKey:    clientKey,
		UserId:       userId,
		events:       make(chan *ListenEvent),
		lock:         &sync.Mutex{},
	}
	listener.ListenTimeout = defaultListenTimeout
	listener.lastListenTime = time.Now()
	return
}

// 客户端 监听 数据 默认超时时间 10 分钟
var defaultListenTimeout = 10 * 60 * time.Second

// 监听程序 最后监听时间 超时时间 超过这个时间 未监听 则移除监听器
var listenerLastListenTimeout = 5 * 60 * time.Second

type ClientTabListener struct {
	ClientKey      string `json:"clientKey"`
	ClientTabKey   string `json:"clientTabKey"`
	UserId         int64  `json:"userId"`
	lastListenTime time.Time
	events         chan *ListenEvent
	lock           sync.Locker
	ListenTimeout  time.Duration `json:"listenTimeout"`
	listenIng      bool
}

func (this_ *ClientTabListener) AddEvent(event *ListenEvent) {
	this_.events <- event
}

func (this_ *ClientTabListener) Listen() []*ListenEvent {
	this_.lock.Lock()
	defer this_.lock.Unlock()
	this_.listenIng = true
	this_.lastListenTime = time.Now()
	defer func() {
		this_.listenIng = false
		this_.lastListenTime = time.Now()
	}()
	var timeout = this_.ListenTimeout
	if timeout <= 0 {
		timeout = defaultListenTimeout
	}

	var events []*ListenEvent
	quit := make(chan bool)
	//新开一个协程
	go func() {
		for {
			select {
			case event := <-this_.events: //如果有数据，下面打印。但是有可能ch一直没数据
				events = append(events, event)
				if len(this_.events) == 0 {
					break
				}
			case <-time.After(timeout): //上面的ch如果一直没数据会阻塞，那么select也会检测其他case条件，检测到后 x 秒后超时
				quit <- true //写入
			}
		}
	}()
	<-quit //这里暂时阻塞，直到可读

	return events
}
