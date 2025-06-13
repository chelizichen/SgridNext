import React, { useState, useEffect, use } from 'react';
import { Table, Tag, Space, Button } from 'antd';
import {  getServerConfigList,getConfigContent } from '../../console/api';
import { useLocation } from 'react-router-dom';
import UploadConfigModal from '../../console/UploadConfigModal';

export default function NodeManager() {
  const [nodes, setNodes] = useState([]);
  const [serverId, setServerId] = useState();
  const [serverConfigVisible,setServerConfigVisible] =  useState(false);
  const [fileName,setFileName] = useState('');
  const location = useLocation()
  function getFileContentByName(fileName){
    setServerConfigVisible(true)
    setFileName(fileName)
  }

  function onOk(){
    setServerConfigVisible(false)
  }
  const columns = [
    {
      title: '配置文件名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Button onClick={() => getFileContentByName(record.name)}>详情</Button>
        </Space>
      ),
    },
  ];

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const serverId = Number(params.get("id"))
    setServerId(serverId)
    getServerConfigList({
      serverId: serverId,
    }).then((res) => {
      setNodes(res.data.map(v=>({name:v})));
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
      <UploadConfigModal 
        visible={serverConfigVisible}
        onOk={onOk}
        onCancel={onOk}
        serverId={serverId}
        fileName={fileName}
      ></UploadConfigModal>
    </div>
  );
}