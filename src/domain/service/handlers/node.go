package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sgridnext.com/src/constant"
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
