package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
	"sgridnext.com/src/patchutils"
)

func UpsertConfig(ctx *gin.Context) {
    var req struct{
        FileName string `json:"fileName"`
        FileContent string `json:"fileContent"`
        ServerId int `json:"serverId"`
    }
    if err := ctx.ShouldBindJSON(&req); err!= nil {
        ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
        return
    }

    serverInfo,err := mapper.T_Mapper.GetServerInfo(req.ServerId)
    if err!= nil {
        ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败"})
        return
    }
    serverName := serverInfo.ServerName
    if serverName == "" {
        ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "服务器名称不能为空"})
        return
    }
    // 创建
    configPath,err := patchutils.T_PatchUtils.UpdateConfigFileContent(serverName, req.FileName,req.FileContent)
    // 更新
    if err != nil {
        ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建配置文件失败"})
        return
    }
    ctx.JSON(200, gin.H{"success": true, "msg": "创建配置文件成功|"+configPath})
}

func DeleteConfig(ctx *gin.Context) {

}

func GetServerConfigList(ctx *gin.Context) {
    var req struct{
        ServerId int `json:"serverId"`
    }
    if err := ctx.ShouldBindJSON(&req); err!= nil {
        ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
        return
    }
    serverInfo,err := mapper.T_Mapper.GetServerInfo(req.ServerId)
    if err!= nil {
        ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败"})
        return
    }
    serverName := serverInfo.ServerName
    if serverName == "" {
        ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "服务器名称不能为空"})
        return
    }
    configList := patchutils.T_PatchUtils.GetConfigFileList(serverName)
    ctx.JSON(200, gin.H{"success": true, "msg": "获取配置文件列表成功","data":configList})
}

func BackupConfig(ctx *gin.Context) {
    var req struct{
        ServerId int `json:"serverId"`
        ConfigName string `json:"configName"`
    }
    if err := ctx.ShouldBindJSON(&req); err!= nil {
        ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
        return
    }

    serverInfo,err := mapper.T_Mapper.GetServerInfo(req.ServerId)
    if err!= nil {
        ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败"})
        return
    }
    serverName := serverInfo.ServerName
    if serverName == "" {
        ctx.JSON(400, gin.H{"success": false, "msg": "服务器名称不能为空"})
        return
    }
    err = patchutils.T_PatchUtils.BackConfigFile(serverName, req.ConfigName, req.ConfigName+"_"+time.Now().Format("20060102150405")+".json")
    if err!= nil {
        ctx.JSON(400, gin.H{"success": false, "msg": "备份配置文件失败"})
        return
    }
    ctx.JSON(200, gin.H{"success": true, "msg": "备份配置文件成功"})
}

func GetConfigContent(ctx *gin.Context) {
    var req struct{
        ServerId int `json:"serverId"`
        ConfigName string `json:"configName"`
    }
    if err := ctx.ShouldBindJSON(&req); err!= nil {
        ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
        return
    }
    serverInfo,err := mapper.T_Mapper.GetServerInfo(req.ServerId)
    if err!= nil {
        ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败"})
        return
    }
    serverName := serverInfo.ServerName
    if serverName == "" {
        ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "服务器名称不能为空"})
        return
    }
    configContent,err := patchutils.T_PatchUtils.GetConfigFileContent(serverName, req.ConfigName)
    if err!= nil {
        logger.Config.Errorf("获取配置文件内容失败:%v", err)
        ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取配置文件内容失败"})
        return
    }
    ctx.JSON(200, gin.H{"success": true, "msg": "获取配置文件内容成功","data":configContent})
}
