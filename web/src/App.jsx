import { Routes, Route, useNavigate } from 'react-router-dom';
import { Layout, Menu } from 'antd';
import './App.css';
import Console from './views/console';
import Control from './views/control';
import Login from './views/Login';

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
            <Menu.Item key="1">SgridNext</Menu.Item>
            <Menu.Item key="2" onClick={() => navigate('/console')}>控制台</Menu.Item>
            <Menu.Item key="3" onClick={() => navigate('/control')}>管理中心</Menu.Item>
          </Menu>
        </Header>
        <Content>
          <Routes>
            <Route path="/console" element={<Console />} />
            <Route path="/control" element={<Control />} />
            <Route path="/login" element={<Login />} />
          </Routes>
        </Content>
        <Footer style={{ textAlign: 'center' }}>
          SgridNext ©2025 Created by chelizichen
        </Footer>
      </Layout>
  );
}

