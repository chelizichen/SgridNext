package service

import (
	"fmt"
	"os"
	"path/filepath"

	"sgridnext.com/server/SgridNodeServer/api"
	"sgridnext.com/server/SgridNodeServer/command"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
	"sgridnext.com/src/patchutils"
)

// ActivateContext 激活上下文，包含激活过程中的所有必要信息
type ActivateContext struct {
	Req             *protocol.ActivateReq
	ServerInfo      *entity.Server
	NodeStatFactory *entity.NodeStatFactory
	Cwd             string
	NeedDeploy      bool
	ServerId        int
	ServerNodeIds   []int
	PackageId       int
	LocalNodeId     int
}

// NewActivateContext 创建激活上下文
func NewActivateContext(req *protocol.ActivateReq) (*ActivateContext, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("获取当前工作目录失败: %w", err)
	}

	return &ActivateContext{
		Req:           req,
		Cwd:           cwd,
		NeedDeploy:    req.Type == constant.ACTIVATE_DEPLOY,
		ServerId:      int(req.ServerId),
		ServerNodeIds: constant.ConvertToIntSlice(req.ServerNodeIds),
		PackageId:     int(req.PackageId),
		LocalNodeId:   config.Conf.GetLocalNodeId(),
	}, nil
}

// Activate 激活服务 - 重构后的主函数
func Activate(req *protocol.ActivateReq) (code int32, msg string) {
	// 创建激活上下文
	ctx, err := NewActivateContext(req)
	if err != nil {
		return CODE_FAIL, err.Error()
	}

	logger.Active.Infof("Activate | req=%+v", req)

	// 获取服务信息
	if err := ctx.loadServerInfo(); err != nil {
		return CODE_FAIL, fmt.Sprintf("获取服务信息失败: %v", err)
	}

	// 设置panic恢复
	defer ctx.setupPanicRecovery()

	// 拉取配置文件
	if err := ctx.fetchConfigFiles(); err != nil {
		return CODE_FAIL, fmt.Sprintf("获取配置文件失败: %v", err)
	}


	// 停止现有服务
	if code, msg := ctx.deactivateExistingNodes(); code != CODE_SUCCESS {
		logger.Active.Errorf("Activate | deactivateExistingNodes | error=%s", msg)
		return code, msg
	}

	// 处理部署逻辑
	if ctx.NeedDeploy {
		if err := ctx.handleDeployment(); err != nil {
			return CODE_FAIL, fmt.Sprintf("部署失败: %v", err)
		}
	}

	// 激活新节点
	if err := ctx.activateNodes(); err != nil {
		return CODE_FAIL, fmt.Sprintf("激活节点失败: %v", err)
	}

	// 更新节点版本信息
	if ctx.NeedDeploy {
		if err := ctx.updateNodePatch(); err != nil {
			return CODE_FAIL, fmt.Sprintf("更新节点版本失败: %v", err)
		}
	}

	// 记录成功状态
	ctx.logSuccessStatus()
	return CODE_SUCCESS, MSG_SUCCESS
}

// loadServerInfo 加载服务信息
func (ctx *ActivateContext) loadServerInfo() error {
	serverInfo, err := mapper.T_Mapper.GetServerInfo(ctx.ServerId)
	if err != nil {
		return fmt.Errorf("获取服务信息失败: %w", err)
	}

	ctx.ServerInfo = &serverInfo
	ctx.NodeStatFactory = entity.NewNodeStatFactory(&entity.NodeStat{
		ServerName: serverInfo.ServerName,
		ServerId:   serverInfo.ID,
	})

	logger.Server.Infof("Activate | serverInfo=%+v", serverInfo)
	return nil
}

// setupPanicRecovery 设置panic恢复机制
func (ctx *ActivateContext) setupPanicRecovery() {
	if r := recover(); r != nil {
		errorMsg := fmt.Sprintf("激活服务失败 %s | 主机节点 %d | 版本号 %d | 原因 %v",
			ctx.ServerInfo.ServerName, ctx.LocalNodeId, ctx.PackageId, r)
		
		mapper.T_Mapper.SaveNodeStat(ctx.NodeStatFactory.Assign(&entity.NodeStat{
			TYPE:    entity.TYPE_ERROR,
			Content: errorMsg,
		}))
		logger.Active.Errorf("Activate | panic recovered: %v", r)
	}
}

// fetchConfigFiles 拉取配置文件
func (ctx *ActivateContext) fetchConfigFiles() error {
	err := api.GetConfigList(api.GetConfigListReq{
		ServerId: ctx.ServerId,
	})
	if err != nil {
		logger.Active.Errorf("Activate | GetConfigList | error=%s", err.Error())
		return fmt.Errorf("获取配置文件失败: %w", err)
	}
	return nil
}

// handleDeployment 处理部署逻辑
func (ctx *ActivateContext) handleDeployment() error {
	// 获取包信息
	packageInfo, err := mapper.T_Mapper.GetServerPackageInfo(ctx.PackageId)
	if err != nil {
		logger.Active.Errorf("Activate | GetServerPackageInfo | error=%s", err.Error())
		return fmt.Errorf("获取包信息失败: %w", err)
	}

	logger.Active.Infof("Activate | packageInfo=%+v", packageInfo)

	// 下载包文件（如果不存在）
	if err := ctx.downloadPackageIfNeeded(packageInfo.FileName); err != nil {
		logger.Active.Errorf("Activate | downloadPackageIfNeeded | error=%s", err.Error())
		return fmt.Errorf("下载包失败: %w", err)
	}

	// 解压包到目标目录
	if err := ctx.extractPackage(packageInfo.FileName); err != nil {
		logger.Active.Errorf("Activate | extractPackage | error=%s", err.Error())
		return fmt.Errorf("解压包失败: %w", err)
	}

	return nil
}

// downloadPackageIfNeeded 如果需要则下载包
func (ctx *ActivateContext) downloadPackageIfNeeded(packageFileName string) error {
	tarPath := filepath.Join(constant.TARGET_PACKAGE_DIR, ctx.ServerInfo.ServerName, packageFileName)
	logger.Active.Infof("Activate | downloadPackageIfNeeded | tarPath=%s", tarPath)
	if _, err := os.Stat(tarPath); os.IsNotExist(err) {
		return api.GetFile(api.FileReq{
			FileName: packageFileName,
			Type:     constant.FILE_TYPE_PACKAGE,
			ServerId: ctx.ServerId,
		})
	}
	return nil
}

// extractPackage 解压包
func (ctx *ActivateContext) extractPackage(packageFileName string) error {
	tarPath := filepath.Join(constant.TARGET_PACKAGE_DIR, ctx.ServerInfo.ServerName, packageFileName)
	serverDir := filepath.Join(constant.TARGET_SERVANT_DIR, ctx.ServerInfo.ServerName)
	
	logger.Active.Infof("Activate | extracting tarPath=%s to serverDir=%s", tarPath, serverDir)
	return patchutils.T_PatchUtils.Tar2Dest(tarPath, serverDir)
}

// deactivateExistingNodes 停止现有节点
func (ctx *ActivateContext) deactivateExistingNodes() (int32, string) {
	deactivateCode, deactivateMsg := Deactivate(&protocol.ActivateReq{
		ServerId:      ctx.Req.ServerId,
		ServerNodeIds: ctx.Req.ServerNodeIds,
	})

	if deactivateCode != CODE_SUCCESS {
		errorMsg := fmt.Sprintf("停止服务失败 %s | 节点 %v | 版本号 %d | 原因 %s",
			ctx.ServerInfo.ServerName, ctx.Req.ServerNodeIds, ctx.PackageId, deactivateMsg)
		logger.Active.Errorf("Activate | deactivateExistingNodes | errorMsg=%s", errorMsg)
		
		mapper.T_Mapper.SaveNodeStat(ctx.NodeStatFactory.Assign(&entity.NodeStat{
			TYPE:    entity.TYPE_ERROR,
			Content: errorMsg,
		}))
	}

	return deactivateCode, deactivateMsg
}

// activateNodes 激活节点
func (ctx *ActivateContext) activateNodes() error {
	nodes, err := mapper.T_Mapper.GetServerNodes(ctx.ServerId, ctx.LocalNodeId)
	if err != nil {
		logger.Active.Errorf("Activate | GetServerNodes | error=%s", err.Error())
		return fmt.Errorf("获取服务节点失败: %w", err)
	}

	logger.Active.Infof("Activate | nodes=%+v", nodes)

	for _, node := range nodes {
		if !patchutils.T_PatchUtils.Contains(ctx.ServerNodeIds, node.Id) {
			continue
		}

		if err := ctx.activateSingleNode(&node); err != nil {
			logger.Active.Errorf("Activate | activateSingleNode | error=%s", err.Error())
			return fmt.Errorf("激活节点 %d 失败: %w", node.Id, err)
		}
	}

	return nil
}

// activateSingleNode 激活单个节点
func (ctx *ActivateContext) activateSingleNode(node *mapper.ServerNodesVo) error {
	// 记录开始状态
	ctx.logNodeStartStatus(node)

	// 创建服务器信息
	patchServerInfo := ctx.createServerInfo(node)
	logger.Active.Infof("Activate | patchServerInfo=%+v", patchServerInfo)

	// 创建命令
	cmd, err := patchServerInfo.CreateCommand()
	if err != nil {
		return fmt.Errorf("创建命令失败: %w", err)
	}

	logger.Active.Infof("Activate | command args=%v", cmd.GetCmd().Args)

	// 启动服务
	if err := cmd.Start(); err != nil {
		logger.Active.Errorf("Activate | Start | error=%s", err.Error())
		ctx.logNodeErrorStatus(node, "启动服务器失败", err)
		return fmt.Errorf("启动服务失败: %w", err)
	}

	// 添加到命令管理器
	command.CenterManager.AddCommand(node.Id, cmd)
	
	// 设置cgroup,安卓就不设置了
	if config.Conf.GetOs() != "android" {
		if err := command.UseCgroup(cmd); err != nil {
			ctx.logNodeErrorStatus(node, "设置cgroup失败", err)
			return fmt.Errorf("设置cgroup失败: %w", err)
		}
	}

	logger.Active.Infof("Activate | node %d activated successfully", node.Id)
	return nil
}

// createServerInfo 创建服务器信息
func (ctx *ActivateContext) createServerInfo(node *mapper.ServerNodesVo) *command.ServerInfo {
	targetFile := filepath.Join(ctx.Cwd, constant.TARGET_SERVANT_DIR, 
		ctx.ServerInfo.ServerName, ctx.ServerInfo.ExecFilePath)

	return &command.ServerInfo{
		ServerType:     ctx.ServerInfo.ServerType,
		ServerName:     ctx.ServerInfo.ServerName,
		TargetFile:     targetFile,
		BindPort:       node.Port,
		BindHost:       node.Host,
		NodeId:         node.Id,
		AdditionalArgs: node.AdditionalArgs,
		ServerRunType:  node.ServerRunType,
		ServerId:       ctx.ServerId,
		DockerName:     ctx.ServerInfo.DockerName,
	}
}

// updateNodePatch 更新节点版本
func (ctx *ActivateContext) updateNodePatch() error {
	return mapper.T_Mapper.UpdateNodePatch(ctx.ServerNodeIds, ctx.PackageId)
}

// logNodeStartStatus 记录节点开始状态
func (ctx *ActivateContext) logNodeStartStatus(node *mapper.ServerNodesVo) {
	content := fmt.Sprintf("开始部署服务器 %s | 节点 %d | 版本号 %d | 端口号 %d",
		ctx.ServerInfo.ServerName, node.Id, ctx.PackageId, node.Port)
	
	mapper.T_Mapper.SaveNodeStat(ctx.NodeStatFactory.Assign(&entity.NodeStat{
		TYPE:         entity.TYPE_INFO,
		ServerNodeId: node.Id,
		Content:      content,
	}))
}

// logNodeErrorStatus 记录节点错误状态
func (ctx *ActivateContext) logNodeErrorStatus(node *mapper.ServerNodesVo, operation string, err error) {
	content := fmt.Sprintf("%s %s | 节点 %d | 版本号 %d | 端口号 %d | 原因 %s",
		operation, ctx.ServerInfo.ServerName, node.Id, ctx.PackageId, node.Port, err.Error())
	
	mapper.T_Mapper.SaveNodeStat(ctx.NodeStatFactory.Assign(&entity.NodeStat{
		TYPE:         entity.TYPE_ERROR,
		ServerNodeId: node.Id,
		Content:      content,
	}))
}

// logSuccessStatus 记录成功状态
func (ctx *ActivateContext) logSuccessStatus() {
	content := fmt.Sprintf("部署服务器成功 %s | 节点 %v | 版本号 %d",
		ctx.ServerInfo.ServerName, ctx.ServerNodeIds, ctx.PackageId)
	
	mapper.T_Mapper.SaveNodeStat(ctx.NodeStatFactory.Assign(&entity.NodeStat{
		TYPE:    entity.TYPE_SUCCESS,
		Content: content,
	}))
}
