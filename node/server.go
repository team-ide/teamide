package node

import (
	"fmt"
	"net"
	"sync"
	"teamide/pkg/system"
)

var tokenByteSize = 128

type Server struct {
	Id             string
	BindAddress    string
	BindToken      string
	ConnAddress    string
	ConnToken      string
	ConnSize       int
	serverListener net.Listener

	OnNodeStatusChange    func(id string, status int8)
	OnNetProxyInnerChange func(id string, status int8)
	OnNetProxyOuterChange func(id string, status int8)

	connNodeListenerKeepAliveLock sync.Mutex
	*Worker
}

func (this_ *Server) Start() (err error) {
	this_.Worker = &Worker{
		server:      this_,
		Space:       newSpace(),
		MonitorData: &MonitorData{},
	}

	if this_.BindAddress != "" {
		go this_.serverListenerKeepAlive()
	}
	if this_.ConnAddress != "" {
		this_.connNodeListenerKeepAlive(this_.ConnAddress, this_.ConnToken, this_.ConnSize)
	}
	go system.StartCollectMonitorData()
	return
}

func (this_ *Server) Stop() {
	if this_.serverListener != nil {
		_ = this_.serverListener.Close()
	}
	if this_.Worker != nil {
		this_.Worker.Stop()
	}

}

func (this_ *Server) GetServerInfo() (str string) {
	return fmt.Sprintf("节点服务[%s][%s]", this_.Id, this_.BindAddress)
}

func (this_ *Server) SystemGetInfo(lineNodeIdList []string) (info *system.Info) {
	res := this_.systemGetInfo(lineNodeIdList)
	if res != nil {
		info = res.Info
	}
	return
}

func (this_ *Server) SystemQueryMonitorData(lineNodeIdList []string, request *system.QueryRequest) (response *system.QueryResponse) {
	res := this_.systemQueryMonitorData(lineNodeIdList, &SystemData{
		QueryRequest: request,
	})
	if res != nil {
		response = res.QueryResponse
	}
	return
}

func (this_ *Server) SystemCleanMonitorData(lineNodeIdList []string) {
	_ = this_.systemCleanMonitorData(lineNodeIdList)
	return
}

func (this_ *Server) GetNodeVersion(lineNodeIdList []string) (version string) {
	version = this_.getVersion(lineNodeIdList)
	return
}

func (this_ *Server) GetNodeStatus(lineNodeIdList []string) (status int8) {
	status = this_.getNodeStatus(lineNodeIdList)
	return
}

func (this_ *Server) GetNodeMonitorData(lineNodeIdList []string) (monitorData *MonitorData) {
	monitorData = this_.getNodeMonitorData(lineNodeIdList)
	return
}

func (this_ *Server) GetNetProxyInnerMonitorData(lineNodeIdList []string, netProxyId string) (monitorData *MonitorData) {
	monitorData = this_.getNetProxyInnerMonitorData(lineNodeIdList, netProxyId)
	return
}

func (this_ *Server) GetNetProxyOuterMonitorData(lineNodeIdList []string, netProxyId string) (monitorData *MonitorData) {
	monitorData = this_.getNetProxyOuterMonitorData(lineNodeIdList, netProxyId)
	return
}

func (this_ *Server) AddToNodeList(lineNodeIdList []string, toNodeList []*ToNode) (err error) {
	this_.addToNodeList(lineNodeIdList, toNodeList)
	return
}

func (this_ *Server) RemoveToNodeList(lineNodeIdList []string, toNodeIdList []string) (err error) {
	this_.removeToNodeList(lineNodeIdList, toNodeIdList)
	return
}

func (this_ *Server) GetNetProxyInnerStatus(lineNodeIdList []string, netProxyId string) (status int8) {
	status = this_.getNetProxyInnerStatus(lineNodeIdList, netProxyId)
	return
}

func (this_ *Server) AddNetProxyInnerList(lineNodeIdList []string, netProxyList []*NetProxyInner) (err error) {
	this_.addNetProxyInnerList(lineNodeIdList, netProxyList)
	return
}

func (this_ *Server) RemoveNetProxyInnerList(lineNodeIdList []string, netProxyIdList []string) (err error) {
	this_.removeNetProxyInnerList(lineNodeIdList, netProxyIdList)
	return
}

func (this_ *Server) GetNetProxyOuterStatus(lineNodeIdList []string, netProxyId string) (status int8) {
	status = this_.getNetProxyOuterStatus(lineNodeIdList, netProxyId)
	return
}

func (this_ *Server) AddNetProxyOuterList(lineNodeIdList []string, netProxyList []*NetProxyOuter) (err error) {
	this_.addNetProxyOuterList(lineNodeIdList, netProxyList)
	return
}

func (this_ *Server) RemoveNetProxyOuterList(lineNodeIdList []string, netProxyIdList []string) (err error) {
	this_.removeNetProxyOuterList(lineNodeIdList, netProxyIdList)
	return
}
