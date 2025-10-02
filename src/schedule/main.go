package schedule

import (
	"strings"
	"time"

	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
	"sgridnext.com/src/probe"
)

func wrapFunc(t *time.Ticker, cb func()) {
	for range t.C {
		cb()
	}
}

func runProbeCallback() {
	arr := config.Conf.GetNewest("networkPrefixs")
	results := probe.Probe(strings.Split(arr, ","))
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
			if _, ok := nodeMap[result.IP]; ok {
				if nodeMap[result.IP].NodeStatus == constant.COMM_STATUS_ONLINE {
					continue
				} else {
					mapper.T_Mapper.UpdateNodeStatus(nodeMap[result.IP].ID, constant.COMM_STATUS_ONLINE)
				}
			} else {
				mapper.T_Mapper.CreateNode(&entity.Node{
					Host:       result.IP,
					NodeStatus: constant.COMM_STATUS_ONLINE,
					CreateTime: constant.GetCurrentTime(),
					UpdateTime: constant.GetCurrentTime(),
					Cpus:       0,
					Memory:     0,
					Alias:      "[Probe]" + result.IP,
					Os:         "[Probe]" + "OS",
				})
			}
		}
	}
}

func LoadProbe() {
	go wrapFunc(time.NewTicker(3*time.Hour), runProbeCallback)
}
