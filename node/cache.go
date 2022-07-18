package node

import (
	"sync"
	"teamide/pkg/util"
)

type Cache struct {
	nodeListenerPoolCache map[string]*MessageListenerPool
	nodeListenerLock      sync.Mutex

	fromNodeIdList []string
	fromNodeIdLock sync.Mutex

	childrenNodeList              []*Info
	nodeList                      []*Info
	childrenNodeListenerPoolCache map[string]*MessageListenerPool
	fromNodeListenerPoolCache     map[string]*MessageListenerPool

	netProxyInnerCache map[string]*InnerServer
	netProxyOuterCache map[string]*OuterListener
	netProxyList       []*NetProxy

	callbackCache            map[string]func(msg *Message)
	nodeLock                 sync.Mutex
	childrenNodeListenerLock sync.Mutex
	fromNodeListenerLock     sync.Mutex
	callbackCacheLock        sync.Mutex
	netProxyLock             sync.Mutex
	netProxyInnerLock        sync.Mutex
	netProxyOuterLock        sync.Mutex
}

func newCache() *Cache {
	return &Cache{
		nodeListenerPoolCache: make(map[string]*MessageListenerPool),
		callbackCache:         make(map[string]func(msg *Message)),
		netProxyInnerCache:    make(map[string]*InnerServer),
		netProxyOuterCache:    make(map[string]*OuterListener),
	}
}

func (this_ *Cache) newIfAbsentNodeListenerPool(fromNodeId string, toNodeId string) (pool *MessageListenerPool) {
	this_.nodeListenerLock.Lock()
	defer this_.nodeListenerLock.Unlock()

	pool, ok := this_.nodeListenerPoolCache[fromNodeId+":"+toNodeId]
	if !ok {
		pool = &MessageListenerPool{}
		this_.nodeListenerPoolCache[fromNodeId+":"+toNodeId] = pool
	}
	return
}

func (this_ *Cache) getNodeListenerPool(fromNodeId string, toNodeId string) (pool *MessageListenerPool) {
	this_.nodeListenerLock.Lock()
	defer this_.nodeListenerLock.Unlock()

	pool, _ = this_.nodeListenerPoolCache[fromNodeId+":"+toNodeId]
	return
}

func (this_ *Cache) removeNodeListenerPool(fromNodeId string, toNodeId string) (pool *MessageListenerPool) {
	this_.nodeListenerLock.Lock()
	defer this_.nodeListenerLock.Unlock()

	pool, ok := this_.nodeListenerPoolCache[fromNodeId+":"+toNodeId]
	if ok {
		delete(this_.nodeListenerPoolCache, fromNodeId+":"+toNodeId)
		pool.Stop()
	}
	return
}

func (this_ *Cache) addFromNodeId(fromNodeId string) {
	this_.fromNodeIdLock.Lock()
	defer this_.fromNodeIdLock.Unlock()

	if util.ContainsString(this_.fromNodeIdList, fromNodeId) < 0 {
		this_.fromNodeIdList = append(this_.fromNodeIdList, fromNodeId)
	}
	return
}

func (this_ *Cache) removeFromNodeId(fromNodeId string) {
	this_.fromNodeIdLock.Lock()
	defer this_.fromNodeIdLock.Unlock()
	var list = this_.fromNodeIdList
	var newList []string
	for _, one := range list {
		if one != fromNodeId {
			newList = append(newList, one)
		}
	}
	this_.fromNodeIdList = newList
	return
}

func (this_ *Cache) getCallback(id string) (callback func(msg *Message), ok bool) {
	this_.callbackCacheLock.Lock()
	defer this_.callbackCacheLock.Unlock()
	callback, ok = this_.callbackCache[id]
	return
}

func (this_ *Cache) setCallback(id string, callback func(msg *Message)) {
	this_.callbackCacheLock.Lock()
	defer this_.callbackCacheLock.Unlock()

	this_.callbackCache[id] = callback
}

func (this_ *Cache) removeCallback(id string) {
	this_.callbackCacheLock.Lock()
	defer this_.callbackCacheLock.Unlock()

	delete(this_.callbackCache, id)
}
