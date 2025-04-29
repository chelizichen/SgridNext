package config

import (
	"encoding/json"
	"io"
	"os"

	"sgridnext.com/src/logger"
)

type Config map[string]interface{}

func (c Config) Get(args... string) string {
	defer func() {
		if err := recover(); err!= nil {
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

var Conf = make(Config)

func LoadConfig(fileName string)Config{
	Conf = ReadJson(fileName)
	return Conf
}

func ReadJson(filePath string)Config {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var result map[string]interface{}
	json.Unmarshal(byteValue, &result)

	return result
}
