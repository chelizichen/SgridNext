package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	protocol "sgridnext.com/test/rpcserver/grpc_go_server/proto"
)

type GreetService struct {
	protocol.UnimplementedGreetServiceServer
}

func (s *GreetService) SayHello(ctx context.Context, in *protocol.SayHelloReq) (*protocol.SayHelloRes, error) {
	fmt.Println("SayHello: ", in.GetName())
	res := &protocol.SayHelloRes{
		Message: "Hello " + in.GetName(),
	}
	return res, nil
}

func main() {
	port := os.Getenv("SGRID_TARGET_PORT")
	fmt.Println("SGRID_TARGET_PORT: ", port)
	if port == "" {
		fmt.Println("SGRID_TARGET_PORT is empty")
		port = "10010"
	}
	host := os.Getenv("SGRID_TARGET_HOST")
	fmt.Println("SGRID_TARGET_HOST: ", port)
	if host == "" {
		fmt.Println("SGRID_TARGET_HOST is empty")
		host = "0.0.0.0"
	}
	BIND_ADDR := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", BIND_ADDR)
	if err != nil {
		fmt.Println("监听失败: ", err)
	}
	var opts []grpc.ServerOption
	opts = append(opts,
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    5 * time.Second,
			Timeout: 1 * time.Second,
		}),
	)
	srv := grpc.NewServer(opts...)
	protocol.RegisterGreetServiceServer(srv, &GreetService{})

	fmt.Println("节点服务启动在 :" + BIND_ADDR)
	if err := srv.Serve(lis); err != nil {
		fmt.Println("服务启动失败: ", err)
	}
}
