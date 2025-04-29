package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
)

func CreateServer(ctx *gin.Context) {
	var req struct {
		ServerName   string `json:"serverName"`
		GroupId      int    `json:"groupId"`
		ServerType   int    `json:"serverType"`
		Description  string `json:"description"`
		ExecFilePath string `json:"execFilePath"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
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
	}
	if _, err := mapper.T_Mapper.CreateServer(server); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建服务失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "创建服务成功"})
}

func CreatePackage(ctx *gin.Context) {

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

func DeployServer(ctx *gin.Context) {

}

func StopServer(ctx *gin.Context) {

}

func RestartServer(ctx *gin.Context) {

}

func GetServerNodesStatus(ctx *gin.Context) {

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
	res, err := mapper.T_Mapper.GetServerNodes(req.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}

func GetServerPackageList(ctx *gin.Context) {

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
