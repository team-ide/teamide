package node

import (
	"go.uber.org/zap"
	"teamide/pkg/util"
)

func (this_ *Worker) findRootNode() (find *Info) {
	var list = this_.cache.nodeList
	for _, one := range list {
		if one.ParentId == "" {
			find = one
		}
	}
	return
}

func (this_ *Worker) findNode(id string) (find *Info) {
	var list = this_.cache.nodeList
	for _, one := range list {
		if one.Id == id {
			find = one
		}
	}
	return
}

func (this_ *Worker) findChildrenNodeList(id string) (nodeList []*Info) {
	var list = this_.cache.nodeList
	for _, one := range list {
		if one.ParentId == id {
			nodeList = append(nodeList, one)
		}
	}
	return
}

func (this_ *Worker) findChildrenNode(id string) (find *Info) {
	var list = this_.cache.childrenNodeList
	for _, one := range list {
		if one.Id == id {
			find = one
		}
	}
	return
}

func (this_ *Worker) AddNodeList(nodeList []*Info) (err error) {
	err = this_.addNodeList(nodeList)
	return
}

func (this_ *Worker) RemoveNodeList(nodeIdList []string) (err error) {
	err = this_.removeNodeList(nodeIdList)
	return
}

func (this_ *Worker) addNodeList(nodeList []*Info) (err error) {
	if len(nodeList) == 0 {
		return
	}
	_ = this_.callChildrenNodePoolList(&Message{
		Method:   methodNodeAdd,
		NodeList: nodeList,
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

	_ = this_.callChildrenNodePoolList(&Message{
		Method:   methodNodeAdd,
		NodeList: nodeList,
	})

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

func (this_ *Worker) removeNodeList(nodeIdList []string) (err error) {
	if len(nodeIdList) == 0 {
		return
	}
	_ = this_.callChildrenNodePoolList(&Message{
		Method:     methodNodeRemove,
		NodeIdList: nodeIdList,
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
		var list = this_.cache.childrenNodeList
		var newList []*Info
		for _, one := range list {
			if one.Id != nodeId {
				newList = append(newList, one)
			} else {
				this_.cache.removeNodeListenerPool(this_.server.Id, nodeId)
			}
		}
		this_.cache.childrenNodeList = newList

		list = this_.cache.nodeList
		newList = []*Info{}
		for _, one := range list {
			if one.Id != nodeId {
				newList = append(newList, one)
			}
		}
		this_.cache.nodeList = newList
	}

	this_.refreshNodeList()

	return
}

func (this_ *Worker) findNodeLineToRoot(nodeId string) (list []string) {
	rootNode := this_.findRootNode()
	if rootNode == nil {
		return
	}
	if rootNode.Id == nodeId {
		list = append(list, nodeId)
		return
	}
	var id = nodeId
	var loadedIdList []string

	for {
		if util.ContainsString(loadedIdList, id) >= 0 {
			break
		}
		loadedIdList = append(loadedIdList, id)

		var find = this_.findNode(id)
		if find == nil {
			break
		}
		list = append(list, find.Id)
		if find.Id == rootNode.Id {
			break
		}
		id = find.ParentId
	}
	if len(list) > 0 {
		if list[0] != nodeId || list[len(list)-1] != rootNode.Id {
			list = []string{}
			return
		}
	}
	return
}

func (this_ *Worker) GetNodeLineByFromTo(fromNodeId, toNodeId string) (lineIdList []string) {

	fromToRootList := this_.findNodeLineToRoot(fromNodeId)
	if len(fromToRootList) == 0 {
		return
	}
	toToRootList := this_.findNodeLineToRoot(toNodeId)
	if len(toToRootList) == 0 {
		return
	}

	if len(fromToRootList) == 1 {
		for i := len(toToRootList) - 1; i >= 0; i-- {
			lineIdList = append(lineIdList, toToRootList[i])
		}
		return
	} else if len(toToRootList) == 1 {
		lineIdList = append(lineIdList, fromToRootList...)
		return
	} else {
		for i, id := range toToRootList {
			var find = false
			for n, one := range fromToRootList {
				if one == id {
					find = true
					if n > 0 {
						lineIdList = append(lineIdList, fromToRootList[:n]...)
					}
					break
				}

			}
			if find {
				for m := i; m >= 0; m-- {
					lineIdList = append(lineIdList, toToRootList[m])
				}
				break
			}
		}
	}
	return
}

func (this_ *Worker) refreshNodeList() {
	var list = this_.cache.nodeList
	for _, one := range list {
		var find = this_.findChildrenNode(one.Id)

		if find == nil && one.ParentId == this_.server.Id {
			this_.addChildrenNode(one)
		}
	}
	if this_.server.OnNodeListChange != nil {
		this_.server.OnNodeListChange(this_.cache.nodeList)
	}
	return
}
