import React, { useState } from 'react';
import { Modal, Table, Form, InputNumber, Button, Card, message } from 'antd';
import { createServerNode } from './api';
import { getServerType } from './constant';
import { Descriptions } from 'antd';
import _ from 'lodash';
export default function ScaleModal({ visible, onOk, onCancel, nodes,serverInfo }) {
    const [messageApi, contextHolder] = message.useMessage();
    const [selectedRowKeys, setSelectedRowKeys] = useState([]);
    const [portMap, setPortMap] = useState({});
    const [form] = Form.useForm();

    const handleSelectChange = (keys) => {
      console.log('keys',keys);
      setSelectedRowKeys(keys);
    };

    const handlePortChange = (key, value) => {
      setPortMap(prev => ({ ...prev, [key]: value }));
      console.log(key,value);
    };

    const handleOk = () => {
      form
        .validateFields()
        .then(() => {
          let selectedNodes = nodes.filter(node => selectedRowKeys.includes(node.Host));
          selectedNodes = selectedNodes.map(node => {
            node.port = portMap[node.Host];
            node.server_id = serverInfo.server_id;
            node.node_id = node.ID;
            node.patch_id = 0
            return _.pick(node,["server_id","node_id","patch_id","port"])
          });
          createServerNode(selectedNodes).then(res=>{
            if(res.success){
              messageApi.info(res.msg);
              onOk()
            }else{
              messageApi.error(res.msg);
            }
          })
        });
    };

    const columns = [
      { title: '主机地址', dataIndex: 'Host', key: 'Host' },
      { title: '内存大小（G）', dataIndex: 'Memory', key: 'Memory' },
      { title: 'CPU核心数', dataIndex: 'Cpus', key: 'Cpus' },
      {
        title: '分配端口',
        key: 'port',
        render: (_, record) => (
          <Form.Item
            name={`port_${record.Host}`}
            rules={[{ required: selectedRowKeys.includes(record.Host), message: '请输入端口' }]}
            style={{ margin: 0 }}
          >
            <InputNumber
              min={10001}
              max={65535}
              disabled={!selectedRowKeys.includes(record.Host)}
              value={portMap[record.key]}
              onChange={value => handlePortChange(record.Host, value)}
              placeholder="请输入端口"
            />
          </Form.Item>
        )
      }
    ];

  return (
    <>
      {contextHolder}
      <Modal
        title="扩容节点选择"
        open={visible}
        onOk={handleOk}
        onCancel={onCancel}
        destroyOnClose
        width={800}
      >
        <Form form={form} layout="vertical">
          <Card style={{marginBottom:"16px"}}>
            <Descriptions>
                <Descriptions.Item label="服务名">{serverInfo.server_name}</Descriptions.Item>
                <Descriptions.Item label="服务号">{serverInfo.server_id}</Descriptions.Item>
                <Descriptions.Item label="服务类型">{getServerType(serverInfo.server_type)}</Descriptions.Item>
                <Descriptions.Item label="服务描述">{serverInfo.desc}</Descriptions.Item>
                <Descriptions.Item label="创建时间">{serverInfo.create_time}</Descriptions.Item>
            </Descriptions>
          </Card>
          <Table
            rowSelection={{
              type: 'checkbox',
              selectedRowKeys,
              onChange: handleSelectChange
            }}
            columns={columns}
            dataSource={nodes}
            rowKey="Host"
            pagination={false}
          />
        </Form>
      </Modal>
    </>

  );
}