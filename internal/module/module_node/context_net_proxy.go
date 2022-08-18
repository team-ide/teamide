package module_node

import (
	"encoding/json"
	"errors"
	"teamide/node"
)

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

	var find bool
	var list = this_.netProxyModelList
	for _, one := range list {
		if one.NetProxyId == id {
			find = true
		}
	}
	if find {
		this_.netProxyModelList = append(this_.netProxyModelList, netProxyModel)
	}
}

func (this_ *NodeContext) removeNetProxyModel(id int64) {
	this_.netProxyIdModelCacheLock.Lock()
	defer this_.netProxyIdModelCacheLock.Unlock()
	delete(this_.netProxyIdModelCache, id)

	var list = this_.netProxyModelList
	var newList []*NetProxyModel
	for _, one := range list {
		if one.NetProxyId != id {
			newList = append(newList, one)
		}
	}
	this_.netProxyModelList = newList
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

func (this_ *NodeContext) formatNetProxy(netProxyModel *NetProxyModel) (err error) {

	netProxyModel.LineNodeIdList = GetStringList(netProxyModel.LineServerIds)
	if netProxyModel.Code == "" {
		err = errors.New("网络代理编号不能为空")
		return
	}
	if netProxyModel.InnerAddress == "" {
		err = errors.New("网络代理输入地址不能为空")
		return
	}
	if netProxyModel.OuterAddress == "" {
		err = errors.New("网络代理输出地址不能为空")
		return
	}
	if len(netProxyModel.LineNodeIdList) == 0 {
		netProxyModel.LineNodeIdList = this_.GetNodeLineByFromTo(netProxyModel.InnerServerId, netProxyModel.OuterServerId)
		if len(netProxyModel.LineNodeIdList) == 0 {
			err = errors.New("无法正确解析输入输出节点关系")
			return
		}
	}
	netProxyModel.ReverseLineNodeIdList = []string{}
	for i := len(netProxyModel.LineNodeIdList) - 1; i >= 0; i-- {
		netProxyModel.ReverseLineNodeIdList = append(netProxyModel.ReverseLineNodeIdList, netProxyModel.LineNodeIdList[i])
	}

	bs, _ := json.Marshal(&netProxyModel.LineNodeIdList)
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

	err := this_.formatNetProxy(netProxyModel)
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

	err := this_.formatNetProxy(netProxyModel)
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

	err := this_.formatNetProxy(netProxyModel)
	if err == nil {
		_ = this_.server.AddNetProxyList([]*node.NetProxy{
			one,
		})
	}
}
