import React, { useEffect, useRef, useState, useCallback } from 'react';
import { Card, Button, message, Select } from 'antd';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import '@xterm/xterm/css/xterm.css';
import './WebShell.css';
import { getNodeList } from '../console/api';

const { Option } = Select;

const WebShell = () => {
  const terminalRef = useRef(null);
  const wsRef = useRef(null);
  const terminal = useRef(null);
  const fitAddon = useRef(null);
  const [connected, setConnected] = useState(false);
  const [selectedNode, setSelectedNode] = useState(null); // null 表示主控节点
  const [nodes, setNodes] = useState([]);

  const connectWebSocket = useCallback((term, nodeHost = null) => {
    // 获取 WebSocket URL
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    let wsUrl;
    
    if (nodeHost) {
      // 连接到节点服务器，使用固定端口 25529
      wsUrl = `${protocol}//${nodeHost}:25529/webshell/ws`;
    } else {
      // 连接到主控服务器
      const host = window.location.host;
      wsUrl = `${protocol}//${host}/api/webshell/ws`;
    }

    try {
      const ws = new WebSocket(wsUrl);
      // 设置二进制类型为 arraybuffer，这样接收到的二进制消息会是 ArrayBuffer
      ws.binaryType = 'arraybuffer';

      ws.onopen = () => {
        setConnected(true);
        message.success('WebShell 连接成功');
        // 调整终端大小
        if (fitAddon.current) {
          fitAddon.current.fit();
        }
        // 确保终端可见
        if (term) {
          term.focus();
        }
      };

      ws.onmessage = async (event) => {
        try {
          let data;
          if (event.data instanceof ArrayBuffer) {
            // ArrayBuffer 类型
            data = new Uint8Array(event.data);
            console.log('收到 ArrayBuffer 消息，长度:', data.length);
            term.write(data);
          } else if (event.data instanceof Blob) {
            // Blob 类型（WebSocket 二进制消息可能是 Blob）
            const arrayBuffer = await event.data.arrayBuffer();
            data = new Uint8Array(arrayBuffer);
            console.log('收到 Blob 消息，长度:', data.length);
            term.write(data);
          } else if (typeof event.data === 'string') {
            // 文本数据
            console.log('收到文本消息:', event.data);
            term.write(event.data);
          } else {
            // 其他类型，尝试转换为字符串
            console.warn('未知的消息类型:', typeof event.data, event.data);
            term.write(String(event.data));
          }
        } catch (error) {
          console.error('处理 WebSocket 消息失败:', error, event);
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket 错误:', error);
        message.error('WebSocket 连接错误');
        setConnected(false);
      };

      ws.onclose = () => {
        setConnected(false);
        message.warning('WebShell 连接已断开');
        // 不自动重连，让用户手动选择节点
      };

      wsRef.current = ws;

      // 处理终端输入
      term.onData((data) => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(data);
        }
      });
    } catch (error) {
      console.error('创建 WebSocket 连接失败:', error);
      message.error('创建 WebSocket 连接失败');
    }
  }, []);

  useEffect(() => {
    // 初始化终端
    const term = new Terminal({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'Consolas, "Courier New", monospace',
      theme: {
        background: '#1e1e1e',
        foreground: '#d4d4d4',
        cursor: '#aeafad',
        selection: '#264f78',
      },
      convertEol: true, // 自动转换换行符
      disableStdin: false, // 确保输入启用
    });

    const fit = new FitAddon();
    term.loadAddon(fit);
    fitAddon.current = fit;

    if (terminalRef.current) {
      term.open(terminalRef.current);
      fit.fit();
      terminal.current = term;
      // 确保终端获得焦点
      term.focus();
    }

    // 加载节点列表
    getNodeList().then((res) => {
      if (res.success && res.data) {
        setNodes(res.data);
      }
    }).catch((err) => {
      console.error('获取节点列表失败:', err);
    });

    // 初始连接（主控节点）- 延迟连接以确保终端已初始化
    setTimeout(() => {
      connectWebSocket(term);
    }, 100);

    // 窗口大小改变时调整终端大小
    const handleResize = () => {
      if (fitAddon.current) {
        fitAddon.current.fit();
      }
    };

    window.addEventListener('resize', handleResize);

    // 清理函数
    return () => {
      window.removeEventListener('resize', handleResize);
      if (wsRef.current) {
        wsRef.current.close();
      }
      if (terminal.current) {
        terminal.current.dispose();
      }
    };
  }, [connectWebSocket]);

  const handleReconnect = () => {
    if (wsRef.current) {
      wsRef.current.close();
    }
    if (terminal.current) {
      const nodeHost = selectedNode ? selectedNode.Host : null;
      connectWebSocket(terminal.current, nodeHost);
    }
  };

  const handleNodeChange = (value) => {
    // 关闭当前连接
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    setConnected(false);
    
    // 清空终端
    if (terminal.current) {
      terminal.current.clear();
    }
    
    const selectedNodeObj = value === 'master' ? null : nodes.find(n => n.ID === value);
    setSelectedNode(selectedNodeObj);
    
    // 重新连接
    if (terminal.current) {
      const nodeHost = selectedNodeObj ? selectedNodeObj.Host : null;
      connectWebSocket(terminal.current, nodeHost);
    }
  };

  const handleClear = () => {
    if (terminal.current) {
      terminal.current.clear();
    }
  };

  return (
    <div style={{ height: 'calc(100vh - 64px)' }}>
      <Card
        title="WebShell 终端"
        extra={
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <Select
              style={{ width: 200 }}
              placeholder="选择节点"
              value={selectedNode ? selectedNode.ID : 'master'}
              onChange={handleNodeChange}
            >
              <Option value="master">主控节点</Option>
              {nodes.map((node) => (
                <Option key={node.ID} value={node.ID}>
                  {node.Alias || node.Host} ({node.NodeStatus === 1 ? '在线' : '离线'})
                </Option>
              ))}
            </Select>
            <Button
              type="primary"
              onClick={handleReconnect}
            >
              {connected ? '重新连接' : '连接'}
            </Button>
            <Button onClick={handleClear}>清屏</Button>
            <span style={{ marginLeft: 8, color: connected ? '#52c41a' : '#ff4d4f' }}>
              {connected ? '● 已连接' : '● 未连接'}
            </span>
          </div>
        }
        style={{ height: '100%', display: 'flex', flexDirection: 'column' }}
        bodyStyle={{ flex: 1, padding: 0, display: 'flex', flexDirection: 'column' }}
      >
        <div
          ref={terminalRef}
          style={{
            flex: 1,
            padding: '16px',
            backgroundColor: '#1e1e1e',
            overflow: 'hidden',
            minHeight: '400px',
            width: '100%',
          }}
        />
      </Card>
    </div>
  );
};

export default WebShell;

