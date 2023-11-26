package node

import (
	"fmt"
	"net"
	"sync"
	"teamide/pkg/system"
)

var tokenByteSize = 128

type LocalNode struct {
	Id             string `json:"id"`
	BindAddress    string `json:"bindAddress"`
	BindToken      string `json:"-"`
	ConnAddress    string `json:"connAddress"`
	ConnToken      string `json:"-"`
	ConnSize       int    `json:"connSize"`
	IsStop         bool   `json:"isStop"`
	serverListener net.Listener
}

func (this_ *LocalNode) GetServerInfo() (str string) {
	return fmt.Sprintf("节点服务[%s][%s]", this_.Id, this_.BindAddress)
}

type Server struct {
	localNodeList []*LocalNode
	serverInfo    string

	OnNodeStatusChange    func(id string, status int8)
	OnNetProxyInnerChange func(id string, status int8)
	OnNetProxyOuterChange func(id string, status int8)

	connNodeListenerKeepAliveLock sync.Mutex
	*Worker
}

func (this_ *Server) GetLocalNodeIdList() (localNodeIdList []string) {

	for _, one := range this_.localNodeList {
		localNodeIdList = append(localNodeIdList, one.Id)
	}
	return
}

func (this_ *Server) AddLocalNode(localNode *LocalNode) {
	this_.localNodeList = append(this_.localNodeList, localNode)

	var serverInfo string

	for _, one := range this_.localNodeList {
		serverInfo += fmt.Sprintf("节点服务[%s][%s] ", one.Id, one.BindAddress)
	}
	this_.serverInfo = serverInfo
	if localNode.BindAddress != "" {
		go this_.serverListenerKeepAlive(localNode)
	}

	if localNode.ConnAddress != "" {
		this_.connNodeListenerKeepAlive(localNode.ConnAddress, localNode.ConnToken, localNode.ConnSize)
	}
}

func (this_ *Server) RemoveLocalNode(id string) {
	var newList []*LocalNode
	for _, one := range this_.localNodeList {
		if one.Id != id {
			newList = append(newList, one)
		} else {
			one.IsStop = true
			if one.serverListener != nil {
				_ = one.serverListener.Close()
			}
		}
	}
	this_.localNodeList = newList

	var serverInfo string

	for _, one := range this_.localNodeList {
		serverInfo += fmt.Sprintf("节点服务[%s][%s] ", one.Id, one.BindAddress)
	}
	this_.serverInfo = serverInfo

}
func (this_ *Server) Start() {
	this_.Worker = &Worker{
		server:      this_,
		Space:       newSpace(),
		MonitorData: &MonitorData{},
	}
	go system.StartCollectMonitorData()
	return
}

func (this_ *Server) Stop() {
	for _, one := range this_.localNodeList {
		one.IsStop = true
		if one.serverListener != nil {
			_ = one.serverListener.Close()
		}
	}
	if this_.Worker != nil {
		this_.Worker.Stop()
	}

}

func (this_ *Server) GetServerInfo() (str string) {
	return this_.serverInfo
}

func (this_ *Server) SystemGetInfo(lineNodeIdList []string) (info *system.Info) {
	res := this_.systemGetInfo(lineNodeIdList)
	if res != nil {
		info = res.Info
	}
	return
}

func (this_ *Server) SystemMonitorData(lineNodeIdList []string) (monitorData *system.MonitorData) {
	res := this_.systemMonitorData(lineNodeIdList)
	if res != nil {
		monitorData = res.MonitorData
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
