package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"sgridnext.com/src/logger"
)

func Kill(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		logger.App.Errorf("查找进程失败: %v", err)
		return fmt.Errorf("查找进程失败: %w", err)
	}

	// 先尝试发送SIGTERM信号，优雅终止
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		logger.App.Warnf("发送SIGTERM信号失败，尝试强制终止: %v", err)
		// 如果SIGTERM失败，使用SIGKILL强制终止
		err = process.Kill()
		if err != nil {
			if errors.Is(err, os.ErrProcessDone) {
				logger.App.Infof("进程 %d 已经终止", pid)
				return nil // 进程已结束
			}
			logger.App.Errorf("终止进程 %d 失败: %v", pid, err)
			return fmt.Errorf("终止进程失败: %w", err)
		}
	}

	logger.App.Infof("成功终止进程 %d", pid)
	return nil
}

func AddPerm(path string)error{
    var cmd *exec.Cmd = exec.Command("chmod","+x",path)
    err := cmd.Run()
	if err != nil {
		logger.App.Errorf("添加权限失败: %v", err)
		return err
	}
	return nil
}