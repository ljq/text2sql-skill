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

package utils

import (
	"context"
	"time"
)

func WithResourceLimits(ctx context.Context, memoryMB int, timeoutMs int) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutMs)*time.Millisecond)

	// Memory limit enforcement would be implemented at the OS level in production
	// This is a placeholder for the concept

	return ctx, func() {
		cancel()
		// Cleanup memory resources if needed
	}
}
