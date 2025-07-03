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
	"sgridnext.com/server/SgridNodeServer/util"
	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/db"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
)

var BIND_ADDR = ""

func init() {
	config.LoadConfig("./config.json")
	ormDb, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	mapper.LoadMapper(ormDb)
	snsList := command.LoadStatList()
	if snsList != nil {
		command.InitCommands(snsList.StatList)
	}else{
		// 创建该文件
		command.InitStatFile()
	}
	defer schedule.LoadTick()
	HOST := util.GetHost()
	BIND_ADDR = fmt.Sprintf("%s:%s", HOST, constant.NODE_PORT)
}

func main() {
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
	protocol.RegisterNodeServantServer(srv, &service.NodeServer{})

	if err := srv.Serve(lis); err != nil {
		logger.App.Fatal("服务启动失败: ", err)
	} else {
		logger.App.Info("节点服务启动在 :" + BIND_ADDR)
	}
}
