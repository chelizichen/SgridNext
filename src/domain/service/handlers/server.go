package handlers

import (
	"context"
	"fmt"
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

func CreateServer(ctx *gin.Context) {
	var req struct {
		ServerName   string `json:"serverName"`
		GroupId      int    `json:"groupId"`
		ServerType   int    `json:"serverType"`
		Description  string `json:"description"`
		ExecFilePath string `json:"execFilePath"`
		LogPath 	 string `json:"logPath"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	err := patchutils.T_PatchUtils.InitDir(req.ServerName)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建服务失败", "error": err.Error()})
		return
	}
	server := &entity.Server{
		ID:           0,
		ServerName:   req.ServerName,
		ServerType:   req.ServerType,
		Status:       constant.COMM_STATUS_ONLINE,
		ExecFilePath: req.ExecFilePath,
		CreateTime:   constant.GetCurrentTime(),
		GroupId:      req.GroupId,
		Description:  req.Description,
		LogPath:      req.LogPath,
	}

	if _, err := mapper.T_Mapper.CreateServer(server); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建服务失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "创建服务成功"})
}

func CreatePackage(ctx *gin.Context) {
	cwd, _ := os.Getwd()
	file, err := ctx.FormFile("file")
	serverName := ctx.PostForm("serverName")
	commit := ctx.PostForm("commit")
	serverId, _ := strconv.Atoi(ctx.PostForm("serverId")) // 转成 int
	logger.Server.Infof("serverName | %s | commit | %s", serverName, commit)
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
		Commit:     commit,
		FileName:   newFileName,
	}
	rsp, err := mapper.T_Mapper.CreateServerPackage(serverPackage)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建服务包失败", "error": err.Error()})
		return
	}
	mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
		ServerName: serverName,
		ServerId:   serverId,
		TYPE:       entity.TYPE_PATCH,
		Content:    fmt.Sprintf("已部署服务包 %s | 版本号 %d", serverName, rsp),
	})

	go func ()  {
		// 全量调用
		proxy.ProxyMap.FullDispatch(func(nsc *protocol.NodeServantClient) error {
			_,err := (*nsc).SyncServicePackage(context.Background(),&protocol.SyncReq{
				FileName: newFileName,
				Type: constant.FILE_TYPE_PACKAGE,
				ServerId: int32(serverId),
			})
			if err  != nil{
				logger.RPC.Infof("调用失败 | SyncServicePackage | err | %s",err.Error())
			}
			return err
		})	
	}()

	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "创建服务包成功", "hash": hash})
}

func GetServerPackageList(ctx *gin.Context) {
	var req struct {
		ServerId int `json:"id"`
		Offset int `json:"offset"`
		Size int `json:"size"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	packages, total, err := mapper.T_Mapper.GetServerPackageList(req.ServerId, req.Offset, req.Size)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务包列表失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": packages,"total":total})
}

func CreateServerNode(ctx *gin.Context) {
	var req []struct {
		NodeId   int `json:"node_id"`
		PatchId  int `json:"patch_id"`
		Port     int `json:"port"`
		ServerId int `json:"server_id"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	var reqs []*entity.ServerNode
	for _, r := range req {
		reqs = append(reqs, &entity.ServerNode{
			NodeId:           r.NodeId,
			PatchId:          r.PatchId,
			Port:             r.Port,
			ServerId:         r.ServerId,
			CreateTime:       constant.GetCurrentTime(),
			ServerNodeStatus: constant.COMM_STATUS_ONLINE,
			ID:               0,
		})
	}
	if err := mapper.T_Mapper.CreateServerNode(reqs); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建服务节点失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "创建服务节点成功"})
}

func CreateGroup(ctx *gin.Context) {
	var req struct {
		GroupName        string `json:"groupName"`
		GroupEnglishName string `json:"groupEnglishName"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	groups, err := mapper.T_Mapper.GetGroupList()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务组列表失败", "error": err.Error()})
		return
	}
	for _, group := range groups {
		if group.Name == req.GroupName {
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "服务组已存在"})
			return
		}
		if group.EngLishName == req.GroupEnglishName {
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "服务组英文名称已存在"})
			return
		}
	}
	createTime := constant.GetCurrentTime()
	group := &entity.ServerGroup{
		Name:        req.GroupName,
		EngLishName: req.GroupEnglishName,
		Status:      constant.COMM_STATUS_ONLINE,
		ID:          0,
		CreateTime:  createTime,
	}
	logger.Server.Info("创建服务组：", group)
	if _, err := mapper.T_Mapper.CreateGroup(group); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建服务组失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "创建服务组成功"})
}

type RUN_SERVER_REQ struct {
	ServerId  int   `json:"serverId"`
	PackageId int   `json:"packageId"`
	NodeIds   []int `json:"serverNodeIds"`
}

func DeployServer(ctx *gin.Context) {
	var req RUN_SERVER_REQ
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	var callRsp []*protocol.BasicRes
	proxy.ProxyMap.FullDispatch(func(client *protocol.NodeServantClient) error {
		rsp, err := (*client).ActivateServant(context.Background(), &protocol.ActivateReq{
			ServerId:      int32(req.ServerId),
			ServerNodeIds: constant.ConvertToInt32Slice(req.NodeIds),
			Type:          constant.ACITVATE_DEPLOY,
			PackageId: 		int32(req.PackageId),
		})
		callRsp = append(callRsp, rsp)
		return err
	})
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "部署服务成功", "data": callRsp})
}

func RestartServer(ctx *gin.Context) {
	var req RUN_SERVER_REQ
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	var callRsp []*protocol.BasicRes
	proxy.ProxyMap.FullDispatch(func(client *protocol.NodeServantClient) error {
		rsp, err := (*client).ActivateServant(context.Background(), &protocol.ActivateReq{
			ServerId:      int32(req.ServerId),
			ServerNodeIds: constant.ConvertToInt32Slice(req.NodeIds),
			Type:          constant.ACTIVATE_RESTART,
		})
		callRsp = append(callRsp, rsp)
		return err
	})
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "部署服务成功", "data": callRsp})
}

func StopServer(ctx *gin.Context) {
	var req struct {
		ServerId int   `json:"serverId"`
		NodeIds  []int `json:"nodeIds"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	var callRsp []*protocol.BasicRes
	proxy.ProxyMap.FullDispatch(func(client *protocol.NodeServantClient) error {
		rsp, err := (*client).DeactivateServant(context.Background(), &protocol.ActivateReq{
			ServerId:      int32(req.ServerId),
			ServerNodeIds: constant.ConvertToInt32Slice(req.NodeIds),
			Type:          constant.ACTIVATE_STOP,
		})
		callRsp = append(callRsp, rsp)
		return err
	})
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "部署服务成功", "data": callRsp})
}

func GetServerNodesLog(ctx *gin.Context) {

}

func GetServerNodes(ctx *gin.Context) {
	var req struct {
		Id int `json:"id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	res, err := mapper.T_Mapper.GetServerNodes(req.Id, 0)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}


func UpdateServerNode(ctx *gin.Context){
	var req struct {
		Ids         []int    `json:"ids"`
		ServerRunType  int    `json:"server_run_type"`
		AdditionalArgs string `json:"additional_args"`
		ViewPage       string `json:"view_page"`
 	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	for _, id := range req.Ids {
		err := mapper.T_Mapper.UpdateServerNode(entity.ServerNode{
			ID: id,
			ServerRunType: req.ServerRunType,
			AdditionalArgs: req.AdditionalArgs,
			ViewPage: req.ViewPage,
		})
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "更新失败"})
			return 
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "更新成功"})
}

func DeleteServerNode(ctx *gin.Context) { 
	var req struct {
		Ids         []int    `json:"ids"`
 	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	 err := mapper.T_Mapper.DeleteServerNode(req.Ids)
	 if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "删除成功"})
}

func GetServerList(ctx *gin.Context) {
	servers, err := mapper.T_Mapper.GetServerListWithGroup()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "获取服务器列表失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": servers})
}

func GetGroupList(ctx *gin.Context) {
	groups, err := mapper.T_Mapper.GetGroupList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "获取服务组列表失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": groups})
}

func GetServerInfo(ctx *gin.Context) {
	var req struct {
		Id int `json:"id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	res, err := mapper.T_Mapper.GetServerInfo(req.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}


func GetLog(ctx *gin.Context){
	var req struct {
		ServerName string `json:"serverName"`
		ServerId int `json:"serverId"`
		NodeId int `json:"nodeId"`
		Len int `json:"len"`
		Keyword string `json:"keyword"`
		Host string `json:"host"`
		LogType int `json:"logType"`
		FileName string `json:"fileName"`
	}
	if err := ctx.ShouldBindJSON(&req); err!= nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	proxy.ProxyMap.DispatchByHost(req.Host, func(client *protocol.NodeServantClient) error {
		rsp, err := (*client).GetLog(context.Background(), &protocol.GetLogReq{
			ServerName: req.ServerName,
			ServerId: int32(req.ServerId),
			Len: int32(req.Len),
			Keyword: req.Keyword,
			LogType:int32(req.LogType),
			FileName: req.FileName,
		})
		if err != nil{
			logger.RPC.Infof("调用失败 | QueryLog | err | %s",err.Error())
			return err
		}
		ctx.JSON(http.StatusOK, gin.H{"success": true, "data": rsp.Data})
		return nil
	})
}