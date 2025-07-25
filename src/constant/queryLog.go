package constant

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const (
	HEAD = 1
	TAIL = 2
)

func QueryLog(logFile string, logType int32, keyword string, len int32) ([]string, error) {
	log_type := ""
	log_cmd := ""
	if logType == HEAD {
		log_type = "head"
	}
	if logType == TAIL {
		log_type = "tail"
	}
	var cmd *exec.Cmd
	if keyword == "" {
		log_cmd = fmt.Sprintf("%s -500000 %s|%s -%d | iconv -c -f UTF-8 -t UTF-8|sed 's/[cA-cZ]//g'", log_type, logFile, log_type, len)
		// 如果没有提供关键词，只截取文件末尾的内容
		cmd = exec.Command("sh", "-c", log_cmd)
	} else {
		log_cmd = fmt.Sprintf("%s -500000 %s |%s -a %s|tail -%d | iconv -c -f UTF-8 -t UTF-8|sed 's/[\\cA-\\cZ]//g'", log_type, logFile, "grep", keyword, len)
		// 如果提供了关键词，先截取文件末尾内容，再筛选包含关键词的行
		cmd = exec.Command("sh", "-c", log_cmd)
	}

	// 创建一个字节缓冲区来存储命令执行的输出
	var out bytes.Buffer
	cmd.Stdout = &out

	// 执行命令
	err := cmd.Run()
	if err != nil {
		fmt.Printf("执行命令时出错: %v\n", err)
		return nil, err
	}

	// 将输出按行分割并过滤空行
	output := out.String()
	lines := strings.Split(output, "\n")
	result := make([]string, 0)

	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}

	fmt.Println(log_cmd)
	return result, nil
}
