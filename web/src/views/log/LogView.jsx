import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { Button, Card, Col, Input, Row, Select, Table, Typography } from 'antd';
import { DownloadOutlined, SearchOutlined } from '@ant-design/icons';
import { downloadFile, getFileList, getLog } from '../console/api';
import { _constant } from '../../common/constant';

const { Title, Text } = Typography;
const { Option } = Select;

const LogView = () => {
  const location = useLocation();
  const [files, setFiles] = useState([]);
  const [selectedFile, setSelectedFile] = useState('');
  const [logContent, setLogContent] = useState([]);
  const [loading, setLoading] = useState(false);
  const [previewParams, setPreviewParams] = useState({
    serverName: '',
    serverId: 0,
    nodeId: 0,
    len: 100,
    keyword: '',
    host: '',
    logType: 1,
    fileName: ''
  });

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const host = params.get('host');
    const serverId = params.get('serverId');
    const serverName = params.get('serverName');
    const nodeId = params.get('nodeId');
    if (host && serverId) {
      setLoading(true);
      getFileList({ host, serverId : Number(serverId),type:_constant.FILE_TYPE_LOG }).then(res => {
        console.log('res',res);
        setPreviewParams({
          ...previewParams,
          host,
          serverName,
          serverId : Number(serverId),
          nodeId:Number(nodeId)
        })
        setFiles(res.data.map(file => ({ key: file, name: file })));
        setLoading(false);
      });
    }
  }, [location.search]);

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

  const handlePreview = () => {
    if (selectedFile) {
      setLoading(true);
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
        setLoading(false);
      });
    }
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
            <Row gutter={16} style={{ marginBottom: '16px' }}>
              <Col span={4}>
                <Text>行数</Text>
                <Input 
                  type="number" 
                  value={previewParams.len}
                  onChange={(e) => setPreviewParams({...previewParams, len: parseInt(e.target.value)})}
                />
              </Col>
              <Col span={4}>
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
              <Col span={8}>
                <Text>关键字</Text>
                <Input 
                  value={previewParams.keyword}
                  onChange={(e) => setPreviewParams({...previewParams, keyword: e.target.value})}
                />
              </Col>
              <Col span={4}>
                <Button 
                  type="primary" 
                  icon={<SearchOutlined />} 
                  onClick={handlePreview}
                  loading={loading}
                  style={{ marginTop: '24px' }}
                >
                  预览
                </Button>
              </Col>
            </Row>
              <pre style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word',background: 'black',color: 'white',padding: '10px',minHeight:'500px' }}>
                {logContent.map((v,i)=>{
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