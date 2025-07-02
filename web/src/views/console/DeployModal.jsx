import React, { useState,useEffect } from 'react';
import { Modal, Upload, Input, Button, List, message,Divider } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import { deployServer, getServerPackageList, uploadPackage } from './api';

export default function DeployModal({ visible, onOk, onCancel, serverInfo,nodes }) {
  const [fileList, setFileList] = useState([]);
  const [commitMsg, setCommitMsg] = useState('');
  const [uploading, setUploading] = useState(false);
  const [selectedPublishNodes, setSelectedPublishNodes] = useState([]);
  const [messageApi, contextHolder] = message.useMessage();
  const [uploadPackageList,setUploadPackageList] = useState([]);
  const [packageId,setPackageId] = useState(0);

  function initUploadPackageList(){
    getServerPackageList({id:serverInfo.server_id}).then(res=>{
      if(res.success){
        setUploadPackageList(res.data);
      }else{
        messageApi.error(res.msg);
      }})
  }

  useEffect(() => {
    if(!serverInfo.server_id) return;
    initUploadPackageList()
    setSelectedPublishNodes([])
  }, [serverInfo.server_id])

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
            <List
              header={<div>节点列表（可多选）</div>}
              bordered
              dataSource={nodes}
              renderItem={item => (
                <List.Item
                  actions={[
                    <input
                      type="checkbox"
                      checked={selectedPublishNodes.includes(item.id)}
                      onChange={e => {
                        if (e.target.checked) {
                          setSelectedPublishNodes([...selectedPublishNodes, item.id]);
                        } else {
                          setSelectedPublishNodes(selectedPublishNodes.filter(k => k !== item.id));
                        }
                      }}
                    />,
                    <span>host</span>,
                    <span>{item.host}</span>,
                    <span>patch_id</span>,
                    <span>{item.patch_id}</span>,
                  ]}
                >
                  {item.title}
                </List.Item>
              )}
            />
          </>
        )}
        <Divider />
        {
          uploadPackageList.length > 0 && (
              <>
                <List
                  header={<div>发布列表</div>}
                  bordered
                  dataSource={uploadPackageList}
                  renderItem={item => (
                    <List.Item
                      actions={[
                        <input
                          type="checkbox"
                          checked={packageId == item.ID }
                          onChange={e => {
                            if (e.target.checked) {
                              setPackageId(item.ID)
                            } else {
                              setPackageId(0)
                            }
                          }}
                        />,
                        <span>id:{item.ID}</span>,
                        <span>{item.Hash.slice(0,5)}</span>,
                        <span>{item.CreateTime}</span>,
                        <span>{item.Commit}</span>
                      ]}
                    >
                      {item.title}
                    </List.Item>
                  )}
                />
              </>
          )
        }
        <Button
          type="primary"
          style={{ marginTop: 16 }}
          onClick={handlePublish}
          disabled={!selectedPublishNodes.length && !packageId}
          block
        >
          发布
        </Button>
      </Modal>
    </>
  );
}