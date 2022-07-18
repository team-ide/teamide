package module_node

import (
	"teamide/node"
	"teamide/pkg/util"
)

type NodeContext struct {
	server      *node.Server
	nodeService *NodeService
}

func (this_ *NodeService) InitContext() (err error) {
	if this_.nodeContext == nil {
		this_.nodeContext = &NodeContext{
			nodeService: this_,
		}
	}
	err = this_.nodeContext.initContext()
	if err != nil {
		return
	}
	return
}

func (this_ *NodeContext) initContext() (err error) {
	var list []*NodeModel
	list, err = this_.nodeService.Query(&NodeModel{})
	if err != nil {
		return
	}

	var root *NodeModel
	for _, one := range list {
		if one.ParentServerId == "" {
			root = one
			break
		}
	}
	if root == nil {
		root = &NodeModel{
			ServerId: "root",
			Name:     "ROOT",
			Comment:  "根节点",
			Address:  ":21090",
			Token:    util.UUID(),
		}
		_, err = this_.nodeService.Insert(root)
		if err != nil {
			return
		}
		list = append(list, root)
	}

	this_.server = &node.Server{
		Id:      root.ServerId,
		Token:   root.Token,
		Address: root.Address,
	}
	err = this_.server.Start()
	if err != nil {
		return
	}
	return
}
