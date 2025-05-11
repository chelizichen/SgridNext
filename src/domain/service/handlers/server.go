package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/command"
	"sgridnext.com/src/domain/patchutils"
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
	rsp, err := mapper.T_Mapper.CreateServerPackage(serverPackage);
	if  err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建服务包失败", "error": err.Error()})
		return
	}
	mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
		ServerName: serverName,
		ServerId:   serverId,
		TYPE:       entity.TYPE_PATCH,
		Content:    fmt.Sprintf("已部署服务包 %s | 版本号 %d", serverName, rsp),
	})
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "创建服务包成功", "hash": hash})
}

func GetServerPackageList(ctx *gin.Context) {
	var req struct {
		ServerId int `json:"id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	packages, err := mapper.T_Mapper.GetServerPackageList(req.ServerId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务包列表失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": packages})
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

type RUN_SERVER_REQ struct{
	ServerId  int   `json:"serverId"`
	PackageId int   `json:"packageId"`
	NodeIds   []int `json:"serverNodeIds"`
}

func RunServer(ctx *gin.Context,req RUN_SERVER_REQ, needDeploy bool) {
	cwd, _ := os.Getwd()
	logger.Server.Infof("DeployServer | req | %v", req)
	serverInfo, err := mapper.T_Mapper.GetServerInfo(req.ServerId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败", "error": err.Error()})
		return
	}
	logger.Server.Infof("DeployServer | serverInfo | %v", serverInfo)
	execPath := serverInfo.ExecFilePath
	serverName := serverInfo.ServerName
	serverType := serverInfo.ServerType
	// 不是重启，需要解压文件到目标目录，然后再启动execPath
	if needDeploy {
		packageInfo, err := mapper.T_Mapper.GetServerPackageInfo(req.PackageId)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器包信息失败", "error": err.Error()})
			return
		}
		logger.Server.Infof("DeployServer | packageInfo | %v", packageInfo)
		packageFileName := packageInfo.FileName
		tarPath := filepath.Join(constant.TARGET_PACKAGE_DIR, serverName, packageFileName)
		serverDir := filepath.Join(constant.TARGET_SERVANT_DIR, serverName)
		logger.Server.Infof("DeployServer | tarPath | %s | serverDir | %s", tarPath, serverDir)
		err = patchutils.T_PatchUtils.Tar2Dest(tarPath, serverDir)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "解压服务器包失败", "error": err.Error()})
			return
		}
	}

	nodes, err := mapper.T_Mapper.GetServerNodes(req.ServerId,0)
	logger.Server.Infof("DeployServer | nodes | %v", nodes)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器节点信息失败", "error": err.Error()})
		return
	}
	nodeStatFactory := entity.NewNodeStatFactory(&entity.NodeStat{
		ServerName: serverInfo.ServerName,
		ServerId:   serverInfo.ID,
	})
	for _, node := range nodes {
		if !patchutils.T_PatchUtils.Contains(req.NodeIds, node.Id) {
			continue
		}
		mapper.T_Mapper.SaveNodeStat(nodeStatFactory.Assign(&entity.NodeStat{
			TYPE:       entity.TYPE_INFO,
			ServerNodeId: node.Id,
			Content:    fmt.Sprintf("开始部署服务器 %s | 节点 %d | 版本号 | %d | 端口号 %d",serverInfo.ServerName, node.Id,req.PackageId,node.Port),
		}))
		logger.Server.Infof("DeployServer | node | %v", node)
		currentCommand := command.CenterManager.GetCommand(node.Id)
		if currentCommand != nil {
			err := currentCommand.Stop()
			if err!= nil {
				ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "停止服务器失败", "error": err.Error()})
				mapper.T_Mapper.SaveNodeStat(nodeStatFactory.Assign(&entity.NodeStat{
					TYPE:       entity.TYPE_ERROR,
					ServerNodeId: node.Id,
					Content:    fmt.Sprintf(
						"部署服务器失败 %s | 节点 %d | 版本号 | %d | 端口号 %d | 原因 %s", 
					serverInfo.ServerName,
						node.Id,
						req.PackageId,
						node.Port,
						err.Error(),
					),
				}))
				return
			}
		}
		targetFile := filepath.Join(cwd,constant.TARGET_SERVANT_DIR, serverName,execPath)
		patchServerInfo := &patchutils.ServerInfo{
			ServerType: serverType,
			ServerName: serverName,
			TargetFile: targetFile,
			BindPort:   node.Port,
			BindHost:   node.Host,
			NodeId:     node.Id,
		}
		logger.Server.Infof("DeployServer | patchServerInfo | %v", patchServerInfo)
		cmd,err := patchServerInfo.CreateCommand()
		if err != nil {
			logger.Server.Infof("DeployServer | err | %v", err)
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建服务器命令失败", "error": err.Error()})
			return
		}
		args := cmd.GetCmd().Args
		cmd.AppendEnv([]string{
			fmt.Sprintf("%s=%s", constant.SGRID_TARGET_HOST,node.Host),
			fmt.Sprintf("%s=%v", constant.SGRID_TARGET_PORT,node.Port),
		})
		logger.Server.Infof("DeployServer | args | %v", args)
		err = cmd.Start()
		command.CenterManager.AddCommand(node.Id, cmd)
		if err!= nil {
			mapper.T_Mapper.SaveNodeStat(nodeStatFactory.Assign(&entity.NodeStat{
				TYPE:       entity.TYPE_ERROR,
				ServerNodeId: node.Id,
				Content:    fmt.Sprintf(
					"启动服务器失败 %s | 节点 %d | 版本号 | %d | 端口号 %d | 原因 %s", 
				serverInfo.ServerName,
					node.Id,
					req.PackageId,
					node.Port,
					err.Error(),
				),
			}))
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "启动服务器失败", "error": err.Error()})
			return
		}
		err = command.UseCgroup(cmd)
		if err != nil {
			mapper.T_Mapper.SaveNodeStat(nodeStatFactory.Assign(&entity.NodeStat{
				TYPE:       entity.TYPE_ERROR,
				ServerNodeId: node.Id,
				Content:    fmt.Sprintf(
					"设置cgroup失败 %s | 节点 %d | 版本号 | %d | 端口号 %d | 原因 %s", 
				serverInfo.ServerName,
					node.Id,
					req.PackageId,
					node.Port,
					err.Error(),
				),
			}))
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "设置cgroup失败", "error": err.Error()})
			return
		}
		logger.Server.Infof("DeployServer | cmd | %v", cmd)
	}
	mapper.T_Mapper.SaveNodeStat(nodeStatFactory.Assign(&entity.NodeStat{
		TYPE:       entity.TYPE_SUCCESS,
		Content:    fmt.Sprintf("部署服务器成功 %s | 节点 %s | 版本号 %d", serverInfo.ServerName, req.NodeIds, req.PackageId),
	}))
	err = mapper.T_Mapper.UpdateNodePatch(req.NodeIds, req.PackageId)
	if err!= nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "更新服务器节点版本号失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "部署服务器成功"})
}

func DeployServer(ctx *gin.Context) {
	var req RUN_SERVER_REQ
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	RunServer(ctx, req, true)
}

func RestartServer(ctx *gin.Context) {
	var req RUN_SERVER_REQ
	if err := ctx.ShouldBindJSON(&req); err!= nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	RunServer(ctx, req, false)
}


func StopServer(ctx *gin.Context) {
	var req struct {
		ServerId  int   `json:"serverId"`
		NodeIds   []int `json:"nodeIds"`
	}
	if err := ctx.ShouldBindJSON(&req); err!= nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	serverInfo, err := mapper.T_Mapper.GetServerInfo(req.ServerId)
	if err!= nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败", "error": err.Error()})
		return
	}

	for _, nodeId := range req.NodeIds {
		currentCommand := command.CenterManager.GetCommand(nodeId)
		if currentCommand!= nil {
			err := currentCommand.Stop()
			if err!= nil {
				mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
					ServerName: serverInfo.ServerName,
					ServerId:   serverInfo.ID,
					TYPE:       entity.TYPE_ERROR,
					Content:    fmt.Sprintf("关停服务失败 %s | 节点 %s | 原因 %s", serverInfo.ServerName, req.NodeIds,err.Error()),
					ServerNodeId: nodeId,
				})
				ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "停止服务器失败", "error": err.Error()})
				return
			}
			command.CenterManager.RemoveCommand(nodeId)
			mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
				ServerName: serverInfo.ServerName,
				ServerId:   serverInfo.ID,
				TYPE:       entity.TYPE_SUCCESS,
				Content:    fmt.Sprintf("关停服务 %s | 节点 %s", serverInfo.ServerName, req.NodeIds),
			})
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "停止服务器成功"})
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
	res, err := mapper.T_Mapper.GetServerNodes(req.Id,0)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": res})
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
