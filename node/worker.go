package node

import (
	"sync"
)

type Worker struct {
	Node                      *Info
	childrenNodeList          []*Info
	isStop                    bool
	nodeList                  []*Info
	childrenNodeListenerCache map[string]*MessageListener

	inPortForwardingList  []*NetProxy
	outPortForwardingList []*NetProxy
	netProxyList          []*NetProxy

	callbackCache     map[string]func(msg *Message)
	nodeLock          sync.RWMutex
	nodeListenerLock  sync.RWMutex
	callbackCacheLock sync.RWMutex
	netProxyLock      sync.RWMutex
}

func (this_ *Worker) isStopped() bool {

	return this_.isStop
}
