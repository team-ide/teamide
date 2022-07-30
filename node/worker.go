package node

import (
	"errors"
	"teamide/pkg/util"
)

type Worker struct {
	server      *Server
	cache       *Cache
	MonitorData *MonitorData
}

func (this_ *Worker) Stop() {
	var list = this_.cache.getNodeListenerPoolListByFromNodeId(this_.server.Id)
	for _, pool := range list {
		this_.cache.removeNodeListenerPool(pool.fromNodeId, pool.toNodeId)
	}
	list = this_.cache.getNodeListenerPoolListByToNodeId(this_.server.Id)
	for _, pool := range list {
		this_.cache.removeNodeListenerPool(pool.fromNodeId, pool.toNodeId)
	}
}

func (this_ *Worker) notifyDo(msg *NotifyChange) {
	if msg == nil {
		return
	}
	if len(msg.NodeStatusChangeList) > 0 {
		_ = this_.doChangeNodeStatus(msg.NodeStatusChangeList)
	}
	if len(msg.NetProxyInnerStatusChangeList) > 0 {
		_ = this_.doChangeNetProxyInnerStatus(msg.NetProxyInnerStatusChangeList)
	}
	if len(msg.NetProxyOuterStatusChangeList) > 0 {
		_ = this_.doChangeNetProxyOuterStatus(msg.NetProxyOuterStatusChangeList)
	}
	if msg.NodeId != "" && len(msg.RemoveConnNodeIdList) > 0 {
		_ = this_.doRemoveNodeConnNodeIdList(msg.NodeId, msg.RemoveConnNodeIdList)
	}
	if len(msg.RemoveNodeIdList) > 0 {
		_ = this_.doRemoveNodeList(msg.RemoveNodeIdList)
	}
	if len(msg.RemoveNetProxyIdList) > 0 {
		_ = this_.doRemoveNetProxyList(msg.RemoveNetProxyIdList)
	}
	if len(msg.NodeList) > 0 {
		_ = this_.doAddNodeList(msg.NodeList)
	}
	if len(msg.NetProxyList) > 0 {
		_ = this_.doAddNetProxyList(msg.NetProxyList)
	}
}
func (this_ *Worker) notifyParent(msg *Message) {
	if msg.NotifyChange == nil {
		return
	}
	msg.NotifyChange.NotifyParent = true
	this_.notifyAllFrom(msg)
	this_.notifyDo(msg.NotifyChange)
}

func (this_ *Worker) notifyChildren(msg *Message) {
	if msg.NotifyChange == nil {
		return
	}
	msg.NotifyChange.NotifyChildren = true
	this_.notifyAllTo(msg)
	this_.notifyDo(msg.NotifyChange)
}

func (this_ *Worker) notifyAll(msg *Message) {
	if msg.NotifyChange == nil {
		return
	}
	msg.NotifyChange.NotifyAll = true
	this_.notifyAllTo(msg)
	this_.notifyAllFrom(msg)
	this_.notifyDo(msg.NotifyChange)
}

func (this_ *Worker) notifyOther(msg *Message) {
	if msg.NotifyChange == nil {
		return
	}
	msg.NotifyChange.NotifyAll = true
	this_.notifyAllTo(msg)
	this_.notifyAllFrom(msg)
}

func (this_ *Worker) getNode(nodeId string, NotifiedNodeIdList []string) (find *Info) {
	if nodeId == "" {
		return
	}
	find = this_.findNode(nodeId)
	if find != nil {
		return
	}

	var callPools = this_.getOtherPool(&NotifiedNodeIdList)

	for _, pool := range callPools {
		if find == nil {
			_ = pool.Do("", func(listener *MessageListener) (e error) {
				msg := &Message{
					NotifiedNodeIdList: NotifiedNodeIdList,
					NodeWorkData: &NodeWorkData{
						NodeId: nodeId,
					},
				}
				res, _ := this_.Call(listener, methodGetNode, msg)
				if res != nil && res.NodeWorkData != nil && res.NodeWorkData.Node != nil {
					find = res.NodeWorkData.Node
				}
				return
			})
		}
	}
	return
}

func (this_ *Worker) getVersion(nodeId string, NotifiedNodeIdList []string) string {
	if nodeId == "" {
		return ""
	}
	if this_.server.Id == nodeId {
		return util.GetVersion()
	}

	var callPools = this_.getOtherPool(&NotifiedNodeIdList)

	var version string
	for _, pool := range callPools {
		if version == "" {
			_ = pool.Do("", func(listener *MessageListener) (e error) {
				msg := &Message{
					NotifiedNodeIdList: NotifiedNodeIdList,
					NodeWorkData: &NodeWorkData{
						NodeId: nodeId,
					},
				}
				res, _ := this_.Call(listener, methodGetNode, msg)
				if res != nil && res.NodeWorkData != nil {
					version = res.NodeWorkData.Version
				}
				return
			})
		}
	}
	return version
}

func (this_ *Worker) getOtherPool(NotifiedNodeIdList *[]string) (callPools []*MessageListenerPool) {
	if util.ContainsString(*NotifiedNodeIdList, this_.server.Id) < 0 {
		*NotifiedNodeIdList = append(*NotifiedNodeIdList, this_.server.Id)
	}
	var list = this_.cache.getNodeListenerPoolListByToNodeId(this_.server.Id)
	for _, pool := range list {
		if util.ContainsString(*NotifiedNodeIdList, pool.fromNodeId) >= 0 {
			continue
		}
		*NotifiedNodeIdList = append(*NotifiedNodeIdList, pool.fromNodeId)
		callPools = append(callPools, pool)
	}

	list = this_.cache.getNodeListenerPoolListByFromNodeId(this_.server.Id)

	for _, pool := range list {
		if util.ContainsString(*NotifiedNodeIdList, pool.toNodeId) >= 0 {
			continue
		}
		*NotifiedNodeIdList = append(*NotifiedNodeIdList, pool.toNodeId)
		callPools = append(callPools, pool)
	}
	return
}

func (this_ *Worker) getNodeMonitorData(nodeId string, NotifiedNodeIdList []string) (monitorData *MonitorData) {
	if nodeId == "" {
		return
	}
	if this_.server.rootNode.Id == nodeId {
		monitorData = this_.MonitorData
		return
	}

	var callPools = this_.getOtherPool(&NotifiedNodeIdList)

	for _, pool := range callPools {
		if monitorData == nil {
			_ = pool.Do("", func(listener *MessageListener) (e error) {
				res, _ := this_.Call(listener, methodGetNodeMonitorData, &Message{
					NotifiedNodeIdList: NotifiedNodeIdList,
					NodeWorkData: &NodeWorkData{
						NodeId: nodeId,
					},
				})
				if res != nil && res.NodeWorkData != nil && res.NodeWorkData.MonitorData != nil {
					monitorData = res.NodeWorkData.MonitorData
				}
				return
			})
		}
	}
	return
}

func (this_ *Worker) getNetProxyInnerMonitorData(netProxyId string, NotifiedNodeIdList []string) (monitorData *MonitorData) {
	if netProxyId == "" {
		return
	}
	var find = this_.getNetProxyInner(netProxyId)
	if find != nil {
		monitorData = find.MonitorData
		return
	}

	var callPools = this_.getOtherPool(&NotifiedNodeIdList)

	for _, pool := range callPools {
		if monitorData == nil {
			_ = pool.Do("", func(listener *MessageListener) (e error) {
				res, _ := this_.Call(listener, methodNetProxyGetInnerMonitorData, &Message{
					NetProxyWorkData: &NetProxyWorkData{
						NetProxyId: netProxyId,
					},
					NotifiedNodeIdList: NotifiedNodeIdList,
				})
				if res != nil && res.NetProxyWorkData != nil && res.NetProxyWorkData.MonitorData != nil {
					monitorData = res.NetProxyWorkData.MonitorData
				}
				return
			})
		}
	}
	return
}

func (this_ *Worker) getNetProxyOuterMonitorData(netProxyId string, NotifiedNodeIdList []string) (monitorData *MonitorData) {
	if netProxyId == "" {
		return
	}
	var find = this_.getNetProxyOuter(netProxyId)
	if find != nil {
		monitorData = find.MonitorData
		return
	}

	var callPools = this_.getOtherPool(&NotifiedNodeIdList)

	for _, pool := range callPools {
		if monitorData == nil {
			_ = pool.Do("", func(listener *MessageListener) (e error) {
				res, _ := this_.Call(listener, methodNetProxyGetOuterMonitorData, &Message{
					NetProxyWorkData: &NetProxyWorkData{
						NetProxyId: netProxyId,
					},
					NotifiedNodeIdList: NotifiedNodeIdList,
				})
				if res != nil && res.NetProxyWorkData != nil && res.NetProxyWorkData.MonitorData != nil {
					monitorData = res.NetProxyWorkData.MonitorData
				}
				return
			})
		}
	}
	return
}

func (this_ *Worker) refresh() {

	this_.refreshNodeList()
	this_.refreshNetProxy()
	return
}

func (this_ *Worker) sendToNext(lineNodeIdList []string, key string, doSend func(listener *MessageListener) (err error)) (send bool, err error) {
	if len(lineNodeIdList) == 0 {
		err = errors.New("节点线不存在")
		return
	}
	var thisNodeIndex = util.ContainsString(lineNodeIdList, this_.server.Id)
	if thisNodeIndex != len(lineNodeIdList)-1 {
		nodeId := lineNodeIdList[thisNodeIndex+1]

		var pool = this_.cache.getNodeListenerPool(this_.server.Id, nodeId)
		if pool == nil {
			pool = this_.cache.getNodeListenerPool(nodeId, this_.server.Id)
		}
		if pool == nil {
			err = errors.New(this_.server.GetServerInfo() + " 与节点 [" + nodeId + "] 暂无通讯渠道")
			return
		}
		err = pool.Do(key, func(listener *MessageListener) (e error) {
			e = doSend(listener)
			return
		})
		if err != nil {
			return
		}
		send = true
		return
	}
	return
}
