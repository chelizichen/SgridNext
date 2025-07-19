import logger from "../components/logger/main";
import bugService from "./bugs/bugService";
import ModuleRegistry from "../client/modules/registry";

// API路由处理器类型
interface ApiHandler {
  GET?: (req: Request) => Promise<Response> | Response;
  POST?: (req: Request) => Promise<Response> | Response;
  PUT?: (req: Request) => Promise<Response> | Response;
  DELETE?: (req: Request) => Promise<Response> | Response;
}

// API路由配置
const apiRoutes: Record<string, ApiHandler> = {
  // Bug相关API
  '/api/bugs': {
    GET: async (req: Request) => {
      try {
        const bugs = bugService.getBugs();
        return new Response(JSON.stringify(bugs), {
          headers: { "Content-Type": "application/json" }
        });
      } catch (error: any) {
        return new Response(JSON.stringify({
          success: false,
          message: "获取Bug列表失败",
          error: error?.message || "未知错误"
        }), {
          status: 500,
          headers: { "Content-Type": "application/json" }
        });
      }
    }
  },
  
  '/api/bugs/:id/fix': {
    POST: async (req: Request) => {
      const url = new URL(req.url);
      const id = url.pathname.split('/')[3] || '';
      
      try {
        const result = bugService.fix(parseInt(id));
        return new Response(JSON.stringify({
          success: true,
          message: "Bug修复成功",
          data: result
        }), {
          headers: { "Content-Type": "application/json" }
        });
      } catch (error: any) {
        return new Response(JSON.stringify({
          success: false,
          message: "Bug修复失败",
          error: error?.message || "未知错误"
        }), {
          status: 400,
          headers: { "Content-Type": "application/json" }
        });
      }
    }
  },
  
  '/api/bugs/:id/close': {
    POST: async (req: Request) => {
      const url = new URL(req.url);
      const id = url.pathname.split('/')[3] || '';
      
      try {
        const result = bugService.close(parseInt(id));
        return new Response(JSON.stringify({
          success: true,
          message: "Bug关闭成功",
          data: result
        }), {
          headers: { "Content-Type": "application/json" }
        });
      } catch (error: any) {
        return new Response(JSON.stringify({
          success: false,
          message: "Bug关闭失败",
          error: error?.message || "未知错误"
        }), {
          status: 400,
          headers: { "Content-Type": "application/json" }
        });
      }
    }
  },

  // 管理API
  '/api/admin/modules': {
    GET: async (req: Request) => {
      try {
        const moduleRegistry = ModuleRegistry.getInstance();
        const modules = moduleRegistry.getAllModules();
        const stats = moduleRegistry.getStats();
        
        return new Response(JSON.stringify({
          success: true,
          data: { modules, stats }
        }), {
          headers: { "Content-Type": "application/json" }
        });
      } catch (error: any) {
        return new Response(JSON.stringify({
          success: false,
          message: "获取模块列表失败",
          error: error?.message || "未知错误"
        }), {
          status: 500,
          headers: { "Content-Type": "application/json" }
        });
      }
    }
  },

  '/api/admin/modules/:id/toggle': {
    POST: async (req: Request) => {
      const url = new URL(req.url);
      const id = url.pathname.split('/')[4] || '';
      
      try {
        const body = await req.json();
        const { enabled } = body;
        
        const moduleRegistry = ModuleRegistry.getInstance();
        const success = enabled 
          ? moduleRegistry.enableModule(id)
          : moduleRegistry.disableModule(id);
        
        if (success) {
          return new Response(JSON.stringify({
            success: true,
            message: `模块${enabled ? '启用' : '禁用'}成功`
          }), {
            headers: { "Content-Type": "application/json" }
          });
        } else {
          return new Response(JSON.stringify({
            success: false,
            message: "模块不存在"
          }), {
            status: 404,
            headers: { "Content-Type": "application/json" }
          });
        }
      } catch (error: any) {
        return new Response(JSON.stringify({
          success: false,
          message: "切换模块状态失败",
          error: error?.message || "未知错误"
        }), {
          status: 400,
          headers: { "Content-Type": "application/json" }
        });
      }
    }
  }
};

// API路由管理器
class ApiRouterManager {
  private static instance: ApiRouterManager;
  
  private constructor() {}
  
  static getInstance(): ApiRouterManager {
    if (!ApiRouterManager.instance) {
      ApiRouterManager.instance = new ApiRouterManager();
    }
    return ApiRouterManager.instance;
  }

  // 获取API路由配置
  getApiRoutes() {
    return apiRoutes;
  }

  // 检查API路由是否存在
  hasApiRoute(path: string): boolean {
    return path in apiRoutes;
  }

  // 获取API处理器
  getApiHandler(path: string): ApiHandler | null {
    return apiRoutes[path] || null;
  }

  // 处理API请求
  async handleApiRequest(req: Request): Promise<Response> {
    const url = new URL(req.url);
    const path = url.pathname;
    const method = req.method.toUpperCase();

    // 查找匹配的路由
    const handler = this.findMatchingHandler(path);
    
    if (!handler) {
      return new Response(JSON.stringify({
        success: false,
        message: "API路由不存在"
      }), {
        status: 404,
        headers: { "Content-Type": "application/json" }
      });
    }

    // 检查HTTP方法是否支持
    const methodHandler = handler[method as keyof ApiHandler];
    if (!methodHandler) {
      return new Response(JSON.stringify({
        success: false,
        message: `不支持的HTTP方法: ${method}`
      }), {
        status: 405,
        headers: { "Content-Type": "application/json" }
      });
    }

    try {
      // 调用对应的处理器
      const response = await methodHandler(req);
      return response;
    } catch (error: any) {
      logger.data.error('API处理失败: %s', error);
      return new Response(JSON.stringify({
        success: false,
        message: "服务器内部错误",
        error: error?.message || "未知错误"
      }), {
        status: 500,
        headers: { "Content-Type": "application/json" }
      });
    }
  }

  // 查找匹配的路由处理器（支持参数路由）
  private findMatchingHandler(path: string): ApiHandler | null {
    // 直接匹配
    if (apiRoutes[path]) {
      return apiRoutes[path];
    }

    // 参数路由匹配
    for (const routePath in apiRoutes) {
      if (this.matchRoute(routePath, path)) {
        return apiRoutes[routePath] || null;
      }
    }

    return null;
  }

  // 路由匹配（支持参数路由如 /api/bugs/:id/fix）
  private matchRoute(routePath: string, requestPath: string): boolean {
    const routeParts = routePath.split('/');
    const requestParts = requestPath.split('/');

    if (routeParts.length !== requestParts.length) {
      return false;
    }

    for (let i = 0; i < routeParts.length; i++) {
      if (routeParts[i].startsWith(':')) {
        // 参数路由，跳过检查
        continue;
      }
      if (routeParts[i] !== requestParts[i]) {
        return false;
      }
    }

    return true;
  }

  // 获取所有API路由路径
  getAllApiPaths(): string[] {
    return Object.keys(apiRoutes);
  }
}

export default ApiRouterManager; 