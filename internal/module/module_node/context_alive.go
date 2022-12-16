package module_node

import (
	"encoding/json"
	"go.uber.org/zap"
	"teamide/pkg/node"
	"time"
)

func (this_ *NodeContext) doAlive() {
	if this_.doAliveIng {
		return
	}

	this_.doAliveLock.Lock()
	defer this_.doAliveLock.Unlock()

	this_.doAliveIng = true
	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("doAlive error", zap.Any("error", e))
		}
		this_.doAliveIng = false
		go func() {
			time.Sleep(time.Second * 5)
			this_.doAlive()
		}()
	}()

	if len(this_.localNodeList) == 0 || len(this_.serverList) == 0 {
		return
	}
	//this_.Logger.Info("node do alive")

	var nodeModelIdList = this_.nodeModelIdList

	for _, id := range nodeModelIdList {
		find := this_.getNodeModel(id)
		if find == nil {
			continue
		}
		this_.toAddNodeModel(find)
	}

	var netProxyModelIdList = this_.netProxyModelIdList
	for _, id := range netProxyModelIdList {
		find := this_.getNetProxyModel(id)
		if find == nil {
			continue
		}
		this_.toAddNetProxyModel(find)
	}

	var countData = newNodeCountData()
	this_.doCountData(countData)
	this_.countData = countData

	this_.checkChangeOut()
	return
}

func (this_ *NodeContext) doCountData(countData *NodeCountData) {
	countData.NodeCount = 0
	countData.NodeSuccessCount = 0
	var nodeModelIdList = this_.nodeModelIdList
	for _, id := range nodeModelIdList {
		find := this_.getNodeModel(id)
		if find == nil {
			continue
		}
		countData.NodeCount++

		localServerId, lineNodeIdList := this_.GetNodeLineTo(find.ServerId)
		server := this_.GetServer(localServerId)
		if len(lineNodeIdList) > 0 && server != nil {
			//this_.Logger.Info("toAddNodeModel", zap.Any("to node", nodeModel.ServerId), zap.Any("lineNodeIdList", lineNodeIdList))
			status := server.GetNodeStatus(lineNodeIdList)
			find.Status = status
			find.IsStarted = status == node.StatusStarted
			if find.IsStarted {
				countData.NodeSuccessCount++
			}
		} else {
			find.Status = 0
		}
	}
	countData.NodeNetProxyCount = 0
	countData.NodeNetProxyInnerSuccessCount = 0
	countData.NodeNetProxyOuterSuccessCount = 0

	var netProxyModelIdList = this_.netProxyModelIdList
	for _, id := range netProxyModelIdList {
		find := this_.getNetProxyModel(id)
		if find == nil {
			continue
		}

		localServerId, lineNodeIdList := this_.GetNodeLineTo(find.InnerServerId)
		server := this_.GetServer(localServerId)
		if len(lineNodeIdList) > 0 && server != nil {
			status := server.GetNetProxyInnerStatus(lineNodeIdList, find.Code)

			find.InnerStatus = status
			find.InnerIsStarted = status == node.StatusStarted
			if find.InnerIsStarted {
				countData.NodeNetProxyInnerSuccessCount++
			}
		} else {
			find.InnerStatus = 0
		}

		localServerId, lineNodeIdList = this_.GetNodeLineTo(find.OuterServerId)
		server = this_.GetServer(localServerId)
		if len(lineNodeIdList) > 0 && server != nil {
			status := server.GetNetProxyOuterStatus(lineNodeIdList, find.Code)
			find.OuterStatus = status
			find.OuterIsStarted = status == node.StatusStarted
			if find.OuterIsStarted {
				countData.NodeNetProxyOuterSuccessCount++
			}
		} else {
			find.OuterStatus = 0
		}
	}

	return
}

func (this_ *NodeContext) getNodeModelList() []*NodeModel {
	var nodeModelList []*NodeModel

	var nodeModelIdList = this_.nodeModelIdList

	for _, id := range nodeModelIdList {
		find := this_.getNodeModel(id)
		if find == nil {
			continue
		}
		nodeModelList = append(nodeModelList, find)
	}
	return nodeModelList
}

func (this_ *NodeContext) getUserNodeModelList(userId int64) []*NodeModel {
	var nodeModelList []*NodeModel

	var list = this_.getNodeModelList()

	for _, one := range list {
		if one.UserId != userId {
			continue
		}
		nodeModelList = append(nodeModelList, one)
	}
	return nodeModelList
}

func (this_ *NodeContext) getNetProxyModelList() []*NetProxyModel {
	var netProxyModelList []*NetProxyModel

	var netProxyModelIdList = this_.netProxyModelIdList
	for _, id := range netProxyModelIdList {
		find := this_.getNetProxyModel(id)
		if find == nil {
			continue
		}
		netProxyModelList = append(netProxyModelList, find)
	}
	return netProxyModelList
}

func (this_ *NodeContext) getUserNetProxyModelList(userId int64) []*NetProxyModel {
	var netProxyModelList []*NetProxyModel

	var list = this_.getNetProxyModelList()

	for _, one := range list {
		if one.UserId != userId {
			continue
		}
		netProxyModelList = append(netProxyModelList, one)
	}
	return netProxyModelList
}

func (this_ *NodeContext) checkChangeOut() {
	var oldNodeListStr = this_.oldNodeListStr
	var oldNetProxyListStr = this_.oldNetProxyListStr
	var oldCountDataStr = this_.oldCountDataStr

	var nodeModelList = this_.getNodeModelList()
	var netProxyModelList = this_.getNetProxyModelList()

	var countData = this_.countData

	newBs, _ := json.Marshal(countData)
	newCountDataStr := string(newBs)
	//this_.Logger.Info("count data validate", zap.Any("old", string(oldCountBs)), zap.Any("new", string(newBs)))
	if oldCountDataStr != newCountDataStr {
		//this_.Logger.Info("data count change", zap.Any("old", oldCountDataStr), zap.Any("new", newCountDataStr))
		this_.oldCountDataStr = newCountDataStr
		this_.countData = countData
		this_.callNodeCountDataChange(countData)
	}

	newBs, _ = json.Marshal(nodeModelList)
	newNodeListStr := string(newBs)
	//this_.Logger.Info("node list validate", zap.Any("old", string(oldNodeBs)), zap.Any("new", string(newBs)))
	if oldNodeListStr != newNodeListStr {
		//this_.Logger.Info("node list change", zap.Any("old", oldNodeListStr), zap.Any("new", newNodeListStr))
		this_.oldNodeListStr = newNodeListStr
		this_.callNodeListChange(nodeModelList)
	}

	newBs, _ = json.Marshal(netProxyModelList)
	newNetProxyListStr := string(newBs)
	if oldNetProxyListStr != newNetProxyListStr {
		//this_.Logger.Info("net proxy list change", zap.Any("old", oldNetProxyListStr), zap.Any("new", newNetProxyListStr))
		this_.oldNetProxyListStr = newNetProxyListStr
		this_.callNetProxyListChange(netProxyModelList)
	}
	return
}
