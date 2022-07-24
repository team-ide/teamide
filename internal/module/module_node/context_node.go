package module_node

import (
	"encoding/json"
	"go.uber.org/zap"
	"teamide/node"
	"teamide/pkg/util"
)

type NodeInfo struct {
	Info      *node.Info `json:"info,omitempty"`
	Model     *NodeModel `json:"model,omitempty"`
	IsStarted bool       `json:"isStarted,omitempty"`
}

func (this_ *NodeContext) getNodeInfo(id string) (res *NodeInfo) {
	this_.nodeListLock.Lock()
	defer this_.nodeListLock.Unlock()

	var list = this_.nodeList
	for _, one := range list {
		if one.Info != nil && id == one.Info.Id {
			res = one
			return
		}
	}
	return
}

func (this_ *NodeContext) getNodeModel(id int64) (res *NodeModel) {
	this_.nodeIdModelCacheLock.Lock()
	defer this_.nodeIdModelCacheLock.Unlock()

	res = this_.nodeIdModelCache[id]
	return
}

func (this_ *NodeContext) setNodeModel(id int64, nodeModel *NodeModel) {
	this_.nodeIdModelCacheLock.Lock()
	defer this_.nodeIdModelCacheLock.Unlock()

	this_.nodeIdModelCache[id] = nodeModel
}

func (this_ *NodeContext) removeNodeModel(id int64) {
	this_.nodeIdModelCacheLock.Lock()
	defer this_.nodeIdModelCacheLock.Unlock()
	delete(this_.nodeIdModelCache, id)
}

func (this_ *NodeContext) getNodeModelByServerId(id string) (res *NodeModel) {
	this_.serverIdModelCacheLock.Lock()
	defer this_.serverIdModelCacheLock.Unlock()
	res = this_.serverIdModelCache[id]
	return
}

func (this_ *NodeContext) setNodeModelByServerId(id string, nodeModel *NodeModel) {
	this_.serverIdModelCacheLock.Lock()
	defer this_.serverIdModelCacheLock.Unlock()

	this_.serverIdModelCache[id] = nodeModel
}

func (this_ *NodeContext) removeNodeModelByServerId(id string) {
	this_.serverIdModelCacheLock.Lock()
	defer this_.serverIdModelCacheLock.Unlock()
	delete(this_.serverIdModelCache, id)
}

func (this_ *NodeContext) onAddNodeModel(nodeModel *NodeModel) {
	if nodeModel == nil {
		return
	}
	this_.setNodeModel(nodeModel.NodeId, nodeModel)
	this_.setNodeModelByServerId(nodeModel.ServerId, nodeModel)
	var err error
	if nodeModel.IsROOT() {
		err = this_.initRoot(nodeModel)
		if err != nil {
			this_.nodeService.Logger.Error("node context init root error", zap.Error(err))
		}
	}
	var connNodeIdList []string
	if nodeModel.ConnServerIds != "" {
		_ = json.Unmarshal([]byte(nodeModel.ConnServerIds), &connNodeIdList)
	}

	_ = this_.server.AddNodeList([]*node.Info{
		{
			Id:             nodeModel.ServerId,
			BindToken:      nodeModel.BindToken,
			BindAddress:    nodeModel.BindAddress,
			ConnAddress:    nodeModel.ConnAddress,
			ConnToken:      nodeModel.ConnToken,
			ConnNodeIdList: connNodeIdList,
		},
	})
}

func (this_ *NodeContext) onUpdateNodeModel(nodeModel *NodeModel) {
	if nodeModel == nil {
		return
	}
	var err error
	if nodeModel.IsROOT() {
		err = this_.initRoot(nodeModel)
		if err != nil {
			this_.nodeService.Logger.Error("node context init root error", zap.Error(err))
		}
	}
	this_.setNodeModel(nodeModel.NodeId, nodeModel)
	this_.setNodeModelByServerId(nodeModel.ServerId, nodeModel)
}

func (this_ *NodeContext) onUpdateNodeConnServerIds(nodeId int64, connServerIds string) {
	nodeModel := this_.getNodeModel(nodeId)
	if nodeModel == nil {
		return
	}
	nodeModel.ConnServerIds = connServerIds

	var connNodeIdList []string
	if connServerIds != "" {
		_ = json.Unmarshal([]byte(connServerIds), &connNodeIdList)
	}
	_ = this_.server.UpdateNodeConnNodeIdList(nodeModel.ServerId, connNodeIdList)
}

func (this_ *NodeContext) onUpdateNodeHistoryConnServerIds(nodeId int64, historyConnServerIds string) {
	nodeModel := this_.getNodeModel(nodeId)
	if nodeModel == nil {
		return
	}
	nodeModel.HistoryConnServerIds = historyConnServerIds

}

func (this_ *NodeContext) onRemoveNodeModel(id int64) {
	var nodeModel = this_.getNodeModel(id)
	if nodeModel == nil {
		return
	}
	this_.removeNodeModel(nodeModel.NodeId)
	this_.removeNodeModelByServerId(nodeModel.ServerId)
	_ = this_.server.RemoveNodeList([]string{nodeModel.ServerId})
}

func (this_ *NodeContext) onEnableNodeModel(id int64) {
	var nodeModel = this_.getNodeModel(id)
	if nodeModel == nil {
		return
	}
	nodeModel.Enabled = 1
	this_.setNodeModel(nodeModel.NodeId, nodeModel)
	this_.setNodeModelByServerId(nodeModel.ServerId, nodeModel)
	_ = this_.server.AddNodeList([]*node.Info{
		{
			Id:          nodeModel.ServerId,
			BindToken:   nodeModel.BindToken,
			BindAddress: nodeModel.BindAddress,
			ConnAddress: nodeModel.ConnAddress,
			ConnToken:   nodeModel.ConnToken,
			Enabled:     nodeModel.Enabled,
		},
	})
}

func (this_ *NodeContext) onDisableNodeModel(id int64) {
	var nodeModel = this_.getNodeModel(id)
	if nodeModel == nil {
		return
	}
	nodeModel.Enabled = 2
	this_.setNodeModel(nodeModel.NodeId, nodeModel)
	this_.setNodeModelByServerId(nodeModel.ServerId, nodeModel)
	_ = this_.server.AddNodeList([]*node.Info{
		{
			Id:          nodeModel.ServerId,
			BindToken:   nodeModel.BindToken,
			BindAddress: nodeModel.BindAddress,
			ConnAddress: nodeModel.ConnAddress,
			ConnToken:   nodeModel.ConnToken,
			Enabled:     nodeModel.Enabled,
		},
	})
}

func (this_ *NodeContext) onNodeListChange(nodeList []*node.Info) {

	this_.nodeListLock.Lock()
	defer this_.nodeListLock.Unlock()

	var nodeInfoList []*NodeInfo
	for _, one := range nodeList {
		var find = this_.getNodeModelByServerId(one.Id)
		var historyConnServerIdList []string
		if find != nil && find.HistoryConnServerIds != "" {
			_ = json.Unmarshal([]byte(find.HistoryConnServerIds), &historyConnServerIdList)
		}
		for _, one := range one.ConnNodeIdList {
			if util.ContainsString(historyConnServerIdList, one) < 0 {
				historyConnServerIdList = append(historyConnServerIdList, one)
			}
		}
		var historyConnServerIds string
		bs, _ := json.Marshal(historyConnServerIdList)
		if bs != nil {
			historyConnServerIds = string(bs)
		}
		if find == nil {
			find = &NodeModel{
				ServerId:             one.Id,
				Name:                 one.Id,
				BindToken:            one.BindToken,
				BindAddress:          one.BindAddress,
				ConnToken:            one.ConnToken,
				ConnAddress:          one.ConnAddress,
				HistoryConnServerIds: historyConnServerIds,
			}
			_, err := this_.nodeService.Insert(find)
			if err != nil {
				find = nil
			} else {
				find, _ = this_.nodeService.Get(find.NodeId)
				if find != nil {
					this_.setNodeModel(find.NodeId, find)
					this_.setNodeModelByServerId(one.Id, find)
				}
			}
		} else {
			_, _ = this_.nodeService.UpdateHistoryConnServerIds(find.NodeId, historyConnServerIds)
		}
		nodeInfo := &NodeInfo{
			Info:      one,
			IsStarted: one.Status == node.StatusStarted,
			Model:     this_.getNodeModelByServerId(one.Id),
		}
		nodeInfoList = append(nodeInfoList, nodeInfo)
	}
	this_.nodeList = nodeInfoList
	this_.callMessage(&Message{
		Method:   "refresh_node_list",
		NodeList: this_.nodeList,
	})
}
