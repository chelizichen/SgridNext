import React, { useState, useEffect } from 'react';
import { Table, Tag, Space, Button } from 'antd';
import { getNodeList } from '../../console/api';
import { useNavigate } from 'react-router-dom';

export default function NodeManager() {
  const [nodes, setNodes] = useState([]);
  const navigate = useNavigate();
  const columns = [
    {
      title: '节点ID',
      dataIndex: 'ID',
      key: 'ID',
    },
    {
      title: '主机地址',
      dataIndex: 'Host',
      key: 'Host',
    },
    {
      title: '操作系统',
      dataIndex: 'Os',
      key: 'Os',
    },
    {
      title: '处理器核心数(CORE)',
      dataIndex: 'Cpus',
      key: 'Cpus',
    },
    {
      title: '内存大小(G)',
      dataIndex: 'Memory',
      key: 'Memory',
    },
    {
      title: '状态',
      key: 'NodeStatus',
      render: (_, record) => (
        <Tag color={record.NodeStatus === 1 ? 'green' : 'red'}>
          {record.NodeStatus === 1 ? '在线' : '离线'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Button>详情</Button>
          <Button onClick={() => navigate(`/control/nodestat_list?id=${record.ID}`)}>节点服务</Button>
          {
            record.NodeStatus === 1 ? (
              <Button danger>下线</Button>
            ) : (
              <Button>上线</Button>
            )
          }
        </Space>
      ),
    },
  ];

  useEffect(() => {
    getNodeList().then((res) => {
      setNodes(res.data);
    });
  }, []);

  return (
    <div>
      {/* <div style={{ marginBottom: 16 }}> */}
        {/* <Button type="primary">新增节点</Button> */}
      {/* </div> */}
      <Table 
        columns={columns} 
        dataSource={nodes} 
        rowKey="id"
        pagination={false}
      />
    </div>
  );
}