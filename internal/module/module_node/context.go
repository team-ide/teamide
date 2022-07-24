package module_node

import (
	"go.uber.org/zap"
	"sync"
	"teamide/node"
)

func (this_ *NodeService) InitContext() {
	if this_.nodeContext == nil {
		this_.nodeContext = &NodeContext{
			nodeService: this_,
			wsCache:     make(map[string]*WSConn),

			nodeIdModelCache:   make(map[int64]*NodeModel),
			serverIdModelCache: make(map[string]*NodeModel),

			netProxyIdModelCache: make(map[int64]*NetProxyModel),
			codeModelCache:       make(map[string]*NetProxyModel),
		}
	}
	err := this_.nodeContext.initContext()
	if err != nil {
		this_.Logger.Error("节点上下文初始化异常", zap.Error(err))
		return
	}
	return
}

type NodeContext struct {
	server      *node.Server
	nodeService *NodeService
	root        *NodeModel

	wsCache     map[string]*WSConn
	wsCacheLock sync.Mutex

	nodeList               []*NodeInfo
	nodeListLock           sync.Mutex
	nodeIdModelCache       map[int64]*NodeModel
	nodeIdModelCacheLock   sync.Mutex
	serverIdModelCache     map[string]*NodeModel
	serverIdModelCacheLock sync.Mutex

	netProxyList             []*NetProxyInfo
	netProxyListLock         sync.Mutex
	netProxyIdModelCache     map[int64]*NetProxyModel
	netProxyIdModelCacheLock sync.Mutex
	codeModelCache           map[string]*NetProxyModel
	codeModelCacheLock       sync.Mutex
}

type NetProxyInfo struct {
	Info           *node.NetProxy `json:"info,omitempty"`
	Model          *NetProxyModel `json:"model,omitempty"`
	InnerIsStarted bool           `json:"innerIsStarted,omitempty"`
	OuterIsStarted bool           `json:"outerIsStarted,omitempty"`
}

func (this_ *NodeContext) initContext() (err error) {
	var list []*NodeModel
	list, _ = this_.nodeService.Query(&NodeModel{})
	for _, one := range list {
		this_.setNodeModel(one.NodeId, one)
		this_.setNodeModelByServerId(one.ServerId, one)
		if one.IsROOT() {
			this_.root = one
		}
	}

	if this_.root != nil {
		this_.onAddNodeModel(this_.root)
	}
	for _, one := range list {
		if !one.IsROOT() {
			this_.onAddNodeModel(one)
		}
	}

	var netProxyList []*NetProxyModel
	netProxyList, _ = this_.nodeService.QueryNetProxy(&NetProxyModel{})
	for _, one := range netProxyList {
		this_.onAddNetProxyModel(one)
	}
	return
}

func (this_ *NodeContext) initRoot(root *NodeModel) (err error) {
	if this_.server != nil {
		this_.server.Stop()
	}
	this_.root = root
	this_.server = &node.Server{
		Id:                   this_.root.ServerId,
		BindToken:            this_.root.BindToken,
		BindAddress:          this_.root.BindAddress,
		OnNodeListChange:     this_.onNodeListChange,
		OnNetProxyListChange: this_.onNetProxyListChange,
	}
	err = this_.server.Start()
	if err != nil {
		return
	}

	return
}
