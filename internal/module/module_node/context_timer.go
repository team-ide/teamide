package module_node

import (
	"encoding/json"
	"teamide/node"
	"time"
)

func (this_ *NodeContext) doAlive() {

	defer func() {
		time.Sleep(time.Second * 5)
		this_.doAlive()
	}()

	if this_.root == nil || this_.server == nil {
		return
	}

	var nodeModelList = this_.nodeModelList
	for _, nodeModel := range nodeModelList {
		this_.toAddNodeModel(nodeModel)
	}

	var netProxyModelList = this_.netProxyModelList
	for _, netProxyModel := range netProxyModelList {
		this_.toAddNetProxyModel(netProxyModel)
	}

	var countData = newNodeCountData()
	this_.doCountData(countData)

	oldBs, _ := json.Marshal(countData)
	newBs, _ := json.Marshal(this_.countData)
	if string(oldBs) != string(newBs) {
		this_.countData = countData
		this_.callNodeCountDataChange(countData)
	}

	return
}

func (this_ *NodeContext) doCountData(countData *NodeCountData) {
	countData.NodeCount = 0
	countData.NodeSuccessCount = 0
	var nodeModelList = this_.nodeModelList
	for _, nodeModel := range nodeModelList {
		countData.NodeCount++

		lineNodeIdList := this_.GetNodeLineTo(nodeModel.ServerId)
		if len(lineNodeIdList) > 0 {
			status := this_.server.GetNodeStatus(lineNodeIdList)
			if status == node.StatusStarted {
				countData.NodeSuccessCount++
			}
		}
	}
	countData.NodeNetProxyCount = 0
	countData.NodeNetProxyInnerSuccessCount = 0
	countData.NodeNetProxyOuterSuccessCount = 0

	var netProxyModelList = this_.netProxyModelList
	for _, netProxyModel := range netProxyModelList {

		lineNodeIdList := this_.GetNodeLineTo(netProxyModel.InnerServerId)
		if len(lineNodeIdList) > 0 {
			status := this_.server.GetNetProxyInnerStatus(lineNodeIdList, netProxyModel.Code)
			if status == node.StatusStarted {
				countData.NodeNetProxyInnerSuccessCount++
			}
		}

		lineNodeIdList = this_.GetNodeLineTo(netProxyModel.OuterServerId)
		if len(lineNodeIdList) > 0 {
			status := this_.server.GetNetProxyOuterStatus(lineNodeIdList, netProxyModel.Code)
			if status == node.StatusStarted {
				countData.NodeNetProxyOuterSuccessCount++
			}
		}
	}

	return
}
