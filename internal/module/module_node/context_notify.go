package module_node

import "teamide/internal/context"

type NodeListChange struct {
	Type     string       `json:"type,omitempty"`
	NodeList []*NodeModel `json:"nodeList,omitempty"`
}

func (this_ *NodeContext) callNodeListChange(userId int64, nodeList []*NodeModel) {
	context.CallUserEvent(userId, context.NewListenEvent("node-data-change", &NodeListChange{
		Type:     "node-list",
		NodeList: nodeList,
	}))
}

type NetProxyListChange struct {
	Type         string           `json:"type,omitempty"`
	NetProxyList []*NetProxyModel `json:"netProxyList,omitempty"`
}

func (this_ *NodeContext) callNetProxyListChange(userId int64, netProxyList []*NetProxyModel) {
	context.CallUserEvent(userId, context.NewListenEvent("node-data-change", &NetProxyListChange{
		Type:         "net-proxy-list",
		NetProxyList: netProxyList,
	}))
}
