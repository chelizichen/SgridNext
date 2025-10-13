import React, { useEffect, useState } from 'react';
import { Modal, Card, Row, Col, Progress, Table, Tag, Button, message } from 'antd';
import { getNodeResource } from './api';

export default function NodeResourceModal({ visible, onCancel, nodeId }) {
    const [messageApi, contextHolder] = message.useMessage();
    const [loading, setLoading] = useState(false);
    const [resourceData, setResourceData] = useState(null);

    const fetchResourceData = async () => {
        setLoading(true);
        try {
            const response = await getNodeResource(nodeId);
            if (response.success) {
                setResourceData(response.data);
            } else {
                messageApi.error(response.msg || '获取资源信息失败');
            }
        } catch (error) {
            messageApi.error('获取资源信息失败: ' + error.message);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (visible) {
            fetchResourceData();
        }
    }, [visible]);

    const processColumns = [
        {
            title: 'PID',
            dataIndex: 'PID',
            key: 'PID',
            width: 80,
        },
        {
            title: '进程名',
            dataIndex: 'Name',
            key: 'Name',
            width: 150,
        },
        {
            title: 'CPU使用率',
            dataIndex: 'CPUPercent',
            key: 'CPUPercent',
            width: 120,
            render: (value) => (
                <Progress 
                    percent={Math.round(value * 100) / 100} 
                    size="small" 
                    status={value > 80 ? 'exception' : value > 60 ? 'active' : 'normal'}
                />
            ),
        },
        {
            title: '内存使用率',
            dataIndex: 'MemoryPercent',
            key: 'MemoryPercent',
            width: 120,
            render: (value) => (
                <Progress 
                    percent={Math.round(value * 100) / 100} 
                    size="small" 
                    status={value > 80 ? 'exception' : value > 60 ? 'active' : 'normal'}
                />
            ),
        },
        {
            title: '内存使用量(MB)',
            dataIndex: 'RSS',
            key: 'RSS',
            width: 120,
            render: (value) => `${Math.round(value / 1024 / 1024 * 100) / 100}`,
        },
    ];

    return (
        <>
            {contextHolder}
            <Modal
                title="节点资源监控"
                open={visible}
                onCancel={onCancel}
                width={1200}
                footer={[
                    <Button key="refresh" type="primary" loading={loading} onClick={fetchResourceData}>
                        刷新
                    </Button>,
                    <Button key="close" onClick={onCancel}>
                        关闭
                    </Button>
                ]}
            >
                {resourceData && (
                    <div>
                        {/* 系统资源概览 */}
                        <Row gutter={16} style={{ marginBottom: 16 }}>
                            <Col span={12}>
                                <Card title="CPU使用率" size="small">
                                    <Progress 
                                        percent={Math.round(resourceData.systemInfo.cpuInfo.usage * 100) / 100} 
                                        status={resourceData.systemInfo.cpuInfo.usage > 80 ? 'exception' : resourceData.systemInfo.cpuInfo.usage > 60 ? 'active' : 'normal'}
                                    />
                                    <div style={{ marginTop: 8, fontSize: '12px', color: '#666' }}>
                                        {Math.round(resourceData.systemInfo.cpuInfo.usage * 100) / 100}%
                                    </div>
                                </Card>
                            </Col>
                            <Col span={12}>
                                <Card title="内存使用率" size="small">
                                    <Progress 
                                        percent={Math.round(resourceData.systemInfo.memoryInfo.usage * 100) / 100} 
                                        status={resourceData.systemInfo.memoryInfo.usage > 80 ? 'exception' : resourceData.systemInfo.memoryInfo.usage > 60 ? 'active' : 'normal'}
                                    />
                                    <div style={{ marginTop: 8, fontSize: '12px', color: '#666' }}>
                                        已用: {Math.round(resourceData.systemInfo.memoryInfo.used * 100) / 100} MB / 
                                        总计: {Math.round(resourceData.systemInfo.memoryInfo.total * 100) / 100} MB
                                    </div>
                                </Card>
                            </Col>
                        </Row>

                        {/* 进程列表 */}
                        <Card title="进程列表" size="small">
                            <Table
                                columns={processColumns}
                                dataSource={resourceData.processInfo.slice(0, 20)} // 只显示前20个进程
                                rowKey="PID"
                                size="small"
                                pagination={{
                                    pageSize: 10,
                                    showSizeChanger: false,
                                    showQuickJumper: true,
                                }}
                                scroll={{ y: 400 }}
                            />
                        </Card>
                    </div>
                )}
            </Modal>
        </>
    );
}