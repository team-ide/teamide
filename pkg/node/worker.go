package node

import (
	"errors"
	"teamide/pkg/util"
)

type Worker struct {
	server *Server
	*Space
	MonitorData *MonitorData
}

func (this_ *Worker) Stop() {
	this_.removeToNodeListenerPoolList()
	this_.removeFromNodeListenerPoolList()
}

func (this_ *Worker) getVersion(lineNodeIdList []string) (version string) {

	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodGetVersion, &Message{
			LineNodeIdList: lineNodeIdList,
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
	version = util.GetVersion()
	return
}

func (this_ *Worker) getNodeStatus(lineNodeIdList []string) (status int8) {

	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodNodeGetStatus, &Message{
			LineNodeIdList: lineNodeIdList,
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		if resMsg != nil && resMsg.NodeWorkData != nil {
			status = resMsg.NodeWorkData.Status
		}
		return
	}
	status = StatusStarted
	return
}

func (this_ *Worker) addToNodeList(lineNodeIdList []string, toNodeList []*ToNode) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodNodeAddToNodeList, &Message{
			LineNodeIdList: lineNodeIdList,
			NodeWorkData: &WorkData{
				ToNodeList: toNodeList,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}
	_ = this_.doAddToNodeList(toNodeList)
	return
}

func (this_ *Worker) removeToNodeList(lineNodeIdList []string, toNodeIdList []string) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodNodeRemoveToNodeList, &Message{
			LineNodeIdList: lineNodeIdList,
			NodeWorkData: &WorkData{
				ToNodeIdList: toNodeIdList,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}
	_ = this_.doRemoveToNodeList(toNodeIdList)
	return
}

func (this_ *Worker) getNodeMonitorData(lineNodeIdList []string) (monitorData *MonitorData) {
	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodNodeGetNodeMonitorData, &Message{
			LineNodeIdList: lineNodeIdList,
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

	monitorData = this_.MonitorData
	return
}

func (this_ *Worker) getNetProxyInnerStatus(lineNodeIdList []string, netProxyId string) (status int8) {

	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodNetProxyGetInnerStatus, &Message{
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
			status = resMsg.NetProxyWorkData.Status
		}
		return
	}

	inner := this_.getNetProxyInner(netProxyId)
	if inner != nil {
		status = inner.status
	}
	return
}

func (this_ *Worker) addNetProxyInnerList(lineNodeIdList []string, netProxyInnerList []*NetProxyInner) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodNetProxyAddNetProxyInnerList, &Message{
			LineNodeIdList: lineNodeIdList,
			NetProxyWorkData: &NetProxyWorkData{
				NetProxyInnerList: netProxyInnerList,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}
	_ = this_.doAddNetProxyInnerList(netProxyInnerList)
	return
}

func (this_ *Worker) removeNetProxyInnerList(lineNodeIdList []string, netProxyIdList []string) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodNetProxyRemoveNetProxyInnerList, &Message{
			LineNodeIdList: lineNodeIdList,
			NetProxyWorkData: &NetProxyWorkData{
				NetProxyIdList: netProxyIdList,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}
	_ = this_.doRemoveNetProxyInnerList(netProxyIdList)
	return
}

func (this_ *Worker) getNetProxyOuterStatus(lineNodeIdList []string, netProxyId string) (status int8) {

	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodNetProxyGetOuterStatus, &Message{
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
			status = resMsg.NetProxyWorkData.Status
		}
		return
	}

	outer := this_.getNetProxyOuter(netProxyId)
	if outer != nil {
		status = StatusStarted
	}
	return
}

func (this_ *Worker) addNetProxyOuterList(lineNodeIdList []string, netProxyOuterList []*NetProxyOuter) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodNetProxyAddNetProxyOuterList, &Message{
			LineNodeIdList: lineNodeIdList,
			NetProxyWorkData: &NetProxyWorkData{
				NetProxyOuterList: netProxyOuterList,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}
	_ = this_.doAddNetProxyOuterList(netProxyOuterList)
	return
}

func (this_ *Worker) removeNetProxyOuterList(lineNodeIdList []string, netProxyIdList []string) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodNetProxyRemoveNetProxyOuterList, &Message{
			LineNodeIdList: lineNodeIdList,
			NetProxyWorkData: &NetProxyWorkData{
				NetProxyIdList: netProxyIdList,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}
	_ = this_.doRemoveNetProxyOuterList(netProxyIdList)
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

	var find = this_.getNetProxyInner(netProxyId)
	//Logger.Info("getNetProxyInnerMonitorData", zap.Any("netProxyId", netProxyId), zap.Any("find", find))
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

	var find = this_.getNetProxyOuter(netProxyId)
	if find != nil {
		monitorData = find.MonitorData
		return
	}

	return
}

func (this_ *Worker) sendToNext(lineNodeIdList []string, key string, doSend func(listener *MessageListener) (err error)) (send bool, err error) {
	if len(lineNodeIdList) == 0 {
		err = errors.New("节点线不存在")
		return
	}
	var thisNodeIndex = -1
	for _, localNode := range this_.server.localNodeList {
		index := util.ContainsString(lineNodeIdList, localNode.Id)
		if index >= 0 {
			thisNodeIndex = index
			if index == len(lineNodeIdList)-1 {
				break
			}
		}
	}
	if thisNodeIndex < 0 {
		err = errors.New(this_.server.GetServerInfo() + " 与节点 [" + util.ToJSON(lineNodeIdList) + "] 暂无通讯渠道")
		return
	}
	if thisNodeIndex != len(lineNodeIdList)-1 {
		nodeId := lineNodeIdList[thisNodeIndex+1]

		var pool = this_.getToNodeListenerPool(nodeId)
		if pool == nil {
			pool = this_.getFromNodeListenerPool(nodeId)
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
