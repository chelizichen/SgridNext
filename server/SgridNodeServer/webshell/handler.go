package webshell

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sgridnext.com/src/config"
	"sgridnext.com/src/logger"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域，生产环境应该检查来源
	},
}

// HandleWebSocket 处理 websocket 连接
func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.App.Errorf("WebSocket 升级失败: %v", err)
		return
	}
	defer conn.Close()

	// 根据操作系统选择 shell
	var shell string
	var shellArgs []string
	if runtime.GOOS == "windows" {
		// Windows 不支持 pty，使用普通管道
		shell = "cmd.exe"
	} else if config.Conf.GetOs() == "android" {
		shell = "/bin/sh"
	} else {
		shell = "/bin/bash"
	}

	// 创建命令
	cmd := exec.Command(shell, shellArgs...)
	// 继承环境变量并添加必要的终端设置
	cmd.Env = append(os.Environ(),
		"TERM=xterm-256color",
		"LANG=en_US.UTF-8",
	)

	// Windows 使用普通管道，其他系统使用 pty
	if runtime.GOOS == "windows" {
		handleWindowsShell(conn, cmd)
	} else {
		handleUnixShell(conn, cmd)
	}
}

// handleUnixShell 处理 Unix 系统的 shell（使用 pty）
func handleUnixShell(conn *websocket.Conn, cmd *exec.Cmd) {
	// 创建伪终端
	ptmx, err := pty.Start(cmd)
	if err != nil {
		logger.App.Errorf("启动 pty 失败: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("启动 shell 失败: "+err.Error()+"\r\n"))
		return
	}
	defer func() {
		ptmx.Close()
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		cmd.Wait()
	}()

	// 设置终端大小
	if err := pty.Setsize(ptmx, &pty.Winsize{
		Rows: 24,
		Cols: 80,
	}); err != nil {
		logger.App.Warnf("设置终端大小失败: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// 从 websocket 读取数据并写入 pty
	go func() {
		defer wg.Done()
		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.App.Errorf("WebSocket 读取错误: %v", err)
				}
				return
			}
			// 处理文本和二进制消息
			if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
				if _, err := ptmx.Write(data); err != nil {
					logger.App.Errorf("写入 pty 失败: %v", err)
					return
				}
			}
		}
	}()

	// 从 pty 读取数据并发送到 websocket
	go func() {
		defer wg.Done()
		buffer := make([]byte, 4096) // 增大缓冲区
		for {
			n, err := ptmx.Read(buffer)
			if n > 0 {
				// 发送二进制消息
				if err := conn.WriteMessage(websocket.BinaryMessage, buffer[:n]); err != nil {
					logger.App.Errorf("WebSocket 写入错误: %v", err)
					return
				}
			}
			if err != nil {
				if err != io.EOF {
					logger.App.Errorf("读取 pty 错误: %v", err)
				}
				return
			}
		}
	}()

	// 等待命令结束
	go func() {
		cmd.Wait()
		wg.Wait()
		conn.Close()
	}()

	// 等待所有 goroutine 完成
	wg.Wait()
}

// handleWindowsShell 处理 Windows 系统的 shell（使用普通管道）
func handleWindowsShell(conn *websocket.Conn, cmd *exec.Cmd) {
	// 设置标准输入输出
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.App.Errorf("创建标准输入管道失败: %v", err)
		return
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.App.Errorf("创建标准输出管道失败: %v", err)
		return
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.App.Errorf("创建标准错误管道失败: %v", err)
		return
	}
	defer stderr.Close()

	// 启动命令
	if err := cmd.Start(); err != nil {
		logger.App.Errorf("启动命令失败: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("启动 shell 失败: "+err.Error()+"\r\n"))
		return
	}

	var wg sync.WaitGroup
	wg.Add(3)

	// 从 websocket 读取数据并写入命令的标准输入
	go func() {
		defer wg.Done()
		defer stdin.Close()
		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.App.Errorf("WebSocket 读取错误: %v", err)
				}
				return
			}
			// 处理文本和二进制消息
			if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
				if _, err := stdin.Write(data); err != nil {
					logger.App.Errorf("写入标准输入失败: %v", err)
					return
				}
			}
		}
	}()

	// 从命令的标准输出读取数据并发送到 websocket
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		for {
			n, err := stdout.Read(buffer)
			if n > 0 {
				if err := conn.WriteMessage(websocket.BinaryMessage, buffer[:n]); err != nil {
					logger.App.Errorf("WebSocket 写入错误: %v", err)
					return
				}
			}
			if err != nil {
				if err != io.EOF {
					logger.App.Errorf("读取标准输出错误: %v", err)
				}
				return
			}
		}
	}()

	// 从命令的标准错误读取数据并发送到 websocket
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		for {
			n, err := stderr.Read(buffer)
			if n > 0 {
				if err := conn.WriteMessage(websocket.BinaryMessage, buffer[:n]); err != nil {
					logger.App.Errorf("WebSocket 写入错误: %v", err)
					return
				}
			}
			if err != nil {
				if err != io.EOF {
					logger.App.Errorf("读取标准错误错误: %v", err)
				}
				return
			}
		}
	}()

	// 等待命令结束
	go func() {
		cmd.Wait()
		wg.Wait()
		conn.Close()
	}()

	// 等待所有 goroutine 完成
	wg.Wait()
}
