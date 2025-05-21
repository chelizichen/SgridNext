package main

// GOOS=linux GOARCH=amd64 go build -o $ServerName
import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"sgridnext.com/src/logger"
)

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func main() {
	port := os.Getenv("SGRID_TARGET_PORT")
	logger.App.Info("SGRID_TARGET_PORT: ", port)
	if port == "" {
		logger.App.Info("SGRID_TARGET_PORT is empty")
	}
	port = "10010"
	host := os.Getenv("SGRID_TARGET_HOST")
	logger.App.Info("SGRID_TARGET_HOST: ", port)
	if host == "" {
		logger.App.Info("SGRID_TARGET_HOST is empty")
	}
	host = "0.0.0.0"

	engine := gin.Default()
	engine.GET("/fib", func(ctx *gin.Context) {
		n, _ := strconv.Atoi(ctx.Query("n"))
		go fib(n)
		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})
	bind_addr := host + ":" + port
	engine.Run(bind_addr)
}
