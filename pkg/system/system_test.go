package system

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"testing"
	"time"
)

func TestSystem(t *testing.T) {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Info()
	cc, _ := cpu.Percent(time.Second, false)
	d, _ := disk.Usage("/")
	n, _ := host.Info()
	nv, _ := net.IOCounters(true)
	bootTime, _ := host.BootTime()
	bTime := time.Unix(int64(bootTime), 0).Format("2006-01-02 15:04:05")
	fmt.Printf("Mem : %v MB Free: %v MB Used:%v Usage:%f%%\n", v.Total/1024/1024, v.Available/1024/1024, v.Used/1024/1024, v.UsedPercent)
	for _, subCpu := range c {
		modelName := subCpu.ModelName
		cores := subCpu.Cores
		fmt.Printf("CPU : %v %v cores \n", modelName, cores)
	}
	fmt.Printf("SystemBoot:%v\n", bTime)
	fmt.Printf("Network: %v bytes / %v bytes\n", nv[0].BytesRecv, nv[0].BytesSent)
	fmt.Printf("CPU Used : used %f%% \n", cc[0])
	fmt.Printf("SSD : %v GB Free: %v GB Usage:%f%%\n", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)
	fmt.Printf("OS : %v(%v) %v \n", n.Platform, n.PlatformFamily, n.PlatformVersion)
	fmt.Printf("Hostname : %v \n", n.Hostname)
}

func TestSystemInfo(t *testing.T) {
	info, _ := GetMonitorData()
	bs, _ := json.Marshal(info)
	println(len(bs))
	println(string(bs))
}
