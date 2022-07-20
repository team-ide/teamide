package module_node

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync"
	"teamide/node"
)

func (this_ *NodeService) InitContext() (err error) {
	if this_.nodeContext == nil {
		this_.nodeContext = &NodeContext{
			nodeService:            this_,
			wsCache:                make(map[string]*websocket.Conn),
			nodeIdNodeModelCache:   make(map[int64]*NodeModel),
			serverIdNodeModelCache: make(map[string]*NodeModel),
		}
	}
	err = this_.nodeContext.initContext()
	if err != nil {
		return
	}
	return
}

func (this_ *NodeService) addWS(id string, ws *websocket.Conn) (err error) {
	this_.nodeContext.setWS(id, ws)
	return
}

type Message struct {
	Method       string           `json:"method,omitempty"`
	NodeList     []*NodeInfo      `json:"nodeList,omitempty"`
	NetProxyList []*node.NetProxy `json:"netProxyList,omitempty"`
}
type NodeContext struct {
	server                     *node.Server
	nodeService                *NodeService
	nodeList                   []*NodeInfo
	root                       *NodeModel
	nodeIdNodeModelCache       map[int64]*NodeModel
	nodeIdNodeModelCacheLock   sync.Mutex
	serverIdNodeModelCache     map[string]*NodeModel
	serverIdNodeModelCacheLock sync.Mutex
	wsCache                    map[string]*websocket.Conn
	wsCacheLock                sync.Mutex
	netProxyList               []*node.NetProxy
}

func (this_ *NodeContext) getNodeModel(id int64) (res *NodeModel) {
	this_.nodeIdNodeModelCacheLock.Lock()
	defer this_.nodeIdNodeModelCacheLock.Unlock()
	res = this_.nodeIdNodeModelCache[id]
	return
}

func (this_ *NodeContext) setNodeModel(id int64, nodeModel *NodeModel) {
	this_.nodeIdNodeModelCacheLock.Lock()
	defer this_.nodeIdNodeModelCacheLock.Unlock()

	this_.nodeIdNodeModelCache[id] = nodeModel
}

func (this_ *NodeContext) removeNodeModel(id int64) {
	this_.nodeIdNodeModelCacheLock.Lock()
	defer this_.nodeIdNodeModelCacheLock.Unlock()
	delete(this_.nodeIdNodeModelCache, id)
}

func (this_ *NodeContext) getNodeModelByServerId(id string) (res *NodeModel) {
	this_.serverIdNodeModelCacheLock.Lock()
	defer this_.serverIdNodeModelCacheLock.Unlock()
	res = this_.serverIdNodeModelCache[id]
	return
}

func (this_ *NodeContext) setNodeModelByServerId(id string, nodeModel *NodeModel) {
	this_.serverIdNodeModelCacheLock.Lock()
	defer this_.serverIdNodeModelCacheLock.Unlock()

	this_.serverIdNodeModelCache[id] = nodeModel
}

func (this_ *NodeContext) removeNodeModelByServerId(id string) {
	this_.serverIdNodeModelCacheLock.Lock()
	defer this_.serverIdNodeModelCacheLock.Unlock()
	delete(this_.serverIdNodeModelCache, id)
}

func (this_ *NodeContext) onAddNodeModel(nodeModel *NodeModel) {
	if nodeModel == nil {
		return
	}
	this_.setNodeModel(nodeModel.NodeId, nodeModel)
	this_.setNodeModelByServerId(nodeModel.ServerId, nodeModel)
	var err error
	if nodeModel.IsROOT() {
		err = this_.initRoot(nodeModel)
		if err != nil {
			this_.nodeService.Logger.Error("node context init root error", zap.Error(err))
		}
	} else {
		_ = this_.server.AddNodeList([]*node.Info{
			{
				Id:          nodeModel.ServerId,
				BindToken:   nodeModel.BindToken,
				BindAddress: nodeModel.BindAddress,
				ConnAddress: nodeModel.ConnAddress,
				ConnToken:   nodeModel.ConnToken,
			},
		})
	}
}

func (this_ *NodeContext) onUpdateNodeModel(nodeModel *NodeModel) {
	if nodeModel == nil {
		return
	}
	var err error
	if nodeModel.IsROOT() {
		err = this_.initRoot(nodeModel)
		if err != nil {
			this_.nodeService.Logger.Error("node context init root error", zap.Error(err))
		}
	}
	this_.setNodeModel(nodeModel.NodeId, nodeModel)
	this_.setNodeModelByServerId(nodeModel.ServerId, nodeModel)
}

func (this_ *NodeContext) onRemoveNodeModel(id int64) {
	var nodeModel = this_.getNodeModel(id)
	if nodeModel == nil {
		return
	}
}

func (this_ *NodeContext) getWS(id string) (ws *websocket.Conn) {
	this_.wsCacheLock.Lock()
	defer this_.wsCacheLock.Unlock()
	ws = this_.wsCache[id]
	return
}

func (this_ *NodeContext) setWS(id string, ws *websocket.Conn) {
	this_.wsCacheLock.Lock()
	defer this_.wsCacheLock.Unlock()

	this_.wsCache[id] = ws
}

func (this_ *NodeContext) removeWS(id string) {
	this_.wsCacheLock.Lock()
	defer this_.wsCacheLock.Unlock()

	find, ok := this_.wsCache[id]
	if ok {
		_ = find.Close()
	}
	delete(this_.wsCache, id)

}

func (this_ *NodeContext) removeWSByWS(ws *websocket.Conn) {
	this_.wsCacheLock.Lock()
	defer this_.wsCacheLock.Unlock()
	var idList []string
	for id, one := range this_.wsCache {
		if one == ws {
			idList = append(idList, id)
		}
	}
	for _, id := range idList {
		find, ok := this_.wsCache[id]
		if ok {
			_ = find.Close()
		}
		delete(this_.wsCache, id)
	}
}

func (this_ *NodeContext) getWSList() (list []*websocket.Conn) {
	this_.wsCacheLock.Lock()
	defer this_.wsCacheLock.Unlock()
	for _, one := range this_.wsCache {
		list = append(list, one)
	}
	return
}

type NodeInfo struct {
	Info      *node.Info `json:"info,omitempty"`
	NodeModel *NodeModel `json:"nodeModel,omitempty"`
}

func (this_ *NodeContext) initContext() (err error) {
	var list []*NodeModel
	list, _ = this_.nodeService.Query(&NodeModel{})
	for _, one := range list {
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
	return
}

func (this_ *NodeContext) initRoot(root *NodeModel) (err error) {
	if this_.server != nil {
		this_.server.Stop()
	}
	this_.root = root
	this_.server = &node.Server{
		Id:               this_.root.ServerId,
		BindToken:        this_.root.BindToken,
		BindAddress:      this_.root.BindAddress,
		OnNodeListChange: this_.onNodeListChange,
	}
	err = this_.server.Start()
	if err != nil {
		return
	}
	_ = this_.server.AddNodeList([]*node.Info{
		{
			Id:          this_.root.ServerId,
			BindToken:   this_.root.BindToken,
			BindAddress: this_.root.BindAddress,
			ConnAddress: this_.root.ConnAddress,
			ConnToken:   this_.root.ConnToken,
		},
	})
	return
}
func (this_ *NodeContext) callMessage(msg *Message) {
	bs, err := json.Marshal(msg)
	if err != nil {
		return
	}
	var list = this_.getWSList()
	for _, one := range list {
		err = one.WriteMessage(websocket.TextMessage, bs)
		if err != nil {
			this_.removeWSByWS(one)
		}
	}
}

func (this_ *NodeContext) onNodeListChange(nodeList []*node.Info) {

	var nodeInfoList []*NodeInfo
	for _, one := range nodeList {
		nodeInfo := &NodeInfo{
			Info:      one,
			NodeModel: this_.getNodeModelByServerId(one.Id),
		}
		nodeInfoList = append(nodeInfoList, nodeInfo)
	}
	this_.nodeList = nodeInfoList
	this_.callMessage(&Message{
		Method:   "refresh_node_list",
		NodeList: this_.nodeList,
	})
}

func (this_ *NodeContext) onNetProxyListChange(netProxyList []*node.NetProxy) {
	this_.netProxyList = netProxyList
	this_.callMessage(&Message{
		Method:       "refresh_net_proxy_list",
		NetProxyList: this_.netProxyList,
	})
}
