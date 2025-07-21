import logger from "../components/logger/main";
import RouterManager from "./clientManager";
import ApiRouterManager from "./apiManager";

// 获取路由管理器实例
const routerManager = RouterManager.getInstance();
const apiRouterManager = ApiRouterManager.getInstance();

// 主路由处理函数
async function handleRequest(req: Request): Promise<Response> {
  const url = new URL(req.url);
  const path = url.pathname;

  try {
    // 处理API请求
    if (path.startsWith('/api/')) {
      return await apiRouterManager.handleApiRequest(req);
    }

    // 处理页面请求
    if (path.startsWith('/view/')) {
      if (routerManager.hasRoute(path)) {
        const { html } = await routerManager.renderPage(path);
        return new Response(html, {
          headers: { "Content-Type": "text/html" }
        });
      }
    }

    // 处理客户端脚本请求
    if (path === '/client.js') {
      return await handleClientScript(req);
    }

    // 处理根路径
    if (path === '/' || path === '') {
      return new Response(`
        <html>
          <head>
            <title>SgridNext 管理系统</title>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <style>
              body {
                font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
                background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
                margin: 0;
                padding: 0;
                min-height: 100vh;
                display: flex;
                align-items: center;
                justify-content: center;
              }
              .container {
                background: white;
                padding: 40px;
                border-radius: 12px;
                box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
                text-align: center;
                max-width: 500px;
                width: 90%;
              }
              h1 {
                color: #2c3e50;
                margin-bottom: 20px;
              }
              .module-list {
                list-style: none;
                padding: 0;
                margin: 30px 0;
              }
              .module-list li {
                margin: 15px 0;
              }
              .module-link {
                display: inline-block;
                padding: 12px 24px;
                background: #667eea;
                color: white;
                text-decoration: none;
                border-radius: 6px;
                transition: all 0.3s ease;
                font-weight: 500;
              }
              .module-link:hover {
                background: #5a6fd8;
                transform: translateY(-2px);
                box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
              }
              .admin-link {
                background: #e74c3c;
              }
              .admin-link:hover {
                background: #c0392b;
                box-shadow: 0 4px 12px rgba(231, 76, 60, 0.4);
              }
            </style>
          </head>
          <body>
            <div class="container">
              <h1>🚀 SgridNext 管理系统</h1>
              <p>欢迎使用模块化管理系统</p>
              <ul class="module-list">
                <li><a href="/view/bugs" class="module-link">🐛 问题追踪</a></li>
                <li><a href="/view/admin" class="module-link admin-link">🔧 系统管理</a></li>
              </ul>
            </div>
          </body>
        </html>
      `, {
        headers: { "Content-Type": "text/html" }
      });
    }

    // 404 处理
    return new Response(`
      <html>
        <head>
          <title>页面未找到</title>
          <meta charset="UTF-8">
        </head>
        <body>
          <h1>404 - 页面未找到</h1>
          <p>请求的路径 ${path} 不存在</p>
          <a href="/">返回首页</a>
        </body>
      </html>
    `, {
      status: 404,
      headers: { "Content-Type": "text/html" }
    });

  } catch (error: any) {
    logger.data.error('请求处理失败: %s', error);
    return new Response(`
      <html>
        <head>
          <title>服务器错误</title>
          <meta charset="UTF-8">
        </head>
        <body>
          <h1>500 - 服务器内部错误</h1>
          <p>${error?.message || '未知错误'}</p>
          <a href="/">返回首页</a>
        </body>
      </html>
    `, {
      status: 500,
      headers: { "Content-Type": "text/html" }
    });
  }
}

const rootDir = process.cwd();
// 处理客户端脚本
async function handleClientScript(req:Request): Promise<Response> {
  try {
    // 使用Bun构建客户端脚本
    const clientPath = `${rootDir}/src/client.tsx`;
    logger.data.info('clientPath %s',clientPath);
    const clientBundle = await Bun.build({
      entrypoints: [clientPath],
      target: 'browser',
      minify: false,
      sourcemap: 'inline',
      define: {
        'process.env.NODE_ENV': '"development"',
        'process.env': '{}',
        'process': '{}'
      }
    });
    
    const output = await clientBundle.outputs[0]?.text() || '';
    logger.data.info('处理水合完成');
    return new Response(output, {
      headers: {
        "Content-Type": "application/javascript"
      }
    });
  } catch (error) {
    console.error('构建客户端脚本失败:', error);
    // 回退到简单的客户端脚本
    const fallbackScript = `
      console.log('客户端脚本加载完成');
      document.addEventListener('DOMContentLoaded', function() {
        console.log('DOM加载完成');
      });
    `;
    
    return new Response(fallbackScript, {
      headers: {
        "Content-Type": "application/javascript"
      }
    });
  }
}

// 导出路由配置
const routes = {
  // 使用通配符处理所有请求
  "/*": handleRequest
};

export default routes;