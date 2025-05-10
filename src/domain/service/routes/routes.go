package routes

import (
	"github.com/gin-gonic/gin"
	"sgridnext.com/src/domain/service/handlers"
)

func LoadRoutes(engine *gin.Engine) {
	// 创建服务
	engine.POST("/server/createServer", handlers.CreateServer)
	// 获取服务信息
	engine.POST("/server/getServerInfo", handlers.GetServerInfo)
	// 上传服务包
	engine.POST("/server/uploadPackage", handlers.CreatePackage)
	// 上传/修改配置文件
	engine.POST("/server/upsertConfig", handlers.UpsertConfig)
	// 获取配置文件内容
	engine.POST("/server/getConfigContent", handlers.GetConfigContent)
	// 创建服务部署节点
	engine.POST("/server/createServerNode", handlers.CreateServerNode)
	// 创建服务组
	engine.POST("/server/createGroup", handlers.CreateGroup)
	// 创建机器节点
	engine.POST("/server/createNode", handlers.CreateNode)
	// 部署服务
	engine.POST("/server/deployServer", handlers.DeployServer)
	// 停止服务
	engine.POST("/server/stopServer", handlers.StopServer)
	// 重启服务「不更改所在包文件」
	engine.POST("/server/restartServer", handlers.RestartServer)
	// 获取服务节点状态
	engine.POST("/server/getServerNodesStatus", handlers.GetServerNodesStatus)
	// 检查服务节点状态
	engine.POST("/server/checkServerNodesStatus", handlers.CheckServerNodesStatus)
	// 获取服务节点状态日志
	engine.POST("/server/getServerNodesLog", handlers.GetServerNodesLog)
	// 获取服务节点
	engine.POST("/server/getServerNodes", handlers.GetServerNodes)
	// 获取服务配置文件列表
	engine.POST("/server/getServerConfigList", handlers.GetServerConfigList)
	// 获取服务包列表
	engine.POST("/server/getServerPackageList", handlers.GetServerPackageList)
	// 获取服务列表
	engine.POST("/server/getServerList", handlers.GetServerList)
	// 获取机器节点列表
	engine.POST("/server/getNodeList", handlers.GetNodeList)
	// 获取机器一段时间的负载详情
	engine.POST("/server/getNodeLoadDetail", handlers.GetNodeLoadDetail)
	// 获取服务组列表
	engine.POST("/server/getGroupList", handlers.GetGroupList)
	// CGROUP 设置 服务 LIMIt CPU
	engine.POST("/server/cgroup/setCpuLimit", handlers.SetCpuLimit)
	// CGROUP 设置 服务 LIMIt MEM
	engine.POST("/server/cgroup/setMemLimit", handlers.SetMemoryLimit)
	// CGROUP 获取 服务 STATUS
	engine.POST("/server/cgroup/getStatus", handlers.GetStatus)
	// 登录
	engine.POST("/login", handlers.Login)
}
