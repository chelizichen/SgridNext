package service

import (
	"fmt"

	"sgridnext.com/server/SgridNodeServer/command"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/patchutils"
)

func CheckStat(in *protocol.CheckStatReq) (code int, msg string) {
	serverId := int(in.ServerId)
	serverNodeIds := constant.ConvertToIntSlice(in.NodeIds)
	serverInfo, err := mapper.T_Mapper.GetServerInfo(serverId)
	if err != nil {
		return CODE_FAIL, "server not found" + err.Error()
	}
	localBindServerNodes, err := mapper.T_Mapper.GetServerNodes(serverId, config.Conf.GetLocalNodeId())
	for _, v := range localBindServerNodes {
		// 检查本地是否存在
		if !patchutils.T_PatchUtils.Contains(serverNodeIds, v.Id) {
			continue
		}
		c := command.CenterManager.GetCommand(v.Id)
		if c == nil {
			mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
				ServerId:     serverId,
				ServerNodeId: v.Id,
				TYPE:         entity.TYPE_CHECK,
				Content:      fmt.Sprintf("nodeId %d is not alive", v),
				ServerName:   serverInfo.ServerName,
			})
			mapper.T_Mapper.UpdateNodeStatus(v.Id, constant.COMM_STATUS_OFFLINE)
			continue
		}
		pid, alive, err := c.CheckStat()
		if err != nil {
			mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
				ServerId:     serverId,
				ServerNodeId: v.Id,
				TYPE:         entity.TYPE_CHECK,
				Content:      fmt.Sprintf("nodeId %v is not alive ,error %s", v, err.Error()),
				ServerName:   serverInfo.ServerName,
			})
			mapper.T_Mapper.UpdateNodeStatus(v.Id, constant.COMM_STATUS_OFFLINE)
			continue
		}
		mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
			ServerId:     serverId,
			ServerNodeId: v.Id,
			TYPE:         entity.TYPE_CHECK,
			Content: fmt.Sprintf("serverName: %s, nodeId: %d, pid: %d, alive: %v",
				serverInfo.ServerName, v, pid, alive),
			ServerName: serverInfo.ServerName,
		})
		mapper.T_Mapper.UpdateNodeStatus(v.Id, constant.COMM_STATUS_ONLINE)
	}
	return CODE_SUCCESS, ""
}
