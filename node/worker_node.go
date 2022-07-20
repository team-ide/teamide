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

func (this_ *Worker) addNodeList(nodeList []*Info, calledNodeIdList []string) (err error) {
	if len(nodeList) == 0 {
		return
	}
	_ = this_.callAllTo(&Message{
		Method:           methodNodeAdd,
		NodeList:         nodeList,
		CalledNodeIdList: calledNodeIdList,
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
	Logger.Info(this_.server.GetServerInfo()+" 添加节点 ", zap.Any("nodeList", nodeList))

	for _, one := range nodeList {
		var find = this_.findNode(one.Id)

		if find == nil {
			this_.cache.nodeList = append(this_.cache.nodeList, one)
		} else {
			copyNode(one, find)
		}
	}

	this_.refreshNodeList()
	return
}

func (this_ *Worker) removeNodeList(nodeIdList []string, calledNodeIdList []string) (err error) {
	if len(nodeIdList) == 0 {
		return
	}
	_ = this_.callAllTo(&Message{
		Method:           methodNodeRemove,
		NodeIdList:       nodeIdList,
		CalledNodeIdList: calledNodeIdList,
	})

	err = this_.doRemoveNodeList(nodeIdList)

	return
}

func (this_ *Worker) doRemoveNodeList(nodeIdList []string) (err error) {
	if len(nodeIdList) == 0 {
		return
	}
	this_.cache.nodeLock.Lock()
	defer this_.cache.nodeLock.Unlock()

	Logger.Info(this_.server.GetServerInfo()+" 移除节点 ", zap.Any("nodeIdList", nodeIdList))

	for _, nodeId := range nodeIdList {

		var list = this_.cache.nodeList
		var newList []*Info
		for _, one := range list {
			if one.Id != nodeId {
				newList = append(newList, one)
			} else {
				this_.cache.removeNodeListenerPool(this_.server.Id, nodeId)
				this_.cache.removeNodeListenerPool(nodeId, this_.server.Id)
			}
		}
		this_.cache.nodeList = newList
	}

	this_.refreshNodeList()

	return
}

func (this_ *Worker) appendNodeLineList(loadedIdList *[]string, lineList *[][]string, parentLine []string, nodeList []*Info) {

	for _, one := range nodeList {
		var line []string
		line = append(line, parentLine...)

		if util.ContainsString(line, one.Id) >= 0 {
			continue
		}
		line = append(line, one.Id)

		*lineList = append(*lineList, line)

		if util.ContainsString(*loadedIdList, one.Id) >= 0 {
			continue
		}
		*loadedIdList = append(*loadedIdList, one.Id)

		var children = this_.findNodeList(one.ConnNodeIdList)
		this_.appendNodeLineList(loadedIdList, lineList, line, children)

		var parentList = this_.findParents(one.Id)
		this_.appendNodeLineList(loadedIdList, lineList, line, parentList)
	}
}

/**




 */
func (this_ *Worker) findNodeLineList(nodeId string) (lineList [][]string) {
	Logger.Info("查询节点所有线", zap.Any("nodeId", nodeId))

	var loadedIdList []string
	loadedIdList = append(loadedIdList, nodeId)
	var line []string
	line = append(line, nodeId)

	var find = this_.findNode(nodeId)
	if find != nil {
		var children = this_.findNodeList(find.ConnNodeIdList)
		this_.appendNodeLineList(&loadedIdList, &lineList, line, children)
	}
	var parentList = this_.findParents(nodeId)
	this_.appendNodeLineList(&loadedIdList, &lineList, line, parentList)

	return
}

func (this_ *Worker) getNodeLineByFromTo(fromNodeId, toNodeId string) (lineIdList []string) {

	Logger.Info("查询节点线", zap.Any("fromNodeId", fromNodeId), zap.Any("toNodeId", toNodeId))
	var lineList = this_.findNodeLineList(fromNodeId)

	for _, line := range lineList {
		Logger.Info("已查询的连线", zap.Any("fromNodeId", fromNodeId), zap.Any("line", line))
		if util.ContainsString(line, toNodeId) >= 0 {
			if len(lineIdList) == 0 || len(line) < len(lineIdList) {
				lineIdList = line
			}
		}
	}
	return
}

func (this_ *Worker) refreshNodeList() {
	var list = this_.cache.nodeList
	for _, one := range list {
		if one.Id == this_.server.Id {
			var connIdList = one.ConnNodeIdList
			for _, connToId := range connIdList {
				pool := this_.cache.getNodeListenerPool(this_.server.Id, connToId)
				if pool == nil {
					var find = this_.findNode(connToId)
					if find != nil && find.ConnAddress != "" {
						this_.server.connNodeListenerKeepAlive(find.ConnAddress, find.ConnToken, find.ConnSize)
					}
				}
			}
		}
	}
	if this_.server.OnNodeListChange != nil {
		this_.server.OnNodeListChange(this_.cache.nodeList)
	}
	return
}

func (this_ *Worker) callAllTo(msg *Message) (err error) {
	var list = this_.cache.getNodeListenerPoolListByFromNodeId(this_.server.Id)
	for _, pool := range list {
		if util.ContainsString(msg.CalledNodeIdList, pool.toNodeId) >= 0 {
			continue
		}
		msg.CalledNodeIdList = append(msg.CalledNodeIdList, pool.toNodeId)
		err = pool.Do(func(listener *MessageListener) (e error) {
			_, e = this_.Call(listener, msg.Method, msg)
			return
		})
	}
	return
}

func (this_ *Worker) callAllFrom(msg *Message) (err error) {
	var list = this_.cache.getNodeListenerPoolListByToNodeId(this_.server.Id)
	for _, pool := range list {
		if util.ContainsString(msg.CalledNodeIdList, pool.fromNodeId) >= 0 {
			continue
		}
		msg.CalledNodeIdList = append(msg.CalledNodeIdList, pool.fromNodeId)
		err = pool.Do(func(listener *MessageListener) (e error) {
			_, e = this_.Call(listener, msg.Method, msg)
			return
		})
	}
	return
}
