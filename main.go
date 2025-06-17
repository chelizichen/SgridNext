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
	ormDb,err := db.InitDB(conf.Get("db"),conf.Get("dbtype"))
	if err != nil {
		panic(err)
	}
	mapper.LoadMapper(ormDb)
	proxy.LoadProxy()
	// 初始化路由
	engine := gin.Default()
	routes.LoadRoutes(engine)
	port := fmt.Sprintf(":%s",conf.Get("httpPort"))
	engine.Run(port)
}
