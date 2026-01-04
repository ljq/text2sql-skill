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

package drivers

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	// 引入官方 PostgreSQL 驱动
	_ "github.com/lib/pq"
)

// RegisterPostgresDriver 注册官方 PostgreSQL 驱动
func RegisterPostgresDriver() string {
	// lib/pq 驱动已经自动注册为 "postgres"
	return "postgres"
}

// ValidatePostgresDSN 验证 PostgreSQL DSN 格式
func ValidatePostgresDSN(dsn string) error {
	// 基本验证：检查是否包含必要的连接参数
	if dsn == "" {
		return fmt.Errorf("PostgreSQL DSN cannot be empty")
	}

	// 检查是否包含 host 和 dbname
	if !containsParam(dsn, "host") {
		return fmt.Errorf("PostgreSQL DSN must contain 'host' parameter")
	}
	if !containsParam(dsn, "dbname") && !containsParam(dsn, "database") {
		return fmt.Errorf("PostgreSQL DSN must contain 'dbname' or 'database' parameter")
	}

	// 检查 SSL 模式（生产环境推荐 require 或 verify-full）
	if !containsParam(dsn, "sslmode") {
		return fmt.Errorf("PostgreSQL DSN should specify 'sslmode' parameter for security")
	}

	return nil
}

// containsParam 检查 DSN 是否包含指定参数
func containsParam(dsn, param string) bool {
	paramLower := param + "="
	return containsCaseInsensitive(dsn, paramLower)
}

// containsCaseInsensitive 不区分大小写地检查子字符串
func containsCaseInsensitive(s, substr string) bool {
	sLower := strings.ToLower(s)
	substrLower := strings.ToLower(substr)
	return strings.Contains(sLower, substrLower)
}

// CreatePostgresConnection 创建 PostgreSQL 连接
func CreatePostgresConnection(dsn string) (*sql.DB, error) {
	if err := ValidatePostgresDSN(dsn); err != nil {
		return nil, fmt.Errorf("invalid PostgreSQL DSN: %w", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	// 验证连接
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping PostgreSQL database: %w", err)
	}

	return db, nil
}
