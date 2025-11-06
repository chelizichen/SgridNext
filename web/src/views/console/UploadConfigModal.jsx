import React, { useEffect, useState, useRef } from 'react';
import { Modal, Upload, Button, message, Form, Input } from 'antd';
import { getConfigContent, upsertConfig } from './api';
import { createTwoFilesPatch } from 'diff';
import { Diff2HtmlUI } from 'diff2html/lib/ui/js/diff2html-ui';
import 'highlight.js/styles/googlecode.css';  // 代码高亮样式
import 'diff2html/bundles/css/diff2html.min.css';  // diff2html 样式
import ButtonGroup from 'antd/es/button/button-group';

export default function UploadConfigModal({ visible, onOk, onCancel,serverId,fileName }) {
    const [messageApi, contextHolder] = message.useMessage();
    const [form] = Form.useForm()
    const [loading, setLoading] = useState(false);
    const [diffModalVisible, setDiffModalVisible] = useState(false);
    const [prevData, setPrevData] = useState();
    const [curData, setCurData] = useState();
    useEffect(()=>{
        if(fileName){
            getConfigContent({
                serverId,
                configName:fileName
            }).then(res=>{
                form.setFieldsValue({
                    fileContent:res.data,
                    fileName
                })
                setPrevData(res.data)
            })

        }
    },[fileName,serverId])
    const handleUpsertConfig = () => {
        setDiffModalVisible(false)
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
    const handleDiff = ()=>{
        setCurData(form.getFieldValue("fileContent"))
        setDiffModalVisible(true)
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
                    onClick={handleDiff}
                    loading={loading}
                    style={{ marginTop: 16 }}
                    block
                    disabled={fileName.includes("_") }
                >
                    对比
                </Button>
                <DiffModal handleUpsertConfig={handleUpsertConfig} isVisible={diffModalVisible} onClose={() => setDiffModalVisible(false)} prevData={prevData} curData={curData} />
            </Modal>
        </>
    );
}


const DiffComponent = ({ prevData, curData, prevFileName, curFileName }) => {
    const diffRef = useRef(null);
    useEffect(() => {
      const formatForDiff = (data) => {
        if (typeof data === 'string') {
          try {
            const parsed = JSON.parse(data);
            return JSON.stringify(parsed, null, 2);
          } catch (e) {
            return data; // 非 JSON 字符串则原样返回
          }
        }
        try {
          return JSON.stringify(data, null, 2);
        } catch (e) {
          return String(data);
        }
      };

      const diffOutput = createTwoFilesPatch(
        prevFileName,   // 左侧文件名
        curFileName,    // 右侧文件名
        formatForDiff(prevData || ''),  // 左侧内容（尽量格式化为多行）
        formatForDiff(curData),   // 右侧内容（尽量格式化为多行）
        '', '',                                 // 标题和前缀（可选）
      );
  
      const targetElement = diffRef.current;
      const configuration = {
        drawFileList: true,
        outputFormat: 'side-by-side',         // 并排显示差异
        highlight: true,
        synchronizeScroll: true,
      };
  
      const diff2htmlUi = new Diff2HtmlUI(targetElement, diffOutput, configuration);
      diff2htmlUi.draw();
      diff2htmlUi.highlightCode();
    }, [prevData, curData, prevFileName, curFileName]);
  
    return <div ref={diffRef}></div>;
  };

const DiffModal = ({ isVisible, onClose, prevData, curData, handleUpsertConfig }) => (
    <Modal width={1400} open={isVisible} onCancel={onClose} footer={null}>
        <DiffComponent prevData={prevData} curData={curData} prevFileName="上次配置" curFileName="当前配置" />
        <Button type="primary" onClick={handleUpsertConfig} style={{marginTop: 16,width:'100%'}}>上传</Button>
    </Modal>
);
