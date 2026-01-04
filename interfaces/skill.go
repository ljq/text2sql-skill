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

package interfaces

import (
	"context"
	"time"
)

type Skill interface {
	CapabilityID() string
	Execute(ctx context.Context, input string) (SkillResult, error)
	SafeShutdown() error
}

type SkillResult struct {
	QueryID   string    `json:"query_id"`
	Result    []byte    `json:"result"`
	Meta      []byte    `json:"meta"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}
