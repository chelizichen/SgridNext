package main

import (
	"sgridnext.com/src/cgroupmanager"
	"sgridnext.com/src/command"
)

func setCpuTest(){
	m, err := cgroupmanager.LoadCgroupManager("sgrid-TestNodeServer-2")
	if err != nil {
		panic(err)
	}
	m.SetCPULimit(1.5)
}

func useHookTest(){
	cmd := command.NewServerCommand("TestNodeServer")
	cmd.SetNodeId(2)
	command.UseCgroup(cmd)
}

func main() {
	// useHookTest()
	// setCpuTest()
}
