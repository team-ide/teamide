package module_node

import (
	"errors"
	"teamide/pkg/node"
	"teamide/pkg/terminal"
)

func NewTerminalService(nodeId string, nodeService *NodeService) (res *terminalService) {
	res = &terminalService{
		nodeId:         nodeId,
		nodeService:    nodeService,
		bytesChan:      make(chan []byte),
		errorBytesChan: make(chan []byte),
	}
	return
}

type terminalService struct {
	nodeId         string
	key            string
	nodeLine       []string
	nodeService    *NodeService
	bytesChan      chan []byte
	errorBytesChan chan []byte
}

func (this_ *terminalService) getServer() (server *node.Server, err error) {
	if this_.nodeService.GetContext() == nil {
		err = errors.New("node上下文未初始化")
		return
	}
	server = this_.nodeService.GetContext().GetServer()
	if server == nil {
		err = errors.New("node服务未初始化")
		return
	}
	nodeLine := this_.nodeService.GetContext().GetNodeLineTo(this_.nodeId)
	if len(nodeLine) == 0 {
		err = errors.New("无法连接到节点[" + this_.nodeId + "]")
		return
	}
	this_.nodeLine = nodeLine
	return
}

func (this_ *terminalService) IsWindows() (isWindows bool, err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	isWindows, err = server.TerminalIsWindows(this_.nodeLine)
	return
}

func (this_ *terminalService) Stop() {
	var err error
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	close(this_.bytesChan)
	close(this_.errorBytesChan)

	err = server.TerminalStop(this_.nodeLine, this_.key)

	return
}

func (this_ *terminalService) ChangeSize(size *terminal.Size) (err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	err = server.TerminalChangeSize(this_.nodeLine, this_.key, size)
	return
}

func (this_ *terminalService) Start(size *terminal.Size) (err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	this_.key, err = server.TerminalStart(this_.nodeLine, size,
		func(buf []byte) (err error) {
			this_.bytesChan <- buf
			return
		}, func(buf []byte) (err error) {
			this_.errorBytesChan <- buf
			return
		})
	return

}

func (this_ *terminalService) Write(buf []byte) (n int, err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	err = server.TerminalWrite(this_.nodeLine, this_.key, buf)
	return
}

func (this_ *terminalService) Read(buf []byte) (n int, err error) {
	bs := <-this_.bytesChan
	n = len(bs)
	for i, b := range bs {
		buf[i] = b
	}
	return
}

func (this_ *terminalService) ReadError(buf []byte) (n int, err error) {
	bs := <-this_.errorBytesChan
	n = len(bs)
	for i, b := range bs {
		buf[i] = b
	}
	return
}
