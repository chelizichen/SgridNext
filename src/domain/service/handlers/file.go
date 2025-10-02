package handlers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
	"sgridnext.com/src/proxy"
)

func GetFile(ctx *gin.Context) {
	var req struct {
		ServerId int    `json:"serverId"`
		FileName string `json:"fileName"`
		Type     int    `json:"type"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	logger.App.Info("GetFile %v", req)
	serverInfo, err := mapper.T_Mapper.GetServerInfo(req.ServerId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务信息失败"})
		return
	}
	var file_path string
	cwd, _ := os.Getwd()
	if req.Type == constant.FILE_TYPE_PACKAGE {
		file_path = filepath.Join(cwd, constant.TARGET_PACKAGE_DIR, serverInfo.ServerName, req.FileName)
	}
	if req.Type == constant.FILE_TYPE_CONFIG {
		file_path = filepath.Join(cwd, constant.TARGET_CONF_DIR, serverInfo.ServerName, req.FileName)
	}
	logger.App.Infof("获取文件 %s", file_path)
	if _, err := os.Stat(file_path); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文件不存在"})
		return
	}
	ctx.File(file_path)
}

// 屏蔽 含有 _ 的历史文件
func GetConfigList(ctx *gin.Context) {
	var req struct {
		ServerId int `json:"serverId"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	logger.App.Info("GetConfigList %v", req)
	serverInfo, err := mapper.T_Mapper.GetServerInfo(req.ServerId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务信息失败"})
		return
	}
	cwd, _ := os.Getwd()
	config_dir := filepath.Join(cwd, constant.TARGET_CONF_DIR, serverInfo.ServerName)
	files, err := os.ReadDir(config_dir)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取配置文件列表失败"})
		return
	}
	var file_list []string
	for _, file := range files {
		if strings.Contains(file.Name(), "_") {
			continue
		}
		file_list = append(file_list, file.Name())
	}
	logger.App.Info("GetConfigList | file_list | %v", file_list)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取配置文件列表成功", "data": file_list})
}

func DownloadFile(ctx *gin.Context) {
	var req struct {
		ServerId int    `json:"serverId"`
		FileName string `json:"fileName"`
		Type     int    `json:"type"`
		Host     string `json:"host"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	logger.App.Info("DownloadFile %v", req)
		// 包就不用走 服务端了， 走本地就行了，主控直接下载
	if req.Type == constant.FILE_TYPE_PACKAGE {
		serverInfo, err := mapper.T_Mapper.GetServerInfo(req.ServerId)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务信息失败"})
			return
		}
		cwd, _ := os.Getwd()
		file_path := filepath.Join(cwd, constant.TARGET_PACKAGE_DIR, serverInfo.ServerName, req.FileName)
		ctx.File(file_path)
		return
	}else{
		proxy.ProxyMap.DispatchByHost(req.Host, func(client *protocol.NodeServantClient) error {
			rsp, err := (*client).DownloadFile(ctx, &protocol.DownloadFileRequest{
				ServerId: int32(req.ServerId),
				FileName: req.FileName,
				Type:     int32(req.Type),
			})
			if err != nil {
				logger.App.Errorf("下载文件失败 %s ", err.Error())
				return err
			}
			for {
				rsp, err := rsp.Recv()
				if err != nil {
					logger.App.Errorf("接受下载文件失败 %s ", err.Error())
					return err
				}
				if rsp.IsEnd {
					logger.App.Info("文件下载完成")
					return nil
				}
				ctx.Data(http.StatusOK, "application/octet-stream", rsp.Data)
			}
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "下载文件成功"})
}

func GetFileList(ctx *gin.Context) {
	var req struct {
		ServerId int    `json:"serverId"`
		Type     int    `json:"type"`
		Host     string `json:"host"`
		LogCategory  int    `json:"logCategory"` // 新增日志类型参数
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	logger.App.Info("GetFileList %v", req)
	
	// 如果是主控日志，直接在本地处理
	if req.Type == constant.FILE_TYPE_LOG && req.LogCategory == constant.LOG_TYPE_MASTER {
		cwd, _ := os.Getwd()
		masterLogDir := filepath.Join(cwd, "logs")
		
		// 如果目录不存在，创建它
		if _, err := os.Stat(masterLogDir); os.IsNotExist(err) {
			os.MkdirAll(masterLogDir, 0755)
		}
		
		files, err := os.ReadDir(masterLogDir)
		if err != nil {
			logger.App.Errorf("读取主控日志目录失败 %s", err.Error())
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "读取主控日志目录失败"})
			return
		}
		
		fileList := make([]string, 0)
		for _, file := range files {
			if !file.IsDir() {
				fileList = append(fileList, file.Name())
			}
		}
		
		ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取主控日志文件列表成功", "data": fileList})
		return
	}
	
	// 业务日志和节点日志通过RPC调用节点服务获取
	proxy.ProxyMap.DispatchByHost(req.Host, func(client *protocol.NodeServantClient) error {
		rsp, err := (*client).GetFileList(ctx, &protocol.GetFileListReq{
			ServerId: int32(req.ServerId),
			Type:     int32(req.Type),
			LogCategory:  int32(req.LogCategory), // 传递日志类型
		})
		if err != nil {
			logger.App.Errorf("获取文件列表失败 %s", err.Error())
			return err
		}
		ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取文件列表成功", "data": rsp.FileList})
		return nil
	})
}

// 同步服务包
func SyncUploadFile(ctx *gin.Context){
	var req struct {
		ServerId int    `json:"serverId"`
		PackageId int   `json:"packageId"`
		NodeIds   []int `json:"serverNodeIds"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	packageInfo, err := mapper.T_Mapper.GetServerPackageInfo(req.PackageId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务包信息失败"})
		return
	}
	for _, nodeId := range req.NodeIds {
		if nodeId == 0 {
			continue
		}
		nodeInfoArr, err := mapper.T_Mapper.GetServerNodes(req.ServerId,nodeId)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取节点信息失败"})
			return
		}
		if len(nodeInfoArr) == 0 {
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取节点信息失败"})
			return
		}
		err = proxy.ProxyMap.DispatchByHost(nodeInfoArr[0].Host, func(client *protocol.NodeServantClient) error {
				_, err := (*client).SyncServicePackage(context.Background(), &protocol.SyncReq{
					FileName: packageInfo.FileName,
					Type: constant.FILE_TYPE_PACKAGE,
				})
				if err != nil {
					logger.Server.Errorf("DeployScripts | SyncServicePackage | %v", err)
					return err
				}
				return nil
			})
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "同步服务包失败: " + err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "同步服务包成功"})
}
