package schedule

import (
	"fmt"
	"time"

	"sgridnext.com/server/SgridNodeServer/command"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/server/SgridNodeServer/service"
	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
	"sgridnext.com/src/patchutils"
)

func runRestartTick() {
	ticker := time.NewTicker(30 * time.Second)
	localNodeId := config.Conf.GetLocalNodeId()
	go func() {
		for range ticker.C {
			// 定时查需要需要重启的，然后在 stat.json 里面看是否存在，不存在的话就需要启动
			nodes, err := mapper.T_Mapper.GetServerNodes(0, localNodeId)
			if err != nil {
				logger.Alive.Errorf("get server node error %s", err.Error())
				continue
			}
			statList := command.LoadStatList()
			if statList == nil {
				// 首次启动没有 stat.json 忽略
				continue
			}
			currOnlineNodeIds := make([]int, 0)
			for _,v := range statList.StatList{
                currOnlineNodeIds = append(currOnlineNodeIds,v.NodeId )
            }
            needRestartNodeMap := make(map[int][]int)

			for _, v := range nodes {
				if v.ServerRunType != constant.SERVER_RUN_TYPE_RESTART_ALWAYS {
					continue
				}
                if patchutils.T_PatchUtils.Contains(currOnlineNodeIds, v.Id) {
                    continue
                }
                needRestartNodeMap[v.ServerId] = append(needRestartNodeMap[v.ServerId], v.Id)
			}
			logger.Alive.Infof("需要重启的节点 | %v", needRestartNodeMap)

			for svrId, svrNodeIds := range needRestartNodeMap {
				logger.Alive.Infof("开始重启服务 | %d |  重启节点ID %v", svrId, svrNodeIds)
				mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
					TYPE:     entity.TYPE_WARN,
					Content:  fmt.Sprintf("开始重启服务 svrId %d | nodeIds %v", svrId, svrNodeIds),
					ServerId: svrId,
				})

				code, errMsg := service.Acitvate(&protocol.ActivateReq{
					ServerId:      int32(svrId),
					ServerNodeIds: constant.ConvertToInt32Slice(svrNodeIds),
				})
				if code != service.CODE_SUCCESS {
					mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
						TYPE:     entity.TYPE_ERROR,
						Content:  fmt.Sprintf("重启服务失败 svrId %d | nodeIds %v | cause %s", svrId, svrNodeIds, errMsg),
						ServerId: svrId,
					})
					logger.Alive.Errorf("重启服务失败 | %d | %v", svrId, errMsg)
					continue
				}
				mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
					TYPE:     entity.TYPE_ERROR,
					Content:  fmt.Sprintf("重启服务成功 svrId %d | nodeIds %v", svrId, svrNodeIds),
					ServerId: svrId,
				})
				logger.Alive.Infof("重启服务成功 | %d", svrId)
			}
		}
	}()
}

func LoadTick() {
	go runRestartTick()
}
