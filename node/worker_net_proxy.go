package node

import (
	"errors"
	"go.uber.org/zap"
	"teamide/pkg/util"
)

func (this_ *Worker) findNetProxy(id string) (find *NetProxy) {
	var list = this_.cache.netProxyList
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
		netProxy.LineNodeIdList = this_.getNodeLineByFromTo(netProxy.Inner.NodeId, netProxy.Outer.NodeId)
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

func (this_ *Worker) addNetProxyList(netProxyList []*NetProxy, calledNodeIdList []string) (err error) {
	if len(netProxyList) == 0 {
		return
	}
	this_.sendAllTo(&Message{
		Method:           methodNetProxyAdd,
		NetProxyList:     netProxyList,
		CalledNodeIdList: calledNodeIdList,
	})
	err = this_.doAddNetProxyList(netProxyList)
	return
}

func (this_ *Worker) doAddNetProxyList(netProxyList []*NetProxy) (err error) {
	if len(netProxyList) == 0 {
		return
	}
	this_.cache.netProxyLock.Lock()
	defer this_.cache.netProxyLock.Unlock()
	Logger.Info(this_.server.GetServerInfo()+" 添加网络代理 ", zap.Any("netProxyList", netProxyList))

	for _, netProxy := range netProxyList {

		var find = this_.findNetProxy(netProxy.Id)
		if find == nil {
			this_.cache.netProxyList = append(this_.cache.netProxyList, netProxy)
		}
	}

	this_.refreshNetProxy()
	return
}

func (this_ *Worker) removeNetProxyList(netProxyIdList []string, calledNodeIdList []string) (err error) {
	if len(netProxyIdList) == 0 {
		return
	}

	this_.sendAllTo(&Message{
		Method:           methodNetProxyRemove,
		NetProxyIdList:   netProxyIdList,
		CalledNodeIdList: calledNodeIdList,
	})

	err = this_.doRemoveNetProxyList(netProxyIdList)

	return
}

func (this_ *Worker) doRemoveNetProxyList(netProxyIdList []string) (err error) {
	if len(netProxyIdList) == 0 {
		return
	}
	this_.cache.netProxyLock.Lock()
	defer this_.cache.netProxyLock.Unlock()
	Logger.Info(this_.server.GetServerInfo()+" 移除网络代理 ", zap.Any("netProxyIdList", netProxyIdList))

	for _, netProxyId := range netProxyIdList {

		var find = this_.findNetProxy(netProxyId)
		if find != nil {
			if this_.server.Id == find.Inner.NodeId {
				_ = this_.removeNetProxyInner(netProxyId)
			}
			if this_.server.Id == find.Outer.NodeId {
				_ = this_.removeNetProxyOuter(netProxyId)
			}

			list := this_.cache.netProxyList
			var newList []*NetProxy
			for _, one := range list {
				if one != find {
					newList = append(newList, one)
				}
			}
			this_.cache.netProxyList = newList
		}

	}
	this_.refreshNetProxy()

	return
}
func (this_ *Worker) refreshNetProxy() {
	var list = this_.cache.netProxyList
	for _, one := range list {
		if this_.server.Id == one.Inner.NodeId {
			_ = this_.getNetProxyInner(one)
		}
		if this_.server.Id == one.Outer.NodeId {
			_ = this_.getNetProxyOuter(one)
		}
	}

	if this_.server.OnNetProxyListChange != nil {
		this_.server.OnNetProxyListChange(this_.cache.netProxyList)
	}
	return
}

func (this_ *Worker) getNetProxyInner(netProxy *NetProxy) (inner *InnerServer) {
	this_.cache.netProxyInnerLock.Lock()
	defer this_.cache.netProxyInnerLock.Unlock()

	inner, ok := this_.cache.netProxyInnerCache[netProxy.Id]
	if !ok {
		inner = &InnerServer{
			netProxy: netProxy,
			Worker:   this_,
		}
		inner.Start()
		this_.cache.netProxyInnerCache[netProxy.Id] = inner
	}
	return
}

func (this_ *Worker) removeNetProxyInner(netProxyId string) (inner *InnerServer) {
	this_.cache.netProxyInnerLock.Lock()
	defer this_.cache.netProxyInnerLock.Unlock()

	inner, ok := this_.cache.netProxyInnerCache[netProxyId]
	if ok {
		inner.Stop()
		delete(this_.cache.netProxyInnerCache, netProxyId)
	}
	return
}

func (this_ *Worker) getNetProxyOuter(netProxy *NetProxy) (outer *OuterListener) {
	this_.cache.netProxyOuterLock.Lock()
	defer this_.cache.netProxyOuterLock.Unlock()

	outer, ok := this_.cache.netProxyOuterCache[netProxy.Id]
	if !ok {
		outer = &OuterListener{
			netProxy: netProxy,
			Worker:   this_,
		}
		outer.Start()
		this_.cache.netProxyOuterCache[netProxy.Id] = outer
	}
	return
}

func (this_ *Worker) removeNetProxyOuter(netProxyId string) (inner *InnerServer) {
	this_.cache.netProxyOuterLock.Lock()
	defer this_.cache.netProxyOuterLock.Unlock()

	outer, ok := this_.cache.netProxyOuterCache[netProxyId]
	if ok {
		outer.Stop()
		delete(this_.cache.netProxyOuterCache, netProxyId)
	}
	return
}

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
	outer := this_.getNetProxyOuter(netProxy)
	err = outer.newConn(connId)
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
