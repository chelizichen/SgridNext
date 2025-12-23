import React, { useState, useEffect } from 'react';
import {
  Card,
  Table,
  Button,
  Modal,
  Form,
  Input,
  Upload,
  message,
  Space,
  Tag,
  Select,
  Drawer,
  Typography,
  Divider,
  Popconfirm,
} from 'antd';
import {
  PlusOutlined,
  UploadOutlined,
  EditOutlined,
  DeleteOutlined,
  DownloadOutlined,
  EyeOutlined,
  LinkOutlined,
  FileTextOutlined,
} from '@ant-design/icons';
import ReactMarkdown from 'react-markdown';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';
import {
  uploadDocument,
  createDocument,
  updateDocument,
  deleteDocument,
  getDocumentList,
  getDocument,
  downloadDocument,
  linkDocumentToServer,
  getDocumentServerRelations,
} from '../console/api';
import { getServerList } from '../console/api';

const { TextArea } = Input;
const { Option } = Select;
const { Title, Paragraph } = Typography;

const DocumentManager = () => {
  const [documents, setDocuments] = useState([]);
  const [loading, setLoading] = useState(false);
  const [editingDoc, setEditingDoc] = useState(null);
  const [viewingDoc, setViewingDoc] = useState(null);
  const [linkModalVisible, setLinkModalVisible] = useState(false);
  const [currentDocId, setCurrentDocId] = useState(null);
  const [servers, setServers] = useState([]);
  const [selectedServers, setSelectedServers] = useState([]);
  const [form] = Form.useForm();
  const [viewForm] = Form.useForm();

  useEffect(() => {
    loadDocuments();
    loadServers();
  }, []);

  const loadDocuments = async () => {
    setLoading(true);
    try {
      const res = await getDocumentList();
      if (res.success) {
        setDocuments(res.data || []);
      } else {
        message.error(res.msg || '加载文档列表失败');
      }
    } catch (error) {
      message.error('加载文档列表失败: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  const loadServers = async () => {
    try {
      const res = await getServerList();
      if (res.success) {
        setServers(res.data || []);
      }
    } catch (error) {
      console.error('加载服务列表失败:', error);
    }
  };

  const handleUpload = async (file) => {
    const formData = new FormData();
    formData.append('file', file);
    
    try {
      const res = await uploadDocument(formData);
      if (res.success) {
        message.success('上传成功');
        loadDocuments();
      } else {
        message.error(res.msg || '上传失败');
      }
    } catch (error) {
      message.error('上传失败: ' + error.message);
    }
    return false; // 阻止自动上传
  };

  const handleCreate = () => {
    setEditingDoc({});
    form.resetFields();
    form.setFieldsValue({
      title: '',
      content: '',
      description: '',
    });
  };

  const handleEdit = async (record) => {
    try {
      const res = await getDocument(record.id);
      if (res.success) {
        setEditingDoc(res.data);
        form.setFieldsValue({
          id: res.data.id,
          title: res.data.title,
          content: res.data.content,
          description: res.data.description,
        });
      } else {
        message.error(res.msg || '获取文档失败');
      }
    } catch (error) {
      message.error('获取文档失败: ' + error.message);
    }
  };

  const handleView = async (record) => {
    try {
      const res = await getDocument(record.id);
      if (res.success) {
        setViewingDoc(res.data);
        viewForm.setFieldsValue({
          title: res.data.title,
          content: res.data.content,
          description: res.data.description,
        });
      } else {
        message.error(res.msg || '获取文档失败');
      }
    } catch (error) {
      message.error('获取文档失败: ' + error.message);
    }
  };

  const handleSave = async (values) => {
    try {
      let res;
      if (editingDoc && editingDoc.id) {
        res = await updateDocument({
          id: editingDoc.id,
          ...values,
        });
      } else {
        res = await createDocument(values);
      }
      if (res.success) {
        message.success(editingDoc ? '更新成功' : '创建成功');
        form.resetFields();
        setEditingDoc(null);
        loadDocuments();
      } else {
        message.error(res.msg || '保存失败');
      }
    } catch (error) {
      message.error('保存失败: ' + error.message);
    }
  };

  const handleDelete = async (id) => {
    try {
      const res = await deleteDocument({ id });
      if (res.success) {
        message.success('删除成功');
        loadDocuments();
      } else {
        message.error(res.msg || '删除失败');
      }
    } catch (error) {
      message.error('删除失败: ' + error.message);
    }
  };

  const handleDownload = async (record) => {
    try {
      const response = await downloadDocument(record.id);
      const blob = new Blob([response.data], { type: 'text/markdown' });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = record.fileName || `${record.title}.md`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      message.success('下载成功');
    } catch (error) {
      message.error('下载失败: ' + error.message);
    }
  };

  const handleLink = async (record) => {
    setCurrentDocId(record.id);
    try {
      const res = await getDocumentServerRelations(record.id);
      if (res.success) {
        setSelectedServers(res.data || []);
      }
    } catch (error) {
      console.error('获取关联服务失败:', error);
    }
    setLinkModalVisible(true);
  };

  const handleLinkSave = async () => {
    try {
      const res = await linkDocumentToServer({
        documentId: currentDocId,
        serverIds: selectedServers,
      });
      if (res.success) {
        message.success('关联成功');
        setLinkModalVisible(false);
        loadDocuments();
      } else {
        message.error(res.msg || '关联失败');
      }
    } catch (error) {
      message.error('关联失败: ' + error.message);
    }
  };

  const getServerNames = (serverIds) => {
    if (!serverIds || serverIds.length === 0) return '未关联';
    return serverIds
      .map((id) => {
        const server = servers.find((s) => s.server_id === id);
        return server ? server.server_name : id;
      })
      .join(', ');
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      render: (text) => text || '-',
    },
    {
      title: '关联服务',
      dataIndex: 'serverIds',
      key: 'serverIds',
      width: 200,
      ellipsis: true,
      render: (serverIds) => (
        <Tag color={serverIds && serverIds.length > 0 ? 'blue' : 'default'}>
          {getServerNames(serverIds)}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createTime',
      key: 'createTime',
      width: 180,
    },
    {
      title: '更新时间',
      dataIndex: 'updateTime',
      key: 'updateTime',
      width: 180,
    },
    {
      title: '操作',
      key: 'action',
      width:540,
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handleView(record)}
          >
            查看
          </Button>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            icon={<DownloadOutlined />}
            onClick={() => handleDownload(record)}
          >
            下载
          </Button>
          <Button
            type="link"
            icon={<LinkOutlined />}
            onClick={() => handleLink(record)}
          >
            关联
          </Button>
          <Popconfirm
            title="确定要删除这个文档吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              danger
              icon={<DeleteOutlined />}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: '24px', background: '#f0f2f5', minHeight: '100vh' }}>
      <Card
        title={
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <FileTextOutlined />
            <span>文档管理中心</span>
          </div>
        }
        extra={
          <Space>
            <Upload
              accept=".md"
              beforeUpload={handleUpload}
              showUploadList={false}
            >
              <Button type="default" icon={<UploadOutlined />}>
                上传文档
              </Button>
            </Upload>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={handleCreate}
            >
              编写文档
            </Button>
          </Space>
        }
        style={{ boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}
      >
        <Table
          columns={columns}
          dataSource={documents}
          rowKey="id"
          loading={loading}
          pagination={{
            pageSize: 10,
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 条`,
          }}
          scroll={{ x: 1200 }}
        />
      </Card>

      {/* 编辑/创建文档 Modal */}
      <Modal
        title={editingDoc && editingDoc.id ? '编辑文档' : '编写文档'}
        open={editingDoc !== null}
        onCancel={() => {
          setEditingDoc(null);
          form.resetFields();
        }}
        onOk={() => form.submit()}
        width={900}
        okText="保存"
        cancelText="取消"
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSave}
        >
          <Form.Item
            name="title"
            label="文档标题"
            rules={[{ required: true, message: '请输入文档标题' }]}
          >
            <Input placeholder="请输入文档标题" />
          </Form.Item>
          <Form.Item
            name="description"
            label="文档描述"
          >
            <Input.TextArea
              rows={2}
              placeholder="请输入文档描述（可选）"
            />
          </Form.Item>
          <Form.Item
            name="content"
            label="文档内容 (Markdown)"
            rules={[{ required: true, message: '请输入文档内容' }]}
          >
            <TextArea
              rows={20}
              placeholder="请输入 Markdown 格式的文档内容..."
              style={{ fontFamily: 'Monaco, Consolas, monospace' }}
            />
          </Form.Item>
        </Form>
      </Modal>

      {/* 查看文档 Drawer */}
      <Drawer
        title={
          viewingDoc ? (
            <div>
              <Title level={4} style={{ margin: 0 }}>
                {viewingDoc.title}
              </Title>
              {viewingDoc.description && (
                <Paragraph type="secondary" style={{ marginTop: '8px', marginBottom: 0 }}>
                  {viewingDoc.description}
                </Paragraph>
              )}
            </div>
          ) : (
            '查看文档'
          )
        }
        placement="right"
        width={800}
        open={viewingDoc !== null}
        onClose={() => {
          setViewingDoc(null);
          viewForm.resetFields();
        }}
        extra={
          viewingDoc && (
            <Space>
              <Button
                icon={<LinkOutlined />}
                onClick={() => {
                  handleLink(viewingDoc);
                  setViewingDoc(null);
                }}
              >
                关联服务
              </Button>
              <Button
                icon={<EditOutlined />}
                onClick={() => {
                  handleEdit(viewingDoc);
                  setViewingDoc(null);
                }}
              >
                编辑
              </Button>
              <Button
                icon={<DownloadOutlined />}
                onClick={() => handleDownload(viewingDoc)}
              >
                下载
              </Button>
            </Space>
          )
        }
      >
        {viewingDoc && (
          <div style={{ padding: '16px 0' }}>
            {/* 关联服务信息 */}
            {viewingDoc.serverIds && viewingDoc.serverIds.length > 0 && (
              <div style={{ marginBottom: '16px', padding: '12px', background: '#f0f2f5', borderRadius: '4px' }}>
                <div style={{ fontWeight: 'bold', marginBottom: '8px' }}>
                  <LinkOutlined style={{ marginRight: '4px' }} />
                  关联的服务：
                </div>
                <div>
                  {viewingDoc.serverIds.map((serverId) => {
                    const server = servers.find((s) => s.server_id === serverId);
                    return (
                      <Tag key={serverId} color="blue" style={{ marginRight: '8px', marginBottom: '4px' }}>
                        {server ? server.server_name : serverId}
                      </Tag>
                    );
                  })}
                </div>
                <Button
                  type="link"
                  size="small"
                  icon={<LinkOutlined />}
                  onClick={() => {
                    handleLink(viewingDoc);
                    setViewingDoc(null);
                  }}
                  style={{ marginTop: '8px', padding: 0 }}
                >
                  修改关联
                </Button>
              </div>
            )}
            <div
              style={{
                background: '#fff',
                padding: '24px',
                borderRadius: '8px',
                minHeight: '400px',
              }}
            >
              <ReactMarkdown
                components={{
                  code({ inline, className, children, ...props }) {
                    const match = /language-(\w+)/.exec(className || '');
                    return !inline && match ? (
                      <SyntaxHighlighter
                        style={vscDarkPlus}
                        language={match[1]}
                        PreTag="div"
                        {...props}
                      >
                        {String(children).replace(/\n$/, '')}
                      </SyntaxHighlighter>
                    ) : (
                      <code className={className} {...props}>
                        {children}
                      </code>
                    );
                  },
                }}
              >
                {viewingDoc.content}
              </ReactMarkdown>
            </div>
          </div>
        )}
      </Drawer>

      {/* 关联服务 Modal */}
      <Modal
        title={
          <div>
            <LinkOutlined style={{ marginRight: '8px' }} />
            关联服务
          </div>
        }
        open={linkModalVisible}
        onOk={handleLinkSave}
        onCancel={() => {
          setLinkModalVisible(false);
          setSelectedServers([]);
        }}
        okText="保存关联"
        cancelText="取消"
        width={600}
      >
        <div style={{ marginBottom: '16px' }}>
          <Paragraph>
            选择一个或多个服务与此文档关联。关联后，可以在服务详情中查看相关文档。
          </Paragraph>
        </div>
        <Select
          mode="multiple"
          style={{ width: '100%' }}
          placeholder="请选择要关联的服务（可多选）"
          value={selectedServers}
          onChange={setSelectedServers}
          showSearch
          filterOption={(input, option) =>
            option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
          }
          size="large"
        >
          {servers.map((server) => (
            <Option key={server.server_id} value={server.server_id}>
              {server.server_name}
            </Option>
          ))}
        </Select>
        {selectedServers.length > 0 && (
          <div style={{ marginTop: '16px', padding: '12px', background: '#f0f2f5', borderRadius: '4px' }}>
            <div style={{ fontWeight: 'bold', marginBottom: '8px' }}>已选择 {selectedServers.length} 个服务：</div>
            <div>
              {selectedServers.map((serverId) => {
                const server = servers.find((s) => s.server_id === serverId);
                return (
                  <Tag key={serverId} color="blue" style={{ marginBottom: '4px' }}>
                    {server ? server.server_name : serverId}
                  </Tag>
                );
              })}
            </div>
          </div>
        )}
        <div style={{ marginTop: '16px', padding: '12px', background: '#e6f7ff', borderRadius: '4px', border: '1px solid #91d5ff' }}>
          <div style={{ color: '#1890ff', fontSize: '12px', lineHeight: '1.6' }}>
            <strong>💡 提示：</strong>
            <ul style={{ margin: '8px 0 0 20px', padding: 0 }}>
              <li>一个文档可以关联多个服务</li>
              <li>一个服务也可以关联多个文档</li>
              <li>关联后，可以在服务管理页面查看相关文档</li>
            </ul>
          </div>
        </div>
      </Modal>
    </div>
  );
};

export default DocumentManager;

