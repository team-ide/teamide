package module_node

import "teamide/internal/context"

type NodeDataChange struct {
	Type         string          `json:"type,omitempty"`
	NodeList     []*NodeInfo     `json:"nodeList,omitempty"`
	NetProxyList []*NetProxyInfo `json:"netProxyList,omitempty"`
}

func (this_ *NodeContext) callNodeDataChange(data *NodeDataChange) {
	context.ServerWebsocketOutEvent("node-data-change", data)
}
