import React, { useEffect, useState } from 'react';
import { Modal, Upload, Button, message, Form, Input } from 'antd';
import { getConfigContent, upsertConfig } from './api';

export default function UploadConfigModal({ visible, onOk, onCancel,serverId,fileName }) {
    const [messageApi, contextHolder] = message.useMessage();
    const [form] = Form.useForm()
    const [loading, setLoading] = useState(false);
    useEffect(()=>{
        console.log('111');
        if(fileName){
            getConfigContent({
                serverId,
                configName:fileName
            }).then(res=>{
                form.setFieldsValue({
                    fileContent:res.data,
                    fileName
                })
            })

        }
    },[fileName,serverId])
    const handleUpsertConfig = () => {
        setLoading(true);
        form.validateFields().then(async (values) => {
            let body = {
                serverId,
                ...values
            }
            let res = await upsertConfig(body);
            if (res.success) {
                messageApi.success('上传成功');
                onOk();
            } else {
                messageApi.error(res.msg);
            }
            setLoading(false);
        }).catch(() => {
            setLoading(false);
        })
    }
    return (
        <>
            {contextHolder}
            <Modal
                title="上传配置文件"
                open={visible}
                onCancel={onCancel}
                footer={null}
                width={800}
            >
                <Form form={form} layout="vertical">
                    <Form.Item name="fileName" label="配置文件名称" rules={[{ required: true, message: '请输入服务名称' }]}> 
                        <Input style={{width:'100%'}} /> 
                    </Form.Item>
                    <Form.Item name="fileContent" label="内容" rules={[{ required: true, message: '请输入服务名称' }]}> 
                        <Input.TextArea style={{width:'100%'}} rows={20} /> 
                    </Form.Item>
                </Form>
                <Button
                    type="primary"
                    onClick={handleUpsertConfig}
                    loading={loading}
                    style={{ marginTop: 16 }}
                    block
                    disabled={fileName.includes("_") }
                >
                    上传
                </Button>
            </Modal>
        </>
    );
}