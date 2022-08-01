package module_node

import (
	"encoding/json"
	"errors"
	"teamide/node"
	"teamide/pkg/util"
)

type NetProxyInfo struct {
	Info             *node.NetProxy     `json:"info,omitempty"`
	Model            *NetProxyModel     `json:"model,omitempty"`
	InnerIsStarted   bool               `json:"innerIsStarted,omitempty"`
	OuterIsStarted   bool               `json:"outerIsStarted,omitempty"`
	InnerMonitorData *MonitorDataFormat `json:"innerMonitorData,omitempty"`
	OuterMonitorData *MonitorDataFormat `json:"outerMonitorData,omitempty"`
}

func (this_ *NodeContext) getNetProxyInfo(id string) (res *NetProxyInfo) {
	this_.netProxyListLock.Lock()
	defer this_.netProxyListLock.Unlock()

	var list = this_.netProxyList
	for _, one := range list {
		if one.Info != nil && id == one.Info.Id {
			res = one
			return
		}
	}
	return
}

func (this_ *NodeContext) getNetProxyModel(id int64) (res *NetProxyModel) {
	this_.netProxyIdModelCacheLock.Lock()
	defer this_.netProxyIdModelCacheLock.Unlock()

	res = this_.netProxyIdModelCache[id]
	return
}

func (this_ *NodeContext) setNetProxyModel(id int64, netProxyModel *NetProxyModel) {
	this_.netProxyIdModelCacheLock.Lock()
	defer this_.netProxyIdModelCacheLock.Unlock()

	this_.netProxyIdModelCache[id] = netProxyModel
}

func (this_ *NodeContext) removeNetProxyModel(id int64) {
	this_.netProxyIdModelCacheLock.Lock()
	defer this_.netProxyIdModelCacheLock.Unlock()
	delete(this_.netProxyIdModelCache, id)
}

func (this_ *NodeContext) getNetProxyModelByCode(id string) (res *NetProxyModel) {
	this_.codeModelCacheLock.Lock()
	defer this_.codeModelCacheLock.Unlock()
	res = this_.codeModelCache[id]
	return
}

func (this_ *NodeContext) setNetProxyModelByCode(id string, netProxyModel *NetProxyModel) {
	this_.codeModelCacheLock.Lock()
	defer this_.codeModelCacheLock.Unlock()

	this_.codeModelCache[id] = netProxyModel
}

func (this_ *NodeContext) removeNetProxyModelByCode(id string) {
	this_.codeModelCacheLock.Lock()
	defer this_.codeModelCacheLock.Unlock()
	delete(this_.codeModelCache, id)
}

func (this_ *NodeContext) formatNetProxy(netProxyModel *NetProxyModel) (netProxy *node.NetProxy, err error) {
	netProxy = &node.NetProxy{
		Id: netProxyModel.Code,
		Inner: &node.NetConfig{
			NodeId:  netProxyModel.InnerServerId,
			Type:    netProxyModel.InnerType,
			Address: netProxyModel.InnerAddress,
		},
		Outer: &node.NetConfig{
			NodeId:  netProxyModel.OuterServerId,
			Type:    netProxyModel.OuterType,
			Address: netProxyModel.OuterAddress,
		},
		Enabled: netProxyModel.Enabled,
	}
	if netProxyModel.LineServerIds != "" {
		_ = json.Unmarshal([]byte(netProxyModel.LineServerIds), &netProxy.LineNodeIdList)
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
		var nodeIdConnNodeIdListCache = make(map[string][]string)
		var list = this_.nodeList
		for _, one := range list {
			var id string
			if one.Info != nil {
				id = one.Info.Id
			} else if one.Model != nil {
				id = one.Model.ServerId
			}
			var connNodeIdList []string
			if one.Info != nil {
				connNodeIdList = append(connNodeIdList, one.Info.ConnNodeIdList...)
			}
			if one.Model != nil {
				if one.Model.ConnServerIds != "" {
					var connServerIdList []string
					_ = json.Unmarshal([]byte(one.Model.ConnServerIds), &connServerIdList)
					for _, connNodeId := range connServerIdList {
						if util.ContainsString(connNodeIdList, connNodeId) < 0 {
							connNodeIdList = append(connNodeIdList, connNodeId)
						}
					}
				}
				if one.Model.HistoryConnServerIds != "" {
					var historyConnServerIdList []string
					_ = json.Unmarshal([]byte(one.Model.HistoryConnServerIds), &historyConnServerIdList)
					for _, connNodeId := range historyConnServerIdList {
						if util.ContainsString(connNodeIdList, connNodeId) < 0 {
							connNodeIdList = append(connNodeIdList, connNodeId)
						}
					}
				}
			}

			nodeIdConnNodeIdListCache[id] = connNodeIdList
		}
		netProxy.LineNodeIdList = this_.server.GetNodeLineByFromTo(netProxy.Inner.NodeId, netProxy.Outer.NodeId, nodeIdConnNodeIdListCache)
		if len(netProxy.LineNodeIdList) == 0 {
			err = errors.New("无法正确解析输入输出节点关系")
			return
		}
	}
	netProxy.ReverseLineNodeIdList = []string{}
	for i := len(netProxy.LineNodeIdList) - 1; i >= 0; i-- {
		netProxy.ReverseLineNodeIdList = append(netProxy.ReverseLineNodeIdList, netProxy.LineNodeIdList[i])
	}

	bs, _ := json.Marshal(&netProxy.LineNodeIdList)
	if len(bs) > 0 {
		netProxyModel.LineServerIds = string(bs)
	}
	return
}

func (this_ *NodeContext) onAddNetProxyModel(netProxyModel *NetProxyModel) {
	if netProxyModel == nil {
		return
	}
	this_.setNetProxyModel(netProxyModel.NetProxyId, netProxyModel)
	this_.setNetProxyModelByCode(netProxyModel.Code, netProxyModel)

	one, err := this_.formatNetProxy(netProxyModel)
	if err == nil {
		_ = this_.server.AddNetProxyList([]*node.NetProxy{
			one,
		})
	}
}

func (this_ *NodeContext) onUpdateNetProxyModel(netProxyModel *NetProxyModel) {
	if netProxyModel == nil {
		return
	}
	this_.setNetProxyModel(netProxyModel.NetProxyId, netProxyModel)
	this_.setNetProxyModelByCode(netProxyModel.Code, netProxyModel)
}

func (this_ *NodeContext) onRemoveNetProxyModel(id int64) {
	var netProxyModel = this_.getNetProxyModel(id)
	if netProxyModel == nil {
		return
	}
	this_.removeNetProxyModel(netProxyModel.NetProxyId)
	this_.removeNetProxyModelByCode(netProxyModel.Code)
	_ = this_.server.RemoveNetProxyList([]string{netProxyModel.Code})
}

func (this_ *NodeContext) onEnableNetProxyModel(id int64) {
	var netProxyModel = this_.getNetProxyModel(id)
	if netProxyModel == nil {
		return
	}
	netProxyModel.Enabled = 1
	this_.setNetProxyModel(netProxyModel.NetProxyId, netProxyModel)
	this_.setNetProxyModelByCode(netProxyModel.Code, netProxyModel)

	one, err := this_.formatNetProxy(netProxyModel)
	if err == nil {
		_ = this_.server.AddNetProxyList([]*node.NetProxy{
			one,
		})
	}
}

func (this_ *NodeContext) onDisableNetProxyModel(id int64) {
	var netProxyModel = this_.getNetProxyModel(id)
	if netProxyModel == nil {
		return
	}
	netProxyModel.Enabled = 2
	this_.setNetProxyModel(netProxyModel.NetProxyId, netProxyModel)
	this_.setNetProxyModelByCode(netProxyModel.Code, netProxyModel)

	one, err := this_.formatNetProxy(netProxyModel)
	if err == nil {
		_ = this_.server.AddNetProxyList([]*node.NetProxy{
			one,
		})
	}
}

func (this_ *NodeContext) onNetProxyListChange(netProxyList []*node.NetProxy) {

	this_.netProxyListLock.Lock()
	defer this_.netProxyListLock.Unlock()

	this_.onNetProxyListChangeIng = true
	defer func() { this_.onNetProxyListChangeIng = false }()

	var netProxyInfoList []*NetProxyInfo
	for _, one := range netProxyList {
		var find = this_.getNetProxyModelByCode(one.Id)
		if find == nil {
			continue
			//find = &NetProxyModel{
			//	Code:          one.Id,
			//	Name:          one.Id,
			//	InnerServerId: one.Inner.NodeId,
			//	InnerType:     one.Inner.Type,
			//	InnerAddress:  one.Inner.Address,
			//	OuterServerId: one.Outer.NodeId,
			//	OuterType:     one.Outer.Type,
			//	OuterAddress:  one.Outer.Address,
			//}
			//_, err := this_.nodeService.InsertNetProxy(find)
			//if err != nil {
			//	find = nil
			//} else {
			//	this_.setNetProxyModel(find.NetProxyId, find)
			//	this_.setNetProxyModelByCode(one.Id, find)
			//}
		}
		netProxyInfo := &NetProxyInfo{
			Info:             one,
			InnerIsStarted:   one.InnerStatus == node.StatusStarted,
			OuterIsStarted:   one.OuterStatus == node.StatusStarted,
			Model:            this_.getNetProxyModelByCode(one.Id),
			InnerMonitorData: ToMonitorDataFormat(nil),
			OuterMonitorData: ToMonitorDataFormat(nil),
		}
		netProxyInfoList = append(netProxyInfoList, netProxyInfo)
	}
	this_.netProxyList = netProxyInfoList
	this_.refreshNetProxyList(this_.netProxyList)
}

func (this_ *NodeContext) refreshNetProxyList(netProxyList []*NetProxyInfo) {
	this_.callNodeDataChange(&NodeDataChange{
		Type:         "netProxyList",
		NetProxyList: netProxyList,
	})
}
