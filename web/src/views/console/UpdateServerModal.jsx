import React from 'react';
import { Modal, Form, Input, Select, message } from 'antd';
import { updateServer } from './api';

export default function UpdateServerModal({ visible, onOk, onCancel, serverInfo }) {
    const [messageApi, contextHolder] = message.useMessage();
    const [form] = Form.useForm();

    // 当模态框打开时，设置表单的初始值
    React.useEffect(() => {
        if (visible && serverInfo) {
            form.setFieldsValue({
                serverName: serverInfo.server_name,
                dockerName: serverInfo.docker_name,
                serverType: serverInfo.server_type,
                execFilePath: serverInfo.exec_path,
                logPath: serverInfo.log_path,
                description: serverInfo.desc,
                configPath: serverInfo.config_path
            });
        }
    }, [visible, serverInfo, form]);

    const handleOk = () => {
        form.validateFields().then(values => {
            const updateData = {
                id: serverInfo.server_id,
                ...values
            };
            
            updateServer(updateData)
                .then(response => {
                    if (response.success) {
                        messageApi.success('服务更新成功');
                        onOk(values);
                    } else {
                        messageApi.error(response.msg);
                    }
                })
                .catch(error => {
                    console.error('更新服务错误:', error);
                    messageApi.error('更新服务失败');
                });
        });
    };

    return (
        <>
            {contextHolder}
            <Modal
                title="更新服务"
                open={visible}
                onOk={handleOk}
                onCancel={onCancel}
                width={600}
            >
                <Form form={form} layout="vertical">
                    <Form.Item name="serverName" label="服务名称" rules={[{ required: true, message: '请输入服务名称' }]}> 
                        <Input style={{width:'100%'}} /> 
                    </Form.Item>
                    <Form.Item name="dockerName" label="docker名称" rules={[{ required: false, message: '请输入docker名称' }]}> 
                        <Input style={{width:'100%'}} /> 
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
                        <Input style={{width:'100%'}} />
                    </Form.Item>
                    <Form.Item name="logPath" label="日志文件" rules={[{ required: false, message: '请输入日志文件路径' }]}> 
                        <Input style={{width:'100%'}} />
                    </Form.Item>
                    <Form.Item name="configPath" label="配置文件" rules={[{ required: false, message: '请输入配置文件路径' }]}> 
                        <Input style={{width:'100%'}} />
                    </Form.Item>
                    <Form.Item name="description" label="服务描述" rules={[{ required: true, message: '请输入服务描述' }]}> 
                        <Input.TextArea style={{width:'100%'}} rows={3} />
                    </Form.Item>
                </Form>
            </Modal>
        </>
    );
}
