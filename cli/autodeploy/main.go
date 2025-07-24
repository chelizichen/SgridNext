package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 自动发布
type T_Config map[string]string

func parseConf(conf_path string) T_Config {
	config := make(T_Config)
	cwd, _ := os.Getwd()
	conf_path = filepath.Join(cwd, conf_path)
	fmt.Println("conf_path: ", conf_path)
	file, err := os.Open(conf_path)
	if err != nil {
		fmt.Println("open file error: ", err)
		// 如果配置文件不存在，返回空配置
		return config
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return config
}

func main() {
	// 定义命令行参数
	var (
		deployPath  = flag.String("DEPLOY_PATH", "", "部署路径 (例如: http://localhost:8080)")
		serverId    = flag.String("SERVER_ID", "", "服务器ID")
		serverName  = flag.String("SERVER_NAME", "", "服务器名称")
		packagePath = flag.String("PACKAGE_PATH", "", "包文件路径")
		commit      = flag.String("COMMIT", "", "提交信息 (可选)")
		configFile  = flag.String("config", "sgridnext.release", "配置文件路径")
	)
	
	flag.Parse()
	
	// 读取配置文件
	conf := parseConf(*configFile)
	
	// 命令行参数优先级高于配置文件
	finalDeployPath := *deployPath
	if finalDeployPath == "" {
		finalDeployPath = conf["DEPLOY_PATH"]
	}
	
	finalServerId := *serverId
	if finalServerId == "" {
		finalServerId = conf["SERVER_ID"]
	}
	
	finalServerName := *serverName
	if finalServerName == "" {
		finalServerName = conf["SERVER_NAME"]
	}
	
	finalPackagePath := *packagePath
	if finalPackagePath == "" {
		finalPackagePath = conf["PACKAGE_PATH"]
	}
	
	finalCommit := *commit
	if finalCommit == "" {
		finalCommit = conf["COMMIT"]
	}
	if finalCommit == "" {
		finalCommit = "auto deploy" // 默认提交信息
	}
	
	// 验证必需参数
	if finalDeployPath == "" {
		fmt.Println("错误: DEPLOY_PATH 参数是必需的")
		flag.Usage()
		os.Exit(1)
	}
	if finalServerId == "" {
		fmt.Println("错误: SERVER_ID 参数是必需的")
		flag.Usage()
		os.Exit(1)
	}
	if finalServerName == "" {
		fmt.Println("错误: SERVER_NAME 参数是必需的")
		flag.Usage()
		os.Exit(1)
	}
	if finalPackagePath == "" {
		fmt.Println("错误: PACKAGE_PATH 参数是必需的")
		flag.Usage()
		os.Exit(1)
	}
	
	cwd, _ := os.Getwd()
	deployFilePath := finalPackagePath
	// 如果是相对路径，则相对于当前工作目录
	if !filepath.IsAbs(deployFilePath) {
		deployFilePath = filepath.Join(cwd, finalPackagePath)
	}
	
	// 检查文件是否存在
	if _, err := os.Stat(deployFilePath); os.IsNotExist(err) {
		fmt.Printf("错误: 文件不存在: %s\n", deployFilePath)
		os.Exit(1)
	}
	
	// http 请求
	fmt.Println("开始部署服务: ", finalServerName, "  服务ID: ", finalServerId, " 部署路径: ", finalDeployPath)
	fmt.Println("包文件路径: ", deployFilePath)
	fmt.Println("提交信息: ", finalCommit)
	
	client := &http.Client{}
	apiPath := fmt.Sprintf("%s/api/server/scripts/deploy", finalDeployPath)
	fmt.Printf("apiPath: %s \n", apiPath)
	
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("serverName", finalServerName)
	writer.WriteField("serverId", finalServerId)
	writer.WriteField("commit", finalCommit)
	
	file, err := os.Open(deployFilePath)
	if err != nil {
		fmt.Println("打开文件失败: ", err)
		return
	}
	defer file.Close()
	
	part, err := writer.CreateFormFile("file", filepath.Base(deployFilePath))
	if err != nil {
		fmt.Println("创建文件失败: ", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("复制文件失败: ", err)
		return
	}
	
	writer.Close()
	req, err := http.NewRequest("POST", apiPath, body)
	if err != nil {
		fmt.Println("创建请求失败: ", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("部署失败: ", err)
		return
	}
	defer resp.Body.Close()
	rsp, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败: ", err)
		return
	}
	fmt.Println("响应: ", string(rsp))
}
