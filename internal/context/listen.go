package context

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
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
		one.AddEvent(event)
	}
}

func CallClientKeyEvent(clientKey string, event *ListenEvent) {
	list := GetClientListeners(clientKey)
	if len(list) == 0 {
		return
	}
	for _, one := range list {
		one.AddEvent(event)
	}
}

func CallClientTabKeyEvent(clientTabKey string, event *ListenEvent) {
	listener := GetListener(clientTabKey)
	if listener == nil {
		return
	}
	listener.AddEvent(event)
}

func listenerInit() {

	go checkListener()
}

func checkListener() {
	var ticker = time.NewTicker(time.Second * 60) // 60 秒检测一次
	for {
		select {
		case <-ticker.C:
			list := GetListeners()
			for _, one := range list {
				nowTimeSecond := util.GetNowSecond()

				// 最后监听时间 在此 之前的 都为超时
				outTimeSecond := nowTimeSecond - listenerLastListenTimeoutSecond
				lastListenTimeSecond := util.GetSecondByTime(one.lastListenTime)
				if lastListenTimeSecond > outTimeSecond {
					continue
				}
				util.Logger.Debug("remove listener", zap.Any("listener", one))
				RemoveListener(one)
			}
		}
	}

}

type ListenEvent struct {
	KeyForRemoveDuplicates string      `json:"-"` // 用于剔除重复
	Event                  string      `json:"event"`
	Data                   interface{} `json:"data"`
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
	events         []*ListenEvent
	eventsLock     sync.Locker
}

func (this_ *ClientTabListener) AddEvent(event *ListenEvent) {
	if event == nil {
		return
	}
	this_.eventsLock.Lock()
	defer this_.eventsLock.Unlock()

	this_.events = append(this_.events, event)
}

func (this_ *ClientTabListener) Listen() []*ListenEvent {

	var lastListenTime = time.Now()
	this_.lastListenTime = lastListenTime
	defer func() {
		this_.lastListenTime = time.Now()
	}()
	var eventsList []*ListenEvent
	var ticker = time.NewTicker(time.Millisecond * 100)
	expireAt := lastListenTime.UnixMilli() + 1000*60 // 超时时间为 60 秒
	for {
		this_.lastListenTime = time.Now()
		select {
		case <-ticker.C:
			this_.eventsLock.Lock()
			eventsList = this_.events
			this_.events = make([]*ListenEvent, 0)
			this_.eventsLock.Unlock()
		}
		if len(eventsList) > 0 {
			break
		}
		// 判断是否已经等待10分钟
		if time.Now().UnixMilli() >= expireAt {
			break
		}
	}
	ticker.Stop()
	var events []*ListenEvent
	var eventCache = make(map[string]*ListenEvent)
	for _, event := range eventsList {
		if event.KeyForRemoveDuplicates == "" {
			events = append(events, event)
		} else {
			key := event.Event + "@" + event.KeyForRemoveDuplicates
			find := eventCache[key]
			if find != nil {
				find.Data = event.Data
			} else {
				eventCache[key] = event
				events = append(events, event)
			}
		}
	}

	return events
}
