import React, { useEffect, useState } from "react";
function BugItem(props: { bug: any }) {
  const [isLoading, setIsLoading] = useState(false);
  const getStatusStyle = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'fixed':
        return { ...styles.statusBadge, ...styles.statusFixed };
      case 'closed':
        return { ...styles.statusBadge, ...styles.statusClosed };
      default:
        return { ...styles.statusBadge, ...styles.statusOpen };
    }
  };

  const getPriorityStyle = (priority: string) => {
    switch (priority?.toLowerCase()) {
      case 'high':
        return { ...styles.priorityBadge, ...styles.priorityHigh };
      case 'medium':
        return { ...styles.priorityBadge, ...styles.priorityMedium };
      default:
        return { ...styles.priorityBadge, ...styles.priorityLow };
    }
  };

  const handleFix = async () => {
    setIsLoading(true);
    try {
      const response = await fetch(`/api/bugs/${props.bug.id}/fix`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        }
      });
      const data = await response.json();
      if (data.success) {
        console.log('修复成功:', data);
        // 可以更新UI或刷新页面
        window.location.reload();
      } else {
        console.error('修复失败:', data.message);
      }
    } catch (error) {
      console.error('修复失败:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = async () => {
    setIsLoading(true);
    try {
      const response = await fetch(`/api/bugs/${props.bug.id}/close`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        }
      });
      const data = await response.json();
      if (data.success) {
        console.log('关闭成功:', data);
        // 可以更新UI或刷新页面
        window.location.reload();
      } else {
        console.error('关闭失败:', data.message);
      }
    } catch (error) {
      console.error('关闭失败:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div style={styles.bugCard} data-bug-id={props.bug.id}>
      <div style={styles.bugHeader}>
        <h3 style={styles.bugTitle}>{props.bug.title || '未命名问题'}</h3>
        <span style={styles.bugId}>#{props.bug.id || 'N/A'}</span>
      </div>
      
      <div style={styles.bugMeta}>
        <span style={getStatusStyle(props.bug.status)}>
          {props.bug.status || 'Open'}
        </span>
        <span style={getPriorityStyle(props.bug.priority)}>
          {props.bug.priority || 'Low'}
        </span>
        {props.bug.createdAt && (
          <span style={{ fontSize: '12px', color: '#7f8c8d' }}>
            创建于: {new Date(props.bug.createdAt).toLocaleDateString('zh-CN')}
          </span>
        )}
      </div>
      
      {props.bug.description && (
        <div style={styles.bugDescription}>
          {props.bug.description}
        </div>
      )}
      
      <div style={styles.bugActions}>
        <button 
          style={{...styles.actionButton, ...styles.fixButton}}
          onClick={handleFix}
          disabled={isLoading}
        >
          {isLoading ? '⏳ 处理中...' : '🔧 修复'}
        </button>
        <button 
          style={{...styles.actionButton, ...styles.closeButton}}
          onClick={handleClose}
          disabled={isLoading}
        >
          {isLoading ? '⏳ 处理中...' : '❌ 关闭'}
        </button>
      </div>
    </div>
  );
}

function Bugs(props: { bugs: any[] }) {
  const [title,setTitle] = useState("");
  useEffect(()=>{
    setTitle("SGRIDNODE_HYBRID");
  },[]);
  return (
      <div style={styles.container}>
        <div style={styles.header}>
          <h1 style={styles.headerTitle}>🐛 {title} 问题追踪</h1>
          <p style={styles.headerSubtitle}>
            管理和跟踪项目中的问题与缺陷
          </p>
        </div>
        <div style={styles.bugsContainer}>
          {props.bugs && props.bugs.length > 0 ? (
            props.bugs.map((bug) => (
              <BugItem key={bug.id} bug={bug} />
            ))
          ) : (
            <div style={styles.emptyState}>
              <div style={styles.emptyIcon}>📋</div>
              <h3>暂无问题</h3>
              <p>当前没有待处理的问题</p>
            </div>
          )}
        </div>
      </div>
  );
}

export default Bugs;


const styles: Record<string, React.CSSProperties> = {
  container: {
    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
    backgroundColor: '#fafafa',
    minHeight: '100vh',
    padding: '20px'
  },
  header: {
    background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    color: 'white',
    padding: '20px 30px',
    borderRadius: '8px',
    marginBottom: '30px',
    boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)'
  },
  headerTitle: {
    fontSize: '28px',
    fontWeight: '600',
    margin: '0 0 8px 0'
  },
  headerSubtitle: {
    fontSize: '14px',
    opacity: '0.9',
    margin: '0'
  },
  bugsContainer: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '16px'
  },
  bugCard: {
    backgroundColor: 'white',
    border: '1px solid #e1e5e9',
    borderRadius: '8px',
    padding: '20px',
    boxShadow: '0 2px 4px rgba(0, 0, 0, 0.05)',
    transition: 'all 0.2s ease',
    // cursor: 'pointer'
  },
  bugCardHover: {
    boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
    transform: 'translateY(-2px)'
  },
  bugHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: '12px'
  },
  bugTitle: {
    fontSize: '16px',
    fontWeight: '600',
    color: '#2c3e50',
    margin: '0',
    flex: '1'
  },
  bugId: {
    fontSize: '12px',
    color: '#7f8c8d',
    backgroundColor: '#ecf0f1',
    padding: '4px 8px',
    borderRadius: '4px',
    marginLeft: '12px'
  },
  bugMeta: {
    display: 'flex',
    alignItems: 'center',
    gap: '16px',
    marginBottom: '16px'
  },
  statusBadge: {
    padding: '4px 12px',
    borderRadius: '20px',
    fontSize: '12px',
    fontWeight: '500',
    textTransform: 'uppercase' as const
  },
  statusOpen: {
    backgroundColor: '#fff3cd',
    color: '#856404',
    border: '1px solid #ffeaa7'
  },
  statusFixed: {
    backgroundColor: '#d4edda',
    color: '#155724',
    border: '1px solid #c3e6cb'
  },
  statusClosed: {
    backgroundColor: '#f8d7da',
    color: '#721c24',
    border: '1px solid #f5c6cb'
  },
  priorityBadge: {
    padding: '4px 8px',
    borderRadius: '4px',
    fontSize: '11px',
    fontWeight: '600'
  },
  priorityHigh: {
    backgroundColor: '#ff6b6b',
    color: 'white'
  },
  priorityMedium: {
    backgroundColor: '#feca57',
    color: '#2c3e50'
  },
  priorityLow: {
    backgroundColor: '#48dbfb',
    color: 'white'
  },
  bugDescription: {
    color: '#5a6c7d',
    fontSize: '14px',
    lineHeight: '1.5',
    marginBottom: '16px'
  },
  bugActions: {
    display: 'flex',
    gap: '8px',
    justifyContent: 'flex-end'
  },
  actionButton: {
    padding: '8px 16px',
    borderRadius: '6px',
    border: 'none',
    fontSize: '13px',
    fontWeight: '500',
    // cursor: 'pointer',
    transition: 'all 0.2s ease'
  },
  fixButton: {
    backgroundColor: '#28a745',
    color: 'white'
  },
  fixButtonHover: {
    backgroundColor: '#218838'
  },
  closeButton: {
    backgroundColor: '#6c757d',
    color: 'white'
  },
  closeButtonHover: {
    backgroundColor: '#5a6268'
  },
  emptyState: {
    textAlign: 'center' as const,
    padding: '60px 20px',
    color: '#7f8c8d'
  },
  emptyIcon: {
    fontSize: '48px',
    marginBottom: '16px',
    opacity: '0.5'
  }
};
