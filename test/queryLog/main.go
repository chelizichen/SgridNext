package main

import (
	"fmt"

	"sgridnext.com/src/constant"
)

func main() {
	logRsp, err := constant.QueryLog("./a.log", constant.HEAD, "心跳正常+b07e1966", 100)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("logRsp >> ", logRsp)
}
