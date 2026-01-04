# Text2SQL 技能引擎

[![Go 版本](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![许可证](https://img.shields.io/badge/许可证-MIT-blue.svg)](LICENSE)
[![Go 报告卡](https://goreportcard.com/badge/github.com/yourusername/text2sql-skill)](https://goreportcard.com/report/github.com/yourusername/text2sql-skill)

一个生产就绪、安全且高性能的文本到 SQL 技能引擎，适用于企业级应用。将自然语言查询转换为安全的 SQL 查询，提供全面的安全防护和审计能力。

**作者**: Jaco Liu | **主页**: https://github.com/ljq | **邮箱**: ljqlab@gmail.com | **微信**: labsec

## ✨ 特性

### 🔒 **安全第一**
- **五层防护系统**：语义分析、权限控制、执行控制、模式演进和审计日志
- **输入验证**：最大长度、熵分析和禁止关键字检测
- **资源限制**：严格的内存使用、行数和结果大小控制
- **只读模式**：可配置的执行模式，防止数据修改

### 🔐 **身份认证与授权**
- **MCP API 身份认证**：可配置的基于令牌的 MCP API 身份认证
- **授权头支持**：支持自定义 Authorization 头名称
- **灵活验证**：可选的令牌验证模式（`validate_only`）
- **安全通信**：公网环境推荐使用 TLS/HTTPS

#### **安全配置示例：**
```yaml
authentication:
  enabled: true        # 启用 MCP API 身份认证
  token: "your-secure-token-here"  # 身份认证令牌
  header_name: "Authorization"  # Token 的 HTTP 头名称
  validate_only: false  # 仅验证 Token 但不强制要求

# 安全注意事项：
# 1. 当 enabled=true 时，所有 MCP 请求必须在 Authorization 头中包含 token
# 2. 生产环境请使用强随机生成的 token
# 3. 公网访问时，请启用 TLS/HTTPS 确保通信安全
# 4. 对于敏感应用，考虑实现更安全的身份认证协议
```

#### **带身份认证的 MCP 请求示例：**
```bash
# 使用身份认证的 curl 请求示例
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: your-secure-token-here" \
  -d '{
    "jsonrpc": "2.0",
    "method": "text2sql/execute",
    "params": {
      "query": "查询销售部门的所有员工"
    },
    "id": 1
  }'
```

#### **安全最佳实践：**
1. **生产环境始终启用身份认证**
2. **使用强随机生成的令牌**（最少 32 个字符）
3. **公网通信启用 TLS/HTTPS**
4. **定期轮换令牌**以增强安全性
5. **监控审计日志**检测未授权访问尝试
6. **实施速率限制**防止暴力破解攻击
7. **使用网络隔离**保护敏感数据库连接

### ⚡ **高性能**
- **智能缓存**：最近最少使用/先进先出/最不经常使用策略，可配置生存时间
- **异步处理**：非阻塞操作，工作池支持
- **连接池**：优化的数据库连接管理
- **批处理**：高效处理多个查询
- **结果压缩**：ZLIB/GZIP 压缩大型结果

### 📊 **可观测性**
- **全面审计日志**：基于文件的存储，支持轮转和压缩
- **指标收集**：内置 Prometheus 指标端点
- **健康检查**：HTTP 健康检查端点
- **结构化日志**：JSON/文本格式，可配置级别

### 🔧 **企业就绪**
- **YAML 配置**：人类可读的配置，支持验证
- **多数据库支持**：MySQL 和 PostgreSQL，支持驱动特定优化
- **优雅关闭**：终止信号时的正确资源清理
- **并发安全**：互斥锁保护的共享资源
- **错误恢复**：全面的错误处理和重试机制

## 🚀 快速开始

### 先决条件
- Go 1.21 或更高版本
- MySQL 5.7+ 或 PostgreSQL 12+
- 最低 2GB RAM

### 安装

```bash
# 克隆仓库
git clone https://github.com/yourusername/text2sql-skill.git
cd text2sql-skill

# 安装依赖
go mod download

# 构建项目
go build -o text2sql-skill main.go
```

### 配置

创建 `config.yaml` 文件：

```yaml
app:
  name: "text2sql-skill"
  version: "1.0.0"
  environment: "production"
  description: "企业级 Text2SQL 技能引擎"

database:
  driver: "mysql"  # 或 "postgres"
  mysql:
    dsn: "user:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"
    pool:
      max_open_connections: 20
      max_idle_connections: 5
    timeout:
      query: "5s"
      connection: "3s"

# 查看 config.yaml.example 获取完整配置选项
```

### 运行

```bash
# 使用默认配置运行
./text2sql-skill

# 使用自定义配置运行
./text2sql-skill -config /path/to/config.yaml

# 开发模式运行
./text2sql-skill -config config-dev.yaml
```

## 📖 架构

### 核心组件

```
┌─────────────────────────────────────────────────────────────┐
│                    Text2SQL 技能引擎                         │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  输入层   │  │ 语义层   │  │ 查询层   │  │ 结果层   │    │
│  │  Input   │  │ Semantic │  │  Query   │  │  Result  │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────┐   │
│  │               五层防护系统                           │   │
│  │ 1. 语义拓扑  │ 2. 权限控制                          │   │
│  │ 3. 执行控制  │ 4. 模式演进                          │   │
│  │ 5. 审计日志  │                                     │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 安全层

1. **语义拓扑**：分析输入结构和模式
2. **权限控制器**：验证允许的操作和关键字
3. **执行控制器**：管理隔离级别和超时
4. **模式演进器**：适应变化的数据库模式
5. **审计记录器**：记录所有操作以供合规

## ⚙️ 配置

### 数据库配置

```yaml
database:
  driver: "mysql"  # 选择 "mysql" 或 "postgres"
  
  # MySQL 专用配置
  mysql:
    dsn: "user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local"
    pool:
      max_open_connections: 20
      max_idle_connections: 5
      connection_max_lifetime: "30m"
      connection_max_idle_time: "10m"
    timeout:
      query: "5s"
      connection: "3s"
  
  # PostgreSQL 专用配置  
  postgres:
    dsn: "postgres://user:password@host:port/database?sslmode=disable"
    ssl_mode: "disable"  # disable, require, verify-ca, verify-full
    binary_parameters: "yes"
```

### 安全配置

```yaml
security:
  mode: "read_only"  # read_only 或 read_write
  allowed_operations:
    - "SELECT"
  forbidden_keywords:
    - "DROP"
    - "DELETE"
    - "INSERT"
    - "UPDATE"
  input_validation:
    max_length: 2048
    min_entropy: 2.5
    max_entropy: 6.0
  resource_limits:
    max_memory_mb: 50
    max_rows: 1000
    max_result_size_mb: 10
```

### 性能配置

```yaml
performance:
  async_processing: true
  worker_pool_size: 4
  batch_processing:
    enabled: true
    batch_size: 100
    flush_interval: "1s"
  compression:
    enabled: true
    algorithm: "zlib"  # zlib, gzip, none
```

## 🔍 监控

### 健康检查
```
GET http://localhost:8080/health
```

响应：
```json
{
  "status": "healthy",
  "timestamp": "2024-01-04T03:59:38Z",
  "database": "connected",
  "cache": "enabled",
  "uptime": "5m30s"
}
```

### 指标端点
```
GET http://localhost:8080/metrics
```

可用指标：
- `text2sql_requests_total`
- `text2sql_errors_total`
- `text2sql_duration_seconds`
- `text2sql_cache_hits_total`
- `text2sql_cache_misses_total`
- `text2sql_db_connections`
- `text2sql_memory_usage_bytes`

## 🧪 测试与示例

### 单元测试
```bash
# 运行所有测试
go test ./...

# 运行特定测试套件
go test ./tests/...

# 运行覆盖率测试
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 集成测试
```bash
# 需要数据库连接
go test -tags=integration ./tests/...
```

### 端到端测试
```bash
# 完整系统测试
go test ./tests/endtoend_test.go
```

### 示例与演示

项目包含完整的示例代码，位于 `examples/` 目录：

#### 1. 技能演示
```bash
# 运行技能演示
cd examples
go run skill_demo.go
```

此演示展示：
- 基础技能初始化和配置
- 自然语言查询示例
- 安全防护系统实战
- 性能特性演示

#### 2. MCP 服务器
```bash
# 启动 MCP 服务器
cd examples
go run mcp_server.go
```

MCP 服务器提供：
- HTTP JSON-RPC 接口：`http://localhost:8080/mcp`
- 健康检查端点：`http://localhost:8080/health`
- 支持多种 MCP 方法：
  - `text2sql/execute` - 执行自然语言查询
  - `text2sql/capabilities` - 获取技能能力
  - `text2sql/health` - 健康检查
  - `text2sql/config` - 获取配置

#### 3. MCP 客户端演示
```bash
# 运行 MCP 客户端演示
cd examples
go run mcp_client_demo.go
```

此客户端演示展示：
- 连接到 MCP 服务器
- 通过 MCP 协议执行查询
- 测试安全特性
- 性能基准测试
- 批量查询处理

### MCP 协议支持

Text2SQL 技能引擎支持 Model Context Protocol (MCP)，用于标准化 AI 工具集成：

#### MCP 方法：
- **text2sql/execute**：执行自然语言 SQL 查询
- **text2sql/capabilities**：获取技能元数据和能力
- **text2sql/health**：健康检查端点
- **text2sql/config**：获取当前配置

#### 集成示例：
```json
{
  "jsonrpc": "2.0",
  "method": "text2sql/execute",
  "params": {
    "query": "查询销售部门的所有员工"
  },
  "id": 1
}
```

#### 响应格式：
```json
{
  "jsonrpc": "2.0",
  "result": {
    "query_id": "q_abc123",
    "status": "success",
    "timestamp": "2024-01-04T03:59:38Z",
    "duration_ms": 125,
    "result_size": 2048,
    "metadata": {
      "input_length": 15,
      "template_used": "SELECT * FROM employees WHERE department = 'sales'",
      "row_count": 42
    }
  },
  "id": 1
}
```

## 📁 项目结构

```
text2sql-skill/
├── main.go                 # 应用入口点
├── config.yaml             # 示例配置
├── config.yaml.example     # 完整配置示例
├── config/
│   ├── config.go          # 配置结构
│   └── validator.go       # 配置验证
├── core/                   # 核心引擎组件
│   ├── skill_impl.go      # 主技能实现
│   ├── guard_system.go    # 五层防护系统
│   ├── permission_controller.go
│   ├── execution_controller.go
│   ├── schema_evolver.go
│   ├── audit_logger.go
│   ├── query_cache.go
│   └── semantic_topology.go
├── drivers/               # 数据库驱动
│   ├── mysql_driver.go
│   └── postgres_driver.go
├── interfaces/           # 公共接口
│   └── skill.go
├── utils/               # 工具函数
│   ├── crypto.go
│   ├── id_generator.go
│   └── resource_limiter.go
├── tests/              # 测试套件
│   ├── audit_test.go
│   ├── guard_test.go
│   ├── permission_test.go
│   ├── semantic_test.go
│   └── endtoend_test.go
└── examples/           # 示例代码和演示
    ├── skill_demo.go      # 技能使用演示
    ├── mcp_server.go      # MCP 服务器实现
    ├── mcp_client_demo.go # MCP 客户端演示
    └── go.mod            # 示例模块
```

## 🔧 开发

### 从源码构建

```bash
# 克隆和设置
git clone https://github.com/yourusername/text2sql-skill.git
cd text2sql-skill

# 安装依赖
go mod tidy

# 构建
go build -o text2sql-skill .

# 运行测试
go test ./...
```

### 代码风格

- 遵循 Go 标准格式化：`gofmt -w .`
- 使用有意义的变量和函数名
- 为导出的函数和类型添加注释
- 为新功能编写单元测试

### 添加新功能

1. 创建功能分支：`git checkout -b feature/your-feature`
2. 实现更改并添加测试
3. 更新文档
4. 运行测试：`go test ./...`
5. 提交拉取请求

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎贡献！请随时提交拉取请求。

1. Fork 仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开拉取请求

### 贡献指南
- 编写清晰的提交信息
- 为新功能添加测试
- 根据需要更新文档
- 遵循现有的代码风格

## 🐛 故障排除

### 常见问题

**数据库连接失败**
```bash
ERROR: Failed to connect to MySQL database:
可能的问题：
1. MySQL 服务器未运行
2. DSN 格式不正确
3. 网络连接问题
4. 认证失败
```

**配置错误**
```bash
ERROR: Failed to load config: database.driver is required but not configured.
```

**性能问题**
- 检查连接池设置
- 监控缓存命中率
- 根据负载调整工作池大小

### 日志
默认情况下，日志写入 stdout。要使用文件日志，请配置：
```yaml
logging:
  output: "file"
  file:
    path: "/var/log/text2sql/app.log"
```

## 📞 支持

- **问题**：[GitHub Issues](https://github.com/yourusername/text2sql-skill/issues)
- **文档**：[项目 Wiki](https://github.com/yourusername/text2sql-skill/wiki)
- **邮箱**：ljqlab@gmail.com

## 🙏 致谢

- 感谢所有帮助塑造此项目的贡献者
- 受企业安全要求和最佳实践启发
- 在 Go 编程语言社区支持下构建

---

<div align="center">
  为安全高效的数据访问而构建 ❤️
</div>
