package terminal

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type Size struct {
	Cols int `json:"cols"`
	Rows int `json:"rows"`
}
type Service interface {
	Start(size *Size) (err error)
	Write(buf []byte) (n int, err error)
	Read(buf []byte) (n int, err error)
	ChangeSize(size *Size) (err error)
	Stop()
	IsWindows() (isWindows bool, err error)

	GetDiskStats() (res map[string]disk.IOCountersStat, err error)
	GetMemInfo() (res *mem.VirtualMemoryStat, err error)
	GetCpuInfo() (res []cpu.InfoStat, err error)
	GetCpuPercent() (res []float64, err error)
	GetCpuStats() (res []cpu.TimesStat, err error)
}
