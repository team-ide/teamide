package node

import (
	"fmt"
	"net"
	"sync"
	"teamide/pkg/system"
)

var tokenByteSize = 128

type Server struct {
	Id                            string
	BindAddress                   string
	BindToken                     string
	ConnAddress                   string
	ConnToken                     string
	ConnSize                      int
	serverListener                net.Listener
	cache                         *Cache
	worker                        *Worker
	OnNodeListChange              func([]*Info)
	OnNetProxyListChange          func([]*NetProxy)
	rootNode                      *Info
	connNodeListenerKeepAliveLock sync.Mutex
}

func (this_ *Server) GetNode(nodeId string) (node *Info) {
	node = this_.worker.getNode(nodeId, []string{})
	return
}

func (this_ *Server) SystemGetInfo(nodeId string) (info *system.Info) {
	res := this_.worker.systemGetInfo(nodeId, []string{})
	if res != nil {
		info = res.Info
	}
	return
}

func (this_ *Server) SystemQueryMonitorData(nodeId string, request *system.QueryRequest) (response *system.QueryResponse) {
	res := this_.worker.systemQueryMonitorData(nodeId, []string{}, &SystemData{
		QueryRequest: request,
	})
	if res != nil {
		response = res.QueryResponse
	}
	return
}

func (this_ *Server) SystemCleanMonitorData(nodeId string) {
	_ = this_.worker.systemCleanMonitorData(nodeId, []string{})
	return
}

func (this_ *Server) GetNodeMonitorData(nodeId string) (monitorData *MonitorData) {
	monitorData = this_.worker.getNodeMonitorData(nodeId, []string{})
	return
}

func (this_ *Server) GetNetProxyInnerMonitorData(netProxyId string) (monitorData *MonitorData) {
	monitorData = this_.worker.getNetProxyInnerMonitorData(netProxyId, []string{})
	return
}

func (this_ *Server) GetNetProxyOuterMonitorData(netProxyId string) (monitorData *MonitorData) {
	monitorData = this_.worker.getNetProxyOuterMonitorData(netProxyId, []string{})
	return
}

func (this_ *Server) AddNodeList(nodeList []*Info) (err error) {
	err = this_.worker.addNodeList(nodeList)
	return
}

func (this_ *Server) UpdateNodeConnNodeIdList(id string, connNodeIdList []string) (err error) {
	err = this_.worker.updateNodeConnNodeIdList(id, connNodeIdList)
	return
}

func (this_ *Server) RemoveNodeList(nodeIdList []string) (err error) {
	err = this_.worker.removeNodeList(nodeIdList)
	return
}

func (this_ *Server) GetNodeLineByFromTo(fromNodeId, toNodeId string, nodeIdConnNodeIdListCache map[string][]string) (lineIdList []string) {

	return this_.worker.getNodeLineByFromTo(fromNodeId, toNodeId, nodeIdConnNodeIdListCache)
}

func (this_ *Server) AddNetProxyList(netProxyList []*NetProxy) (err error) {
	err = this_.worker.addNetProxyList(netProxyList)
	return
}

func (this_ *Server) RemoveNetProxyList(netProxyIdList []string) (err error) {
	err = this_.worker.removeNetProxyList(netProxyIdList)
	return
}

func (this_ *Server) Start() (err error) {
	this_.rootNode = &Info{
		Id:          this_.Id,
		BindAddress: this_.BindAddress,
		BindToken:   this_.BindToken,
	}
	this_.cache = newCache()
	this_.worker = &Worker{
		server:      this_,
		cache:       this_.cache,
		MonitorData: &MonitorData{},
	}
	_ = this_.worker.doAddNodeList([]*Info{this_.rootNode})

	if this_.BindAddress != "" {
		go this_.serverListenerKeepAlive()
	}
	if this_.ConnAddress != "" {
		this_.connNodeListenerKeepAlive(nil, this_.ConnAddress, this_.ConnToken, this_.ConnSize)
	}
	go system.StartCollectMonitorData()
	return
}

func (this_ *Server) Stop() {
	if this_.serverListener != nil {
		_ = this_.serverListener.Close()
	}
	if this_.worker != nil {
		this_.worker.Stop()
	}

}

func (this_ *Server) GetServerInfo() (str string) {
	return fmt.Sprintf("节点服务[%s][%s]", this_.Id, this_.BindAddress)
}
