package node

import (
	"go.uber.org/zap"
	"teamide/pkg/util"
)

func (this_ *Worker) doAddNetProxyInnerList(netProxyList []*NetProxyInner) (err error) {
	if len(netProxyList) == 0 {
		return
	}
	this_.netProxyInnerListLock.Lock()
	defer this_.netProxyInnerListLock.Unlock()

	for _, netProxy := range netProxyList {

		var find = this_.findInnerNetProxy(netProxy.Id)
		if find == nil {
			Logger.Info(this_.server.GetServerInfo()+" doAddNetProxyInnerList ", zap.Any("netProxy", netProxy))
			this_.netProxyInnerList = append(this_.netProxyInnerList, netProxy)

			if netProxy.IsEnabled() {
				_ = this_.getNetProxyInnerIfAbsentCreate(netProxy, this_)
			}
		} else {

			var hasChange bool
			if netProxy.Enabled != 0 {
				if netProxy.IsEnabled() != find.IsEnabled() {
					hasChange = true
				}
				find.Enabled = netProxy.Enabled
			}
			if netProxy.Type != find.Type {
				hasChange = true
				find.Type = netProxy.Type
			}
			if netProxy.Address != find.Address {
				hasChange = true
				find.Address = netProxy.Address
			}

			if hasChange {
				Logger.Info(this_.server.GetServerInfo()+" 更新网络代理 ", zap.Any("netProxy", netProxy))
				_ = this_.removeNetProxyInner(netProxy.Id)
				if find.IsEnabled() {
					_ = this_.getNetProxyInnerIfAbsentCreate(netProxy, this_)
				}
			}

		}
	}

	return
}

func (this_ *Worker) doAddNetProxyOuterList(netProxyList []*NetProxyOuter) (err error) {
	if len(netProxyList) == 0 {
		return
	}
	this_.netProxyOuterListLock.Lock()
	defer this_.netProxyOuterListLock.Unlock()

	for _, netProxy := range netProxyList {

		var find = this_.findOuterNetProxy(netProxy.Id)
		if find == nil {
			Logger.Info(this_.server.GetServerInfo()+" doAddNetProxyOuterList ", zap.Any("netProxy", netProxy))
			this_.netProxyOuterList = append(this_.netProxyOuterList, netProxy)

			if netProxy.IsEnabled() {
				_ = this_.getNetProxyOuterIfAbsentCreate(netProxy, this_)
			}
		} else {

			var hasChange bool
			if netProxy.Enabled != 0 {
				if netProxy.IsEnabled() != find.IsEnabled() {
					hasChange = true
				}
				find.Enabled = netProxy.Enabled
			}
			if netProxy.Type != find.Type {
				hasChange = true
				find.Type = netProxy.Type
			}
			if netProxy.Address != find.Address {
				hasChange = true
				find.Address = netProxy.Address
			}

			if hasChange {
				Logger.Info(this_.server.GetServerInfo()+" 更新网络代理 ", zap.Any("netProxy", netProxy))
				_ = this_.removeNetProxyOuter(netProxy.Id)
				if find.IsEnabled() {
					_ = this_.getNetProxyOuterIfAbsentCreate(netProxy, this_)
				}
			}

		}
	}

	return
}

func (this_ *Worker) doRemoveNetProxyInnerList(removeNetProxyIdList []string) (err error) {
	if len(removeNetProxyIdList) == 0 {
		return
	}
	this_.netProxyInnerListLock.Lock()
	defer this_.netProxyInnerListLock.Unlock()
	Logger.Info(this_.server.GetServerInfo()+" 移除网络代理 ", zap.Any("removeNetProxyIdList", removeNetProxyIdList))

	for _, netProxyId := range removeNetProxyIdList {
		_ = this_.removeNetProxyInner(netProxyId)
	}

	var list = this_.netProxyInnerList
	var newList []*NetProxyInner
	for _, one := range list {
		if util.ContainsString(removeNetProxyIdList, one.Id) >= 0 {
		} else {
			newList = append(newList, one)
		}
	}
	this_.netProxyInnerList = newList

	return
}

func (this_ *Worker) doRemoveNetProxyOuterList(removeNetProxyIdList []string) (err error) {
	if len(removeNetProxyIdList) == 0 {
		return
	}
	this_.netProxyOuterListLock.Lock()
	defer this_.netProxyOuterListLock.Unlock()
	Logger.Info(this_.server.GetServerInfo()+" 移除网络代理 ", zap.Any("removeNetProxyIdList", removeNetProxyIdList))

	for _, netProxyId := range removeNetProxyIdList {
		_ = this_.removeNetProxyOuter(netProxyId)
	}

	var list = this_.netProxyOuterList
	var newList []*NetProxyOuter
	for _, one := range list {
		if util.ContainsString(removeNetProxyIdList, one.Id) >= 0 {
		} else {
			newList = append(newList, one)
		}
	}
	this_.netProxyOuterList = newList
	return
}
