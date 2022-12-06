package module_node

import (
	"encoding/json"
	"time"
)

const (
	// ModuleNode 节点模块
	ModuleNode = "node"
	// TableNode 节点表
	TableNode        = "TM_NODE"
	TableNodeComment = "节点"
	// TableNodeNetProxy 节点网络代理表
	TableNodeNetProxy        = "TM_NODE_NET_PROXY"
	TableNodeNetProxyComment = "节点网络代理"
)

// NodeModel 节点模型，和节点表对应
type NodeModel struct {
	NodeId               int64     `json:"nodeId,omitempty"`
	ServerId             string    `json:"serverId,omitempty"`
	Name                 string    `json:"name,omitempty"`
	Comment              string    `json:"comment,omitempty"`
	BindAddress          string    `json:"bindAddress,omitempty"`
	BindToken            string    `json:"bindToken,omitempty"`
	ConnAddress          string    `json:"connAddress,omitempty"`
	ConnToken            string    `json:"connToken,omitempty"`
	ConnServerIds        string    `json:"connServerIds,omitempty"`
	HistoryConnServerIds string    `json:"historyConnServerIds,omitempty"`
	Option               string    `json:"option,omitempty"`
	IsLocal              int8      `json:"isLocal"`
	UserId               int64     `json:"userId,omitempty"`
	Enabled              int8      `json:"enabled"`
	DeleteUserId         int64     `json:"deleteUserId,omitempty"`
	Deleted              int8      `json:"deleted"`
	CreateTime           time.Time `json:"createTime,omitempty"`
	UpdateTime           time.Time `json:"updateTime,omitempty"`
	DeleteTime           time.Time `json:"deleteTime,omitempty"`

	ConnServerIdList        []string `json:"connServerIdList,omitempty"`
	HistoryConnServerIdList []string `json:"historyConnServerIdList,omitempty"`
	IsStarted               bool     `json:"isStarted"`
	Status                  int8     `json:"status"`
}

func GetStringList(str string) []string {
	var list []string
	if str != "" {
		_ = json.Unmarshal([]byte(str), &list)
	}
	return list
}

func GetListToString(list []string) string {
	if len(list) > 0 {
		bs, _ := json.Marshal(list)
		return string(bs)
	}
	return ""
}

func (entity *NodeModel) IsROOT() bool {
	return entity.IsLocal == 1
}

func (entity *NodeModel) GetTableName() string {
	return TableNode
}

func (entity *NodeModel) GetPKColumnName() string {
	return ""
}

// NetProxyModel 节点网络代理
type NetProxyModel struct {
	NetProxyId    int64     `json:"netProxyId,omitempty"`
	Name          string    `json:"name,omitempty"`
	Comment       string    `json:"comment,omitempty"`
	Code          string    `json:"code,omitempty"`
	InnerServerId string    `json:"innerServerId,omitempty"`
	InnerType     string    `json:"innerType,omitempty"`
	InnerAddress  string    `json:"innerAddress,omitempty"`
	OuterServerId string    `json:"outerServerId,omitempty"`
	OuterType     string    `json:"outerType,omitempty"`
	OuterAddress  string    `json:"outerAddress,omitempty"`
	LineServerIds string    `json:"lineServerIds,omitempty"`
	Option        string    `json:"option,omitempty"`
	UserId        int64     `json:"userId,omitempty"`
	Enabled       int8      `json:"enabled"`
	DeleteUserId  int64     `json:"deleteUserId,omitempty"`
	Deleted       int8      `json:"deleted"`
	CreateTime    time.Time `json:"createTime,omitempty"`
	UpdateTime    time.Time `json:"updateTime,omitempty"`
	DeleteTime    time.Time `json:"deleteTime,omitempty"`

	InnerStatus           int8     `json:"innerStatus"`
	InnerIsStarted        bool     `json:"innerIsStarted"`
	OuterStatus           int8     `json:"outerStatus"`
	OuterIsStarted        bool     `json:"outerIsStarted"`
	LineNodeIdList        []string `json:"lineNodeIdList,omitempty"`
	ReverseLineNodeIdList []string `json:"reverseLineNodeIdList,omitempty"`
}

func (entity *NetProxyModel) GetTableName() string {
	return TableNodeNetProxy
}

func (entity *NetProxyModel) GetPKColumnName() string {
	return ""
}
