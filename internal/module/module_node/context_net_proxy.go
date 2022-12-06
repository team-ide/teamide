package module_node

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"teamide/pkg/node"
	"teamide/pkg/util"
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
}

func (this_ *NodeContext) removeNetProxyModel(id int64) {
	this_.netProxyIdModelCacheLock.Lock()
	defer this_.netProxyIdModelCacheLock.Unlock()
	delete(this_.netProxyIdModelCache, id)

	var list = this_.netProxyModelIdList
	var newList []int64
	for _, one := range list {
		if one != id {
			newList = append(newList, one)
		}
	}
	this_.netProxyModelIdList = newList
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

func (this_ *NodeContext) addNetProxyModel(netProxyModel *NetProxyModel) {
	if netProxyModel == nil {
		return
	}
	this_.setNetProxyModel(netProxyModel.NetProxyId, netProxyModel)
	this_.setNetProxyModelByCode(netProxyModel.Code, netProxyModel)

	var list = this_.netProxyModelIdList
	if util.ContainsInt64(list, netProxyModel.NetProxyId) < 0 {
		list = append(list, netProxyModel.NetProxyId)
	}
	this_.netProxyModelIdList = list

}

func (this_ *NodeContext) onAddNetProxyModel(netProxyModel *NetProxyModel) {
	if netProxyModel == nil {
		return
	}
	this_.addNetProxyModel(netProxyModel)

	this_.toAddNetProxyModel(netProxyModel)
	this_.doAlive()

}

func (this_ *NodeContext) toAddNetProxyModel(netProxyModel *NetProxyModel) {
	if netProxyModel == nil {
		return
	}
	err := this_.formatNetProxy(netProxyModel)
	if err != nil {
		this_.Logger.Error("toAddNetProxyModel formatNetProxy error", zap.Error(err))
		return
	}
	lineNodeIdList := this_.GetNodeLineTo(netProxyModel.InnerServerId)
	if len(lineNodeIdList) > 0 {
		err = this_.server.AddNetProxyInnerList(lineNodeIdList, []*node.NetProxyInner{
			{
				Id:             netProxyModel.Code,
				NodeId:         netProxyModel.InnerServerId,
				Type:           netProxyModel.InnerType,
				Address:        netProxyModel.InnerAddress,
				Enabled:        netProxyModel.Enabled,
				LineNodeIdList: netProxyModel.LineNodeIdList,
			},
		})
		if err != nil {
			this_.Logger.Error("toAddNetProxyModel AddNetProxyInnerList error", zap.Error(err))
		}
	}
	lineNodeIdList = this_.GetNodeLineTo(netProxyModel.OuterServerId)
	if len(lineNodeIdList) > 0 {
		err = this_.server.AddNetProxyOuterList(lineNodeIdList, []*node.NetProxyOuter{
			{
				Id:                    netProxyModel.Code,
				NodeId:                netProxyModel.OuterServerId,
				Type:                  netProxyModel.OuterType,
				Address:               netProxyModel.OuterAddress,
				Enabled:               netProxyModel.Enabled,
				ReverseLineNodeIdList: netProxyModel.ReverseLineNodeIdList,
			},
		})
		if err != nil {
			this_.Logger.Error("toAddNetProxyModel AddNetProxyOuterList error", zap.Error(err))
		}
	}
}

func (this_ *NodeContext) onUpdateNetProxyModel(netProxyModel *NetProxyModel) {
	if netProxyModel == nil {
		return
	}
	this_.setNetProxyModel(netProxyModel.NetProxyId, netProxyModel)
	this_.setNetProxyModelByCode(netProxyModel.Code, netProxyModel)

	this_.toAddNetProxyModel(netProxyModel)
	this_.doAlive()
}

func (this_ *NodeContext) onRemoveNetProxyModel(id int64) {
	var netProxyModel = this_.getNetProxyModel(id)
	if netProxyModel == nil {
		return
	}
	this_.removeNetProxyModel(netProxyModel.NetProxyId)
	this_.removeNetProxyModelByCode(netProxyModel.Code)

	this_.toRemoveNetProxyModel(netProxyModel)
	this_.doAlive()
}

func (this_ *NodeContext) toRemoveNetProxyModel(netProxyModel *NetProxyModel) {
	if netProxyModel == nil {
		return
	}

	lineNodeIdList := this_.GetNodeLineTo(netProxyModel.InnerServerId)
	if len(lineNodeIdList) > 0 {
		_ = this_.server.RemoveNetProxyInnerList(lineNodeIdList, []string{
			netProxyModel.Code,
		})
	}
	lineNodeIdList = this_.GetNodeLineTo(netProxyModel.OuterServerId)
	if len(lineNodeIdList) > 0 {
		_ = this_.server.RemoveNetProxyOuterList(lineNodeIdList, []string{
			netProxyModel.Code,
		})
	}
}

func (this_ *NodeContext) onEnableNetProxyModel(id int64) {
	var netProxyModel = this_.getNetProxyModel(id)
	if netProxyModel == nil {
		return
	}
	netProxyModel.Enabled = 1

	this_.toAddNetProxyModel(netProxyModel)
	this_.doAlive()
}

func (this_ *NodeContext) onDisableNetProxyModel(id int64) {
	var netProxyModel = this_.getNetProxyModel(id)
	if netProxyModel == nil {
		return
	}
	netProxyModel.Enabled = 2

	this_.toAddNetProxyModel(netProxyModel)
	this_.doAlive()
}
