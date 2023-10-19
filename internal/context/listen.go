package context

import (
	"github.com/team-ide/go-tool/util"
	"sync"
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

	go listener.AddEvent(&ListenEvent{})
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
	go listener.AddEvent(&ListenEvent{})
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

	go listener.AddEvent(&ListenEvent{})

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
	startTimeSecond := util.GetNowSecond()

	defer func() {
		endTimeSecond := util.GetNowSecond()
		useSecond := endTimeSecond - startTimeSecond
		waitSecond := 60 - useSecond
		if waitSecond > 0 {
			time.Sleep(time.Second * time.Duration(waitSecond))
		}
		checkListener()
	}()

	list := GetListeners()
	for _, one := range list {
		nowTimeSecond := util.GetNowSecond()

		// 最后监听时间 在此 之前的 都为超时
		outTimeSecond := nowTimeSecond - listenerLastListenTimeoutSecond
		lastListenTimeSecond := util.GetSecondByTime(one.lastListenTime)
		if lastListenTimeSecond > outTimeSecond {
			if one.listenIng {
				go one.AddEvent(&ListenEvent{})
			}
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
		events:       make(chan *ListenEvent, 100),
		listenLock:   &sync.Mutex{},
		eventsLock:   &sync.Mutex{},
	}
	listener.lastListenTime = time.Now()
	return
}

// 监听程序 最后监听时间 超时时间 超过这个时间 未监听 则移除监听器
var listenerLastListenTimeoutSecond int64 = 10 * 60

type ClientTabListener struct {
	ClientKey      string `json:"clientKey"`
	ClientTabKey   string `json:"clientTabKey"`
	UserId         int64  `json:"userId"`
	lastListenTime time.Time
	events         chan *ListenEvent
	listenLock     sync.Locker
	eventsLock     sync.Locker
	listenIng      bool
}

func (this_ *ClientTabListener) AddEvent(event *ListenEvent) {
	if event == nil {
		return
	}
	this_.eventsLock.Lock()
	defer this_.eventsLock.Unlock()

	this_.events <- event
}

func (this_ *ClientTabListener) Listen() (events []*ListenEvent) {
	this_.listenLock.Lock()
	defer this_.listenLock.Unlock()
	//fmt.Println("Listen start")

	this_.listenIng = true
	this_.lastListenTime = time.Now()
	defer func() {
		//fmt.Println("Listen end")
		this_.listenIng = false
		this_.lastListenTime = time.Now()
	}()

	event := <-this_.events
	if event != nil && event.Event != "" {
		events = append(events, event)
	}

	return
}
