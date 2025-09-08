package constant

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// escapeGrepKeyword escapes special characters in grep keywords
func escapeGrepKeyword(keyword string) string {
	// Characters that need to be escaped in grep: . * ? + [ ] ( ) { } ^ $ \ |
	specialChars := []string{".", "*", "?", "+", "[", "]", "(", ")", "{", "}", "^", "$", "\\", "|"}
	
	result := keyword
	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, "\\"+char)
	}
	
	return result
}

const (
	HEAD = 1
	TAIL = 2
)

func QueryLog(logFile string, logType int32, keyword string, lineCount int32) ([]string, error) {
	log_type := ""
	if logType == HEAD {
		log_type = "head"
	} else if logType == TAIL {
		log_type = "tail"
	} else {
		return nil, fmt.Errorf("无效的日志类型: %d", logType)
	}

	// 构建基础命令
	baseCmd := fmt.Sprintf("%s -500000 %s", log_type, logFile)
	var cmdStr string

	if keyword == "" {
		// 如果没有提供关键词，只截取文件内容
		cmdStr = fmt.Sprintf("%s | %s -%d | iconv -c -f UTF-8 -t UTF-8 | sed 's/[cA-cZ]//g'", baseCmd, log_type, lineCount)
	} else if strings.Contains(keyword, "+") {
		// 处理多个grep查询
		keywords := strings.Split(keyword, "+")
		
		// 构建grep管道命令
		grepCmd := baseCmd
		for _, kw := range keywords {
			// 去除关键词前后的空格并确保关键词不为空
			kw = strings.TrimSpace(kw)
			if kw != "" {
				// 转义grep关键词中的特殊字符
				escapedKw := escapeGrepKeyword(kw)
				grepCmd = fmt.Sprintf("%s | grep -a '%s'", grepCmd, escapedKw)
			}
		}
		
		// 添加尾部处理
		cmdStr = fmt.Sprintf("%s | tail -%d | iconv -c -f UTF-8 -t UTF-8 | sed 's/[\\cA-\\cZ]//g'", grepCmd, lineCount)
	} else {
		// 单个关键词查询
		escapedKeyword := escapeGrepKeyword(keyword)
		cmdStr = fmt.Sprintf("%s | grep -a '%s' | tail -%d | iconv -c -f UTF-8 -t UTF-8 | sed 's/[\\cA-\\cZ]//g'", 
			baseCmd, escapedKeyword, lineCount)
	}

	// 执行命令
	cmd := exec.Command("sh", "-c", cmdStr)

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
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}

	// 打印执行的命令，便于调试
	fmt.Println(cmdStr)
	return result, nil
}
