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
	initialization           bool
}

func (this_ *Worker) isStopped() bool {

	return this_.isStop
}

func (this_ *Worker) initialize(nodeList []*Info, netProxyList []*NetProxy) {
	if this_.initialization {
		return
	}
	this_.initialization = true

	for _, one := range nodeList {
		_ = this_.AddNode(one)
	}
	for _, one := range netProxyList {
		_ = this_.AddNetProxy(one)
	}
}
