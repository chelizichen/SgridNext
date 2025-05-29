package handlers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
	"sgridnext.com/src/patchutils"
	"sgridnext.com/src/proxy"
)

func DeployScripts(ctx *gin.Context) {
    // 一键部署脚本
    // 1. 上传文件
    // 2. 发布到节点
    // 3. 激活节点
    cwd, _ := os.Getwd()
	file, err := ctx.FormFile("file")
	serverName := ctx.PostForm("serverName")
	serverId, _ := strconv.Atoi(ctx.PostForm("serverId")) // 转成 int
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文件上传失败", "error": err.Error()})
		return
	}
	hash, err := patchutils.T_PatchUtils.CalcPackageHash(file)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "计算文件hash失败", "error": err.Error()})
		return
	}
	storePath := filepath.Join(cwd, constant.TARGET_PACKAGE_DIR, serverName, file.Filename)
	ctx.SaveUploadedFile(file, storePath)
	newFileName, err := patchutils.T_PatchUtils.RenamePackageWithHash(storePath, hash)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "重命名文件失败", "error": err.Error()})
		return
	}
	serverPackage := &entity.ServerPackage{
		ID:         0,
		ServerId:   serverId,
		Hash:       hash,
		CreateTime: constant.GetCurrentTime(),
		FileName:   newFileName,
		Commit:     "AutoDeploy",
	}
	id, err := mapper.T_Mapper.CreateServerPackage(serverPackage)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建服务包失败", "error": err.Error()})
		return
	}
    // 获取节点列表
    nodes, err := mapper.T_Mapper.GetServerNodes(serverId,0)
    if err != nil {
        ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取节点列表失败", "error": err.Error()})
        return
    }
    serverNodeIds := []int32{}
    for _, node := range nodes {
        serverNodeIds = append(serverNodeIds, int32(node.Id))
    }
    logger.Server.Infof("DeployScripts | CreateServerPackage | %d", id)
    // 发布到节点
    proxy.ProxyMap.FullDispatch(func(client *protocol.NodeServantClient) error {
        _, err := (*client).SyncServicePackage(context.Background(), &protocol.SyncReq{
            FileName: newFileName,
            Type: constant.FILE_TYPE_PACKAGE,
        })
        if err != nil {
            logger.Server.Errorf("DeployScripts | SyncServicePackage | %v", err)
            return err
        }
        return nil
    })
    proxy.ProxyMap.FullDispatch(func(client *protocol.NodeServantClient) error {
        _, err := (*client).ActivateServant(context.Background(), &protocol.ActivateReq{
            ServerId: int32(serverId),
            PackageId: int32(id),
            ServerNodeIds: serverNodeIds,
            Type: constant.ACITVATE_DEPLOY,
        })
        if err != nil {
            logger.Server.Errorf("DeployScripts | ActivateServant | %v", err)
            return err
        }
        return nil
    })
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "部署成功"})
}