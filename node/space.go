package node

import (
	"sync"
)

type Space struct {
	spaceId string

	toNodeList     []*ToNode
	toNodeListLock sync.Mutex

	fromNodeList     []*FromNode
	fromNodeListLock sync.Mutex

	toNodeListenerPoolCache     map[string]*MessageListenerPool
	toNodeListenerPoolCacheLock sync.Mutex

	fromNodeListenerPoolCache     map[string]*MessageListenerPool
	fromNodeListenerPoolCacheLock sync.Mutex

	netProxyInnerList     []*NetProxy
	netProxyInnerListLock sync.Mutex

	netProxyOuterList     []*NetProxy
	netProxyOuterListLock sync.Mutex

	netProxyInnerCache     map[string]*InnerServer
	netProxyInnerCacheLock sync.Mutex

	netProxyOuterCache     map[string]*OuterListener
	netProxyOuterCacheLock sync.Mutex

	callbackCache     map[string]func(msg *Message)
	callbackCacheLock sync.Mutex
}

func newSpace(spaceId string) *Space {
	return &Space{
		spaceId:                   spaceId,
		toNodeListenerPoolCache:   make(map[string]*MessageListenerPool),
		fromNodeListenerPoolCache: make(map[string]*MessageListenerPool),
		callbackCache:             make(map[string]func(msg *Message)),
		netProxyInnerCache:        make(map[string]*InnerServer),
		netProxyOuterCache:        make(map[string]*OuterListener),
	}
}

type ToNode struct {
	Id          string `json:"id,omitempty"`
	ConnAddress string `json:"connAddress,omitempty"`
	ConnToken   string `json:"connToken,omitempty"`
	ConnSize    int    `json:"connSize,omitempty"`
	Enabled     int8   `json:"enabled,omitempty"`
}

func (this_ *ToNode) IsEnabled() bool {
	return this_.Enabled != 2
}

type FromNode struct {
	Id        string `json:"id,omitempty"`
	FromToken string `json:"connToken,omitempty"`
	Enabled   int8   `json:"enabled,omitempty"`
}

func (this_ *FromNode) IsEnabled() bool {
	return this_.Enabled != 2
}

func (this_ *Space) findToNode(id string) (find *ToNode) {
	var list = this_.toNodeList
	for _, one := range list {
		if one.Id == id {
			find = one
		}
	}
	return
}

func (this_ *Space) findFromNode(id string) (find *FromNode) {
	var list = this_.fromNodeList
	for _, one := range list {
		if one.Id == id {
			find = one
		}
	}
	return
}

func (this_ *Space) findInnerNetProxy(id string) (find *NetProxy) {
	var list = this_.netProxyInnerList
	for _, one := range list {
		if one.Id == id {
			find = one
		}
	}
	return
}

func (this_ *Space) findOuterNetProxy(id string) (find *NetProxy) {
	var list = this_.netProxyOuterList
	for _, one := range list {
		if one.Id == id {
			find = one
		}
	}
	return
}

func (this_ *Space) newToNodeListenerPoolIfAbsent(toNodeId string) (pool *MessageListenerPool) {
	this_.toNodeListenerPoolCacheLock.Lock()
	defer this_.toNodeListenerPoolCacheLock.Unlock()

	pool, ok := this_.toNodeListenerPoolCache[toNodeId]
	if !ok {
		pool = &MessageListenerPool{}
		this_.toNodeListenerPoolCache[toNodeId] = pool
	}
	return
}

func (this_ *Space) getToNodeListenerPool(toNodeId string) (pool *MessageListenerPool) {
	this_.toNodeListenerPoolCacheLock.Lock()
	defer this_.toNodeListenerPoolCacheLock.Unlock()

	pool, _ = this_.toNodeListenerPoolCache[toNodeId]
	return
}

func (this_ *Space) removeToNodeListenerPool(toNodeId string) (pool *MessageListenerPool) {
	this_.toNodeListenerPoolCacheLock.Lock()
	defer this_.toNodeListenerPoolCacheLock.Unlock()

	pool, ok := this_.toNodeListenerPoolCache[toNodeId]
	if ok {
		delete(this_.toNodeListenerPoolCache, toNodeId)
		pool.Stop()
	}
	return
}

func (this_ *Space) getToNodeListenerPoolList() (poolList []*MessageListenerPool) {
	this_.toNodeListenerPoolCacheLock.Lock()
	defer this_.toNodeListenerPoolCacheLock.Unlock()

	for _, one := range this_.toNodeListenerPoolCache {
		poolList = append(poolList, one)
	}
	return
}

func (this_ *Space) newFromNodeListenerPoolIfAbsent(fromNodeId string) (pool *MessageListenerPool) {
	this_.fromNodeListenerPoolCacheLock.Lock()
	defer this_.fromNodeListenerPoolCacheLock.Unlock()

	pool, ok := this_.fromNodeListenerPoolCache[fromNodeId]
	if !ok {
		pool = &MessageListenerPool{}
		this_.fromNodeListenerPoolCache[fromNodeId] = pool
	}
	return
}

func (this_ *Space) getFromNodeListenerPool(fromNodeId string) (pool *MessageListenerPool) {
	this_.fromNodeListenerPoolCacheLock.Lock()
	defer this_.fromNodeListenerPoolCacheLock.Unlock()

	pool, _ = this_.fromNodeListenerPoolCache[fromNodeId]
	return
}

func (this_ *Space) removeFromNodeListenerPool(fromNodeId string) (pool *MessageListenerPool) {
	this_.fromNodeListenerPoolCacheLock.Lock()
	defer this_.fromNodeListenerPoolCacheLock.Unlock()

	pool, ok := this_.fromNodeListenerPoolCache[fromNodeId]
	if ok {
		delete(this_.fromNodeListenerPoolCache, fromNodeId)
		pool.Stop()
	}
	return
}

func (this_ *Space) getFromNodeListenerPoolList() (poolList []*MessageListenerPool) {
	this_.fromNodeListenerPoolCacheLock.Lock()
	defer this_.fromNodeListenerPoolCacheLock.Unlock()

	for _, one := range this_.fromNodeListenerPoolCache {
		poolList = append(poolList, one)
	}
	return
}

func (this_ *Space) getCallback(id string) (callback func(msg *Message), ok bool) {
	this_.callbackCacheLock.Lock()
	defer this_.callbackCacheLock.Unlock()

	callback, ok = this_.callbackCache[id]
	return
}

func (this_ *Space) setCallback(id string, callback func(msg *Message)) {
	this_.callbackCacheLock.Lock()
	defer this_.callbackCacheLock.Unlock()

	this_.callbackCache[id] = callback
}

func (this_ *Space) removeCallback(id string) {
	this_.callbackCacheLock.Lock()
	defer this_.callbackCacheLock.Unlock()

	delete(this_.callbackCache, id)
}
