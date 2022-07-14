package node

import (
	"fmt"
	"net"
)

type Client struct {
	Ip       string
	Port     int
	Node     *Info
	doMethod func(msg *MethodMessage) (body []byte, err error)
	*Processor
}

func (this_ *Client) Start() (err error) {
	address := fmt.Sprintf("%s:%d", this_.Ip, this_.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("client dial err:", err.Error())
		return
	}
	this_.Processor = &Processor{
		conn:     conn,
		Node:     this_.Node,
		doMethod: this_.doMethod,
	}
	this_.Processor.listen()
	return
}

func (this_ *Client) Stop() {
	if this_.Processor == nil {
		return
	}
	this_.Processor.stopListen()
}

func (this_ *Client) isStopped() bool {
	return this_.Processor == nil || this_.needStop()
}
