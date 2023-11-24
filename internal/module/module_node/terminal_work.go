package module_node

import (
	"errors"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/node"
	"teamide/pkg/terminal"
	"time"
)

func NewTerminalService(nodeId string, nodeService *NodeService) (res *terminalService) {
	res = &terminalService{
		nodeId:      nodeId,
		nodeService: nodeService,
		bytesChan:   make(chan []byte),
	}
	return
}

type terminalService struct {
	nodeId      string
	key         string
	nodeLine    []string
	nodeService *NodeService
	bytesChan   chan []byte
}

func (this_ *terminalService) getServer() (server *node.Server, err error) {
	if this_.nodeService.GetContext() == nil {
		err = errors.New("node上下文未初始化")
		return
	}
	server = this_.nodeService.GetContext().GetServer()
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
	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("Stop err", zap.Any("err", e))
		}
	}()

	var err error
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	if this_.bytesChan != nil {
		close(this_.bytesChan)
		this_.bytesChan = nil
	}

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
		})
	return

}

func (this_ *terminalService) Write(buf []byte) (n int, err error) {
	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("Write err", zap.Any("err", e))
		}
	}()

	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	err = server.TerminalWrite(this_.nodeLine, this_.key, buf)
	return
}

func (this_ *terminalService) Read(buf []byte) (n int, err error) {
	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("Read err", zap.Any("err", e))
		}
	}()
	if this_.bytesChan == nil {
		err = errors.New("bytesChan is close")
		return
	}

	bs := <-this_.bytesChan
	n = len(bs)
	for i, b := range bs {
		buf[i] = b
	}
	return
}

func (this_ *terminalService) GetDiskStats() (res map[string]disk.IOCountersStat, err error) {
	return disk.IOCounters()
}
func (this_ *terminalService) GetMemInfo() (res *mem.VirtualMemoryStat, err error) {
	return mem.VirtualMemory()
}
func (this_ *terminalService) GetCpuInfo() (res []cpu.InfoStat, err error) {
	return cpu.Info()
}
func (this_ *terminalService) GetCpuPercent() (res []float64, err error) {
	return cpu.Percent(time.Second, true)
}
func (this_ *terminalService) GetCpuStats() (res []cpu.TimesStat, err error) {
	return cpu.Times(true)
}
