package cgroupmanager

import (
	"github.com/containerd/cgroups"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type CgroupManager struct {
	cgroup cgroups.Cgroup
}

func NewCgroupManager(name string) (*CgroupManager, error) {
	control, err := cgroups.New(cgroups.V1, cgroups.StaticPath(name), &specs.LinuxResources{})
	if err != nil {
		return nil, err
	}
	return &CgroupManager{cgroup: control}, nil
}

func (cm *CgroupManager) SetCPULimit(cpuShares uint64) error {
	return cm.cgroup.Update(&specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Shares: &cpuShares,
		},
	})
}

func (cm *CgroupManager) SetMemoryLimit(memoryLimit int64) error {
	return cm.cgroup.Update(&specs.LinuxResources{
		Memory: &specs.LinuxMemory{
			Limit: &memoryLimit,
		},
	})
}

func (cm *CgroupManager) AddProcess(pid int) error {
	return cm.cgroup.Add(cgroups.Process{Pid: pid})
}
