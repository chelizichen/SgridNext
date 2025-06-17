import React, { useEffect, useState } from "react";
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
} from "antd";
import { useMediaQuery } from "react-responsive";
import ButtonGroup from "antd/es/button/button-group";
import {
  checkServerNodesStatus,
  getGroupList,
  getNodeList,
  getServerConfigList,
  getServerInfo,
  getServerList,
  getServerNodes,
  getServerNodesStatus,
  restartServer,
  stopServer,
  updateMachineNodeStatus,
} from "./api";
import { getServerNodeStatusType, getServerType } from "./constant";
import ResourceModal from "./ResourceModal";
import GroupModal from "./GroupModal";
import ServerModal from "./ServerModal";
import ScaleModal from "./ScaleModal";
import DeployModal from "./DeployModal";
import { Descriptions } from "antd";
import _ from "lodash-es";
import AddNodeModal from "./AddNodeForm";
import UploadConfigModal from "./UploadConfigModal";
import HistoryModal from "./HistoryModal";
import { useNavigate } from 'react-router-dom'

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
  const nagivate = useNavigate();
  const [messageApi, contextHolder] = message.useMessage();
  const [resourceModalVisible, setResourceModalVisible] = useState(false);
  const [groupModalVisible, setGroupModalVisible] = useState(false);
  const [serverModalVisible, setServerModalVisible] = useState(false);
  const [scaleModalVisible, setScaleModalVisible] = useState(false);
  const [deployModalVisible, setDeployModalVisible] = useState(false);
  const [serverConfigVisible, setServerConfigVisible] = useState(false);
  const [addNodeVisible, setAddNodeVisible] = useState(false);

  const [form] = Form.useForm();
  const [groupForm] = Form.useForm();
  const [serverForm] = Form.useForm();
  const [groupOptions, setGroupOptions] = useState([]);
  const [selectNodes, setSelectNodes] = useState([]);

  const NodeColumns = [
    { title: "主机地址", dataIndex: "Host", key: "Host" },
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
    { title: "创建时间", key: "CreateTime", dataIndex: "CreateTime" },
    {
      title: "操作",
      key: "action",
      dataIndex: "action",
      render: (text, record) => {
        return (
          <ButtonGroup size="middle">
            <Button  onClick={()=>handleUpdateNodeStatus(record,1)}>上线</Button>
            <Button danger  onClick={()=>handleUpdateNodeStatus(record,2)}>下线</Button>
          </ButtonGroup>
        );
      }
    }
  ];

  const ServerNodesColumns = [
    { title: "主机地址", dataIndex: "host", key: "host",render:(text,record)=>{
        return (
          <div style={{"cursor":"pointer",color:"#1677ff"}} onClick={()=>toLogPage(record,serverInfo.server_name,serverInfo.server_id)}>
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
        return <span style={{ color: "black" }}>已删除</span>;
      },
    },
    { title: "版本号", key: "patch_id", dataIndex: "patch_id",render:(text,record)=>toViewPage(record) },
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
  ];

  function toLogPage(record,serverName,serverId){
    let new_path = `/log?host=${record.host}&serverId=${serverId}&serverName=${serverName}&nodeId=${record.id}`
    let newPath = location.pathname + "/#" + new_path
    window.open(newPath,"_blank")
  }

  function toViewPage(record){
    return (
      <div  style={{"cursor":"pointer",color:"#1677ff"}} onClick={()=>window.open(record.view_page,"_blank")}>{record.patch_id}</div>
    )
  }

  

  function handleUpdateNodeStatus(record,status) {
    console.log("record",record);
    updateMachineNodeStatus({
      id: record.ID,
      status
    }).then(res=>{
      if(res.success){
        messageApi.success("操作成功");
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
    getServerInfo({ id: serverId }).then((data) => {
      setServerInfo({
        server_name: data.data.ServerName,
        server_id: data.data.ID,
        server_type: data.data.ServerType,
        exec_path: data.data.ExecFilePath,
        desc: data.data.Description,
        create_time: data.data.CreateTime,
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
        setNodes(data.data);
      }
    });
  }
  function initServerTreeData() {
    getServerList().then((data) => {
      console.log("data", data);
      let serverGroup = _.groupBy(data.data, "group_name");
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

  function initServersAndNodes() {
    initServerTreeData();
    initNodes();
  }

  useEffect(() => {
    initServersAndNodes();
  }, []);

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
      disabled: record.name === "Disabled User", // Column configuration not to be checked
      name: record.name,
    }),
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
  return (
    <div style={{ padding: 24 }}>
      {contextHolder}
      <Row gutter={isMobile ? 8 : 16}>
        <Col span={isMobile ? 24 : 8}>
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
            extra={<Button onClick={handleRefresh}>刷新</Button>}
          >
            {serverInfo ? (
              <Descriptions>
                <Descriptions.Item label="服务名">
                  {serverInfo.server_name}
                </Descriptions.Item>
                <Descriptions.Item label="服务号">
                  {serverInfo.server_id}
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
                  <Button onClick={handleRefresh} danger>
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
        onOk={(values) => {
          console.log("创建组:", values);
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
          console.log("添加节点:", values);
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
    </div>
  );
}
