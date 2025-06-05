package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)


func main(){
    // removePerm()
    addPerm()
    run()
}

func addPerm(){
    var cmd *exec.Cmd = exec.Command("chmod","+x","./TestBash.sh")
    cmd.Run()
}

func removePerm(){
    var cmd *exec.Cmd = exec.Command("chmod","-x","./TestBash.sh")
    cmd.Run()
}

func runBash(){
    var cmd *exec.Cmd = exec.Command("./TestBash.sh")
    cmd.Env = append(cmd.Env, "SGRID_TARGET_PORT=8080")
    cmd.Env = append(cmd.Env, os.Environ()...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Run()
}

func runNode(){
    var cmd *exec.Cmd = exec.Command("../nodeserver/bashnode.sh")
    cmd.Env = append(cmd.Env, "SGRID_TARGET_PORT=8089")
    cmd.Env = append(cmd.Env, os.Environ()...)
    cmd.Dir = "../nodeserver/"
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Start()
    fmt.Println("cmd.Process.Pid >> ",cmd.Process.Pid)
}

func run(){
    // runNode()
    findProcess()
}

func findProcess(){
    nodePid := 19419

	// 1. 获取进程的 PGID
	cmd := exec.Command("ps", "-o", "pgid=", "-p", fmt.Sprintf("%d", nodePid))
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("获取 PGID 失败: %v\n", err)
		return
	}

	pgid := strings.TrimSpace(string(out))
	fmt.Printf("Node 进程 (PID=%d) 的 PGID: %s\n", nodePid, pgid)

	// 2. 查找该 PGID 下的所有进程
	cmd = exec.Command("ps", "-o", "pid=", "-g", pgid)
	out, err = cmd.Output()
	if err != nil {
		fmt.Printf("查找进程组失败: %v\n", err)
		return
	}

	pids := strings.Split(strings.TrimSpace(string(out)), "\n")
	fmt.Printf("进程组 %s 下的所有进程:\n", pgid)
	for _, pid := range pids {
		fmt.Printf("- PID: %s\n", strings.TrimSpace(pid))
	}
}

// 18140 kill -9
// 18142