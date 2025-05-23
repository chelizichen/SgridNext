// 节点服务
// 1. 需要保证与主节点网络互通
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
	"sgridnext.com/server/SgridNodeServer/api"
	"sgridnext.com/server/SgridNodeServer/command"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/server/SgridNodeServer/service"
	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/db"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
)

type NodeServer struct {
	protocol.UnimplementedNodeServantServer
}

func (n *NodeServer) KeepAlive(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	logger.Alive.Info("alive called")
	now := constant.GetCurrentTime()
	command.CenterManager.SyncStat(now)
	return &emptypb.Empty{}, nil
}

func (s *NodeServer) GetNodeStat(ctx context.Context, in *emptypb.Empty) (*protocol.GetNodeStatRsp, error){
	cwd, _ := os.Getwd()
	stat_path := filepath.Join(cwd, "stat.json")
	jsonStr, err := os.ReadFile(stat_path)
	if err!= nil {
		logger.App.Errorf("读取stat.json文件失败: %v", err)
		return &protocol.GetNodeStatRsp{
			Code: service.CODE_FAIL,
			Msg: service.MSG_FAIL,
		}, nil
	}
	return &protocol.GetNodeStatRsp{
		Code: service.CODE_SUCCESS,
		Msg:  service.MSG_SUCCESS,
		Data: string(jsonStr),
	}, nil
}

func (s *NodeServer) SyncAllNodeStat(ctx context.Context, in *protocol.SyncStatReq) (*protocol.BasicRes, error) {
	cwd, _ := os.Getwd()
	stat_remote_path := filepath.Join(cwd, "stat-remote.json")
	outFile, err := os.OpenFile(stat_remote_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logger.App.Errorf("创建文件失败: SyncStat |%v", err)
		return &protocol.BasicRes{
			Code: service.CODE_FAIL,
			Msg: service.MSG_FAIL,
		}, nil
	}
	defer outFile.Close()
	if _, err := outFile.Write([]byte(in.Data)); err != nil {
		logger.App.Errorf("文件写入失败: SyncStat | %v", err)
		return &protocol.BasicRes{
			Code: service.CODE_FAIL,
			Msg: service.MSG_FAIL,
		}, nil
	}
	return &protocol.BasicRes{
		Code: service.CODE_SUCCESS,
		Msg:  service.MSG_SUCCESS,
	}, nil
}

func (s *NodeServer) ActivateServant(ctx context.Context, in *protocol.ActivateReq) (*protocol.BasicRes, error) {
	logger.App.Info("服务激活 %v ", in.String())
	code, msg := service.Acitvate(in)
	return &protocol.BasicRes{
		Code: int32(code),
		Msg:  msg,
	}, nil
}

func (s *NodeServer) DeactivateServant(ctx context.Context, in *protocol.ActivateReq) (*protocol.BasicRes, error) {
	logger.App.Info("服务关闭 %v", in.String())
	code, msg := service.Deactivate(in)
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
			Code: service.CODE_FAIL,
			Msg:  err.Error(),
		}, nil
	}
	return &protocol.BasicRes{
		Code: service.CODE_SUCCESS,
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
			Code: service.CODE_FAIL,
			Msg:  err.Error(),
		}, nil
	}
	return &protocol.BasicRes{
		Code: service.CODE_SUCCESS,
		Msg:  "下载成功",
	}, nil
}

func (s *NodeServer) CgroupLimit(ctx context.Context, in *protocol.CgroupLimitReq) (*protocol.BasicRes, error) {
	logger.App.Info("设置CGroup %v ", in.String())
	code, msg := service.CgroupLimit(in)
	return &protocol.BasicRes{
		Code: int32(code),
		Msg:  msg,
	}, nil
}

func (s *NodeServer) CheckStat(ctx context.Context, in *protocol.CheckStatReq) (*protocol.BasicRes, error) {
	logger.App.Info("获取CGroup %v ", in.String())
	code, msg := service.CheckStat(in)
	return &protocol.BasicRes{
		Code: int32(code),
		Msg:  msg,
	}, nil
}

func init() {
	conf := config.LoadConfig("./config.json")
	ormDb, err := db.InitDB(conf.Get("db"))
	if err != nil {
		panic(err)
	}
	mapper.LoadMapper(ormDb)
	snsList := command.LoadStatList()
	if snsList != nil{
		command.InitCommands(snsList.StatList)
	}
}

func main() {
	HOST, err := mapper.T_Mapper.GetHost(config.Conf.GetLocalNodeId())
	if err != nil {
		panic(fmt.Sprintf("获取主机信息失败 %v", err))
	}
	BIND_ADDR := fmt.Sprintf("%s:%s", HOST, constant.NODE_PORT)
	lis, err := net.Listen("tcp", BIND_ADDR)
	if err != nil {
		logger.App.Fatal("监听失败: ", err)
	}
	var opts []grpc.ServerOption
	opts = append(opts,
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    5 * time.Second,
			Timeout: 1 * time.Second,
		}),
	)
	srv := grpc.NewServer(opts...)
	protocol.RegisterNodeServantServer(srv, &NodeServer{})

	fmt.Println("节点服务启动在 :" + BIND_ADDR)
	if err := srv.Serve(lis); err != nil {
		logger.App.Fatal("服务启动失败: ", err)
	}
}
