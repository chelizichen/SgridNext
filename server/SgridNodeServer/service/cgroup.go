package service

import (
	"sgridnext.com/server/SgridNodeServer/cgroupmanager"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
	"sgridnext.com/src/patchutils"
)

func CgroupLimit(req *protocol.CgroupLimitReq) (code int32, msg string) {
	serverId := int(req.ServerId)
	serverNodeIds := constant.ConvertToIntSlice(req.NodeIds)
	localBindServerNodes, err := mapper.T_Mapper.GetServerNodes(serverId, config.Conf.GetLocalNodeId())
	if err != nil {
		logger.Server.Infof("SetCpuLimit Error | %v", err)
		return CODE_FAIL, "获取本地节点绑定服务节点列表信息失败" + err.Error()
	}
	serverInfo, err := mapper.T_Mapper.GetServerInfo(serverId)
	if err != nil {
		logger.Server.Infof("SetCpuLimit Error | %v", err)
		return
	}
	logger.Server.Infof("SetCpuLimit args | %v", req)
	for _, node := range localBindServerNodes {
		if !patchutils.T_PatchUtils.Contains(serverNodeIds, node.Id) {
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
			return CODE_FAIL, "获取CgroupManager失败" + err.Error()
		}
		if req.Type == constant.CGROUP_TYPE_CPU {
			err = cgroup.SetCPULimit(float64(req.Value))
			if err != nil {
				logger.Server.Infof("SetCpuLimit Error | %v", err)
				return CODE_FAIL, "设置CPU限制失败" + err.Error()
			}
			// 保存到数据库
			mapper.T_Mapper.UpsertServerNodeLimit(&entity.ServerNodeLimit{
				ServerId:     serverId,
				ServerNodeId: node.Id,
				CpuLimit:     float64(req.Value),
			})
		}
		if req.Type == constant.CGROUP_TYPE_MEMORY {
			err = cgroup.SetMemoryLimit(int64(req.Value))
			if err != nil {
				logger.Server.Infof("SetMemoryLimit Error | %v", err)
				return CODE_FAIL, "设置内存限制失败" + err.Error()
			}
			// 保存到数据库
			mapper.T_Mapper.UpsertServerNodeLimit(&entity.ServerNodeLimit{
				ServerId:     int(req.ServerId),
				ServerNodeId: node.Id,
				MemoryLimit:  int64(req.Value),
			})
		}
		if req.Type == constant.CGROUP_TYPE_DELETE {
			name := scg.GetCgroupName()
			cgroup, err := cgroupmanager.NewCgroupManager(name)
			if err != nil {
				logger.Server.Infof("DeleteCgroupLimit Error | %v", err)
				return CODE_FAIL, "获取CgroupManager失败" + err.Error()
			}
			err = cgroup.Remove()
			if err != nil {
				logger.Server.Infof("DeleteCgroupLimit Error | %v", err)
				return CODE_FAIL, "Cgroup删除失败" + err.Error()
			}
		}

	}
	return CODE_SUCCESS, "success"
}
