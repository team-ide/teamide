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
