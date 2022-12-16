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

func (this_ *NodeContext) callNodeCountDataChange(userId int64, data *NodeCountData) {
	context.CallUserEvent(userId, context.NewListenEvent("node-data-change", data))
}

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
