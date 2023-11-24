package terminal

import (
	"errors"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"sync"
	"time"
)

func NewLocalService() (res *localService) {
	res = &localService{}
	return
}

type terminalStart struct {
	Stop       func()
	Write_     func(p []byte) (n int, err error)
	Read_      func(b []byte) (int, error)
	ChangeSize func(size *Size) (err error)
}

func (this_ *terminalStart) Write(p []byte) (n int, err error) {
	return this_.Write_(p)
}

func (this_ *terminalStart) Read(p []byte) (n int, err error) {
	return this_.Read_(p)
}

type localService struct {
	terminalStart *terminalStart
	onClose       func()
	readeLock     sync.Mutex
	writeLock     sync.Mutex
}

func (this_ *localService) IsWindows() (isWindows bool, err error) {
	return IsWindows(), nil
}

func (this_ *localService) Stop() {
	if this_.terminalStart != nil {
		this_.terminalStart.Stop()
		this_.terminalStart = nil
	}
}

func (this_ *localService) ChangeSize(size *Size) (err error) {
	if this_.terminalStart == nil {
		return
	}
	if this_.terminalStart.ChangeSize != nil {
		err = this_.terminalStart.ChangeSize(size)
	}
	return
}

func (this_ *localService) Start(size *Size) (err error) {

	this_.terminalStart, err = start(size)

	if err != nil {
		util.Logger.Error("terminal local start error", zap.Error(err))
		return
	}

	util.Logger.Info("terminal local start success")

	return
}

func (this_ *localService) Write(buf []byte) (n int, err error) {

	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("Write err", zap.Any("err", e))
		}
	}()
	if this_.terminalStart == nil {
		err = errors.New("stdin is close")
		return
	}

	this_.writeLock.Lock()
	defer this_.writeLock.Unlock()

	//util.Logger.Info("local terminal write start")

	n = len(buf)
	err = util.Write(this_.terminalStart, buf, nil)
	return
}

func (this_ *localService) Read(buf []byte) (n int, err error) {

	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("Read err", zap.Any("err", e))
		}
	}()
	if this_.terminalStart == nil {
		err = errors.New("stdout is close")
		return
	}

	this_.readeLock.Lock()
	defer this_.readeLock.Unlock()

	//util.Logger.Info("local terminal read start")
	n, err = this_.terminalStart.Read(buf)
	return
}

func (this_ *localService) GetDiskStats() (res map[string]disk.IOCountersStat, err error) {
	return disk.IOCounters()
}
func (this_ *localService) GetMemInfo() (res *mem.VirtualMemoryStat, err error) {
	return mem.VirtualMemory()
}
func (this_ *localService) GetCpuInfo() (res []cpu.InfoStat, err error) {
	return cpu.Info()
}
func (this_ *localService) GetCpuPercent() (res []float64, err error) {
	return cpu.Percent(time.Second, true)
}
func (this_ *localService) GetCpuStats() (res []cpu.TimesStat, err error) {
	return cpu.Times(true)
}
