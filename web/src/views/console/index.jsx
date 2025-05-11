import React, { useEffect, useState } from 'react';
import { Row, Col, Card, List, Divider, Tree, Button, Modal, Form, InputNumber, Table, Input, message } from 'antd';
import ButtonGroup from 'antd/es/button/button-group';
import {  checkServerNodesStatus, getGroupList, getNodeList, getServerConfigList, getServerInfo, getServerList, getServerNodes,getServerNodesStatus,getStatus, restartServer, stopServer } from './api';
import { getServerNodeStatusType, getServerType } from './constant';
import ResourceModal from './ResourceModal';
import GroupModal from './GroupModal';
import ServerModal from './ServerModal';
import ScaleModal from './ScaleModal';
import DeployModal from './DeployModal';
import { Descriptions } from 'antd';
import _ from 'lodash'
import AddNodeModal from './AddNodeForm';
import UploadConfigModal from './UploadConfigModal';


const ServerNodesColumns = [
    { title: '主机地址', dataIndex: 'host', key: 'host' },
    { title: '端口号', key: 'port', dataIndex:"port" },
    { title: '状态', key: 'server_node_status', dataIndex:"server_node_status",render:(text,record)=>{
        if(record.server_node_status === 1){
            return <span style={{color:'green'}}>online</span>
        }
        if(record.server_node_status === 2){
            return <span style={{color:'red'}}>offline</span>
        }
        return <span style={{color:'black'}}>已删除</span>
    }},
    { title: '版本号', key: 'patch_id', dataIndex:"patch_id" },
    { title: '节点资源限制', key: 'id', dataIndex:"id" ,
        render:(text,record)=>{
            return (
               <>
                    <span>MAX_CPU: {record.cpu_limit} (CORE)</span>
                    <br />
                    <span>MAX_MEMORY: {record.memory_limit} (M)</span>
                    <br/>
               </>
            )
        }
    },
    { title: '创建时间', dataIndex: 'node_create_time', key: 'node_create_time' },

]

const NodeColumns = [
    { title: '主机地址', dataIndex: 'Host', key: 'Host' },
    { title: '状态', key:'NodeStatus', dataIndex:"NodeStatus",render:(text,record)=>{
        if(record.NodeStatus === 1){
            return <span style={{color:'green'}}>online</span>
        }
        if(record.NodeStatus === 2){
            return <span style={{color:'red'}}>offline</span>
        }
    }},
    { title: '节点配置', key: 'ID', dataIndex: 'ID', render: (text, record) => {
        return (
            <div>
                <span>CPU: {record.Cpus} (CORE)</span>
                <br />
                <span>MEMORY: {record.Memory} (G)</span>
            </div>
        )
    } },
    { title: '创建时间', key: 'CreateTime', dataIndex: 'CreateTime' },
]

export default function Console(){
    const [messageApi, contextHolder] = message.useMessage();
    const [resourceModalVisible, setResourceModalVisible] = useState(false);
    const [groupModalVisible, setGroupModalVisible] = useState(false);
    const [serverModalVisible, setServerModalVisible] = useState(false);
    const [scaleModalVisible, setScaleModalVisible] = useState(false);
    const [deployModalVisible, setDeployModalVisible] = useState(false);
    const [serverConfigVisible,setServerConfigVisible] = useState(false);
    const [addNodeVisible, setAddNodeVisible] = useState(false);

    const [form] = Form.useForm();
    const [groupForm] = Form.useForm();
    const [serverForm] = Form.useForm();
    const [groupOptions, setGroupOptions] = useState([]);
    const [selectNodes, setSelectNodes] = useState([])
    const handleRefresh = () => {
        setTimeout(()=>{
            handleTreeNodeClick({
                isGroup:false,
                key: serverInfo.server_id
            })
            messageApi.info('刷新成功');
        },0)
    };
    
    let [serverInfo,setServerInfo] = useState({
        "server_name":"",
        "server_id":"",
        "desc":"",
        "create_time":"",
        "server_type":"",
        "exec_path":""
    })
    let [serverNodes,setServerNodes] = useState([]);
    let [serverConfigList,setServerConfigList] = useState([]);
    const [serverNodePage,setServerNodePage] = useState({
        offset:1,
        size:10
    })
    const [serverNodeStatusList,setServerNodeStatusList] = useState([])
    const [serverNodeStatusTotal,setServerNodeStatusTotal] = useState(0)

    function handleTreeNodeClick(node){
        if(node.isGroup){
            return;
        }
        let serverId = node.key;
        getServerInfo({id:serverId}).then(data => {
            setServerInfo({
                "server_name":data.data.ServerName,
                "server_id":data.data.ID,
                "server_type":data.data.ServerType,
                "exec_path":data.data.ExecFilePath,
                "desc":data.data.Description,
                "create_time":data.data.CreateTime
            })
        })
        getServerNodes({id:serverId}).then(res=>{
            if(!res.data){
                res.data = []
            }
            setServerNodes(res.data);
            console.log('res',res);
            getStatus({
                nodeIds:res.data.map(v=>v.id),
                serverId:serverId
            }).then(res=>{
                console.log('getStatus.res',res);
            })

            checkServerNodesStatus({
                server_id:serverId,
                server_node_ids:res.data.map(v=>v.id),
            }).then(res=>{
                console.log('getStatus.res',res);
            })

        });

        getServerConfigList({serverId:serverId}).then(res=>{
            setServerConfigList(res.data || []);
            console.log('getServerConfigList >> ',res);
        })

        getServerNodesStatus({
            server_id:serverId,
            offset:serverNodePage.offset,
            size:serverNodePage.size
        }).then(res=>{
            setServerNodeStatusTotal(res.data.total)
            setServerNodeStatusList(res.data.list);
        })

    }

    const [nodes,setNodes] = useState([])
    const [treeData, setTreeData] = useState([]);
    function initNodes(){
        getNodeList().then(data => {
            if(data.success){
                setNodes(data.data);
            }
        })
    }
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

    function initServersAndNodes(){
        initServerTreeData()
        initNodes()
    }

    useEffect(() => {
        initServersAndNodes()
    },[])

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

    const [fileName,setFileName] = useState('');
    function handleUpsertConfig(STATE,fileName){
        // 关闭
        if(STATE === -1){
            setServerConfigVisible(false);
            return
        }
        // 创建，不做多的
        if(STATE === 2){
            setServerConfigVisible(true);
            return
        }
        // 查看
        if(STATE == 1){
            setServerConfigVisible(true);
            setFileName(fileName);
        }
    }

    function handleStopServerNodes(){
        if(selectNodes.length === 0){
            messageApi.warning('请选择至少一个节点');
            return
        }
        let nodeIds = selectNodes.map(v=>v.id);
        stopServer({
            nodeIds,
            serverId:serverInfo.server_id
        }).then(res=>{
            if(res.success){
                messageApi.success('停止成功');
                handleRefresh();
            }else{
                messageApi.error(res.msg);
            }
        })
        console.log('nodeIds',nodeIds);
    }

    function handleRestartServerNodes(){
        if(selectNodes.length === 0){
            messageApi.warning('请选择至少一个节点');
            return
        }
        const serverNodeIds = selectNodes.map(v=>v.id);
        const serverId = serverInfo.server_id;
        console.log('nodeIds',serverNodeIds);
        restartServer({
            serverNodeIds,
            serverId,
            packgeId:0,
        }).then(res=>{
            if(res.success){
                messageApi.success('重启成功');
                handleRefresh();
            }else{
                messageApi.error(res.msg);
            }
        })
    }

    const handleSetResourceModalVisible = () => {
        if (selectNodes.length === 0) {
            messageApi.warning('请选择至少一个节点');
            return;
        }
        setResourceModalVisible(true);
    }

    const handleSetDeployModalVisible = () => {
        if (selectNodes.length === 0) {
            messageApi.warning('请选择至少一个节点');
            return;
        }
        setDeployModalVisible(true);
    }
    return (
        <div style={{ padding: 24 }}>
            {contextHolder}
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
                    <Card title="配置文件" extra={
                        <ButtonGroup>
                            <Button onClick={initServersAndNodes}>刷新</Button>
                            <Button onClick={()=>handleUpsertConfig(2)}>上传</Button>
                        </ButtonGroup>
                    }>
                        <List
                            dataSource={serverConfigList}
                            renderItem={item => (
                                <List.Item>
                                    <div>
                                        {item}
                                    </div>
                                    <ButtonGroup size={"small"} style={{float:'right'}}>
                                        <Button onClick={()=>handleUpsertConfig(1,item)} >查看</Button>
                                        <Button danger>删除</Button>
                                    </ButtonGroup>
                                </List.Item>
                            )}
                        />
                    </Card>
                    <Divider />
                    <Card title="节点列表" extra={
                        <>
                            <ButtonGroup>
                                <Button onClick={initServersAndNodes}>刷新</Button>
                                <Button onClick={()=>setAddNodeVisible(true)}>新增节点</Button>
                            </ButtonGroup>
                        </>

                    }>
                        <Table
                            bordered
                            dataSource={nodes}
                            columns={NodeColumns}
                        >
                        </Table>
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
                    <Card title="服务节点列表" variant={false} extra={
                        <div>
                            <Button onClick={handleRefresh}>刷新</Button>
                            <ButtonGroup style={{marginLeft:"16px"}}>
                                <Button onClick={handleSetDeployModalVisible}>部署</Button>
                                <Button onClick={handleRestartServerNodes}>重启</Button>
                                <Button onClick={handleSetResourceModalVisible}>资源配置</Button>
                                <Button onClick={() => setScaleModalVisible(true)}>扩容</Button>
                                <Button onClick={() => handleStopServerNodes()}>停止</Button>
                                <Button onClick={handleRefresh} danger>删除</Button>
                            </ButtonGroup>
                        </div>
                    }>
                        <Table 
                            rowSelection={Object.assign({ type: "checkbox" }, rowSelection)}
                            bordered
                            columns={ServerNodesColumns}
                            dataSource={serverNodes}
                            rowKey="id"
                        />
                    </Card>
                    <Divider />
                    <Card title="服务状态日志" bo={false} extra={<Button onClick={handleRefresh}>刷新</Button>}>
                        <List
                            pagination={{
                                pageSize: 10,
                                total: serverNodeStatusTotal,
                                showSizeChanger: true,
                                onChange: (page, pageSize) => {
                                    let offset = (page - 1) * pageSize;
                                    setServerNodePage({
                                        offset:offset,
                                        size:pageSize
                                    })
                                    getServerNodesStatus({
                                        server_id:serverInfo.server_id,
                                        offset:offset,
                                        size:pageSize
                                    }).then(res=>{
                                        setServerNodeStatusTotal(res.data.total)
                                        setServerNodeStatusList(res.data.list);
                                    })
                                }
                            }}
                            bordered
                            dataSource={serverNodeStatusList}
                            renderItem={item => (
                                <List.Item>
                                    {item.CreateTime}-{getServerNodeStatusType(item.TYPE)}-
                                    <span style={
                                        item.TYPE === 2 ? {color:'red'} : {color:'black'}
                                    } >{item.Content}</span>
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
                serverId={serverInfo.server_id}
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
            <UploadConfigModal 
                visible={serverConfigVisible} 
                serverId={serverInfo.server_id}
                onCancel={()=>{
                    handleUpsertConfig(-1);
                }}
                onOk={()=>{
                    handleUpsertConfig(-1);
                }}
                fileName={fileName}
            >
            </UploadConfigModal>
        </div>
    );
}

