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

const (
	KEY_CONFIG_PATH = "__configPath__"
	KEY_NODE_INDEX = "nodeIndex"
	KEY_HOST = "host"
)

func (c Config) Get(args ...string) string {
	defer func() {
		if err := recover(); err != nil {
			logger.App.Errorf("config get string error: %v, key: %v", err, args)
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

func (c Config) GetFloat64(args ...string) float64 {
	defer func() {
		if err := recover(); err != nil {
			logger.App.Errorf("config get float64 error: %v, key: %v", err, args)
		}
	}()
	if len(args) == 0 {
		return 0
	}
	conf := c
	for i, arg := range args {
		if i == len(args)-1 {
			return conf[arg].(float64)
		}
		conf = conf[arg].(map[string]interface{})
	}
	return 0
}

// 获取最新配置的值
func (c Config)GetNewest(key string) string {
	filePath := c.Get(KEY_CONFIG_PATH)
	// 如果 filePath 为空，则返回空
	if filePath == "" {
		return ""
	}
	conf := loadJsonConfig(filePath)
	return conf.Get(key)
}

func (c Config) Set(key string, value interface{}) {
	c[key] = value
	// 保存到文件
	c.Save()
}

func (c Config) Save() {
	// 保存到文件
	filePath := c.Get(KEY_CONFIG_PATH)
	if filePath == "" {
		return
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	encoder.Encode(c)
}

func (c Config) GetLocalNodeId() int {
	nodeId, err := strconv.Atoi(c.Get(KEY_NODE_INDEX))
	if err != nil {
		panic("本地节点ID获取失败")
	}
	return nodeId
}

var Conf = make(Config)

func LoadConfig(fileName string) Config {
	Conf = loadJsonConfig(fileName)
	Conf.Set(KEY_CONFIG_PATH, fileName)
	return Conf
}

func loadJsonConfig(filePath string) Config {
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
