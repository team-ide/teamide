package node

type Worker struct {
	server *Server
	cache  *Cache
}

func (this_ *Worker) initialize(nodeList []*Info, netProxyList []*NetProxy) {

	var oldList = this_.cache.nodeList
	for _, one := range oldList {
		var find *Info
		for _, one_ := range nodeList {
			if one_.Id == one.Id {
				find = one_
			}
		}
		if find == nil {
			_ = this_.RemoveNode(one.Id)
		}
	}
	for _, one := range nodeList {
		_ = this_.AddNode(one)
	}

	var oldNetProxyList = this_.cache.netProxyList
	for _, one := range oldNetProxyList {
		var find *NetProxy
		for _, one_ := range netProxyList {
			if one_.Id == one.Id {
				find = one_
			}
		}
		if find == nil {
			_ = this_.RemoveNetProxy(one.Id)
		}
	}
	for _, one := range netProxyList {
		_ = this_.AddNetProxy(one)
	}
}
