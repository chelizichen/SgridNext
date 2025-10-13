package resource

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryInfo struct {
	Total float64 `json:"total"`
	Used  float64 `json:"used"`
	Usage float64 `json:"usage"`
}

type CPUInfo struct {
	Usage float64 `json:"usage"`
}

type SystemInfo struct {
	MemoryInfo MemoryInfo `json:"memoryInfo"`
	CPUInfo    CPUInfo    `json:"cpuInfo"`
}

func GetSystemStats() SystemInfo {

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
	return systemInfo
}
