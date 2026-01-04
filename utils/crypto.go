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
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/binary"
	"math"
)

func EncryptResult(data []map[string]interface{}, compress bool) []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(0x7F) // magic number

	for _, row := range data {
		buf.WriteByte(0x01) // row start
		for k, v := range row {
			// Key obfuscation
			ks := []byte(k)
			for i := range ks {
				ks[i] ^= 0xAA
			}
			buf.Write(ks)
			buf.WriteByte(0x1F) // key-value separator

			// Value encoding
			switch val := v.(type) {
			case int64:
				buf.WriteByte(0x02) // INT type
				var b [8]byte
				binary.LittleEndian.PutUint64(b[:], uint64(val))
				buf.Write(b[:])
			case float64:
				buf.WriteByte(0x03) // FLOAT type
				var b [8]byte
				binary.LittleEndian.PutUint64(b[:], math.Float64bits(val))
				buf.Write(b[:])
			case string:
				buf.WriteByte(0x04) // STRING type
				buf.Write([]byte(val))
			}
			buf.WriteByte(0x1E) // field end
		}
		buf.WriteByte(0x00) // row end
	}

	// Add checksum
	checksum := sha256.Sum256(buf.Bytes())
	buf.Write(checksum[:4])

	if compress && buf.Len() > 1024 {
		return compressData(buf.Bytes())
	}

	return buf.Bytes()
}

func compressData(data []byte) []byte {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write(data)
	w.Close()
	return buf.Bytes()
}
