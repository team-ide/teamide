package node

import (
	"errors"
	"teamide/pkg/util"
)

func (this_ *Worker) sendToNext(lineNodeIdList []string, doSend func(listener *MessageListener) (err error)) (send bool, err error) {
	if len(lineNodeIdList) == 0 {
		err = errors.New("节点线不存在")
		return
	}
	var thisNodeIndex = util.ContainsString(lineNodeIdList, this_.server.Id)
	if thisNodeIndex != len(lineNodeIdList)-1 {
		nodeId := lineNodeIdList[thisNodeIndex+1]

		var pool = this_.cache.getNodeListenerPool(this_.server.Id, nodeId)
		if pool == nil {
			pool = this_.cache.getNodeListenerPool(nodeId, this_.server.Id)
		}
		if pool == nil {
			err = errors.New(this_.server.GetServerInfo() + " 与节点 [" + nodeId + "] 暂无通讯渠道")
			return
		}
		err = pool.Do(func(listener *MessageListener) (e error) {
			e = doSend(listener)
			return
		})
		if err != nil {
			return
		}
		send = true
		return
	}
	return
}
func (this_ *Worker) netProxyNewConn(lineNodeIdList []string, netProxyId string, connId string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodNetProxyNewConn, &Message{
			LineNodeIdList: lineNodeIdList,
			NetProxyId:     netProxyId,
			ConnId:         connId,
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}
	var netProxy = this_.findNetProxy(netProxyId)
	if netProxy == nil {
		err = errors.New("网络代理[" + netProxyId + "]不存在")
		return
	}
	outer := this_.getNetProxyOuter(netProxy.Id)
	if outer != nil {
		err = outer.newConn(connId)
	}
	if err != nil {
		return
	}
	return
}

func (this_ *Worker) netProxyCloseConn(isReverse bool, lineNodeIdList []string, netProxyId string, connId string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodNetProxyCloseConn, &Message{
			LineNodeIdList: lineNodeIdList,
			NetProxyId:     netProxyId,
			ConnId:         connId,
			IsReverse:      isReverse,
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}
	var netProxy = this_.findNetProxy(netProxyId)
	if netProxy == nil {
		err = errors.New("网络代理[" + netProxyId + "]不存在")
		return
	}
	if isReverse {
		inner := this_.getNetProxyInner(netProxy.Id)
		if inner != nil {
			err = inner.closeConn(connId)
		}
	} else {
		outer := this_.getNetProxyOuter(netProxy.Id)
		if outer != nil {
			err = outer.closeConn(connId)
		}
	}
	if err != nil {
		return
	}
	return
}

func (this_ *Worker) netProxySend(isReverse bool, lineNodeIdList []string, netProxyId string, connId string, bytes []byte) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodNetProxySend, &Message{
			LineNodeIdList: lineNodeIdList,
			NetProxyId:     netProxyId,
			ConnId:         connId,
			Bytes:          bytes,
			IsReverse:      isReverse,
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}
	var netProxy = this_.findNetProxy(netProxyId)
	if netProxy == nil {
		err = errors.New("网络代理[" + netProxyId + "]不存在")
		return
	}
	if isReverse {
		inner := this_.getNetProxyInner(netProxy.Id)
		if inner != nil {
			err = inner.send(connId, bytes)
		}
	} else {
		outer := this_.getNetProxyOuter(netProxy.Id)
		if outer != nil {
			err = outer.send(connId, bytes)
		}
	}
	if err != nil {
		return
	}
	return
}
