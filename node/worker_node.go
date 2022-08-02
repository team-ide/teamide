package node

import (
	"go.uber.org/zap"
	"teamide/pkg/util"
)

func (this_ *Worker) findNode(id string) (find *Info) {
	var list = this_.cache.nodeList
	for _, one := range list {
		if one.Id == id {
			find = one
		}
	}
	return
}

func (this_ *Worker) findNodeList(idList []string) (findList []*Info) {
	var list = this_.cache.nodeList
	for _, one := range list {
		if util.ContainsString(idList, one.Id) >= 0 {
			findList = append(findList, one)
		}
	}
	return
}

func (this_ *Worker) findParents(id string) (parents []*Info) {
	var list = this_.cache.nodeList
	for _, one := range list {
		if util.ContainsString(one.ConnNodeIdList, id) >= 0 {
			parents = append(parents, one)
		}
	}
	return
}

func (this_ *Worker) updateNodeConnNodeIdList(id string, connNodeIdList []string) (err error) {
	var find = this_.findNode(id)

	if find == nil {
		return
	}
	var newOne = &Info{
		Id: id,
	}
	_ = copyNode(find, newOne)
	newOne.ConnNodeIdList = connNodeIdList
	_ = this_.addNodeList([]*Info{newOne})

	return
}

func (this_ *Worker) addNodeList(nodeList []*Info) (err error) {
	if len(nodeList) == 0 {
		return
	}
	this_.notifyAll(&Message{
		NotifyChange: &NotifyChange{
			NodeList: nodeList,
		},
	})

	err = this_.doAddNodeList(nodeList)

	return
}

func (this_ *Worker) doAddNodeList(nodeList []*Info) (err error) {
	if len(nodeList) == 0 {
		return
	}
	this_.cache.nodeLock.Lock()
	defer this_.cache.nodeLock.Unlock()

	var findChanged = false
	for _, node := range nodeList {
		var find = this_.findNode(node.Id)

		if find == nil {
			Logger.Info(this_.server.GetServerInfo()+" 添加节点 ", zap.Any("node", node))
			findChanged = true
			this_.cache.nodeList = append(this_.cache.nodeList, node)
		} else {
			Logger.Info(this_.server.GetServerInfo()+" 更新节点 ", zap.Any("node", node))

			if copyNode(node, find) {
				findChanged = true
			}
		}
	}

	if findChanged {
		this_.refresh()

		if this_.server.OnNodeListChange != nil {
			this_.server.OnNodeListChange(this_.cache.nodeList)
		}
	}
	return
}

func (this_ *Worker) removeNodeList(removeNodeIdList []string) (err error) {
	if len(removeNodeIdList) == 0 {
		return
	}
	this_.notifyAll(&Message{
		NotifyChange: &NotifyChange{
			RemoveNodeIdList: removeNodeIdList,
		},
	})
	err = this_.doRemoveNodeList(removeNodeIdList)

	return
}

func (this_ *Worker) doRemoveNodeList(removeNodeIdList []string) (err error) {
	if len(removeNodeIdList) == 0 {
		return
	}
	this_.cache.nodeLock.Lock()
	defer this_.cache.nodeLock.Unlock()

	Logger.Info(this_.server.GetServerInfo()+" 移除节点 ", zap.Any("removeNodeIdList", removeNodeIdList))

	var findChanged = false

	var list = this_.cache.nodeList
	var newList []*Info
	for _, nodeId := range removeNodeIdList {
		this_.cache.removeNodeListenerPool(this_.server.Id, nodeId)
		this_.cache.removeNodeListenerPool(nodeId, this_.server.Id)
	}
	for _, one := range list {
		if util.ContainsString(removeNodeIdList, one.Id) >= 0 {
			findChanged = true
		} else {
			newList = append(newList, one)
		}
		var newConnNodeIdList []string
		for _, nodeId := range one.ConnNodeIdList {
			if util.ContainsString(removeNodeIdList, nodeId) >= 0 {
				findChanged = true
			} else {
				newConnNodeIdList = append(newConnNodeIdList, nodeId)
			}
		}
		one.ConnNodeIdList = newConnNodeIdList
	}
	this_.cache.nodeList = newList

	if findChanged {
		this_.refresh()

		if this_.server.OnNodeListChange != nil {
			this_.server.OnNodeListChange(this_.cache.nodeList)
		}
	}

	return
}

func (this_ *Worker) removeNodeConnNodeIdList(id string, removeConnNodeIdList []string) (err error) {
	if len(removeConnNodeIdList) == 0 {
		return
	}
	this_.notifyAll(&Message{
		NotifyChange: &NotifyChange{
			NodeId:               id,
			RemoveConnNodeIdList: removeConnNodeIdList,
		},
	})

	_ = this_.doRemoveNodeConnNodeIdList(id, removeConnNodeIdList)

	return
}

func (this_ *Worker) doRemoveNodeConnNodeIdList(id string, removeConnNodeIdList []string) (err error) {
	if len(removeConnNodeIdList) == 0 {
		return
	}
	var find = this_.findNode(id)
	if find == nil {
		return
	}
	var removeNetProxyIdList []string
	for _, nodeId := range removeConnNodeIdList {
		this_.cache.removeNodeListenerPool(find.Id, nodeId)
		this_.cache.removeNodeListenerPool(nodeId, find.Id)
		removeNetProxyIdList = append(removeNetProxyIdList, this_.findNetProxyIdListByNodeId(nodeId)...)
	}
	var newConnNodeIdList []string
	var findChanged bool
	for _, nodeId := range find.ConnNodeIdList {
		if util.ContainsString(removeConnNodeIdList, nodeId) >= 0 {
			findChanged = true
		} else {
			newConnNodeIdList = append(newConnNodeIdList, nodeId)
		}
	}
	find.ConnNodeIdList = newConnNodeIdList

	if findChanged {
		this_.refresh()
		if this_.server.OnNodeListChange != nil {
			this_.server.OnNodeListChange(this_.cache.nodeList)
		}
	}
	_ = this_.doRemoveNetProxyList(removeNetProxyIdList)

	return
}

func (this_ *Worker) doChangeNodeStatus(statusChangeList []*StatusChange) (err error) {
	if len(statusChangeList) == 0 {
		return
	}

	var findChanged = false
	for _, one := range statusChangeList {
		if one.Id == "" || one.Status == 0 {
			continue
		}
		var find = this_.findNode(one.Id)
		if find == nil {
			continue
		}
		if copyNode(&Info{
			Status:      one.Status,
			StatusError: one.StatusError,
		}, find) {
			findChanged = true
		}
	}
	if findChanged && this_.server.OnNodeListChange != nil {
		this_.server.OnNodeListChange(this_.cache.nodeList)
	}
	return
}

func (this_ *Worker) refreshNodeList() {

	var statusChangeList []*StatusChange
	var connIdList = this_.server.rootNode.ConnNodeIdList
	for _, connToId := range connIdList {
		var find = this_.findNode(connToId)
		if find == nil {
			pool := this_.cache.getNodeListenerPool(this_.server.Id, connToId)
			if pool != nil {
				pool.Stop()
			}
		} else {
			pool := this_.cache.getNodeListenerPool(this_.server.Id, connToId)
			if pool == nil { //  && find.IsEnabled()
				if find != nil && find.ConnAddress != "" {
					this_.server.connNodeListenerKeepAlive(find, find.ConnAddress, find.ConnToken, find.ConnSize)
				}
			}
		}

		pool := this_.cache.getNodeListenerPool(this_.server.Id, connToId)

		if pool != nil && len(pool.listeners) > 0 {
			statusChangeList = append(statusChangeList, &StatusChange{
				Id:          connToId,
				Status:      StatusStarted,
				StatusError: "",
			})
		}
	}
	statusChangeList = append(statusChangeList, &StatusChange{
		Id:          this_.server.Id,
		Status:      StatusStarted,
		StatusError: "",
	})

	if len(statusChangeList) > 0 {
		this_.notifyOther(&Message{
			NotifyChange: &NotifyChange{
				NodeStatusChangeList: statusChangeList,
			},
		})
	}

	return
}

func (this_ *Worker) notifyAllTo(msg *Message) {
	if util.ContainsString(msg.NotifiedNodeIdList, this_.server.Id) < 0 {
		msg.NotifiedNodeIdList = append(msg.NotifiedNodeIdList, this_.server.Id)
	}
	var list = this_.cache.getNodeListenerPoolListByFromNodeId(this_.server.Id)
	var callPools []*MessageListenerPool
	for _, pool := range list {
		if util.ContainsString(msg.NotifiedNodeIdList, pool.toNodeId) >= 0 {
			continue
		}
		msg.NotifiedNodeIdList = append(msg.NotifiedNodeIdList, pool.toNodeId)
		callPools = append(callPools, pool)
	}
	for _, pool := range callPools {
		if pool.isStop {
			continue
		}
		_ = pool.Do("", func(listener *MessageListener) (e error) {
			_ = listener.Send(msg, this_.MonitorData)
			return
		})
	}
	return
}

func (this_ *Worker) notifyAllFrom(msg *Message) {
	if util.ContainsString(msg.NotifiedNodeIdList, this_.server.Id) < 0 {
		msg.NotifiedNodeIdList = append(msg.NotifiedNodeIdList, this_.server.Id)
	}
	var list = this_.cache.getNodeListenerPoolListByToNodeId(this_.server.Id)
	var callPools []*MessageListenerPool
	for _, pool := range list {
		if util.ContainsString(msg.NotifiedNodeIdList, pool.fromNodeId) >= 0 {
			continue
		}
		msg.NotifiedNodeIdList = append(msg.NotifiedNodeIdList, pool.fromNodeId)
		callPools = append(callPools, pool)
	}
	for _, pool := range callPools {
		if pool.isStop {
			continue
		}
		_ = pool.Do("", func(listener *MessageListener) (e error) {
			_ = listener.Send(msg, this_.MonitorData)
			return
		})
	}
	return
}
