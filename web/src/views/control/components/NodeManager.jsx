import React, { useState, useEffect } from 'react';
import { Table, Tag, Space, Button } from 'antd';

export default function NodeManager() {
  const [nodes, setNodes] = useState([]);

  const columns = [
    {
      title: '节点ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: '主机地址',
      dataIndex: 'host',
      key: 'host',
    },
    {
      title: '端口',
      dataIndex: 'port',
      key: 'port',
    },
    {
      title: '状态',
      key: 'status',
      render: (_, record) => (
        <Tag color={record.status === 'online' ? 'green' : 'red'}>
          {record.status}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Button type="link">详情</Button>
          <Button danger type="link">下线</Button>
        </Space>
      ),
    },
  ];

  useEffect(() => {
    // TODO: 对接后端节点状态接口
    const mockData = [
      { id: 'node-001', host: '192.168.1.101', port: 8080, status: 'online' },
      { id: 'node-002', host: '192.168.1.102', port: 8081, status: 'offline' },
    ];
    setNodes(mockData);
  }, []);

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary">新增节点</Button>
      </div>
      <Table 
        columns={columns} 
        dataSource={nodes} 
        rowKey="id"
        pagination={false}
      />
    </div>
  );
}