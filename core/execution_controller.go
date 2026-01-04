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

package core

import (
	"context"
	"time"

	"text2sql-skill/config"
)

type ExecutionController struct {
	cfg *config.Config
}

func NewExecutionController(cfg *config.Config) *ExecutionController {
	return &ExecutionController{cfg: cfg}
}

func (e *ExecutionController) GetExecutionContext(parent context.Context) (context.Context, context.CancelFunc) {
	// 解析超时时间
	timeout, err := time.ParseDuration(e.cfg.Execution.Timeout.Total)
	if err != nil {
		// 如果解析失败，使用默认值 10 秒
		timeout = 10 * time.Second
	}
	return context.WithTimeout(parent, timeout)
}

func (e *ExecutionController) CheckResourceLimits(inputSize int, estimatedRows int, estimatedMemoryMB float64) bool {
	return inputSize <= 10240 && // 10KB
		estimatedRows <= e.cfg.Security.ResourceLimits.MaxRows &&
		estimatedMemoryMB <= float64(e.cfg.Security.ResourceLimits.MaxMemoryMB)
}

func (e *ExecutionController) GetIsolationLevel() string {
	return e.cfg.Execution.IsolationLevel
}
