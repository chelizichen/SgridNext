import React from 'react';
import { hydrateRoot } from 'react-dom/client';

// 声明全局变量类型
declare global {
  interface Window {
    __INITIAL_DATA__: any;
    __PAGE_TITLE__: string;
    __CURRENT_ROUTE__: string;
    __AVAILABLE_ROUTES__: string[];
  }
}

// 页面组件映射（支持动态导入）
const pageComponentMap: Record<string, () => Promise<{ default: React.ComponentType<any> }>> = {
  '/view/bugs': () => import('./client/views/bugs/bugView'),
  '/view/admin': () => import('./client/views/admin/view'),
  // 可以添加更多页面组件
  // '/view/users': () => import('./domain/view/users/view'),
  // '/view/projects': () => import('./domain/view/projects/view'),
};

// 预加载的组件（用于快速访问）
const preloadedComponents: Record<string, React.ComponentType<any>> = {};

// 获取当前路由
const getCurrentRoute = (): string => {
  // 优先使用服务端注入的路由信息
  if (window.__CURRENT_ROUTE__) {
    return window.__CURRENT_ROUTE__;
  }
  
  // 回退到从URL获取
  return window.location.pathname;
};

// 动态加载页面组件
const loadPageComponent = async (route: string): Promise<React.ComponentType<any> | null> => {
  // 检查是否已预加载
  if (preloadedComponents[route]) {
    console.log(`✅ 使用预加载的组件: ${route}`);
    return preloadedComponents[route];
  }

  // 检查是否有动态导入配置
  const importFn = pageComponentMap[route];
  if (importFn) {
    try {
      console.log(`📦 动态加载组件: ${route}`);
      console.log('module',importFn);
      const module = await importFn();
      const component = module.default;
      
      // 缓存组件
      preloadedComponents[route] = component;
      
      console.log(`✅ 组件加载成功: ${route}`);
      return component;
    } catch (error) {
      console.error(`❌ 组件加载失败: ${route}`, error);
      return null;
    }
  }

  console.warn(`❌ 未找到组件配置: ${route}`);
  return null;
};

// 错误页面组件
const ErrorPage = ({ route, availableRoutes }: { route: string; availableRoutes: string[] }) => (
  <div style={{ 
    padding: '40px 20px', 
    textAlign: 'center',
    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    backgroundColor: '#fafafa',
    minHeight: '100vh',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center'
  }}>
    <div style={{
      background: 'white',
      padding: '40px',
      borderRadius: '12px',
      boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
      maxWidth: '500px',
      width: '100%'
    }}>
      <h1 style={{ color: '#e74c3c', marginBottom: '20px' }}>⚠️ 页面未找到</h1>
      <p style={{ color: '#5a6c7d', marginBottom: '16px' }}>
        路由 <code style={{ background: '#f8f9fa', padding: '2px 6px', borderRadius: '4px' }}>{route}</code> 对应的组件未找到
      </p>
      <p style={{ color: '#7f8c8d', fontSize: '14px', marginBottom: '24px' }}>
        请检查路由配置或联系管理员
      </p>
      
      {availableRoutes.length > 0 && (
        <div style={{ marginBottom: '24px' }}>
          <h3 style={{ fontSize: '16px', color: '#2c3e50', marginBottom: '12px' }}>可用的路由：</h3>
          <ul style={{ 
            listStyle: 'none', 
            padding: 0, 
            margin: 0,
            textAlign: 'left'
          }}>
            {availableRoutes.map((availableRoute) => (
              <li key={availableRoute} style={{ marginBottom: '8px' }}>
                <a 
                  href={availableRoute}
                  style={{ 
                    color: '#667eea', 
                    textDecoration: 'none',
                    padding: '8px 12px',
                    display: 'block',
                    borderRadius: '6px',
                    transition: 'background-color 0.2s'
                  }}
                  onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#f8f9fa'}
                  onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                >
                  {availableRoute}
                </a>
              </li>
            ))}
          </ul>
        </div>
      )}
      
      <a 
        href="/" 
        style={{ 
          color: 'white',
          backgroundColor: '#667eea',
          textDecoration: 'none',
          padding: '12px 24px',
          borderRadius: '6px',
          display: 'inline-block',
          fontWeight: '500',
          transition: 'background-color 0.2s'
        }}
        onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#5a6fd8'}
        onMouseLeave={(e) => e.currentTarget.style.backgroundColor = '#667eea'}
      >
        返回首页
      </a>
    </div>
  </div>
);

// 加载页面组件
const loadPageComponentSync = (route: string): React.ComponentType<any> | null => {
  // 检查是否已预加载
  if (preloadedComponents[route]) {
    console.log(`✅ 使用预加载的组件: ${route}`);
    return preloadedComponents[route];
  }

  console.warn(`❌ 未找到页面组件: ${route}`);
  return null;
};

// 主水合函数
const hydratePage = async () => {
  const container = document.getElementById('root');
  if (!container) {
    console.error('❌ 找不到根元素 #root');
    return;
  }

  console.log('🚀 开始React水合...');
  
  const currentRoute = getCurrentRoute();
  console.log('📍 当前路由:', currentRoute);
  
  // 获取服务端渲染的数据
  const serverData = window.__INITIAL_DATA__ || {};
  console.log('📦 服务端数据:', serverData);
  
  // 获取可用的路由列表
  const availableRoutes = window.__AVAILABLE_ROUTES__ || Object.keys(pageComponentMap);
  
  // 尝试同步加载页面组件
  let PageComponent = loadPageComponentSync(currentRoute);
  
  // 如果同步加载失败，尝试异步加载
  if (!PageComponent) {
    console.log('🔄 尝试异步加载组件...');
    PageComponent = await loadPageComponent(currentRoute);
  }
  
  if (PageComponent) {
    try {
      console.log('🔄 开始水合页面组件...');
      hydrateRoot(container, React.createElement(PageComponent, serverData));
      console.log('✅ React水合完成');
      
      // 设置页面标题
      if (window.__PAGE_TITLE__) {
        document.title = window.__PAGE_TITLE__;
      }
      
      // 预加载其他路由的组件（可选）
      // setTimeout(() => {
      //   console.log('🔄 开始预加载其他组件...');
      //   availableRoutes.forEach(route => {
      //     if (route !== currentRoute && !preloadedComponents[route]) {
      //       loadPageComponent(route).then(component => {
      //         if (component) {
      //           console.log(`✅ 预加载完成: ${route}`);
      //         }
      //       });
      //     }
      //   });
      // }, 1000);
      
    } catch (error) {
      console.error('❌ 水合过程中发生错误:', error);
      hydrateRoot(container, React.createElement(ErrorPage, { 
        route: currentRoute, 
        availableRoutes 
      }));
    }
  } else {
    console.warn('⚠️ 显示错误页面');
    hydrateRoot(container, React.createElement(ErrorPage, { 
      route: currentRoute, 
      availableRoutes 
    }));
  }
};

// 启动水合
hydratePage(); 