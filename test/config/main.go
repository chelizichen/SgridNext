package main

import (
	"fmt"
	"time"

	"sgridnext.com/src/domain/patchutils"
)

func Test_Config() {
	p1, err := patchutils.T_PatchUtils.UpdateConfigFileContent("TestHighCpuServer", "config.json", `{"hello":"world"}`)
	fmt.Println(p1, err)
	time.Sleep(time.Second * 1)
	p2, err := patchutils.T_PatchUtils.UpdateConfigFileContent("TestHighCpuServer", "config.json", `{"hello":"world222"}`)
	fmt.Println(p2, err)
}

func Test_Back() {
	patchutils.T_PatchUtils.BackConfigFile("TestHighCpuServer", "config.json", `config_1746610013.json`)
}

func main() {
	// Test_Config()
	Test_Back()
}
