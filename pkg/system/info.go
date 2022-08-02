package system

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"go.uber.org/zap"
	"sync"
	"teamide/pkg/util"
	"time"
)

var (
	infoList       []*Info
	WaitSecond     = 5
	CollectMaxSize = 3600
	collectLock    sync.Mutex
	infoListLock   sync.Mutex
	isStart        bool
)

type Info struct {
	BootTime          time.Time              `json:"bootTime,omitempty"`
	VirtualMemoryStat *mem.VirtualMemoryStat `json:"virtualMemoryStat,omitempty"`
	Cpus              []cpu.InfoStat         `json:"cpus,omitempty"`
	CpuPercents       []float64              `json:"cpuPercents,omitempty"`
	DiskUsageStat     *disk.UsageStat        `json:"diskUsageStat,omitempty"`
	HostInfoStat      *host.InfoStat         `json:"hostInfoStat,omitempty"`
	NetIOCountersStat []net.IOCountersStat   `json:"netIOCountersStat,omitempty"`
	StartTime         time.Time              `json:"startTime,omitempty"`
	EndTime           time.Time              `json:"endTime,omitempty"`
}

func StopCollect() {
	isStart = false
}

func StartCollect() {
	collectLock.Lock()
	defer collectLock.Unlock()
	if isStart {
		return
	}
	isStart = true
	go start()
}

func start() {
	if !isStart {
		return
	}
	defer func() {
		time.Sleep(time.Second * time.Duration(WaitSecond))
		start()
	}()

	infoListLock.Lock()
	defer infoListLock.Unlock()

	info, err := GetInfo()
	if err != nil {
		util.Logger.Error("system get info error", zap.Error(err))
		return
	}
	if len(infoList) >= CollectMaxSize {
		infoList = infoList[len(infoList)-CollectMaxSize+1:]
	}
	infoList = append(infoList, info)
}

func GetInfoList() (info []*Info) {
	return infoList
}

func Clean() {
	infoListLock.Lock()
	defer infoListLock.Unlock()

	infoList = []*Info{}
}

func GetInfo() (info *Info, err error) {
	info = &Info{}
	startTime := time.Now()
	defer func() {
		endTime := time.Now()
		info.StartTime = startTime
		info.EndTime = endTime
	}()

	bootTime, err := host.BootTime()
	if err != nil {
		return
	}
	info.BootTime = time.Unix(int64(bootTime), 0)

	info.VirtualMemoryStat, err = mem.VirtualMemory()
	if err != nil {
		return
	}

	info.Cpus, err = cpu.Info()
	if err != nil {
		return
	}

	info.CpuPercents, err = cpu.Percent(time.Second, true)
	if err != nil {
		return
	}

	info.DiskUsageStat, err = disk.Usage("/")
	if err != nil {
		return
	}

	info.HostInfoStat, err = host.Info()
	if err != nil {
		return
	}

	info.NetIOCountersStat, err = net.IOCounters(true)
	if err != nil {
		return
	}

	return
}
