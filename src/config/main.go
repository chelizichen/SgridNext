package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"sgridnext.com/src/logger"
)

type Config map[string]interface{}

func (c Config) Get(args ...string) string {
	defer func() {
		if err := recover(); err != nil {
			logger.App.Errorf("config get error: %v", err)
		}
	}()
	if len(args) == 0 {
		return ""
	}
	conf := c
	for i, arg := range args {
		if i == len(args)-1 {
			return conf[arg].(string)
		}
		conf = conf[arg].(map[string]interface{})
	}
	return ""
}

func (c Config) GetLocalNodeId() int {
	nodeId, err := strconv.Atoi(c.Get("nodeIndex"))
	if err != nil {
		panic("本地节点ID获取失败")
	}
	return nodeId
}

var Conf = make(Config)

func LoadConfig(fileName string) Config {
	Conf = ReadJson(fileName)
	return Conf
}

func ReadJson(filePath string) Config {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var result map[string]interface{}
	json.Unmarshal(byteValue, &result)
	fmt.Printf("Read Json Conf filepath | %s, value | %s", filePath, result)
	return result
}
