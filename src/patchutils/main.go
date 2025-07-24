package patchutils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/c4milo/unpackit"
	"sgridnext.com/server/SgridNodeServer/command"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/logger"
)

type pathUtils struct{}

var T_PatchUtils = &pathUtils{}

// 初始化目录
func (p *pathUtils) InitDir(serverName string) error {
	cwd, _ := os.Getwd()
	dirs := []string{
		filepath.Join(cwd, constant.TARGET_CONF_DIR, serverName),
		filepath.Join(cwd, constant.TARGET_LOG_DIR, serverName),
		filepath.Join(cwd, constant.TARGET_PACKAGE_DIR, serverName),
		filepath.Join(cwd, constant.TARGET_SERVANT_DIR, serverName),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// // 添加配置文件
// func (p *pathUtils) AddConfigFile(serverName, configName, configContent string) error {
// 	cwd, _ := os.Getwd()
// 	configPath := filepath.Join(cwd, constant.TARGET_CONF_DIR, serverName, configName)
// 	return os.WriteFile(configPath, []byte(configContent), 0644)
// }

// 获取配置文件列表
func (p *pathUtils) GetConfigFileList(serverName string) []string {
	cwd, _ := os.Getwd()
	configPath := filepath.Join(cwd, constant.TARGET_CONF_DIR, serverName)
	configList, _ := os.ReadDir(configPath)
	var configListRsp []string
	for _, config := range configList {
		configListRsp = append(configListRsp, config.Name())
	}
	return configListRsp
}

// 获取配置文件内容
func (p *pathUtils) GetConfigFileContent(serverName, configName string) (string, error) {
	cwd, _ := os.Getwd()
	configPath := filepath.Join(cwd, constant.TARGET_CONF_DIR, serverName, configName)
	configContent, err := os.ReadFile(configPath)
	logger.Config.Infof("configPath: %s", configPath)
	if err != nil {
		return "", err
	}
	return string(configContent), nil
}

// 新配置文件进来，先创建时间备份的配置文件
// 再删除旧的配置文件
// 最后更新新的配置文件
func (p *pathUtils) UpdateConfigFileContent(serverName, configName, configContent string) (string, error) {
	cwd, _ := os.Getwd()
	configPath := filepath.Join(cwd, constant.TARGET_CONF_DIR, serverName, configName)
	logger.Config.Infof("configPath: %s", configPath)
	// 备份配置文件
	timestamp := time.Now().Unix()
	ext := filepath.Ext(configName)
	backupName := configName[:len(configName)-len(ext)] + "_" + fmt.Sprintf("%d", timestamp) + ext
	backupPath := filepath.Join(cwd, constant.TARGET_CONF_DIR, serverName, backupName)
	logger.Config.Infof("backupPath: ", backupPath)
	err := os.WriteFile(backupPath, []byte(configContent), 0644)
	if err != nil {
		logger.Config.Errorf("failed to backup config: %v", err)
		return "", err
	}
	err = os.Remove(configPath)
	if err != nil {
		logger.Config.Errorf("failed to remove config: %v", err)
		// return err
	}
	return backupPath, os.WriteFile(configPath, []byte(configContent), 0644)
}

// 计算文件hash
func (p *pathUtils) CalcPackageHash(file *multipart.FileHeader) (string, error) {
	h := sha256.New()
	// 打开文件以获取 io.Reader
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	if _, err := io.Copy(h, src); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %v", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (p *pathUtils) CalcPackageHashFromReader(reader io.Reader) (string, error) {
	hash := sha256.New()
	_, err := io.Copy(hash, reader)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// 重命名包
func (p *pathUtils) RenamePackage(oldPath string, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// SgridTestServer.tar.gz 改名成 SgridTestServer_1234567890.tar.gz
func (p *pathUtils) RenamePackageWithHash(oldPath string, hash string) (string, error) {
	ext := filepath.Ext(oldPath)
	newPath := oldPath[:len(oldPath)-len(ext)] + "_" + hash + ext
	err := os.Rename(oldPath, newPath)
	newFileName := filepath.Base(newPath)
	return newFileName, err
}

// func (p *pathUtils) InitServer(serverInfo *ServerInfo) (*command.Command,error) {
// 	logger.PatchUtils.Infof("InitServer: %v", serverInfo)
// 	server_cmd,err := serverInfo.CreateCommand()
// 	server_cmd.AppendEnv([]string{
// 		fmt.Sprintf("%s=%s", constant.SGRID_TARGET_HOST, serverInfo.BindHost),
// 		fmt.Sprintf("%s=%s", constant.SGRID_TARGET_PORT, serverInfo.BindPort),
// 	})
// 	return server_cmd,err
// }

func (p *pathUtils) StartServer(cmd *command.Command) (int, error) {
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	return cmd.GetCmd().Process.Pid, nil
}

func (p *pathUtils) Tar2Dest(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		fmt.Println("Open Error", err.Error())
		return err
	}
	defer file.Close()
	err = unpackit.Unpack(file, dest)
	if err != nil {
		fmt.Println("Unpackit Error", err.Error())
		return err
	}
	return nil
}

func (p *pathUtils) Contains(nodes []int, node int) bool {
	for _, n := range nodes {
		if n == node {
			return true
		}
	}
	return false
}

// 回退
// args | serverName | configName | newConfigName
// TestServer | sgrid.yml | sgrid_1234567890.yml
func (p *pathUtils) BackConfigFile(serverName, originConfigName, newConfigName string) error {
	cwd, _ := os.Getwd()
	newConfig := filepath.Join(cwd, constant.TARGET_CONF_DIR, serverName, newConfigName)
	originPath := filepath.Join(cwd, constant.TARGET_CONF_DIR, serverName, originConfigName)
	if _, err := os.Stat(originPath); err == nil {
		if err := os.Remove(originPath); err != nil {
			return fmt.Errorf("failed to remove originPath config: %v", err)
		}
	}
	bytes, err := os.ReadFile(newConfig)
	if err != nil {
		logger.Config.Errorf("failed to backup config: %v", err)
		return err
	}
	return os.WriteFile(originPath, bytes, 0644)
}
