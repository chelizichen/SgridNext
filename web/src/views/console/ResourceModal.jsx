import React, { useEffect } from 'react';
import { Modal, Form, InputNumber, List, message } from 'antd';
import { setCpuLimit } from './api';
export default function ResourceModal({ visible, onOk, onCancel, form, nodes,serverId }) {
    const [messageApi, contextHolder] = message.useMessage();
    return (
        <>
            {contextHolder}
            <Modal
                title="资源配置"
                open={visible}
                onOk={() => {
                    form.validateFields()
                        .then(values => {
                            setCpuLimit({
                                cpuLimit: values.cpu,
                                nodeIds: nodes.map(node => node.id),
                                serverId,
                            });
                            onOk(values);
                        })
                        .catch(info => {
                            // 可以在这里处理验证失败的逻辑
                        });
                }}
                onCancel={onCancel}
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
                        <InputNumber min={1} max={32768} style={{ width: '100%' }} />
                    </Form.Item>
                    <Form.Item
                        name="cpu"
                        label="CPU核数"
                        rules={[{ required: true, message: '请输入CPU核数' }]}
                    >
                        <InputNumber min={0.1} max={32} style={{ width: '100%' }} />
                    </Form.Item>
                </Form>
            </Modal>
        </>
    );
}