// Copyright 2024 Text2SQL Skill Engine
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Jaco Liu (Jianqiu Liu) <ljqlab@gmail.com>
// GitHub: https://github.com/ljq

// Text2SQL MCP 客户端演示
// 演示如何通过 MCP 协议调用 Text2SQL 技能

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// MCPRequest MCP 协议请求结构
type MCPRequest struct {
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
}

// MCPResponse MCP 协议响应结构
type MCPResponse struct {
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	ID      int         `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
}

// Text2SQLMCPClient Text2SQL MCP 客户端
type Text2SQLMCPClient struct {
	baseURL string
	client  *http.Client
}

// NewText2SQLMCPClient 创建新的 MCP 客户端
func NewText2SQLMCPClient(baseURL string) *Text2SQLMCPClient {
	return &Text2SQLMCPClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Call 调用 MCP 方法
func (c *Text2SQLMCPClient) Call(method string, params interface{}) (interface{}, error) {
	req := MCPRequest{
		Method:  method,
		Params:  params,
		ID:      1,
		JSONRPC: "2.0",
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	resp, err := c.client.Post(c.baseURL+"/mcp", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("HTTP 请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var mcpResp MCPResponse
	if err := json.Unmarshal(respBody, &mcpResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if mcpResp.Error != nil {
		return nil, fmt.Errorf("MCP 错误: %v", mcpResp.Error)
	}

	return mcpResp.Result, nil
}

// Execute 执行 Text2SQL 查询
func (c *Text2SQLMCPClient) Execute(query string) (interface{}, error) {
	params := map[string]interface{}{
		"query": query,
	}
	return c.Call("text2sql/execute", params)
}

// GetCapabilities 获取技能能力
func (c *Text2SQLMCPClient) GetCapabilities() (interface{}, error) {
	return c.Call("text2sql/capabilities", nil)
}

// GetHealth 获取健康状态
func (c *Text2SQLMCPClient) GetHealth() (interface{}, error) {
	return c.Call("text2sql/health", nil)
}

// GetConfig 获取配置信息
func (c *Text2SQLMCPClient) GetConfig() (interface{}, error) {
	return c.Call("text2sql/config", nil)
}

// main 客户端演示主函数
func main() {
	fmt.Println("==========================================")
	fmt.Println("🔌 Text2SQL MCP 客户端演示")
	fmt.Println("==========================================")

	// 创建客户端
	client := NewText2SQLMCPClient("http://localhost:8080")

	// 1. 测试连接和健康检查
	fmt.Println("\n1. 🩺 健康检查...")
	health, err := client.GetHealth()
	if err != nil {
		fmt.Printf("❌ 健康检查失败: %v\n", err)
		fmt.Println("💡 提示: 请确保 MCP 服务器正在运行")
		fmt.Println("   运行: go run examples/mcp_server.go")
		return
	}
	fmt.Printf("✅ 健康状态: %v\n", health)

	// 2. 获取能力信息
	fmt.Println("\n2. 📋 获取技能能力...")
	capabilities, err := client.GetCapabilities()
	if err != nil {
		fmt.Printf("❌ 获取能力失败: %v\n", err)
		return
	}
	capJSON, _ := json.MarshalIndent(capabilities, "", "  ")
	fmt.Printf("✅ 技能能力:\n%s\n", string(capJSON))

	// 3. 获取配置信息
	fmt.Println("\n3. ⚙️ 获取配置信息...")
	config, err := client.GetConfig()
	if err != nil {
		fmt.Printf("❌ 获取配置失败: %v\n", err)
		return
	}
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("✅ 配置信息:\n%s\n", string(configJSON))

	// 4. 执行示例查询
	fmt.Println("\n4. 🔍 执行示例查询...")
	examples := []string{
		"查询销售部门的所有员工",
		"获取上个月的销售额",
		"找出销售额最高的10个产品",
		"统计每个地区的客户数量",
	}

	for i, query := range examples {
		fmt.Printf("\n  示例 %d: %s\n", i+1, query)

		start := time.Now()
		result, err := client.Execute(query)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Printf("  ❌ 执行失败: %v\n", err)
			continue
		}

		resultJSON, _ := json.MarshalIndent(result, "", "    ")
		fmt.Printf("  ✅ 执行成功 (耗时: %v)\n", elapsed)
		fmt.Printf("  结果:\n%s\n", string(resultJSON))

		// 添加延迟
		time.Sleep(500 * time.Millisecond)
	}

	// 5. 测试安全防护
	fmt.Println("\n5. 🔒 测试安全防护...")
	forbiddenQueries := []string{
		"DROP TABLE users",
		"DELETE FROM customers",
		"SELECT * FROM users; DROP TABLE users",
	}

	for i, query := range forbiddenQueries {
		fmt.Printf("\n  测试 %d: %s\n", i+1, query)

		result, err := client.Execute(query)
		if err != nil {
			fmt.Printf("  ❌ 请求失败: %v\n", err)
			continue
		}

		resultJSON, _ := json.MarshalIndent(result, "", "    ")
		fmt.Printf("  响应:\n%s\n", string(resultJSON))
	}

	// 6. 性能测试
	fmt.Println("\n6. ⚡ 性能测试...")
	testQuery := "查询所有产品信息"
	iterations := 5
	var totalTime time.Duration

	for i := 0; i < iterations; i++ {
		start := time.Now()
		_, err := client.Execute(testQuery)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Printf("  迭代 %d: ❌ 失败 - %v\n", i+1, err)
		} else {
			fmt.Printf("  迭代 %d: ✅ 成功 - 耗时: %v\n", i+1, elapsed)
			totalTime += elapsed
		}

		time.Sleep(200 * time.Millisecond)
	}

	if iterations > 0 {
		avgTime := totalTime / time.Duration(iterations)
		fmt.Printf("\n  📊 平均响应时间: %v\n", avgTime)
	}

	// 7. 批量查询演示
	fmt.Println("\n7. 📦 批量查询演示...")
	batchQueries := []string{
		"查询产品库存",
		"获取客户列表",
		"统计订单数量",
	}

	fmt.Println("  发送批量查询请求...")
	for _, query := range batchQueries {
		go func(q string) {
			result, err := client.Execute(q)
			if err != nil {
				fmt.Printf("  查询 '%s': ❌ 失败\n", q)
			} else {
				fmt.Printf("  查询 '%s': ✅ 成功\n", q)
				_ = result // 忽略结果，仅用于演示
			}
		}(query)
	}

	// 等待批量查询完成
	time.Sleep(2 * time.Second)

	fmt.Println("\n==========================================")
	fmt.Println("🎉 MCP 客户端演示完成!")
	fmt.Println("==========================================")
	fmt.Println("\n📚 使用说明:")
	fmt.Println("  1. 启动 MCP 服务器:")
	fmt.Println("     go run examples/mcp_server.go")
	fmt.Println("  2. 运行客户端演示:")
	fmt.Println("     go run examples/mcp_client_demo.go")
	fmt.Println("  3. 或直接调用技能演示:")
	fmt.Println("     go run examples/skill_demo.go")
	fmt.Println("\n🔗 MCP 端点:")
	fmt.Println("  - HTTP: http://localhost:8080/mcp")
	fmt.Println("  - 健康检查: http://localhost:8080/health")
	fmt.Println("==========================================")
}
