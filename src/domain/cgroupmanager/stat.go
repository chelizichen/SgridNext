package cgroupmanager

import (
	"fmt"
	"time"

	"github.com/containerd/cgroups"
)

// CgroupStats 统一的 cgroup 统计信息结构体
type CgroupStats struct {
	CPU     CPUStats    `json:"cpu"`
	Memory  MemoryStats `json:"memory"`
	IO      IOStats     `json:"io,omitempty"` // v2 可能没有单独的 IO 统计
	Pids    PidsStats   `json:"pids,omitempty"`
	Version string      `json:"version"` // "v1" 或 "v2"
	Time    time.Time   `json:"time"`
}

type CPUStats struct {
	Usage       uint64  `json:"usage"`       // CPU 使用时间（纳秒）
	UsagePerSec float64 `json:"usagePerSec"` // 每秒 CPU 使用率
	Shares      uint64  `json:"shares"`      // CPU 权重（v1 shares 或 v2 weight）
	Throttled   uint64  `json:"throttled"`   // 被限制的次数
}

type MemoryStats struct {
	Usage     uint64 `json:"usage"`     // 当前内存使用量（字节）
	Limit     uint64 `json:"limit"`     // 内存限制（字节）
	Cache     uint64 `json:"cache"`     // 缓存使用量
	SwapUsage uint64 `json:"swapUsage"` // Swap 使用量
	SwapLimit uint64 `json:"swapLimit"` // Swap 限制
	OOMEvents uint64 `json:"oomEvents"` // OOM 发生次数
}

type IOStats struct {
	ReadBytes  uint64 `json:"readBytes"`
	WriteBytes uint64 `json:"writeBytes"`
}

type PidsStats struct {
	Current uint64 `json:"current"`
	Limit   uint64 `json:"limit"`
}

// ... 其他现有方法保持不变 ...

// Stat 获取 cgroup 统计信息
func (cm *CgroupManager) Stat() (*CgroupStats, error) {
	stats := &CgroupStats{
		Time: time.Now(),
	}

	if cm.isV2 {
		stats.Version = "v2"
		v2Stats, err := cm.cgroup2.Stat()
		if err != nil {
			return nil, fmt.Errorf("failed to get cgroup v2 stats: %v", err)
		}

		// CPU 统计
		if v2Stats.CPU != nil {
			stats.CPU.Usage = v2Stats.CPU.UsageUsec * 1000 // 微秒转纳秒
			stats.CPU.Throttled = v2Stats.CPU.ThrottledUsec

		}

		// 内存统计
		if v2Stats.Memory != nil {
			stats.Memory.Usage = v2Stats.Memory.Usage
			stats.Memory.Cache = v2Stats.Memory.File
			stats.Memory.SwapUsage = v2Stats.Memory.SwapUsage

		}

		// Pids 统计
		if v2Stats.Pids != nil {
			stats.Pids.Current = v2Stats.Pids.Current
			stats.Pids.Limit = v2Stats.Pids.Limit
		}

	} else {
		stats.Version = "v1"
		v1Stats, err := cm.cgroup.Stat(cgroups.IgnoreNotExist)
		if err != nil {
			return nil, fmt.Errorf("failed to get cgroup v1 stats: %v", err)
		}

		// CPU 统计
		if v1Stats.CPU != nil {
			stats.CPU.Usage = v1Stats.CPU.Usage.Total
			stats.CPU.Throttled = v1Stats.CPU.Throttling.ThrottledTime

		}

		// 内存统计
		if v1Stats.Memory != nil {
			stats.Memory.Usage = v1Stats.Memory.Usage.Usage
			stats.Memory.Cache = v1Stats.Memory.Cache
			stats.Memory.SwapUsage = v1Stats.Memory.Swap.Usage

		}

		// IO 统计 (v1 特有)
		if v1Stats.Blkio != nil {
			var readBytes, writeBytes uint64
			for _, entry := range v1Stats.Blkio.IoServiceBytesRecursive {
				if entry.Op == "Read" {
					readBytes += entry.Value
				} else if entry.Op == "Write" {
					writeBytes += entry.Value
				}
			}
			stats.IO = IOStats{
				ReadBytes:  readBytes,
				WriteBytes: writeBytes,
			}
		}

		// Pids 统计
		if v1Stats.Pids != nil {
			stats.Pids.Current = v1Stats.Pids.Current
			stats.Pids.Limit = v1Stats.Pids.Limit
		}
	}

	// 计算 CPU 使用率
	stats.CPU.UsagePerSec = calculateCPUUsagePerSec(stats.CPU.Usage)

	return stats, nil
}

// calculateCPUUsagePerSec 计算 CPU 使用率（假设上次调用是1秒前）
func calculateCPUUsagePerSec(usage uint64) float64 {
	// 这里需要实现实际的 CPU 使用率计算
	// 简单示例：需要保存上次的 usage 值和时间戳
	// 这里只是返回一个示例值
	return float64(usage) / 1e9 // 转换为秒
}

// ... 其他方法保持不变 ...
