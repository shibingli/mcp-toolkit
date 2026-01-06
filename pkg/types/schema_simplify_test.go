// Copyright 2024 MCP Toolkit Authors
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

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSimplifyDescription 测试精简描述功能 / Test simplify description function
func TestSimplifyDescription(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "空描述 / Empty description",
			input:    "",
			expected: "",
		},
		{
			name:     "去掉关键词 / Remove keywords",
			input:    "Create a new file. Keywords: create, new file, write file",
			expected: "Create a new file.",
		},
		{
			name:     "去掉使用场景 / Remove use cases",
			input:    "Create a new file. Use this tool when you need to: 1) Create files, 2) Write content",
			expected: "Create a new file.",
		},
		{
			name:     "去掉中文部分 / Remove Chinese part",
			input:    "Create a new file / 创建新文件",
			expected: "Create a new file",
		},
		{
			name:     "去掉IMPORTANT部分 / Remove IMPORTANT part",
			input:    "Create a new file. IMPORTANT: This is the primary tool.",
			expected: "Create a new file.",
		},
		{
			name:     "复杂描述精简 / Complex description simplification",
			input:    "CREATE A NEW FILE with specified content. Use this tool when you need to: 1) Create a new file from scratch, 2) Write initial content. IMPORTANT: This is the primary tool. Keywords: create, new file, write file. / 创建新文件",
			expected: "CREATE A NEW FILE with specified content.",
		},
		{
			name:     "去掉Examples部分 / Remove Examples part",
			input:    "The file path. Examples: file.txt. / 文件路径",
			expected: "The file path.",
		},
		{
			name:     "超长描述截断 / Truncate long description",
			input:    "This is a very long description that exceeds the maximum length limit and should be truncated to ensure it doesn't consume too many tokens when sent to the language model for processing",
			expected: "This is a very long description that exceeds the maximum length limit and should be truncated to ensure it doesn't consume too many tokens when sent t...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SimplifyDescription(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSimplifyProperty 测试精简属性功能 / Test simplify property function
func TestSimplifyProperty(t *testing.T) {
	prop := Property{
		Type:        "string",
		Description: "The file path to create. Keywords: path, file. / 文件路径",
		MinLength:   intPtr(1),
		MaxLength:   intPtr(255),
		Examples:    []any{"file1.txt", "file2.txt", "file3.txt"},
		Pattern:     "^[a-z]+$",
		Format:      "uri",
	}

	simplified := SimplifyProperty(prop)

	// 验证保留的字段 / Verify retained fields
	assert.Equal(t, "string", simplified.Type)
	assert.Equal(t, "The file path to create.", simplified.Description)

	// 验证去掉的字段 / Verify removed fields
	assert.Nil(t, simplified.MinLength)
	assert.Nil(t, simplified.MaxLength)
	assert.Empty(t, simplified.Examples)
	assert.Empty(t, simplified.Pattern)
	assert.Empty(t, simplified.Format)
}

// TestSimplifySchema 测试精简Schema功能 / Test simplify schema function
func TestSimplifySchema(t *testing.T) {
	schema := JSONSchema{
		Type:        "object",
		Description: "Create a new file with content. Keywords: create, file. / 创建文件",
		Properties: map[string]Property{
			"path": {
				Type:        "string",
				Description: "The file path. Examples: file.txt. / 文件路径",
				MinLength:   intPtr(1),
				Examples:    []any{"file1.txt", "file2.txt"},
			},
			"content": {
				Type:        "string",
				Description: "The file content. / 文件内容",
				MaxLength:   intPtr(1000),
			},
		},
		Required:             []string{"path", "content"},
		AdditionalProperties: boolPtr(false),
	}

	simplified := SimplifySchema(schema)

	// 验证基本字段 / Verify basic fields
	assert.Equal(t, "object", simplified.Type)
	assert.Equal(t, "Create a new file with content.", simplified.Description)
	assert.Equal(t, []string{"path", "content"}, simplified.Required)
	assert.Nil(t, simplified.AdditionalProperties)

	// 验证属性被精简 / Verify properties are simplified
	assert.Equal(t, 2, len(simplified.Properties))
	pathProp := simplified.Properties["path"]
	assert.Equal(t, "string", pathProp.Type)
	assert.Equal(t, "The file path.", pathProp.Description)
	assert.Nil(t, pathProp.MinLength)
	assert.Empty(t, pathProp.Examples)

	contentProp := simplified.Properties["content"]
	assert.Equal(t, "string", contentProp.Type)
	assert.Equal(t, "The file content.", contentProp.Description)
	assert.Nil(t, contentProp.MaxLength)
}

// boolPtr 返回bool指针 / Returns bool pointer
func boolPtr(b bool) *bool {
	return &b
}
