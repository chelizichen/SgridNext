import React from 'react';
import { Modal, Form, Input, message } from 'antd';
import { createGroup } from './api';

export default function GroupModal({ visible, onCancel,onOk, form }) {
    const [messageApi, contextHolder] = message.useMessage();
    return (
        <>
            {contextHolder}
            <Modal
                title="创建服务组"
                open={visible}
                onOk={() => {
                    form.validateFields().then(values => {
                        createGroup(values)
                            .then(response => {
                                if (response.success) {
                                    messageApi.info(response.msg);
                                    onOk(values);
                                } else {
                                    messageApi.error(response.msg);
                                }
                            })
                            .catch(error => {
                                console.error('请求错误:', error);
                            });
                    });
                }}
                onCancel={onCancel}
            >
                <Form form={form} layout="vertical">
                    <Form.Item name="groupName" label="服务组名称" rules={[{ required: true, message: '请输入服务组名称' }]}> 
                        <Input style={{width:'100%'}} /> 
                    </Form.Item>
                    <Form.Item name="groupEnglishName" label="服务组英文名称" rules={[{ required: true, message: '请输入服务组名称' }]}> 
                        <Input style={{width:'100%'}} /> 
                    </Form.Item>
                </Form>
            </Modal>
        </>

    );
}