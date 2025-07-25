package util

import (
	"fmt"
	"os/exec"
	"strings"

	"sgridnext.com/src/logger"
)

func DockerGetAlive(dockerName string) (bool, error) {
	fmt.Printf("check docker alive: %s \n", dockerName)

	// 使用 docker ps 命令获取容器列表
	cmd := exec.Command("docker", "ps", "--format", "'{{.Names}}'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("check docker alive error: ", err)
		return false, err
	}
	// 检查输出中是否包含目标容器名称
	outputStr := string(output)
	names := strings.Split(outputStr, "\n")
	fmt.Println("names: ", names)
	dockerName = strings.Replace(dockerName, "'", "", -1)
	for _, name := range names {
		name = strings.Replace(name, "'", "", -1)
		if name == dockerName {
			fmt.Println("matched : ", name, "dockerName: ", dockerName)
			return true, nil
		} else {
			fmt.Println("not match: ", name, "dockerName: ", dockerName)
		}
	}
	return false, nil
}

func DockerStop(dockerName string) error {
	cmd := exec.Command("docker", "stop", dockerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Server.Errorf("docker stop error: %v", err)
		return err
	}
	logger.Server.Infof("docker stop output: %s", string(output))
	return nil
}
