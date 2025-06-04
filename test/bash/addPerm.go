package main

import (
	"os"
	"os/exec"
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

func run(){
    var cmd *exec.Cmd = exec.Command("./TestBash.sh")
    cmd.Env = append(cmd.Env, "SGRID_TARGET_PORT=8080")
    cmd.Env = append(cmd.Env, os.Environ()...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Run()
}