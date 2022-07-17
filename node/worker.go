package node

import (
	"sync"
)

type Worker struct {
	Node                          *Info
	childrenNodeList              []*Info
	isStop                        bool
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

func (this_ *Worker) isStopped() bool {

	return this_.isStop
}

func (this_ *Worker) initialize(nodeList []*Info, netProxyList []*NetProxy) {

	var oldList = this_.nodeList
	for _, one := range oldList {
		var find *Info
		for _, one_ := range nodeList {
			if one_.Id == one.Id {
				find = one_
			}
		}
		if find == nil {
			_ = this_.RemoveNode(one)
		}
	}
	for _, one := range nodeList {
		_ = this_.AddNode(one)
	}

	var oldNetProxyList = this_.netProxyList
	for _, one := range oldNetProxyList {
		var find *NetProxy
		for _, one_ := range netProxyList {
			if one_.Id == one.Id {
				find = one_
			}
		}
		if find == nil {
			_ = this_.RemoveNetProxy(one)
		}
	}
	for _, one := range netProxyList {
		_ = this_.AddNetProxy(one)
	}
}
