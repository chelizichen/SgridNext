package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryInfo struct {
	Total float64
	Used  float64
	Usage float64
}

type CPUInfo struct {
	Usage float64
}

type SystemInfo struct {
	MemoryInfo MemoryInfo
	CPUInfo    CPUInfo
}

func getSystemStats() {

	// 获取内存信息
	v, _ := mem.VirtualMemory()
	memoryInfo := MemoryInfo{
		Total: float64(v.Total) / (1024 * 1024),
		Used:  float64(v.Used) / (1024 * 1024),
		Usage: v.UsedPercent,
	}

	// 获取 CPU 整体占用率
	// percpu=false 表示获取总的占用率，间隔 1 秒
	cpuPercent, _ := cpu.Percent(time.Second, false)
	var cpuInfo CPUInfo
	if len(cpuPercent) > 0 {
		cpuInfo = CPUInfo{
			Usage: cpuPercent[0],
		}
	}
	systemInfo := SystemInfo{
		MemoryInfo: memoryInfo,
		CPUInfo:    cpuInfo,
	}
	fmt.Printf("System Info | %v", systemInfo)
}

func main() {
	getSystemStats()
}
