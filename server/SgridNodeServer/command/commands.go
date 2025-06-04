package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sgridnext.com/src/logger"
)

func CreateNodeCommand(serverName string, targetFile string) (*Command ,error){
	cmd := NewServerCommand(serverName)
	err := cmd.SetCommand("node", targetFile)
	return cmd,err
}

func CreateBinaryCommand(serverName string, targetFile string) (*Command ,error){
	cmd := NewServerCommand(serverName)
	err := cmd.SetCommand(targetFile)
	if err != nil{
		logger.App.Error("create command error:", err)
		return nil,err
	}
	err = AddPerm(targetFile)
	if err != nil{
		logger.App.Error("add perm error:", err)
		return nil,err
	}
	return cmd,err
}

func CreateJavaJarCommand(serverName string, targetDir string) (*Command ,error){
	logger.CMD.Infof("CreateJavaJarCommand | %s | targetDir %s ",serverName,targetDir)
	// 通过targetDir 去扫路径下的 jar文件
	cmd := NewServerCommand(serverName)
	// 在 目录下寻找 以 .jar 结尾的文件
	jarFile, err := FindFile(targetDir, ".jar")
	if err != nil {
		logger.CMD.Errorf("未找到jar文件: %v", err)
		return nil, err
	}
	if jarFile != "" {
		logger.CMD.Infof("找到jar文件: %s", jarFile)
		err = cmd.SetCommand("java", "-jar", jarFile)
		return cmd,err
	}
	warFile, err := FindFile(targetDir, ".war")
	if err!= nil {
		logger.CMD.Errorf("未找到war文件: %v", err)
		return nil, err
	}
	if warFile != "" {
		logger.CMD.Infof("找到war文件: %s", warFile)
		err = cmd.SetCommand("java", "-war", warFile)
		return cmd,err
	}
	return nil, fmt.Errorf("未找到.jar或.war后缀文件")
}

func FindFile(targetDir string, suffix string) (string, error) {
	var foundFile string
	logger.CMD.Infof("开始在目录下查找文件: %s", targetDir)
	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.CMD.Errorf("文件遍历错误: %v", err)
			return err
		}
		logger.CMD.Infof("正在遍历文件: %s", path)
		if !info.IsDir() && strings.HasSuffix(path, suffix) {
			foundFile = path
			return filepath.SkipDir
		}
		return nil
	})

	if foundFile == "" {
		return "", fmt.Errorf("未找到%s后缀文件", suffix)
	}
	return foundFile, err
}