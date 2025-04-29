import React, { useState } from 'react';
import { Modal, Upload, Input, Button, List, message } from 'antd';
import { UploadOutlined } from '@ant-design/icons';

export default function DeployModal({ visible, onOk, onCancel, nodes }) {
  const [fileList, setFileList] = useState([]);
  const [commitMsg, setCommitMsg] = useState('');
  const [uploading, setUploading] = useState(false);
  const [publishList, setPublishList] = useState([]);
  const [selectedPublishes, setSelectedPublishes] = useState([]);

  const handleUpload = () => {
    if (!fileList.length) {
      message.warning('请先选择文件');
      return;
    }
    setUploading(true);
    // 模拟上传
    setTimeout(() => {
      setUploading(false);
      // 假设上传后返回发布列表
      setPublishList(nodes.map(node => ({
        key: node.key,
        title: node.title,
        status: '待发布'
      })));
      message.success('文件上传成功');
    }, 1000);
  };

  const handlePublish = () => {
    if (!selectedPublishes.length) {
      message.warning('请选择要发布的节点');
      return;
    }
    // 模拟发布
    setPublishList(publishList.map(item =>
      selectedPublishes.includes(item.key)
        ? { ...item, status: '已发布' }
        : item
    ));
    message.success('发布成功');
    onOk && onOk();
  };

  return (
    <Modal
      title="部署发布"
      open={visible}
      onCancel={onCancel}
      footer={null}
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
      {publishList.length > 0 && (
        <>
          <List
            header={<div>发布列表（可多选）</div>}
            bordered
            dataSource={publishList}
            renderItem={item => (
              <List.Item
                actions={[
                  <input
                    type="checkbox"
                    checked={selectedPublishes.includes(item.key)}
                    onChange={e => {
                      if (e.target.checked) {
                        setSelectedPublishes([...selectedPublishes, item.key]);
                      } else {
                        setSelectedPublishes(selectedPublishes.filter(k => k !== item.key));
                      }
                    }}
                  />,
                  <span>{item.status}</span>
                ]}
              >
                {item.title}
              </List.Item>
            )}
          />
          <Button
            type="primary"
            style={{ marginTop: 16 }}
            onClick={handlePublish}
            disabled={!selectedPublishes.length}
            block
          >
            发布
          </Button>
        </>
      )}
    </Modal>
  );
}