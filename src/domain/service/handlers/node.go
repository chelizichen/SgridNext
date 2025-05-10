package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/command"
	"sgridnext.com/src/domain/config"
	"sgridnext.com/src/domain/patchutils"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/domain/service/mapper"
)

func GetNodeList(ctx *gin.Context) {
	nodes, err := mapper.T_Mapper.GetNodeList()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取节点列表失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": nodes})
}

func GetNodeLoadDetail(ctx *gin.Context) {

}

func CreateNode(ctx *gin.Context) {
	var req struct {
		Host   string `json:"Host"`
		Os     string `json:"Os"`
		Memory int    `json:"Memory"`
		Cpus   int    `json:"cpus"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	node := &entity.Node{
		Host:       req.Host,
		Os:         req.Os,
		Memory:     req.Memory,
		Cpus:       req.Cpus,
		CreateTime: constant.GetCurrentTime(),
		ID:         0,
		NodeStatus: constant.COMM_STATUS_ONLINE,
	}
	id, err := mapper.T_Mapper.CreateNode(node)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "创建节点成功", "data": id})
}


func GetServerNodesStatus(ctx *gin.Context) {
	var req struct {
		NodeId       int    `json:"node_id,omitempty"`
		ServerId     int    `json:"server_id,omitempty"`
		ServerNodeId int    `json:"server_node_id,omitempty"`
		TYPE         int    `json:"type,omitempty"`
		Offset       int    `json:"offset,omitempty"`
		Size         int    `json:"size,omitempty"`
	}
	if err := ctx.ShouldBindJSON(&req); err!= nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	
	rsp,err := mapper.T_Mapper.GetNodeStatList(&entity.NodeStat{
		ServerId:     req.ServerId,
		ServerNodeId: req.ServerNodeId,
		NodeId:       req.NodeId,
		TYPE:         req.TYPE,
	}, req.Offset, req.Size)

	if err!= nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取节点状态列表失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": rsp})
}

func CheckServerNodesStatus(ctx *gin.Context) {
	var req struct {
		ServerId     int    `json:"server_id,omitempty"`
		ServerNodeIds []int    `json:"server_node_ids,omitempty"`
	}
	if err := ctx.ShouldBindJSON(&req); err!= nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	serverInfo,err := mapper.T_Mapper.GetServerInfo(req.ServerId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务信息失败", "error": "server not found"})
		return
	}
	localBindServerNodes,err := mapper.T_Mapper.GetServerNodes(req.ServerId,config.Conf.GetLocalNodeId())
	for _, v := range localBindServerNodes {
		// 检查本地是否存在
		if !patchutils.T_PatchUtils.Contains(req.ServerNodeIds, v.Id) {
			continue
		}
		c := command.CenterManager.GetCommand(v.Id)
		if c == nil {
			mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
				ServerId:     req.ServerId,
				ServerNodeId: v.Id,
				TYPE:         entity.TYPE_CHECK,
				Content: 		fmt.Sprintf("nodeId %d is not alive", v),
				ServerName:   serverInfo.ServerName,
			})
			mapper.T_Mapper.UpdateNodeStatus(v.Id, constant.COMM_STATUS_OFFLINE)
			continue
		}
		pid,alive,err := c.CheckStat()
		if err != nil {
			mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
				ServerId:     req.ServerId,
				ServerNodeId: v.Id,
				TYPE:         entity.TYPE_CHECK,
				Content: 		fmt.Sprintf("nodeId %d is not alive", v),
				ServerName:   serverInfo.ServerName,
			})
			mapper.T_Mapper.UpdateNodeStatus(v.Id, constant.COMM_STATUS_OFFLINE)
			continue
		}
		mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
			ServerId:     req.ServerId,
			ServerNodeId: v.Id,
			TYPE:         entity.TYPE_CHECK,
			Content: 	fmt.Sprintf("serverName: %s, nodeId: %d, pid: %d, alive: %v", 
													serverInfo.ServerName,v,pid, alive),
			ServerName:   serverInfo.ServerName,
		})
		mapper.T_Mapper.UpdateNodeStatus(v.Id, constant.COMM_STATUS_ONLINE)
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "已检查完成"})
}