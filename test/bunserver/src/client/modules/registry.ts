// 浏览器兼容的 logger
const logger = {
  data: {
    info: (message: string, ...args: any[]) => console.log(`[INFO] ${message}`, ...args),
    warn: (message: string, ...args: any[]) => console.warn(`[WARN] ${message}`, ...args),
    error: (message: string, ...args: any[]) => console.error(`[ERROR] ${message}`, ...args)
  }
};

// 模块接口定义
export interface Module {
  id: string;
  name: string;
  description: string;
  version: string;
  routes: string[];
  apis: string[];
  icon?: string;
  enabled: boolean;
}

// 模块管理器
class ModuleRegistry {
  private static instance: ModuleRegistry;
  private modules: Map<string, Module> = new Map();

  private constructor() {
    this.registerDefaultModules();
  }

  static getInstance(): ModuleRegistry {
    if (!ModuleRegistry.instance) {
      ModuleRegistry.instance = new ModuleRegistry();
    }
    return ModuleRegistry.instance;
  }

  // 注册默认模块
  private registerDefaultModules() {
    this.registerModule({
      id: 'bugs',
      name: '问题追踪',
      description: '管理和跟踪项目中的问题与缺陷',
      version: '1.0.0',
      routes: ['/view/bugs'],
      apis: ['/api/bugs', '/api/bugs/:id/fix', '/api/bugs/:id/close'],
      icon: '🐛',
      enabled: true
    });

    // 可以添加更多默认模块
    this.registerModule({
      id: 'admin',
      name: '管理',
      description: '管理系统用户和权限',
      version: '1.0.0',
      routes: ['/view/admin'],
      apis: ['/api/bugs'],
      icon: '👥',
      enabled: true
    });

    // this.registerModule({
    //   id: 'projects',
    //   name: '项目管理',
    //   description: '管理项目信息和配置',
    //   version: '1.0.0',
    //   routes: ['/view/projects'],
    //   apis: ['/api/projects'],
    //   icon: '📁',
    //   enabled: true
    // });
  }

  // 注册模块
  registerModule(module: Module): void {
    if (this.modules.has(module.id)) {
      logger.data.warn('模块已存在，将被覆盖: %s', module.id);
    }
    
    this.modules.set(module.id, module);
    logger.data.info('模块注册成功: %s (%s)', module.name, module.id);
  }

  // 获取模块
  getModule(id: string): Module | undefined {
    return this.modules.get(id);
  }

  // 获取所有模块
  getAllModules(): Module[] {
    return Array.from(this.modules.values());
  }

  // 获取启用的模块
  getEnabledModules(): Module[] {
    return Array.from(this.modules.values()).filter(module => module.enabled);
  }

  // 启用模块
  enableModule(id: string): boolean {
    const module = this.modules.get(id);
    if (module) {
      module.enabled = true;
      logger.data.info('模块已启用: %s', id);
      return true;
    }
    return false;
  }

  // 禁用模块
  disableModule(id: string): boolean {
    const module = this.modules.get(id);
    if (module) {
      module.enabled = false;
      logger.data.info('模块已禁用: %s', id);
      return true;
    }
    return false;
  }

  // 检查路由是否属于某个模块
  getModuleByRoute(route: string): Module | undefined {
    for (const module of this.modules.values()) {
      if (module.enabled && module.routes.includes(route)) {
        return module;
      }
    }
    return undefined;
  }

  // 检查API是否属于某个模块
  getModuleByApi(api: string): Module | undefined {
    for (const module of this.modules.values()) {
      if (module.enabled && module.apis.some(moduleApi => 
        this.matchApiPattern(moduleApi, api)
      )) {
        return module;
      }
    }
    return undefined;
  }

  // API模式匹配（支持参数路由）
  private matchApiPattern(pattern: string, api: string): boolean {
    const patternParts = pattern.split('/');
    const apiParts = api.split('/');

    if (patternParts.length !== apiParts.length) {
      return false;
    }

    for (let i = 0; i < patternParts.length; i++) {
      if (patternParts[i].startsWith(':')) {
        // 参数路由，跳过检查
        continue;
      }
      if (patternParts[i] !== apiParts[i]) {
        return false;
      }
    }

    return true;
  }

  // 获取模块统计信息
  getStats() {
    const allModules = this.getAllModules();
    const enabledModules = this.getEnabledModules();
    
    return {
      total: allModules.length,
      enabled: enabledModules.length,
      disabled: allModules.length - enabledModules.length,
      modules: allModules.map(module => ({
        id: module.id,
        name: module.name,
        enabled: module.enabled,
        routeCount: module.routes.length,
        apiCount: module.apis.length
      }))
    };
  }
}

export default ModuleRegistry; 