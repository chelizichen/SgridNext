package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/proxy"
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
	var callRsp []*protocol.BasicRes
	proxy.ProxyMap.FullDispatch(func(client *protocol.NodeServantClient) error {
		rsp, err := (*client).CheckStat(context.Background(), &protocol.CheckStatReq{
			ServerId:      int32(req.ServerId),
			NodeIds: constant.ConvertToInt32Slice(req.ServerNodeIds),
		})
		callRsp = append(callRsp, rsp)
		return err
	})
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "已检查完成"})
}