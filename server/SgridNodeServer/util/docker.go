package util

import (
	"os/exec"
	"strings"

	"sgridnext.com/src/logger"
)

func DockerGetAlive(dockerName string) (bool, error) {
	logger.Docker.Infof("check docker alive: %s \n", dockerName)

	// 使用 docker ps 命令获取容器列表
	cmd := exec.Command("docker", "ps", "--format", "'{{.Names}}'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Docker.Errorf("check docker alive error: %s \n", err)
		return false, err
	}
	// 检查输出中是否包含目标容器名称
	outputStr := string(output)
	names := strings.Split(outputStr, "\n")
	logger.Docker.Infof("docker ps names: %s \n", names)
	dockerName = strings.Replace(dockerName, "'", "", -1)
	for _, name := range names {
		name = strings.Replace(name, "'", "", -1)
		if name == dockerName {
			logger.Docker.Infof("matched : %s \n", name)
			return true, nil
		} else {
			logger.Docker.Infof("not match: %s \n", name)
		}
	}
	return false, nil
}

func DockerStop(dockerName string) error {
	cmd := exec.Command("docker", "stop", dockerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Docker.Errorf("docker stop error: %v", err)
		return err
	}
	logger.Docker.Infof("docker stop output: %s", string(output))
	return nil
}
