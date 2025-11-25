import React, { useState, useEffect, useRef } from 'react';
import { useLocation } from 'react-router-dom';
import { Button, Card, Col, Input, Row, Select, Table, Typography, Radio, Checkbox } from 'antd';
import { DownloadOutlined, SearchOutlined } from '@ant-design/icons';
import { downloadFile, getFileList, getLog } from '../console/api';
import { _constant } from '../../common/constant';

const { Title, Text } = Typography;
const { Option } = Select;
const LOG_CATEGORIES = {
  BUSINESS: 1, // 业务日志
  MASTER: 2,   // 主控日志
  NODE: 3      // 节点日志
};

const LogView = () => {
  const location = useLocation();
  const [files, setFiles] = useState([]);
  const [selectedFile, setSelectedFile] = useState('');
  const [logContent, setLogContent] = useState([]);
  const [loading, setLoading] = useState(false);
  const [logOrder, setLogOrder] = useState(2); // 1: 从下到上(倒序), 2: 从上到下(正序)
  const [autoRefresh, setAutoRefresh] = useState(false); // 自动刷新开关
  const [refreshInterval, setRefreshInterval] = useState(5); // 刷新间隔（秒）
  const intervalRef = useRef(null); // 定时器引用
  const [previewParams, setPreviewParams] = useState({
    serverName: '',
    serverId: 0,
    nodeId: 0,
    len: 100,
    keyword: '',
    host: '',
    logType: 1,
    fileName: '',
    logCategory: LOG_CATEGORIES.BUSINESS // 新增日志分类参数
  });

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const host = params.get('host');
    const serverId = params.get('serverId');
    const serverName = params.get('serverName');
    const nodeId = params.get('nodeId');
    const logCategory = params.get('logCategory');
    if (host && serverId) {
      setLoading(true);
      getFileList({ 
        host, 
        serverId : Number(serverId),
        type:_constant.FILE_TYPE_LOG,
        logCategory: Number(logCategory)
      }).then(res => {
        console.log('res',res);
        setPreviewParams({
          ...previewParams,
          host,
          serverName:serverName || '',
          serverId : Number(serverId),
          nodeId:Number(nodeId),
          logCategory:Number(logCategory)
        })
        
        // 过滤文件：只保留包含 "log" 的文件，过滤掉 .json 和 .gz 文件
        const filteredFiles = res.data.filter(file => {
          const lowerFile = file.toLowerCase();
          return lowerFile.includes('log') && !lowerFile.endsWith('.json') && !lowerFile.endsWith('.gz');
        });
        
        // 添加排序逻辑
        const sortedFiles = filteredFiles.sort((a, b) => {
          // 提取文件名中的日期部分
          const dateRegex = /\.(\d{4}-\d{2}-\d{2})/;
          const dateA = a.match(dateRegex);
          const dateB = b.match(dateRegex);
          
          // 如果A没有日期，B有日期，A排在前面
          if (!dateA && dateB) return -1;
          // 如果A有日期，B没有日期，B排在前面
          if (dateA && !dateB) return 1;
          // 如果都没有日期，按文件名字母顺序
          if (!dateA && !dateB) return a.localeCompare(b);
          
          // 如果都有日期，按日期倒序（新的在前）
          const dateStrA = dateA[1];
          const dateStrB = dateB[1];
          return dateStrB.localeCompare(dateStrA);
        });
        
        setFiles(sortedFiles.map(file => ({ key: file, name: file })));
        setLoading(false);
      });
    }
  }, [location.search]);

  // 自动刷新效果
  useEffect(() => {
    if (autoRefresh && selectedFile) {
      // 启动定时器
      intervalRef.current = setInterval(() => {
        handlePreview(true); // 传入true表示是自动刷新
      }, refreshInterval * 1000);
    } else {
      // 清除定时器
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    }

    // 清理函数
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    };
  }, [autoRefresh, selectedFile, refreshInterval, previewParams]);

  const handleDownload = () => {
    if (selectedFile) {
      downloadFile({
        serverId:previewParams.serverId,
        fileName:selectedFile,
        type:_constant.FILE_TYPE_LOG,
        host:previewParams.host
      })
      // window.open(`/api/downloadFile?file=${selectedFile}`);
    }
  };

  const handlePreview = (isAutoRefresh = false) => {
    if (selectedFile) {
      // 如果是自动刷新，不显示loading状态
      if (!isAutoRefresh) {
        setLoading(true);
      }
      let keyword = previewParams.keyword
      if(!keyword){
        keyword = `''`
      }
      getLog({
        ...previewParams,
        fileName: selectedFile,
        keyword:keyword,
      }).then(res => {
        setLogContent(res.data || []);
        if (!isAutoRefresh) {
          setLoading(false);
        }
      }).catch(err => {
        console.error('获取日志失败:', err);
        if (!isAutoRefresh) {
          setLoading(false);
        }
      });
    }
  };

  // 处理日志顺序变化
  const handleLogOrderChange = (e) => {
    setLogOrder(e.target.value);
  };

  // 处理自动刷新开关变化
  const handleAutoRefreshChange = (e) => {
    setAutoRefresh(e.target.checked);
  };

  // 处理刷新间隔变化
  const handleRefreshIntervalChange = (e) => {
    const value = parseInt(e.target.value);
    if (value > 0) {
      setRefreshInterval(value);
    }
  };

  // 根据选择的顺序显示日志内容
  const getDisplayLogContent = () => {
    if (logOrder === 1) {
      // 从下到上（倒序）
      return [...logContent].reverse();
    }
    // 从上到下（正序）
    return logContent;
  };

  const columns = [
    {
      title: '文件名',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Button 
          icon={<DownloadOutlined />} 
          onClick={() => {
            setSelectedFile(record.name);
            handleDownload();
          }}
        />
      ),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Row gutter={16} style={{ marginBottom: '24px' }}>
        <Col span={4}>
          <Card title="日志文件列表" bordered={false} style={{maxHeight:"80vh", overflowY: "auto"}}>
            <Table 
              columns={columns}
              dataSource={files}
              loading={loading}
              rowSelection={{
                type: 'radio',
                selectedRowKeys: selectedFile ? [selectedFile] : [],
                onChange: (selectedRowKeys) => {
                  setSelectedFile(selectedRowKeys[0]);
                },
              }}
              pagination={false}
            />
          </Card>
        </Col>
        <Col span={20}>
          <Card title="日志预览" bordered={false}>
            <Row gutter={24} style={{ marginBottom: '16px' }}>
              <Col span={3}>
                <Text>行数</Text>
                <Input 
                  type="number" 
                  value={previewParams.len}
                  onChange={(e) => setPreviewParams({...previewParams, len: parseInt(e.target.value)})}
                />
              </Col>
              <Col span={3}>
                <Text>类型</Text>
                <Select 
                  style={{ width: '100%' }}
                  value={previewParams.logType}
                  onChange={(value) => setPreviewParams({...previewParams, logType: value})}
                >
                  <Option value={1}>head</Option>
                  <Option value={2}>tail</Option>
                </Select>
              </Col>
              <Col span={4}>
                <Text>关键字</Text>
                <Input 
                  value={previewParams.keyword}
                  onChange={(e) => setPreviewParams({...previewParams, keyword: e.target.value})}
                />
              </Col>
              <Col span={2}>
                <Button 
                  type="primary" 
                  icon={<SearchOutlined />} 
                  onClick={() => handlePreview()}
                  loading={loading}
                  style={{ marginTop: '24px' }}
                >
                  预览
                </Button>
              </Col>
              <Col span={4}>
                <Text>显示顺序</Text>
                <Radio.Group 
                  value={logOrder} 
                  onChange={handleLogOrderChange}
                  style={{ marginTop: '8px' }}
                >
                  <Radio value={1}>从下到上</Radio>
                  <Radio value={2}>从上到下</Radio>
                </Radio.Group>
              </Col>
              <Col span={6}>
                <Text>自动刷新</Text>
                <div style={{ marginTop: '8px' }}>
                  <Checkbox 
                    checked={autoRefresh}
                    onChange={handleAutoRefreshChange}
                  >
                    启用自动刷新（5秒间隔）
                  </Checkbox>
                </div>
              </Col>
            </Row>
              <pre style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word',background: 'black',color: 'white',padding: '10px',minHeight:'500px' }}>
                {getDisplayLogContent().map((v,i)=>{  
                  return <p key={i}>{i} {v}</p>
                })}
              </pre>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default LogView;