import React, { useEffect, useState } from 'react';
import { Row, Col, Card, List, Divider, Tree, Button, Modal, Form, InputNumber, Table, Input, message } from 'antd';
import ButtonGroup from 'antd/es/button/button-group';
import {  getGroupList, getNodeList, getServerInfo, getServerList, getServerNodes,getServerType } from './api';
import ResourceModal from './ResourceModal';
import GroupModal from './GroupModal';
import ServerModal from './ServerModal';
import ScaleModal from './ScaleModal';
import DeployModal from './DeployModal';
import { Descriptions } from 'antd';
import _ from 'lodash'
import AddNodeModal from './AddNodeForm';
export default function Console(){
    const [resourceModalVisible, setResourceModalVisible] = useState(false);
    const [form] = Form.useForm();
    const [groupModalVisible, setGroupModalVisible] = useState(false);
    const [serverModalVisible, setServerModalVisible] = useState(false);
    const [scaleModalVisible, setScaleModalVisible] = useState(false);
    const [deployModalVisible, setDeployModalVisible] = useState(false);
    const [addNodeVisible, setAddNodeVisible] = useState(false);
    const [groupForm] = Form.useForm();
    const [serverForm] = Form.useForm();
    const [groupOptions, setGroupOptions] = useState([]);
    const [selectNodes, setSelectNodes] = useState([])
    const handleRefresh = () => {
        // 这里可以添加实际的刷新逻辑
        console.log('刷新数据');
    };
    
    const [treeData, setTreeData] = useState([]);
    function initServerTreeData() {
        getServerList().then(data => {
            console.log('data',data);
            let serverGroup = _.groupBy(data.data, 'group_name');
            let treeStructure = Object.keys(serverGroup).map(groupName => ({
                title: groupName,
                key: groupName,
                isGroup: true,
                children: serverGroup[groupName].map(server => ({
                    title: server.server_name,
                    key: server.server_id,
                    isGroup: false
                }))
            }));
            setTreeData(treeStructure);
            console.log('treeStructure', treeStructure);
        });
    }
    let [serverInfo,setServerInfo] = useState({
        "server_name":"",
        "server_id":"",
        "desc":"",
        "create_time":"",
        "server_type":"",
        "exec_path":""
    })
    let [serverNodes,setServerNodes] = useState([]);
    // 1. 拿服务信息
    // 2. 拿节点信息
    function handleTreeNodeClick(node){
        if(node.isGroup){
            return;
        }
        getServerInfo({id:node.key}).then(data => {
            setServerInfo({
                "server_name":data.data.ServerName,
                "server_id":data.data.ID,
                "server_type":data.data.ServerType,
                "exec_path":data.data.ExecFilePath,
                "desc":data.data.Description,
                "create_time":data.data.CreateTime
            })
        })
        getServerNodes({id:node.key}).then(res=>{
            if(!res.data){
                res.data = []
            }
            setServerNodes(res.data);
            console.log('res',res);
        });
    }

    useEffect(() => {
        initServerTreeData()
    },[])

    const [nodes,setNodes] = useState([])
    function initNodes(){
        getNodeList().then(data => {
            if(data.success){
                setNodes(data.data);
            }
        })
    }
    useEffect(() => {
        initNodes()
    },[])

    // 示例日志数据
    const logData = [
        { time: '10:30', status: '运行中', message: '节点1启动完成' },
        { time: '10:25', status: '警告', message: '节点2负载过高' },
        { time: '10:20', status: '正常', message: '节点3连接成功' },
    ];
    const rowSelection = {
        onChange: (selectedRowKeys, selectedRows) => {
          console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
          setSelectNodes(selectedRows);
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
                    <Card title="服务总揽"
                        variant={false} extra={
                        <ButtonGroup>
                            <Button onClick={handleRefresh}  style={{marginLeft:'16px'}}>刷新</Button>
                            <Button onClick={()=>setGroupModalVisible(true)}>创建组</Button>
                            <Button onClick={()=>setServerModalVisible(true)}>添加服务</Button>
                        </ButtonGroup>
                    }>
                        <Tree
                            treeData={treeData}
                            onSelect={(keys, { node }) => handleTreeNodeClick(node)}
                        />
                    </Card>
                    <Divider />
                    <Card title="配置文件" extra={<Button onClick={handleRefresh}>刷新</Button>}>
                        <List
                            dataSource={logData}
                            renderItem={item => (
                                <List.Item>
                                    [{item.time}] [{item.status}] {item.message}
                                </List.Item>
                            )}
                        />
                    </Card>
                    <Divider />
                    <Card title="节点状态" extra={
                        <>
                            <ButtonGroup>
                                <Button onClick={handleRefresh}>刷新</Button>
                                <Button onClick={()=>setAddNodeVisible(true)}>新增节点</Button>
                            </ButtonGroup>
                        </>

                    }>
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
                        {serverInfo ? (
                            <Descriptions>
                                <Descriptions.Item label="服务名">{serverInfo.server_name}</Descriptions.Item>
                                <Descriptions.Item label="服务号">{serverInfo.server_id}</Descriptions.Item>
                                <Descriptions.Item label="服务类型">{getServerType(serverInfo.server_type)}</Descriptions.Item>
                                <Descriptions.Item label="服务描述">{serverInfo.desc}</Descriptions.Item>
                                <Descriptions.Item label="创建时间">{serverInfo.create_time}</Descriptions.Item>
                            </Descriptions>
                        ) : (
                            <p>请从左侧选择节点</p>
                        )}
                    </Card>
                    <Divider />
                    <Card title="节点列表" variant={false} extra={
                        <div>
                            <Button onClick={handleRefresh}>刷新</Button>
                            <ButtonGroup style={{marginLeft:"16px"}}>
                                <Button onClick={() => {
                                    if (serverInfo) {
                                        setDeployModalVisible(true);
                                    } else {
                                        console.log('111');
                                        
                                        message.warning('请先选择至少一个节点');
                                    }
                                }}>部署</Button>
                                <Button onClick={() => setResourceModalVisible(true)}>资源配置</Button>
                                <Button onClick={() => setScaleModalVisible(true)}>扩容</Button>
                                <Button onClick={handleRefresh}>删除</Button>
                            </ButtonGroup>
                        </div>

                    }>
                        <Table 
                            rowSelection={Object.assign({ type: "checkbox" }, rowSelection)}
                            columns={[
                                { title: '主机地址', dataIndex: 'host', key: 'host' },
                                { title: '端口号', key: 'port', dataIndex:"port" },
                                { title: '创建时间', dataIndex: 'node_create_time', key: 'node_create_time' },
                                { title: '版本号', key: 'patch_id', dataIndex:"patch_id" },
                                { title: '端口号', key: 'port', dataIndex:"port" },
                            ]}
                            dataSource={serverNodes}
                            rowKey="id"
                            pagination={false}
                        />
                    </Card>
                    <Divider />
                    <Card title="服务状态日志" bo={false} extra={<Button onClick={handleRefresh}>刷新</Button>}>
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
            <ResourceModal
                visible={resourceModalVisible}
                onOk={(values) => {
                    console.log('资源配置:', values);
                    setResourceModalVisible(false);
                }}
                onCancel={() => setResourceModalVisible(false)}
                form={form}
                nodes={selectNodes}
            />
            <GroupModal
                visible={groupModalVisible}
                onOk={(values) => {
                    console.log('创建组:', values);
                    setGroupModalVisible(false);
                    handleRefresh();
                }}
                onCancel={() => setGroupModalVisible(false)}
                form={groupForm}
            />
            <ServerModal
                visible={serverModalVisible}
                onOk={(values) => {
                    setServerModalVisible(false);
                    handleRefresh();
                }}
                onCancel={() => setServerModalVisible(false)}
                afterOpenChange={(open) => {
                    if (open) {
                        getGroupList().then(data => {
                            console.log('data',data);
                            
                            if (data.success) {
                                setGroupOptions(data.data.map(item => ({ label: item.Name, value: item.ID })));
                            }
                        });
                    }
                }}
                form={serverForm}
                groupOptions={groupOptions}
            />
            <ScaleModal
                visible={scaleModalVisible}
                onOk={(selectedNodes) => {
                    setScaleModalVisible(false);
                    handleRefresh();
                }}
                onCancel={() => setScaleModalVisible(false)}
                nodes={nodes}
                serverInfo={serverInfo}
            />
            <DeployModal
                visible={deployModalVisible}
                onOk={() => setDeployModalVisible(false)}
                onCancel={() => setDeployModalVisible(false)}
                nodes={serverNodes}
                serverInfo={serverInfo}
            />
            <AddNodeModal
                visible={addNodeVisible}
                onOk={(values) => {
                    console.log('添加节点:', values);
                    setAddNodeVisible(false);
                }}
                onCancel={() => setAddNodeVisible(false)}
                form={form}
            ></AddNodeModal>
        </div>
    );
}

