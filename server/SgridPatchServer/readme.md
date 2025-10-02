# 框架部署

## SgridNode服务部署

````js
服务名: sgridnode
创建端口: 25528
服务类型: Binary
执行地址: update
其他参数可自由设置
````

每台服务器都需要部署一个 sgridnode 节点，方便做运维和自升级

升级步骤

执行 ./deploy.sh（node 要 18）

然后将打好的 sgridnode.tar.gz 部署上去即可
