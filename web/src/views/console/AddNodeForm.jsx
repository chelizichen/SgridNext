import React from 'react';
import { Modal, Form, InputNumber, Select, Input,message } from 'antd';
import { createNode } from './api';

export default function AddNodeModal({ visible, onOk, onCancel, form,afterOpenChange }) {
    const [messageApi, contextHolder] = message.useMessage();
    return (
        <>
            {contextHolder}
            <Modal
                title="创建节点"
                open={visible}
                onOk={() => {
                    form.validateFields().then(values => {
                      createNode(values)
                            .then(response => {
                                console.log('响应数据:', response.data);
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
                    <Form.Item name="Host" label="主机地址" rules={[{ required: true, message: '请输入服务名称' }]}> 
                        <Input style={{width:'100%'}} /> 
                    </Form.Item>
                    <Form.Item name="Cpus" label="CPU核心数" rules={[{ required: true, message: '请选择所属服务组' }]}> 
                      <InputNumber step={1} style={{width:'100%'}} /> 
                    </Form.Item>
                    <Form.Item name="Memory" label="内存大小" rules={[{ required: true, message: '请选择所属服务组' }]}> 
                      <InputNumber step={1} style={{width:'100%'}} /> 
                    </Form.Item>
                    <Form.Item name="Os" label="操作系统" rules={[{ required: true, message: '请选择所属服务组' }]}> 
                      <Input style={{width:'100%'}} /> 
                    </Form.Item>
                </Form>
            </Modal>
        </>
    );
}