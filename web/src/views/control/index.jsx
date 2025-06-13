import NodeManager from "./components/NodeManager";
import { Routes, Route, useNavigate, Outlet } from "react-router-dom";

import {
  AppstoreOutlined,
  ContainerOutlined,
  DesktopOutlined,
  MailOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  PieChartOutlined,
} from "@ant-design/icons";
import { Button, Menu } from "antd";
const items = [
  { key: "server_list", icon: <PieChartOutlined />, label: "服务列表" },
  { key: "machine_list", icon: <DesktopOutlined />, label: "节点列表" },
  { key: "nodestat_list", icon: <DesktopOutlined />, label: "节点状态列表" },
  // {
  //   key: 'sub1',
  //   label: 'Navigation One',
  //   icon: <MailOutlined />,
  //   children: [
  //     { key: '5', label: 'Option 5' },
  //     { key: '6', label: 'Option 6' },
  //     { key: '7', label: 'Option 7' },
  //     { key: '8', label: 'Option 8' },
  //   ],
  // },
];
const Control = () => {
  const nagivate = useNavigate();
  function handleClick(e) {
    console.log("click ", e);
    nagivate({
      pathname: e.key,
    });
  }

  return (
    <div style={{ display: "flex" }}>
      <div style={{ width: 218 }}>
        <Menu
          mode="inline"
          theme="dark"
          items={items}
          style={{ height: "95vh",overflow:'scroll' }}
          onClick={(e) => {
            handleClick(e);
          }}
        />
      </div>
      <div style={{ flex: 1, padding: 20 }}>
        <Outlet />
      </div>
    </div>
  );
};

export default Control;
