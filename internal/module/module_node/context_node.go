package module_node

import (
	"go.uber.org/zap"
	"teamide/pkg/node"
	"teamide/pkg/system"
	"teamide/pkg/util"
)

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

	var list = this_.nodeModelIdList
	var newList []int64
	for _, one := range list {
		if one != id {
			newList = append(newList, one)
		}
	}
	this_.nodeModelIdList = newList
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

func (this_ *NodeContext) addNodeModel(nodeModel *NodeModel) {
	if nodeModel == nil {
		return
	}

	this_.Logger.Info("addNodeModel", zap.Any("serverId", nodeModel.NodeId))
	this_.setNodeModel(nodeModel.NodeId, nodeModel)
	this_.setNodeModelByServerId(nodeModel.ServerId, nodeModel)

	var list = this_.nodeModelIdList
	var newList []int64
	var find bool
	for _, one := range list {
		if one == nodeModel.NodeId {
			find = true
			newList = append(newList, one)
		} else {
			newList = append(newList, one)
		}
	}
	if !find {
		newList = append(newList, nodeModel.NodeId)
	}

	this_.nodeModelIdList = newList
}
func (this_ *NodeContext) onAddNodeModel(nodeModel *NodeModel) {
	if nodeModel == nil {
		return
	}
	this_.addNodeModel(nodeModel)
	var err error
	if nodeModel.IsROOT() {
		this_.Logger.Info("onAddNodeModel init root")
		err = this_.initRoot(nodeModel)
		if err != nil {
			this_.Logger.Error("node context init root error", zap.Error(err))
		}
	}
	this_.cleanNodeLine()
	this_.toAddNodeModel(nodeModel)
	this_.doAlive()
}

func (this_ *NodeContext) onUpdateNodeModel(nodeModel *NodeModel) {
	if nodeModel == nil {
		return
	}
	this_.addNodeModel(nodeModel)
	var err error
	if nodeModel.IsROOT() {
		err = this_.initRoot(nodeModel)
		if err != nil {
			this_.Logger.Error("node context init root error", zap.Error(err))
		}
	}

	this_.cleanNodeLine()
	this_.toAddNodeModel(nodeModel)
	this_.doAlive()
}

func (this_ *NodeContext) onUpdateNodeConnServerIds(nodeId int64, connServerIds string) {
	nodeModel := this_.getNodeModel(nodeId)
	if nodeModel == nil {
		return
	}
	var oldConnServerIdList = nodeModel.ConnServerIdList
	nodeModel.ConnServerIds = connServerIds
	nodeModel.ConnServerIdList = GetStringList(connServerIds)

	var newConnServerIdList = nodeModel.ConnServerIdList
	for _, oldConnServerId := range oldConnServerIdList {
		if util.ContainsString(newConnServerIdList, oldConnServerId) < 0 {
			lineNodeIdList := this_.GetNodeLineTo(nodeModel.ServerId)
			_ = this_.server.RemoveToNodeList(lineNodeIdList, []string{
				oldConnServerId,
			})
		}
	}

	var toNodeList []*node.ToNode
	for _, newConnServerId := range newConnServerIdList {
		var toNodeModel = this_.getNodeModelByServerId(newConnServerId)
		if toNodeModel == nil {
			continue
		}
		toNodeList = append(toNodeList, &node.ToNode{
			Id:          nodeModel.ServerId,
			ConnAddress: nodeModel.ConnAddress,
			ConnToken:   nodeModel.ConnToken,
			Enabled:     nodeModel.Enabled,
		})
	}
	lineNodeIdList := this_.GetNodeLineTo(nodeModel.ServerId)
	_ = this_.server.AddToNodeList(lineNodeIdList, toNodeList)
	this_.doAlive()
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

	var list = this_.nodeModelIdList
	for _, one := range list {
		var find = this_.getNodeModel(one)
		if find == nil {
			continue
		}
		if util.ContainsString(find.ConnServerIdList, nodeModel.ServerId) < 0 {
			continue
		}
		var connServerIdList []string
		for _, serverId := range find.ConnServerIdList {
			if serverId != nodeModel.ServerId {
				connServerIdList = append(connServerIdList, serverId)
			}
		}
		find.ConnServerIds = GetListToString(connServerIdList)
		find.ConnServerIdList = connServerIdList

		lineNodeIdList := this_.GetNodeLineTo(find.ServerId)

		_ = this_.server.RemoveToNodeList(lineNodeIdList, []string{
			nodeModel.ServerId,
		})
	}
	this_.removeNodeModel(nodeModel.NodeId)
	this_.removeNodeModelByServerId(nodeModel.ServerId)
	this_.cleanNodeLine()
	this_.doAlive()

}

func (this_ *NodeContext) onEnableNodeModel(id int64) {
	var nodeModel = this_.getNodeModel(id)
	if nodeModel == nil {
		return
	}
	nodeModel.Enabled = 1

	this_.toAddNodeModel(nodeModel)
	this_.doAlive()
}

func (this_ *NodeContext) onDisableNodeModel(id int64) {
	var nodeModel = this_.getNodeModel(id)
	if nodeModel == nil {
		return
	}
	nodeModel.Enabled = 2

	this_.toAddNodeModel(nodeModel)
	this_.doAlive()
}

func (this_ *NodeContext) toAddNodeModel(nodeModel *NodeModel) {
	if nodeModel == nil {
		return
	}
	//this_.Logger.Info("toAddNodeModel", zap.Any("serverId", nodeModel.ServerId))
	var list = this_.nodeModelIdList
	for _, one := range list {
		var find = this_.getNodeModel(one)
		if find == nil {
			continue
		}
		if util.ContainsString(find.ConnServerIdList, nodeModel.ServerId) < 0 {
			continue
		}

		lineNodeIdList := this_.GetNodeLineTo(find.ServerId)
		_ = this_.server.AddToNodeList(lineNodeIdList, []*node.ToNode{
			{
				Id:          nodeModel.ServerId,
				ConnAddress: nodeModel.ConnAddress,
				ConnToken:   nodeModel.ConnToken,
				Enabled:     nodeModel.Enabled,
			},
		})
	}
}

func (this_ *NodeContext) SystemGetInfo(nodeId string) (info *system.Info) {
	lineNodeIdList := this_.GetNodeLineTo(nodeId)
	return this_.server.SystemGetInfo(lineNodeIdList)
}

func (this_ *NodeContext) SystemQueryMonitorData(nodeId string, request *system.QueryRequest) (info *system.QueryResponse) {
	lineNodeIdList := this_.GetNodeLineTo(nodeId)
	return this_.server.SystemQueryMonitorData(lineNodeIdList, request)
}

func (this_ *NodeContext) SystemCleanMonitorData(nodeId string) {
	lineNodeIdList := this_.GetNodeLineTo(nodeId)
	this_.server.SystemCleanMonitorData(lineNodeIdList)
}
