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

	"text2sql-skill/core"
)

func TestSemanticTopology(t *testing.T) {
	topology := core.NewSemanticTopology()

	tests := []struct {
		input    string
		expected int
	}{
		{"北京销售额", 2},
		{"2025年北京客户", 3},
		{"上海", 1},
		{"", 1}, // Root node
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			node := topology.BuildTopology([]byte(tt.input))
			count := countNodes(node)
			if count != tt.expected {
				t.Errorf("Input: %q, Expected nodes: %d, Got: %d", tt.input, tt.expected, count)
			}
		})
	}
}

func TestTopologyBalance(t *testing.T) {
	topology := core.NewSemanticTopology()

	tests := []struct {
		input  string
		minBal float32
		maxBal float32
	}{
		{"北京销售额", -0.5, 0.5},
		{"2025年北京客户", -0.5, 0.5},
		{"上海", -0.1, 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			node := topology.BuildTopology([]byte(tt.input))
			balance := topology.CalculateTopologyBalance(node)
			if balance < tt.minBal || balance > tt.maxBal {
				t.Errorf("Input: %q, Balance: %f, Expected range: [%f, %f]",
					tt.input, balance, tt.minBal, tt.maxBal)
			}
		})
	}
}

func countNodes(node *core.SemanticNode) int {
	if node == nil {
		return 0
	}
	count := 1
	for _, link := range node.Links {
		count += countNodes(link)
	}
	return count
}
