# 分布式系统

## Q: 一个服务如何通过名字服务找到对应的多台服务地址？

1. 要有服务发现的机制, 服务发现目前在 /usr/sgridnode/stat.json 保存了当前机器上服务节点的信息，但是无法获取到其他机器的信息
2. proxy 节点 30秒会更新一次 stat.json 文件,可以通过 proxy 节点获取到其他机器的信息，具体拉取所有节点信息步骤如下
   1. 当 Proxy 节点访问 当前节点时，返回当前节点的 stat 信息， 保存到 proxy 节点上，proxy节点存储所有sgridnode节点信息
   2. Proxy 节点向 所有 sgridnode 节点发送 rpc 请求，携带 所有节点的 stat 信息, sgridnode 节点收到请求后，将请求的 stat 信息保存到 stat-remote.json 文件中
   3. 业务服务 启动时，访问 stat-remote.json 文件，寻找名字节点
   4. 拉取成一个 数组，同时进行节点代理
   5. 请求时，随机选择一个节点进行请求，优先选择 健康的、离的近的节点
   6. 返回结果

stat-remote.json 示例

````json
{  
   "update_time":"2025-05-23 22:53:15",
   "stat_list":[
      {"node_id":5,"server_name":"SgridTestGoServer","pid":910613,"host":"10.0.12.17","port":16501,"machine_id":1,"server_id":3},
      {"node_id":3,"server_name":"SgridTestNodeServer","pid":910813,"host":"10.0.12.17","port":19421,"machine_id":1,"server_id":2},
      {"node_id":6,"server_name":"SgridTestGoServer","pid":31852,"host":"10.0.16.2","port":16501,"machine_id":2,"server_id":3},
      {"node_id":4,"server_name":"SgridTestNodeServer","pid":32069,"host":"10.0.16.2","port":19421,"machine_id":2,"server_id":2},
   ]
}
````

stat.json 示例

````json
{
   "update_time":"2025-05-23 22:54:15",
   "stat_list":[
      {"node_id":6,"server_name":"SgridTestGoServer","pid":31852,"host":"10.0.16.2","port":16501,"machine_id":2,"server_id":3},
      {"node_id":4,"server_name":"SgridTestNodeServer","pid":32069,"host":"10.0.16.2","port":19421,"machine_id":2,"server_id":2},
   ]
}
````
