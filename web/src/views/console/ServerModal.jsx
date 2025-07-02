import React from 'react';
import { Modal, Form, InputNumber, Select, Input,message } from 'antd';
import { createServer } from './api';

export default function ServerModal({ visible, onOk, onCancel, form, groupOptions,afterOpenChange }) {
    const [messageApi, contextHolder] = message.useMessage();
    return (
        <>
            {contextHolder}
            <Modal
                title="创建服务"
                open={visible}
                onOk={() => {
                    form.validateFields().then(values => {
                    createServer(values)
                            .then(response => {
                                if (response.success) {
                                    onOk(values);
                                    messageApi.info(response.msg);
                                } else {
                                    messageApi.error(response.msg);
                                }
                            })
                            .catch(error => {
                                console.error('请求错误:', error);
                            });
                    });
                }}
                afterOpenChange={afterOpenChange}
                onCancel={onCancel}
            >
                <Form form={form} layout="vertical">
                    <Form.Item name="serverName" label="服务名称" rules={[{ required: true, message: '请输入服务名称' }]}> 
                        <Input style={{width:'100%'}} /> 
                    </Form.Item>
                    <Form.Item name="groupId" label="所属服务组" rules={[{ required: true, message: '请选择所属服务组' }]}> 
                        <Select options={groupOptions} style={{width:'100%'}} placeholder="请选择服务组" /> 
                    </Form.Item>
                    <Form.Item name="serverType" label="服务类型" rules={[{ required: true, message: '请选择服务类型' }]}> 
                        <Select style={{width:'100%'}} placeholder="请选择服务类型"
                            options={[
                                { label: 'NODE', value: 1 },
                                { label: 'JAVA', value: 2 },
                                { label: 'BINARY', value: 3 }
                            ]}
                        />
                    </Form.Item>
                    <Form.Item name="execFilePath" label="执行路径" rules={[{ required: true, message: '请输入执行文件路径' }]}> 
                        <Input style={{width:'100%'}} min={1} />
                    </Form.Item>
                    <Form.Item name="logPath" label="日志文件" rules={[{ required: false, message: '请输入日志文件路径' }]}> 
                        <Input style={{width:'100%'}} />
                    </Form.Item>
                    <Form.Item name="confPath" label="配置文件（TODO）" rules={[{ required: false, message: '请输入配置文件路径' }]}> 
                        <Input style={{width:'100%'}} />
                    </Form.Item>
                    <Form.Item name="description" label="服务描述" rules={[{ required: true, message: '请输入执行文件路径' }]}> 
                        <Input.TextArea style={{width:'100%'}} min={1} />
                    </Form.Item>

                </Form>
            </Modal>
        </>
    );
}