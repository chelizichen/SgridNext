package probe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
)

func RunProbeTask() {
	arr := config.Conf.GetNewest("networkPrefixs")
	results := Probe(strings.Split(arr, ","))
	nodes, err := mapper.T_Mapper.GetNodeList()
	nodeMap := make(map[string]entity.Node)
	for _, node := range nodes {
		nodeMap[node.Host] = node
	}
	if err != nil {
		logger.Probe.Errorf("获取节点列表失败: %s", err)
		return
	}
	for _, result := range results {
		// 如果探针成功 并且 nodes 表存在，且状态为已下线，则改为已上线
		// 如果探针成功 并且 nodes 表不存在，则创建
		// 如果探针成功 并且 nodes 表存在，且状态为已上线，则不处理
		// 如果探针失败，则不处理
		if result.Status == "成功" {
			logger.Probe.Infof("探针成功: %s", result.IP)
			if _, ok := nodeMap[result.IP]; ok {
				if nodeMap[result.IP].NodeStatus == constant.COMM_STATUS_ONLINE {
					logger.Probe.Infof("节点已上线，无需处理: %s", result.IP)
					continue
				} else {
					mapper.T_Mapper.UpdateNodeStatus(nodeMap[result.IP].ID, constant.COMM_STATUS_ONLINE)
					logger.Probe.Infof("上线节点: %s", result.IP)
				}
			} else {
				logger.Probe.Infof("节点不存在，创建节点: %s", result.IP)
				id, _ := mapper.T_Mapper.CreateNode(&entity.Node{
					Host:       result.IP,
					NodeStatus: constant.COMM_STATUS_ONLINE,
					CreateTime: constant.GetCurrentTime(),
					UpdateTime: constant.GetCurrentTime(),
					Cpus:       0,
					Memory:     0,
					Alias:      "[Probe]" + result.IP,
					Os:         "[Probe]" + "OS",
				})
				confObj := constant.ConfObj{
					Host:      result.IP,
					Db:        config.Conf.Get("db"),
					DbType:    config.Conf.Get("dbtype"),
					NodeIndex: fmt.Sprint(id),
					MainNode:  fmt.Sprintf("http://%s:15872", config.Conf.Get("host")),
				}
				confStr, _ := json.Marshal(confObj)

				// 创建gRPC连接
				address := fmt.Sprintf("%s:%s", result.IP, constant.NODE_PORT)

				conn, err := grpc.NewClient(address,
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				)
				if err != nil {
					return
				}
				defer conn.Close()

				// 创建客户端
				client := protocol.NewNodeServantClient(conn)

				// 调用Probe接口
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()

				_, err = client.Probe(ctx, &protocol.ProbeReq{
					Conf: string(confStr),
					Type: 2,
				})
				if err != nil {
					logger.Probe.Errorf("创建节点失败: %s", err)
					return
				}

				logger.Probe.Infof("创建节点: %s", result.IP)
			}
		} else {
			logger.Probe.Infof("探针失败: %s", result.IP)
		}
	}
}
