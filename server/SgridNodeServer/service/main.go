package service

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/types/known/emptypb"
	"sgridnext.com/server/SgridNodeServer/api"
	"sgridnext.com/server/SgridNodeServer/command"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/server/SgridNodeServer/state"
	"sgridnext.com/server/SgridNodeServer/util"
	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
	"sgridnext.com/src/resource"
)

const (
	MSG_SUCCESS  = "请求成功"
	CODE_SUCCESS = 0
	MSG_FAIL     = "请求失败"
	CODE_FAIL    = -1
)



type NodeServer struct {
	protocol.UnimplementedNodeServantServer
}

func (n *NodeServer) Probe(ctx context.Context, req *protocol.ProbeReq) (*protocol.BasicRes, error) {
	if req.Type == 1 {
		logger.Probe.Info("探针请求，不做任何操作")
		// 探针请求，不做任何操作
		return &protocol.BasicRes{
			Code: CODE_SUCCESS,
			Msg: MSG_SUCCESS,
		}, nil
	}
	logger.Probe.Info("探针请求，写入配置文件" + req.String())
	defer state.NodeServerState.Store(state.NODE_STATE_ONLINE)
	// 覆盖配置文件
	confStr := req.Conf
	var confObj constant.ConfObj

	err := json.Unmarshal([]byte(confStr), &confObj)
	if err != nil {
		logger.Probe.Errorf("解析配置文件失败: %v", err)
		return &protocol.BasicRes{
			Code: CODE_FAIL,
			Msg:  MSG_FAIL + " | " + err.Error(),
		}, nil
	}
	config.Conf.Set("db", confObj.Db)
	config.Conf.Set("dbtype", confObj.DbType)
	config.Conf.Set("nodeIndex", confObj.NodeIndex)
	config.Conf.Set("mainNode", confObj.MainNode)
	config.Conf.Set("host", confObj.Host)
	config.Conf.Set("nodeStatus", state.NODE_STATE_ONLINE)
	logger.Probe.Info("探针请求完成，节点状态设置为在线")
	return &protocol.BasicRes{
		Code: CODE_SUCCESS,
		Msg: MSG_SUCCESS,
	}, nil
}

func (n *NodeServer) KeepAlive(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if state.NodeServerState.Load() == state.NODE_STATE_STAYBY {
		logger.App.Info("节点状态为stayby，不进行心跳")
		return &emptypb.Empty{}, nil
	}
	logger.Alive.Info("alive called")
	err := mapper.T_Mapper.UpdateNodeUpdateTime(config.Conf.GetLocalNodeId())
	if err != nil {
		logger.App.Errorf("更新节点更新时间失败: %v", err)
	}
	logger.App.Infof("更新节点更新时间成功: %v", config.Conf.GetLocalNodeId())
	now := constant.GetCurrentTime()
	command.CenterManager.SyncStat(now)
	return &emptypb.Empty{}, nil
}

func (s *NodeServer) GetNodeStat(ctx context.Context, in *emptypb.Empty) (*protocol.GetNodeStatRsp, error) {
	if state.NodeServerState.Load() == state.NODE_STATE_STAYBY {
		logger.App.Info("节点状态为stayby，不进行心跳")
		return &protocol.GetNodeStatRsp{
			Code: CODE_SUCCESS,
			Msg: MSG_SUCCESS,
		}, nil
	}
	cwd, _ := os.Getwd()
	stat_path := filepath.Join(cwd, "stat.json")
	jsonStr, err := os.ReadFile(stat_path)
	if err != nil {
		logger.App.Errorf("读取stat.json文件失败: %v", err)
		return &protocol.GetNodeStatRsp{
			Code: CODE_FAIL,
			Msg:  MSG_FAIL,
		}, nil
	}
	return &protocol.GetNodeStatRsp{
		Code: CODE_SUCCESS,
		Msg:  MSG_SUCCESS,
		Data: string(jsonStr),
	}, nil
}

func (s *NodeServer) SyncAllNodeStat(ctx context.Context, in *protocol.SyncStatReq) (*protocol.BasicRes, error) {
	if state.NodeServerState.Load() == state.NODE_STATE_STAYBY {
		logger.App.Info("节点状态为stayby，不进行心跳")
		return &protocol.BasicRes{
			Code: CODE_SUCCESS,
			Msg: MSG_SUCCESS,
		}, nil
	}	
	cwd, _ := os.Getwd()
	stat_remote_path := filepath.Join(cwd, "stat-remote.json")
	outFile, err := os.OpenFile(stat_remote_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logger.App.Errorf("创建文件失败: SyncStat |%v", err)
		return &protocol.BasicRes{
			Code: CODE_FAIL,
			Msg:  MSG_FAIL,
		}, nil
	}
	defer outFile.Close()
	if _, err := outFile.Write([]byte(in.Data)); err != nil {
		logger.App.Errorf("文件写入失败: SyncStat | %v", err)
		return &protocol.BasicRes{
			Code: CODE_FAIL,
			Msg:  MSG_FAIL,
		}, nil
	}
	return &protocol.BasicRes{
		Code: CODE_SUCCESS,
		Msg:  MSG_SUCCESS,
	}, nil
}

func (s *NodeServer) ActivateServant(ctx context.Context, in *protocol.ActivateReq) (*protocol.BasicRes, error) {
	logger.App.Info("服务激活 %v ", in.String())
	code, msg := Activate(in)
	return &protocol.BasicRes{
		Code: int32(code),
		Msg:  msg,
	}, nil
}

func (s *NodeServer) DeactivateServant(ctx context.Context, in *protocol.ActivateReq) (*protocol.BasicRes, error) {
	logger.App.Info("服务关闭 %v", in.String())
	code, msg := Deactivate(in)
	return &protocol.BasicRes{
		Code: int32(code),
		Msg:  msg,
	}, nil
}

func (s *NodeServer) SyncConfigFile(ctx context.Context, in *protocol.SyncReq) (*protocol.BasicRes, error) {
	logger.App.Info("配置同步 %v", in.String())
	err := api.GetFile(api.FileReq{
		FileName: in.FileName,
		ServerId: int(in.ServerId),
		Type:     int(in.Type),
	})
	if err != nil {
		return &protocol.BasicRes{
			Code: CODE_FAIL,
			Msg:  err.Error(),
		}, nil
	}
	return &protocol.BasicRes{
		Code: CODE_SUCCESS,
		Msg:  "下载成功",
	}, nil
}

func (s *NodeServer) SyncServicePackage(ctx context.Context, in *protocol.SyncReq) (*protocol.BasicRes, error) {
	logger.App.Info("服务包同步 %v", in.String())
	err := api.GetFile(api.FileReq{
		FileName: in.FileName,
		ServerId: int(in.ServerId),
		Type:     int(in.Type),
	})
	if err != nil {
		return &protocol.BasicRes{
			Code: CODE_FAIL,
			Msg:  err.Error(),
		}, nil
	}
	return &protocol.BasicRes{
		Code: CODE_SUCCESS,
		Msg:  "下载成功",
	}, nil
}

func (s *NodeServer) CgroupLimit(ctx context.Context, in *protocol.CgroupLimitReq) (*protocol.BasicRes, error) {
	logger.App.Info("设置CGroup %v ", in.String())
	code, msg := CgroupLimit(in)
	return &protocol.BasicRes{
		Code: int32(code),
		Msg:  msg,
	}, nil
}

func (s *NodeServer) CheckStat(ctx context.Context, in *protocol.CheckStatReq) (*protocol.BasicRes, error) {
	logger.App.Info("获取CGroup %v ", in.String())
	code, msg := CheckStat(in)
	return &protocol.BasicRes{
		Code: int32(code),
		Msg:  msg,
	}, nil
}

func (s *NodeServer) DownloadFile(in *protocol.DownloadFileRequest, stream protocol.NodeServant_DownloadFileServer) error {
	logger.App.Info("下载文件 %v ", in.String())
	if in.Type == constant.FILE_TYPE_LOG {
		cwd, _ := os.Getwd()
		serverInfo, err := mapper.T_Mapper.GetServerInfo(int(in.ServerId))
		if err != nil {
			return err
		}
		serverName := serverInfo.ServerName
		fileName := in.FileName
		filePath := filepath.Join(cwd, constant.TARGET_LOG_DIR, serverName, fileName)
		if serverInfo.LogPath != "" {
			filePath = filepath.Join(serverInfo.LogPath,in.FileName)
		}
		logger.App.Info("下载文件路径 %s ", filePath)
		file, err := os.Open(filePath)
		if err != nil {
			logger.App.Errorf("打开文件失败 %s ", err.Error())
			return err
		}
		stat, err := os.Stat(filePath)
		if err != nil {
			logger.App.Errorf("获取文件大小失败 %s ", err.Error())
			return err
		}
		fileSize := stat.Size()
		reader := bufio.NewReader(file)
		buffer := make([]byte, 1024)
		for {
			n, err := reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					logger.App.Info("文件发送完成")
					break
				}
				logger.App.Errorf("读取文件失败 %s ", err.Error())
				return err
			}
			stream.Send(&protocol.DownloadFileResponse{
				Code:  CODE_SUCCESS,
				Msg:   MSG_SUCCESS,
				Data:  buffer[:n],
				IsEnd: false,
			})
			logger.App.Info("发送文件 %d 字节", n)
			if fileSize <= int64(n) {
				logger.App.Info("文件发送完成")
				break
			}
		}
		stream.Send(&protocol.DownloadFileResponse{
			Code:  CODE_SUCCESS,
			Msg:   MSG_SUCCESS,
			Data:  nil,
			IsEnd: true,
		})
		return nil
	}

	return nil
}

// 可以指定LogPath，兼容旧服务进行日志查询
// 修改 GetFileList 方法
func (s *NodeServer) GetFileList(ctx context.Context, in *protocol.GetFileListReq) (*protocol.GetFileListResponse, error) {
	logger.App.Info("获取文件列表 %v ", in.String())
	if in.Type == constant.FILE_TYPE_LOG {
		// 根据日志类型决定读取哪个目录
		var logDir string
		
		if in.LogCategory == constant.LOG_TYPE_NODE {
			// 节点日志：读取当前工作目录下的logs目录
			cwd, _ := os.Getwd()
			logDir = filepath.Join(cwd, "logs")
		} else {
			// 业务日志：使用原有逻辑
			serverInfo, err := mapper.T_Mapper.GetServerInfo(int(in.ServerId))
			if err != nil {
				return nil, err
			}
			logDir = util.GetLogDir(&serverInfo)
		}
		
		files, err := os.ReadDir(logDir)
		if err != nil {
			logger.App.Errorf("读取目录失败 %s ", err.Error())
			return nil, err
		}
		fileList := make([]string, 0)
		for _, file := range files {
			if !file.IsDir() {
				fileList = append(fileList, file.Name())
			}
		}
		logger.App.Info("获取到文件列表 %v", fileList)
		return &protocol.GetFileListResponse{
			Code:     CODE_SUCCESS,
			Msg:      MSG_SUCCESS,
			FileList: fileList,
		}, nil
	}
	return nil, nil
}

// 修改 GetLog 方法
func (s *NodeServer) GetLog(ctx context.Context, in *protocol.GetLogReq) (*protocol.GetLogRes, error) {
	logger.App.Info("获取日志 %v ", in.String())
	
	var logPath string
	
	if in.LogCategory == constant.LOG_TYPE_NODE {
		// 节点日志：直接从当前工作目录下的logs目录读取
		cwd, _ := os.Getwd()
		logPath = filepath.Join(cwd, "logs", in.FileName)
	} else {
		// 业务日志：使用原有逻辑
		serverInfo, err := mapper.T_Mapper.GetServerInfo(int(in.ServerId))
		if err != nil {
			return nil, err
		}
		logPath = util.GetLogPath(&serverInfo, in.FileName)
	}
	
	rsp, err := QueryLog(logPath, in.LogType, in.Keyword, in.Len)
	if err != nil {
		return nil, err
	}
	return &protocol.GetLogRes{
		Code: CODE_SUCCESS,
		Msg:  MSG_SUCCESS,
		Data: rsp,
	},nil
}

func (s *NodeServer) GetNodeResource(ctx context.Context, in *emptypb.Empty) (*protocol.GetNodeResourceRes, error) {
	logger.App.Info("获取节点资源信息")
	cwd, _ := os.Getwd()
	nodeResource := resource.GetNodeResource(cwd)
	jsonStr, err := json.Marshal(nodeResource)
	if err != nil {
		logger.App.Errorf("序列化节点资源信息失败: %v", err)
		return nil, err
	}
	logger.App.Infof("节点资源信息: %s", string(jsonStr))
	return &protocol.GetNodeResourceRes{
		Code: CODE_SUCCESS,
		Msg: MSG_SUCCESS,
		Data: string(jsonStr),
	}, nil
}