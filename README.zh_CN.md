# Text2SQL 技能引擎为 AI Agent 而生

[English](README.md)  

* 一个生产就绪、安全且高性能的文本到 SQL 技能引擎，适用于企业级应用。将自然语言查询转换为安全的 SQL 查询，提供全面的安全防护和审计能力。

* 编写此项目的目的是希望在不借助大语言模型语义能力的情况下，为AI大语言模型提供一种相对可靠的企业级Text2SQL解决方案，这也是Agent最关键的能力引擎。

## 为什么Text2SQL是Agent的“难点”和巨大挑战？
* **复杂性高**：自然语言充满歧义（比如“最近的订单”可能指时间、地理位置或优先级），而SQL要求绝对精确。Agent需要同时理解用户意图、数据库模式（schema）、表关系、数据类型，甚至业务规则。一个微小错误（如JOIN条件写错）就可能导致查询失败、返回错误数据，甚至引发安全风险（虽然现代框架会做注入防护）。
* **数据类Agent门槛高**：市面上大多数Agent（如客服机器人、任务自动化工具）偏“工具类”，因为它们聚焦于明确定义的动作（如调用API、发送邮件），规则简单、容错性高。而“数据类”Agent（如BI助手、数据分析Agent）需要深度集成数据源，处理动态schema、数据漂移、性能优化等问题，开发和维护成本更高。企业往往优先选择低垂果实（low-hanging fruit），导致数据类Agent生态还不成熟。
* **评估困难**：Text2SQL的准确率不能只看语法正确性，更要看**执行准确率**（Execution Accuracy）——生成的SQL是否在真实数据库上返回正确结果。这需要大量标注数据和测试环境，而工具类Agent的评估（如任务完成率）更直观。

### 交流反馈联系方式

* **作者**: Jaco Liu  
* **主页**: https://www.wdft.com  
* **邮箱**: ljqlab@gmail.com   
* **微信**: labsec  

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

#### **安全最佳实践（部分增强功能可自行实现）：**
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
git clone https://github.com/ljq/text2sql-skill.git
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
    dsn: "postgres://user:password@localhost:port/database?sslmode=disable"
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
  "timestamp": "2025-01-04 03:59:38",
  "database": "connected",
  "cache": "enabled",
  "uptime": "5m30s"
}
```

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
    "timestamp": "2025-01-04 03:59:38",
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

## LICENSE 许可证
[Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0)
