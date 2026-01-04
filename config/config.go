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
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 表示应用程序的完整配置
type Config struct {
	App            AppConfig            `yaml:"app"`
	Database       DatabaseConfig       `yaml:"database"`
	Security       SecurityConfig       `yaml:"security"`
	Execution      ExecutionConfig      `yaml:"execution"`
	Cache          CacheConfig          `yaml:"cache"`
	Audit          AuditConfig          `yaml:"audit"`
	Performance    PerformanceConfig    `yaml:"performance"`
	Monitoring     MonitoringConfig     `yaml:"monitoring"`
	Logging        LoggingConfig        `yaml:"logging"`
	Authentication AuthenticationConfig `yaml:"authentication"`
}

// AppConfig 应用程序基础配置
type AppConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
	Description string `yaml:"description"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string         `yaml:"driver"` // mysql, postgres
	MySQL    MySQLConfig    `yaml:"mysql"`
	Postgres PostgresConfig `yaml:"postgres"`
}

// MySQLConfig MySQL 数据库配置
type MySQLConfig struct {
	DSN     string        `yaml:"dsn"`
	Pool    PoolConfig    `yaml:"pool"`
	Timeout TimeoutConfig `yaml:"timeout"`
}

// PostgresConfig PostgreSQL 数据库配置
type PostgresConfig struct {
	DSN              string        `yaml:"dsn"`
	Pool             PoolConfig    `yaml:"pool"`
	Timeout          TimeoutConfig `yaml:"timeout"`
	SSLMode          string        `yaml:"ssl_mode"`
	BinaryParameters string        `yaml:"binary_parameters"`
}

// PoolConfig 连接池配置
type PoolConfig struct {
	MaxOpenConnections int    `yaml:"max_open_connections"`
	MaxIdleConnections int    `yaml:"max_idle_connections"`
	ConnMaxLifetime    string `yaml:"connection_max_lifetime"`
	ConnMaxIdleTime    string `yaml:"connection_max_idle_time"`
}

// TimeoutConfig 超时配置
type TimeoutConfig struct {
	Query      string `yaml:"query"`
	Connection string `yaml:"connection"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	Mode              string          `yaml:"mode"`
	AllowedOperations []string        `yaml:"allowed_operations"`
	ForbiddenKeywords []string        `yaml:"forbidden_keywords"`
	InputValidation   InputValidation `yaml:"input_validation"`
	ResourceLimits    ResourceLimits  `yaml:"resource_limits"`
}

// InputValidation 输入验证配置
type InputValidation struct {
	MaxLength  int     `yaml:"max_length"`
	MinEntropy float32 `yaml:"min_entropy"`
	MaxEntropy float32 `yaml:"max_entropy"`
}

// ResourceLimits 资源限制配置
type ResourceLimits struct {
	MaxMemoryMB     int `yaml:"max_memory_mb"`
	MaxRows         int `yaml:"max_rows"`
	MaxResultSizeMB int `yaml:"max_result_size_mb"`
}

// ExecutionConfig 执行配置
type ExecutionConfig struct {
	IsolationLevel string           `yaml:"isolation_level"`
	Timeout        ExecutionTimeout `yaml:"timeout"`
	Retry          RetryConfig      `yaml:"retry"`
}

// ExecutionTimeout 执行超时配置
type ExecutionTimeout struct {
	Total        string `yaml:"total"`
	QueryBuild   string `yaml:"query_build"`
	QueryExecute string `yaml:"query_execute"`
	ResultScan   string `yaml:"result_scan"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	Enabled           bool    `yaml:"enabled"`
	MaxAttempts       int     `yaml:"max_attempts"`
	InitialBackoff    string  `yaml:"initial_backoff"`
	MaxBackoff        string  `yaml:"max_backoff"`
	BackoffMultiplier float64 `yaml:"backoff_multiplier"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Size     int    `yaml:"size"`
	TTL      string `yaml:"ttl"`
	Strategy string `yaml:"strategy"`
}

// AuditConfig 审计配置
type AuditConfig struct {
	Enabled bool         `yaml:"enabled"`
	Level   string       `yaml:"level"`
	Storage AuditStorage `yaml:"storage"`
}

// AuditStorage 审计存储配置
type AuditStorage struct {
	Type     string        `yaml:"type"`
	Path     string        `yaml:"path"`
	Rotation AuditRotation `yaml:"rotation"`
}

// AuditRotation 审计日志轮转配置
type AuditRotation struct {
	MaxSizeMB  int  `yaml:"max_size_mb"`
	MaxAgeDays int  `yaml:"max_age_days"`
	MaxBackups int  `yaml:"max_backups"`
	Compress   bool `yaml:"compress"`
}

// PerformanceConfig 性能配置
type PerformanceConfig struct {
	AsyncProcessing bool              `yaml:"async_processing"`
	WorkerPoolSize  int               `yaml:"worker_pool_size"`
	BatchProcessing BatchProcessing   `yaml:"batch_processing"`
	Compression     CompressionConfig `yaml:"compression"`
}

// BatchProcessing 批处理配置
type BatchProcessing struct {
	Enabled       bool   `yaml:"enabled"`
	BatchSize     int    `yaml:"batch_size"`
	FlushInterval string `yaml:"flush_interval"`
}

// CompressionConfig 压缩配置
type CompressionConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Algorithm string `yaml:"algorithm"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Enabled     bool              `yaml:"enabled"`
	Metrics     MetricsConfig     `yaml:"metrics"`
	HealthCheck HealthCheckConfig `yaml:"health_check"`
}

// MetricsConfig 指标配置
type MetricsConfig struct {
	Enabled      bool     `yaml:"enabled"`
	Endpoint     string   `yaml:"endpoint"`
	PushInterval string   `yaml:"push_interval"`
	Collect      []string `yaml:"collect"`
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	Enabled       bool   `yaml:"enabled"`
	Port          int    `yaml:"port"`
	Path          string `yaml:"path"`
	CheckInterval string `yaml:"check_interval"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  string        `yaml:"level"`
	Format string        `yaml:"format"`
	Output string        `yaml:"output"`
	File   FileLogConfig `yaml:"file"`
}

// AuthenticationConfig 身份认证配置
type AuthenticationConfig struct {
	Enabled      bool   `yaml:"enabled"`
	Token        string `yaml:"token"`
	HeaderName   string `yaml:"header_name"`
	ValidateOnly bool   `yaml:"validate_only"`
}

// FileLogConfig 文件日志配置
type FileLogConfig struct {
	Path       string `yaml:"path"`
	MaxSizeMB  int    `yaml:"max_size_mb"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAgeDays int    `yaml:"max_age_days"`
	Compress   bool   `yaml:"compress"`
}

// LoadConfig 从 YAML 文件加载配置
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := ValidateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// GetActiveDatabaseConfig 获取当前激活的数据库配置
func (c *Config) GetActiveDatabaseConfig() (string, string) {
	switch c.Database.Driver {
	case "mysql":
		return "mysql", c.Database.MySQL.DSN
	case "postgres":
		return "postgres", c.Database.Postgres.DSN
	default:
		// 默认使用 MySQL
		return "mysql", c.Database.MySQL.DSN
	}
}

// GetActivePoolConfig 获取当前激活的连接池配置
func (c *Config) GetActivePoolConfig() PoolConfig {
	switch c.Database.Driver {
	case "mysql":
		return c.Database.MySQL.Pool
	case "postgres":
		return c.Database.Postgres.Pool
	default:
		return c.Database.MySQL.Pool
	}
}

// GetActiveTimeoutConfig 获取当前激活的超时配置
func (c *Config) GetActiveTimeoutConfig() TimeoutConfig {
	switch c.Database.Driver {
	case "mysql":
		return c.Database.MySQL.Timeout
	case "postgres":
		return c.Database.Postgres.Timeout
	default:
		return c.Database.MySQL.Timeout
	}
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:        "text2sql-skill",
			Version:     "1.0.0",
			Environment: "production",
			Description: "企业级 Text2SQL 技能引擎",
		},
		Database: DatabaseConfig{
			Driver: "mysql",
			MySQL: MySQLConfig{
				DSN: "app_user:secure_password@tcp(db-prod.example.com:3306)/sales_analytics?charset=utf8mb4&parseTime=True&loc=Local",
				Pool: PoolConfig{
					MaxOpenConnections: 20,
					MaxIdleConnections: 5,
					ConnMaxLifetime:    "30m",
					ConnMaxIdleTime:    "10m",
				},
				Timeout: TimeoutConfig{
					Query:      "5s",
					Connection: "3s",
				},
			},
			Postgres: PostgresConfig{
				DSN: "postgres://postgres:password@localhost:5432/mydb?sslmode=disable",
				Pool: PoolConfig{
					MaxOpenConnections: 20,
					MaxIdleConnections: 5,
					ConnMaxLifetime:    "30m",
					ConnMaxIdleTime:    "10m",
				},
				Timeout: TimeoutConfig{
					Query:      "5s",
					Connection: "3s",
				},
				SSLMode:          "disable",
				BinaryParameters: "yes",
			},
		},
		Security: SecurityConfig{
			Mode:              "read_only",
			AllowedOperations: []string{"SELECT"},
			ForbiddenKeywords: []string{
				"DROP", "DELETE", "INSERT", "UPDATE", "ALTER", "EXEC",
				"TRUNCATE", "CREATE", "GRANT", "REVOKE",
			},
			InputValidation: InputValidation{
				MaxLength:  2048,
				MinEntropy: 2.5,
				MaxEntropy: 6.0,
			},
			ResourceLimits: ResourceLimits{
				MaxMemoryMB:     50,
				MaxRows:         1000,
				MaxResultSizeMB: 10,
			},
		},
		Execution: ExecutionConfig{
			IsolationLevel: "full",
			Timeout: ExecutionTimeout{
				Total:        "10s",
				QueryBuild:   "2s",
				QueryExecute: "7s",
				ResultScan:   "1s",
			},
			Retry: RetryConfig{
				Enabled:           true,
				MaxAttempts:       3,
				InitialBackoff:    "100ms",
				MaxBackoff:        "2s",
				BackoffMultiplier: 1.5,
			},
		},
		Cache: CacheConfig{
			Enabled:  true,
			Size:     1000,
			TTL:      "5m",
			Strategy: "lru",
		},
		Audit: AuditConfig{
			Enabled: true,
			Level:   "detailed",
			Storage: AuditStorage{
				Type: "file",
				Path: "/var/log/text2sql/audit",
				Rotation: AuditRotation{
					MaxSizeMB:  1024,
					MaxAgeDays: 30,
					MaxBackups: 10,
					Compress:   true,
				},
			},
		},
		Performance: PerformanceConfig{
			AsyncProcessing: true,
			WorkerPoolSize:  4,
			BatchProcessing: BatchProcessing{
				Enabled:       true,
				BatchSize:     100,
				FlushInterval: "1s",
			},
			Compression: CompressionConfig{
				Enabled:   true,
				Algorithm: "zlib",
			},
		},
		Monitoring: MonitoringConfig{
			Enabled: true,
			Metrics: MetricsConfig{
				Enabled:      true,
				Endpoint:     "http://metrics.example.com/v1/metrics",
				PushInterval: "30s",
				Collect: []string{
					"requests_total",
					"errors_total",
					"duration_seconds",
					"cache_hits_total",
					"cache_misses_total",
					"db_connections",
					"memory_usage_bytes",
				},
			},
			HealthCheck: HealthCheckConfig{
				Enabled:       true,
				Port:          8080,
				Path:          "/health",
				CheckInterval: "30s",
			},
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
			File: FileLogConfig{
				Path:       "/var/log/text2sql/app.log",
				MaxSizeMB:  100,
				MaxBackups: 10,
				MaxAgeDays: 30,
				Compress:   true,
			},
		},
		Authentication: AuthenticationConfig{
			Enabled:      false,
			Token:        "your-secure-token-here",
			HeaderName:   "Authorization",
			ValidateOnly: false,
		},
	}
}

// parseDuration 解析持续时间字符串
func parseDuration(durationStr string) (time.Duration, error) {
	return time.ParseDuration(durationStr)
}
