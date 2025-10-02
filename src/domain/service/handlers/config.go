package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
	"sgridnext.com/src/patchutils"
)

// 获取主控配置文件
func GetMainConfig(ctx *gin.Context) {
	cwd, err := os.Getwd()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "获取当前工作目录失败"})
		return
	}

	configPath := filepath.Join(cwd, "config.json")

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "msg": "配置文件不存在"})
		return
	}

	// 读取配置文件
	configData, err := os.ReadFile(configPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "读取配置文件失败"})
		return
	}

	// 解析JSON
	var config map[string]interface{}
	if err := json.Unmarshal(configData, &config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "解析配置文件失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"configPath": configPath,
			"config":     config,
		},
	})
}

// 更新主控配置文件
func UpdateMainConfig(ctx *gin.Context) {
	var req struct {
		Config map[string]interface{} `json:"config"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	
	// 调试信息
	logger.Server.Infof("收到配置更新请求: %+v", req.Config)

	cwd, err := os.Getwd()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "获取当前工作目录失败"})
		return
	}

	configPath := filepath.Join(cwd, "config.json")

	// 创建备份
	backupPath := configPath + ".backup"
	if err := copyFile(configPath, backupPath); err != nil {
		logger.Server.Warnf("创建配置备份失败: %v", err)
	}

	// 写入新配置
	configData, err := json.MarshalIndent(req.Config, "", "    ")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "序列化配置失败"})
		return
	}

	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "写入配置文件失败"})
		return
	}

	logger.Server.Infof("配置文件已更新: %s", configPath)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "配置更新成功"})
}

// 获取配置项
func GetConfigItem(ctx *gin.Context) {
	var req struct {
		Key string `json:"key"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "获取当前工作目录失败"})
		return
	}

	configPath := filepath.Join(cwd, "config.json")

	// 读取配置文件
	configData, err := os.ReadFile(configPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "读取配置文件失败"})
		return
	}

	// 解析JSON
	var config map[string]interface{}
	if err := json.Unmarshal(configData, &config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "解析配置文件失败"})
		return
	}

	// 获取配置项
	value, exists := config[req.Key]
	if !exists {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "配置项不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"key":   req.Key,
			"value": value,
		},
	})
}

// 设置配置项
func SetConfigItem(ctx *gin.Context) {
	var req struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "获取当前工作目录失败"})
		return
	}

	configPath := filepath.Join(cwd, "config.json")

	// 读取现有配置
	var config map[string]interface{}
	if _, err := os.Stat(configPath); err == nil {
		configData, err := os.ReadFile(configPath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "读取配置文件失败"})
			return
		}

		if err := json.Unmarshal(configData, &config); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "解析配置文件失败"})
			return
		}
	} else {
		config = make(map[string]interface{})
	}

	// 更新配置项
	config[req.Key] = req.Value

	// 写入配置文件
	configData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "序列化配置失败"})
		return
	}

	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "写入配置文件失败"})
		return
	}

	logger.Server.Infof("配置项已更新: %s = %v", req.Key, req.Value)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "配置项更新成功"})
}

// 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

func UpsertConfig(ctx *gin.Context) {
	var req struct {
		FileName    string `json:"fileName"`
		FileContent string `json:"fileContent"`
		ServerId    int    `json:"serverId"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	serverInfo, err := mapper.T_Mapper.GetServerInfo(req.ServerId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败"})
		return
	}
	serverName := serverInfo.ServerName
	if serverName == "" {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "服务器名称不能为空"})
		return
	}
	// 创建
	configPath, err := patchutils.T_PatchUtils.UpdateConfigFileContent(serverName, req.FileName, req.FileContent)
	// 更新
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建配置文件失败"})
		return
	}
	ctx.JSON(200, gin.H{"success": true, "msg": "创建配置文件成功|" + configPath})
}

func DeleteConfig(ctx *gin.Context) {

}

func GetServerConfigList(ctx *gin.Context) {
	var req struct {
		ServerId int `json:"serverId"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	serverInfo, err := mapper.T_Mapper.GetServerInfo(req.ServerId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败"})
		return
	}
	serverName := serverInfo.ServerName
	if serverName == "" {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "服务器名称不能为空"})
		return
	}
	configList := patchutils.T_PatchUtils.GetConfigFileList(serverName)
	ctx.JSON(200, gin.H{"success": true, "msg": "获取配置文件列表成功", "data": configList})
}

func BackupConfig(ctx *gin.Context) {
	var req struct {
		ServerId   int    `json:"serverId"`
		ConfigName string `json:"configName"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	serverInfo, err := mapper.T_Mapper.GetServerInfo(req.ServerId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败"})
		return
	}
	serverName := serverInfo.ServerName
	if serverName == "" {
		ctx.JSON(400, gin.H{"success": false, "msg": "服务器名称不能为空"})
		return
	}
	err = patchutils.T_PatchUtils.BackConfigFile(serverName, req.ConfigName, req.ConfigName+"_"+time.Now().Format("20060102150405")+".json")
	if err != nil {
		ctx.JSON(400, gin.H{"success": false, "msg": "备份配置文件失败"})
		return
	}
	ctx.JSON(200, gin.H{"success": true, "msg": "备份配置文件成功"})
}

func GetConfigContent(ctx *gin.Context) {
	var req struct {
		ServerId   int    `json:"serverId"`
		ConfigName string `json:"configName"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	serverInfo, err := mapper.T_Mapper.GetServerInfo(req.ServerId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败"})
		return
	}
	serverName := serverInfo.ServerName
	if serverName == "" {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "服务器名称不能为空"})
		return
	}
	configContent, err := patchutils.T_PatchUtils.GetConfigFileContent(serverName, req.ConfigName)
	if err != nil {
		logger.Config.Errorf("获取配置文件内容失败:%v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取配置文件内容失败"})
		return
	}
	ctx.JSON(200, gin.H{"success": true, "msg": "获取配置文件内容成功", "data": configContent})
}
