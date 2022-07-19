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

func (this_ *Worker) initialize(nodeList []*Info, netProxyList []*NetProxy) {

	var oldList = this_.cache.nodeList
	var removeNodeIdList []string
	for _, one := range oldList {
		var find *Info
		for _, one_ := range nodeList {
			if one_.Id == one.Id {
				find = one_
			}
		}
		if find == nil {
			removeNodeIdList = append(removeNodeIdList, one.Id)
		}
	}
	_ = this_.removeNodeList(removeNodeIdList)
	_ = this_.addNodeList(nodeList)

	var oldNetProxyList = this_.cache.netProxyList
	var removeNetProxyIdList []string
	for _, one := range oldNetProxyList {
		var find *NetProxy
		for _, one_ := range netProxyList {
			if one_.Id == one.Id {
				find = one_
			}
		}
		if find == nil {
			removeNetProxyIdList = append(removeNetProxyIdList, one.Id)
		}
	}
	_ = this_.removeNetProxyList(removeNetProxyIdList)
	_ = this_.addNetProxyList(netProxyList)
}

func (this_ *Worker) notifyParentRefresh(nodeList []*Info, netProxyList []*NetProxy) {

	_ = this_.doAddNodeList(nodeList)
	_ = this_.doAddNetProxyList(netProxyList)
}
