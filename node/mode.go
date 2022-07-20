package node

import (
	"fmt"
	"sync"
	"teamide/pkg/util"
)

var (
	StatusStarted = 1
	StatusStopped = 2
	StatusError   = 3
)

type Info struct {
	Id             string   `json:"id,omitempty"`
	BindAddress    string   `json:"bindAddress,omitempty"`
	BindToken      string   `json:"bindToken,omitempty"`
	ConnAddress    string   `json:"connAddress,omitempty"`
	ConnToken      string   `json:"connToken,omitempty"`
	ConnNodeIdList []string `json:"connNodeIdList,omitempty"`
	ConnSize       int      `json:"connSize,omitempty"`
	Status         int      `json:"status,omitempty"`
	StatusError    string   `json:"statusError,omitempty"`
	connIdListLock sync.Mutex
}

func (this_ *Info) GetNodeStr() (str string) {
	return fmt.Sprintf("节点[%s][%s]", this_.Id, this_.ConnAddress)
}

func (this_ *Info) addConnNodeId(connNodeId string) {
	this_.connIdListLock.Lock()
	defer this_.connIdListLock.Unlock()

	if connNodeId == "" || connNodeId == this_.Id {
		return
	}
	if util.ContainsString(this_.ConnNodeIdList, connNodeId) < 0 {
		this_.ConnNodeIdList = append(this_.ConnNodeIdList, connNodeId)
	}
	return
}

func (this_ *Info) removeConnNodeId(connNodeId string) {
	this_.connIdListLock.Lock()
	defer this_.connIdListLock.Unlock()

	var list = this_.ConnNodeIdList
	var newList []string
	for _, one := range list {
		if one != connNodeId {
			newList = append(newList, one)
		}
	}
	this_.ConnNodeIdList = newList
	return
}

type NetProxy struct {
	Id                    string     `json:"id,omitempty"`
	Inner                 *NetConfig `json:"inner,omitempty"`
	Outer                 *NetConfig `json:"outer,omitempty"`
	LineNodeIdList        []string   `json:"lineNodeIdList,omitempty"`
	ReverseLineNodeIdList []string   `json:"reverseLineNodeIdList,omitempty"`
}

type NetConfig struct {
	NodeId  string `json:"nodeId,omitempty"`
	Network string `json:"network,omitempty"`
	Address string `json:"address,omitempty"`
}

func (this_ *NetConfig) GetInfoStr() (str string) {
	return fmt.Sprintf("[%s][%s]", this_.Network, this_.Address)
}

func (this_ *NetConfig) GetNetwork() (str string) {
	return GetNetwork(this_.Network)
}

func (this_ *NetConfig) GetAddress() (str string) {
	return GetAddress(this_.Address)
}

func GetNetwork(network string) (str string) {
	if network == "" {
		return "tcp"
	}
	return network
}

func GetAddress(address string) (str string) {
	if address == "" {
		return ""
	}
	return address
}

func copyNode(source, target *Info) {
	target.Id = source.Id
	target.BindAddress = source.BindAddress
	target.BindToken = source.BindToken
	target.ConnAddress = source.ConnAddress
	target.ConnToken = source.ConnToken
	target.Status = source.Status
	target.StatusError = source.StatusError
	var list = source.ConnNodeIdList
	for _, one := range list {
		target.addConnNodeId(one)
	}
}
