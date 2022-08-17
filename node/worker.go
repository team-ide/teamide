package node

import (
	"errors"
	"teamide/pkg/util"
)

type Worker struct {
	server      *Server
	space       *Space
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
	if len(msg.NetProxyInnerStatusChangeList) > 0 || len(msg.NetProxyOuterStatusChangeList) > 0 {
		_ = this_.doChangeNetProxyStatus(msg.NetProxyInnerStatusChangeList, msg.NetProxyOuterStatusChangeList)
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
	this_.notifyDo(msg.NotifyChange)
	this_.notifyAllFrom(msg)
}

func (this_ *Worker) notifyChildren(msg *Message) {
	if msg.NotifyChange == nil {
		return
	}
	msg.NotifyChange.NotifyChildren = true
	this_.notifyDo(msg.NotifyChange)
	this_.notifyAllTo(msg)
}

func (this_ *Worker) notifyAll(msg *Message) {
	if msg.NotifyChange == nil {
		return
	}
	msg.NotifyChange.NotifyAll = true
	this_.notifyDo(msg.NotifyChange)
	this_.notifyAllTo(msg)
	this_.notifyAllFrom(msg)
}

func (this_ *Worker) notifyOther(msg *Message) {
	if msg.NotifyChange == nil {
		return
	}
	msg.NotifyChange.NotifyAll = true
	this_.notifyAllTo(msg)
	this_.notifyAllFrom(msg)
}

func (this_ *Worker) getVersion(lineNodeIdList []string, nodeId string) (version string) {

	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodGetVersion, &Message{
			LineNodeIdList: lineNodeIdList,
			NodeWorkData: &NodeWorkData{
				NodeId: nodeId,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		if resMsg != nil && resMsg.NodeWorkData != nil {
			version = resMsg.NodeWorkData.Version
		}
		return
	}

	if nodeId == "" {
		return ""
	}
	if this_.server.Id == nodeId {
		return util.GetVersion()
	}

	return
}

func (this_ *Worker) getNode(lineNodeIdList []string, nodeId string) (find *Info) {
	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodGetNode, &Message{
			LineNodeIdList: lineNodeIdList,
			NodeWorkData: &NodeWorkData{
				NodeId: nodeId,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		if resMsg != nil && resMsg.NodeWorkData != nil {
			find = resMsg.NodeWorkData.Node
		}
		return
	}

	if nodeId == "" {
		return
	}
	if this_.server.rootNode.Id == nodeId {
		find = this_.server.rootNode
		return
	}
	return
}

func (this_ *Worker) getOtherPool1(NotifiedNodeIdList *[]string) (callPools []*MessageListenerPool) {
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

func (this_ *Worker) getNodeMonitorData(lineNodeIdList []string, nodeId string) (monitorData *MonitorData) {
	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodGetNodeMonitorData, &Message{
			LineNodeIdList: lineNodeIdList,
			NodeWorkData: &NodeWorkData{
				NodeId: nodeId,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		if resMsg != nil && resMsg.NodeWorkData != nil {
			monitorData = resMsg.NodeWorkData.MonitorData
		}
		return
	}

	if nodeId == "" {
		return
	}
	if this_.server.rootNode.Id == nodeId {
		monitorData = this_.MonitorData
		return
	}

	return
}

func (this_ *Worker) getNetProxyInnerMonitorData(lineNodeIdList []string, netProxyId string) (monitorData *MonitorData) {
	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodNetProxyGetInnerMonitorData, &Message{
			LineNodeIdList: lineNodeIdList,
			NetProxyWorkData: &NetProxyWorkData{
				NetProxyId: netProxyId,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		if resMsg != nil && resMsg.NetProxyWorkData != nil {
			monitorData = resMsg.NetProxyWorkData.MonitorData
		}
		return
	}

	if netProxyId == "" {
		return
	}
	var find = this_.getNetProxyInner(netProxyId)
	if find != nil {
		monitorData = find.MonitorData
		return
	}

	return
}

func (this_ *Worker) getNetProxyOuterMonitorData(lineNodeIdList []string, netProxyId string) (monitorData *MonitorData) {
	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodNetProxyGetOuterMonitorData, &Message{
			LineNodeIdList: lineNodeIdList,
			NetProxyWorkData: &NetProxyWorkData{
				NetProxyId: netProxyId,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		if resMsg != nil && resMsg.NetProxyWorkData != nil {
			monitorData = resMsg.NetProxyWorkData.MonitorData
		}
		return
	}

	if netProxyId == "" {
		return
	}
	var find = this_.getNetProxyOuter(netProxyId)
	if find != nil {
		monitorData = find.MonitorData
		return
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
		var listener *MessageListener
		listener, err = pool.GetOne(key)
		if err != nil {
			return
		}
		err = doSend(listener)
		if err != nil {
			return
		}
		send = true
		return
	}
	return
}
