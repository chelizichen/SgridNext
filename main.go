package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"sgridnext.com/src/config"
	"sgridnext.com/src/db"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/domain/service/routes"
	"sgridnext.com/src/proxy"
)

func main() {
	conf := config.LoadConfig("./config.json")
	ormDb, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	mapper.LoadMapper(ormDb)
	proxy.LoadProxy()
	// 初始化路由
	engine := gin.Default()
	routes.LoadRoutes(engine)
	port := fmt.Sprintf(":%s", conf.Get("httpPort"))
	err = engine.Run(port)
	if err != nil {
		panic(fmt.Sprintf("启动HTTP服务失败: %v", err))
	} else {
		fmt.Printf("HTTP服务已启动，监听端口: %s\n", port)
	}
}
