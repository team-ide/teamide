package node

func (this_ *Worker) addChildrenNode(childrenNode *Info) {
	Logger.Info(this_.server.GetServerInfo() + " 添加子节点 " + childrenNode.GetNodeStr())
	var find = this_.findChildrenNode(childrenNode.Id)
	if find == nil {
		pool := this_.cache.getNodeListenerPool(this_.server.Id, childrenNode.Id)
		if pool == nil {
			this_.server.connNodeListenerKeepAlive(childrenNode.ConnAddress, childrenNode.ConnToken, childrenNode.ConnSize)
		}
		this_.cache.childrenNodeList = append(this_.cache.childrenNodeList, childrenNode)
	} else {
		copyNode(childrenNode, find)
	}
}

func (this_ *Worker) callChildrenNodePoolList(msg *Message) (err error) {
	var list = this_.cache.getNodeListenerPoolListByFromNodeId(this_.server.Id)
	for _, pool := range list {
		err = pool.Do(func(listener *MessageListener) (e error) {
			_, e = this_.Call(listener, msg.Method, msg)
			return
		})
	}
	return
}

func (this_ *Worker) callParentNodePoolList(msg *Message) (err error) {
	var list = this_.cache.getNodeListenerPoolListByToNodeId(this_.server.Id)
	for _, pool := range list {
		err = pool.Do(func(listener *MessageListener) (e error) {
			_, e = this_.Call(listener, msg.Method, msg)
			return
		})
	}
	return
}
