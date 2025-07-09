package util

import (
	"fmt"
	"os/exec"

	"sgridnext.com/src/logger"
)

func DockerGetAlive(dockerName string) (bool, error) {
	logger.Server.Infof("check docker alive: %s", dockerName)
	cmd := exec.Command(
		"docker",
		"ps",
		"--filter",
		fmt.Sprintf("name=%s", dockerName),
		" | grep ",
		dockerName,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Server.Errorf("check docker alive error: %v", err)
		return false, err
	}
	if len(output) == 0 {
		logger.Server.Infof("check docker alive output is empty")
		return false, nil
	}
	logger.Server.Infof("check docker alive output: %s", string(output))
	return true, nil
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
