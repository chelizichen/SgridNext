package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"sgridnext.com/server/SgridNodeServer/cgroupmanager"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
	"sgridnext.com/src/patchutils"
	"sgridnext.com/src/proxy"
)

func SetCpuLimit(ctx *gin.Context) {
	var req struct {
		NodeIds  []int   `json:"nodeIds"`
		CpuLimit float64 `json:"cpuLimit"`
		ServerId int     `json:"serverId"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Server.Infof("SetCpuLimit Error | %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	var callRsp []*protocol.BasicRes
	proxy.ProxyMap.FullDispatch(func(client *protocol.NodeServantClient) error {
		rsp, err := (*client).CgroupLimit(context.Background(), &protocol.CgroupLimitReq{
			ServerId:      int32(req.ServerId),
			NodeIds: constant.ConvertToInt32Slice(req.NodeIds),
			Type:          constant.CGROUP_TYPE_CPU,
			Value:         float32(req.CpuLimit),
		})
		callRsp = append(callRsp, rsp)
		return err
	})
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "设置CPU限制成功"})
}

func GetStatus(ctx *gin.Context) {
	var req struct {
		NodeIds  []int `json:"nodeIds"`
		ServerId int   `json:"serverId"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	var statsRsp []*cgroupmanager.CgroupStats

	localBindServerNodes, err := mapper.T_Mapper.GetServerNodes(req.ServerId, config.Conf.GetLocalNodeId())
	if err != nil {
		logger.Server.Infof("GetStatus Error | %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取本地节点绑定服务节点列表信息失败"})
		return
	}
	serverInfo, err := mapper.T_Mapper.GetServerInfo(req.ServerId)
	if err != nil {
		logger.Server.Infof("GetStatus Error | %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取服务器信息失败"})
		return
	}
	logger.Server.Infof("SetCpuLimit args | %v", req)
	for _, node := range localBindServerNodes {
		if !patchutils.T_PatchUtils.Contains(req.NodeIds, node.Id) {
			continue
		}
		scg := &cgroupmanager.SgridCgroup{
			ServerName: serverInfo.ServerName,
			NodeId:     node.Id,
		}
		name := scg.GetCgroupName()
		cgroup, err := cgroupmanager.NewCgroupManager(name)
		if err != nil {
			logger.Server.Infof("SetCpuLimit Error | %v", err)
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取CgroupManager失败"})
			return
		}
		stat, err := cgroup.Stat()
		if err != nil {
			logger.Server.Infof("SetCpuLimit Error | %v", err)
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "设置CPU限制失败"})
			return
		}
		statsRsp = append(statsRsp, stat)
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "获取信息成功", "data": statsRsp})
}

// 删除节点 Cgroup 限制
func DeleteCgroupLimit(ctx *gin.Context) {
	var req struct {
		NodeIds  []int `json:"nodeIds"`
		ServerId int   `json:"serverId"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	var callRsp []*protocol.BasicRes
	proxy.ProxyMap.FullDispatch(func(client *protocol.NodeServantClient) error {
		rsp, err := (*client).CgroupLimit(context.Background(), &protocol.CgroupLimitReq{
			ServerId:      int32(req.ServerId),
			NodeIds: constant.ConvertToInt32Slice(req.NodeIds),
			Type:          constant.CGROUP_TYPE_DELETE,
			Value:         0,
		})
		callRsp = append(callRsp, rsp)
		return err
	})
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "删除成功"})
}

func SetMemoryLimit(ctx *gin.Context) {
	var req struct {
		NodeIds     []int `json:"nodeIds"`
		MemoryLimit int64 `json:"memoryLimit"`
		ServerId    int   `json:"serverId"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	var callRsp []*protocol.BasicRes
	proxy.ProxyMap.FullDispatch(func(client *protocol.NodeServantClient) error {
		rsp, err := (*client).CgroupLimit(context.Background(), &protocol.CgroupLimitReq{
			ServerId:      int32(req.ServerId),
			NodeIds: constant.ConvertToInt32Slice(req.NodeIds),
			Type:          constant.CGROUP_TYPE_MEMORY,
			Value:         float32(req.MemoryLimit),
		})
		callRsp = append(callRsp, rsp)
		return err
	})
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "设置内存限制成功"})
}
