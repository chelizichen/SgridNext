package main

import (
	"fmt"

	"github.com/gin-contrib/pprof"
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
	if conf.Get("pprof") == "enable" {
		pprof.Register(engine)
	}
	routes.LoadRoutes(engine)
	addr := fmt.Sprintf("%s:%s", conf.Get("httpHost"), conf.Get("httpPort"))
	err = engine.Run(addr)
	if err != nil {
		panic(fmt.Sprintf("启动HTTP服务失败: %v", err))
	} else {
		fmt.Printf("HTTP服务已启动，监听端口: %s\n", addr)
	}
}
