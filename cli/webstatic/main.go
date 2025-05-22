package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// go run main.go -static_public_path="./dist/" -relative_path="/sgridtest"
func main() {
	// 前端设置的公共路径
	var cli_static_public_path string
	// 后台的路径
	var cli_relative_path string

	flag.StringVar(&cli_static_public_path, "static_public_path", "./dist/", "static public path")
	flag.StringVar(&cli_relative_path, "relative_path", "/web", "relative path")

	flag.Parse()
	engine := gin.Default()
	cwd, _ := os.Getwd()
	root_path := filepath.Join(cwd, cli_static_public_path)
	engine.Static(cli_relative_path, root_path)
	fmt.Println("SGRID_WEB_STATIC_SERVER STARTED ", "relative_path: ", cli_relative_path, "root_path: ", root_path)
	port := os.Getenv("SGRID_TARGET_PORT")
	if port == "" {
		fmt.Println("SGRID_TARGET_PORT is empty")
		port = "8080"
	}
	host := os.Getenv("SGRID_TARGET_HOST")
	if host == "" {
		fmt.Println("SGRID_TARGET_HOST is empty")
		host = "0.0.0.0"
	}

	bind_addr := host + ":" + port
	engine.Run(bind_addr)
}
