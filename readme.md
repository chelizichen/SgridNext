# SgridNext

## 简介

极简的运维部署平台，借鉴 Tars的部署功能 ，重构旧版 Sgrid平台，同时支持 Cgroup 资源限制

该平台的首要目的是为了搭建一个测试环境，同时摆脱框架依赖，降低框架耦合度

## 编译部署

### 后台

执行 ./build.sh 编译后台程序

### 前端

### SgridWeb 部署

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
2. 执行 ./build.sh 编译 Node 服务
3. 编写配置文件
   1. db 为主库地址
   2. nodeIndex 为 节点的 **ID**
   3. nodePort 为绑定的端口号 默认为 25528
   4. mainNode 为主节点地址，需要 pin 的通，不然没法拉取配置和包
4. 将 sgridnode 文件 拷贝至 **/usr/sgridnode/** 目录下
5. 编写 systemctl 启动文件 **/usr/lib/systemd/system/sgridnode.service**

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
