# BUN STUDY

## 文档

https://bun.net.cn/docs/bundler/executables

## CLI

1. 安装依赖 bun i knex
2. 普通打包 bun build ./index.ts --outdir ./dist
3. 开发模式 bun run --watch ./index.ts
4. 本机编译 bun build ./index.ts --compile --outfile app
5. 跨平台编译 bun build --compile --minify --sourcemap --target=bun-linux-x64 ./index.ts --outfile app

## 🏗️ 整体架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   前端界面开发   │    │  服务端渲染     │    │   后端接口开发   │
│   (React组件)   │    │  (SSR + 水合)   │    │   (API路由)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   统一路由管理   │
                    │  (页面 + API)   │
                    └─────────────────┘
```

## 🔄 开发流程

### 1. 前端界面开发

#### 1.1 创建React组件
```typescript
// src/client/views/your-module/view.tsx
import React from "react";

interface YourModuleProps {
  data: any;
}

function YourModuleView({ data }: YourModuleProps) {
  return (
    <div style={{ padding: '20px' }}>
      <h1>你的模块</h1>
      {/* 你的组件内容 */}
    </div>
  );
}

export default YourModuleView;
```

#### 1.2 组件设计原则
- **服务端渲染友好**: 避免使用 `window`、`document` 等浏览器API
- **数据驱动**: 通过 props 接收数据，避免直接调用API
- **样式内联**: 使用内联样式或CSS-in-JS，避免外部CSS依赖
- **错误边界**: 添加错误处理，提供友好的错误提示

### 2. 服务端渲染配置

#### 2.1 注册页面路由
```typescript
// src/client/views/router/index.tsx
import YourModuleView from "../your-module/view";

const routes: Record<string, PageComponent> = {
  '/view/your-module': {
    component: YourModuleView,
    getData: () => yourService.getData(), // 获取初始数据
    title: '你的模块'
  }
};
```

#### 2.2 数据获取函数
```typescript
// 在对应的服务文件中
const yourService = {
  async getData() {
    // 从数据库、API或其他数据源获取数据
    return { items: [], total: 0 };
  }
};
```

#### 2.3 水合配置
```typescript
// src/client.tsx - 客户端入口
const pageComponentMap = {
  '/view/your-module': () => import('./client/views/your-module/view'),
};
```

### 3. 后端接口开发

#### 3.1 创建业务服务
```typescript
// src/domain/your-module/yourService.ts
const yourService = {
  async getItems() {
    // 业务逻辑
    return { items: [], total: 0 };
  },
  
  async createItem(data: any) {
    // 创建逻辑
    return { success: true, id: 1 };
  },
  
  async updateItem(id: number, data: any) {
    // 更新逻辑
    return { success: true };
  }
};

export default yourService;
```

#### 3.2 注册API路由
```typescript
// src/domain/router.ts
import yourService from "./your-module/yourService";

const apiRoutes = {
  '/api/your-module': {
    GET: async (req: Request) => {
      const data = await yourService.getItems();
      return new Response(JSON.stringify(data), {
        headers: { "Content-Type": "application/json" }
      });
    },
    POST: async (req: Request) => {
      const body = await req.json();
      const result = await yourService.createItem(body);
      return new Response(JSON.stringify(result), {
        headers: { "Content-Type": "application/json" }
      });
    }
  },
  '/api/your-module/:id': {
    PUT: async (req: Request) => {
      const id = req.params.id;
      const body = await req.json();
      const result = await yourService.updateItem(parseInt(id), body);
      return new Response(JSON.stringify(result), {
        headers: { "Content-Type": "application/json" }
      });
    }
  }
};
```

#### 3.3 注册模块
```typescript
// src/client/modules/registry.ts
this.registerModule({
  id: 'your-module',
  name: '你的模块',
  description: '模块描述',
  version: '1.0.0',
  routes: ['/view/your-module'],
  apis: ['/api/your-module', '/api/your-module/:id'],
  icon: '📦',
  enabled: true
});
```

## 🎯 开发步骤总结

### 步骤1: 规划模块
1. 确定模块功能
2. 设计数据模型
3. 规划页面路由和API路由

### 步骤2: 开发后端
1. 创建业务服务 (`src/domain/your-module/yourService.ts`)
2. 实现数据操作逻辑
3. 注册API路由 (`src/domain/router.ts`)

### 步骤3: 开发前端
1. 创建React组件 (`src/client/views/your-module/view.tsx`)
2. 设计用户界面
3. 实现交互逻辑

### 步骤4: 配置路由
1. 注册页面路由 (`src/client/views/router/index.tsx`)
2. 配置数据获取函数
3. 添加客户端水合配置 (`src/client.tsx`)

### 步骤5: 注册模块
1. 在模块注册表中添加模块 (`src/client/modules/registry.ts`)
2. 配置模块信息（名称、描述、路由等）

### 步骤6: 测试验证
1. 启动服务器 (`bun run index.ts`)
2. 访问页面路由 (`/view/your-module`)
3. 测试API接口 (`/api/your-module`)
4. 验证水合功能

## 🔧 关键技术点

### 1. 服务端渲染 (SSR)
```typescript
// 服务端渲染流程
const html = renderToString(React.createElement(Component, data));
const fullHtml = htmlTemplate.replace(
  `<div id="root"></div>`,
  `<div id="root">${html}</div>`
);
```

### 2. 客户端水合
```typescript
// 客户端水合流程
const container = document.getElementById('root');
const serverData = window.__INITIAL_DATA__;
hydrateRoot(container, React.createElement(PageComponent, serverData));
```

### 3. 数据注入
```typescript
// 服务端数据注入
const script = `
  window.__INITIAL_DATA__ = ${JSON.stringify(data)};
  window.__CURRENT_ROUTE__ = ${JSON.stringify(path)};
`;
```

### 4. 动态路由
```typescript
// 参数路由支持
'/api/items/:id': {
  GET: async (req: Request) => {
    const id = req.params.id; // 自动解析参数
    // 处理逻辑
  }
}
```

## 📝 开发规范

### 1. 文件命名
- 组件文件: `view.tsx` 或 `component.tsx`
- 服务文件: `service.ts`
- 类型文件: `types.ts`

### 2. 目录结构
```
src/
├── client/views/your-module/     # 前端组件
├── domain/your-module/           # 后端服务
└── components/                   # 通用组件
```

### 3. 代码规范
- 使用 TypeScript 类型定义
- 添加错误处理
- 记录关键日志
- 遵循模块化原则

### 4. 测试要点
- 服务端渲染是否正常
- 客户端水合是否成功
- API接口是否响应
- 数据流是否正确

## 🚀 快速开始

1. **克隆项目**
```bash
git clone <your-repo>
cd bunserver
```

2. **安装依赖**
```bash
bun install
```

3. **启动开发服务器**
```bash
bun run index.ts
```

4. **访问应用**
- 首页: http://localhost:3000/
- 管理界面: http://localhost:3000/view/admin
