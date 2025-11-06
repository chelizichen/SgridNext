package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
	"sgridnext.com/src/proxy"
	"sgridnext.com/src/resource"
)

// GetNodeResource 获取节点资源信息
func GetNodeResource(c *gin.Context) {
	var req struct {
		NodeId int `json:"nodeId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}
	if req.NodeId == 0 {
		cwd, _ := os.Getwd()
		logger.App.Infof("工作目录: %s", cwd)
		// 获取节点资源信息
		nodeResource := resource.GetNodeResource(cwd)
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "获取节点资源信息成功",
			"data":    nodeResource,
		})
	} else {
		host, err := mapper.T_Mapper.GetHost(req.NodeId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取节点信息失败"})
			return
		}
		proxy.ProxyMap.DispatchByHost(host, func(client *protocol.NodeServantClient) error {
			nodeResource, err := (*client).GetNodeResource(context.Background(), &emptypb.Empty{})
			if err != nil {
				logger.Alive.Errorf("节点 | %s | 获取节点资源信息异常 | %s", req.NodeId, err.Error())
				return err
			}
			var nodeResourceData resource.NodeResource
			err = json.Unmarshal([]byte(nodeResource.Data), &nodeResourceData)
			if err != nil {
				logger.Alive.Errorf("节点 | %s | 获取节点资源信息异常 | %s", req.NodeId, err.Error())
				return err
			}
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"success": true,
				"msg":     "获取节点资源信息成功",
				"data":    nodeResourceData,
			})
			logger.Alive.Infof("节点 | %s | 获取节点资源信息成功 | %s", req.NodeId, nodeResourceData)
			return nil
		})
	}

}
