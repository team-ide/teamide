package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Worker struct {
	Node             *Info
	ChildrenNodeList []*Info
	ParentNode       *Info
	Server           *Server
	isStop           bool
}

func (this_ *Worker) AddChildrenNode(childrenNode *Info) (err error) {
	if childrenNode.ParentCode == this_.Node.Code {
		var find *Info

		for _, one := range this_.ChildrenNodeList {
			if one.Token == childrenNode.Token {
				find = one
			}
		}
		if find == nil {
			this_.ChildrenNodeList = append(this_.ChildrenNodeList, childrenNode)
		}
	} else {
		var callNodeList []*Info
		for _, one := range this_.ChildrenNodeList {
			if one.Token == childrenNode.ParentCode {
				callNodeList = append(callNodeList, one)
			}
		}
		if len(callNodeList) == 0 {
			callNodeList = append(callNodeList, this_.ChildrenNodeList...)
		}
		var bs []byte
		bs, err = json.Marshal(childrenNode)
		if err != nil {
			return
		}
		var lastErr error
		var added bool
		for _, one := range callNodeList {

			client := &Client{
				Ip:       one.Ip,
				Port:     one.Port,
				Node:     one,
				doMethod: this_.doMethod,
			}
			connectErr := client.Start()
			if connectErr != nil {
				if lastErr == nil {
					lastErr = connectErr
				}
				continue
			}

			_, err = client.Call(OK, bs)
			if err != nil {
				client.Stop()
				return
			}
			client.Stop()
			added = true
		}
		if added {
			return
		}
		if lastErr != nil {
			err = lastErr
			return
		}
		err = errors.New(childrenNode.GetNodeStr() + " 无法寻址到父节点，请检查父节点是否存在并且在线")
		return
	}
	return
}

func (this_ *Worker) AddPortForwarding(portForwarding *PortForwarding) (err error) {
	if portForwarding.InNode == this_.Node.Code {

	}
	return

}

func (this_ *Worker) Start() (err error) {
	this_.Server = &Server{
		Ip:     this_.Node.Ip,
		Port:   this_.Node.Port,
		Worker: this_,
	}
	err = this_.Server.Start()
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
		for _, one := range this_.ChildrenNodeList {
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

	client := &Client{
		Ip:       childrenNode.Ip,
		Port:     childrenNode.Port,
		Node:     childrenNode,
		doMethod: this_.doMethod,
	}

	err := client.Start()
	if err != nil {
		childrenNode.Status = StatusError
		childrenNode.StatusError = err.Error()
		return
	}

	defer client.Stop()

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
	AddChildrenNode   = "addChildrenNode"
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
	case AddChildrenNode:
		data := &Info{}

		err = json.Unmarshal(msg.Body, data)
		if err != nil {
			return
		}
		err = this_.AddChildrenNode(data)
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
