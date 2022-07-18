package node

import (
	"bytes"
	"fmt"
)

var (
	StatusStarted = 1
	StatusStopped = 2
	StatusError   = 3
)

type Info struct {
	Id          string `json:"id,omitempty"`
	ConnAddress string `json:"connAddress,omitempty"`
	Token       string `json:"token,omitempty"`
	ParentId    string `json:"parentId,omitempty"`
	ConnSize    int    `json:"connSize,omitempty"`
	Status      int    `json:"status,omitempty"`
	StatusError string `json:"statusError,omitempty"`
}

func (this_ *Info) GetNodeStr() (str string) {
	return fmt.Sprintf("节点[%s][%s]", this_.Id, this_.ConnAddress)
}

func (this_ *Info) checkToken(token []byte) bool {
	nodeToken := []byte(this_.Token)
	if len(nodeToken) != len(token) {
		Logger.Error(this_.GetNodeStr() + " Token check field")
		return false
	}
	if !bytes.Contains(token, nodeToken) {
		Logger.Error(this_.GetNodeStr() + " Token check field")
		return false
	}
	return true
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
	target.ConnAddress = source.ConnAddress
	target.Token = source.Token
	target.ParentId = source.ParentId
}
