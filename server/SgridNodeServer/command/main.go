package command

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"

	"sgridnext.com/server/SgridNodeServer/util"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/logger"
)

type Command struct {
	cmd        *exec.Cmd
	serverName string
	nodeId     int
	mu         sync.Locker
	pid    		int
	host string
	port int
	localMachineId int
	serverId int
	additionalArgs string
	serverRunType int
	dockerName string
}


func NewServerCommand(serverName string) *Command {
	return &Command{
		serverName: serverName,
		mu:         &sync.Mutex{},
	}
}

func NewPidCommand(pid int,serverName string,nodeId int) *Command {
	os.FindProcess(pid)
	cmd := &Command{
		serverName: serverName,
		mu:         &sync.Mutex{},
	}
	cmd.SetPid(pid)
	cmd.SetNodeId(nodeId)
	return cmd
}

func (c *Command) SetHost(host string) {
	c.host = host
}
func (c *Command) GetHost() string {
	return c.host
}

func (c *Command) SetPort(port int) {
	c.port = port
}

func (c *Command) GetPort() int {
	return c.port
}

func (c *Command) SetLocalMachineId(localMachineId int) {
	c.localMachineId = localMachineId
}

func (c *Command) GetLocalMachineId() int {
	return c.localMachineId
}

func (c *Command) SetServerId(serverId int) {
	c.serverId = serverId
}

func (c *Command) GetServerId() int {
	return c.serverId
}

func (c *Command) GetCmd() *exec.Cmd {
	return c.cmd
}

func (c *Command) SetPid(pid int) {
	c.pid = pid
}

func (c *Command) GetPid() int {
	return c.pid
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

func (c *Command) GetAdditionalArgs() string {
	return c.additionalArgs
}


func (c *Command) SetAdditionalArgs(additionalArgs string) {
	c.additionalArgs = additionalArgs
}

func (c *Command) GetRunServerType() int {
	return c.serverRunType
}

func (c *Command) SetRunServerType(t int)  {
	 c.serverRunType = t
}

func (c *Command) SetDockerName(dockerName string) {
	c.dockerName = dockerName
}

func (c *Command) GetDockerName() string {
	return c.dockerName
}

func (c *Command) SetCommand(cmd string, args ...string) error {
	logger.CMD.Infof("s.cmd: %s | args: %s \n", cmd, args)
	cwd, _ := os.Getwd()
	c.cmd = exec.Command(cmd, args...)
	c.cmd.Env = append(c.cmd.Env,os.Environ()...)
	logger.CMD.Info("debug >> c.host %s| c.port %d",c.host,c.port)
	// 初始化环境变量
	c.cmd.Env = append(c.cmd.Env,
		fmt.Sprintf("%s=%s", constant.SGRID_LOG_DIR, filepath.Join(cwd, constant.TARGET_LOG_DIR, c.serverName)),
		fmt.Sprintf("%s=%s", constant.SGRID_CONF_DIR, filepath.Join(cwd, constant.TARGET_CONF_DIR, c.serverName)),
		fmt.Sprintf("%s=%s", constant.SGRID_PACKAGE_DIR, filepath.Join(cwd, constant.TARGET_PACKAGE_DIR, c.serverName)),
		fmt.Sprintf("%s=%s", constant.SGRID_SERVANT_DIR, filepath.Join(cwd, constant.TARGET_SERVANT_DIR, c.serverName)),
		fmt.Sprintf("%s=%s", constant.SGRID_NODE_DIR, cwd),
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
	if c.additionalArgs != "" {
		var additionalArgs []string
		err := json.Unmarshal([]byte(c.additionalArgs), &additionalArgs)
		if err != nil {
			return err
		}
		c.AppendEnv(additionalArgs)
	}
	c.AppendEnv([]string{
		fmt.Sprintf("%s=%s", constant.SGRID_TARGET_HOST, c.host),
		fmt.Sprintf("%s=%v", constant.SGRID_TARGET_PORT, c.port),
	})
	if c.cmd == nil {
		return fmt.Errorf("command not initialized")
	}
	cwd, _ := os.Getwd()
	redirectFilePath := filepath.Join(cwd, constant.TARGET_LOG_DIR, c.serverName, c.serverName+".log")
	if err := os.MkdirAll(filepath.Dir(redirectFilePath), 0755); err != nil {
		logger.App.Errorf("创建目录失败: %v", err)
		return fmt.Errorf("创建目录失败: %v", err)
	}
	outFile, err := os.Create(redirectFilePath)
	if err != nil {
		logger.CMD.Errorf("failed to create output file: %v", err)
		return err
	}
	defer outFile.Close()
	c.cmd.Stdout = outFile
	c.cmd.Stderr = outFile

	err = c.cmd.Start()
	c.SetPid(c.cmd.Process.Pid)

	// 将命令输出重定向到文件
	if err != nil {
		logger.CMD.Errorf("failed to start command: %v", err)
		return err
	}
	return nil
}


func (c *Command) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.GetDockerName() != "" {
		err := util.DockerStop(c.GetDockerName())
		if err != nil {
			logger.CMD.Errorf("停止docker失败: %v", err)
			return err
		}
		return nil
	}
	if c.cmd == nil || c.cmd.Process == nil {
		logger.CMD.Infof("command not initialized")
		// 进程未启动，直接返回nil
		return nil
	}
	logger.CMD.Infof("正在停止进程: %d", c.GetPid())
	// groups,err := FindProcessGroup(c.GetPid())
	// if err != nil {
	// 	logger.CMD.Errorf("failed to kill process: %v", err)
	// 	return err
	// }
	// for _,pid := range groups {
	// 	logger.CMD.Infof("正在停止进程: %d", pid)
	// 	if err := Kill(pid); err != nil {
	// 		logger.CMD.Errorf("failed to kill process: %v", err)
	// 		return err
	// 	}
	// }
	// if err := c.cmd.Process.Kill(); err != nil {
	// 	if errors.Is(err, os.ErrProcessDone) {
	// 		return nil // 进程已结束
	// 	}
	// 	return fmt.Errorf("kill process failed: %w", err)
	// }
	err := Kill(c.GetPid())
	if err != nil {
		logger.CMD.Errorf("删除节点失败: %v", err)
		return err
	}
	return err
}

func (c *Command) CheckStat() (pid int, alive bool, err error) {
	if c.GetDockerName() != "" {
		alive, err = util.DockerGetAlive(c.GetDockerName())
		return 0, alive, err
	}
	// if c.cmd == nil || c.cmd.Process == nil {
	// 	return 0, false, fmt.Errorf("command not initialized")
	// }
	pid = c.GetPid()
	process, err := os.FindProcess(pid)
	if err != nil {
		return pid, false, fmt.Errorf("find process failed: %w", err)
	}
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		// 如果发送信号失败，说明进程不存在或没有权限
		logger.CMD.Infof("进程 %d 不存活: %v", pid, err)
		return pid, false, nil
	}
	return pid, true, nil
}
