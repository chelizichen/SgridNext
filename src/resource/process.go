package resource

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo 用于存储我们需要的进程信息
type ProcessInfo struct {
	PID           int32
	Name          string
	CPUPercent    float64
	MemoryPercent float32
	RSS           uint64 // Resident Set Size (常驻内存大小)
	Cmdline       string // 命令行
}

func GetProcessInfo(workDir string) []ProcessInfo {
	// 获取所有进程
	processes, err := process.Processes()
	if err != nil {
		fmt.Println("Error getting processes:", err)
		return nil
	}

	var processList []ProcessInfo

	for _, p := range processes {
		name, _ := p.Name()
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()
		memInfo, _ := p.MemoryInfo()
		cmdline, _ := p.Cmdline()
		var rss uint64
		if memInfo != nil {
			rss = memInfo.RSS
		}

		processInfo := ProcessInfo{
			PID:           p.Pid,
			Name:          name,
			CPUPercent:    cpuPercent,
			MemoryPercent: memPercent,
			RSS:           rss,
			Cmdline:       cmdline,
		}
		// 确保是业务服务，将其他的屏蔽
		if strings.Contains(cmdline, workDir) {
			processList = append(processList, processInfo)
		}
	}

	sort.Slice(processList, func(i, j int) bool {
		return processList[i].CPUPercent > processList[j].CPUPercent
	})

	sort.Slice(processList, func(i, j int) bool {
		return processList[i].RSS > processList[j].RSS
	})
	return processList
}
