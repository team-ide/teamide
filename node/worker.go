package node

import "teamide/pkg/util"

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

func (this_ *Worker) notifyDo(msg *Message) {
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
	msg.NotifyParent = true
	this_.notifyAllFrom(msg)
	this_.notifyDo(msg)
}

func (this_ *Worker) notifyChildren(msg *Message) {
	msg.NotifyChildren = true
	this_.notifyAllTo(msg)
	this_.notifyDo(msg)
}

func (this_ *Worker) notifyAll(msg *Message) {
	msg.NotifyAll = true
	this_.notifyAllTo(msg)
	this_.notifyAllFrom(msg)
	this_.notifyDo(msg)
}

func (this_ *Worker) notifyOther(msg *Message) {
	msg.NotifyAll = true
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
					NodeId:             nodeId,
					NotifiedNodeIdList: NotifiedNodeIdList,
				}
				res, _ := this_.Call(listener, methodGetNode, msg)
				if res != nil && res.Node != nil {
					find = res.Node
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
					NodeId:             nodeId,
					NotifiedNodeIdList: NotifiedNodeIdList,
				}
				res, _ := this_.Call(listener, methodGetNode, msg)
				if res != nil {
					version = res.Version
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
					NodeId:             nodeId,
					NotifiedNodeIdList: NotifiedNodeIdList,
				})
				if res != nil && res.MonitorData != nil {
					monitorData = res.MonitorData
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
					NetProxyId:         netProxyId,
					NotifiedNodeIdList: NotifiedNodeIdList,
				})
				if res != nil && res.MonitorData != nil {
					monitorData = res.MonitorData
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
					NetProxyId:         netProxyId,
					NotifiedNodeIdList: NotifiedNodeIdList,
				})
				if res != nil && res.MonitorData != nil {
					monitorData = res.MonitorData
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
