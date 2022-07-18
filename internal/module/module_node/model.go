package module_node

import "time"

const (
	// ModuleNode 节点模块
	ModuleNode = "node"
	// TableNode 节点表
	TableNode        = "TM_NODE"
	TableNodeComment = "节点"
)

// NodeModel 节点模型，和节点表对应
type NodeModel struct {
	NodeId       int64     `json:"nodeId,omitempty"`
	NodeServerId string    `json:"nodeServerId,omitempty"`
	Name         string    `json:"name,omitempty"`
	Comment      string    `json:"comment,omitempty"`
	Option       string    `json:"option,omitempty"`
	UserId       int64     `json:"userId,omitempty"`
	DeleteUserId int64     `json:"deleteUserId,omitempty"`
	Deleted      int8      `json:"deleted,omitempty"`
	CreateTime   time.Time `json:"createTime,omitempty"`
	UpdateTime   time.Time `json:"updateTime,omitempty"`
	DeleteTime   time.Time `json:"deleteTime,omitempty"`
}

func (entity *NodeModel) GetTableName() string {
	return TableNode
}

func (entity *NodeModel) GetPKColumnName() string {
	return ""
}
