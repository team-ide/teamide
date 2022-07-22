package module_node

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync"
	"teamide/node"
	"teamide/pkg/util"
)

func (this_ *NodeService) InitContext() {
	if this_.nodeContext == nil {
		this_.nodeContext = &NodeContext{
			nodeService:            this_,
			wsCache:                make(map[string]*WSConn),
			nodeIdNodeModelCache:   make(map[int64]*NodeModel),
			serverIdNodeModelCache: make(map[string]*NodeModel),
		}
	}
	err := this_.nodeContext.initContext()
	if err != nil {
		this_.Logger.Error("节点上下文初始化异常", zap.Error(err))
		return
	}
	return
}

type WSConn struct {
	id       string
	conn     *websocket.Conn
	connLock sync.Mutex
}

func (this_ *WSConn) Close() {
	_ = this_.conn.Close()
	return
}

func (this_ *WSConn) WriteMessage(bytes []byte) (err error) {
	this_.connLock.Lock()
	defer this_.connLock.Unlock()

	err = this_.conn.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		return
	}
	return
}

func (this_ *NodeService) addWS(id string, ws *websocket.Conn) (err error) {
	wsConn := &WSConn{
		id:   id,
		conn: ws,
	}
	this_.nodeContext.setWS(wsConn)
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
	wsCache                    map[string]*WSConn
	wsCacheLock                sync.Mutex
	netProxyList               []*node.NetProxy
	nodeListLock               sync.Mutex
	nodeListChangeLock         sync.Mutex
}

func (this_ *NodeContext) getNodeInfo(id string) (res *NodeInfo) {
	this_.nodeListLock.Lock()
	defer this_.nodeListLock.Unlock()
	var list = this_.nodeList
	for _, one := range list {
		if one.Info != nil && id == one.Info.Id {
			res = one
			return
		}
	}
	return
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
	}
	var connNodeIdList []string
	if nodeModel.ConnServerIds != "" {
		_ = json.Unmarshal([]byte(nodeModel.ConnServerIds), &connNodeIdList)
	}

	_ = this_.server.AddNodeList([]*node.Info{
		{
			Id:             nodeModel.ServerId,
			BindToken:      nodeModel.BindToken,
			BindAddress:    nodeModel.BindAddress,
			ConnAddress:    nodeModel.ConnAddress,
			ConnToken:      nodeModel.ConnToken,
			ConnNodeIdList: connNodeIdList,
		},
	})
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

func (this_ *NodeContext) onUpdateNodeConnServerIds(nodeId int64, connServerIds string) {
	nodeModel := this_.getNodeModel(nodeId)
	if nodeModel == nil {
		return
	}
	nodeModel.ConnServerIds = connServerIds

	var connNodeIdList []string
	if connServerIds != "" {
		_ = json.Unmarshal([]byte(connServerIds), &connNodeIdList)
	}
	_ = this_.server.UpdateNodeConnNodeIdList(nodeModel.ServerId, connNodeIdList)
}

func (this_ *NodeContext) onUpdateNodeHistoryConnServerIds(nodeId int64, historyConnServerIds string) {
	nodeModel := this_.getNodeModel(nodeId)
	if nodeModel == nil {
		return
	}
	nodeModel.HistoryConnServerIds = historyConnServerIds

}

func (this_ *NodeContext) onRemoveNodeModel(id int64) {
	var nodeModel = this_.getNodeModel(id)
	if nodeModel == nil {
		return
	}
	this_.removeNodeModel(nodeModel.NodeId)
	this_.removeNodeModelByServerId(nodeModel.ServerId)
	_ = this_.server.RemoveNodeList([]string{nodeModel.ServerId})
}

func (this_ *NodeContext) getWS(id string) (ws *WSConn) {
	this_.wsCacheLock.Lock()
	defer this_.wsCacheLock.Unlock()
	ws = this_.wsCache[id]
	return
}

func (this_ *NodeContext) setWS(ws *WSConn) {
	this_.wsCacheLock.Lock()
	defer this_.wsCacheLock.Unlock()

	this_.wsCache[ws.id] = ws
}

func (this_ *NodeContext) removeWS(id string) {
	this_.wsCacheLock.Lock()
	defer this_.wsCacheLock.Unlock()

	find, ok := this_.wsCache[id]
	if ok {
		find.Close()
	}
	delete(this_.wsCache, id)

}

func (this_ *NodeContext) getWSList() (list []*WSConn) {
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

	return
}
func (this_ *NodeContext) callMessage(msg *Message) {
	bs, err := json.Marshal(msg)
	if err != nil {
		return
	}
	var list = this_.getWSList()
	for _, one := range list {
		err = one.WriteMessage(bs)
		if err != nil {
			this_.removeWS(one.id)
		}
	}
}

func (this_ *NodeContext) onNodeListChange(nodeList []*node.Info) {

	this_.nodeListChangeLock.Lock()
	defer this_.nodeListChangeLock.Unlock()

	var nodeInfoList []*NodeInfo
	for _, one := range nodeList {
		var find = this_.getNodeModelByServerId(one.Id)
		var historyConnServerIdList []string
		if find != nil && find.HistoryConnServerIds != "" {
			_ = json.Unmarshal([]byte(find.HistoryConnServerIds), &historyConnServerIdList)
		}
		for _, one := range one.ConnNodeIdList {
			if util.ContainsString(historyConnServerIdList, one) < 0 {
				historyConnServerIdList = append(historyConnServerIdList, one)
			}
		}
		var historyConnServerIds string
		bs, _ := json.Marshal(historyConnServerIdList)
		if bs != nil {
			historyConnServerIds = string(bs)
		}
		if find == nil {
			find = &NodeModel{
				ServerId:             one.Id,
				Name:                 one.Id,
				BindToken:            one.BindToken,
				BindAddress:          one.BindAddress,
				ConnToken:            one.ConnToken,
				ConnAddress:          one.ConnAddress,
				HistoryConnServerIds: historyConnServerIds,
			}
			_, err := this_.nodeService.Insert(find)
			if err != nil {
				find = nil
			} else {
				this_.setNodeModel(find.NodeId, find)
				this_.setNodeModelByServerId(one.Id, find)
			}
		} else {
			_, _ = this_.nodeService.UpdateHistoryConnServerIds(find.NodeId, historyConnServerIds)
		}
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
