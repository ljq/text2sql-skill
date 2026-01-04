package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"text2sql-skill/config"
	"text2sql-skill/core"
	"text2sql-skill/drivers"
	"text2sql-skill/interfaces"
)

func main() {
	configPath := flag.String("config", "./config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("ERROR: Failed to load config: %v", err)
	}

	log.Printf("INFO: Config loaded: %s v%s", cfg.App.Name, cfg.App.Version)
	log.Printf("INFO: Environment: %s", cfg.App.Environment)

	// Get active database configuration (获取激活的数据库配置)
	dbDriver, dbDSN := cfg.GetActiveDatabaseConfig()
	poolConfig := cfg.GetActivePoolConfig()
	timeoutConfig := cfg.GetActiveTimeoutConfig()

	log.Printf("INFO: Database driver: %s", dbDriver)

	// Register database driver (注册数据库驱动)
	var driverName string
	switch dbDriver {
	case "mysql":
		driverName = drivers.RegisterMySQLDriver()
		log.Println("INFO: Using MySQL driver")
	case "postgres":
		driverName = drivers.RegisterPostgresDriver()
		log.Println("INFO: Using PostgreSQL driver")
	default:
		log.Fatalf("ERROR: Unsupported database driver: %s. Supported drivers: mysql, postgres", dbDriver)
	}

	// Create database connection (创建数据库连接)
	log.Printf("INFO: Connecting to database...")
	db, err := sql.Open(driverName, dbDSN)
	if err != nil {
		log.Fatalf("ERROR: Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Set connection pool settings (设置连接池配置)
	db.SetMaxOpenConns(poolConfig.MaxOpenConnections)
	db.SetMaxIdleConns(poolConfig.MaxIdleConnections)

	// Test database connection with timeout (测试数据库连接，带超时)
	pingCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("INFO: Testing database connection...")
	if err := db.PingContext(pingCtx); err != nil {
		// Provide helpful error messages based on database type (根据数据库类型提供有用的错误信息)
		switch dbDriver {
		case "mysql":
			log.Fatalf(`ERROR: Failed to connect to MySQL database:
	Error: %v
	Possible issues:
	1. MySQL server is not running
	2. Incorrect DSN format: %s
	3. Network connectivity issues
	4. Authentication failed
	
	DSN format: user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local
	Example: root:password@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local`, err, dbDSN)
		case "postgres":
			log.Fatalf(`ERROR: Failed to connect to PostgreSQL database:
	Error: %v
	Possible issues:
	1. PostgreSQL server is not running
	2. Incorrect DSN format: %s
	3. Network connectivity issues
	4. Authentication failed
	
	DSN format: postgres://user:password@host:port/database?sslmode=disable
	Example: postgres://postgres:password@localhost:5432/mydb?sslmode=disable`, err, dbDSN)
		default:
			log.Fatalf("ERROR: Failed to ping database: %v", err)
		}
	}

	log.Println("INFO: Database connection successful")

	// Create skill (创建技能)
	log.Println("INFO: Creating Text2SQL skill...")
	skill, err := core.NewText2SQLSkill(cfg, db)
	if err != nil {
		log.Fatalf("ERROR: Failed to create skill: %v", err)
	}

	// Setup graceful shutdown (设置优雅关闭)
	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-shutdownCtx.Done()
		log.Println("INFO: Shutting down gracefully...")
		if err := skill.(interfaces.Skill).SafeShutdown(); err != nil {
			log.Printf("WARN: Error during shutdown: %v", err)
		}
	}()

	// Display startup information (显示启动信息)
	log.Println("============================================================")
	log.Printf("INFO: Text2SQL Skill Engine Started")
	log.Printf("INFO: Name: %s v%s", cfg.App.Name, cfg.App.Version)
	log.Printf("INFO: Environment: %s", cfg.App.Environment)
	log.Printf("INFO: Description: %s", cfg.App.Description)
	log.Println("----------------------------------------")
	log.Printf("INFO: Database: %s", dbDriver)
	log.Printf("INFO: Connections: %d open, %d idle", poolConfig.MaxOpenConnections, poolConfig.MaxIdleConnections)
	log.Printf("INFO: Query timeout: %s", timeoutConfig.Query)
	log.Println("----------------------------------------")
	log.Printf("INFO: Security mode: %s", cfg.Security.Mode)
	log.Printf("INFO: Allowed operations: %v", cfg.Security.AllowedOperations)
	log.Printf("INFO: Forbidden keywords: %d keywords", len(cfg.Security.ForbiddenKeywords))
	log.Println("----------------------------------------")
	log.Printf("INFO: Cache: %v (size: %d, TTL: %s)", cfg.Cache.Enabled, cfg.Cache.Size, cfg.Cache.TTL)
	log.Printf("INFO: Audit logging: %v (level: %s)", cfg.Audit.Enabled, cfg.Audit.Level)
	log.Printf("INFO: Performance: async=%v, workers=%d", cfg.Performance.AsyncProcessing, cfg.Performance.WorkerPoolSize)
	log.Println("============================================================")
	log.Println("INFO: Service is ready and waiting for requests...")
	log.Println("INFO: Press Ctrl+C to shutdown")

	// Keep the service running (保持服务运行)
	<-shutdownCtx.Done()
	log.Println("INFO: Service stopped gracefully")
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}
