import React, { useState } from 'react';
import { Row, Col, Card, List, Divider, Tree, Button, Modal, Form, InputNumber, Table } from 'antd';
import ButtonGroup from 'antd/es/button/button-group';

export default function Console(){
    const [selectedNode, setSelectedNode] = useState(null);
    const [resourceModalVisible, setResourceModalVisible] = useState(false);
    const [form] = Form.useForm();

    const handleRefresh = () => {
        // 这里可以添加实际的刷新逻辑
        console.log('刷新数据');
    };
    
    // 示例树数据
    const treeData = [
        {
            title: '服务组1',
            key: 'group1',
            children: [
                { title: '节点1', key: 'node1' },
                { title: '节点2', key: 'node2' },
            ],
        },
        {
            title: '服务组2',
            key: 'group2',
            children: [
                { title: '节点3', key: 'node3' },
            ],
        },
    ];
    
    // 示例日志数据
    const logData = [
        { time: '10:30', status: '运行中', message: '节点1启动完成' },
        { time: '10:25', status: '警告', message: '节点2负载过高' },
        { time: '10:20', status: '正常', message: '节点3连接成功' },
    ];
    const rowSelection = {
        onChange: (selectedRowKeys, selectedRows) => {
          console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
        },
        getCheckboxProps: record => ({
          disabled: record.name === 'Disabled User', // Column configuration not to be checked
          name: record.name,
        }),
      };
    return (
        <div style={{ padding: 24 }}>
            <Row gutter={16}>
                <Col span={8}>
                    <Card title="服务列表"
                        variant={false} extra={
                        <ButtonGroup>
                            <Button onClick={handleRefresh}  style={{marginLeft:'16px'}}>刷新</Button>
                            <Button onClick={handleRefresh}>添加服务</Button>
                        </ButtonGroup>
                    }>
                        <Tree
                            treeData={treeData}
                            onSelect={(keys, { node }) => setSelectedNode(node)}
                        />
                    </Card>
                    <Divider />
                    <Card title="节点信息" extra={<Button onClick={handleRefresh}>刷新</Button>}>
                    <List
                            dataSource={logData}
                            renderItem={item => (
                                <List.Item>
                                    [{item.time}] [{item.status}] {item.message}
                                </List.Item>
                            )}
                        />
                    </Card>
                </Col>
                <Col span={16}>
                    <Card title="服务信息" variant={true} extra={<Button onClick={handleRefresh}>刷新</Button>}>
                        {selectedNode ? (
                            <div>
                                <p>名称: {selectedNode.title}</p>
                                <p>ID: {selectedNode.key}</p>
                                <p>状态: 运行中</p>
                            </div>
                        ) : (
                            <p>请从左侧选择节点</p>
                        )}
                    </Card>
                    <Divider />
                    <Card title="节点列表" variant={false} extra={
                        <div>
                            <Button onClick={handleRefresh}>刷新</Button>
                            <ButtonGroup style={{marginLeft:"16px"}}>
                                <Button onClick={handleRefresh}>部署</Button>
                                <Button onClick={() => setResourceModalVisible(true)}>资源配置</Button>
                                <Button onClick={handleRefresh}>扩容</Button>
                                <Button onClick={handleRefresh}>删除</Button>
                            </ButtonGroup>
                        </div>

                    }>
                        <Table 
                                rowSelection={Object.assign({ type: "checkbox" }, rowSelection)}
                            columns={[
                                { title: '节点名称', dataIndex: 'title', key: 'title' },
                                { title: '节点ID', dataIndex: 'key', key: 'key' },
                                { title: '状态', key: 'status', render: () => '运行中' }
                            ]}
                            dataSource={treeData.flatMap(group => group.children)}
                            rowKey="key"
                            pagination={false}
                        />
                    </Card>
                    <Divider />
                    <Card title="状态日志" bo={false} extra={<Button onClick={handleRefresh}>刷新</Button>}>
                        <List
                            dataSource={logData}
                            renderItem={item => (
                                <List.Item>
                                    [{item.time}] [{item.status}] {item.message}
                                </List.Item>
                            )}
                        />
                    </Card>
                </Col>
            </Row>
            
            <Modal
                title="资源配置"
                open={resourceModalVisible}
                onOk={() => {
                    form.validateFields()
                        .then(values => {
                            console.log('资源配置:', values);
                            setResourceModalVisible(false);
                        })
                        .catch(info => {
                            console.log('验证失败:', info);
                        });
                }}
                onCancel={() => setResourceModalVisible(false)}
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="memory"
                        label="内存限制(MB)"
                        rules={[{ required: true, message: '请输入内存限制' }]}
                    >
                        <InputNumber min={1} max={32768} style={{ width: '100%' }} />
                    </Form.Item>
                    <Form.Item
                        name="cpu"
                        label="CPU核数"
                        rules={[{ required: true, message: '请输入CPU核数' }]}
                    >
                        <InputNumber min={1} max={32} style={{ width: '100%' }} />
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
}


