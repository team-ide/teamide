package node

type Worker struct {
	server *Server
	cache  *Cache
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

func (this_ *Worker) initialize(msg *Message) {

	var oldList = this_.cache.nodeList
	var removeNodeIdList []string
	for _, one := range oldList {
		var find *Info
		for _, one_ := range msg.NodeList {
			if one_.Id == one.Id {
				find = one_
			}
		}
		if find == nil {
			removeNodeIdList = append(removeNodeIdList, one.Id)
		}
	}
	_ = this_.removeNodeList(removeNodeIdList, msg.CalledNodeIdList)
	_ = this_.addNodeList(msg.NodeList, msg.CalledNodeIdList)

	var oldNetProxyList = this_.cache.netProxyList
	var removeNetProxyIdList []string
	for _, one := range oldNetProxyList {
		var find *NetProxy
		for _, one_ := range msg.NetProxyList {
			if one_.Id == one.Id {
				find = one_
			}
		}
		if find == nil {
			removeNetProxyIdList = append(removeNetProxyIdList, one.Id)
		}
	}
	_ = this_.removeNetProxyList(removeNetProxyIdList, msg.CalledNodeIdList)
	_ = this_.addNetProxyList(msg.NetProxyList, msg.CalledNodeIdList)
}

func (this_ *Worker) notifyParentRefresh(msg *Message) {
	_ = this_.callAllFrom(msg)
	_ = this_.doAddNodeList(msg.NodeList)
	_ = this_.doAddNetProxyList(msg.NetProxyList)
}
