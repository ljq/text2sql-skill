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

package config

import (
	"fmt"
)

// ValidateConfig 验证配置的合法性
func ValidateConfig(cfg *Config) error {
	// 验证应用配置
	if cfg.App.Name == "" {
		return fmt.Errorf("app.name cannot be empty")
	}
	if cfg.App.Version == "" {
		return fmt.Errorf("app.version cannot be empty")
	}
	if cfg.App.Environment == "" {
		return fmt.Errorf("app.environment cannot be empty")
	}

	// Validate database configuration (验证数据库配置)
	if cfg.Database.Driver == "" {
		return fmt.Errorf(`database.driver is required but not configured.

Please configure a database driver in your config.yaml:
1. For MySQL:   database.driver: "mysql"
2. For PostgreSQL: database.driver: "postgres"

If you don't have a database, consider:
1. Setting up a local MySQL/PostgreSQL instance
2. Using a cloud database service
3. Or implement file-based storage if needed

Example configuration:
database:
  driver: "mysql"
  mysql:
    dsn: "user:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"`)
	}

	// Validate configuration based on driver (根据驱动验证相应的配置)
	switch cfg.Database.Driver {
	case "mysql":
		if cfg.Database.MySQL.DSN == "" {
			return fmt.Errorf(`mysql.dsn is required but not configured.

Please configure a MySQL connection string in your config.yaml:

Format: user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local

Examples:
1. Local: root:password@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local
2. Cloud: user:password@tcp(db.localhost:3306)/production?charset=utf8mb4&parseTime=True&loc=Local`)
		}
		pool := cfg.Database.MySQL.Pool
		timeout := cfg.Database.MySQL.Timeout

		// Validate connection pool configuration (验证连接池配置)
		if pool.MaxOpenConnections <= 0 || pool.MaxOpenConnections > 100 {
			return fmt.Errorf("mysql.pool.max_open_connections must be between 1 and 100")
		}
		if pool.MaxIdleConnections < 0 || pool.MaxIdleConnections > pool.MaxOpenConnections {
			return fmt.Errorf("mysql.pool.max_idle_connections must be between 0 and max_open_connections")
		}

		// Validate duration strings (验证持续时间字符串)
		if _, err := parseDuration(pool.ConnMaxLifetime); err != nil {
			return fmt.Errorf("mysql.pool.connection_max_lifetime: %v", err)
		}
		if _, err := parseDuration(pool.ConnMaxIdleTime); err != nil {
			return fmt.Errorf("mysql.pool.connection_max_idle_time: %v", err)
		}
		if _, err := parseDuration(timeout.Query); err != nil {
			return fmt.Errorf("mysql.timeout.query: %v", err)
		}
		if _, err := parseDuration(timeout.Connection); err != nil {
			return fmt.Errorf("mysql.timeout.connection: %v", err)
		}

	case "postgres":
		if cfg.Database.Postgres.DSN == "" {
			return fmt.Errorf(`postgres.dsn is required but not configured.

Please configure a PostgreSQL connection string in your config.yaml:

Format: postgres://user:password@localhost:port/database?sslmode=disable

Examples:
1. Local: postgres://postgres:password@localhost:5432/mydb?sslmode=disable
2. Cloud: postgres://user:password@:5432/production?sslmode=require`)
		}
		pool := cfg.Database.Postgres.Pool
		timeout := cfg.Database.Postgres.Timeout

		// Validate connection pool configuration (验证连接池配置)
		if pool.MaxOpenConnections <= 0 || pool.MaxOpenConnections > 100 {
			return fmt.Errorf("postgres.pool.max_open_connections must be between 1 and 100")
		}
		if pool.MaxIdleConnections < 0 || pool.MaxIdleConnections > pool.MaxOpenConnections {
			return fmt.Errorf("postgres.pool.max_idle_connections must be between 0 and max_open_connections")
		}

		// Validate duration strings (验证持续时间字符串)
		if _, err := parseDuration(pool.ConnMaxLifetime); err != nil {
			return fmt.Errorf("postgres.pool.connection_max_lifetime: %v", err)
		}
		if _, err := parseDuration(pool.ConnMaxIdleTime); err != nil {
			return fmt.Errorf("postgres.pool.connection_max_idle_time: %v", err)
		}
		if _, err := parseDuration(timeout.Query); err != nil {
			return fmt.Errorf("postgres.timeout.query: %v", err)
		}
		if _, err := parseDuration(timeout.Connection); err != nil {
			return fmt.Errorf("postgres.timeout.connection: %v", err)
		}

		// Validate PostgreSQL specific configuration (验证 PostgreSQL 特定配置)
		switch cfg.Database.Postgres.SSLMode {
		case "disable", "require", "verify-ca", "verify-full":
			// Valid values (有效值)
		default:
			return fmt.Errorf("postgres.ssl_mode must be 'disable', 'require', 'verify-ca', or 'verify-full'")
		}

	default:
		return fmt.Errorf("unsupported database driver: %s. Supported drivers: mysql, postgres", cfg.Database.Driver)
	}

	// 验证安全配置
	switch cfg.Security.Mode {
	case "read_only", "read_write":
	default:
		return fmt.Errorf("security.mode must be 'read_only' or 'read_write'")
	}

	if len(cfg.Security.AllowedOperations) == 0 {
		return fmt.Errorf("security.allowed_operations cannot be empty")
	}

	if cfg.Security.InputValidation.MaxLength <= 0 {
		return fmt.Errorf("security.input_validation.max_length must be positive")
	}
	if cfg.Security.InputValidation.MinEntropy < 0 {
		return fmt.Errorf("security.input_validation.min_entropy cannot be negative")
	}
	if cfg.Security.InputValidation.MaxEntropy < cfg.Security.InputValidation.MinEntropy {
		return fmt.Errorf("security.input_validation.max_entropy must be >= min_entropy")
	}

	if cfg.Security.ResourceLimits.MaxMemoryMB <= 0 {
		return fmt.Errorf("security.resource_limits.max_memory_mb must be positive")
	}
	if cfg.Security.ResourceLimits.MaxRows <= 0 {
		return fmt.Errorf("security.resource_limits.max_rows must be positive")
	}
	if cfg.Security.ResourceLimits.MaxResultSizeMB <= 0 {
		return fmt.Errorf("security.resource_limits.max_result_size_mb must be positive")
	}

	// 验证执行配置
	switch cfg.Execution.IsolationLevel {
	case "none", "basic", "full":
	default:
		return fmt.Errorf("execution.isolation_level must be 'none', 'basic', or 'full'")
	}

	// 验证超时配置
	if _, err := parseDuration(cfg.Execution.Timeout.Total); err != nil {
		return fmt.Errorf("execution.timeout.total: %v", err)
	}
	if _, err := parseDuration(cfg.Execution.Timeout.QueryBuild); err != nil {
		return fmt.Errorf("execution.timeout.query_build: %v", err)
	}
	if _, err := parseDuration(cfg.Execution.Timeout.QueryExecute); err != nil {
		return fmt.Errorf("execution.timeout.query_execute: %v", err)
	}
	if _, err := parseDuration(cfg.Execution.Timeout.ResultScan); err != nil {
		return fmt.Errorf("execution.timeout.result_scan: %v", err)
	}

	// 验证重试配置
	if cfg.Execution.Retry.Enabled {
		if cfg.Execution.Retry.MaxAttempts <= 0 {
			return fmt.Errorf("execution.retry.max_attempts must be positive when retry is enabled")
		}
		if _, err := parseDuration(cfg.Execution.Retry.InitialBackoff); err != nil {
			return fmt.Errorf("execution.retry.initial_backoff: %v", err)
		}
		if _, err := parseDuration(cfg.Execution.Retry.MaxBackoff); err != nil {
			return fmt.Errorf("execution.retry.max_backoff: %v", err)
		}
		if cfg.Execution.Retry.BackoffMultiplier <= 0 {
			return fmt.Errorf("execution.retry.backoff_multiplier must be positive")
		}
	}

	// 验证缓存配置
	if cfg.Cache.Enabled {
		if cfg.Cache.Size <= 0 {
			return fmt.Errorf("cache.size must be positive when cache is enabled")
		}
		if _, err := parseDuration(cfg.Cache.TTL); err != nil {
			return fmt.Errorf("cache.ttl: %v", err)
		}
		switch cfg.Cache.Strategy {
		case "lru", "fifo", "lfu":
		default:
			return fmt.Errorf("cache.strategy must be 'lru', 'fifo', or 'lfu'")
		}
	}

	// 验证审计配置
	if cfg.Audit.Enabled {
		switch cfg.Audit.Level {
		case "none", "basic", "detailed":
		default:
			return fmt.Errorf("audit.level must be 'none', 'basic', or 'detailed'")
		}
		switch cfg.Audit.Storage.Type {
		case "file", "console":
		default:
			return fmt.Errorf("audit.storage.type must be 'file' or 'console'")
		}
		if cfg.Audit.Storage.Type == "file" && cfg.Audit.Storage.Path == "" {
			return fmt.Errorf("audit.storage.path cannot be empty when type is 'file'")
		}
	}

	// 验证性能配置
	if cfg.Performance.WorkerPoolSize <= 0 {
		return fmt.Errorf("performance.worker_pool_size must be positive")
	}
	if cfg.Performance.BatchProcessing.Enabled {
		if cfg.Performance.BatchProcessing.BatchSize <= 0 {
			return fmt.Errorf("performance.batch_processing.batch_size must be positive")
		}
		if _, err := parseDuration(cfg.Performance.BatchProcessing.FlushInterval); err != nil {
			return fmt.Errorf("performance.batch_processing.flush_interval: %v", err)
		}
	}
	if cfg.Performance.Compression.Enabled {
		switch cfg.Performance.Compression.Algorithm {
		case "zlib", "gzip", "none":
		default:
			return fmt.Errorf("performance.compression.algorithm must be 'zlib', 'gzip', or 'none'")
		}
	}

	// 验证监控配置
	if cfg.Monitoring.Enabled {
		if cfg.Monitoring.HealthCheck.Enabled {
			if cfg.Monitoring.HealthCheck.Port <= 0 || cfg.Monitoring.HealthCheck.Port > 65535 {
				return fmt.Errorf("monitoring.health_check.port must be between 1 and 65535")
			}
			if cfg.Monitoring.HealthCheck.Path == "" {
				return fmt.Errorf("monitoring.health_check.path cannot be empty")
			}
			if _, err := parseDuration(cfg.Monitoring.HealthCheck.CheckInterval); err != nil {
				return fmt.Errorf("monitoring.health_check.check_interval: %v", err)
			}
		}
	}

	// 验证日志配置
	switch cfg.Logging.Level {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("logging.level must be 'debug', 'info', 'warn', or 'error'")
	}
	switch cfg.Logging.Format {
	case "json", "text":
	default:
		return fmt.Errorf("logging.format must be 'json' or 'text'")
	}
	switch cfg.Logging.Output {
	case "stdout", "file":
	default:
		return fmt.Errorf("logging.output must be 'stdout' or 'file'")
	}
	if cfg.Logging.Output == "file" && cfg.Logging.File.Path == "" {
		return fmt.Errorf("logging.file.path cannot be empty when output is 'file'")
	}

	return nil
}
