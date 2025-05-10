package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"sgridnext.com/src/constant"
	"sgridnext.com/src/logger"
)

type Command struct {
	cmd        *exec.Cmd
	serverName string
	nodeId int
	mu sync.Locker
}

func NewServerCommand(serverName string) *Command {
	return &Command{
		serverName: serverName,
		mu: &sync.Mutex{},
	}
}

func (c *Command) GetCmd() *exec.Cmd {
	return c.cmd
}


func (c *Command) GetPid() int {
	// return 9999
	return c.cmd.Process.Pid
}

func (c *Command) SetNodeId(nodeId int) {
	c.nodeId = nodeId
}

func (c *Command) GetNodeId() int {
	return c.nodeId
}

func (c *Command) GetServerName() string {
	return c.serverName
}


func (c *Command) SetCommand(cmd string, args ...string) error {
	logger.CMD.Infof("s.cmd: %s | args: %s \n", cmd, args)
	cwd, _ := os.Getwd()
	c.cmd = exec.Command(cmd, args...)
	c.cmd.Env = append(c.cmd.Env,
		fmt.Sprintf("%s=%s", constant.SGRID_LOG_DIR, filepath.Join(cwd, constant.TARGET_LOG_DIR,c.serverName)),
		fmt.Sprintf("%s=%s", constant.SGRID_CONF_DIR, filepath.Join(cwd, constant.TAGET_CONF_DIR,c.serverName)),
		fmt.Sprintf("%s=%s", constant.SGRID_PACKAGE_DIR, filepath.Join(cwd, constant.TARGET_PACKAGE_DIR,c.serverName)),
		fmt.Sprintf("%s=%s", constant.SGRID_SERVANT_DIR, filepath.Join(cwd, constant.TARGET_SERVANT_DIR,c.serverName)),
	)
	logger.CMD.Infof("s.cmd.Env: %s \n", c.cmd.Env)
	c.cmd.Dir = filepath.Join(cwd, constant.TARGET_SERVANT_DIR, c.serverName)
	logger.CMD.Infof("s.cmd.Dir: %s \n", c.cmd.Dir)
	return nil
}

func (c *Command) AppendEnv(kvarr []string) {
	c.cmd.Env = append(c.cmd.Env, kvarr...)
}

func (c *Command) Start() error {
	if c.cmd == nil {
		return fmt.Errorf("command not initialized")
	}
	err := c.cmd.Start()
	if err!= nil {
		logger.CMD.Errorf("failed to start command: %v", err)
		return err
	}
	return nil
}


func (c *Command) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()
    if c.cmd == nil || c.cmd.Process == nil {
        return fmt.Errorf("command not initialized")
    }
    if err := c.cmd.Process.Kill(); err != nil {
        if errors.Is(err, os.ErrProcessDone) {
            return nil // 进程已结束
        }
        return fmt.Errorf("kill process failed: %w", err)
    }
    return nil
}

func (c *Command) CheckStat() (pid int,alive bool,err error) {
	// 检查进程状态
	if c.cmd == nil || c.cmd.Process == nil {
		return 0,false,fmt.Errorf("command not initialized")
	}
	pid = c.cmd.Process.Pid
	_, err = os.FindProcess(pid)
	if err != nil {
		return 0,false,fmt.Errorf("find process failed: %w", err)
	}
	return pid,true,nil
}