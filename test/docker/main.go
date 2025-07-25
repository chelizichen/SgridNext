package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GOOS=linux GOARCH=amd64 go build -o docker_alive
// docker ps --filter name=gk-collector | grep gk-collector
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
	dockerName = strings.Replace(dockerName, "'", "", -1)
	fmt.Println("names: ", names)
	for _, name := range names {
		fmt.Println("name: ", name)
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

func main() {
	// 名称通过参数传递
	if len(os.Args) < 2 {
		fmt.Println("请输入容器名称")
		return
	}
	DockerGetAlive(os.Args[1])
}
