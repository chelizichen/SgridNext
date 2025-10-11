import React, { useEffect, useState, useCallback } from "react";
import {
  Row,
  Col,
  Card,
  List,
  Divider,
  Tree,
  Button,
  Modal,
  Form,
  InputNumber,
  Table,
  Input,
  message,
  Tag
} from "antd";
import { useMediaQuery } from "react-responsive";
import ButtonGroup from "antd/es/button/button-group";
import {
  checkServerNodesStatus,
  deleteServerNode,
  getGroupList,
  getNodeList,
  getServerConfigList,
  getServerInfo,
  getServerList,
  getServerNodes,
  getServerNodesStatus,
  restartServer,
  stopServer,
  updateMachineNodeAlias,
  updateMachineNodeStatus,
  runProbeTask,
  getMainConfig,
} from "./api";
import { getServerNodeStatusType, getServerType } from "./constant";
import ResourceModal from "./ResourceModal";
import GroupModal from "./GroupModal";
import ServerModal from "./ServerModal";
import UpdateServerModal from "./UpdateServerModal";
import ConfigModal from "./ConfigModal";
import ScaleModal from "./ScaleModal";
import DeployModal from "./DeployModal";
import { Descriptions } from "antd";
import _ from "lodash-es";
import AddNodeModal from "./AddNodeForm";
import UploadConfigModal from "./UploadConfigModal";
import HistoryModal from "./HistoryModal";
// import { useNavigate } from 'react-router-dom'

function conversionConf(files) {
  const result = {};
  const baseFiles = files.filter((file) => !file.includes("_"));
  baseFiles.forEach((baseFile) => {
    const baseName = baseFile.split(".")[0]; // 如 "application" 或 "server"
    const timestampFiles = files.filter((file) =>
      file.startsWith(`${baseName}_`),
    );
    result[baseFile] = timestampFiles;
  });
  return result;
}

export default function Console() {
  // 检测是否为移动设备
  const isMobile = useMediaQuery({ maxWidth: 1000 });
  // const nagivate = useNavigate();
  const [messageApi, contextHolder] = message.useMessage();
  const [resourceModalVisible, setResourceModalVisible] = useState(false);
  const [groupModalVisible, setGroupModalVisible] = useState(false);
  const [serverModalVisible, setServerModalVisible] = useState(false);
  const [updateServerVisible, setUpdateServerVisible] = useState(false);
  const [configModalVisible, setConfigModalVisible] = useState(false);
  const [probeLoading, setProbeLoading] = useState(false);
  const [scaleModalVisible, setScaleModalVisible] = useState(false);
  const [deployModalVisible, setDeployModalVisible] = useState(false);
  const [serverConfigVisible, setServerConfigVisible] = useState(false);
  const [addNodeVisible, setAddNodeVisible] = useState(false);
  
  const [form] = Form.useForm();
  const [groupForm] = Form.useForm();
  const [serverForm] = Form.useForm();
  const [groupOptions, setGroupOptions] = useState([]);
  const [selectNodes, setSelectNodes] = useState([]);

  const [serverTotal, setServerTotal] = useState(0);
  const [nodeTotal, setNodeTotal] = useState(0);
  const [mainConfig, setMainConfig] = useState({});
  const NodeColumns = [
    { title: "别名", dataIndex: "Alias", key: "Alias", render: (text, record) => {
      return (
        <div style={{width:"48px"}}>
          {record.Alias}
        </div>
      )
    } },
    { title: "主机地址", dataIndex: "Host", key: "Host" ,render:(text,record)=>{
      return (
        <div style={{width:"120px",cursor:'pointer',color:'#1677ff'}} onClick={()=>toApplicationLog(record.Host,3)}>
          {record.Host}
        </div>
      )
    }},
    {
      title: "状态",
      key: "NodeStatus",
      dataIndex: "NodeStatus",
      render: (text, record) => {
        if (record.NodeStatus === 1) {
          return <span style={{ color: "green" }}>online</span>;
        }
        if (record.NodeStatus === 2) {
          return <span style={{ color: "red" }}>offline</span>;
        }
      },
    },
    {
      title: "节点配置",
      key: "ID",
      dataIndex: "ID",
      render: (text, record) => {
        return (
          <div>
            <span>CPU: {record.Cpus} (CORE)</span>
            <br />
            <span>MEMORY: {record.Memory} (G)</span>
          </div>
        );
      },
    },
    { title: "上报时间", key: "UpdateTime", dataIndex: "UpdateTime"},
    {
      title: "操作",
      key: "action",
      dataIndex: "action",
      render: (text, record) => {
        return (
          <ButtonGroup size="small">
            <Button onClick={()=>handleUpdateNodeAlias(record)}>别名</Button>
            <Button onClick={()=>handleUpdateNodeStatus(record,1)}>上线</Button>
            <Button danger  onClick={()=>handleUpdateNodeStatus(record,2)}>下线</Button>
          </ButtonGroup>
        );
      }
    }
  ];

  const ServerNodesColumns = [
    { title: "主机地址", dataIndex: "host", key: "host",render:(text,record)=>{
        return (
          <div style={{color:"#1677ff"}}>
            {record.host}
          </div>
        )
      } 
    },
    { title: "端口号", key: "port", dataIndex: "port" },
    {
      title: "状态",
      key: "server_node_status",
      dataIndex: "server_node_status",
      render: (text, record) => {
        if (record.server_node_status === 1) {
          return <span style={{ color: "green" }}>online</span>;
        }
        if (record.server_node_status === 2) {
          return <span style={{ color: "red" }}>offline</span>;
        }
        return <span style={{ color: "black" }}>deleted</span>;
      },
    },
    { title: "版本号", key: "patch_id", dataIndex: "patch_id"},
    { title: "运行类型", key: "server_run_type", dataIndex: "server_run_type",  
      render: (text, record) => {
        if(record.server_run_type === 0){
          return "手动重启"
        }
        if(record.server_run_type === 12){
          return "自动重启"
        }
      } 
    },
    {
      title: "节点资源限制",
      key: "id",
      dataIndex: "id",
      render: (text, record) => {
        return (
          <>
            <span>MAX_CPU: {record.cpu_limit} (CORE)</span>
            <br />
            <span>MAX_MEMORY: {record.memory_limit} (M)</span>
            <br />
          </>
        );
      },
    },
    { title: "创建时间", dataIndex: "node_create_time", key: "node_create_time" },
    { title: "操作", dataIndex: "node_create_time", key: "node_create_time",   render: (text, record) => {
      return (
        <>
           <Button size="small" type="primary" onClick={()=>toLogPage(record,serverInfo.server_name,serverInfo.server_id)} style={{marginRight:"16px"}}>日志</Button>
          {
            record.view_page && <Button size="small" onClick={()=>window.open(record.view_page,"_blank")}>预览</Button>
          }
        </>
      );
    }, },
  ];

  function toLogPage(record,serverName,serverId){
    let new_path = `/log?host=${record.host}&serverId=${serverId}&serverName=${serverName}&nodeId=${record.id}&logCategory=1`
    let newPath = location.pathname + "/#" + new_path
    window.open(newPath,"_blank")
  }

  function toApplicationLog(host,type){
    if(type == 2){
      host = 'sgridnext'
    }
    let new_path = `/log?host=${host}&logCategory=${type}&serverId=-25528`
    let newPath = location.pathname + "/#" + new_path
    window.open(newPath,"_blank")
  }
  

  function handleUpdateNodeStatus(record,status) {
    console.log("record",record);
    updateMachineNodeStatus({
      id: record.ID,
      status
    }).then(res=>{
      if(res.success){
        messageApi.success("操作成功");
        initServersAndNodes()
      }else{
        messageApi.error(res.msg);
      }
    })
  }

  const handleRefresh = () => {
    setTimeout(() => {
      handleTreeNodeClick({
        isGroup: false,
        key: serverInfo.server_id,
      });
      messageApi.info("刷新成功");
    }, 0);
  };

  let [serverInfo, setServerInfo] = useState({
    server_name: "",
    server_id: "",
    desc: "",
    create_time: "",
    server_type: "",
    exec_path: "",
    log_path: "",
  });
  let [serverNodes, setServerNodes] = useState([]);
  let [serverConfigList, setServerConfigList] = useState({});
  const [serverNodePage, setServerNodePage] = useState({
    offset: 1,
    size: 10,
  });
  const [serverNodeStatusList, setServerNodeStatusList] = useState([]);
  const [serverNodeStatusTotal, setServerNodeStatusTotal] = useState(0);

  function handleTreeNodeClick(node) {
    if (node.isGroup) {
      return;
    }
    let serverId = node.key;
    // 将上次全选的置空
    console.log('clear select nodes');
    setSelectNodes([]);
    getServerInfo({ id: serverId }).then((data) => {
      setServerInfo({
        server_name: data.data.ServerName,
        server_id: data.data.ID,
        server_type: data.data.ServerType,
        exec_path: data.data.ExecFilePath,
        desc: data.data.Description,
        create_time: data.data.CreateTime,
        log_path: data.data.LogPath,
        docker_name: data.data.DockerName,
        config_path: data.data.ConfigPath,
      });
    });
    getServerNodes({ id: serverId }).then((res) => {
      if (!res.data) {
        res.data = [];
      }
      setServerNodes(res.data);
      console.log("res", res);

      checkServerNodesStatus({
        server_id: serverId,
        server_node_ids: res.data.map((v) => v.id),
      }).then((res) => {
        console.log("getStatus.res", res);
      });
    });
    getServerConfigList({ serverId: serverId }).then((res) => {
      setServerConfigList(conversionConf(res.data || []));
      console.log("getServerConfigList >> ", res);
    });
    getServerNodesStatus({
      server_id: serverId,
      offset: serverNodePage.offset,
      size: serverNodePage.size,
    }).then((res) => {
      setServerNodeStatusTotal(res.data.total);
      setServerNodeStatusList(res.data.list);
    });
  }

  const [nodes, setNodes] = useState([]);
  const [treeData, setTreeData] = useState([]);
  function initNodes() {
    getNodeList().then((data) => {
      if (data.success) {
        setNodeTotal(data.data.length);
        setNodes(data.data);
      }
    });
  }

  const [serverIdToGroupMap, setServerIdToGroupMap] = useState({});

  function initServerTreeData() {
    getServerList().then((data) => {
      setServerTotal(data.data.length);
      console.log("data", data);
      let serverGroup = _.groupBy(data.data, "group_name");
      setServerIdToGroupMap(_.keyBy(data.data, "server_id"))
      let treeStructure = Object.keys(serverGroup).map((groupName) => ({
        title: groupName,
        key: groupName,
        isGroup: true,
        children: serverGroup[groupName].map((server) => ({
          title: server.server_name,
          key: server.server_id,
          isGroup: false,
        })),
      }));
      setTreeData(treeStructure);
      console.log("treeStructure", treeStructure);
    });
  }

  const initServersAndNodes = useCallback(() => {
    initServerTreeData();
    initNodes();
    getMainConfig().then(res=>{
      setMainConfig(res.data);
    });
  }, []);

  useEffect(() => {
    initServersAndNodes();
  }, [initServersAndNodes]);

  const rowSelection = {
    onChange: (selectedRowKeys, selectedRows) => {
      console.log(
        `selectedRowKeys: ${selectedRowKeys}`,
        "selectedRows: ",
        selectedRows,
      );
      setSelectNodes(selectedRows);
    },
    getCheckboxProps: (record) => ({
      name: record.name,
    }),
    selectedRowKeys: selectNodes.map(node => node.id), // 关键修复：同步选中状态
  };

  const [fileName, setFileName] = useState("");
  function handleUpsertConfig(STATE, fileName) {
    // 关闭
    if (STATE === -1) {
      setServerConfigVisible(false);
      return;
    }
    // 创建，不做多的
    if (STATE === 2) {
      setServerConfigVisible(true);
      return;
    }
    // 查看
    if (STATE == 1) {
      setServerConfigVisible(true);
      setFileName(fileName);
    }
  }

  function handleStopServerNodes() {
    if (selectNodes.length === 0) {
      messageApi.warning("请选择至少一个节点");
      return;
    }
    let nodeIds = selectNodes.map((v) => v.id);
    stopServer({
      nodeIds,
      serverId: serverInfo.server_id,
    }).then((res) => {
      if (res.success) {
        messageApi.success("停止成功");
        handleRefresh();
      } else {
        messageApi.error(res.msg);
      }
    });
    console.log("nodeIds", nodeIds);
  }

  function handleRestartServerNodes() {
    if (selectNodes.length === 0) {
      messageApi.warning("请选择至少一个节点");
      return;
    }
    const serverNodeIds = selectNodes.map((v) => v.id);
    const serverId = serverInfo.server_id;
    console.log("nodeIds", serverNodeIds);
    restartServer({
      serverNodeIds,
      serverId,
      packgeId: 0,
    }).then((res) => {
      if (res.success) {
        messageApi.success("重启成功");
        handleRefresh();
      } else {
        messageApi.error(res.msg);
      }
    });
  }

  // 运行探针任务
  const handleRunProbe = async () => {
    setProbeLoading(true);
    try {
      const response = await runProbeTask();
      if (response.success) {
        messageApi.success(response.msg);
        // 刷新节点列表
        initNodes();
      } else {
        messageApi.error(response.msg);
      }
    } catch (error) {
      messageApi.error('运行探针失败');
      console.error('运行探针错误:', error);
    } finally {
      setProbeLoading(false);
    }
  };

  const handleSetResourceModalVisible = () => {
    if (selectNodes.length === 0) {
      messageApi.warning("请选择至少一个节点");
      return;
    }
    setResourceModalVisible(true);
  };

  const handleSetDeployModalVisible = () => {
    if (selectNodes.length === 0) {
      messageApi.warning("请选择至少一个节点");
      return;
    }
    setDeployModalVisible(true);
  };

  const [deleteModalVisible, setDeleteModalVisible] = useState(false);
  const handleDeleteServerNodes = (type)=>{
    if(type == 1){
      setDeleteModalVisible(true)
    }
    if(type == 2){
      let body = {
        ids: selectNodes.map(v=>v.id)
      }
      deleteServerNode(body).then(res=>{ 
        if(res.success){
          messageApi.success("删除成功");
          getServerNodes();
        }else{
          messageApi.error(res.msg);
        }
      })
    }

  }

  const [historyModelVisible, setHistoryModelVisible] = useState(false);
  const [historyData, setHistoryData] = useState([]);
  const handleCheckhistory = (fileName) => {
    setHistoryData(serverConfigList[fileName]);
    setHistoryModelVisible(true);
  };
  function switchType(type) {
    switch (type) {
      case 1:
        return { color: "green" };
      case 2:
        return { color: "red" };
      default:
        return { color: "black" };
    }
  }

  function exportSgridReleaseConf(){
    let apiPath = window.location.origin
    let content = `
    # SgridNext 自动发布配置
    # 执行 sgridnext 从 当前目录寻找 sgridnext.release 文件, 读取配置项, 进行发布

    # 发布地址
    DEPLOY_PATH = ${apiPath}
    
    # 服务ID
    SERVER_ID = ${serverInfo.server_id}
    
    # 服务名
    SERVER_NAME = ${serverInfo.server_name}

    # 包地址 此项需要手动填写
    # 例如：/archive/SgridTestJavaServer.tar.gz
    PACKAGE_PATH = /archive/${serverInfo.server_name}.tar.gz

    `
    let blob = new Blob([content], { type: "text/plain" });
    let url = URL.createObjectURL(blob);
    let a = document.createElement("a");
    a.href = url;
    a.download = "sgridnext.release";
    a.click();
    URL.revokeObjectURL(url);
  }

  const [updateNodeAliasVisible, setUpdateNodeAliasVisible] = useState(false);
  const [updateNodeAliasForm] = Form.useForm();
  const handleUpdateNodeAlias = (record) => {
    setUpdateNodeAliasVisible(true);
    updateNodeAliasForm.resetFields();
    console.log('record',record);
    updateNodeAliasForm.setFieldValue("id",record.ID)
    updateNodeAliasForm.setFieldValue("host",record.Host)
  }
  const handleUpdateNodeAliasOk = () => {
    updateNodeAliasForm.validateFields().then((values) => {
      updateMachineNodeAlias({
        id: updateNodeAliasForm.getFieldValue("id"),
        alias: values.alias,
      }).then(res=>{
        if(res.success){
          messageApi.success("更新别名成功");
          setUpdateNodeAliasVisible(false);
          initServersAndNodes();
        }else{
          messageApi.error(res.msg);
        }
      })
    });
  }

  return (
    <div style={{ padding: 24 }}>
      {contextHolder}
      <Row gutter={isMobile ? 8 : 16}>
        <Col span={isMobile ? 24 : 8}>
          <Card
            title="主控信息"
            variant={false}
            extra={
              <ButtonGroup
                style={
                  isMobile ? { display: "flex", flexDirection: "column" } : {}
                }
              >
                <Button onClick={() => setConfigModalVisible(true)}>
                  配置管理
                </Button>
                <Button onClick={()=>toApplicationLog('',2)}>
                    日志管理
                </Button>
                <Button 
                  type="primary" 
                  onClick={handleRunProbe}
                  loading={probeLoading}
                >
                  运行探针
                </Button>
              </ButtonGroup>
            }
          > 
          <Descriptions column={4}>
            <Descriptions.Item  span={2} label="服务总数">
              {serverTotal}
            </Descriptions.Item>
            <Descriptions.Item span={2} label="节点总数">
              {nodeTotal}
            </Descriptions.Item>
            <Descriptions.Item span={4} label="主控主机"> 
              {mainConfig.config?.host}
            </Descriptions.Item>
            <Descriptions.Item span={4} label="主控配置文件">
              {mainConfig.configPath}
            </Descriptions.Item>
          </Descriptions>
          </Card>
          <Divider />
          <Card
            title="服务总揽"
            variant={false}
            extra={
              <ButtonGroup
                style={
                  isMobile ? { display: "flex", flexDirection: "column" } : {}
                }
              >
                <Button
                  onClick={handleRefresh}
                  style={{
                    marginLeft: isMobile ? 0 : "16px",
                    marginBottom: isMobile ? "8px" : 0,
                  }}
                >
                  刷新
                </Button>
                <Button
                  onClick={() => setGroupModalVisible(true)}
                  style={{ marginBottom: isMobile ? "8px" : 0 }}
                >
                  创建组
                </Button>
                <Button onClick={() => setServerModalVisible(true)}>
                  添加服务
                </Button>
              </ButtonGroup>
            }
          >
            <Tree
              treeData={treeData}
              onSelect={(keys, { node }) => handleTreeNodeClick(node)}
            />
          </Card>
          <Divider />
          <Card
            title="配置文件"
            extra={
              <ButtonGroup
                style={
                  isMobile ? { display: "flex", flexDirection: "column" } : {}
                }
              >
                <Button
                  onClick={initServersAndNodes}
                  style={{ marginBottom: isMobile ? "8px" : 0 }}
                >
                  刷新
                </Button>
                <Button onClick={() => handleUpsertConfig(2)}>上传</Button>
              </ButtonGroup>
            }
          >
            <List
              dataSource={Object.keys(serverConfigList)}
              renderItem={(item) => (
                <List.Item>
                  <div>{item}</div>
                  <ButtonGroup size={"small"} style={{ float: "right" }}>
                    <Button onClick={() => handleUpsertConfig(1, item)}>
                      查看
                    </Button>
                    <Button onClick={() => handleCheckhistory(item)}>
                      查看历史记录
                    </Button>
                    {/* <Button danger>删除</Button> */}
                  </ButtonGroup>
                </List.Item>
              )}
            />
          </Card>
          <Divider />
          <Card
            title="节点列表"
            extra={
              <>
                <ButtonGroup
                  style={
                    isMobile ? { display: "flex", flexDirection: "column" } : {}
                  }
                >
                  <Button
                    onClick={initServersAndNodes}
                    style={{ marginBottom: isMobile ? "8px" : 0 }}
                  >
                    刷新
                  </Button>
                  <Button onClick={() => setAddNodeVisible(true)}>
                    新增节点
                  </Button>
                </ButtonGroup>
              </>
            }
          >
            <Table
              bordered
              dataSource={nodes}
              columns={NodeColumns}
              scroll={isMobile ? { x: "max-content" } : undefined}
              pagination={isMobile ? { pageSize: 5 } : undefined}
              style={{"overflowX": "scroll"}}
            ></Table>
          </Card>
        </Col>
        <Col span={isMobile ? 24 : 16}>
          <Card
            title="服务信息"
            variant={true}
            extra={
            <ButtonGroup> 
              <Button onClick={()=>{
                if (!serverInfo.server_id) {
                  messageApi.warning("请先选择一个服务");
                  return;
                }
                setUpdateServerVisible(true);
              }}>更新服务</Button>
              <Button onClick={handleRefresh}>刷新</Button>
              <Button onClick={exportSgridReleaseConf} disabled={!serverInfo.server_id}>导出配置</Button>
            </ButtonGroup>
            }
          >
            {serverInfo ? (
              <Descriptions>
                <Descriptions.Item label="服务名">
                  {serverInfo.server_name}
                </Descriptions.Item>
                <Descriptions.Item label="docker名称">
                  {serverInfo.docker_name}
                </Descriptions.Item>
                <Descriptions.Item label="服务号">
                  {serverInfo.server_id}
                </Descriptions.Item>
                <Descriptions.Item label="服务组" >
                  {serverIdToGroupMap[serverInfo.server_id]?.group_name} 
                </Descriptions.Item>
                <Descriptions.Item label="服务类型">
                  {getServerType(serverInfo.server_type)}
                </Descriptions.Item>
                <Descriptions.Item label="服务描述">
                  {serverInfo.desc}
                </Descriptions.Item>
                <Descriptions.Item label="创建时间">
                  {serverInfo.create_time}
                </Descriptions.Item>
                <Descriptions.Item label="执行地址">
                  {serverInfo.exec_path}
                </Descriptions.Item>
                <Descriptions.Item label="日志地址">
                  {
                  serverInfo.log_path || (
                  serverInfo.server_id ? 
                   "${cwd}/server/SgridPatchServer/log/"+serverInfo.server_name
                    :
                   ''
                  )}
                </Descriptions.Item>
                <Descriptions.Item label="配置文件地址">
                  {
                  serverInfo.config_path || (
                  serverInfo.server_id ? 
                   "${cwd}/server/SgridPatchServer/config/"+serverInfo.server_name
                    :
                   ''
                  )}
                </Descriptions.Item>
              </Descriptions>
            ) : (
              <p>请从左侧选择节点</p>
            )}
          </Card>
          <Divider />
          <Card
            title="服务节点列表"
            variant={false}
            extra={
              <div
                style={
                  isMobile ? { display: "flex", flexDirection: "column" } : {}
                }
              >
                <Button
                  onClick={handleRefresh}
                  style={{ marginBottom: isMobile ? "8px" : 0 }}
                >
                  刷新
                </Button>
                <ButtonGroup
                  style={
                    isMobile
                      ? {
                          display: "flex",
                          flexDirection: "column",
                          width: "100%",
                        }
                      : { marginLeft: "16px" }
                  }
                >
                  <Button
                    onClick={handleSetDeployModalVisible}
                    style={{ marginBottom: isMobile ? "8px" : 0 }}
                  >
                    部署
                  </Button>
                  <Button
                    onClick={handleRestartServerNodes}
                    style={{ marginBottom: isMobile ? "8px" : 0 }}
                  >
                    重启
                  </Button>
                  <Button
                    onClick={handleSetResourceModalVisible}
                    style={{ marginBottom: isMobile ? "8px" : 0 }}
                  >
                    资源配置
                  </Button>
                  <Button
                    onClick={() => setScaleModalVisible(true)}
                    style={{ marginBottom: isMobile ? "8px" : 0 }}
                  >
                    扩容
                  </Button>
                  <Button
                    onClick={() => handleStopServerNodes()}
                    style={{ marginBottom: isMobile ? "8px" : 0 }}
                  >
                    停止
                  </Button>
                  <Button onClick={()=>handleDeleteServerNodes(1)} danger>
                    删除
                  </Button>
                </ButtonGroup>
              </div>
            }
          >
            <Table
              rowSelection={Object.assign({ type: "checkbox" }, rowSelection)}
              bordered
              columns={ServerNodesColumns}
              dataSource={serverNodes}
              rowKey="id"
              scroll={isMobile ? { x: "max-content" } : undefined}
              pagination={isMobile ? { pageSize: 5 } : undefined}
            />
          </Card>
          <Divider />
          <Card
            title="服务状态日志"
            bo={false}
            extra={
              <Button
                onClick={handleRefresh}
                style={isMobile ? { width: "100%" } : {}}
              >
                刷新
              </Button>
            }
          >
            <List
              pagination={{
                pageSize: 10,
                total: serverNodeStatusTotal,
                showSizeChanger: true,
                onChange: (page, pageSize) => {
                  let offset = (page - 1) * pageSize;
                  setServerNodePage({
                    offset: offset,
                    size: pageSize,
                  });
                  getServerNodesStatus({
                    server_id: serverInfo.server_id,
                    offset: offset,
                    size: pageSize,
                  }).then((res) => {
                    setServerNodeStatusTotal(res.data.total);
                    setServerNodeStatusList(res.data.list);
                  });
                },
              }}
              bordered
              dataSource={serverNodeStatusList}
              renderItem={(item) => (
                <List.Item>
                  {item.CreateTime}-{getServerNodeStatusType(item.TYPE)}-
                  <span style={switchType(item.TYPE)}>{item.Content}</span>
                </List.Item>
              )}
            />
          </Card>
        </Col>
      </Row>
      <ResourceModal
        visible={resourceModalVisible}
        onOk={(values) => {
          console.log("资源配置:", values);
          setResourceModalVisible(false);
        }}
        onCancel={() => setResourceModalVisible(false)}
        form={form}
        nodes={selectNodes}
        serverId={serverInfo.server_id}
      />
      <GroupModal
        visible={groupModalVisible}
        onOk={() => {
          setGroupModalVisible(false);
          handleRefresh();
        }}
        onCancel={() => setGroupModalVisible(false)}
        form={groupForm}
      />
      <ServerModal
        visible={serverModalVisible}
        onOk={() => {
          setServerModalVisible(false);
          handleRefresh();
        }}
        onCancel={() => setServerModalVisible(false)}
        afterOpenChange={(open) => {
          if (open) {
            getGroupList().then((data) => {
              console.log("data", data);

              if (data.success) {
                setGroupOptions(
                  data.data.map((item) => ({
                    label: item.Name,
                    value: item.ID,
                  })),
                );
              }
            });
          }
        }}
        form={serverForm}
        groupOptions={groupOptions}
      />
      <ScaleModal
        visible={scaleModalVisible}
        onOk={() => {
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
      <UpdateServerModal
        visible={updateServerVisible}
        onOk={()=>{
          setUpdateServerVisible(false);
          handleRefresh();
        }}
        onCancel={()=>setUpdateServerVisible(false)}
        serverInfo={serverInfo}
        groupOptions={groupOptions}
      />
      <ConfigModal
        visible={configModalVisible}
        onOk={()=>{
          setConfigModalVisible(false);
        }}
        onCancel={()=>setConfigModalVisible(false)}
      />
      <AddNodeModal
        visible={addNodeVisible}
        onOk={() => {
          setAddNodeVisible(false);
        }}
        onCancel={() => setAddNodeVisible(false)}
        form={form}
      ></AddNodeModal>
      <UploadConfigModal
        visible={serverConfigVisible}
        serverId={serverInfo.server_id}
        onCancel={() => {
          handleUpsertConfig(-1);
        }}
        onOk={() => {
          handleUpsertConfig(-1);
        }}
        fileName={fileName}
      ></UploadConfigModal>
      <HistoryModal
        visible={historyModelVisible}
        onCancel={() => setHistoryModelVisible(false)}
        historyData={historyData}
        checkFileCb={(fileName) => {
          handleUpsertConfig(1, fileName);
        }}
      ></HistoryModal>

      <Modal
        title="确认删除以下服务节点吗"
        open={deleteModalVisible}
        onOk={()=>handleDeleteServerNodes(2)}
        onCancel={() => setDeleteModalVisible(false)}
      >
        <div>
          {selectNodes.map((item) => (
            <div>
              <Tag color="error" bordered>
                {item.host}:{item.port}
              </Tag>
            </div>
          ))}
        </div>
      </Modal>
      <Modal
        title="更新节点别名"
        open={updateNodeAliasVisible}
        onOk={handleUpdateNodeAliasOk}
        onCancel={() => setUpdateNodeAliasVisible(false)}
      >
        <Form form={updateNodeAliasForm} layout="vertical">
          <Form.Item name="host" label="主机地址">
            <Input disabled style={{width:'100%'}} value={updateNodeAliasForm.getFieldValue("host")} />
          </Form.Item>
          <Form.Item name="alias" label="别名" rules={[{ required: true, message: '请输入别名' }]}>
            <Input style={{width:'100%'}} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
