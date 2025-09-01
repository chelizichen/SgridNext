// sgridnext patch server
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// 由于 SgridNode 本身就是一个服务，所以在 systemctl restart sgridnode 时会丢失，导致后续日志无法打印
// 可以在 web界面查看心跳是否正常

func main() {
	// 执行脚本
	nodeDir := os.Getenv("SGRID_NODE_DIR")
	fmt.Println("nodeDir:", nodeDir)
	cmd := exec.Command("./update.sh", nodeDir)
	// 执行脚本
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("执行脚本失败:", err)
		os.Exit(1)
	}
	fmt.Println("脚本输出结果")
	fmt.Println(out.String())
	fmt.Println("脚本执行成功")
}
