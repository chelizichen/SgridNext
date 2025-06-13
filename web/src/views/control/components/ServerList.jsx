import React, { useState, useEffect } from 'react';
import { Table, Tag, Space, Button } from 'antd';
import { getServerList } from '../../console/api';
import { useNavigate } from 'react-router-dom';

export default function NodeManager() {
  const [nodes, setNodes] = useState([]);
  const navigate = useNavigate()
  const columns = [
    {
      title: '服务ID',
      dataIndex: 'server_id',
      key: 'server_id',
    },
    {
      title: '服务名称',
      dataIndex: 'server_name',
      key: 'server_name',
    },
    {
      title: '服务组',
      dataIndex: 'group_name',
      key: 'group_name',
    },
    // {
    //   title: '状态',
    //   key: 'status',
    //   render: (_, record) => (
    //     <Tag color={record.status === 'online' ? 'green' : 'red'}>
    //       {record.status}
    //     </Tag>
    //   ),
    // },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Button onClick={()=>toConfList(record)}>配置文件历史</Button>
          <Button danger>下线</Button>
        </Space>
      ),
    },
  ];
  function toConfList(record) {
    console.log('record',record);
    navigate('/control/config_list?id='+record.server_id);
  }
  useEffect(() => {
    getServerList().then((res) => {
      console.log('res',res);
      setNodes(res.data);
      // setServerList(res.data);
    });

    // TODO: 对接后端节点状态接口
    // const mockData = [
    //   { id: 'node-001', host: '192.168.1.101', port: 8080, status: 'online' },
    //   { id: 'node-002', host: '192.168.1.102', port: 8081, status: 'offline' },
    // ];
  }, []);

  return (
    <div>
      {/* <div style={{ marginBottom: 16 }}> */}
        {/* <Button type="primary">新增节点</Button> */}
      {/* </div> */}
      <Table 
        bordered
        columns={columns} 
        dataSource={nodes} 
        rowKey="id"
        pagination={false}
      />
    </div>
  );
}