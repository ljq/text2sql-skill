# Text2SQL-skill engine is born for AI agent

[简体中文](README.zh_CN.md)  

* A production-ready, secure, and high-performance Text-to-SQL skill engine for enterprise applications. Convert natural language queries into secure SQL queries with comprehensive safety guards and audit capabilities.

* The purpose of developing this project is to provide a relatively reliable enterprise-level Text2SQL solution for AI large language models without relying on the semantic capabilities of large language models. This is also the most critical capability engine for Agent.


## Why is Text2SQL the "difficulty" and a huge challenge for Agents?
* **High complexity**: Natural language is full of ambiguity (for example, "the latest order" may refer to time, geographical location, or priority), while SQL requires absolute precision. Agents need to simultaneously understand user intent, database schema, table relationships, data types, and even business rules. A minor error (such as a wrong JOIN condition) can lead to query failure, return of incorrect data, or even pose security risks (although modern frameworks provide injection protection).
* **High threshold for data-based Agents**: Most Agents on the market (such as customer service bots and task automation tools) are more "tool-oriented" because they focus on clearly defined actions (such as calling APIs and sending emails), with simple rules and high fault tolerance. However, "data-based" Agents (such as BI assistants and data analysis Agents) require deep integration with data sources, handling issues such as dynamic schemas, data drift, and performance optimization, leading to higher development and maintenance costs. Enterprises often prioritize low-hanging fruit, resulting in an immature ecosystem for data-based Agents.
* **Evaluation Challenges**: The accuracy of Text2SQL should not only be assessed based on syntactic correctness, but also on **Execution Accuracy** - whether the generated SQL returns correct results on a real database. This requires a large amount of annotated data and a testing environment, whereas the evaluation of tool-like Agents (such as task completion rate) is more intuitive.

## Feedback contact information
 
* **Author**: Jaco Liu   
* **Blog**: https://www.wdft.com   
* **Email**: ljqlab@gmail.com   
* **WeChat**: labsec  

## ✨ Features

### 🔒 **Security First**
- **Five-Layer Guard System**: Semantic analysis, permission control, execution control, schema evolution, and audit logging
- **Input Validation**: Maximum length, entropy analysis, and forbidden keyword detection
- **Resource Limits**: Strict control over memory usage, row counts, and result sizes
- **Read-Only Mode**: Configurable execution mode to prevent data modification

### 🔐 **Authentication**
- **MCP API Authentication**: Configurable token-based authentication for MCP API calls
- **Authorization Header**: Support for custom Authorization header names
- **Flexible Validation**: Optional token validation with `validate_only` mode
- **Secure Communication**: Recommendations for TLS/HTTPS in public network environments

#### **Security Configuration Example:**
```yaml
authentication:
  enabled: true        # Enable authentication for MCP API
  token: "your-secure-token-here"  # Authentication token
  header_name: "Authorization"  # HTTP header name for token
  validate_only: false  # Only validate token without requiring it

# Security Notes:
# 1. When enabled=true, all MCP requests must include the token in the Authorization header
# 2. For production environments, use strong, randomly generated tokens
# 3. For public network access, enable TLS/HTTPS for secure communication
# 4. Consider implementing more secure authentication protocols for sensitive applications
```

#### **MCP Request with Authentication:**
```bash
# Example curl request with authentication
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -H "Authorization: your-secure-token-here" \
  -d '{
    "jsonrpc": "2.0",
    "method": "text2sql/execute",
    "params": {
      "query": "Query all employees in the sales department"
    },
    "id": 1
  }'
```

#### **Security Best Practices:**
1. **Always enable authentication** when deploying in production environments
2. **Use strong, randomly generated tokens** (minimum 32 characters)
3. **Enable TLS/HTTPS** for all public network communications
4. **Regularly rotate tokens** for enhanced security
5. **Monitor audit logs** for unauthorized access attempts
6. **Implement rate limiting** to prevent brute force attacks
7. **Use network isolation** for sensitive database connections

### ⚡ **High Performance**
- **Intelligent Caching**: LRU/FIFO/LFU strategies with configurable TTL
- **Async Processing**: Non-blocking operations with worker pool
- **Connection Pooling**: Optimized database connection management
- **Batch Processing**: Efficient handling of multiple queries
- **Result Compression**: ZLIB/GZIP compression for large results

### 📊 **Observability**
- **Comprehensive Audit Logging**: File-based storage with rotation and compression
- **Health Checks**: HTTP health check endpoint for monitoring
- **Structured Logging**: JSON/text format with configurable levels

### 🔧 **Enterprise Support**
- **YAML Configuration**: Human-readable configuration with validation
- **Multi-Database Support**: MySQL and PostgreSQL with driver-specific optimizations
- **Graceful Shutdown**: Proper resource cleanup on termination signals
- **Concurrency Safe**: Mutex-protected shared resources
- **Error Recovery**: Comprehensive error handling and retry mechanisms

## 🚀 Quick Start

###  Basic configuration conditions and requirements
- Go 1.21 or later
- MySQL 5.7+ or PostgreSQL 12+
- 1GB RAM minimum (Recommend)

###### 💡Tip: I personally recommend using PostgresSQL as a priority.

### Installation

```bash
# Clone the repository
git clone https://github.com/ljq/text2sql-skill.git
cd text2sql-skill

# Install dependencies
go mod download

# Build the project
go build -o text2sql-skill main.go
```

### Configuration

Create a `config.yaml` file:

```yaml
app:
  name: "text2sql-skill"
  version: "1.0.0"
  environment: "production"
  description: "Enterprise Text2SQL Skill Engine"

database:
  driver: "mysql"  # or "postgres"
  mysql:
    dsn: "user:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"
    pool:
      max_open_connections: 20
      max_idle_connections: 5
    timeout:
      query: "5s"
      connection: "3s"

# See config.yaml.example for full configuration options
```

### Running

```bash
# Run with default config
./text2sql-skill

# Run with custom config
./text2sql-skill -config /path/to/config.yaml

# Run in development mode
./text2sql-skill -config config-dev.yaml
```

## 📖 Architecture

### Core Components

```
┌─────────────────────────────────────────────────────────────┐
│                    Text2SQL Skill Engine                    │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐     │
│  │  Input   │  │ Semantic │  │  Query   │  │  Result  │     │
│  │  Layer   │  │  Layer   │  │  Layer   │  │  Layer   │     │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘     │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────┐   │
│  │               Five-Layer Guard System                │   │
│  │ 1. Semantic Topology  │ 2. Permission Control        │   │
│  │ 3. Execution Control  │ 4. Schema Evolution          │   │
│  │ 5. Audit Logging      │                              │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Security Layers

1. **Semantic Topology**: Analyzes input structure and patterns
2. **Permission Controller**: Validates allowed operations and keywords
3. **Execution Controller**: Manages isolation levels and timeouts
4. **Schema Evolver**: Adapts to changing database schemas
5. **Audit Logger**: Records all operations for compliance

## ⚙️ Configuration

### Database Configuration

```yaml
database:
  driver: "mysql"  # Choose between "mysql" and "postgres"
  
  # MySQL specific configuration
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
  
  # PostgreSQL specific configuration  
  postgres:
    dsn: "postgres://user:password@localhost:port/database?sslmode=disable"
    ssl_mode: "disable"  # disable, require, verify-ca, verify-full
    binary_parameters: "yes"
```

### Security Configuration

```yaml
security:
  mode: "read_only"  # read_only or read_write
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

### Performance Configuration

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

## 🔍 Monitoring

### Health Check
```
GET http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2025-01-04 03:59:38",
  "database": "connected",
  "cache": "enabled",
  "uptime": "5m30s"
}
```

## 🧪 Testing & Examples

### Unit Tests
```bash
# Run all tests
go test ./...

# Run specific test suite
go test ./tests/...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests
```bash
# Requires database connection
go test -tags=integration ./tests/...
```

### End-to-End Tests
```bash
# Full system test
go test ./tests/endtoend_test.go
```

### Examples & Demos

The project includes comprehensive examples in the `examples/` directory:

#### 1. Skill Demo
```bash
# Run the skill demonstration
cd examples
go run skill_demo.go
```

This demo shows:
- Basic skill initialization and configuration
- Example natural language queries
- Security guard system in action
- Performance features demonstration

#### 2. MCP Server
```bash
# Start the MCP server
cd examples
go run mcp_server.go
```

The MCP server provides:
- HTTP JSON-RPC interface at `http://localhost:8080/mcp`
- Health check endpoint at `http://localhost:8080/health`
- Support for multiple MCP methods:
  - `text2sql/execute` - Execute natural language queries
  - `text2sql/capabilities` - Get skill capabilities
  - `text2sql/health` - Health check
  - `text2sql/config` - Get configuration

#### 3. MCP Client Demo
```bash
# Run the MCP client demonstration
cd examples
go run mcp_client_demo.go
```

This client demo shows:
- Connecting to the MCP server
- Executing queries via MCP protocol
- Testing security features
- Performance benchmarking
- Batch query processing

### MCP Protocol Support

Text2SQL Skill Engine supports the Model Context Protocol (MCP) for standardized AI tool integration:

#### MCP Methods:
- **text2sql/execute**: Execute natural language SQL queries
- **text2sql/capabilities**: Get skill metadata and capabilities
- **text2sql/health**: Health check endpoint
- **text2sql/config**: Get current configuration

#### Integration Example:
```json
{
  "jsonrpc": "2.0",
  "method": "text2sql/execute",
  "params": {
    "query": "Please check all employees in the sales department"
  },
  "id": 1
}
```

#### Response Format:
```json
{
  "jsonrpc": "2.0",
  "result": {
    "query_id": "q_ID123",
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

## 📁 Project Structure

```
text2sql-skill/
├── main.go                 # Application entry point
├── config.yaml             # Example configuration
├── config.yaml.example     # Full configuration example
├── config/
│   ├── config.go          # Configuration structures
│   └── validator.go       # Configuration validation
├── core/                   # Core engine components
│   ├── skill_impl.go      # Main skill implementation
│   ├── guard_system.go    # Five-layer guard system
│   ├── permission_controller.go
│   ├── execution_controller.go
│   ├── schema_evolver.go
│   ├── audit_logger.go
│   ├── query_cache.go
│   └── semantic_topology.go
├── drivers/               # Database drivers
│   ├── mysql_driver.go
│   └── postgres_driver.go
├── interfaces/           # Public interfaces
│   └── skill.go
├── utils/               # Utility functions
│   ├── crypto.go
│   ├── id_generator.go
│   └── resource_limiter.go
├── tests/              # Test suites
│   ├── audit_test.go
│   ├── guard_test.go
│   ├── permission_test.go
│   ├── semantic_test.go
│   └── endtoend_test.go
└── examples/           # Example code and demos
    ├── skill_demo.go      # Skill usage demonstration
    ├── mcp_server.go      # MCP server implementation
    ├── mcp_client_demo.go # MCP client demonstration
    └── go.mod            # Examples module
```

## LICENSE
This project adopts the [Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0) LICENSE.


## Disclaimer
This project is licensed under the Apache License 2.0. Use of this software for any illegal, harmful, or unethical purposes—including but not limited to cyber attacks, data theft, privacy violations, or malware distribution—is strictly prohibited. Users are solely responsible for ensuring lawful and ethical use; the author disclaim all liability for misuse.
