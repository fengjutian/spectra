# Spectra - 日志收集与分析系统

Spectra 是一个现代化的日志收集和分析平台，专为实时监控和分析应用程序数据而设计。系统支持错误日志、性能指标、用户行为、自定义事件和页面停留时长等多种数据类型的收集和可视化。

## 🌟 特性

- **多类型数据收集**: 支持错误日志、性能指标、用户行为、自定义事件和页面停留时长
- **实时分析**: 基于 ClickHouse 的高性能数据存储和查询
- **RESTful API**: 提供完整的 API 接口用于数据收集和查询
- **现代化前端**: 使用 React + TypeScript + Vite 构建的用户界面
- **高性能**: 采用 Go 语言开发的后端服务，支持高并发处理
- **灵活配置**: 支持通过配置文件自定义系统参数

## 📁 项目结构

```
spectra/
├── spectra-backend/   # Go 后端服务
│   ├── config/        # 配置管理
│   ├── handlers/      # HTTP 处理器
│   ├── middleware/    # 中间件
│   ├── models/        # 数据模型
│   ├── repository/    # 数据访问层
│   ├── router/        # 路由定义
│   ├── services/      # 业务逻辑层
│   └── SQL/           # 数据库初始化脚本
└── spectra-frontend/  # React 前端应用
    ├── src/           # 源代码
    ├── public/        # 静态资源
    └── package.json   # 依赖配置
```

## 🚀 快速开始

### 后端服务

1. **环境要求**
   - Go 1.19+
   - ClickHouse 数据库
   - Git

2. **安装依赖**
   ```bash
   cd spectra-backend
   go mod download
   ```

3. **数据库初始化**
   ```bash
   # 使用提供的 SQL 脚本初始化 ClickHouse 数据库
   clickhouse-client < SQL/init.sql
   ```

4. **配置系统**
   编辑 `spectra-backend/config/config.yaml` 文件，配置数据库连接等参数。

5. **启动服务**
   ```bash
   go run main.go
   ```

### 前端应用

1. **环境要求**
   - Node.js 16+
   - npm 或 yarn

2. **安装依赖**
   ```bash
   cd spectra-frontend
   npm install
   ```

3. **开发模式**
   ```bash
   npm run dev
   ```

4. **构建生产版本**
   ```bash
   npm run build
   ```

## 📊 数据模型

### 1. 错误日志 (ErrorLog)
- **POST /api/error-logs** - 记录错误日志
- **GET /api/error-logs** - 查询错误日志列表

### 2. 性能指标 (PerformanceMetric)
- **POST /api/performance-metrics** - 记录性能指标
- **GET /api/performance-metrics** - 查询性能指标列表

### 3. 用户行为 (UserAction)
- **POST /api/user-actions** - 记录用户行为
- **GET /api/user-actions** - 查询用户行为列表

### 4. 自定义事件 (CustomEvent)
- **POST /api/custom-events** - 记录自定义事件
- **GET /api/custom-events** - 查询自定义事件列表

### 5. 页面停留时长 (PageStay)
- **POST /api/page-stays** - 记录页面停留时长
- **GET /api/page-stays/average** - 查询平均页面停留时长

## 🔧 配置说明

后端服务配置文件位于 `spectra-backend/config/config.yaml`，主要配置项包括：

```yaml
app:
  name: spectra-backend
  version: 1.0.0
  environment: development

server:
  port: 8080
  host: 0.0.0.0
  read_timeout: 15
  write_timeout: 15

log:
  level: info
  path: ./logs/app.log
  max_size: 500
  max_age: 30
  compress: true

db:
  driver: clickhouse
  host: localhost
  port: 9000
  database: spectra
  username: default
  password: ""
  debug: false
```

## 🛠️ 技术栈

### 后端
- **语言**: Go
- **Web 框架**: Gin
- **数据库**: ClickHouse
- **日志**: Zap
- **配置**: Viper

### 前端
- **框架**: React 18
- **语言**: TypeScript
- **构建工具**: Vite
- **样式**: CSS Modules

## 📚 API 文档

所有查询 API 都支持以下参数：
- `project_id` (必填) - 项目ID
- `start_time` (可选，默认24小时前) - 开始时间 (RFC3339格式)
- `end_time` (可选，默认当前时间) - 结束时间 (RFC3339格式)

详细的 API 文档请参考各子项目的 README 文件。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。

## 📄 许可证

MIT License

## 📞 联系方式

如有问题或建议，请通过以下方式联系：
- 提交 GitHub Issue
- 发送邮件至项目维护者

---

**Spectra** - 让日志分析变得简单而强大！