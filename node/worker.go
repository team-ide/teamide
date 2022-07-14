package node

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	Node                    *Info
	childrenNodeList        []*Info
	server                  *Server
	isStop                  bool
	nodeList                []*Info
	childrenNodeClientCache map[string]*Client
	inPortForwardingList    []*PortForwarding
	outPortForwardingList   []*PortForwarding
	portForwardingList      []*PortForwarding

	addNodeLock sync.RWMutex
}

func (this_ *Worker) findNode(code string) (find *Info) {
	for _, one := range this_.nodeList {
		if one.Code == code {
			find = one
		}
	}
	return
}

func (this_ *Worker) findChildrenNode(code string) (find *Info) {
	for _, one := range this_.childrenNodeList {
		if one.Code == code {
			find = one
		}
	}
	return
}

func (this_ *Worker) addChildrenNode(childrenNode *Info) {
	var find = this_.findChildrenNode(childrenNode.Code)
	if find == nil {
		this_.childrenNodeClientCache[childrenNode.Code] = this_.getChildrenNodeClient(childrenNode)
		this_.childrenNodeList = append(this_.childrenNodeList, childrenNode)
	} else {
		copyNode(childrenNode, find)
	}
}

func (this_ *Worker) getChildrenNodeClient(childrenNode *Info) (client *Client) {
	client, ok := this_.childrenNodeClientCache[childrenNode.Code]
	if !ok {
		client = &Client{
			Ip:       childrenNode.Ip,
			Port:     childrenNode.Port,
			Node:     childrenNode,
			doMethod: this_.doMethod,
		}

		err := client.Start()
		if err != nil {
		}
		this_.childrenNodeClientCache[childrenNode.Code] = client
	}
	if client.isStopped() {
		err := client.Start()
		if err != nil {
			childrenNode.Status = StatusError
			childrenNode.StatusError = err.Error()
			return
		}
		bs, _ := json.Marshal(this_.nodeList)
		if len(bs) > 0 {
			_, _ = client.Call(AddNodeList, bs)
		}
	}
	return
}

func (this_ *Worker) AddNode(node *Info) {
	this_.addNodeLock.Lock()
	defer this_.addNodeLock.Unlock()

	for _, one := range this_.childrenNodeList {
		client := this_.getChildrenNodeClient(one)
		if !client.isStopped() {
			bs, _ := json.Marshal(node)
			if len(bs) > 0 {
				_, _ = client.Call(AddNode, bs)
			}
		}
	}

	var find *Info

	for _, one := range this_.nodeList {
		if one.Code == node.Code {
			find = one
		}
	}
	if find == nil {
		this_.nodeList = append(this_.nodeList, node)
	} else {
		copyNode(node, find)
	}

	this_.refreshNodeList()

	return
}

func (this_ *Worker) refreshNodeList() {
	for _, one := range this_.nodeList {
		var find = this_.findChildrenNode(one.Code)

		if find == nil {
			if one.ParentCode == this_.Node.Code {
				this_.addChildrenNode(one)
			}
		}
	}
	return
}

func (this_ *Worker) AddPortForwarding(portForwarding *PortForwarding) (err error) {
	if portForwarding.InNode == this_.Node.Code {

	}
	return

}

func (this_ *Worker) Start() (err error) {
	this_.childrenNodeClientCache = map[string]*Client{}
	this_.server = &Server{
		Ip:     this_.Node.Ip,
		Port:   this_.Node.Port,
		Worker: this_,
	}
	err = this_.server.Start()
	if err != nil {
		return
	}
	println(fmt.Sprintf(this_.Node.GetNodeStr() + " 服务启动成功"))
	go this_.childrenListen()
	return
}

func (this_ *Worker) childrenListen() {

	for {
		if this_.isStopped() {
			return
		}
		for _, one := range this_.childrenNodeList {
			this_.childrenCheck(one)
			switch one.Status {
			case StatusStarted:
				println(fmt.Sprintf(this_.Node.GetNodeStr() + " 子节点 " + one.GetNodeStr() + " 验证成功"))
				break
			case StatusStopped:
				println(fmt.Sprintf(this_.Node.GetNodeStr() + " 子节点 " + one.GetNodeStr() + " 未启动"))
				break
			case StatusError:
				println(fmt.Sprintf(this_.Node.GetNodeStr()+" 子节点 "+one.GetNodeStr()+" 验证异常 [%s]", one.StatusError))
				break
			}
		}
		if this_.isStopped() {
			return
		}
		time.Sleep(time.Second * 5)
	}

}

func (this_ *Worker) childrenCheck(childrenNode *Info) {

	client := this_.getChildrenNodeClient(childrenNode)
	if client.isStopped() {
		return
	}
	childrenNode.Status = StatusStarted
	childrenNode.StatusError = ""

	sayHello := "Say Hello"

	res, err := client.Call(OK, []byte(sayHello))
	if err != nil {
		childrenNode.Status = StatusError
		childrenNode.StatusError = err.Error()
		return
	}

	if string(res) != sayHello {
		childrenNode.Status = StatusError
		childrenNode.StatusError = "服务节点验证失败"
		return
	}
	childrenNode.Status = StatusStarted
	childrenNode.StatusError = ""

	return
}

var (
	OK                = "OK"
	AddNode           = "addNode"
	AddNodeList       = "addNodeList"
	AddPortForwarding = "addPortForwarding"
)

func (this_ *Worker) doMethod(msg *MethodMessage) (body []byte, err error) {
	if msg == nil {
		return
	}

	switch msg.Method {
	case OK:
		body = msg.Body
		return
	case AddNode:
		data := &Info{}
		err = json.Unmarshal(msg.Body, data)
		if err != nil {
			return
		}
		this_.AddNode(data)
		return
	case AddNodeList:
		var data []*Info
		err = json.Unmarshal(msg.Body, &data)
		if err != nil {
			return
		}
		for _, one := range data {
			this_.AddNode(one)
		}
		return
	case AddPortForwarding:
		data := &PortForwarding{}

		err = json.Unmarshal(msg.Body, data)
		if err != nil {
			return
		}
		err = this_.AddPortForwarding(data)
		return
	}

	return
}

func (this_ *Worker) isStopped() bool {

	return this_.isStop
}
