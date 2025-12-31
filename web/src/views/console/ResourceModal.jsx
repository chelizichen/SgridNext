import React, { useState } from 'react';
import { Modal, Form, InputNumber, List, message, Input, Select, Button, Space } from 'antd';
import { setCpuLimit, setMemoryLimit, updateServerNode } from './api';
export default function ResourceModal({ visible, onOk, onCancel, form, nodes, serverId }) {
    const [messageApi, contextHolder] = message.useMessage();
    const [loading, setLoading] = useState({
        memory: false,
        cpu: false,
        server_run_type: false,
        additional_args: false,
        view_page: false,
    });

    // 设置内存限制
    const handleSetMemory = async () => {
        try {
            const value = form.getFieldValue('memory');
            if (!value) {
                messageApi.error('请输入内存限制');
                return;
            }
            setLoading(prev => ({ ...prev, memory: true }));
            await setMemoryLimit({
                memoryLimit: value,
                nodeIds: nodes.map(node => node.id),
                serverId,
            });
            messageApi.success('内存限制设置成功');
        } catch {
            messageApi.error('设置失败');
        } finally {
            setLoading(prev => ({ ...prev, memory: false }));
        }
    };

    // 设置CPU限制
    const handleSetCpu = async () => {
        try {
            const value = form.getFieldValue('cpu');
            if (!value) {
                messageApi.error('请输入CPU核数');
                return;
            }
            setLoading(prev => ({ ...prev, cpu: true }));
            await setCpuLimit({
                cpuLimit: value,
                nodeIds: nodes.map(node => node.id),
                serverId,
            });
            messageApi.success('CPU核数设置成功');
        } catch {
            messageApi.error('设置失败');
        } finally {
            setLoading(prev => ({ ...prev, cpu: false }));
        }
    };

    // 设置运行类型
    const handleSetRunType = async () => {
        try {
            const value = form.getFieldValue('server_run_type');
            if (value === undefined || value === null) {
                messageApi.error('请选择运行类型');
                return;
            }
            setLoading(prev => ({ ...prev, server_run_type: true }));
            await updateServerNode({
                ids: nodes.map(node => node.id),
                server_run_type: value,
            });
            messageApi.success('运行类型设置成功');
        } catch {
            messageApi.error('设置失败');
        } finally {
            setLoading(prev => ({ ...prev, server_run_type: false }));
        }
    };

    // 设置参数
    const handleSetAdditionalArgs = async () => {
        try {
            const value = form.getFieldValue('additional_args');
            const args = value ? JSON.stringify(value.split(";").filter(v => v)) : JSON.stringify([]);
            setLoading(prev => ({ ...prev, additional_args: true }));
            await updateServerNode({
                ids: nodes.map(node => node.id),
                additional_args: args,
            });
            messageApi.success('参数设置成功');
        } catch {
            messageApi.error('设置失败');
        } finally {
            setLoading(prev => ({ ...prev, additional_args: false }));
        }
    };

    // 设置预览地址
    const handleSetViewPage = async () => {
        try {
            const value = form.getFieldValue('view_page');
            setLoading(prev => ({ ...prev, view_page: true }));
            await updateServerNode({
                ids: nodes.map(node => node.id),
                view_page: value || '',
            });
            messageApi.success('预览地址设置成功');
        } catch {
            messageApi.error('设置失败');
        } finally {
            setLoading(prev => ({ ...prev, view_page: false }));
        }
    };

    return (
        <>
            {contextHolder}
            <Modal
                title="资源配置"
                open={visible}
                onOk={onOk}
                onCancel={onCancel}
                footer={[
                    <Button key="cancel" onClick={onCancel}>
                        关闭
                    </Button>,
                ]}
            >
                <List
                    dataSource={nodes}
                    renderItem={item => (
                        <List.Item>
                            {item.host}:{item.port}
                        </List.Item>
                    )}
                />
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="memory"
                        label="内存限制(MB)"
                        rules={[{ required: true, message: '请输入内存限制' }]}
                    >
                        <Space.Compact style={{ width: '100%' }}>
                            <InputNumber min={1} max={32768} style={{ width: '100%' }} />
                            <Button 
                                type="primary" 
                                onClick={handleSetMemory}
                                loading={loading.memory}
                            >
                                设置
                            </Button>
                        </Space.Compact>
                    </Form.Item>
                    <Form.Item
                        name="cpu"
                        label="CPU核数"
                        rules={[{ required: true, message: '请输入CPU核数' }]}
                    >
                        <Space.Compact style={{ width: '100%' }}>
                            <InputNumber min={0.1} max={32} style={{ width: '100%' }} />
                            <Button 
                                type="primary" 
                                onClick={handleSetCpu}
                                loading={loading.cpu}
                            >
                                设置
                            </Button>
                        </Space.Compact>
                    </Form.Item>
                    <Form.Item 
                        name="server_run_type" 
                        label="运行类型" 
                        rules={[{ required: true, message: '请选择运行类型' }]}
                    > 
                        <Space.Compact style={{ width: '100%' }}>
                            <Select 
                                style={{ width: '100%' }} 
                                placeholder="请选择运行类型"
                                options={[
                                    { label: '手动重启', value: 0 },
                                    { label: '自动重启', value: 12 },
                                ]}
                            />
                            <Button 
                                type="primary" 
                                onClick={handleSetRunType}
                                loading={loading.server_run_type}
                            >
                                设置
                            </Button>
                        </Space.Compact>
                    </Form.Item>
                    <Form.Item
                        name="additional_args"
                        label="参数"
                    >
                        <Space.Compact style={{ width: '100%' }}>
                            <Input style={{ width: '100%' }} />
                            <Button 
                                type="primary" 
                                onClick={handleSetAdditionalArgs}
                                loading={loading.additional_args}
                            >
                                设置
                            </Button>
                        </Space.Compact>
                    </Form.Item>
                    <Form.Item
                        name="view_page"
                        label="预览地址"
                    >
                        <Space.Compact style={{ width: '100%' }}>
                            <Input />
                            <Button 
                                type="primary" 
                                onClick={handleSetViewPage}
                                loading={loading.view_page}
                            >
                                设置
                            </Button>
                        </Space.Compact>
                    </Form.Item>
                </Form>
            </Modal>
        </>
    );
}