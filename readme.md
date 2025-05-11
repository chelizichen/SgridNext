# SgridNext

## 简介

极简的运维部署平台，借鉴 Tars的部署功能 ，重构旧版 Sgrid平台，同时支持 Cgroup 资源限制

该平台的首要目的是为了搭建一个测试环境，同时摆脱框架依赖，降低框架耦合度

## 编译部署

### 后台

执行 ./build.sh 编译后台程序

### 前端

在 web 目录下打包即可

执行 npm run build

### Systemctl 启动

创建指定systemctl 启动文件 /usr/lib/systemd/system/sgridnext.service
ExecStart 为启动文件
WorkingDirectory 为工作目录
Environment 为环境变量

-- 具体配置参考官方文档

配置完之后 执行

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

### 创建节点

1. 节点内网地址 通过 ip show addr etho0 查看
2. 节点外网地址 通过 curl ifconfig.me 查看

### TIPS

1. 纯Express服务，启动需要 20M的内存
2. 纯Gin服务，启动需要8M的内存

可以通过

cat /sys/fs/cgroup/system.slice/xx/memory.current 查看使用内存大小
