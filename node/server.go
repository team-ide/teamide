package node

import (
	"fmt"
	"net"
	"sync"
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
