package command

import (
	"sgridnext.com/server/SgridNodeServer/cgroupmanager"
	"sgridnext.com/src/logger"
)

func UseCgroup(c *Command) error {
	scg := &cgroupmanager.SgridCgroup{
		ServerName: c.serverName,
		NodeId:     c.nodeId,
	}
	cgName := scg.GetCgroupName()
	mgr, err := cgroupmanager.NewCgroupManager(cgName)
	if err != nil {
		logger.Hook_Cgroup.Errorf("failed to create cgroup manager | server %s | error : %s", cgName, err.Error())
		return err
	}
	logger.Hook_Cgroup.Infof("load cgroup manager | server %s", cgName)
	err = mgr.AddProcess(c.GetPid())
	if err != nil {
		logger.Hook_Cgroup.Errorf("failed to add process to cgroup | server %s | error : %s", cgName, err.Error())
		return err
	}
	logger.Hook_Cgroup.Infof("load cgroup manager | server %s | for pid: %d", cgName, c.GetPid())
	return nil
}
