package node

import (
	"fmt"
	"sync"
	"teamide/pkg/terminal"
)

type Space struct {
	toNodeList     []*ToNode
	toNodeListLock sync.Mutex

	fromNodeList     []*FromNode
	fromNodeListLock sync.Mutex

	toNodeListenerPoolCache     map[string]*MessageListenerPool
	toNodeListenerPoolCacheLock sync.Mutex

	fromNodeListenerPoolCache     map[string]*MessageListenerPool
	fromNodeListenerPoolCacheLock sync.Mutex

	netProxyInnerList     []*NetProxyInner
	netProxyInnerListLock sync.Mutex

	netProxyOuterList     []*NetProxyOuter
	netProxyOuterListLock sync.Mutex

	netProxyInnerCache     map[string]*InnerServer
	netProxyInnerCacheLock sync.Mutex

	netProxyOuterCache     map[string]*OuterListener
	netProxyOuterCacheLock sync.Mutex

	callbackCache     map[string]func(msg *Message)
	callbackCacheLock sync.Mutex

	terminalServiceCache     map[string]terminal.Service
	terminalServiceCacheLock sync.Mutex

	toNodeListenerKeepAliveLock sync.Mutex

	onBytesCache     map[string]*OnBytes
	onBytesCacheLock sync.Mutex
}

type OnBytes struct {
	start func() (err error)
	end   func() (err error)
	on    func(buf []byte) (err error)
}

func (this_ *Space) addOnBytesCache(key string, onBytes *OnBytes) {
	this_.onBytesCacheLock.Lock()
	defer this_.onBytesCacheLock.Unlock()

	this_.onBytesCache[key] = onBytes
	return
}

func (this_ *Space) getOnBytesCache(key string) (onBytes *OnBytes) {
	this_.onBytesCacheLock.Lock()
	defer this_.onBytesCacheLock.Unlock()

	onBytes = this_.onBytesCache[key]
	return
}

func (this_ *Space) removeOnBytesCache(key string) {
	this_.onBytesCacheLock.Lock()
	defer this_.onBytesCacheLock.Unlock()

	delete(this_.onBytesCache, key)
	return
}

func (this_ *Space) addTerminalService(key string, one terminal.Service) {
	this_.terminalServiceCacheLock.Lock()
	defer this_.terminalServiceCacheLock.Unlock()

	this_.terminalServiceCache[key] = one
	return
}

func (this_ *Space) getTerminalService(key string) (res terminal.Service) {
	this_.terminalServiceCacheLock.Lock()
	defer this_.terminalServiceCacheLock.Unlock()

	res = this_.terminalServiceCache[key]
	return
}

func (this_ *Space) removeTerminalService(key string) {
	this_.terminalServiceCacheLock.Lock()
	defer this_.terminalServiceCacheLock.Unlock()

	delete(this_.terminalServiceCache, key)
	return
}

func newSpace() *Space {
	return &Space{
		toNodeListenerPoolCache:   make(map[string]*MessageListenerPool),
		fromNodeListenerPoolCache: make(map[string]*MessageListenerPool),
		callbackCache:             make(map[string]func(msg *Message)),
		netProxyInnerCache:        make(map[string]*InnerServer),
		netProxyOuterCache:        make(map[string]*OuterListener),
		onBytesCache:              make(map[string]*OnBytes),
		terminalServiceCache:      make(map[string]terminal.Service),
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

type NetProxyInner struct {
	Id             string   `json:"id,omitempty"`
	NodeId         string   `json:"nodeId,omitempty"`
	Type           string   `json:"type,omitempty"`
	Address        string   `json:"address,omitempty"`
	LineNodeIdList []string `json:"lineNodeIdList,omitempty"`
	Enabled        int8     `json:"enabled,omitempty"`
}

func (this_ *NetProxyInner) IsEnabled() bool {
	return this_.Enabled != 2
}

func (this_ *NetProxyInner) GetInfoStr() (str string) {
	return fmt.Sprintf("[%s][%s]", this_.GetType(), this_.Address)
}

func (this_ *NetProxyInner) GetType() (str string) {
	var t = this_.Type
	if t == "" {
		t = "tcp"
	}
	return t
}

func (this_ *NetProxyInner) GetAddress() (str string) {
	return GetAddress(this_.Address)
}

type NetProxyOuter struct {
	Id                    string   `json:"id,omitempty"`
	NodeId                string   `json:"nodeId,omitempty"`
	Type                  string   `json:"type,omitempty"`
	Address               string   `json:"address,omitempty"`
	ReverseLineNodeIdList []string `json:"reverseLineNodeIdList,omitempty"`
	Enabled               int8     `json:"enabled,omitempty"`
}

func (this_ *NetProxyOuter) IsEnabled() bool {
	return this_.Enabled != 2
}

func (this_ *NetProxyOuter) GetInfoStr() (str string) {
	return fmt.Sprintf("[%s][%s]", this_.GetType(), this_.Address)
}

func (this_ *NetProxyOuter) GetType() (str string) {
	var t = this_.Type
	if t == "" {
		t = "tcp"
	}
	return t
}

func (this_ *NetProxyOuter) GetAddress() (str string) {
	return GetAddress(this_.Address)
}

func GetAddress(address string) (str string) {
	if address == "" {
		return ""
	}
	return address
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

func (this_ *Space) findInnerNetProxy(id string) (find *NetProxyInner) {
	var list = this_.netProxyInnerList
	for _, one := range list {
		if one.Id == id {
			find = one
		}
	}
	return
}

func (this_ *Space) findOuterNetProxy(id string) (find *NetProxyOuter) {
	var list = this_.netProxyOuterList
	for _, one := range list {
		if one.Id == id {
			find = one
		}
	}
	return
}

func (this_ *Space) getToNodeListenerPoolIfAbsentCreate(toNodeId string) (pool *MessageListenerPool) {
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

func (this_ *Space) removeToNodeListenerPoolList() {
	this_.toNodeListenerPoolCacheLock.Lock()
	defer this_.toNodeListenerPoolCacheLock.Unlock()

	for _, pool := range this_.toNodeListenerPoolCache {
		pool.Stop()
	}
	this_.toNodeListenerPoolCache = map[string]*MessageListenerPool{}
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

func (this_ *Space) getFromNodeListenerPoolIfAbsentCreate(fromNodeId string) (pool *MessageListenerPool) {
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

func (this_ *Space) removeFromNodeListenerPoolList() {
	this_.fromNodeListenerPoolCacheLock.Lock()
	defer this_.fromNodeListenerPoolCacheLock.Unlock()

	for _, pool := range this_.fromNodeListenerPoolCache {
		pool.Stop()
	}
	this_.fromNodeListenerPoolCache = map[string]*MessageListenerPool{}
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

func (this_ *Space) getNetProxyInnerIfAbsentCreate(netProxy *NetProxyInner, worker *Worker) (inner *InnerServer) {
	this_.netProxyInnerCacheLock.Lock()
	defer this_.netProxyInnerCacheLock.Unlock()

	inner, ok := this_.netProxyInnerCache[netProxy.Id]
	if !ok {
		inner = &InnerServer{
			netProxy: netProxy,
			worker:   worker,
		}
		inner.Start()
		this_.netProxyInnerCache[netProxy.Id] = inner
	}
	return
}
func (this_ *Space) getNetProxyInner(netProxyId string) (inner *InnerServer) {
	this_.netProxyInnerCacheLock.Lock()
	defer this_.netProxyInnerCacheLock.Unlock()

	inner = this_.netProxyInnerCache[netProxyId]
	return
}

func (this_ *Space) removeNetProxyInner(netProxyId string) (inner *InnerServer) {
	this_.netProxyInnerCacheLock.Lock()
	defer this_.netProxyInnerCacheLock.Unlock()

	inner, ok := this_.netProxyInnerCache[netProxyId]
	if ok {
		delete(this_.netProxyInnerCache, netProxyId)
		inner.Stop()
	}
	return
}

func (this_ *Space) getNetProxyOuterIfAbsentCreate(netProxy *NetProxyOuter, worker *Worker) (outer *OuterListener) {
	this_.netProxyOuterCacheLock.Lock()
	defer this_.netProxyOuterCacheLock.Unlock()

	outer, ok := this_.netProxyOuterCache[netProxy.Id]
	if !ok {
		outer = &OuterListener{
			netProxy: netProxy,
			worker:   worker,
		}
		outer.Start()
		this_.netProxyOuterCache[netProxy.Id] = outer
	}
	return
}

func (this_ *Space) getNetProxyOuter(netProxyId string) (outer *OuterListener) {
	this_.netProxyOuterCacheLock.Lock()
	defer this_.netProxyOuterCacheLock.Unlock()

	outer = this_.netProxyOuterCache[netProxyId]
	return
}

func (this_ *Space) removeNetProxyOuter(netProxyId string) (inner *InnerServer) {
	this_.netProxyOuterCacheLock.Lock()
	defer this_.netProxyOuterCacheLock.Unlock()

	outer, ok := this_.netProxyOuterCache[netProxyId]
	if ok {
		delete(this_.netProxyOuterCache, netProxyId)
		outer.Stop()
	}
	return
}
