import { Routes, Route, useNavigate, Outlet } from 'react-router-dom';
import { Layout, Menu } from 'antd';
import './App.css';
import Console from './views/console';
import Control from './views/control';
import Login from './views/Login';
import LogView from './views/log/LogView';
import WebShell from './views/webshell/WebShell';
import DocumentManager from './views/document/DocumentManager';
import NodeManager from './views/control/components/NodeManager';
import ServerList from './views/control/components/ServerList';
import ConfigList from './views/control/components/ConfigList';
import NodeStat from './views/control/components/NodeStat';
import ProbeList from './views/control/components/ProbeList';

const { Header, Content, Footer } = Layout;

export default function App() {
  const navigate = useNavigate();
  const isLoggedIn = localStorage.getItem('isLoggedIn');

  if (!isLoggedIn) {
    navigate('/login');
  }
  
  return (
      <Layout className="layout">
        <Header >
          <Menu
            theme="dark"
            mode="horizontal"
            defaultSelectedKeys={['2']}
          >
            <Menu.Item key="1">
              <div style={{display:"flex",alignItems:"center"}}>
                <img src="sgridcloud.png" style={{width:"50px",height:"50px","borderRadius":"100%","marginRight":"20px"}}></img>
                <div>SgridNext</div>
              </div>
            </Menu.Item>
            <Menu.Item key="2" onClick={() => navigate('/console')}>控制台</Menu.Item>
            <Menu.Item key="3" onClick={() => navigate('/control')}>管理中心</Menu.Item>
            <Menu.Item key="4" onClick={() => navigate('/webshell')}>WebShell</Menu.Item>
            <Menu.Item key="5" onClick={() => navigate('/document')}>文档管理</Menu.Item>
          </Menu>
        </Header>
        <Content>
          <Routes>
            <Route path="/console" element={<Console />} />
            <Route path="/control" element={<Control />} >
              <Route path="server_list" element={<ServerList />} />
              <Route path="machine_list" element={<NodeManager />} />
              <Route path="config_list" element={<ConfigList />} />
              <Route path="nodestat_list" element={<NodeStat />} />
              <Route path="probe_list" element={<ProbeList />} />
            </Route>
            <Route path="/login" element={<Login />} />
            <Route path="/log" element={<LogView />} />
            <Route path="/webshell" element={<WebShell />} />
            <Route path="/document" element={<DocumentManager />} />
          </Routes>
        </Content>
        <Footer style={{ textAlign: 'center' }}>
          SgridNext ©2025 Created by chelizichen
        </Footer>
      </Layout>
  );
}

