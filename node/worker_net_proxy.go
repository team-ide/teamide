package node

import (
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

func (this_ *Worker) findNetProxyListByNodeId(nodeId string) (findList []*NetProxy) {
	var list = this_.cache.netProxyList
	for _, one := range list {
		if one.Inner.NodeId == nodeId {
			findList = append(findList, one)
		}
		if one.Outer.NodeId == nodeId {
			findList = append(findList, one)
		}
	}
	return
}

func (this_ *Worker) findNetProxyIdListByNodeId(nodeId string) (findList []string) {
	var list = this_.cache.netProxyList
	for _, one := range list {
		if one.Inner.NodeId == nodeId {
			findList = append(findList, one.Id)
		}
		if one.Outer.NodeId == nodeId {
			findList = append(findList, one.Id)
		}
	}
	return
}

func (this_ *Worker) addNetProxyList(netProxyList []*NetProxy) (err error) {
	if len(netProxyList) == 0 {
		return
	}
	this_.notifyAll(&Message{
		NetProxyList: netProxyList,
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

	var findChanged = false
	for _, netProxy := range netProxyList {

		var find = this_.findNetProxy(netProxy.Id)
		if find == nil {
			Logger.Info(this_.server.GetServerInfo()+" 添加网络代理 ", zap.Any("netProxy", netProxy))
			findChanged = true
			this_.cache.netProxyList = append(this_.cache.netProxyList, netProxy)
		} else {
			Logger.Info(this_.server.GetServerInfo()+" 更新网络代理 ", zap.Any("netProxy", netProxy))
			if netProxy.InnerStatus != 0 && netProxy.InnerStatus != find.InnerStatus {
				find.InnerStatus = netProxy.InnerStatus
				findChanged = true
			}
			if netProxy.OuterStatus != 0 && netProxy.OuterStatus != find.OuterStatus {
				find.OuterStatus = netProxy.OuterStatus
				findChanged = true
			}
			if netProxy.Enabled != 0 {
				if netProxy.IsEnabled() && find.IsEnabled() {
					find.Enabled = netProxy.Enabled
				} else {
					find.Enabled = netProxy.Enabled
					findChanged = true
				}
			}
		}
	}

	if findChanged {
		this_.refresh()
		if this_.server.OnNetProxyListChange != nil {
			this_.server.OnNetProxyListChange(this_.cache.netProxyList)
		}
	}
	return
}

func (this_ *Worker) doChangeNetProxyOuterStatus(netProxyId string, status int8, statusError string) (err error) {
	if netProxyId == "" || status == 0 {
		return
	}
	var find = this_.findNetProxy(netProxyId)
	var findChanged = false
	if find != nil {
		if find.OuterStatus != status || find.OuterStatusError != statusError {
			findChanged = true
			find.OuterStatus = status
			find.OuterStatusError = statusError
		}
	}
	if findChanged && this_.server.OnNetProxyListChange != nil {
		this_.server.OnNetProxyListChange(this_.cache.netProxyList)
	}
	return
}

func (this_ *Worker) doChangeNetProxyInnerStatus(netProxyId string, status int8, statusError string) (err error) {
	if netProxyId == "" || status == 0 {
		return
	}
	var find = this_.findNetProxy(netProxyId)
	var findChanged = false
	if find != nil {
		if find.InnerStatus != status || find.InnerStatusError != statusError {
			findChanged = true
			find.InnerStatus = status
			find.InnerStatusError = statusError
		}
	}
	if findChanged && this_.server.OnNetProxyListChange != nil {
		this_.server.OnNetProxyListChange(this_.cache.netProxyList)
	}
	return
}

func (this_ *Worker) removeNetProxyList(removeNetProxyIdList []string) (err error) {
	if len(removeNetProxyIdList) == 0 {
		return
	}

	this_.notifyAll(&Message{
		RemoveNetProxyIdList: removeNetProxyIdList,
	})

	err = this_.doRemoveNetProxyList(removeNetProxyIdList)

	return
}

func (this_ *Worker) doRemoveNetProxyList(removeNetProxyIdList []string) (err error) {
	if len(removeNetProxyIdList) == 0 {
		return
	}
	this_.cache.netProxyLock.Lock()
	defer this_.cache.netProxyLock.Unlock()
	Logger.Info(this_.server.GetServerInfo()+" 移除网络代理 ", zap.Any("removeNetProxyIdList", removeNetProxyIdList))

	for _, netProxyId := range removeNetProxyIdList {
		var find = this_.findNetProxy(netProxyId)
		if find != nil {
			if this_.server.rootNode.Id == find.Inner.NodeId {
				_ = this_.removeNetProxyInner(netProxyId)
			}
			if this_.server.rootNode.Id == find.Outer.NodeId {
				_ = this_.removeNetProxyOuter(netProxyId)
			}
		}
	}

	var findChanged = false
	var list = this_.cache.netProxyList
	var newList []*NetProxy
	for _, one := range list {
		if util.ContainsString(removeNetProxyIdList, one.Id) >= 0 {
			findChanged = true
		} else {
			newList = append(newList, one)
		}
	}
	this_.cache.netProxyList = newList
	if findChanged {
		this_.refresh()
		if this_.server.OnNetProxyListChange != nil {
			this_.server.OnNetProxyListChange(this_.cache.netProxyList)
		}
	}

	return
}
func (this_ *Worker) refreshNetProxy() {
	var list = this_.cache.netProxyList
	for _, netProxy := range list {
		if this_.server.rootNode.Id == netProxy.Inner.NodeId {
			if this_.server.rootNode.IsEnabled() && netProxy.IsEnabled() {
				_ = this_.newIfAbsentNetProxyInner(netProxy)
			} else {
				_ = this_.removeNetProxyInner(netProxy.Id)
			}
		}
		if this_.server.rootNode.Id == netProxy.Outer.NodeId {
			if this_.server.rootNode.IsEnabled() && netProxy.IsEnabled() {
				_ = this_.newIfAbsentNetProxyOuter(netProxy)
			} else {
				_ = this_.removeNetProxyOuter(netProxy.Id)
			}
		}
	}
	return
}

func (this_ *Worker) newIfAbsentNetProxyInner(netProxy *NetProxy) (inner *InnerServer) {
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
func (this_ *Worker) getNetProxyInner(netProxyId string) (inner *InnerServer) {
	this_.cache.netProxyInnerLock.Lock()
	defer this_.cache.netProxyInnerLock.Unlock()

	inner = this_.cache.netProxyInnerCache[netProxyId]
	return
}

func (this_ *Worker) removeNetProxyInner(netProxyId string) (inner *InnerServer) {
	this_.cache.netProxyInnerLock.Lock()
	defer this_.cache.netProxyInnerLock.Unlock()

	inner, ok := this_.cache.netProxyInnerCache[netProxyId]
	if ok {
		delete(this_.cache.netProxyInnerCache, netProxyId)
		inner.Stop()
	}
	return
}

func (this_ *Worker) newIfAbsentNetProxyOuter(netProxy *NetProxy) (outer *OuterListener) {
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

func (this_ *Worker) getNetProxyOuter(netProxyId string) (outer *OuterListener) {
	this_.cache.netProxyOuterLock.Lock()
	defer this_.cache.netProxyOuterLock.Unlock()

	outer = this_.cache.netProxyOuterCache[netProxyId]
	return
}

func (this_ *Worker) removeNetProxyOuter(netProxyId string) (inner *InnerServer) {
	this_.cache.netProxyOuterLock.Lock()
	defer this_.cache.netProxyOuterLock.Unlock()

	outer, ok := this_.cache.netProxyOuterCache[netProxyId]
	if ok {
		delete(this_.cache.netProxyOuterCache, netProxyId)
		outer.Stop()
	}
	return
}
