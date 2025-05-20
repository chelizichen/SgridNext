import React from 'react';
import { Modal, List } from 'antd';

export default function HistoryModal({ visible, onCancel, historyData, checkFileCb }) {
  return (
    <Modal
      title="历史版本"
      open={visible}
      onCancel={onCancel}
      footer={null}
      width={600}
    >
      <List
        dataSource={historyData}
        renderItem={item => (
          <List.Item>
            <div style={{ display: 'flex', justifyContent: 'space-between', width: '100%' }}>
              <span>{item}</span>
              <a onClick={() => checkFileCb(item)}>查看</a>
            </div>
          </List.Item>
        )}
      />
    </Modal>
  );
}