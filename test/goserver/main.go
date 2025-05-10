package main

// GOOS=linux GOARCH=amd64 go build -o $ServerName
import (
	"os"

	"github.com/gin-gonic/gin"
	"sgridnext.com/src/logger"
)

func main() {
	port := os.Getenv("SGRID_TARGET_PORT")
	logger.App.Info("SGRID_TARGET_PORT: ", port)
	if port == "" {
		logger.App.Info("SGRID_TARGET_PORT is empty")
		panic("SGRID_TARGET_PORT is empty")
	}
	port = "10010"
	host := os.Getenv("SGRID_TARGET_HOST")
	logger.App.Info("SGRID_TARGET_HOST: ", port)
	if host == "" {
		logger.App.Info("SGRID_TARGET_HOST is empty")
		panic("SGRID_TARGET_HOST is empty")
	}
	engine := gin.Default()
	engine.GET("/greet", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hello",
		})
	})
	bind_addr := host + ":" + port
	engine.Run(bind_addr)
}
