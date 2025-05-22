package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 自动发布
type T_Config map[string]string

func parseConf(conf_path string) T_Config {
	config := make(T_Config)
	cwd, _ := os.Getwd()
	file, err := os.Open(filepath.Join(cwd, conf_path))
	if err != nil {
		fmt.Println("open file error: ", err)
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
	conf := parseConf("./sgridnext.release")
	fmt.Println(conf)
}
