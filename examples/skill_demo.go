package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"text2sql-skill/config"
	"text2sql-skill/core"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		// 使用默认配置作为示例
		cfg = config.DefaultConfig()
		log.Printf("使用默认配置: %s v%s", cfg.App.Name, cfg.App.Version)
	}

	// 2. 创建数据库连接
	// 注意：这是一个演示，实际使用时需要真实的数据库连接
	// 这里我们创建一个nil连接，实际项目中应该使用：
	// db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database")
	// 或
	// db, err := sql.Open("postgres", "postgres://user:password@localhost/database")
	var db *sql.DB

	// 演示目的：使用nil连接，实际功能需要真实数据库
	log.Println("INFO: 演示模式 - 使用模拟数据库连接")
	// 在实际部署中，请取消注释以下代码并配置正确的数据库连接
	/*
		db, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/database")
		if err != nil {
			log.Fatalf("ERROR: 数据库连接失败: %v", err)
		}
		defer db.Close()
	*/

	// 3. 创建Text2SQL技能实例
	skill, err := core.NewText2SQLSkill(cfg, db)
	if err != nil {
		log.Fatalf("ERROR: 创建技能失败: %v", err)
	}
	defer skill.SafeShutdown()

	// 4. 显示技能信息
	fmt.Println("==========================================")
	fmt.Printf("🎯 Text2SQL 技能演示\n")
	fmt.Printf("📋 技能ID: %s\n", skill.CapabilityID())
	fmt.Printf("🔒 安全模式: %s\n", cfg.Security.Mode)
	fmt.Printf("✅ 允许的操作: %v\n", cfg.Security.AllowedOperations)
	fmt.Println("==========================================")

	// 5. 执行示例查询
	examples := []string{
		"查询销售部门的所有员工",
		"获取上个月的销售额",
		"找出销售额最高的10个产品",
		"统计每个地区的客户数量",
		"分析产品库存情况",
	}

	ctx := context.Background()

	for i, query := range examples {
		fmt.Printf("\n🔍 示例 %d: %s\n", i+1, query)

		start := time.Now()
		result, err := skill.Execute(ctx, query)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Printf("❌ 执行失败: %v\n", err)
			continue
		}

		// 显示结果
		fmt.Printf("✅ 查询ID: %s\n", result.QueryID)
		fmt.Printf("📊 状态: %s\n", result.Status)
		fmt.Printf("⏱️ 执行时间: %v\n", elapsed)
		fmt.Printf("📅 时间戳: %v\n", result.Timestamp.Format("2006-01-02 15:04:05"))

		// 解析元数据
		if len(result.Meta) > 0 {
			fmt.Printf("📋 元数据: %s\n", string(result.Meta))
		}

		// 显示结果大小
		if len(result.Result) > 0 {
			fmt.Printf("💾 结果大小: %d 字节\n", len(result.Result))

			// 在实际使用中，可以解密和显示结果
			if len(result.Result) < 100 {
				fmt.Printf("🔓 结果预览: %s\n", string(result.Result))
			}
		}

		// 添加延迟以便观察
		time.Sleep(500 * time.Millisecond)
	}

	// 6. 演示安全防护
	fmt.Println("\n==========================================")
	fmt.Println("🔒 安全防护演示")
	fmt.Println("==========================================")

	// 尝试执行被禁止的操作
	forbiddenQueries := []string{
		"DROP TABLE users",                      // 包含禁止关键字 DROP
		"DELETE FROM customers",                 // 包含禁止关键字 DELETE
		"SELECT * FROM users; DROP TABLE users", // 注入攻击
		"非常长的查询语句" + string(make([]byte, 3000)), // 超过最大长度限制
	}

	for i, query := range forbiddenQueries {
		fmt.Printf("\n🚫 测试防护 %d: %s\n", i+1, query[:min(50, len(query))])

		result, err := skill.Execute(ctx, query)
		if err != nil {
			fmt.Printf("❌ 错误: %v\n", err)
			continue
		}

		fmt.Printf("📊 状态: %s\n", result.Status)
		if result.Status == "rejected" {
			fmt.Printf("✅ 成功拦截: %s\n", string(result.Meta))
		}
	}

	// 7. 性能演示
	fmt.Println("\n==========================================")
	fmt.Println("⚡ 性能特性演示")
	fmt.Println("==========================================")

	fmt.Printf("💾 缓存: %v\n", cfg.Cache.Enabled)
	if cfg.Cache.Enabled {
		fmt.Printf("  策略: %s, 大小: %d, TTL: %s\n",
			cfg.Cache.Strategy, cfg.Cache.Size, cfg.Cache.TTL)
	}

	fmt.Printf("🚀 异步处理: %v\n", cfg.Performance.AsyncProcessing)
	if cfg.Performance.AsyncProcessing {
		fmt.Printf("  工作池大小: %d\n", cfg.Performance.WorkerPoolSize)
	}

	fmt.Printf("📦 批处理: %v\n", cfg.Performance.BatchProcessing.Enabled)
	if cfg.Performance.BatchProcessing.Enabled {
		fmt.Printf("  批大小: %d, 刷新间隔: %s\n",
			cfg.Performance.BatchProcessing.BatchSize,
			cfg.Performance.BatchProcessing.FlushInterval)
	}

	fmt.Printf("🗜️ 压缩: %v\n", cfg.Performance.Compression.Enabled)
	if cfg.Performance.Compression.Enabled {
		fmt.Printf("  算法: %s\n", cfg.Performance.Compression.Algorithm)
	}

	fmt.Println("\n==========================================")
	fmt.Println("🎉 演示完成!")
	fmt.Println("==========================================")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
