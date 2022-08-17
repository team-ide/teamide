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
	workerCache                   map[string]*Worker
	workerCacheLock               sync.Mutex
	OnNodeStatusChange            func(id string, status int8)
	OnNetProxyInnerChange         func(id string, status int8)
	OnNetProxyOuterChange         func(id string, status int8)
	connNodeListenerKeepAliveLock sync.Mutex
}

func (this_ *Server) Start() (err error) {
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

func (this_ *Server) getWorkerIfAbsentCreate(spaceId string) (worker *Worker) {
	this_.workerCacheLock.Lock()
	defer this_.workerCacheLock.Unlock()

	worker, ok := this_.workerCache[spaceId]
	if !ok {
		worker = &Worker{
			server:      this_,
			space:       newSpace(spaceId),
			MonitorData: &MonitorData{},
		}
		this_.workerCache[spaceId] = worker
	}
	return
}

func (this_ *Server) getWorker(spaceId string) (worker *Worker) {
	this_.workerCacheLock.Lock()
	defer this_.workerCacheLock.Unlock()

	worker, _ = this_.workerCache[spaceId]
	return
}

func (this_ *Server) GetServerInfo() (str string) {
	return fmt.Sprintf("节点服务[%s][%s]", this_.Id, this_.BindAddress)
}

func (this_ *Server) GetNode(spaceId string, lineNodeIdList []string, nodeId string) (node *Info) {
	worker := this_.getWorker(spaceId)
	if worker == nil {
		return
	}
	node = worker.getNode(lineNodeIdList, nodeId)
	return
}

func (this_ *Server) SystemGetInfo(spaceId string, lineNodeIdList []string, nodeId string) (info *system.Info) {
	worker := this_.getWorker(spaceId)
	if worker == nil {
		return
	}

	res := worker.systemGetInfo(lineNodeIdList, nodeId)
	if res != nil {
		info = res.Info
	}
	return
}

func (this_ *Server) SystemQueryMonitorData(spaceId string, lineNodeIdList []string, nodeId string, request *system.QueryRequest) (response *system.QueryResponse) {
	worker := this_.getWorker(spaceId)
	if worker == nil {
		return
	}
	res := worker.systemQueryMonitorData(lineNodeIdList, nodeId, &SystemData{
		QueryRequest: request,
	})
	if res != nil {
		response = res.QueryResponse
	}
	return
}

func (this_ *Server) SystemCleanMonitorData(spaceId string, lineNodeIdList []string, nodeId string) {
	worker := this_.getWorker(spaceId)
	if worker == nil {
		return
	}
	_ = worker.systemCleanMonitorData(lineNodeIdList, nodeId)
	return
}

func (this_ *Server) GetNodeMonitorData(spaceId string, lineNodeIdList []string, nodeId string) (monitorData *MonitorData) {
	worker := this_.getWorker(spaceId)
	if worker == nil {
		return
	}
	monitorData = worker.getNodeMonitorData(lineNodeIdList, nodeId)
	return
}

func (this_ *Server) GetNetProxyInnerMonitorData(spaceId string, lineNodeIdList []string, netProxyId string) (monitorData *MonitorData) {
	worker := this_.getWorker(spaceId)
	if worker == nil {
		return
	}
	monitorData = worker.getNetProxyInnerMonitorData(lineNodeIdList, netProxyId)
	return
}

func (this_ *Server) GetNetProxyOuterMonitorData(spaceId string, lineNodeIdList []string, netProxyId string) (monitorData *MonitorData) {
	worker := this_.getWorker(spaceId)
	if worker == nil {
		return
	}
	monitorData = worker.getNetProxyOuterMonitorData(lineNodeIdList, netProxyId)
	return
}

func (this_ *Server) AddNodeList(spaceId string, nodeList []*Info) (err error) {
	worker := this_.getWorker(spaceId)
	if worker == nil {
		return
	}
	err = worker.addNodeList(nodeList)
	return
}

func (this_ *Server) UpdateNodeConnNodeIdList(spaceId string, id string, connNodeIdList []string) (err error) {
	worker := this_.getWorker(spaceId)
	if worker == nil {
		return
	}
	err = worker.updateNodeConnNodeIdList(id, connNodeIdList)
	return
}

func (this_ *Server) RemoveNodeList(spaceId string, nodeIdList []string) (err error) {
	worker := this_.getWorker(spaceId)
	if worker == nil {
		return
	}
	err = worker.removeNodeList(nodeIdList)
	return
}

func (this_ *Server) AddNetProxyList(spaceId string, netProxyList []*NetProxy) (err error) {
	worker := this_.getWorker(spaceId)
	if worker == nil {
		return
	}
	err = worker.addNetProxyList(netProxyList)
	return
}

func (this_ *Server) RemoveNetProxyList(spaceId string, netProxyIdList []string) (err error) {
	worker := this_.getWorker(spaceId)
	if worker == nil {
		return
	}
	err = worker.removeNetProxyList(netProxyIdList)
	return
}
