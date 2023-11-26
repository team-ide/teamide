package terminal

import (
	"teamide/pkg/system"
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

	SystemInfo() (res *system.Info, err error)
	SystemMonitorData() (res *system.MonitorData, err error)
}
