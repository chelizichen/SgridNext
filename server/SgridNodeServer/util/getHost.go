package util

import (
	"fmt"

	"sgridnext.com/src/config"
	"sgridnext.com/src/domain/service/mapper"
)

func defaultGetHost() string {
	// 先从配置文件取
	host := config.Conf.Get(config.KEY_HOST)
	if host != "" {
		idx := config.Conf.Get(config.KEY_NODE_INDEX)
		if idx == "" {
			nodeId,err := mapper.T_Mapper.GetNodeIdByHost(host)
			if err != nil {
				fmt.Printf("获取节点ID失败 %v",err.Error())
				return host
			}
			config.Conf.Set(config.KEY_NODE_INDEX, fmt.Sprint(nodeId))
		}
		return host
	}

	// 从节点取
	host, err := mapper.T_Mapper.GetHost(config.Conf.GetLocalNodeId())
	if err != nil {
		return ""
	}
	// merge 配置文件和节点的 host
	config.Conf.Set(config.KEY_HOST, host)
	return host
}

func GetHost() string {
	return defaultGetHost()
}
