// Text2SQL MCP 服务器实现
// 提供通过 Model Context Protocol (MCP) 访问 Text2SQL 技能的功能

package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"text2sql-skill/config"
	"text2sql-skill/core"
	"text2sql-skill/interfaces"
)

// MCPRequest MCP 协议请求结构
type MCPRequest struct {
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
}

// MCPResponse MCP 协议响应结构
type MCPResponse struct {
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
}

// MCPError MCP 协议错误结构
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// Text2SQLMCPServer Text2SQL MCP 服务器
type Text2SQLMCPServer struct {
	skill interfaces.Skill
	cfg   *config.Config
}

// NewText2SQLMCPServer 创建新的 MCP 服务器
func NewText2SQLMCPServer(cfg *config.Config, skill interfaces.Skill) *Text2SQLMCPServer {
	return &Text2SQLMCPServer{
		skill: skill,
		cfg:   cfg,
	}
}

// HandleRequest 处理 MCP 请求
func (s *Text2SQLMCPServer) HandleRequest(req MCPRequest) MCPResponse {
	switch req.Method {
	case "text2sql/execute":
		return s.handleExecute(req)
	case "text2sql/capabilities":
		return s.handleCapabilities(req)
	case "text2sql/health":
		return s.handleHealth(req)
	case "text2sql/config":
		return s.handleConfig(req)
	default:
		return MCPResponse{
			ID:      req.ID,
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not found",
				Data:    req.Method,
			},
		}
	}
}

// handleExecute 处理执行请求
func (s *Text2SQLMCPServer) handleExecute(req MCPRequest) MCPResponse {
	var params struct {
		Query string `json:"query"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return MCPResponse{
			ID:      req.ID,
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
				Data:    err.Error(),
			},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	startTime := time.Now()
	result, err := s.skill.Execute(ctx, params.Query)
	elapsed := time.Since(startTime)

	if err != nil {
		return MCPResponse{
			ID:      req.ID,
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32000,
				Message: "Execution failed",
				Data:    err.Error(),
			},
		}
	}

	response := map[string]interface{}{
		"query_id":    result.QueryID,
		"status":      result.Status,
		"timestamp":   result.Timestamp.Format(time.RFC3339),
		"duration_ms": elapsed.Milliseconds(),
		"result_size": len(result.Result),
	}

	if len(result.Meta) > 0 {
		var meta map[string]interface{}
		if err := json.Unmarshal(result.Meta, &meta); err == nil {
			response["metadata"] = meta
		}
	}

	// 如果结果较小，直接包含在响应中
	if len(result.Result) < 1024 {
		response["result"] = string(result.Result)
	}

	return MCPResponse{
		ID:      req.ID,
		JSONRPC: "2.0",
		Result:  response,
	}
}

// handleCapabilities 处理能力查询请求
func (s *Text2SQLMCPServer) handleCapabilities(req MCPRequest) MCPResponse {
	capabilities := map[string]interface{}{
		"skill_id": s.skill.CapabilityID(),
		"methods": []string{
			"text2sql/execute",
			"text2sql/capabilities",
			"text2sql/health",
			"text2sql/config",
		},
		"security": map[string]interface{}{
			"mode":                     s.cfg.Security.Mode,
			"allowed_operations":       s.cfg.Security.AllowedOperations,
			"forbidden_keywords_count": len(s.cfg.Security.ForbiddenKeywords),
		},
		"authentication": map[string]interface{}{
			"enabled":       s.cfg.Authentication.Enabled,
			"header_name":   s.cfg.Authentication.HeaderName,
			"validate_only": s.cfg.Authentication.ValidateOnly,
		},
		"performance": map[string]interface{}{
			"cache_enabled":       s.cfg.Cache.Enabled,
			"async_processing":    s.cfg.Performance.AsyncProcessing,
			"compression_enabled": s.cfg.Performance.Compression.Enabled,
		},
	}

	return MCPResponse{
		ID:      req.ID,
		JSONRPC: "2.0",
		Result:  capabilities,
	}
}

// handleHealth 处理健康检查请求
func (s *Text2SQLMCPServer) handleHealth(req MCPRequest) MCPResponse {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"skill":     s.skill.CapabilityID(),
		"version":   s.cfg.App.Version,
	}

	return MCPResponse{
		ID:      req.ID,
		JSONRPC: "2.0",
		Result:  health,
	}
}

// handleConfig 处理配置查询请求
func (s *Text2SQLMCPServer) handleConfig(req MCPRequest) MCPResponse {
	configInfo := map[string]interface{}{
		"app": map[string]interface{}{
			"name":        s.cfg.App.Name,
			"version":     s.cfg.App.Version,
			"environment": s.cfg.App.Environment,
		},
		"security": map[string]interface{}{
			"mode":               s.cfg.Security.Mode,
			"max_input_length":   s.cfg.Security.InputValidation.MaxLength,
			"max_rows":           s.cfg.Security.ResourceLimits.MaxRows,
			"max_memory_mb":      s.cfg.Security.ResourceLimits.MaxMemoryMB,
			"max_result_size_mb": s.cfg.Security.ResourceLimits.MaxResultSizeMB,
		},
		"authentication": map[string]interface{}{
			"enabled":       s.cfg.Authentication.Enabled,
			"header_name":   s.cfg.Authentication.HeaderName,
			"validate_only": s.cfg.Authentication.ValidateOnly,
		},
		"performance": map[string]interface{}{
			"cache_enabled":       s.cfg.Cache.Enabled,
			"cache_size":          s.cfg.Cache.Size,
			"cache_ttl":           s.cfg.Cache.TTL,
			"async_processing":    s.cfg.Performance.AsyncProcessing,
			"worker_pool_size":    s.cfg.Performance.WorkerPoolSize,
			"compression_enabled": s.cfg.Performance.Compression.Enabled,
		},
	}

	return MCPResponse{
		ID:      req.ID,
		JSONRPC: "2.0",
		Result:  configInfo,
	}
}

// HTTPHandler HTTP 处理器
func (s *Text2SQLMCPServer) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 身份认证验证
	if s.cfg.Authentication.Enabled {
		token := r.Header.Get(s.cfg.Authentication.HeaderName)
		if token == "" {
			// 如果启用了身份认证但没有提供token，返回认证错误
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(MCPResponse{
				JSONRPC: "2.0",
				Error: &MCPError{
					Code:    -32600,
					Message: "Authentication required",
					Data:    "Missing Authorization header",
				},
			})
			return
		}

		// 验证token
		if token != s.cfg.Authentication.Token {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(MCPResponse{
				JSONRPC: "2.0",
				Error: &MCPError{
					Code:    -32600,
					Message: "Authentication failed",
					Data:    "Invalid token",
				},
			})
			return
		}
	}

	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	resp := s.HandleRequest(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// StartServer 启动 MCP 服务器
func (s *Text2SQLMCPServer) StartServer(addr string) error {
	http.HandleFunc("/mcp", s.HTTPHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
		})
	})

	log.Printf("MCP 服务器启动在 %s", addr)
	return http.ListenAndServe(addr, nil)
}

// StartUnixSocketServer 启动 Unix Socket 服务器
func (s *Text2SQLMCPServer) StartUnixSocketServer(socketPath string) error {
	// 尝试创建监听器，如果socket文件已存在，会返回错误
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		// 如果是因为socket文件已存在而失败，记录错误但不删除文件
		log.Printf("ERROR: 无法创建Unix socket监听器: %v", err)
		log.Printf("INFO: 如果socket文件 %s 已存在，请手动删除或使用不同的路径", socketPath)
		return err
	}
	defer listener.Close()

	log.Printf("MCP Unix Socket 服务器启动在 %s", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("接受连接错误: %v", err)
			continue
		}

		go s.handleSocketConnection(conn)
	}
}

// handleSocketConnection 处理 Socket 连接
func (s *Text2SQLMCPServer) handleSocketConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var req MCPRequest
		if err := decoder.Decode(&req); err != nil {
			break
		}

		resp := s.HandleRequest(req)
		if err := encoder.Encode(resp); err != nil {
			break
		}
	}
}

// RunMCPServer 运行 MCP 服务器
func RunMCPServer() {
	// 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		cfg = config.DefaultConfig()
		log.Printf("使用默认配置: %s v%s", cfg.App.Name, cfg.App.Version)
	}

	// 创建技能实例（演示目的，使用nil数据库连接）
	var skill interfaces.Skill
	skill, err = core.NewText2SQLSkill(cfg, nil)
	if err != nil {
		log.Fatalf("创建技能失败: %v", err)
	}

	// 创建 MCP 服务器
	server := NewText2SQLMCPServer(cfg, skill)

	// 启动服务器
	// 可以选择 HTTP 或 Unix Socket
	addr := ":8080"
	log.Printf("启动 Text2SQL MCP 服务器...")
	log.Printf("HTTP 端点: http://localhost%s/mcp", addr)
	log.Printf("健康检查: http://localhost%s/health", addr)
	log.Printf("技能ID: %s", skill.CapabilityID())

	if err := server.StartServer(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
