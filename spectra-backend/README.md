# Spectra Backend

## 项目概述
Spectra Backend 是一个日志收集和分析系统的后端服务，支持错误日志、性能指标、用户行为、自定义事件和页面停留时长等数据的收集和查询。

## 项目结构

```
spectra-backend/
├── config/          # 配置相关
│   ├── config.go    # 配置结构体和加载逻辑
│   └── config.yaml  # 配置文件
├── handlers/        # HTTP处理器
│   └── log_handler.go
├── middleware/      # 中间件
│   └── logger.go
├── models/          # 数据模型
│   └── models.go
├── repository/      # 数据访问层
│   ├── repository.go
│   └── clickhouse_repository.go
├── router/          # 路由
│   └── routes.go
├── services/        # 业务逻辑层
│   └── log_service.go
├── SQL/             # SQL脚本
│   └── init.sql
├── main.go          # 程序入口
├── go.mod
└── go.sum
```

## 数据模型

### 1. ErrorLog (错误日志)
- **POST /api/error-logs** - 记录错误日志
- **GET /api/error-logs** - 查询错误日志列表

### 2. PerformanceMetric (性能指标)
- **POST /api/performance-metrics** - 记录性能指标
- **GET /api/performance-metrics** - 查询性能指标列表

### 3. UserAction (用户行为)
- **POST /api/user-actions** - 记录用户行为
- **GET /api/user-actions** - 查询用户行为列表

### 4. CustomEvent (自定义事件)
- **POST /api/custom-events** - 记录自定义事件
- **GET /api/custom-events** - 查询自定义事件列表

### 5. PageStay (页面停留时长)
- **POST /api/page-stays** - 记录页面停留时长
- **GET /api/page-stays/average** - 查询平均页面停留时长

## 查询参数
所有查询API都支持以下参数：
- `project_id` (必填) - 项目ID
- `start_time` (可选，默认24小时前) - 开始时间 (RFC3339格式)
- `end_time` (可选，默认当前时间) - 结束时间 (RFC3339格式)

## 配置说明
配置文件位于 `config/config.yaml`，主要配置项包括：

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

## 启动服务

1. 确保 ClickHouse 数据库已安装并运行
2. 使用 init.sql 创建必要的数据表
3. 运行以下命令启动服务：

```bash
go run main.go
```

## 依赖说明
- **gin-gonic/gin** - Web框架
- **ClickHouse/clickhouse-go/v2** - ClickHouse驱动
- **spf13/viper** - 配置管理
- **uber-go/zap** - 日志库