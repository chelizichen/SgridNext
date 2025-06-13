import React, { useState, useEffect } from 'react';
import { Table, Tag, Space, Button } from 'antd';
import { getSyncStatus } from '../../console/api';
import { useLocation } from 'react-router-dom';

let mainControlNode = 0

export default function NodeStat() {
  const [nodes, setNodes] = useState({
    update_time:"",
    stat_list: [],
  });
  const location = useLocation()
  const columns = [
    {
      title: '机器节点ID',
      dataIndex: 'machine_id',
      key: 'machine_id', 
    },
    {
      title: '服务名',
      dataIndex: 'server_name',
      key: 'server_name',
    },
    {
      title: '服务节点ID',
      dataIndex: 'node_id',
      key: 'node_id',
    },
    {
      title: '主机地址',
      dataIndex: 'host',
      key: 'host',
    },
    {
      title: '端口号',
      dataIndex: 'port',
      key: 'port',
    },
    {
      title: '进程ID',
      dataIndex: 'pid',
      key: 'pid',
    },
    // {
    //   title: '操作',
    //   key: 'action',
    //   render: (_, record) => (
    //     <Space size="middle">
    //       <Button>详情</Button>
    //       {
    //         record.NodeStatus === 1 ? (
    //           <Button danger>下线</Button>
    //         ) : (
    //           <Button>上线</Button>
    //         )
    //       }
    //     </Space>
    //   ),
    // },
  ];

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    console.log('params',params);
    
    const nodeId = Number(params.get("id"))
    getSyncStatus({
      nodeId: nodeId || mainControlNode,
    }).then((res) => {
      console.log('res',res);
      setNodes(res.data);
    });
  }, [location.search]);

  return (
    <div>
      {/* <div style={{ marginBottom: 16 }}> */}
        {/* <Button type="primary">新增节点</Button> */}
      {/* </div> */}
      <Table 
        columns={columns} 
        title={() => <div>节点状态更新时间  {nodes.update_time}</div>}
        dataSource={nodes.stat_list} 
        rowKey="id"
        pagination={false}
      />
    </div>
  );
}