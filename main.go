package main

import (
	"github.com/gin-gonic/gin"
	"sgridnext.com/src/domain/config"
	"sgridnext.com/src/domain/db"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/domain/service/routes"
)

func main() {
	conf := config.LoadConfig("./config.json")
	ormDb,err := db.InitDB(conf.Get("db"))
	if err != nil {
		panic(err)
	}
	mapper.LoadMapper(ormDb)

	// 初始化路由
	engine := gin.Default()
	routes.LoadRoutes(engine)

	engine.Run(":15872")
}
