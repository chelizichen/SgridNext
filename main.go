package main

import (
	"github.com/gin-gonic/gin"
	"sgridnext.com/src/domain/config"
	"sgridnext.com/src/domain/service"
)

func main() {
	config.LoadConfig("./config.json")
	r := gin.Default()
	service.LoadRoutes(r)
	r.Run(":15872")
}
