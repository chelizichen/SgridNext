package main

// GOOS=linux GOARCH=amd64 go build -o $ServerName
import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"sgridnext.com/distributed"
	protocol "sgridnext.com/test/rpcserver/grpc_go_server/proto"
)

type GreetServicePrx struct{}

func (g *GreetServicePrx) GetAddrs() []*distributed.BaseSvrNodeStat {
	// 从注册中心获取节点列表
	registry := distributed.Registry{}
	nodes, err := registry.FindRegistryByServerName(g.GetServerName())
	if err != nil {
		fmt.Println("获取节点列表失败")
		return nil
	}
	// 转换为[]*command.BaseSvrNodeStat
	addrs := make([]*distributed.BaseSvrNodeStat, 0)
	for _, node := range nodes {
		addr := &distributed.BaseSvrNodeStat{
			ServerName: node.ServerName,
			ServerHost: node.ServerHost,
			ServerPort: node.ServerPort,
		}
		addrs = append(addrs, addr)
	}
	fmt.Println("获取节点列表成功", addrs)
	return addrs
}

func (g *GreetServicePrx) NewClient(conn grpc.ClientConnInterface) *protocol.GreetServiceClient {
	client := protocol.NewGreetServiceClient(conn)
	return &client
}

func (g *GreetServicePrx)GetServerName() string {
	return "SgridTestGrpcGoServer"
}


func LoadProxy() *distributed.PrxManage[*protocol.GreetServiceClient] {
	var prx = &GreetServicePrx{}
	pm, err := distributed.LoadStringToProxy(prx)
	if err != nil {
		fmt.Println("加载代理失败 | err: ", err)
		return nil
	}
	return pm
}

func main() {
	port := os.Getenv("SGRID_TARGET_PORT")
	fmt.Println("SGRID_TARGET_PORT: ", port)
	if port == "" {
		fmt.Println("SGRID_TARGET_PORT is empty")
		port = "12051"
	}
	host := os.Getenv("SGRID_TARGET_HOST")
	fmt.Println("SGRID_TARGET_HOST: ", port)
	if host == "" {
		fmt.Println("SGRID_TARGET_HOST is empty")
		host = "0.0.0.0"
	}

	engine := gin.Default()
	prx := LoadProxy()

	// time.Sleep(time.Second * 30)
	// client := *prx.GetClient()
	// rsp, error := client.SayHello(context.Background(), &protocol.SayHelloReq{
	// 	Name: "grpc_go_client",
	// })
	// if error != nil {
	// 	fmt.Println("调用远程方法失败")
	// 	return
	// }
	// fmt.Println("调用远程方法成功")
	// fmt.Println("rsp: ", rsp.Message)

	engine.GET("/", func(c *gin.Context) {
		client,ok := prx.GetClient()
		if !ok {
			fmt.Println("获取客户端失败")
			c.JSON(200, gin.H{
				"message": "获取客户端失败",
			})
			return
		}
		rsp, error := (*client).SayHello(context.Background(), &protocol.SayHelloReq{
			Name: "grpc_go_client",
		})
		if error != nil {
			fmt.Println("调用远程方法失败")
			c.JSON(200, gin.H{
				"message": "调用远程方法失败",
			})
			return
		}
		fmt.Println("调用远程方法成功")
		fmt.Println("rsp: ", rsp.Message)
		c.JSON(200, gin.H{
			"message": rsp.Message,
		})
	})
	bind_addr := host + ":" + port
	engine.Run(bind_addr)
}
