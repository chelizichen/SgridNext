package config

import (
	"encoding/json"
	"io"
	"os"
)

var Conf = make(map[string]interface{})

func LoadConfig(fileName string) {
	Conf = ReadJson(fileName)
}

func ReadJson(filePath string) map[string]interface{} {
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
