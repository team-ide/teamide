package system

import (
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"go.uber.org/zap"
	"reflect"
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

type VirtualMemoryStat struct {
	Total          uint64  `json:"total,omitempty"`
	Available      uint64  `json:"available,omitempty"`
	Used           uint64  `json:"used,omitempty"`
	UsedPercent    float64 `json:"usedPercent,omitempty"`
	Free           uint64  `json:"free,omitempty"`
	Active         uint64  `json:"active,omitempty"`
	Inactive       uint64  `json:"inactive,omitempty"`
	Wired          uint64  `json:"wired,omitempty"`
	Laundry        uint64  `json:"laundry,omitempty"`
	Buffers        uint64  `json:"buffers,omitempty"`
	Cached         uint64  `json:"cached,omitempty"`
	WriteBack      uint64  `json:"writeBack,omitempty"`
	Dirty          uint64  `json:"dirty,omitempty"`
	WriteBackTmp   uint64  `json:"writeBackTmp,omitempty"`
	Shared         uint64  `json:"shared,omitempty"`
	Slab           uint64  `json:"slab,omitempty"`
	Sreclaimable   uint64  `json:"sreclaimable,omitempty"`
	Sunreclaim     uint64  `json:"sunreclaim,omitempty"`
	PageTables     uint64  `json:"pageTables,omitempty"`
	SwapCached     uint64  `json:"swapCached,omitempty"`
	CommitLimit    uint64  `json:"commitLimit,omitempty"`
	CommittedAS    uint64  `json:"committedAS,omitempty"`
	HighTotal      uint64  `json:"highTotal,omitempty"`
	HighFree       uint64  `json:"highFree,omitempty"`
	LowTotal       uint64  `json:"lowTotal,omitempty"`
	LowFree        uint64  `json:"lowFree,omitempty"`
	SwapTotal      uint64  `json:"swapTotal,omitempty"`
	SwapFree       uint64  `json:"swapFree,omitempty"`
	Mapped         uint64  `json:"mapped,omitempty"`
	VmallocTotal   uint64  `json:"vmallocTotal,omitempty"`
	VmallocUsed    uint64  `json:"vmallocUsed,omitempty"`
	VmallocChunk   uint64  `json:"vmallocChunk,omitempty"`
	HugePagesTotal uint64  `json:"hugePagesTotal,omitempty"`
	HugePagesFree  uint64  `json:"hugePagesFree,omitempty"`
	HugePagesRsvd  uint64  `json:"hugePagesRsvd,omitempty"`
	HugePagesSurp  uint64  `json:"hugePagesSurp,omitempty"`
	HugePageSize   uint64  `json:"hugePageSize,omitempty"`
}

type DiskUsageStat struct {
	Path              string  `json:"path,omitempty"`
	Fstype            string  `json:"fstype,omitempty"`
	Total             uint64  `json:"total,omitempty"`
	Free              uint64  `json:"free,omitempty"`
	Used              uint64  `json:"used,omitempty"`
	UsedPercent       float64 `json:"usedPercent,omitempty"`
	InodesTotal       uint64  `json:"inodesTotal,omitempty"`
	InodesUsed        uint64  `json:"inodesUsed,omitempty"`
	InodesFree        uint64  `json:"inodesFree,omitempty"`
	InodesUsedPercent float64 `json:"inodesUsedPercent,omitempty"`
}

type HostInfoStat struct {
	Hostname             string `json:"hostname,omitempty"`
	Uptime               uint64 `json:"uptime,omitempty"`
	BootTime             uint64 `json:"bootTime,omitempty"`
	Procs                uint64 `json:"procs,omitempty"`           // number of processes
	OS                   string `json:"os,omitempty"`              // ex: freebsd, linux
	Platform             string `json:"platform,omitempty"`        // ex: ubuntu, linuxmint
	PlatformFamily       string `json:"platformFamily,omitempty"`  // ex: debian, rhel
	PlatformVersion      string `json:"platformVersion,omitempty"` // version of the complete OS
	KernelVersion        string `json:"kernelVersion,omitempty"`   // version of the OS kernel (if available)
	KernelArch           string `json:"kernelArch,omitempty"`      // native cpu architecture queried at runtime, as returned by `uname -m` or empty string in case of error
	VirtualizationSystem string `json:"virtualizationSystem,omitempty"`
	VirtualizationRole   string `json:"virtualizationRole,omitempty"` // guest or host
	HostID               string `json:"hostId,omitempty"`             // ex: uuid
}

type NetIOCountersStat struct {
	Name        string `json:"name,omitempty"`        // interface name
	BytesSent   uint64 `json:"bytesSent,omitempty"`   // number of bytes sent
	BytesRecv   uint64 `json:"bytesRecv,omitempty"`   // number of bytes received
	PacketsSent uint64 `json:"packetsSent,omitempty"` // number of packets sent
	PacketsRecv uint64 `json:"packetsRecv,omitempty"` // number of packets received
	Errin       uint64 `json:"errin,omitempty"`       // total number of errors while receiving
	Errout      uint64 `json:"errout,omitempty"`      // total number of errors while sending
	Dropin      uint64 `json:"dropin,omitempty"`      // total number of incoming packets which were dropped
	Dropout     uint64 `json:"dropout,omitempty"`     // total number of outgoing packets which were dropped (always 0 on OSX and BSD)
	Fifoin      uint64 `json:"fifoin,omitempty"`      // total number of FIFO buffers errors while receiving
	Fifoout     uint64 `json:"fifoout,omitempty"`     // total number of FIFO buffers errors while sending
}

type CpuInfoStat struct {
	CPU        int32    `json:"cpu,omitempty"`
	VendorID   string   `json:"vendorId,omitempty"`
	Family     string   `json:"family,omitempty"`
	Model      string   `json:"model,omitempty"`
	Stepping   int32    `json:"stepping,omitempty"`
	PhysicalID string   `json:"physicalId,omitempty"`
	CoreID     string   `json:"coreId,omitempty"`
	Cores      int32    `json:"cores,omitempty"`
	ModelName  string   `json:"modelName,omitempty"`
	Mhz        float64  `json:"mhz,omitempty"`
	CacheSize  int32    `json:"cacheSize,omitempty"`
	Flags      []string `json:"flags,omitempty"`
	Microcode  string   `json:"microcode,omitempty"`
}
type Info struct {
	BootTime           int64                `json:"bootTime,omitempty"`
	VirtualMemoryStat  *VirtualMemoryStat   `json:"virtualMemoryStat,omitempty"`
	CpuInfoStats       []*CpuInfoStat       `json:"cpuInfoStats,omitempty"`
	CpuPercents        []float64            `json:"cpuPercents,omitempty"`
	DiskUsageStat      *DiskUsageStat       `json:"diskUsageStat,omitempty"`
	HostInfoStat       *HostInfoStat        `json:"hostInfoStat,omitempty"`
	NetIOCountersStats []*NetIOCountersStat `json:"netIOCountersStats,omitempty"`
	StartTime          int64                `json:"startTime,omitempty"`
	EndTime            int64                `json:"endTime,omitempty"`
}

type QueryRequest struct {
	Timestamp int64 `json:"timestamp,omitempty"`
	Size      int   `json:"size,omitempty"`
}

type QueryResponse struct {
	LastTimestamp int64   `json:"lastTimestamp,omitempty"`
	InfoList      []*Info `json:"infoList,omitempty"`
	Size          int     `json:"size,omitempty"`
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

func QueryInfo(request *QueryRequest) (response *QueryResponse) {
	response = &QueryResponse{}
	var startTimestamp int64
	var size int
	if request != nil {
		startTimestamp = request.Timestamp
		size = request.Size
	}
	if size <= 0 {
		size = 100
	}

	var list = infoList
	var appendSize = 0
	for _, one := range list {
		if appendSize >= size {
			break
		}
		if one.StartTime <= startTimestamp {
			continue
		}
		response.InfoList = append(response.InfoList, one)
		response.LastTimestamp = one.StartTime
		appendSize++
	}

	return
}

func Clean() {
	infoListLock.Lock()
	defer infoListLock.Unlock()

	infoList = []*Info{}
}

func GetInfo() (info *Info, err error) {
	info = &Info{
		StartTime: util.GetNowTime(),
	}
	defer func() {
		info.EndTime = util.GetNowTime()
	}()

	bootTime, err := host.BootTime()
	if err != nil {
		return
	}
	info.BootTime = util.GetTimeTime(time.Unix(int64(bootTime), 0))

	info.CpuPercents, err = cpu.Percent(time.Second, true)
	if err != nil {
		return
	}

	virtualMemoryStat, err := mem.VirtualMemory()
	if err != nil {
		return
	}
	info.VirtualMemoryStat = &VirtualMemoryStat{}
	err = SimpleCopyProperties(info.VirtualMemoryStat, virtualMemoryStat)
	if err != nil {
		return
	}

	cpus, err := cpu.Info()
	if err != nil {
		return
	}
	for _, one := range cpus {
		cInfo := &CpuInfoStat{}
		err = SimpleCopyProperties(cInfo, one)
		if err != nil {
			return
		}
		info.CpuInfoStats = append(info.CpuInfoStats, cInfo)
	}

	diskUsageStat, err := disk.Usage("/")
	if err != nil {
		return
	}
	info.DiskUsageStat = &DiskUsageStat{}
	err = SimpleCopyProperties(info.DiskUsageStat, diskUsageStat)
	if err != nil {
		return
	}

	hostInfoStat, err := host.Info()
	if err != nil {
		return
	}
	info.HostInfoStat = &HostInfoStat{}
	err = SimpleCopyProperties(info.HostInfoStat, hostInfoStat)
	if err != nil {
		return
	}

	netIOCountersStats, err := net.IOCounters(true)
	if err != nil {
		return
	}
	for _, one := range netIOCountersStats {
		nInfo := &NetIOCountersStat{}
		err = SimpleCopyProperties(nInfo, one)
		if err != nil {
			return
		}
		info.NetIOCountersStats = append(info.NetIOCountersStats, nInfo)
	}

	return
}

func SimpleCopyProperties(dst, src interface{}) (err error) {
	// 防止意外panic
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	// dst必须结构体指针类型
	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return errors.New("dst type should be a struct pointer")
	}

	// src必须为结构体或者结构体指针
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return errors.New("src type should be a struct or a struct pointer")
	}

	// 取具体内容
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	// 属性个数
	propertyNums := dstType.NumField()

	for i := 0; i < propertyNums; i++ {
		// 属性
		property := dstType.Field(i)
		// 待填充属性值
		propertyValue := srcValue.FieldByName(property.Name)

		// 无效，说明src没有这个属性 || 属性同名但类型不同
		if !propertyValue.IsValid() || property.Type != propertyValue.Type() {
			continue
		}

		if dstValue.Field(i).CanSet() {
			dstValue.Field(i).Set(propertyValue)
		}
	}

	return nil
}
