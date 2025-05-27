# SgridNext

## 简介

极简的运维部署平台，借鉴 Tars的部署功能 ，重构旧版 Sgrid平台，同时支持 Cgroup 资源限制

该平台的首要目的是为了搭建一个测试环境，同时摆脱框架依赖，降低框架耦合度

## 编译部署

执行 ./deploy.sh 脚本进行编译部署

### SgridNext 部署

#### Systemctl 启动

创建指定systemctl 启动文件 **/usr/lib/systemd/system/sgridnext.service**

````s
[Unit]
Description = sgrid next,A cloud platform for grid computing

[Service]
Type = simple
ExecStart = /usr/sgridnext/sgridnext
WorkingDirectory = /usr/sgridnext
Environment=PATH=/usr/bin:/usr/local/bin
Restart = no
````

ExecStart 为启动文件
WorkingDirectory 为工作目录
Environment 为环境变量

-- 具体配置参考官方文档
配置完之后 执行 **systemctl start sgridnext**

````shell
# 启动
systemctl start sgridnext
# 重启
systemctl restart sgridnext
# 停止
systemctl stop sgridnext
# 查看状态
systemctl status sgridnext
````

### SgridNode 部署

1. 进入到 server/SgridNodeServer 目录下
2. 编写配置文件
   1. db 为主库地址
   2. nodeIndex 为 节点的 **ID**
   3. nodePort 为绑定的端口号 默认为 25528
   4. mainNode 为主节点地址，需要 pin 的通，不然没法拉取配置和包
3. 将 sgridnode 文件 拷贝至 **/usr/sgridnode/** 目录下
4. 编写 systemctl 启动文件 **/usr/lib/systemd/system/sgridnode.service**

````s
[Unit]
Description = sgrid next,A cloud platform for grid computing

[Service]
Type = simple
ExecStart = /usr/sgridnode/sgridnode
WorkingDirectory = /usr/sgridnode
Environment=PATH=/usr/bin:/usr/local/bin
Restart = no
````

### 创建节点

1. 节点内网地址 通过 ip addr show etho0 查看
2. 节点外网地址 通过 curl ifconfig.me 查看

### TIPS

1. 纯Express服务，启动需要 20M的内存
2. 纯Gin服务，启动需要8M的内存
3. SpringBoot 服务，启动需要200M的内存

### 内存与CPU用量查看

cat /sys/fs/cgroup/system.slice/xx/memory.current 查看使用内存大小

## 服务部署注意事项

### Golang-Gin 服务

端口号 PORT 和 主机地址 HOST 在生产中会以 环境变量的行式传入

````go
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
````

## Golang-Gin-Proxy 服务

如果作为 Proxy代理服务需要调用 其他 GRPC服务，需要实现 `type Proxy[T any] interface` 接口

接口定义如下

````go

type Proxy[T any] interface {
 // 获取服务注册表地址，进行代理链接
 GetAddrs() []*BaseSvrNodeStat
 // 客户端实例函数
 NewClient(conn grpc.ClientConnInterface) T
 // 服务名称
 GetServerName() string
}

````

Proxy代理是如何做到的？

sgridnode 含有服务节点维护的功能，在服务启动时，会将服务信息同步到 /usr/sgridnode/stat.json 中
sgridnext 主控 会读取所有的 节点服务状态，再同步给所有的 sgridnode节点到 /usr/sgridnode/stat-remote.json 中

GetAddrs 在初始化时执行一次，随后每30s执行一次，进行代理节点的更新与删除

接口基本实现示例如下

````go
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
````

grpc-client-proxy 调用远程rpc 方法示例

````go
func LoadProxy() *distributed.PrxManage[*protocol.GreetServiceClient] {
 var prx = &GreetServicePrx{}
 pm, err := distributed.LoadStringToProxy(prx)
 if err != nil {
  fmt.Println("加载代理失败 | err: ", err)
  return nil
 }
 return pm
}

prx := LoadProxy()

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
````

## Golang-Grpc 服务

与Gin服务类似，生产的 PORT 和 HOST 以环境变量的行式传入

````go
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
````
