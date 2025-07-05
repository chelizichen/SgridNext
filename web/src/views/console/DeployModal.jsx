import React, { useState,useEffect } from 'react';
import { Modal, Upload, Input, Button, List, message,Divider, Table } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import { deployServer, getServerPackageList, uploadPackage } from './api';
import { downloadFile } from './api';
import { _constant } from '../../common/constant';

export default function DeployModal({ visible, onOk, onCancel, serverInfo,nodes }) {
  const [fileList, setFileList] = useState([]);
  const [commitMsg, setCommitMsg] = useState('');
  const [uploading, setUploading] = useState(false);
  const [selectedPublishNodes, setSelectedPublishNodes] = useState([]);
  const [messageApi, contextHolder] = message.useMessage();
  const [uploadPackageList,setUploadPackageList] = useState([]);
  const [packageId,setPackageId] = useState(0);
  const [offset,setOffset] = useState(0);
  const [size,setSize] = useState(5);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);

  function initUploadPackageList(){
    setLoading(true);
    getServerPackageList({id:serverInfo.server_id,offset:offset,size:size}).then(res=>{
      if(res.success){
        setUploadPackageList(res.data);
        setTotal(res.total || res.data.length);
      }else{
        messageApi.error(res.msg);
      }
      setLoading(false);
    })
  }

  useEffect(() => {
    if(!serverInfo.server_id) return;
    initUploadPackageList()
    setSelectedPublishNodes([])
  }, [serverInfo.server_id, offset, size])

  const handleTableChange = (pagination) => {
    setOffset((pagination.current - 1) * pagination.pageSize);
    setSize(pagination.pageSize);
  };

  const handleUpload = () => {
    setUploading(true);
    console.log('serverInfo',serverInfo);
    if (!fileList.length) {
      messageApi.warning('请先选择文件');
      return;
    }
    if(!fileList[0].name.includes(serverInfo.server_name)){
      messageApi.warning('请选择正确的文件');
      return;
    }

    const data = new FormData();
    data.append('file', fileList[0]);
    data.append('commit', commitMsg);
    // 发送请求到后端，上传文件并生成发布列
    data.append("serverName",serverInfo.server_name)
    data.append("serverId",serverInfo.server_id)
    
  

    uploadPackage(data).then(res=>{
      if(res.success){
        messageApi.success(res.msg);
        initUploadPackageList()
      }else{
        messageApi.error(res.msg);
      }
      setUploading(false);
    })
  };

  const handlePublish = () => {
    if (!selectedPublishNodes.length) {
      messageApi.warning('请选择要发布的节点');
      return;
    }
    if(!packageId){
      messageApi.warning('请选择要发布的包');
      return;
    }
    const body = {
      serverNodeIds: selectedPublishNodes,
      packageId: packageId,
      serverId: serverInfo.server_id,
    }
    deployServer(body).then(res=>{
      if(!res.success){
        messageApi.error(res.msg);
        return;
      }
      messageApi.success('发布成功');
    })
    onOk && onOk();
  };

  useEffect(()=>{
    setOffset(0);
    setSize(5);
  },[serverInfo.server_id])

  function handleDownload(record){
    console.log('record',record);
    downloadFile({
      serverId:serverInfo.server_id,
      fileName:record.FileName,
      type:_constant.FILE_TYPE_PACKAGE,
      host:record.Host
    })
  }

  return (
    <>
      {contextHolder}
      <Modal
        title="部署发布"
        open={visible}
        onCancel={onCancel}
        footer={null}
        width={800}
        destroyOnClose
      >
        <div style={{ marginBottom: 16 }}>
          <Upload
            fileList={fileList}
            beforeUpload={file => {
              setFileList([file]);
              return false;
            }}
            onRemove={() => setFileList([])}
            maxCount={1}
          >
            <Button icon={<UploadOutlined />}>选择文件</Button>
          </Upload>
        </div>
        <div style={{ marginBottom: 16 }}>
          <Input.TextArea
            rows={2}
            placeholder="请输入commit内容"
            value={commitMsg}
            onChange={e => setCommitMsg(e.target.value)}
          />
        </div>
        <div style={{ marginBottom: 16 }}>
          <Button
            type="primary"
            onClick={handleUpload}
            loading={uploading}
            disabled={!fileList.length || !commitMsg}
            block
          >
            上传并生成发布列表
          </Button>
        </div>
        {nodes.length > 0 && (
          <>
            <div style={{ marginBottom: 16 }}>
              <h4>节点列表（可多选）</h4>
            </div>
            <Table
              dataSource={nodes}
              columns={[
                {
                  title: '选择',
                  key: 'select',
                  width: 80,
                  render: (_, record) => (
                    <input
                      type="checkbox"
                      checked={selectedPublishNodes.includes(record.id)}
                      onChange={e => {
                        if (e.target.checked) {
                          setSelectedPublishNodes([...selectedPublishNodes, record.id]);
                        } else {
                          setSelectedPublishNodes(selectedPublishNodes.filter(k => k !== record.id));
                        }
                      }}
                    />
                  )
                },
                {
                  title: '别名',
                  dataIndex: 'alias',
                  key: 'alias',
                  width: 150,
                },
                {
                  title: 'Host',
                  dataIndex: 'host',
                  key: 'host',
                  width: 150,
                },
                {
                  title: 'Port',
                  dataIndex: 'port',
                  key: 'port',
                  width: 150,
                },
                {
                  title: 'Patch ID',
                  dataIndex: 'patch_id',
                  key: 'patch_id',
                  width: 100,
                }
              ]}
              rowKey="id"
              pagination={false}
              size="small"
              bordered
            />
          </>
        )}
        <div style={{ marginBottom: '16px',marginTop: '16px' }}>
          <h4>发布列表</h4>
        </div>
        {
          uploadPackageList.length > 0 && (
              <>
                <Table
                  dataSource={uploadPackageList}
                  bordered
                  columns={[
                    {
                      title: '选择',
                      key: 'select',
                      width: 80,
                      render: (_, record) => (
                        <input
                          type="checkbox"
                          checked={packageId === record.ID}
                          onChange={e => {
                            if (e.target.checked) {
                              setPackageId(record.ID)
                            } else {
                              setPackageId(0)
                            }
                          }}
                        />
                      )
                    },
                    {
                      title: 'ID',
                      dataIndex: 'ID',
                      key: 'ID',
                      width: 80,
                    },
                    {
                      title: 'Hash',
                      dataIndex: 'Hash',
                      key: 'Hash',
                      width: 200,
                    },
                    {
                      title: '创建时间',
                      dataIndex: 'CreateTime',
                      key: 'CreateTime',
                      width: 150,
                    },
                    {
                      title: 'Commit',
                      dataIndex: 'Commit',
                      key: 'Commit',
                      ellipsis: true,
                    },
                    {
                      title: '操作',
                      key: 'action',
                      width: 100,
                      render: (_, record) => (
                        <Button  onClick={()=>{
                          handleDownload(record)
                        }}>下载</Button>
                      )
                    } 
                  ]}
                  rowKey="ID"
                  pagination={{
                    current: Math.floor(offset / size) + 1,
                    pageSize: size,
                    total: total,
                    showSizeChanger: true,
                    showQuickJumper: true,
                    showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条/共 ${total} 条`,
                  }}
                  onChange={handleTableChange}
                  loading={loading}
                  size="small"
                />
              </>
          )
        }
        <Button
          type="primary"
          style={{ marginTop: 16 }}
          onClick={handlePublish}
          disabled={!selectedPublishNodes.length || !packageId}
          block
        >
          发布
        </Button>
      </Modal>
    </>
  );
}