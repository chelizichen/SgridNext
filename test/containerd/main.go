package main

import (
	"sgridnext.com/src/domain/cgroupmanager"
)

func main() {
	m, err := cgroupmanager.LoadCgroupManager("TestHighCpuServer")
	if err != nil {
		panic(err)
	}
	m.SetCPULimit(1.5)
}
