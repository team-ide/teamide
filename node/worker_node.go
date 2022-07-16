package node

import (
	"teamide/pkg/util"
)

func (this_ *Worker) findRootNode() (find *Info) {
	for _, one := range this_.nodeList {
		if one.ParentId == "" {
			find = one
		}
	}
	return
}

func (this_ *Worker) findNode(id string) (find *Info) {
	for _, one := range this_.nodeList {
		if one.Id == id {
			find = one
		}
	}
	return
}

func (this_ *Worker) findChildrenNodeList(id string) (nodeList []*Info) {
	for _, one := range this_.nodeList {
		if one.ParentId == id {
			nodeList = append(nodeList, one)
		}
	}
	return
}

func (this_ *Worker) findChildrenNode(id string) (find *Info) {
	for _, one := range this_.childrenNodeList {
		if one.Id == id {
			find = one
		}
	}
	return
}

func (this_ *Worker) AddNode(node *Info) {
	this_.nodeLock.Lock()
	defer this_.nodeLock.Unlock()

	var list = this_.childrenNodeList
	for _, one := range list {
		listener := this_.getChildrenNodeListener(one)
		_, _ = this_.Call(one, listener, methodAddNode, &Message{
			Node: node,
		})
	}

	var find = this_.findNode(node.Id)

	if find == nil {
		this_.nodeList = append(this_.nodeList, node)
	} else {
		copyNode(node, find)
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

func (this_ *Worker) RemoveNode(node *Info) {
	this_.nodeLock.Lock()
	defer this_.nodeLock.Unlock()

	var list = this_.childrenNodeList
	var newList []*Info
	for _, one := range list {
		listener := this_.getChildrenNodeListener(one)
		_, _ = this_.Call(one, listener, methodRemoveNode, &Message{
			Node: node,
		})
		if one.Id != node.Id {
			newList = append(newList, one)
		} else {
			this_.removeChildrenNodeListener(one)
		}
	}
	this_.childrenNodeList = newList

	list = this_.nodeList
	newList = []*Info{}
	for _, one := range list {
		if one.Id != node.Id {
			newList = append(newList, one)
		}
	}
	this_.nodeList = newList

	this_.refreshNodeList()
}

func (this_ *Worker) refreshNodeList() {
	var list = this_.nodeList
	for _, one := range list {
		var find = this_.findChildrenNode(one.Id)

		if find == nil {
			if one.ParentId == this_.Node.Id {
				this_.addChildrenNode(one)
			}
		}
	}
	return
}

func (this_ *Worker) resetNodeList(nodeList []*Info) {
	var list = this_.nodeList
	for _, one := range list {
		this_.RemoveNode(one)
	}

	for _, one := range nodeList {
		this_.AddNode(one)
	}
	return
}
