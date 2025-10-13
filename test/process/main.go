package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo 用于存储我们需要的进程信息
type ProcessInfo struct {
	PID           int32
	Name          string
	CPUPercent    float64
	MemoryPercent float32
	RSS           uint64 // Resident Set Size (常驻内存大小)
	Path          string // 进程路径
}

func main() {
	// 获取所有进程
	processes, err := process.Processes()
	if err != nil {
		fmt.Println("Error getting processes:", err)
		return
	}

	var processList []ProcessInfo

	for _, p := range processes {
		name, _ := p.Name()
		// 注意：CPUPercent() 第一次调用通常返回 0.0，因为它需要时间间隔来计算变化量。
		// 在实际应用中，你需要等待一个时间间隔（例如几百毫秒）后再次调用才能获得准确的 CPU 百分比。
		// 为了演示，这里我们先获取一次。为了获取准确的 CPU 占用，你需要先调用一次，等待，再调用一次。
		// 这里我们简化处理，假设我们已经在某个时间点调用过一次。
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()
		memInfo, _ := p.MemoryInfo()
		path, _ := p.Cmdline()
		var rss uint64
		if memInfo != nil {
			rss = memInfo.RSS
		}

		processList = append(processList, ProcessInfo{
			PID:           p.Pid,
			Name:          name,
			CPUPercent:    cpuPercent,
			MemoryPercent: memPercent,
			RSS:           rss,
			Path:          path,
		})
	}

	// ----------------- CPU 占用高的进程 -----------------
	// 按 CPU 百分比降序排序
	sort.Slice(processList, func(i, j int) bool {
		return processList[i].CPUPercent > processList[j].CPUPercent
	})

	fmt.Println("--- Top 10 Processes by CPU Usage ---")
	for i, info := range processList {
		if i >= 10 {
			break
		}
		fmt.Printf("PID: %d, Name: %s, CPU: %.2f%%, Mem: %.2f%%, Path: %s\n",
			info.PID, info.Name, info.CPUPercent, info.MemoryPercent)
	}

	// ----------------- 内存占用高的进程 -----------------
	// 按 RSS (常驻内存大小) 降序排序
	sort.Slice(processList, func(i, j int) bool {
		return processList[i].RSS > processList[j].RSS
	})

	fmt.Println("\n--- Top 10 Processes by Memory (RSS) Usage ---")
	for i, info := range processList {
		if i >= 10 {
			break
		}
		// 将 RSS 转换为 MB 显示
		rssMB := float64(info.RSS) / (1024 * 1024)
		fmt.Printf("PID: %d, Name: %s, CPU: %.2f%%, Mem: %.2f%%, RSS: %.2f MB, Path: %s\n",
			info.PID, info.Name, info.CPUPercent, info.MemoryPercent, rssMB, info.Path)
	}
}

// **注意 CPU 百分比的获取**
// 为了获取准确的 CPU 占用率 (例如在一个 1 秒的间隔内)，你需要做以下操作：
// 1. 调用 p.CPUPercent() (第一次调用)
// 2. time.Sleep(time.Second) (等待一个时间间隔)
// 3. 再次调用 p.CPUPercent() (第二次调用，会返回自上次调用以来计算的百分比)

// 示例：获取单个进程的准确 CPU 占用率
func getAccurateCPUPercent(p *process.Process) (float64, error) {
	// 第一次调用，用于初始化
	_, err := p.CPUPercent()
	if err != nil {
		return 0, err
	}
	time.Sleep(500 * time.Millisecond) // 等待 500ms

	// 第二次调用，获取实际的百分比
	return p.CPUPercent()
}
