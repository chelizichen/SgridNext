package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/cgroupmanager"
	"sgridnext.com/src/logger"
)

type Command struct {
	cmd        *exec.Cmd
	serverName string
	pid        int
	cgroupMgr  *cgroupmanager.CgroupManager
}

func NewServerCommand(serverName string) *Command {
	return &Command{
		serverName: serverName,
	}
}

func (c *Command) GetCmd() *exec.Cmd {
	return c.cmd
}

func (c *Command) SetPid(pid int) {
	c.pid = pid
}

func (c *Command) SetCgroup(name string) error {
	mgr, err := cgroupmanager.NewCgroupManager(name)
	if err != nil {
		return err
	}
	c.cgroupMgr = mgr
	return nil
}

func (c *Command) SetCPULimit(cpuShares uint64) error {
	if c.cgroupMgr == nil {
		return fmt.Errorf("cgroup manager not initialized")
	}
	return c.cgroupMgr.SetCPULimit(cpuShares)
}

func (c *Command) SetMemoryLimit(memoryLimit int64) error {
	if c.cgroupMgr == nil {
		return fmt.Errorf("cgroup manager not initialized")
	}
	return c.cgroupMgr.SetMemoryLimit(memoryLimit)
}

func (c *Command) SetCommand(cmd string, args ...string) {
	logger.CMD.Infof("s.cmd: %s | args: %s \n", cmd, args)
	cwd, _ := os.Getwd()
	c.cmd = exec.Command(cmd, args...)
	c.cmd.Env = append(c.cmd.Env,
		fmt.Sprintf("%s=%s", constant.SGRID_LOG_DIR, filepath.Join(cwd, constant.TARGET_LOG_DIR)),
		fmt.Sprintf("%s=%s", constant.SGRID_CONF_DIR, filepath.Join(cwd, constant.TAGET_CONF_DIR)),
		fmt.Sprintf("%s=%s", constant.SGRID_PACKAGE_DIR, filepath.Join(cwd, constant.TARGET_PACKAGE_DIR)),
		fmt.Sprintf("%s=%s", constant.SGRID_SERVANT_DIR, filepath.Join(cwd, constant.TARGET_SERVANT_DIR)),
	)
	logger.CMD.Infof("s.cmd.Env: %s \n", c.cmd.Env)
	c.cmd.Dir = filepath.Join(cwd, constant.TARGET_SERVANT_DIR, c.serverName)
	logger.CMD.Infof("s.cmd.Dir: %s \n", c.cmd.Dir)
}

func (c *Command) AppendEnv(kvarr []string) {
	c.cmd.Env = append(c.cmd.Env, kvarr...)
}

func (c *Command) Start() error {
	if c.cmd == nil {
		return fmt.Errorf("command not initialized")
	}
	if c.cgroupMgr == nil {
		return fmt.Errorf("cgroup manager not initialized")
	}
	err := c.cgroupMgr.AddProcess(c.pid)
	if err != nil {
		return err
	}
	return c.cmd.Start()
}

// func main() {
// 	cmd := NewServerCommand("test")
// 	cmd.SetCgroup("test-group")           // 初始化cgroup
// 	cmd.SetCPULimit(512)                  // 设置CPU份额
// 	cmd.SetMemoryLimit(1024 * 1024 * 100) // 设置内存限
// }
