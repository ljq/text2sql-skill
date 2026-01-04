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
	"crypto/sha256"
	"encoding/binary"
	"unicode"
	"unicode/utf8"
)

type SemanticNode struct {
	Token     string
	Weight    float32
	Direction int8
	Links     [3]*SemanticNode
}

type SemanticTopology struct{}

func NewSemanticTopology() *SemanticTopology {
	return &SemanticTopology{}
}

func (s *SemanticTopology) BuildTopology(input []byte) *SemanticNode {
	root := &SemanticNode{Weight: 1.0}
	current := root

	for i := 0; i < len(input); {
		r, size := utf8.DecodeRune(input[i:])
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			i += size
			continue
		}

		token := make([]byte, 0, 16)
		j := i
		for j < len(input) {
			rj, sizej := utf8.DecodeRune(input[j:])
			if unicode.IsSpace(rj) || unicode.IsPunct(rj) {
				break
			}
			token = append(token, input[j:j+sizej]...)
			j += sizej
		}

		hash := sha256.Sum256(token)
		weight := float32(binary.LittleEndian.Uint32(hash[:4])) / 0xffffffff

		node := &SemanticNode{
			Token:  string(token),
			Weight: weight,
		}

		if current.Links[1] == nil {
			current.Links[1] = node
			node.Direction = 0
		} else if weight > current.Weight {
			node.Links[0] = current
			current.Links[1] = node
			node.Direction = 1
		} else {
			node.Links[1] = current
			current.Links[0] = node
			node.Direction = -1
		}

		current = node
		i = j
	}

	if root.Links[1] != nil {
		return root.Links[1]
	}
	return root
}

func (s *SemanticTopology) CalculateTopologyBalance(node *SemanticNode) float32 {
	if node == nil {
		return 0
	}
	left := s.CalculateTopologyBalance(node.Links[0])
	right := s.CalculateTopologyBalance(node.Links[1])
	return (left + right + float32(node.Direction)*node.Weight) / 3
}

func (s *SemanticTopology) GenerateTopologyFingerprint(node *SemanticNode) []byte {
	if node == nil {
		return []byte{0}
	}

	buf := make([]byte, 0, 256)
	buf = append(buf, []byte(node.Token)...)
	buf = append(buf, byte(node.Direction+1))
	buf = append(buf, 0) // null terminator

	for _, link := range node.Links {
		if link != nil {
			fingerprint := s.GenerateTopologyFingerprint(link)
			buf = append(buf, fingerprint...)
		}
	}

	hash := sha256.Sum256(buf)
	return hash[:8]
}
