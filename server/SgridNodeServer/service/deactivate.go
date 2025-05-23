package service

import (
	"fmt"

	"sgridnext.com/server/SgridNodeServer/command"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
)

func Deactivate(req *protocol.ActivateReq) (code int32, msg string) {
	defer func() {
		if err := recover(); err != nil {
			logger.CMD.Error("recover | Deactivate error | %s", err)
			code = CODE_FAIL
			msg = "Deactivate error"
		}
	}()
	serverId := int(req.ServerId)
	serverInfo, err := mapper.T_Mapper.GetServerInfo(serverId)
	serverNodeIds := constant.ConvertToIntSlice(req.ServerNodeIds)
	if err != nil {
		return CODE_FAIL, "获取服务器信息失败" + err.Error()
	}
	for _, nodeId := range serverNodeIds {
		currentCommand := command.CenterManager.GetCommand(nodeId)
		if currentCommand != nil {
			logger.CMD.Info("Deactivate | serverName: %s | nodeId: %d | command is nil: %v | pid: %d", 
				serverInfo.ServerName,
			 	nodeId,
			 	currentCommand.GetCmd() == nil,
			 	currentCommand.GetPid(),
			)
			if currentCommand.GetCmd() == nil && currentCommand.GetPid() > 0{
				logger.CMD.Info("Deactivate | Kill pid | %d",currentCommand.GetPid())
				err := command.Kill(currentCommand.GetPid())
				if err!= nil {
					logger.CMD.Error("Deactivate | Kill pid error | %s ",err.Error())
				}
			}else{
				err := currentCommand.Stop()
				if err != nil {
					mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
						ServerName:   serverInfo.ServerName,
						ServerId:     serverInfo.ID,
						TYPE:         entity.TYPE_ERROR,
						Content:      fmt.Sprintf("关停服务失败 %s | 节点 %s | 原因 %s", serverInfo.ServerName, req.ServerNodeIds, err.Error()),
						ServerNodeId: nodeId,
					})
					return CODE_FAIL, "关停服务失败" + err.Error()
				}
				command.CenterManager.RemoveCommand(nodeId)
				mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
					ServerName: serverInfo.ServerName,
					ServerId:   serverInfo.ID,
					TYPE:       entity.TYPE_SUCCESS,
					Content:    fmt.Sprintf("关停服务 %s | 节点 %s", serverInfo.ServerName, req.ServerNodeIds),
				})
			}
		}
	}
	return CODE_SUCCESS, MSG_SUCCESS
}
