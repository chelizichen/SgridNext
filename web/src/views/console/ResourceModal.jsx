import React, { useEffect } from 'react';
import { Modal, Form, InputNumber, List, message, Input,Select } from 'antd';
import { setCpuLimit, setMemoryLimit, updateServerNode } from './api';
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
                        .then(async values => {
                            await setCpuLimit({
                                cpuLimit: values.cpu,
                                nodeIds: nodes.map(node => node.id),
                                serverId,
                            });
                            await setMemoryLimit({
                                memoryLimit: values.memory,
                                nodeIds: nodes.map(node => node.id),
                                serverId,
                            });
                            
                            const args = JSON.stringify(values.additional_args.split(";").filter(v=>v))
                            await updateServerNode({
                                ids: nodes.map(node => node.id),
                                server_run_type: values.server_run_type,
                                additional_args: args,
                            })
                            message.success('修改成功')
                            onOk(values);
                        })
                        .catch(info => {
                            // 可以在这里处理验证失败的逻辑
                            messageApi.error(info.errorFields[0].errors[0]);
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
                    <Form.Item
                        name="additional_args"
                        label="参数"
                        rules={[{ required: true, message: '请输入CPU核数' }]}
                    >
                        <Input.TextArea style={{width:'100%'}} rows={3}/>
                    </Form.Item>
                    <Form.Item 
                        name="server_run_type" 
                        label="运行类型" 
                        rules={[{ required: true, message: '请选择运行类型' }]}
                    > 
                        <Select style={{width:'100%'}} placeholder="请选择运行类型"
                            options={[
                                { label: '手动重启', value: 0 },
                                { label: '自动重启', value: 12 },
                            ]}
                        />
                    </Form.Item>
                </Form>
            </Modal>
        </>
    );
}