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

	if len(this_.localNodeList) == 0 {
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

	for _, id := range nodeModelIdList {
		find := this_.getNodeModel(id)
		if find == nil {
			continue
		}

		lineNodeIdList := this_.GetNodeLineTo(find.ServerId)
		if len(lineNodeIdList) > 0 {
			//this_.Logger.Info("toAddNodeModel", zap.Any("to node", nodeModel.ServerId), zap.Any("lineNodeIdList", lineNodeIdList))
			status := this_.GetServer().GetNodeStatus(lineNodeIdList)
			find.Status = status
			find.IsStarted = status == node.StatusStarted
		} else {
			find.Status = 0
		}
	}

	for _, id := range netProxyModelIdList {
		find := this_.getNetProxyModel(id)
		if find == nil {
			continue
		}

		lineNodeIdList := this_.GetNodeLineTo(find.InnerServerId)
		if len(lineNodeIdList) > 0 {
			status := this_.GetServer().GetNetProxyInnerStatus(lineNodeIdList, find.Code)

			find.InnerStatus = status
			find.InnerIsStarted = status == node.StatusStarted
		} else {
			find.InnerStatus = 0
		}

		lineNodeIdList = this_.GetNodeLineTo(find.OuterServerId)
		if len(lineNodeIdList) > 0 {
			status := this_.GetServer().GetNetProxyOuterStatus(lineNodeIdList, find.Code)
			find.OuterStatus = status
			find.OuterIsStarted = status == node.StatusStarted
		} else {
			find.OuterStatus = 0
		}
	}

	this_.checkChangeOut()
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

	var nodeModelList = this_.getNodeModelList()
	var netProxyModelList = this_.getNetProxyModelList()

	newBs, _ := json.Marshal(nodeModelList)
	newNodeListStr := string(newBs)
	//this_.Logger.Info("node list validate", zap.Any("old", string(oldNodeBs)), zap.Any("new", string(newBs)))
	if oldNodeListStr != newNodeListStr {
		//this_.Logger.Info("node list change", zap.Any("old", oldNodeListStr), zap.Any("new", newNodeListStr))
		this_.oldNodeListStr = newNodeListStr

		userNodeList := make(map[int64][]*NodeModel)
		for _, one := range nodeModelList {
			userNodeList[one.UserId] = append(userNodeList[one.UserId], one)
		}
		for userId, list := range userNodeList {
			this_.callNodeListChange(userId, list)
		}
	}

	newBs, _ = json.Marshal(netProxyModelList)
	newNetProxyListStr := string(newBs)
	if oldNetProxyListStr != newNetProxyListStr {
		//this_.Logger.Info("net proxy list change", zap.Any("old", oldNetProxyListStr), zap.Any("new", newNetProxyListStr))
		this_.oldNetProxyListStr = newNetProxyListStr

		userNetProxyList := make(map[int64][]*NetProxyModel)
		for _, one := range netProxyModelList {
			userNetProxyList[one.UserId] = append(userNetProxyList[one.UserId], one)
		}
		for userId, list := range userNetProxyList {
			this_.callNetProxyListChange(userId, list)
		}
	}
	return
}
