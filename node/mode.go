package node

import (
	"fmt"
	"sync"
	"teamide/pkg/util"
)

var (
	StatusStarted int8 = 1
	StatusStopped int8 = 2
	StatusError   int8 = 3
)

type Info struct {
	Id             string   `json:"id,omitempty"`
	BindAddress    string   `json:"bindAddress,omitempty"`
	BindToken      string   `json:"bindToken,omitempty"`
	ConnAddress    string   `json:"connAddress,omitempty"`
	ConnToken      string   `json:"connToken,omitempty"`
	ConnNodeIdList []string `json:"connNodeIdList,omitempty"`
	ConnSize       int      `json:"connSize,omitempty"`
	Status         int8     `json:"status,omitempty"`
	StatusError    string   `json:"statusError,omitempty"`
	Enabled        int8     `json:"enabled,omitempty"`
	connIdListLock sync.Mutex
}

func (this_ *Info) IsEnabled() bool {
	return this_.Enabled != 2
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
	InnerStatus           int8       `json:"innerStatus,omitempty"`
	InnerStatusError      string     `json:"innerStatusError,omitempty"`
	OuterStatus           int8       `json:"outerStatus,omitempty"`
	OuterStatusError      string     `json:"outerStatusError,omitempty"`
	Enabled               int8       `json:"enabled,omitempty"`
}

func (this_ *NetProxy) IsEnabled() bool {
	return this_.Enabled != 2
}

type NetConfig struct {
	NodeId  string `json:"nodeId,omitempty"`
	Type    string `json:"type,omitempty"`
	Address string `json:"address,omitempty"`
}

func (this_ *NetConfig) GetInfoStr() (str string) {
	return fmt.Sprintf("[%s][%s]", this_.GetType(), this_.Address)
}

func (this_ *NetConfig) GetType() (str string) {
	var t = this_.Type
	if t == "" {
		t = "tcp"
	}
	return t
}

func (this_ *NetConfig) GetAddress() (str string) {
	return GetAddress(this_.Address)
}

func GetAddress(address string) (str string) {
	if address == "" {
		return ""
	}
	return address
}

func copyNode(source, target *Info) (hasChange bool) {
	if source.BindAddress != "" {
		target.BindAddress = source.BindAddress
	}
	if source.BindToken != "" {
		target.BindToken = source.BindToken
	}
	if source.ConnAddress != "" {
		target.ConnAddress = source.ConnAddress
	}
	if source.ConnToken != "" {
		target.ConnToken = source.ConnToken
	}
	if source.Status != 0 {
		if target.Status != source.Status || target.StatusError != source.StatusError {
			hasChange = true
			target.Status = source.Status
			target.StatusError = source.StatusError
		}
	}
	if source.Enabled != 0 {
		if target.IsEnabled() != source.IsEnabled() {
			hasChange = true
			target.Enabled = source.Enabled
		}
	}
	var list = source.ConnNodeIdList
	for _, one := range list {
		if util.ContainsString(target.ConnNodeIdList, one) < 0 {
			hasChange = true
			target.addConnNodeId(one)
		}
	}
	return
}

type MonitorData struct {
	ReadSize  int64 `json:"readSize,omitempty"`
	ReadTime  int64 `json:"readTime,omitempty"`
	readLock  sync.Mutex
	WriteSize int64 `json:"writeSize,omitempty"`
	WriteTime int64 `json:"writeTime,omitempty"`
	writeLock sync.Mutex
}

func (this_ *MonitorData) monitorRead(bytesSize int64, useTime int64) {
	this_.readLock.Lock()
	defer this_.readLock.Unlock()

	this_.ReadSize += bytesSize
	this_.ReadTime += useTime
}

func (this_ *MonitorData) monitorWrite(bytesSize int64, useTime int64) {
	this_.writeLock.Lock()
	defer this_.writeLock.Unlock()

	this_.WriteSize += bytesSize
	this_.WriteTime += useTime
}
