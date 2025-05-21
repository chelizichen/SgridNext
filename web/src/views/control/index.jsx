import NodeManager from './components/NodeManager';

export default function Control(){
    return (
        <div style={{ padding: 24 }}>
          <h2 style={{ marginBottom: 24 }}>集群节点管理</h2>
          <NodeManager />
        </div>
      );
}