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

package tests

import (
	"testing"

	"text2sql-skill/config"
	"text2sql-skill/core"
)

func TestAuditLogger(t *testing.T) {
	// 使用内存中的模拟审计日志器，避免文件系统操作
	cfg := config.DefaultConfig()
	cfg.Audit.Storage.Type = "console" // 使用控制台输出代替文件
	cfg.Audit.Enabled = true
	cfg.Audit.Level = "detailed"

	logger := core.NewAuditLogger(cfg)
	defer logger.Close()

	// Test logging
	queryID := "test_query_123"
	logger.LogEvent(queryID, "test_event", map[string]interface{}{
		"test": "value",
		"num":  42,
	})

	// 验证日志器已初始化
	if logger == nil {
		t.Error("Audit logger failed to initialize")
	}
}
