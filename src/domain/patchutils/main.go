package patchutils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/command"
	"sgridnext.com/src/logger"
)

type pathUtils struct{}

var T_PatchUtils = &pathUtils{}

// 初始化目录
func (p *pathUtils) InitDir(serverName string) error {
	cwd, _ := os.Getwd()
	dirs := []string{
		filepath.Join(cwd, constant.TAGET_CONF_DIR, serverName),
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

// 添加日志文件
func (p *pathUtils) AddConfigFile(serverName, configName, configContent string) error {
	cwd, _ := os.Getwd()
	configPath := filepath.Join(cwd, constant.TAGET_CONF_DIR, serverName, configName)
	return os.WriteFile(configPath, []byte(configContent), 0644)
}

// 获取配置文件列表
func (p *pathUtils) GetConfigFileList(serverName string) string {
	cwd, _ := os.Getwd()
	configPath := filepath.Join(cwd, constant.TAGET_CONF_DIR, serverName)
	configList, _ := os.ReadDir(configPath)
	configListStr := ""
	for _, config := range configList {
		configListStr += config.Name() + "\n"
	}
	return configListStr
}

// 获取配置文件内容
func (p *pathUtils) GetConfigFileContent(serverName, configName string) (string, error) {
	cwd, _ := os.Getwd()
	configPath := filepath.Join(cwd, constant.TAGET_CONF_DIR, serverName, configName)
	configContent, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}
	return string(configContent), nil
}

// 更改配置文件内容，将旧配置文件名进行替换成时间戳的形式
// 原始有一个 sgrid.yml 的配置文件 ，现在新的文件进来了，旧的配置文件名改成 sgrid_1234567890.yml
// 然后将新的配置文件内容写入到 sgrid.yml 文件中
func (p *pathUtils) UpdateConfigFileContent(serverName, configName, configContent string) error {
	cwd, _ := os.Getwd()
	configPath := filepath.Join(cwd, constant.TAGET_CONF_DIR, serverName, configName)
	// 如果旧配置文件存在，则重命名为带时间戳的备份文件
	if _, err := os.Stat(configPath); err == nil {
		timestamp := time.Now().Unix()
		ext := filepath.Ext(configName)
		backupName := configName[:len(configName)-len(ext)] + "_" + fmt.Sprintf("%d", timestamp) + ext
		backupPath := filepath.Join(cwd, constant.TAGET_CONF_DIR, serverName, backupName)
		if err := os.Rename(configPath, backupPath); err != nil {
			return fmt.Errorf("failed to backup config: %v", err)
		}
	}
	// 写入新的配置文件内容
	return os.WriteFile(configPath, []byte(configContent), 0644)
}

// 计算文件hash
func (p *pathUtils) CalcPackageHash(file os.File) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, &file); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %v", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// 重命名包
func (p *pathUtils) RenamePackage(oldPath string, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// SgridTestServer.tar.gz 改名成 SgridTestServer_1234567890.tar.gz
func (p *pathUtils) RenamePackageWithHash(oldPath string, hash string) error {
	ext := filepath.Ext(oldPath)
	newPath := oldPath[:len(oldPath)-len(ext)] + "_" + hash + ext
	return os.Rename(oldPath, newPath)
}

func (p *pathUtils) InitServer(serverInfo *ServerInfo) *command.Command {
	logger.PatchUtils.Infof("InitServer: %v", serverInfo)
	server_cmd := serverInfo.CreateCommand()
	server_cmd.AppendEnv([]string{
		fmt.Sprintf("%s=%s", constant.SGRID_TARGET_HOST, serverInfo.BindHost),
		fmt.Sprintf("%s=%s", constant.SGRID_TARGET_PORT, serverInfo.BindPort),
	})
	return server_cmd
}

func (p *pathUtils) StartServer(cmd *command.Command) (int, error) {
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	cmd.SetPid(cmd.GetCmd().Process.Pid)
	return cmd.GetCmd().Process.Pid, nil
}
