package module_node

import (
	"encoding/json"
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
		this_.doAliveIng = false
		time.Sleep(time.Second * 5)
		this_.doAlive()
	}()

	if this_.root == nil || this_.server == nil {
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

		lineNodeIdList := this_.GetNodeLineTo(find.ServerId)
		if len(lineNodeIdList) > 0 {
			//this_.Logger.Info("toAddNodeModel", zap.Any("to node", nodeModel.ServerId), zap.Any("lineNodeIdList", lineNodeIdList))
			status := this_.server.GetNodeStatus(lineNodeIdList)
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

		lineNodeIdList := this_.GetNodeLineTo(find.InnerServerId)
		if len(lineNodeIdList) > 0 {
			status := this_.server.GetNetProxyInnerStatus(lineNodeIdList, find.Code)

			find.InnerStatus = status
			find.InnerIsStarted = status == node.StatusStarted
			if find.InnerIsStarted {
				countData.NodeNetProxyInnerSuccessCount++
			}
		} else {
			find.InnerStatus = 0
		}

		lineNodeIdList = this_.GetNodeLineTo(find.OuterServerId)
		if len(lineNodeIdList) > 0 {
			status := this_.server.GetNetProxyOuterStatus(lineNodeIdList, find.Code)
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

func (this_ *NodeContext) checkChangeOut() {

	var nodeModelList = this_.getNodeModelList()
	var netProxyModelList = this_.getNetProxyModelList()

	var countData = this_.countData

	newBs, _ := json.Marshal(countData)
	newCountDataStr := string(newBs)
	//this_.Logger.Info("count data validate", zap.Any("old", string(oldCountBs)), zap.Any("new", string(newBs)))
	if this_.oldCountDataStr != newCountDataStr {
		this_.oldCountDataStr = newCountDataStr
		this_.countData = countData
		this_.callNodeCountDataChange(countData)
	}

	newBs, _ = json.Marshal(nodeModelList)
	newNodeListStr := string(newBs)
	//this_.Logger.Info("node list validate", zap.Any("old", string(oldNodeBs)), zap.Any("new", string(newBs)))
	if this_.oldNodeListStr != newNodeListStr {
		this_.oldNodeListStr = newNodeListStr
		this_.callNodeListChange(nodeModelList)
	}

	newBs, _ = json.Marshal(netProxyModelList)
	newNetProxyListStr := string(newBs)
	//this_.Logger.Info("net proxy list validate", zap.Any("old", string(netProxyBs)), zap.Any("new", string(newBs)))
	if this_.oldNetProxyListStr != newNetProxyListStr {
		this_.oldNetProxyListStr = newNetProxyListStr
		this_.callNetProxyListChange(netProxyModelList)
	}
	return
}
