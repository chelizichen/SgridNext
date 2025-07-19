import React from "react";
import ModuleRegistry from "../../modules/registry";
import type { Module } from "../../modules/registry";

const styles = {
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
  statsGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
    gap: '20px',
    marginBottom: '30px'
  },
  statCard: {
    backgroundColor: 'white',
    padding: '20px',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0, 0, 0, 0.05)',
    textAlign: 'center' as const
  },
  statNumber: {
    fontSize: '32px',
    fontWeight: 'bold',
    color: '#667eea',
    margin: '0 0 8px 0'
  },
  statLabel: {
    fontSize: '14px',
    color: '#666',
    margin: '0'
  },
  modulesGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))',
    gap: '20px'
  },
  moduleCard: {
    backgroundColor: 'white',
    border: '1px solid #e1e5e9',
    borderRadius: '8px',
    padding: '20px',
    boxShadow: '0 2px 4px rgba(0, 0, 0, 0.05)',
    transition: 'all 0.2s ease'
  },
  moduleHeader: {
    display: 'flex',
    alignItems: 'center',
    marginBottom: '12px'
  },
  moduleIcon: {
    fontSize: '24px',
    marginRight: '12px'
  },
  moduleTitle: {
    fontSize: '18px',
    fontWeight: '600',
    color: '#2c3e50',
    margin: '0'
  },
  moduleDescription: {
    color: '#5a6c7d',
    fontSize: '14px',
    lineHeight: '1.5',
    marginBottom: '16px'
  },
  moduleMeta: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '16px'
  },
  moduleVersion: {
    fontSize: '12px',
    color: '#7f8c8d',
    backgroundColor: '#ecf0f1',
    padding: '4px 8px',
    borderRadius: '4px'
  },
  moduleStatus: {
    padding: '4px 12px',
    borderRadius: '20px',
    fontSize: '12px',
    fontWeight: '500'
  },
  statusEnabled: {
    backgroundColor: '#d4edda',
    color: '#155724'
  },
  statusDisabled: {
    backgroundColor: '#f8d7da',
    color: '#721c24'
  },
  moduleActions: {
    display: 'flex',
    gap: '8px'
  },
  actionButton: {
    padding: '8px 16px',
    borderRadius: '6px',
    border: 'none',
    fontSize: '13px',
    fontWeight: '500',
    cursor: 'pointer',
    transition: 'all 0.2s ease'
  },
  enableButton: {
    backgroundColor: '#28a745',
    color: 'white'
  },
  disableButton: {
    backgroundColor: '#dc3545',
    color: 'white'
  },
  viewButton: {
    backgroundColor: '#007bff',
    color: 'white'
  },
  moduleRoutes: {
    marginTop: '12px',
    padding: '12px',
    backgroundColor: '#f8f9fa',
    borderRadius: '4px'
  },
  routesTitle: {
    fontSize: '12px',
    fontWeight: '600',
    color: '#495057',
    margin: '0 0 8px 0'
  },
  routeList: {
    fontSize: '11px',
    color: '#6c757d',
    margin: '0',
    fontFamily: 'monospace'
  }
};

function ModuleCard({ module }: { module: Module }) {
  const handleToggle = async () => {
    try {
      const response = await fetch(`/api/admin/modules/${module.id}/toggle`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ enabled: !module.enabled })
      });
      
      if (response.ok) {
        window.location.reload();
      } else {
        console.error('切换模块状态失败');
      }
    } catch (error) {
      console.error('切换模块状态失败:', error);
    }
  };

  return (
    <div style={styles.moduleCard}>
      <div style={styles.moduleHeader}>
        <span style={styles.moduleIcon}>{module.icon || '📦'}</span>
        <h3 style={styles.moduleTitle}>{module.name}</h3>
      </div>
      
      <p style={styles.moduleDescription}>{module.description}</p>
      
      <div style={styles.moduleMeta}>
        <span style={styles.moduleVersion}>v{module.version}</span>
        <span style={{
          ...styles.moduleStatus,
          ...(module.enabled ? styles.statusEnabled : styles.statusDisabled)
        }}>
          {module.enabled ? '已启用' : '已禁用'}
        </span>
      </div>
      
      <div style={styles.moduleActions}>
        <button 
          style={{
            ...styles.actionButton,
            ...(module.enabled ? styles.disableButton : styles.enableButton)
          }}
          onClick={handleToggle}
        >
          {module.enabled ? '禁用' : '启用'}
        </button>
        
        {module.enabled && module.routes.length > 0 && (
          <a 
            href={module.routes[0]}
            style={{
              ...styles.actionButton,
              ...styles.viewButton,
              textDecoration: 'none',
              display: 'inline-block'
            }}
          >
            查看
          </a>
        )}
      </div>
      
      <div style={styles.moduleRoutes}>
        <h4 style={styles.routesTitle}>路由 ({module.routes.length})</h4>
        <ul style={styles.routeList}>
          {module.routes.map((route, index) => (
            <li key={index}>{route}</li>
          ))}
        </ul>
      </div>
    </div>
  );
}

function AdminView() {
  const moduleRegistry = ModuleRegistry.getInstance();
  const stats = moduleRegistry.getStats();
  const modules = moduleRegistry.getAllModules();

  return (
      <div style={styles.container}>
        <div style={styles.header}>
          <h1 style={styles.headerTitle}>🔧 系统管理</h1>
          <p style={styles.headerSubtitle}>
            管理和配置系统模块
          </p>
        </div>
        
        <div style={styles.statsGrid}>
          <div style={styles.statCard}>
            <div style={styles.statNumber}>{stats.total}</div>
            <div style={styles.statLabel}>总模块数</div>
          </div>
          <div style={styles.statCard}>
            <div style={styles.statNumber}>{stats.enabled}</div>
            <div style={styles.statLabel}>已启用</div>
          </div>
          <div style={styles.statCard}>
            <div style={styles.statNumber}>{stats.disabled}</div>
            <div style={styles.statLabel}>已禁用</div>
          </div>
        </div>
        
        <div style={styles.modulesGrid}>
          {modules.map((module) => (
            <ModuleCard key={module.id} module={module} />
          ))}
        </div>
      </div>
  );
}

export default AdminView; 