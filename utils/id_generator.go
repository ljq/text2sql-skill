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
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
	"time"
)

func GenerateQueryID() string {
	now := time.Now().UnixNano()
	pid := int32(os.Getpid())
	tid := int32(runtime.NumGoroutine())

	buf := make([]byte, 16)
	binary.LittleEndian.PutUint64(buf[:8], uint64(now))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(pid))
	binary.LittleEndian.PutUint32(buf[12:16], uint32(tid))

	hash := sha256.Sum256(buf)
	return fmt.Sprintf("%x", hash[:12])
}
