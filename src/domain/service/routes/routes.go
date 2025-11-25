package routes

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"sgridnext.com/src/domain/service/handlers"
	"sgridnext.com/src/webshell"
)

func LoadRoutes(engine *gin.Engine) {
	group := engine.Group("/api")
	// 创建服务
	group.POST("/server/createServer", handlers.CreateServer)
	// 更新服务
	group.POST("/server/updateServer", handlers.UpdateServer)
	// 获取服务信息
	group.POST("/server/getServerInfo", handlers.GetServerInfo)
	// 上传服务包
	group.POST("/server/uploadPackage", handlers.CreatePackage)
	// 上传/修改配置文件
	group.POST("/server/upsertConfig", handlers.UpsertConfig)
	// 获取配置文件内容
	group.POST("/server/getConfigContent", handlers.GetConfigContent)
	// 创建服务部署节点
	group.POST("/server/createServerNode", handlers.CreateServerNode)
	// 更新节点
	group.POST("/server/updateServerNode", handlers.UpdateServerNode)
	// 删除服务节点
	group.POST("/server/deleteServerNode", handlers.DeleteServerNode)
	// 创建服务组
	group.POST("/server/createGroup", handlers.CreateGroup)
	// 创建机器节点
	group.POST("/server/createNode", handlers.CreateNode)
	// 更新机器节点
	group.POST("/server/updateNode", handlers.UpdateMachineNode)
	// 更新节点别名
	group.POST("/server/updateNodeAlias", handlers.UpdateMachineNodeAlias)
	// 部署服务
	group.POST("/server/deployServer", handlers.DeployServer)
	// 停止服务
	group.POST("/server/stopServer", handlers.StopServer)
	// 重启服务「不更改所在包文件」
	group.POST("/server/restartServer", handlers.RestartServer)
	// 获取服务节点状态
	group.POST("/server/getServerNodesStatus", handlers.GetServerNodesStatus)
	// 检查服务节点状态
	group.POST("/server/checkServerNodesStatus", handlers.CheckServerNodesStatus)
	// 获取服务节点状态日志
	group.POST("/server/getServerNodesLog", handlers.GetServerNodesLog)
	// 获取服务节点
	group.POST("/server/getServerNodes", handlers.GetServerNodes)
	// 获取服务配置文件列表
	group.POST("/server/getServerConfigList", handlers.GetServerConfigList)
	// 获取服务包列表
	group.POST("/server/getServerPackageList", handlers.GetServerPackageList)
	// 获取服务列表
	group.POST("/server/getServerList", handlers.GetServerList)
	// 获取机器节点列表
	group.POST("/server/getNodeList", handlers.GetNodeList)
	// 获取机器一段时间的负载详情
	group.POST("/server/getNodeLoadDetail", handlers.GetNodeLoadDetail)
	// 获取服务组列表
	group.POST("/server/getGroupList", handlers.GetGroupList)
	// CGROUP 设置 服务 LIMIT CPU
	group.POST("/server/cgroup/setCpuLimit", handlers.SetCpuLimit)
	// CGROUP 设置 服务 LIMIT MEM
	group.POST("/server/cgroup/setMemLimit", handlers.SetMemoryLimit)
	// CGROUP 获取 服务 STATUS
	group.POST("/server/cgroup/getStatus", handlers.GetStatus)
	// 获取文件
	group.POST("/server/getFile", handlers.GetFile)
	// 获取发布时对应的配置文件列表
	group.POST("/server/getConfigList", handlers.GetConfigList)
	// 获取同步的节点状态 0代表主控
	group.POST("/server/getSyncStatus", handlers.GetSyncStatus)
	// 登录
	group.POST("/login", handlers.Login)

	// 主控配置管理
	group.POST("/config/getMainConfig", handlers.GetMainConfig)
	group.POST("/config/updateMainConfig", handlers.UpdateMainConfig)
	group.POST("/config/getConfigItem", handlers.GetConfigItem)
	group.POST("/config/setConfigItem", handlers.SetConfigItem)

	group.POST("/server/scripts/deploy", handlers.DeployScripts)
	group.POST("/server/downloadFile", handlers.DownloadFile)
	group.POST("/server/getFileList", handlers.GetFileList)
	group.POST("/server/getLog", handlers.GetLog)
	// 上传文件
	group.POST("/server/syncUploadFile", handlers.SyncUploadFile)
	// 探针
	group.POST("/probe/runProbeTask", handlers.RunProbeTask)
	// 获取节点资源信息
	group.POST("/resource/getNodeResource", handlers.GetNodeResource)
	// WebShell WebSocket
	group.GET("/webshell/ws", webshell.HandleWebSocket)
	// 前端静态文件
	cwd, _ := os.Getwd()
	root := filepath.Join(cwd, "dist")
	fmt.Println("web root:", root)
	engine.Static("/sgridnext/", root)
}
