# SgridNext

## 简介

极简的运维部署平台，借鉴 Tars的部署功能 ，重构旧版 Sgrid平台，同时支持 Cgroup 资源限制

该平台的首要目的是为了搭建一个测试环境，同时摆脱框架依赖，降低框架耦合度

特性：

1. 多节点管理与部署
2. 提供服务注册与发现
3. 提供远程配置中心
4. 提供**Cgroup** 的资源限制
5. 支持Java、Node、以及二进制文件的部署与版本管理，也可以通过 sh 脚本对docker 进行更新管理
6. **探针管理**，通过在主控直接添加探针配置（**networkPrefixs**），可以自动搜寻在线的 sgridnode 服务，无需重复编写配置文件

该平台采用 **主控-节点** 分离架构，确保业务服务不受管理服务（主控/节点）故障影响：

- **主控服务（Master Service）**
  - 仅负责 **服务发布、配置管理、Web 管理界面**。
  - 故障影响：**无法发布新服务，但现有业务不受影响**。
  - 心跳检测中断，但节点仍可自主管理业务服务。
- **节点服务（Agent Service）**
  - 负责 **业务服务的启停、状态监控**。
  - 故障影响：**无法远程管理业务服务，但业务进程仍正常运行**。
- **业务服务（Business Service）**
  - 由节点服务管理，但 **独立存活**，即使节点服务崩溃也不受影响。

为确保节点服务停止时 **不影响业务服务**，采用 **Cgroup 逃逸** 技术：

### **（1）默认行为（问题）**

- `systemctl start sgridnode` 启动节点服务时，系统自动创建 cgroup（如 `/sys/fs/cgroup/system.slice/sgridnode.service/`）。
- 若节点服务启动业务进程（如 `business-service`），默认情况下，业务进程会继承节点服务的 cgroup。
- 当 `systemctl stop sgridnode` 时，系统会向该 cgroup 内所有进程（包括业务进程）发送 `SIGTERM`/`SIGKILL`，导致业务服务被误杀。

### **（2）解决方案：Cgroup 逃逸**

- **在业务进程启动后，立即将其移至独立 cgroup, 并提供资源限制功能更**（如 `/sys/fs/cgroup/business.slice/`）。
- 这样，当节点服务停止时，系统仅清理 `sgridnode.service` cgroup，而业务进程因属于其他 cgroup 不受影响。

## 编译部署

执行 ./deploy.sh 脚本进行编译部署

### SgridNext 部署

#### Systemctl 启动

创建指定systemctl 启动文件 **/usr/lib/systemd/system/sgridnext.service**

```toml
[Unit]
Description = sgrid next,A cloud platform for grid computing

[Service]
Type = simple
ExecStart = /usr/sgridnext/sgridnext
WorkingDirectory = /usr/sgridnext
Environment=PATH=/usr/bin:/usr/local/bin
Restart = no
```

ExecStart 为启动文件，为用户自己部署的目录下的启动文件
WorkingDirectory 为工作目录，用户部署目录，自行控制
Environment 为环境变量

-- 具体配置参考官方文档
配置完之后 执行 **systemctl start sgridnext**

```shell
# 启动
systemctl start sgridnext
# 重启
systemctl restart sgridnext
# 停止
systemctl stop sgridnext
# 查看状态
systemctl status sgridnext
```

#### 配置文件示例 config.json

````json
{
    "db": "host=10.111.112.113 port=5432 user=admin password=123456 dbname=sgrid_next sslmode=disable",
    "dbtype": "postgres",
    "host": "10.111.112.113",
    "httpPort": "15872",
    "networkPrefixs": "10.111.112.114,10.111.112.115,10.111.112.116",
    "nodeIndex": "1",
    "plantform": {
        "account": "admin",
        "password": "admin@sgridnext"
    }
}
````



### SgridNode 部署

1. 进入到 server/SgridNodeServer 目录下
2. 编写配置文件 config.json ，首次只需要填写 host 即可，后续通过探针请求来进行配置文件更新
3. 将 sgridnode 文件 拷贝至 **/usr/sgridnode/** 目录下
4. 编写 systemctl 启动文件 **/usr/lib/systemd/system/sgridnode.service**

```s
[Unit]
Description = sgrid next,A cloud platform for grid computing

[Service]
Type = simple
ExecStart = /usr/sgridnode/sgridnode
WorkingDirectory = /usr/sgridnode
Environment=PATH=/usr/bin:/usr/local/bin
Restart = no
```



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

```go
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
```

## Golang-Gin-Proxy 服务

如果作为 Proxy代理服务需要调用 其他 GRPC服务，需要实现 `type Proxy[T any] interface` 接口

接口定义如下

```go

type Proxy[T any] interface {
 // 获取服务注册表地址，进行代理链接
 GetAddrs() []*BaseSvrNodeStat
 // 客户端实例函数
 NewClient(conn grpc.ClientConnInterface) T
 // 服务名称
 GetServerName() string
}

```

Proxy代理是如何做到的？

sgridnode 含有服务节点维护的功能，在服务启动时，会将服务信息同步到 /usr/sgridnode/stat.json 中
sgridnext 主控 会读取所有的 节点服务状态，再同步给所有的 sgridnode节点到 /usr/sgridnode/stat-remote.json 中

GetAddrs 在初始化时执行一次，随后每30s执行一次，进行代理节点的更新与删除

接口基本实现示例如下

```go
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
```

grpc-client-proxy 调用远程rpc 方法示例

```go
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
```

## Golang-Grpc 服务

与Gin服务类似，生产的 PORT 和 HOST 以环境变量的行式传入

```go
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
```

## Docker-Java 服务

如果是 springboot 服务打包成 docker 镜像部署到测试环境的开发方式，可以使用一套shell脚本进行测试环境的快速更新，原理是通过更换镜像内的jar包进行快速更新

1. 编写 docker 更新命令 (update.sh)

```sh
#! /bin/bash

echo  "开始执行 docker 切换jar包"
docker cp ./SpringBootServer.jar docker-id:/app/SpringBootServer.jar
echo "开始执行 docker 重启服务"
docker restart docker-id

echo "部署完成"
```

2. 编写 jar 包打包命令 (build.sh)

````sh
echo "运行外部命令"
echo "Building SpringBootServer"

rm -r SpringBootServer.tar.gz

tar -czf SpringBootServer.tar.gz ./SpringBootServer.jar ./update.sh

echo "构建完成"
````
