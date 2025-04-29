import React, { useState } from 'react';
import { Form, Input, Button, message } from 'antd';
import axios from 'axios';
import { Card } from 'antd';
import { useNavigate } from 'react-router-dom';

const Login = () => {
  const [account, setAccount] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleLogin = async () => {
    try {
      const response = await axios.post('/api/login', { account, password });
      if (response.data.success) {
        message.success('登录成功');
        navigate('/console');
        localStorage.setItem('isLoggedIn', 'true');
        // 跳转到主界面逻辑
      } else {
        message.error('账号或密码错误');
      }
    } catch (error) {
      message.error('登录失败');
    }
  };

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
      <Card style={{ width: 300 }}>
        <Form
          name="login"
          initialValues={{ remember: true }}
          onFinish={handleLogin}
        >
          <Form.Item
            name="account"
            rules={[{ required: true, message: '请输入账号!' }]}
          >
            <Input
              placeholder="账号"
              value={account}
              onChange={(e) => setAccount(e.target.value)}
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码!' }]}
          >
            <Input.Password
              placeholder="密码"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit">
              登录
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default Login;