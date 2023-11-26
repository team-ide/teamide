package system

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"sync"
	"teamide/pkg/task"
	"time"
)

var (
	monitorDataList     []*MonitorData
	CollectMaxSize      = 3600
	collectLock         = &sync.Mutex{}
	monitorDataListLock = &sync.Mutex{}
	monitorDataTask     *task.CronTask
)

//func StopCollectMonitorData() {
//	if monitorDataTask != nil {
//		monitorDataTask.Stop()
//	}
//	monitorDataTask = nil
//}

func StartCollectMonitorData() {
	collectLock.Lock()
	defer collectLock.Unlock()
	if monitorDataTask != nil {
		return
	}
	monitorDataTask = &task.CronTask{
		Spec: "0/10 * * * * *",
		Task: &task.Task{
			Key: monitorDataTaskKey,
			Do: func() {
				monitorData, err := GetMonitorData()
				if err != nil {
					util.Logger.Error("system get info error", zap.Error(err))
					return
				}

				monitorDataListLock.Lock()
				defer monitorDataListLock.Unlock()

				if len(monitorDataList) >= CollectMaxSize {
					monitorDataList = monitorDataList[len(monitorDataList)-CollectMaxSize+1:]
				}
				monitorDataList = append(monitorDataList, monitorData)
			},
		},
	}

	_ = task.AddCronTask(monitorDataTask)
}

const (
	monitorDataTaskKey = "system-monitor-data-task-key"
)

func QueryMonitorData(request *QueryRequest) (response *QueryResponse) {
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
	//util.Logger.Info("QueryMonitorData start", zap.Any("request", request))
	var list = monitorDataList
	var appendSize = 0
	for _, one := range list {
		if appendSize >= size {
			break
		}
		if one.StartTime <= startTimestamp {
			continue
		}
		response.MonitorDataList = append(response.MonitorDataList, one)
		response.LastTimestamp = one.StartTime
		appendSize++
	}

	//util.Logger.Info("QueryMonitorData end", zap.Any("response", response))
	return
}

func CleanMonitorData() {
	monitorDataListLock.Lock()
	defer monitorDataListLock.Unlock()

	monitorDataList = []*MonitorData{}
}

func GetInfo() (info *Info) {
	info = &Info{}

	hostInfoStat, err := host.Info()
	if err != nil {
		return
	}
	info.HostInfoStat = &HostInfoStat{
		Hostname:             hostInfoStat.Hostname,
		Uptime:               hostInfoStat.Uptime,
		BootTime:             hostInfoStat.BootTime,
		Procs:                hostInfoStat.Procs,
		OS:                   hostInfoStat.OS,
		Platform:             hostInfoStat.Platform,
		PlatformFamily:       hostInfoStat.PlatformFamily,
		PlatformVersion:      hostInfoStat.PlatformVersion,
		KernelVersion:        hostInfoStat.KernelVersion,
		KernelArch:           hostInfoStat.KernelArch,
		VirtualizationSystem: hostInfoStat.VirtualizationSystem,
		VirtualizationRole:   hostInfoStat.VirtualizationRole,
		HostID:               hostInfoStat.HostID,
	}

	cpus, _ := cpu.Info()
	if cpus != nil {
		for _, one := range cpus {
			cInfo := &CpuInfoStat{
				CPU:        one.CPU,
				VendorID:   one.VendorID,
				Family:     one.Family,
				Model:      one.Model,
				Stepping:   one.Stepping,
				PhysicalID: one.PhysicalID,
				CoreID:     one.CoreID,
				Cores:      one.Cores,
				ModelName:  one.ModelName,
				Mhz:        one.Mhz,
				CacheSize:  one.CacheSize,
				//Flags:      one.Flags,
				Microcode: one.Microcode,
			}
			info.CpuInfoStats = append(info.CpuInfoStats, cInfo)
		}
	}

	virtualMemoryStat, _ := mem.VirtualMemory()
	if virtualMemoryStat != nil {
		info.Memory = &VirtualMemoryStat{
			Total:          virtualMemoryStat.Total,
			Available:      virtualMemoryStat.Available,
			Used:           virtualMemoryStat.Used,
			UsedPercent:    virtualMemoryStat.UsedPercent,
			Free:           virtualMemoryStat.Free,
			Active:         virtualMemoryStat.Active,
			Inactive:       virtualMemoryStat.Inactive,
			Wired:          virtualMemoryStat.Wired,
			Laundry:        virtualMemoryStat.Laundry,
			Buffers:        virtualMemoryStat.Buffers,
			Cached:         virtualMemoryStat.Cached,
			WriteBack:      virtualMemoryStat.WriteBack,
			Dirty:          virtualMemoryStat.Dirty,
			WriteBackTmp:   virtualMemoryStat.WriteBackTmp,
			Shared:         virtualMemoryStat.Shared,
			Slab:           virtualMemoryStat.Slab,
			Sreclaimable:   virtualMemoryStat.Sreclaimable,
			Sunreclaim:     virtualMemoryStat.Sunreclaim,
			PageTables:     virtualMemoryStat.PageTables,
			SwapCached:     virtualMemoryStat.SwapCached,
			CommitLimit:    virtualMemoryStat.CommitLimit,
			CommittedAS:    virtualMemoryStat.CommittedAS,
			HighTotal:      virtualMemoryStat.HighTotal,
			HighFree:       virtualMemoryStat.HighFree,
			LowTotal:       virtualMemoryStat.LowTotal,
			LowFree:        virtualMemoryStat.LowFree,
			SwapTotal:      virtualMemoryStat.SwapTotal,
			SwapFree:       virtualMemoryStat.SwapFree,
			Mapped:         virtualMemoryStat.Mapped,
			VmallocTotal:   virtualMemoryStat.VmallocTotal,
			VmallocUsed:    virtualMemoryStat.VmallocUsed,
			VmallocChunk:   virtualMemoryStat.VmallocChunk,
			HugePagesTotal: virtualMemoryStat.HugePagesTotal,
			HugePagesFree:  virtualMemoryStat.HugePagesFree,
			HugePagesRsvd:  virtualMemoryStat.HugePagesRsvd,
			HugePagesSurp:  virtualMemoryStat.HugePagesSurp,
			HugePageSize:   virtualMemoryStat.HugePageSize,
			AnonHugePages:  virtualMemoryStat.AnonHugePages,
		}
	}

	ps, _ := disk.Partitions(true)
	if ps != nil {
		for _, p := range ps {
			diskUsageStat, _ := disk.Usage(p.Device)
			if diskUsageStat == nil {
				continue
			}
			info.Disks = append(info.Disks, &DiskUsageStat{
				Path:              diskUsageStat.Path,
				Fstype:            diskUsageStat.Fstype,
				Total:             diskUsageStat.Total,
				Free:              diskUsageStat.Free,
				Used:              diskUsageStat.Used,
				UsedPercent:       diskUsageStat.UsedPercent,
				InodesTotal:       diskUsageStat.InodesTotal,
				InodesUsed:        diskUsageStat.InodesUsed,
				InodesFree:        diskUsageStat.InodesFree,
				InodesUsedPercent: diskUsageStat.InodesUsedPercent,
			})
		}
	}

	return
}

func GetCacheOrNew() (monitorData *MonitorData, err error) {
	monitorDataListLock.Lock()
	size := len(monitorDataList)
	if size > 0 {
		monitorData = monitorDataList[size-1]
		monitorDataListLock.Unlock()
		return
	}
	monitorDataListLock.Unlock()
	return GetMonitorData()
}

func GetMonitorData() (monitorData *MonitorData, err error) {
	monitorData = &MonitorData{
		StartTime: util.GetNowMilli(),
	}
	defer func() {
		monitorData.EndTime = util.GetNowMilli()
	}()

	monitorData.CpuPercents, err = cpu.Percent(0, true)
	if err != nil {
		return
	}

	virtualMemoryStat, err := mem.VirtualMemory()
	if err != nil {
		return
	}
	monitorData.VirtualMemoryStat = &VirtualMemoryStat{
		Total:          virtualMemoryStat.Total,
		Available:      virtualMemoryStat.Available,
		Used:           virtualMemoryStat.Used,
		UsedPercent:    virtualMemoryStat.UsedPercent,
		Free:           virtualMemoryStat.Free,
		Active:         virtualMemoryStat.Active,
		Inactive:       virtualMemoryStat.Inactive,
		Wired:          virtualMemoryStat.Wired,
		Laundry:        virtualMemoryStat.Laundry,
		Buffers:        virtualMemoryStat.Buffers,
		Cached:         virtualMemoryStat.Cached,
		WriteBack:      virtualMemoryStat.WriteBack,
		Dirty:          virtualMemoryStat.Dirty,
		WriteBackTmp:   virtualMemoryStat.WriteBackTmp,
		Shared:         virtualMemoryStat.Shared,
		Slab:           virtualMemoryStat.Slab,
		Sreclaimable:   virtualMemoryStat.Sreclaimable,
		Sunreclaim:     virtualMemoryStat.Sunreclaim,
		PageTables:     virtualMemoryStat.PageTables,
		SwapCached:     virtualMemoryStat.SwapCached,
		CommitLimit:    virtualMemoryStat.CommitLimit,
		CommittedAS:    virtualMemoryStat.CommittedAS,
		HighTotal:      virtualMemoryStat.HighTotal,
		HighFree:       virtualMemoryStat.HighFree,
		LowTotal:       virtualMemoryStat.LowTotal,
		LowFree:        virtualMemoryStat.LowFree,
		SwapTotal:      virtualMemoryStat.SwapTotal,
		SwapFree:       virtualMemoryStat.SwapFree,
		Mapped:         virtualMemoryStat.Mapped,
		VmallocTotal:   virtualMemoryStat.VmallocTotal,
		VmallocUsed:    virtualMemoryStat.VmallocUsed,
		VmallocChunk:   virtualMemoryStat.VmallocChunk,
		HugePagesTotal: virtualMemoryStat.HugePagesTotal,
		HugePagesFree:  virtualMemoryStat.HugePagesFree,
		HugePagesRsvd:  virtualMemoryStat.HugePagesRsvd,
		HugePagesSurp:  virtualMemoryStat.HugePagesSurp,
		HugePageSize:   virtualMemoryStat.HugePageSize,
		AnonHugePages:  virtualMemoryStat.AnonHugePages,
	}

	diskStat, _ := disk.IOCounters("/")
	diskStatCache := lastDiskIOCountersStatCache
	newSce := uint64(time.Now().Unix())
	if diskStat != nil {
		for _, one := range diskStat {
			nInfo := &DiskIOCountersStat{
				Name:             one.Name,
				ReadCount:        one.ReadCount,
				MergedReadCount:  one.MergedReadCount,
				WriteCount:       one.WriteCount,
				MergedWriteCount: one.MergedWriteCount,
				ReadBytes:        one.ReadBytes,
				WriteBytes:       one.WriteBytes,
				ReadTime:         one.ReadTime,
				WriteTime:        one.WriteTime,
				IopsInProgress:   one.IopsInProgress,
				IoTime:           one.IoTime,
				WeightedIO:       one.WeightedIO,
				SerialNumber:     one.SerialNumber,
				Label:            one.Label,
			}
			find := diskStatCache[nInfo.Name]
			if find.ReadBytes > 0 && nInfo.ReadBytes > find.ReadBytes {
				nInfo.ReadBytesSpeed = (nInfo.ReadBytes - find.ReadBytes) / (newSce - lastDiskIOCounters)
			}
			if find.WriteBytes > 0 && nInfo.WriteBytes > find.WriteBytes {
				nInfo.WriteBytesSpeed = (nInfo.WriteBytes - find.WriteBytes) / (newSce - lastDiskIOCounters)
			}
			if find.ReadCount > 0 && nInfo.ReadCount > find.ReadCount {
				nInfo.ReadCountIncrease = nInfo.ReadCount - find.ReadCount
			}
			if find.WriteCount > 0 && nInfo.WriteCount > find.WriteCount {
				nInfo.WriteCountIncrease = nInfo.WriteCount - find.WriteCount
			}
			monitorData.DiskIOCountersStats = append(monitorData.DiskIOCountersStats, nInfo)
		}
	}
	diskStatCache = make(map[string]disk.IOCountersStat)
	if diskStat != nil {
		for _, one := range diskStat {
			diskStatCache[one.Name] = one
		}
	}
	lastDiskIOCounters = newSce
	lastDiskIOCountersStatCache = diskStatCache

	newSce = uint64(time.Now().Unix())
	netIOCountersStats, _ := net.IOCounters(true)
	netStatCache := lastNetIOCountersStatCache
	for _, one := range netIOCountersStats {
		nInfo := &NetIOCountersStat{
			Name:        one.Name,
			BytesSent:   one.BytesSent,
			BytesRecv:   one.BytesRecv,
			PacketsSent: one.PacketsSent,
			PacketsRecv: one.PacketsRecv,
			Errin:       one.Errin,
			Errout:      one.Errout,
			Dropin:      one.Dropin,
			Dropout:     one.Dropout,
			Fifoin:      one.Fifoin,
			Fifoout:     one.Fifoout,
		}
		find := netStatCache[one.Name]
		if find.BytesSent > 0 && nInfo.BytesSent > find.BytesSent {
			nInfo.SpeedSent = (nInfo.BytesSent - find.BytesSent) / (newSce - lastNetIOCounters)
		}
		if find.BytesRecv > 0 && nInfo.BytesRecv > find.BytesRecv {
			nInfo.SpeedRecv = (nInfo.BytesRecv - find.BytesRecv) / (newSce - lastNetIOCounters)
		}
		monitorData.NetIOCountersStats = append(monitorData.NetIOCountersStats, nInfo)
	}
	netStatCache = make(map[string]net.IOCountersStat)
	for _, one := range netIOCountersStats {
		netStatCache[one.Name] = one
	}
	lastNetIOCounters = newSce
	lastNetIOCountersStatCache = netStatCache
	return
}

var lastDiskIOCounters uint64
var lastDiskIOCountersStatCache = make(map[string]disk.IOCountersStat)

var lastNetIOCounters uint64
var lastNetIOCountersStatCache = make(map[string]net.IOCountersStat)
