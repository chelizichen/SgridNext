import React from 'react';
import { renderToString } from "react-dom/server";
import logger from "../components/logger/main";
// 导入所有页面组件
import Bugs from "../client/views/bugs/bugView";
import AdminView from "../client/views/admin/view";

// 导入所有服务
import bugService from "./bugs/bugService";

// 页面组件类型定义
interface PageComponent {
  component: React.ComponentType<any>;
  getData?: () => Promise<any> | any;
  title?: string;
}

// 路由配置
const routes: Record<string, PageComponent> = {
  '/view/bugs': {
    component: Bugs,
    getData: () => bugService.getBugs(),
    title: '问题追踪'
  },
  '/view/admin': {
    component: AdminView,
    getData: () => ({}),
    title: '系统管理'
  }
  // 可以添加更多路由
  // '/view/users': {
  //   component: Users,
  //   getData: () => userService.getUsers(),
  //   title: '用户管理'
  // },
  // '/view/projects': {
  //   component: Projects,
  //   getData: () => projectService.getProjects(),
  //   title: '项目管理'
  // }
};


const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SGRIDNODE_HYBRID</title>
</head>
<body>
    <div id="root"></div>
    <div><!--SGRID_BUN_HYDRATE--></div>
    <script type="module" src="/client.js"></script>
</body>
</html>
`

// 路由管理器
class RouterManager {
  private static instance: RouterManager;

  private constructor() {}

  static getInstance(): RouterManager {
    if (!RouterManager.instance) {
      RouterManager.instance = new RouterManager();
    }
    return RouterManager.instance;
  }

  // 获取路由配置
  getRoutes() {
    return routes;
  }

  // 检查路由是否存在
  hasRoute(path: string): boolean {
    return path in routes;
  }

  // 获取页面组件
  getPageComponent(path: string): PageComponent | null {
    return routes[path] || null;
  }

  // 渲染页面
  async renderPage(path: string): Promise<{ html: string; data?: any }> {
    const pageConfig = this.getPageComponent(path);

    if (!pageConfig) {
      throw new Error(`路由 ${path} 不存在`);
    }

    try {
      // 获取数据
      const data = pageConfig.getData ? await pageConfig.getData() : {};

      // 创建组件实例
      const Component = pageConfig.component;
      
      const componentInstance = React.createElement(Component, data);

      // 渲染为HTML
      const html = renderToString(componentInstance);
      const htmlFile = htmlTemplate
      // 注入数据
      const fullHtml = htmlFile.replace(
        `<div id="root"></div>`,
        `<div id="root">${html}</div>`
      ).replace(
        `<!--SGRID_BUN_HYDRATE-->`,
        `<!--SGRID_BUN_HYDRATE-->
        <script>
            window.__INITIAL_DATA__ = ${JSON.stringify(data)};
            window.__PAGE_TITLE__ = ${JSON.stringify(pageConfig.title || "")};
            window.__CURRENT_ROUTE__ = ${JSON.stringify(path)};
            window.__AVAILABLE_ROUTES__ = ${JSON.stringify(this.getAllPaths())};
        </script>`
      );

      return { html: fullHtml, data };
    } catch (error) {
      logger.data.error("渲染页面失败: %s", error);
      throw error;
    }
  }

  // 获取所有路由路径
  getAllPaths(): string[] {
    return Object.keys(routes);
  }
}

export default RouterManager; 