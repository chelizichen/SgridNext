package service

import (
	"fmt"
	"os"
	"path/filepath"

	"sgridnext.com/server/SgridNodeServer/api"
	"sgridnext.com/server/SgridNodeServer/command"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
	"sgridnext.com/src/patchutils"
)

func Acitvate(req *protocol.ActivateReq) (code int32, msg string) {
	cwd, _ := os.Getwd()
	needDeploy := req.Type == constant.ACITVATE_DEPLOY
	serverId := int(req.ServerId)
	serverNodeIds := constant.ConvertToIntSlice(req.ServerNodeIds)
	packageId := int(req.PackageId)
	localNodeId := config.Conf.GetLocalNodeId()
	logger.Server.Infof("DeployServer | req | %v", req)
	serverInfo, err := mapper.T_Mapper.GetServerInfo(serverId)
	if err != nil {
		return CODE_FAIL, "获取服务信息失败:" + err.Error()
	}
	logger.Server.Infof("DeployServer | serverInfo | %v", serverInfo)
	execPath := serverInfo.ExecFilePath
	serverName := serverInfo.ServerName
	serverType := serverInfo.ServerType

	// 拉取配置文件
	err = api.GetConfigList(api.GetConfigListReq{
		ServerId: serverId,
	})
	if err != nil {
		return CODE_FAIL, "获取配置文件失败: " + err.Error()
	}

	// 部署模式，需要解压文件到目标目录，然后再启动execPath
	if needDeploy {
		// TODO 如果没有包，则需要下载
		packageInfo, err := mapper.T_Mapper.GetServerPackageInfo(packageId)
		if err != nil {
			return CODE_FAIL, err.Error()
		}
		logger.Server.Infof("DeployServer | packageInfo | %v", packageInfo)
		packageFileName := packageInfo.FileName
		tarPath := filepath.Join(constant.TARGET_PACKAGE_DIR, serverName, packageFileName)
		if _, err := os.Stat(tarPath); err != nil {
			err = api.GetFile(api.FileReq{
				FileName: packageFileName,
				Type:     constant.FILE_TYPE_PACKAGE,
				ServerId: serverId,
			})
			if err != nil {
				return CODE_FAIL, "下载服务包失败: " + err.Error()
			}
		}
		serverDir := filepath.Join(constant.TARGET_SERVANT_DIR, serverName)
		logger.Server.Infof("DeployServer | tarPath | %s | serverDir | %s", tarPath, serverDir)
		err = patchutils.T_PatchUtils.Tar2Dest(tarPath, serverDir)
		if err != nil {
			return CODE_FAIL, "解压服务包失败: " + err.Error()
		}
	}

	nodes, err := mapper.T_Mapper.GetServerNodes(serverId, localNodeId)
	logger.Server.Infof("DeployServer | nodes | %v", nodes)
	if err != nil {
		return CODE_FAIL, "获取服务节点失败: " + err.Error()
	}
	nodeStatFactory := entity.NewNodeStatFactory(&entity.NodeStat{
		ServerName: serverInfo.ServerName,
		ServerId:   serverInfo.ID,
	})
	// 遍历 当前节点下的服务节点列表，找出需要激活的节点
	for _, node := range nodes {
		if !patchutils.T_PatchUtils.Contains(serverNodeIds, node.Id) {
			continue
		}
		mapper.T_Mapper.SaveNodeStat(nodeStatFactory.Assign(&entity.NodeStat{
			TYPE:         entity.TYPE_INFO,
			ServerNodeId: node.Id,
			Content:      fmt.Sprintf("开始部署服务器 %s | 节点 %d | 版本号 | %d | 端口号 %d", serverInfo.ServerName, node.Id, req.PackageId, node.Port),
		}))
		logger.Server.Infof("DeployServer | node | %v", node)
		currentCommand := command.CenterManager.GetCommand(node.Id)
		if currentCommand != nil {
			err := currentCommand.Stop()
			if err != nil {
				mapper.T_Mapper.SaveNodeStat(nodeStatFactory.Assign(&entity.NodeStat{
					TYPE:         entity.TYPE_ERROR,
					ServerNodeId: node.Id,
					Content:      fmt.Sprintf("部署服务器失败 %s | 节点 %d | 版本号 | %d | 端口号 %d | 原因 %s", serverInfo.ServerName, node.Id, req.PackageId, node.Port, err.Error()),
				}))
				return CODE_FAIL, "停止服务器失败: " + err.Error()
			}
		}
		targetFile := filepath.Join(cwd, constant.TARGET_SERVANT_DIR, serverName, execPath)
		patchServerInfo := &command.ServerInfo{
			ServerType: serverType,
			ServerName: serverName,
			TargetFile: targetFile,
			BindPort:   node.Port,
			BindHost:   node.Host,
			NodeId:     node.Id,
		}
		logger.Server.Infof("DeployServer | patchServerInfo | %v", patchServerInfo)
		cmd, err := patchServerInfo.CreateCommand()
		cmd.SetHost(node.Host)
		cmd.SetPort(node.Port)
		cmd.SetLocalMachineId(localNodeId)
		cmd.SetServerId(serverId)
		if err != nil {
			logger.Server.Infof("DeployServer | err | %v", err)
			// ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建服务器命令失败", "error": err.Error()})
			return CODE_FAIL, err.Error()
		}
		args := cmd.GetCmd().Args
		cmd.AppendEnv([]string{
			fmt.Sprintf("%s=%s", constant.SGRID_TARGET_HOST, node.Host),
			fmt.Sprintf("%s=%v", constant.SGRID_TARGET_PORT, node.Port),
		})
		logger.Server.Infof("DeployServer | args | %v", args)
		err = cmd.Start()
		command.CenterManager.AddCommand(node.Id, cmd)
		if err != nil {
			mapper.T_Mapper.SaveNodeStat(nodeStatFactory.Assign(&entity.NodeStat{
				TYPE:         entity.TYPE_ERROR,
				ServerNodeId: node.Id,
				Content:      fmt.Sprintf("启动服务器失败 %s | 节点 %d | 版本号 | %d | 端口号 %d | 原因 %s", serverInfo.ServerName, node.Id, req.PackageId, node.Port, err.Error()),
			}))
			return CODE_FAIL, "启动服务器失败: " + err.Error()
		}
		err = command.UseCgroup(cmd)
		if err != nil {
			mapper.T_Mapper.SaveNodeStat(nodeStatFactory.Assign(&entity.NodeStat{
				TYPE:         entity.TYPE_ERROR,
				ServerNodeId: node.Id,
				Content:      fmt.Sprintf("设置cgroup失败 %s | 节点 %d | 版本号 | %d | 端口号 %d | 原因 %s", serverInfo.ServerName, node.Id, req.PackageId, node.Port, err.Error()),
			}))
			return CODE_FAIL, "设置cgroup失败:" + err.Error()
		}
		logger.Server.Infof("DeployServer | cmd | %v", cmd)
	}
	mapper.T_Mapper.SaveNodeStat(nodeStatFactory.Assign(&entity.NodeStat{
		TYPE:    entity.TYPE_SUCCESS,
		Content: fmt.Sprintf("部署服务器成功 %s | 节点 %s | 版本号 %d", serverInfo.ServerName, constant.ConvertToIntSlice(req.ServerNodeIds), req.PackageId),
	}))
	if needDeploy {
		err = mapper.T_Mapper.UpdateNodePatch(serverNodeIds, packageId)
		if err != nil {
			return CODE_FAIL, "更新服务器节点版本号失败 :" + err.Error()
		}
	}
	return CODE_SUCCESS, MSG_SUCCESS
}
