package main

import (
	"bufio"
	"bytes"
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
		panic("open file error:"+err.Error())
		// return nil
	}
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
	cwd, _ := os.Getwd()
	conf := parseConf("sgridnext.release")
	fmt.Println(conf)
	file_path := conf["PACKAGE_PATH"]
	serverName := conf["SERVER_NAME"]
	serverId := conf["SERVER_ID"]
	deploy_path := conf["DEPLOY_PATH"]
	deploy_file_path :=  filepath.Join(cwd, file_path)
	// http 请求
	fmt.Println("开始部署服务: ", serverName, "  服务ID: ", serverId, " 部署路径: ", deploy_path)
	client := &http.Client{}
	apiPath := fmt.Sprintf("%s/api/server/scripts/deploy", deploy_path)
	fmt.Printf("apiPath: %s \n", apiPath)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("serverName", serverName)
	writer.WriteField("serverId", serverId)
	file, err := os.Open(deploy_file_path)
	if err != nil {
		fmt.Println("打开文件失败: ", err)
		return
	}
	part, err := writer.CreateFormFile("file",deploy_file_path)
	if err != nil {
		fmt.Println("创建文件失败: ", err)
		return
	}
	io.Copy(part, file)

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
