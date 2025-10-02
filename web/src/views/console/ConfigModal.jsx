import React, { useState, useEffect, useCallback } from 'react';
import { Modal, Form, Input, Button, message, Card, Row, Col, Divider, Tag, Space, Tooltip, Popconfirm } from 'antd';
import { EditOutlined, SaveOutlined, ReloadOutlined, InfoCircleOutlined, PlusOutlined, DeleteOutlined } from '@ant-design/icons';
import { getMainConfig, updateMainConfig, setConfigItem } from './api';

export default function ConfigModal({ visible, onOk, onCancel }) {
    const [messageApi, contextHolder] = message.useMessage();
    const [form] = Form.useForm();
    const [configData, setConfigData] = useState({});
    const [editingKey, setEditingKey] = useState(null);
    const [addingNew, setAddingNew] = useState(false);
    const [newKey, setNewKey] = useState('');
    const [newValue, setNewValue] = useState('');
    const [loading, setLoading] = useState(false);
    const [saving, setSaving] = useState(false);

    // 加载配置数据
    const loadConfig = useCallback(async () => {
        setLoading(true);
        try {
            const response = await getMainConfig();
            if (response.success) {
                setConfigData(response.data.config);
                form.setFieldsValue(response.data.config);
            } else {
                messageApi.error(response.msg);
            }
        } catch (error) {
            messageApi.error('加载配置失败');
            console.error('加载配置错误:', error);
        } finally {
            setLoading(false);
        }
    }, [form, messageApi]);

    useEffect(() => {
        if (visible) {
            loadConfig();
        }
    }, [visible, loadConfig]);

    // 开始编辑
    const startEdit = (key) => {
        setEditingKey(key);
    };

    // 取消编辑
    const cancelEdit = () => {
        setEditingKey(null);
        form.setFieldsValue(configData);
    };

    // 保存单个配置项
    const saveConfigItem = async (key, value) => {
        setSaving(true);
        try {
            // 先更新本地状态
            setConfigData(prev => ({ ...prev, [key]: value }));
            
            const response = await setConfigItem({ key, value });
            if (response.success) {
                messageApi.success('配置项保存成功');
                setEditingKey(null);
            } else {
                messageApi.error(response.msg);
                // 如果保存失败，恢复原值
                setConfigData(prev => ({ ...prev, [key]: configData[key] }));
            }
        } catch (error) {
            messageApi.error('保存配置项失败');
            console.error('保存配置项错误:', error);
            // 如果保存失败，恢复原值
            setConfigData(prev => ({ ...prev, [key]: configData[key] }));
        } finally {
            setSaving(false);
        }
    };

    // 删除配置项
    const deleteConfigItem = async (key) => {
        setSaving(true);
        try {
            const newConfig = { ...configData };
            delete newConfig[key];
            const response = await updateMainConfig({ config: newConfig });
            if (response.success) {
                messageApi.success('配置项删除成功');
                setConfigData(newConfig);
                form.setFieldsValue(newConfig);
            } else {
                messageApi.error(response.msg);
            }
        } catch (error) {
            messageApi.error('删除配置项失败');
            console.error('删除配置项错误:', error);
        } finally {
            setSaving(false);
        }
    };

    // 添加新配置项
    const addNewConfigItem = async () => {
        if (!newKey.trim()) {
            messageApi.error('请输入配置项名称');
            return;
        }
        if (Object.prototype.hasOwnProperty.call(configData, newKey)) {
            messageApi.error('配置项已存在');
            return;
        }

        setSaving(true);
        try {
            let value = newValue;
            // 尝试解析JSON
            if (newValue.trim().startsWith('{') || newValue.trim().startsWith('[')) {
                try {
                    value = JSON.parse(newValue);
                } catch {
                    // 如果不是有效JSON，保持原值
                }
            }
            
            const response = await setConfigItem({ key: newKey, value });
            if (response.success) {
                messageApi.success('配置项添加成功');
                setConfigData(prev => ({ ...prev, [newKey]: value }));
                setAddingNew(false);
                setNewKey('');
                setNewValue('');
            } else {
                messageApi.error(response.msg);
            }
        } catch (error) {
            messageApi.error('添加配置项失败');
            console.error('添加配置项错误:', error);
        } finally {
            setSaving(false);
        }
    };

    // 保存整个配置
    const handleSaveAll = async () => {
        setSaving(true);
        try {
            // 使用当前的configData而不是form.getFieldsValue()
            // 因为form可能没有正确绑定所有字段
            console.log('保存配置数据:', configData);
            
            if (Object.keys(configData).length === 0) {
                messageApi.warning('没有配置数据可保存');
                setSaving(false);
                return;
            }
            
            const response = await updateMainConfig({ config: configData });
            if (response.success) {
                messageApi.success('配置保存成功');
                onOk && onOk();
            } else {
                messageApi.error(response.msg);
            }
        } catch (error) {
            messageApi.error('保存配置失败');
            console.error('保存配置错误:', error);
        } finally {
            setSaving(false);
        }
    };

    // 渲染配置项
    const renderConfigItem = (key, value) => {
        const isEditing = editingKey === key;
        const isComplexValue = typeof value === 'object' && value !== null;

        return (
            <Card key={key} size="small" style={{ marginBottom: 8 }}>
                <Row align="middle" gutter={16}>
                    <Col span={6}>
                        <Space>
                            <strong>{key}</strong>
                            {isComplexValue && (
                                <Tooltip title="复杂对象，请使用JSON编辑器">
                                    <InfoCircleOutlined style={{ color: '#1890ff' }} />
                                </Tooltip>
                            )}
                        </Space>
                    </Col>
                    <Col span={14}>
                        {isEditing ? (
                            <Form.Item name={key} style={{ margin: 0 }}>
                                {isComplexValue ? (
                                    <Input.TextArea 
                                        rows={4} 
                                        placeholder="请输入JSON格式的值"
                                        defaultValue={JSON.stringify(value, null, 2)}
                                        onBlur={(e) => {
                                            try {
                                                const parsed = JSON.parse(e.target.value);
                                                saveConfigItem(key, parsed);
                                            } catch {
                                                messageApi.error('JSON格式错误');
                                            }
                                        }}
                                    />
                                ) : (
                                    <Input 
                                        defaultValue={String(value)}
                                        onPressEnter={(e) => saveConfigItem(key, e.target.value)}
                                        onBlur={(e) => saveConfigItem(key, e.target.value)}
                                        autoFocus
                                    />
                                )}
                            </Form.Item>
                        ) : (
                            <div style={{ padding: '4px 8px', background: '#f5f5f5', borderRadius: 4 }}>
                                {isComplexValue ? (
                                    <pre style={{ margin: 0, fontSize: '12px' }}>
                                        {JSON.stringify(value, null, 2)}
                                    </pre>
                                ) : (
                                    <span>{String(value)}</span>
                                )}
                            </div>
                        )}
                    </Col>
                    <Col span={4}>
                        <Space>
                            {isEditing ? (
                                <>
                                    <Button 
                                        size="small" 
                                        onClick={cancelEdit}
                                        disabled={saving}
                                    >
                                        取消
                                    </Button>
                                </>
                            ) : (
                                <>
                                    <Button 
                                        size="small" 
                                        icon={<EditOutlined />}
                                        onClick={() => startEdit(key)}
                                    >
                                        编辑
                                    </Button>
                                    <Popconfirm
                                        title="确定要删除这个配置项吗？"
                                        onConfirm={() => deleteConfigItem(key)}
                                        okText="确定"
                                        cancelText="取消"
                                    >
                                        <Button 
                                            size="small" 
                                            danger
                                            icon={<DeleteOutlined />}
                                            disabled={saving}
                                        >
                                            删除
                                        </Button>
                                    </Popconfirm>
                                </>
                            )}
                        </Space>
                    </Col>
                </Row>
            </Card>
        );
    };

    return (
        <>
            {contextHolder}
            <Modal
                title="主控配置管理"
                open={visible}
                onCancel={onCancel}
                width={1000}
                footer={[
                    <Button key="reload" icon={<ReloadOutlined />} onClick={loadConfig} loading={loading}>
                        刷新
                    </Button>,
                    <Button key="save" type="primary" icon={<SaveOutlined />} onClick={handleSaveAll} loading={saving}>
                        保存全部
                    </Button>,
                    <Button key="cancel" onClick={onCancel}>
                        关闭
                    </Button>
                ]}
            >
                <Form form={form} layout="vertical">
                    {/* 调试信息 */}
                    {process.env.NODE_ENV === 'development' && (
                        <div style={{ marginBottom: 16, padding: 8, background: '#f0f0f0', borderRadius: 4, fontSize: '12px' }}>
                            <strong>调试信息:</strong> 当前配置项数量: {Object.keys(configData).length}
                            <br />
                            配置数据: {JSON.stringify(configData, null, 2)}
                        </div>
                    )}
                    <div style={{ maxHeight: '60vh', overflowY: 'auto' }}>
                        {/* 添加新配置项 */}
                        {addingNew && (
                            <Card size="small" style={{ marginBottom: 8, border: '2px dashed #1890ff' }}>
                                <Row align="middle" gutter={16}>
                                    <Col span={6}>
                                        <Input 
                                            placeholder="配置项名称"
                                            value={newKey}
                                            onChange={(e) => setNewKey(e.target.value)}
                                        />
                                    </Col>
                                    <Col span={14}>
                                        <Input.TextArea 
                                            rows={2}
                                            placeholder="配置项值 (支持JSON格式)"
                                            value={newValue}
                                            onChange={(e) => setNewValue(e.target.value)}
                                        />
                                    </Col>
                                    <Col span={4}>
                                        <Space>
                                            <Button 
                                                size="small" 
                                                type="primary"
                                                onClick={addNewConfigItem}
                                                loading={saving}
                                            >
                                                保存
                                            </Button>
                                            <Button 
                                                size="small" 
                                                onClick={() => {
                                                    setAddingNew(false);
                                                    setNewKey('');
                                                    setNewValue('');
                                                }}
                                            >
                                                取消
                                            </Button>
                                        </Space>
                                    </Col>
                                </Row>
                            </Card>
                        )}
                        
                        {/* 现有配置项 */}
                        {Object.entries(configData).map(([key, value]) => 
                            renderConfigItem(key, value)
                        )}
                        
                        {/* 添加按钮 */}
                        {!addingNew && (
                            <Card size="small" style={{ marginTop: 8, textAlign: 'center' }}>
                                <Button 
                                    icon={<PlusOutlined />}
                                    onClick={() => setAddingNew(true)}
                                    disabled={saving}
                                >
                                    添加新配置项
                                </Button>
                            </Card>
                        )}
                    </div>
                </Form>
            </Modal>
        </>
    );
}
