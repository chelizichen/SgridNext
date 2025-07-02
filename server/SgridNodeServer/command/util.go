package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"sgridnext.com/src/constant"
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

func FindProcessGroup(nodePid int)([]int,error){
	var rsp []int = make([]int, 0)
	rsp = append(rsp, nodePid)
	// 1. 获取进程的 PGID
	// ps -o pgid= -p 1615345
	// ps -o pid= -g 1615345 
	// cmd := exec.Command("ps", "-o", "pgid=", "-p", fmt.Sprintf("%d", nodePid))
	// out, err := cmd.Output()
	// if err != nil {
	// 	logger.CMD.Infof("获取 PGID 失败: %v\n", err)
	// 	return rsp,err
	// }

	// pgid := strings.TrimSpace(string(out))
	// logger.CMD.Infof("Node 进程 (PID=%d) 的 PGID: %s\n", nodePid, pgid)
	// 2. 查找该 PGID 下的所有进程
	cmd := exec.Command("ps", "-o", "pid=", "-g", fmt.Sprintf("%d",nodePid))
	out, err := cmd.Output()
	if err != nil {
		logger.CMD.Infof("查找进程组失败: %v\n", err)
		return rsp,err
	}

	pids := strings.Split(strings.TrimSpace(string(out)), "\n")
	sgridnodePid := os.Getpid()
	logger.CMD.Infof("sgridnodePid %s 进程 PID :\n", sgridnodePid)
	logger.CMD.Infof("进程组 %s 下的所有进程:\n", nodePid)
	for _, pid := range pids {
		logger.CMD.Infof("- PID: %s\n", strings.TrimSpace(pid))
		cpid,err := strconv.Atoi(strings.TrimSpace(pid))
		if err !=nil{
			logger.CMD.Infof("PID 转换失败: %v\n", err)
		}
		if cpid == sgridnodePid {
			logger.CMD.Infof("忽略 节点PID: %s\n", cpid)
			continue
		}
		rsp = append(rsp, cpid)
	}
	rsp = constant.DeduplicateInts(rsp)
	logger.CMD.Infof("进程组 %s 下的进程数: %d\n", nodePid, len(rsp))
	return rsp,nil
}