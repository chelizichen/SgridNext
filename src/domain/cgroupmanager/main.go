package cgroupmanager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/containerd/cgroups"
	cgroupsv2 "github.com/containerd/cgroups/v2"
	"github.com/opencontainers/runtime-spec/specs-go"
	"sgridnext.com/src/logger"
)

type CgroupManager struct {
	cgroup  cgroups.Cgroup  // for v1
	cgroup2 *cgroupsv2.Manager // for v2
	isV2    bool
}

const (
	defaultV2MountPath = "/sys/fs/cgroup/system.slice"
)

func (cm *CgroupManager) GetCgroup() interface{} {
	if cm.isV2 {
		return cm.cgroup2
	}
	return cm.cgroup
}

func LoadCgroupManager(name string) (*CgroupManager, error) {
	// 检测 cgroup 版本
	isV2, err := isCgroupV2()
	if err != nil {
		logger.Cgroup.Errorf("failed to detect cgroup version: %v", err)
		return nil, fmt.Errorf("failed to detect cgroup version: %v", err)
	}
	logger.Cgroup.Infof("detect cgroup version is v2: %v", isV2)
	if isV2 {
		// 使用 cgroup v2
		groupPath := filepath.Join("/", name)
		logger.Cgroup.Infof("debug | load cgroup groupPath | %s ", groupPath)
		manager, err := cgroupsv2.LoadManager(defaultV2MountPath,groupPath)
		if err != nil {
			logger.Cgroup.Errorf("failed to load cgroup v2 manager: %v", err)
			return nil, fmt.Errorf("failed to load cgroup v2 manager: %v", err)
		}
		targetDir := filepath.Join(defaultV2MountPath,groupPath)
		if _, err := os.Stat(targetDir); err != nil {
			logger.Cgroup.Errorf("failed to stat cgroup v2 target directory: %v", err)
			return nil,err
		}
		logger.Cgroup.Infof("load cgroup %s ", name)
		return &CgroupManager{cgroup2: manager, isV2: true}, nil
	}
	// 使用 cgroup v1
	mountPath := cgroups.Slice("system.slice", name)
	control, err := cgroups.Load(cgroups.Systemd, mountPath)
	logger.Cgroup.Infof("load cgroup mountPath | %s ", mountPath)
	if err!= nil {
		logger.Cgroup.Errorf("failed to load cgroup v1 manager: %v", err)
		return nil, fmt.Errorf("failed to load cgroup v1 manager: %v", err)
	}
	return &CgroupManager{cgroup: control, isV2: false}, nil
}

func NewCgroupManager(name string) (*CgroupManager, error) {
	// 先加载看有没有 Cgroup目录
	manger, err := LoadCgroupManager(name)
	if err == nil {
		if manger.isV2{
			// if manger.cgroup2
		}
		logger.Cgroup.Infof("cgroup %s already exists", name)
		return manger, nil
	}
	// 检测 cgroup 版本
	isV2, err := isCgroupV2()
	if err != nil {
		return nil, fmt.Errorf("failed to detect cgroup version: %v", err)
	}

	if isV2 {
		// 使用 cgroup v2
		groupPath := filepath.Join("/", name)
		manager, err := cgroupsv2.NewManager(defaultV2MountPath, groupPath, &cgroupsv2.Resources{})
		if err != nil {
			logger.Cgroup.Errorf("failed to create cgroup v2 manager: %v", err)
			return nil, fmt.Errorf("failed to create cgroup v2 manager: %v", err)
		}
		return &CgroupManager{cgroup2: manager, isV2: true}, nil
	}

	// 使用 cgroup v1
	mountPath := cgroups.Slice("system.slice", name)
	control, err := cgroups.New(cgroups.Systemd, mountPath, &specs.LinuxResources{})
	if err != nil {
		return nil, fmt.Errorf("failed to create cgroup v1 manager: %v", err)
	}
	return &CgroupManager{cgroup: control, isV2: false}, nil
}

func (cm *CgroupManager) SetCPULimit(cores float64) error {
	if cores <= 0 {
		return fmt.Errorf("cpu cores must be positive")
	}

	if cm.isV2 {
		// cgroup v2 实现
		// 转换为 quota 和 period 格式
		// 通常 period 默认为 100000 微秒(100ms)
		period := uint64(100000)
		quota := int64(float64(period) * cores)
		logger.Cgroup.Infof("quota: %d | period: %d", quota, period)
		return cm.cgroup2.Update(&cgroupsv2.Resources{
			CPU: &cgroupsv2.CPU{
				Max: cgroupsv2.NewCPUMax(&quota, &period),
			},
		})
	}

	// cgroup v1 实现
	period := uint64(100000)
	quota := int64(float64(period) * cores)
	
	return cm.cgroup.Update(&specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Period: &period,
			Quota:  &quota,
		},
	})
}

func (cm *CgroupManager) SetMemoryLimit(memoryLimit int64) error {
	if cm.isV2 {
		return cm.cgroup2.Update(&cgroupsv2.Resources{
			Memory: &cgroupsv2.Memory{
				Max: &memoryLimit,
			},
		})
	}
	return cm.cgroup.Update(&specs.LinuxResources{
		Memory: &specs.LinuxMemory{
			Limit: &memoryLimit,
		},
	})
}

func (cm *CgroupManager) AddProcess(pid int) error {
	logger.Cgroup.Infof("add process %d to cgroup |  isV2 %s", pid,cm.isV2)
	if cm.isV2 {
		return cm.cgroup2.AddProc(uint64(pid))
	}
	return cm.cgroup.Add(cgroups.Process{Pid: pid})
}

func (cm *CgroupManager) Remove() error {
	if cm.isV2 {
		return cm.cgroup2.Delete()
	}
	return cm.cgroup.Delete()
}

// isCgroupV2 检测系统是否使用 cgroup v2
func isCgroupV2() (bool, error) {
	_, err := os.Stat("/sys/fs/cgroup/cgroup.controllers")
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// convertSharesToWeight 将 v1 的 cpu shares 转换为 v2 的 weight
func convertSharesToWeight(shares uint64) uint64 {
	// v1 shares 范围: 2-262144 (默认 1024)
	// v2 weight 范围: 1-10000 (默认 100)
	if shares == 0 {
		return 100 // 默认值
	}
	weight := shares * 10000 / 262144
	if weight < 1 {
		return 1
	}
	if weight > 10000 {
		return 10000
	}
	return weight
}

type SgridCgroup struct {
	ServerName string
	NodeId int
}

func (s *SgridCgroup)GetCgroupName() string {
	return fmt.Sprintf("sgrid-%s-%d", s.ServerName, s.NodeId)
}