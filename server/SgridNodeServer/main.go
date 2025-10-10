// 节点服务
// 1. 需要保证与主节点网络互通
package main

import (
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"sgridnext.com/server/SgridNodeServer/command"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/server/SgridNodeServer/schedule"
	"sgridnext.com/server/SgridNodeServer/service"
	"sgridnext.com/server/SgridNodeServer/state"
	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/db"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
)

var BIND_ADDR = ""

func doInit(){
	logger.App.Info("[SgridNode] 初始化完成 状态为online，开始执行 INIT 任务")
	// 轮询状态，如果是stayby 则 等待进入下一次循环
	// 如果为已链接 则不再执行，并且将配置文件覆盖
	config.LoadConfig("./config.json")
	ormDb, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	mapper.LoadMapper(ormDb)
	snsList := command.LoadStatList()
	if snsList != nil {
		command.InitCommands(snsList.StatList)
	} else {
		// 创建该文件
		command.InitStatFile()
	}
	schedule.LoadTick()
	logger.App.Info("[SgridNode] 初始化完成， INIT 任务执行完成，退出轮询")
	state.NodeServerState.Store(state.NODE_STATE_DONE_INIT)
}

func init() {
	// 先加载基础配置文件，避免读取到空配置
	conf := config.LoadConfig("./config.json")
	nodeStatus := conf.GetFloat64("nodeStatus")
	// 如果没被初始化过，则设置为stayby
	if nodeStatus == 0 {
		state.NodeServerState.Store(state.NODE_STATE_STAYBY)
	} else {
		state.NodeServerState.Store(int32(nodeStatus))
	}
	if state.NodeServerState.Load() == state.NODE_STATE_ONLINE {
		logger.App.Info("[SgridNode] 节点状态为online，直接进入初始化")
		doInit()
		return
	}else{
		// 依赖外部状态更改
		go func() {
			ticker := time.NewTicker(time.Minute * 1)
			defer ticker.Stop()
			for range ticker.C {
				if state.NodeServerState.Load() == state.NODE_STATE_STAYBY {
					logger.App.Info("[SgridNode] 初始化尚未完成 状态为stayby，等待进入下一次循环")
					continue
				} else if state.NodeServerState.Load() == state.NODE_STATE_ONLINE {
					doInit()
					break
				} else {
					logger.App.Info("[SgridNode] 初始化已完成 状态为done_init，退出轮询")
					break
				}
			}
		}()
	}

}

func main() {
	BIND_ADDR = fmt.Sprintf("%s:%s", config.Conf.Get("host"), constant.NODE_PORT)
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
		grpc.MaxRecvMsgSize(constant.MSG_RECV_SIZE_MAX),
		grpc.MaxSendMsgSize(constant.MSG_CALL_SIZE_MAX),
	)
	srv := grpc.NewServer(opts...)
	protocol.RegisterNodeServantServer(srv, &service.NodeServer{})

	if err := srv.Serve(lis); err != nil {
		logger.App.Fatal("服务启动失败: ", err)
	} else {
		logger.App.Info("节点服务启动在 :" + BIND_ADDR)
	}
}
