package system

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
	AnonHugePages  uint64  `json:"anonHugePages,omitempty"`
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

type DiskIOCountersStat struct {
	Name               string `json:"name,omitempty"` // interface name
	ReadCount          uint64 `json:"readCount,omitempty"`
	MergedReadCount    uint64 `json:"mergedReadCount,omitempty"`
	WriteCount         uint64 `json:"writeCount,omitempty"`
	MergedWriteCount   uint64 `json:"mergedWriteCount,omitempty"`
	ReadBytes          uint64 `json:"readBytes,omitempty"`
	WriteBytes         uint64 `json:"writeBytes,omitempty"`
	ReadTime           uint64 `json:"readTime,omitempty"`
	WriteTime          uint64 `json:"writeTime,omitempty"`
	IopsInProgress     uint64 `json:"iopsInProgress,omitempty"`
	IoTime             uint64 `json:"ioTime,omitempty"`
	WeightedIO         uint64 `json:"weightedIO,omitempty"`
	SerialNumber       string `json:"serialNumber,omitempty"`
	Label              string `json:"label,omitempty"`
	ReadBytesSpeed     uint64 `json:"readBytesSpeed,omitempty"`
	WriteBytesSpeed    uint64 `json:"writeBytesSpeed,omitempty"`
	ReadCountIncrease  uint64 `json:"readCountIncrease,omitempty"`
	WriteCountIncrease uint64 `json:"writeCountIncrease,omitempty"`
}

type NetIOCountersStat struct {
	Name        string `json:"name,omitempty"`        // interface name
	BytesSent   uint64 `json:"bytesSent,omitempty"`   // number of bytes sent
	BytesRecv   uint64 `json:"bytesRecv,omitempty"`   // number of bytes received
	PacketsSent uint64 `json:"packetsSent,omitempty"` // number of packets sent
	PacketsRecv uint64 `json:"packetsRecv,omitempty"` // number of packets received
	SpeedSent   uint64 `json:"speedSent,omitempty"`   // number of packets sent
	SpeedRecv   uint64 `json:"speedRecv,omitempty"`   // number of packets received
	Errin       uint64 `json:"errin,omitempty"`       // total number of errors while receiving
	Errout      uint64 `json:"errout,omitempty"`      // total number of errors while sending
	Dropin      uint64 `json:"dropin,omitempty"`      // total number of incoming packets which were dropped
	Dropout     uint64 `json:"dropout,omitempty"`     // total number of outgoing packets which were dropped (always 0 on OSX and BSD)
	Fifoin      uint64 `json:"fifoin,omitempty"`      // total number of FIFO buffers errors while receiving
	Fifoout     uint64 `json:"fifoout,omitempty"`     // total number of FIFO buffers errors while sending
}

type CpuInfoStat struct {
	CPU        int32   `json:"cpu,omitempty"`
	VendorID   string  `json:"vendorId,omitempty"`
	Family     string  `json:"family,omitempty"`
	Model      string  `json:"model,omitempty"`
	Stepping   int32   `json:"stepping,omitempty"`
	PhysicalID string  `json:"physicalId,omitempty"`
	CoreID     string  `json:"coreId,omitempty"`
	Cores      int32   `json:"cores,omitempty"`
	ModelName  string  `json:"modelName,omitempty"`
	Mhz        float64 `json:"mhz,omitempty"`
	CacheSize  int32   `json:"cacheSize,omitempty"`
	//Flags      []string `json:"flags,omitempty"`
	Microcode string `json:"microcode,omitempty"`
}

type MonitorData struct {
	VirtualMemoryStat   *VirtualMemoryStat    `json:"virtualMemoryStat,omitempty"`
	CpuPercents         []float64             `json:"cpuPercents,omitempty"`
	NetIOCountersStats  []*NetIOCountersStat  `json:"netIOCountersStats,omitempty"`
	DiskIOCountersStats []*DiskIOCountersStat `json:"diskIOCountersStats,omitempty"`
	StartTime           int64                 `json:"startTime,omitempty"`
	EndTime             int64                 `json:"endTime,omitempty"`
}

type Info struct {
	HostInfoStat *HostInfoStat      `json:"hostInfoStat,omitempty"`
	CpuInfoStats []*CpuInfoStat     `json:"cpuInfoStats,omitempty"`
	Disks        []*DiskUsageStat   `json:"disks,omitempty"`
	Memory       *VirtualMemoryStat `json:"memory,omitempty"`
}

type QueryRequest struct {
	Timestamp int64 `json:"timestamp,omitempty"`
	Size      int   `json:"size,omitempty"`
}

type QueryResponse struct {
	LastTimestamp   int64          `json:"lastTimestamp,omitempty"`
	MonitorDataList []*MonitorData `json:"monitorDataList,omitempty"`
	Size            int            `json:"size,omitempty"`
}
