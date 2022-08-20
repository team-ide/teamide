package module_node

import "teamide/internal/context"

type NodeCountData struct {
	Type                          string `json:"type,omitempty"`
	NodeCount                     int    `json:"nodeCount,omitempty"`
	NodeSuccessCount              int    `json:"nodeSuccessCount,omitempty"`
	NodeNetProxyCount             int    `json:"nodeNetProxyCount,omitempty"`
	NodeNetProxyInnerSuccessCount int    `json:"nodeNetProxyInnerSuccessCount,omitempty"`
	NodeNetProxyOuterSuccessCount int    `json:"nodeNetProxyOuterSuccessCount,omitempty"`
}

func newNodeCountData() *NodeCountData {
	return &NodeCountData{
		Type: "count",
	}
}

func (this_ *NodeContext) callNodeCountDataChange(data *NodeCountData) {
	context.ServerWebsocketOutEvent("node-data-change", data)
}

type NodeListChange struct {
	Type     string       `json:"type,omitempty"`
	NodeList []*NodeModel `json:"nodeList,omitempty"`
}

func (this_ *NodeContext) callNodeListChange(nodeList []*NodeModel) {
	context.ServerWebsocketOutEvent("node-data-change", &NodeListChange{
		Type:     "node-list",
		NodeList: nodeList,
	})
}

type NetProxyListChange struct {
	Type         string           `json:"type,omitempty"`
	NetProxyList []*NetProxyModel `json:"netProxyList,omitempty"`
}

func (this_ *NodeContext) callNetProxyListChange(netProxyList []*NetProxyModel) {
	context.ServerWebsocketOutEvent("node-data-change", &NetProxyListChange{
		Type:         "net-proxy-list",
		NetProxyList: netProxyList,
	})
}
