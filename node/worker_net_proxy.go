package node

import (
	"errors"
	"go.uber.org/zap"
	"teamide/pkg/util"
)

func (this_ *Worker) findNetProxy(id string) (find *NetProxy) {
	var list = this_.netProxyList
	for _, one := range list {
		if one.Id == id {
			find = one
		}
	}
	return
}

func (this_ *Worker) formatNetProxy(netProxy *NetProxy) (err error) {
	if netProxy == nil {
		err = errors.New("网络代理配置不能为空")
		return
	}
	if netProxy.Id == "" {
		err = errors.New("网络代理编号不能为空")
		return
	}
	if netProxy.Inner == nil {
		err = errors.New("网络代理输入配置不能为空")
		return
	}
	if netProxy.Inner.GetAddress() == "" {
		err = errors.New("网络代理输入地址不能为空")
		return
	}
	if netProxy.Outer == nil {
		err = errors.New("网络代理输出配置不能为空")
		return
	}
	if netProxy.Outer.GetAddress() == "" {
		err = errors.New("网络代理输出地址不能为空")
		return
	}
	if len(netProxy.LineNodeIdList) == 0 {
		netProxy.LineNodeIdList = this_.GetNodeLineByFromTo(netProxy.Inner.NodeId, netProxy.Outer.NodeId)
		if len(netProxy.LineNodeIdList) == 0 {
			err = errors.New("无法正确解析输入输出节点关系")
			return
		}
		for i := len(netProxy.LineNodeIdList) - 1; i >= 0; i-- {
			netProxy.ReverseLineNodeIdList = append(netProxy.ReverseLineNodeIdList, netProxy.LineNodeIdList[i])
		}

	}

	return
}

func (this_ *Worker) AddNetProxy(netProxy *NetProxy) (err error) {
	this_.netProxyLock.Lock()
	defer this_.netProxyLock.Unlock()
	Logger.Info(this_.Node.GetNodeStr()+" 添加网络代理 ", zap.Any("netProxy", netProxy))

	err = this_.formatNetProxy(netProxy)
	if err != nil {
		return
	}

	var list = this_.childrenNodeList
	for _, one := range list {
		pool := this_.getChildrenNodeListenerPool(one)
		_ = pool.Do(func(listener *MessageListener) (e error) {
			e = listener.Send(&Message{
				Token:    one.Token,
				Method:   methodNetProxyAdd,
				NetProxy: netProxy,
			})
			return
		})
	}

	this_.netProxyList = append(this_.netProxyList, netProxy)

	this_.refreshNetProxy()
	return
}

func (this_ *Worker) RemoveNetProxy(netProxy *NetProxy) (err error) {
	this_.netProxyLock.Lock()
	defer this_.netProxyLock.Unlock()
	Logger.Info(this_.Node.GetNodeStr()+" 移除网络代理 ", zap.Any("netProxy", netProxy))

	var list = this_.childrenNodeList
	for _, one := range list {
		pool := this_.getChildrenNodeListenerPool(one)
		_ = pool.Do(func(listener *MessageListener) (e error) {
			e = listener.Send(&Message{
				Token:    one.Token,
				Method:   methodNetProxyRemove,
				NetProxy: netProxy,
			})
			return
		})
	}

	return
}

func (this_ *Worker) refreshNetProxy() {
	var list = this_.netProxyList
	for _, one := range list {
		if this_.Node.Id == one.Inner.NodeId {
			_ = this_.getNetProxyInner(one)
		}
		if this_.Node.Id == one.Outer.NodeId {
			_ = this_.getNetProxyOuter(one)
		}
	}
	return
}

func (this_ *Worker) getNetProxyInner(netProxy *NetProxy) (inner *InnerServer) {
	this_.netProxyInnerLock.Lock()
	defer this_.netProxyInnerLock.Unlock()

	inner, ok := this_.netProxyInnerCache[netProxy.Id]
	if !ok {
		inner = &InnerServer{
			netProxy: netProxy,
			Worker:   this_,
		}
		inner.Start()
		this_.netProxyInnerCache[netProxy.Id] = inner
	}
	return
}

func (this_ *Worker) getNetProxyOuter(netProxy *NetProxy) (outer *OuterListener) {
	this_.netProxyOuterLock.Lock()
	defer this_.netProxyOuterLock.Unlock()

	outer, ok := this_.netProxyOuterCache[netProxy.Id]
	if !ok {
		outer = &OuterListener{
			netProxy: netProxy,
			Worker:   this_,
		}
		outer.Start()
		this_.netProxyOuterCache[netProxy.Id] = outer
	}
	return
}

func (this_ *Worker) sendToNext(lineNodeIdList []string, doSend func(node *Info, listener *MessageListener) (err error)) (send bool, err error) {
	if len(lineNodeIdList) == 0 {
		err = errors.New("节点线不存在")
		return
	}
	var thisNodeIndex = util.ContainsString(lineNodeIdList, this_.Node.Id)
	if thisNodeIndex != len(lineNodeIdList)-1 {
		nodeId := lineNodeIdList[thisNodeIndex+1]
		var find = this_.findChildrenNode(nodeId)
		var pool *MessageListenerPool
		if find == nil {
			find = this_.findNode(nodeId)
			if find == nil {
				err = errors.New("节点[" + nodeId + "]不存在")
				return
			}
			pool = this_.getFromNodeListenerPool(nodeId)
		} else {
			pool = this_.getChildrenNodeListenerPool(find)
		}
		err = pool.Do(func(listener *MessageListener) (e error) {
			e = doSend(find, listener)
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
	send, err := this_.sendToNext(lineNodeIdList, func(node *Info, listener *MessageListener) (e error) {
		_, e = this_.Call(node, listener, methodNetProxyNewConn, &Message{
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
	outer := this_.getNetProxyOuter(netProxy)
	err = outer.newConn(connId)
	if err != nil {
		return
	}
	return
}

func (this_ *Worker) netProxyCloseConn(isReverse bool, lineNodeIdList []string, netProxyId string, connId string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, func(node *Info, listener *MessageListener) (e error) {
		_, e = this_.Call(node, listener, methodNetProxyCloseConn, &Message{
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
		inner := this_.getNetProxyInner(netProxy)
		err = inner.closeConn(connId)
	} else {
		outer := this_.getNetProxyOuter(netProxy)
		err = outer.closeConn(connId)
	}
	if err != nil {
		return
	}
	return
}

func (this_ *Worker) netProxySend(isReverse bool, lineNodeIdList []string, netProxyId string, connId string, bytes []byte) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, func(node *Info, listener *MessageListener) (e error) {
		_, e = this_.Call(node, listener, methodNetProxySend, &Message{
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
		inner := this_.getNetProxyInner(netProxy)
		err = inner.send(connId, bytes)
	} else {
		outer := this_.getNetProxyOuter(netProxy)
		err = outer.send(connId, bytes)
	}
	if err != nil {
		return
	}
	return
}
