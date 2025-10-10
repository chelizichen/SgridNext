package config_test

import (
	"fmt"
	"testing"

	"sgridnext.com/src/config"
)

func Test_ReadJson(t *testing.T) {
	data := config.LoadConfig("./test.json")
	nodeStatus := data.GetFloat64("nodeStatus")
	fmt.Println("nodeStatus",nodeStatus)
	// nodeStatusInt, err := strconv.Atoi(nodeStatus)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("nodeStatusInt",nodeStatusInt)
}
