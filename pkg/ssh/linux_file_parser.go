package ssh

import (
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"math"
	"os"
	"strconv"
	"strings"
)

func ParseProcCpuInfo(cpuInfoText string, readLines func(filepath string) []string) ([]cpu.InfoStat, error) {
	lines := strings.Split(cpuInfoText, "\n")

	var ret []cpu.InfoStat
	var processorName string

	c := cpu.InfoStat{CPU: -1, Cores: 1}
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])

		switch key {
		case "Processor":
			processorName = value
		case "processor", "cpu number":
			if c.CPU >= 0 {
				finishCPUInfo(&c, readLines)
				ret = append(ret, c)
			}
			c = cpu.InfoStat{Cores: 1, ModelName: processorName}
			t, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return ret, err
			}
			c.CPU = int32(t)
		case "vendorId", "vendor_id":
			c.VendorID = value
			if strings.Contains(value, "S390") {
				processorName = "S390"
			}
		case "CPU implementer":
			if v, err := strconv.ParseUint(value, 0, 8); err == nil {
				switch v {
				case 0x41:
					c.VendorID = "ARM"
				case 0x42:
					c.VendorID = "Broadcom"
				case 0x43:
					c.VendorID = "Cavium"
				case 0x44:
					c.VendorID = "DEC"
				case 0x46:
					c.VendorID = "Fujitsu"
				case 0x48:
					c.VendorID = "HiSilicon"
				case 0x49:
					c.VendorID = "Infineon"
				case 0x4d:
					c.VendorID = "Motorola/Freescale"
				case 0x4e:
					c.VendorID = "NVIDIA"
				case 0x50:
					c.VendorID = "APM"
				case 0x51:
					c.VendorID = "Qualcomm"
				case 0x56:
					c.VendorID = "Marvell"
				case 0x61:
					c.VendorID = "Apple"
				case 0x69:
					c.VendorID = "Intel"
				case 0xc0:
					c.VendorID = "Ampere"
				}
			}
		case "cpu family":
			c.Family = value
		case "model", "CPU part":
			c.Model = value
			// if CPU is arm based, model name is found via model number. refer to: arch/arm64/kernel/cpuinfo.c
			if c.VendorID == "ARM" {
				if v, err := strconv.ParseUint(c.Model, 0, 16); err == nil {
					modelName, exist := armModelToModelName[v]
					if exist {
						c.ModelName = modelName
					} else {
						c.ModelName = "Undefined"
					}
				}
			}
		case "Model Name", "model name", "cpu":
			c.ModelName = value
			if strings.Contains(value, "POWER") {
				c.Model = strings.Split(value, " ")[0]
				c.Family = "POWER"
				c.VendorID = "IBM"
			}
		case "stepping", "revision", "CPU revision":
			val := value

			if key == "revision" {
				val = strings.Split(value, ".")[0]
			}

			t, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return ret, err
			}
			c.Stepping = int32(t)
		case "cpu MHz", "clock", "cpu MHz dynamic":
			// treat this as the fallback value, thus we ignore error
			if t, err := strconv.ParseFloat(strings.Replace(value, "MHz", "", 1), 64); err == nil {
				c.Mhz = t
			}
		case "cache size":
			t, err := strconv.ParseInt(strings.Replace(value, " KB", "", 1), 10, 64)
			if err != nil {
				return ret, err
			}
			c.CacheSize = int32(t)
		case "physical id":
			c.PhysicalID = value
		case "core id":
			c.CoreID = value
		case "flags", "Features":
			c.Flags = strings.FieldsFunc(value, func(r rune) bool {
				return r == ',' || r == ' '
			})
		case "microcode":
			c.Microcode = value
		}
	}
	if c.CPU >= 0 {
		finishCPUInfo(&c, readLines)
		ret = append(ret, c)
	}
	return ret, nil
}

var armModelToModelName = map[uint64]string{
	0x810: "ARM810",
	0x920: "ARM920",
	0x922: "ARM922",
	0x926: "ARM926",
	0x940: "ARM940",
	0x946: "ARM946",
	0x966: "ARM966",
	0xa20: "ARM1020",
	0xa22: "ARM1022",
	0xa26: "ARM1026",
	0xb02: "ARM11 MPCore",
	0xb36: "ARM1136",
	0xb56: "ARM1156",
	0xb76: "ARM1176",
	0xc05: "Cortex-A5",
	0xc07: "Cortex-A7",
	0xc08: "Cortex-A8",
	0xc09: "Cortex-A9",
	0xc0d: "Cortex-A17",
	0xc0f: "Cortex-A15",
	0xc0e: "Cortex-A17",
	0xc14: "Cortex-R4",
	0xc15: "Cortex-R5",
	0xc17: "Cortex-R7",
	0xc18: "Cortex-R8",
	0xc20: "Cortex-M0",
	0xc21: "Cortex-M1",
	0xc23: "Cortex-M3",
	0xc24: "Cortex-M4",
	0xc27: "Cortex-M7",
	0xc60: "Cortex-M0+",
	0xd01: "Cortex-A32",
	0xd02: "Cortex-A34",
	0xd03: "Cortex-A53",
	0xd04: "Cortex-A35",
	0xd05: "Cortex-A55",
	0xd06: "Cortex-A65",
	0xd07: "Cortex-A57",
	0xd08: "Cortex-A72",
	0xd09: "Cortex-A73",
	0xd0a: "Cortex-A75",
	0xd0b: "Cortex-A76",
	0xd0c: "Neoverse-N1",
	0xd0d: "Cortex-A77",
	0xd0e: "Cortex-A76AE",
	0xd13: "Cortex-R52",
	0xd20: "Cortex-M23",
	0xd21: "Cortex-M33",
	0xd40: "Neoverse-V1",
	0xd41: "Cortex-A78",
	0xd42: "Cortex-A78AE",
	0xd43: "Cortex-A65AE",
	0xd44: "Cortex-X1",
	0xd46: "Cortex-A510",
	0xd47: "Cortex-A710",
	0xd48: "Cortex-X2",
	0xd49: "Neoverse-N2",
	0xd4a: "Neoverse-E1",
	0xd4b: "Cortex-A78C",
	0xd4c: "Cortex-X1C",
	0xd4d: "Cortex-A715",
	0xd4e: "Cortex-X3",
}

func finishCPUInfo(c *cpu.InfoStat, readLines func(filepath string) []string) {
	var lines []string
	var err error
	var value float64
	var basePath = fmt.Sprintf("/sys/devices/system/cpu/cpu%d/", c.CPU)

	if len(c.CoreID) == 0 {
		lines = readLines(basePath + "topology/core_id")
		if len(lines) > 0 {
			c.CoreID = lines[0]
		}
	}

	// override the value of c.Mhz with cpufreq/cpuinfo_max_freq regardless
	// of the value from /proc/cpuinfo because we want to report the maximum
	// clock-speed of the CPU for c.Mhz, matching the behaviour of Windows
	lines = readLines(basePath + "cpufreq/cpuinfo_max_freq")
	// if we encounter errors below such as there are no cpuinfo_max_freq file,
	// we just ignore. so let Mhz is 0.
	if err != nil || len(lines) == 0 {
		return
	}
	value, err = strconv.ParseFloat(lines[0], 64)
	if err != nil {
		return
	}
	c.Mhz = value / 1000.0 // value is in kHz
	if c.Mhz > 9999 {
		c.Mhz = c.Mhz / 1000.0 // value in Hz
	}
}
func ParseProcStat(statText string) (res []cpu.TimesStat) {
	lines := strings.Split(statText, "\n")
	res = make([]cpu.TimesStat, 0, len(lines))

	if len(lines) > 1 {
		for _, line := range lines[1:] {
			if !strings.HasPrefix(line, "cpu") {
				break
			}
			ct, err := parseStatLine(line)
			if err != nil {
				continue
			}
			res = append(res, *ct)

		}
	}
	return
}

var ClocksPerSec = float64(100)

func parseStatLine(statText string) (*cpu.TimesStat, error) {
	fields := strings.Fields(statText)

	if len(fields) < 8 {
		return nil, errors.New("stat does not contain cpu info")
	}

	if !strings.HasPrefix(fields[0], "cpu") {
		return nil, errors.New("not contain cpu")
	}

	cpu_ := fields[0]
	if cpu_ == "cpu" {
		cpu_ = "cpu-total"
	}
	user, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, err
	}
	nice, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, err
	}
	system, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, err
	}
	idle, err := strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return nil, err
	}
	iowait, err := strconv.ParseFloat(fields[5], 64)
	if err != nil {
		return nil, err
	}
	irq, err := strconv.ParseFloat(fields[6], 64)
	if err != nil {
		return nil, err
	}
	softirq, err := strconv.ParseFloat(fields[7], 64)
	if err != nil {
		return nil, err
	}

	ct := &cpu.TimesStat{
		CPU:     cpu_,
		User:    user / ClocksPerSec,
		Nice:    nice / ClocksPerSec,
		System:  system / ClocksPerSec,
		Idle:    idle / ClocksPerSec,
		Iowait:  iowait / ClocksPerSec,
		Irq:     irq / ClocksPerSec,
		Softirq: softirq / ClocksPerSec,
	}
	if len(fields) > 8 { // Linux >= 2.6.11
		steal, err := strconv.ParseFloat(fields[8], 64)
		if err != nil {
			return nil, err
		}
		ct.Steal = steal / ClocksPerSec
	}
	if len(fields) > 9 { // Linux >= 2.6.24
		guest, err := strconv.ParseFloat(fields[9], 64)
		if err != nil {
			return nil, err
		}
		ct.Guest = guest / ClocksPerSec
	}
	if len(fields) > 10 { // Linux >= 3.2.0
		guestNice, err := strconv.ParseFloat(fields[10], 64)
		if err != nil {
			return nil, err
		}
		ct.GuestNice = guestNice / ClocksPerSec
	}

	return ct, nil
}

type VirtualMemoryExStat struct {
	ActiveFile   uint64 `json:"activefile"`
	InactiveFile uint64 `json:"inactivefile"`
	ActiveAnon   uint64 `json:"activeanon"`
	InactiveAnon uint64 `json:"inactiveanon"`
	Unevictable  uint64 `json:"unevictable"`
}

func ParseProcMemInfo(memInfoText string, zoneInfoText string) (*mem.VirtualMemoryStat, *VirtualMemoryExStat, error) {

	lines := strings.Split(memInfoText, "\n")

	// flag if MemAvailable is in /proc/meminfo (kernel 3.14+)
	memavail := false
	activeFile := false   // "Active(file)" not available: 2.6.28 / Dec 2008
	inactiveFile := false // "Inactive(file)" not available: 2.6.28 / Dec 2008
	sReclaimable := false // "Sreclaimable:" not available: 2.6.19 / Nov 2006

	ret := &mem.VirtualMemoryStat{}
	retEx := &VirtualMemoryExStat{}

	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		value = strings.Replace(value, " kB", "", -1)

		switch key {
		case "MemTotal":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.Total = t * 1024
		case "MemFree":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.Free = t * 1024
		case "MemAvailable":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			memavail = true
			ret.Available = t * 1024
		case "Buffers":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.Buffers = t * 1024
		case "Cached":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.Cached = t * 1024
		case "Active":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.Active = t * 1024
		case "Inactive":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.Inactive = t * 1024
		case "Active(anon)":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			retEx.ActiveAnon = t * 1024
		case "Inactive(anon)":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			retEx.InactiveAnon = t * 1024
		case "Active(file)":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			activeFile = true
			retEx.ActiveFile = t * 1024
		case "Inactive(file)":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			inactiveFile = true
			retEx.InactiveFile = t * 1024
		case "Unevictable":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			retEx.Unevictable = t * 1024
		case "Writeback":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.WriteBack = t * 1024
		case "WritebackTmp":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.WriteBackTmp = t * 1024
		case "Dirty":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.Dirty = t * 1024
		case "Shmem":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.Shared = t * 1024
		case "Slab":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.Slab = t * 1024
		case "SReclaimable":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			sReclaimable = true
			ret.Sreclaimable = t * 1024
		case "SUnreclaim":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.Sunreclaim = t * 1024
		case "PageTables":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.PageTables = t * 1024
		case "SwapCached":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.SwapCached = t * 1024
		case "CommitLimit":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.CommitLimit = t * 1024
		case "Committed_AS":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.CommittedAS = t * 1024
		case "HighTotal":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.HighTotal = t * 1024
		case "HighFree":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.HighFree = t * 1024
		case "LowTotal":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.LowTotal = t * 1024
		case "LowFree":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.LowFree = t * 1024
		case "SwapTotal":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.SwapTotal = t * 1024
		case "SwapFree":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.SwapFree = t * 1024
		case "Mapped":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.Mapped = t * 1024
		case "VmallocTotal":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.VmallocTotal = t * 1024
		case "VmallocUsed":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.VmallocUsed = t * 1024
		case "VmallocChunk":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.VmallocChunk = t * 1024
		case "HugePages_Total":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.HugePagesTotal = t
		case "HugePages_Free":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.HugePagesFree = t
		case "HugePages_Rsvd":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.HugePagesRsvd = t
		case "HugePages_Surp":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.HugePagesSurp = t
		case "Hugepagesize":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.HugePageSize = t * 1024
		case "AnonHugePages":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, retEx, err
			}
			ret.AnonHugePages = t * 1024
		}
	}

	ret.Cached += ret.Sreclaimable

	if !memavail {
		if activeFile && inactiveFile && sReclaimable {
			ret.Available = calculateAvailVMem(zoneInfoText, ret, retEx)
		} else {
			ret.Available = ret.Cached + ret.Free
		}
	}

	ret.Used = ret.Total - ret.Free - ret.Buffers - ret.Cached
	ret.UsedPercent = float64(ret.Used) / float64(ret.Total) * 100.0

	return ret, retEx, nil
}
func calculateAvailVMem(zoneInfoText string, ret *mem.VirtualMemoryStat, retEx *VirtualMemoryExStat) uint64 {
	var watermarkLow uint64

	lines := strings.Split(zoneInfoText, "\n")
	if len(lines) == 0 {
		return ret.Free + ret.Cached // fallback under kernel 2.6.13
	}

	pagesize := uint64(os.Getpagesize())
	watermarkLow = 0

	for _, line := range lines {
		fields := strings.Fields(line)

		if strings.HasPrefix(fields[0], "low") {
			lowValue, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				lowValue = 0
			}
			watermarkLow += lowValue
		}
	}

	watermarkLow *= pagesize

	availMemory := ret.Free - watermarkLow
	pageCache := retEx.ActiveFile + retEx.InactiveFile
	pageCache -= uint64(math.Min(float64(pageCache/2), float64(watermarkLow)))
	availMemory += pageCache
	availMemory += ret.Sreclaimable - uint64(math.Min(float64(ret.Sreclaimable/2.0), float64(watermarkLow)))

	if availMemory < 0 {
		availMemory = 0
	}

	return availMemory
}

const sectorSize uint64 = 512

func ParseProcDiskStats(diskStatsText string) (map[string]disk.IOCountersStat, error) {

	lines := strings.Split(diskStatsText, "\n")
	ret := make(map[string]disk.IOCountersStat)
	empty := disk.IOCountersStat{}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 14 {
			// malformed line in /proc/diskstats, avoid panic by ignoring.
			continue
		}
		name := fields[2]

		reads, err := strconv.ParseUint((fields[3]), 10, 64)
		if err != nil {
			return ret, err
		}
		mergedReads, err := strconv.ParseUint((fields[4]), 10, 64)
		if err != nil {
			return ret, err
		}
		rbytes, err := strconv.ParseUint((fields[5]), 10, 64)
		if err != nil {
			return ret, err
		}
		rtime, err := strconv.ParseUint((fields[6]), 10, 64)
		if err != nil {
			return ret, err
		}
		writes, err := strconv.ParseUint((fields[7]), 10, 64)
		if err != nil {
			return ret, err
		}
		mergedWrites, err := strconv.ParseUint((fields[8]), 10, 64)
		if err != nil {
			return ret, err
		}
		wbytes, err := strconv.ParseUint((fields[9]), 10, 64)
		if err != nil {
			return ret, err
		}
		wtime, err := strconv.ParseUint((fields[10]), 10, 64)
		if err != nil {
			return ret, err
		}
		iopsInProgress, err := strconv.ParseUint((fields[11]), 10, 64)
		if err != nil {
			return ret, err
		}
		iotime, err := strconv.ParseUint((fields[12]), 10, 64)
		if err != nil {
			return ret, err
		}
		weightedIO, err := strconv.ParseUint((fields[13]), 10, 64)
		if err != nil {
			return ret, err
		}
		d := disk.IOCountersStat{
			ReadBytes:        rbytes * sectorSize,
			WriteBytes:       wbytes * sectorSize,
			ReadCount:        reads,
			WriteCount:       writes,
			MergedReadCount:  mergedReads,
			MergedWriteCount: mergedWrites,
			ReadTime:         rtime,
			WriteTime:        wtime,
			IopsInProgress:   iopsInProgress,
			IoTime:           iotime,
			WeightedIO:       weightedIO,
		}
		if d == empty {
			continue
		}
		d.Name = name

		// Names passed in can be full paths (/dev/sda) or just device names (sda).
		// Since `name` here is already a basename, re-add the /dev path.
		// This is not ideal, but we may break the API by changing how SerialNumberWithContext
		// works.
		//d.SerialNumber, _ = SerialNumberWithContext(ctx, common.HostDevWithContext(ctx, name))
		//d.Label, _ = LabelWithContext(ctx, name)

		ret[name] = d
	}
	return ret, nil
}
func calculateAllBusy(t1, t2 []cpu.TimesStat) ([]float64, error) {
	// Make sure the CPU measurements have the same length.
	if len(t1) != len(t2) {
		return nil, fmt.Errorf(
			"received two CPU counts: %d != %d",
			len(t1), len(t2),
		)
	}

	ret := make([]float64, len(t1))
	for i, t := range t2 {
		ret[i] = calculateBusy(t1[i], t)
	}
	return ret, nil
}
func calculateBusy(t1, t2 cpu.TimesStat) float64 {
	t1All, t1Busy := getAllBusy(t1)
	t2All, t2Busy := getAllBusy(t2)

	if t2Busy <= t1Busy {
		return 0
	}
	if t2All <= t1All {
		return 100
	}
	return math.Min(100, math.Max(0, (t2Busy-t1Busy)/(t2All-t1All)*100))
}

func getAllBusy(t cpu.TimesStat) (float64, float64) {
	tot := t.Total()
	//if runtime.GOOS == "linux" {
	tot -= t.Guest     // Linux 2.6.24+
	tot -= t.GuestNice // Linux 3.2.0+
	//}

	busy := tot - t.Idle - t.Iowait

	return tot, busy
}
