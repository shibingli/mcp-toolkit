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

// JSONSchema JSON Schema 结构体 / JSON Schema structure
// 注意：properties 和 required 字段不使用 omitempty，确保始终输出
// 这是为了兼容 llama.cpp 等工具的 Jinja2 模板
type JSONSchema struct {
	Type        string              `json:"type"`
	Properties  map[string]Property `json:"properties"`
	Required    []string            `json:"required"`
	Description string              `json:"description,omitempty"`
}

// Property JSON Schema 属性 / JSON Schema property
type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description,omitempty"`
	Items       *Items   `json:"items,omitempty"`
	Enum        []string `json:"enum,omitempty"`
	Default     any      `json:"default,omitempty"`
}

// Items JSON Schema 数组项 / JSON Schema array items
type Items struct {
	Type string `json:"type"`
}

// ToolSchemas 工具输入模式定义 / Tool input schema definitions
var ToolSchemas = map[string]JSONSchema{
	"create_file": {
		Type: "object",
		Properties: map[string]Property{
			"path":    {Type: "string", Description: "文件路径 / File path"},
			"content": {Type: "string", Description: "文件内容 / File content"},
		},
		Required: []string{"path", "content"},
	},
	"create_directory": {
		Type: "object",
		Properties: map[string]Property{
			"path": {Type: "string", Description: "目录路径 / Directory path"},
		},
		Required: []string{"path"},
	},
	"read_file": {
		Type: "object",
		Properties: map[string]Property{
			"path": {Type: "string", Description: "文件路径 / File path"},
		},
		Required: []string{"path"},
	},
	"write_file": {
		Type: "object",
		Properties: map[string]Property{
			"path":    {Type: "string", Description: "文件路径 / File path"},
			"content": {Type: "string", Description: "文件内容 / File content"},
		},
		Required: []string{"path", "content"},
	},
	"delete": {
		Type: "object",
		Properties: map[string]Property{
			"path": {Type: "string", Description: "要删除的文件或目录的路径（自动判断类型）/ Path of file or directory to delete (auto-detect type)"},
		},
		Required: []string{"path"},
	},
	"delete_file": {
		Type: "object",
		Properties: map[string]Property{
			"path": {Type: "string", Description: "要删除的文件路径 / Path of file to delete"},
		},
		Required: []string{"path"},
	},
	"delete_directory": {
		Type: "object",
		Properties: map[string]Property{
			"path":      {Type: "string", Description: "要删除的目录路径 / Path of directory to delete"},
			"recursive": {Type: "boolean", Description: "是否递归删除所有子目录和文件，默认为true / Whether to recursively delete all subdirectories and files, default is true"},
		},
		Required: []string{"path"},
	},
	"copy": {
		Type: "object",
		Properties: map[string]Property{
			"source":      {Type: "string", Description: "源文件或目录路径 / Source file or directory path"},
			"destination": {Type: "string", Description: "目标路径 / Destination path"},
		},
		Required: []string{"source", "destination"},
	},
	"copy_file": {
		Type: "object",
		Properties: map[string]Property{
			"source":      {Type: "string", Description: "源文件路径 / Source file path"},
			"destination": {Type: "string", Description: "目标文件路径 / Destination file path"},
		},
		Required: []string{"source", "destination"},
	},
	"copy_directory": {
		Type: "object",
		Properties: map[string]Property{
			"source":      {Type: "string", Description: "源目录路径 / Source directory path"},
			"destination": {Type: "string", Description: "目标目录路径 / Destination directory path"},
		},
		Required: []string{"source", "destination"},
	},
	"move": {
		Type: "object",
		Properties: map[string]Property{
			"source":      {Type: "string", Description: "源文件或目录路径 / Source file or directory path"},
			"destination": {Type: "string", Description: "目标路径 / Destination path"},
		},
		Required: []string{"source", "destination"},
	},
	"move_file": {
		Type: "object",
		Properties: map[string]Property{
			"source":      {Type: "string", Description: "源文件路径 / Source file path"},
			"destination": {Type: "string", Description: "目标文件路径 / Destination file path"},
		},
		Required: []string{"source", "destination"},
	},
	"move_directory": {
		Type: "object",
		Properties: map[string]Property{
			"source":      {Type: "string", Description: "源目录路径 / Source directory path"},
			"destination": {Type: "string", Description: "目标目录路径 / Destination directory path"},
		},
		Required: []string{"source", "destination"},
	},
	"list_directory": {
		Type: "object",
		Properties: map[string]Property{
			"path": {Type: "string", Description: "目录路径 / Directory path"},
		},
		Required: []string{"path"},
	},
	"search_files": {
		Type: "object",
		Properties: map[string]Property{
			"path":    {Type: "string", Description: "搜索路径 / Search path"},
			"pattern": {Type: "string", Description: "搜索模式(支持通配符) / Search pattern (supports wildcards)"},
		},
		Required: []string{"path", "pattern"},
	},
	"batch_delete": {
		Type: "object",
		Properties: map[string]Property{
			"paths": {Type: "array", Description: "文件或目录路径列表 / List of file or directory paths", Items: &Items{Type: "string"}},
		},
		Required: []string{"paths"},
	},
	"file_stat": {
		Type: "object",
		Properties: map[string]Property{
			"path": {Type: "string", Description: "文件路径 / File path"},
		},
		Required: []string{"path"},
	},
	"file_exists": {
		Type: "object",
		Properties: map[string]Property{
			"path": {Type: "string", Description: "文件路径 / File path"},
		},
		Required: []string{"path"},
	},
	"get_current_time": {
		Type: "object",
		Properties: map[string]Property{
			"timezone": {Type: "string", Description: "时区名称（可选），如 'Asia/Shanghai'、'America/New_York'、'Europe/London'。为空则使用系统本地时区。/ Time zone name (optional), e.g. 'Asia/Shanghai', 'America/New_York', 'Europe/London'. Empty means use system local timezone."},
		},
		Required: []string{},
	},
	"execute_command": {
		Type: "object",
		Properties: map[string]Property{
			"command":  {Type: "string", Description: "要执行的命令 / Command to execute"},
			"args":     {Type: "array", Description: "命令参数 / Command arguments", Items: &Items{Type: "string"}},
			"work_dir": {Type: "string", Description: "工作目录(相对于沙箱根目录) / Working directory (relative to sandbox root)"},
			"timeout":  {Type: "integer", Description: "超时时间(秒),0表示不限制 / Timeout in seconds, 0 means no limit"},
		},
		Required: []string{"command", "work_dir"},
	},
	"get_command_blacklist": {
		Type:       "object",
		Properties: map[string]Property{},
		Required:   []string{},
	},
	"update_command_blacklist": {
		Type: "object",
		Properties: map[string]Property{
			"commands":    {Type: "array", Description: "要添加的黑名单命令 / Commands to add to blacklist", Items: &Items{Type: "string"}},
			"directories": {Type: "array", Description: "要添加的黑名单目录 / Directories to add to blacklist", Items: &Items{Type: "string"}},
		},
		Required: []string{},
	},
	"get_working_directory": {
		Type:       "object",
		Properties: map[string]Property{},
		Required:   []string{},
	},
	"change_directory": {
		Type: "object",
		Properties: map[string]Property{
			"path": {Type: "string", Description: "目标目录路径(相对于沙箱根目录) / Target directory path (relative to sandbox root)"},
		},
		Required: []string{"path"},
	},
	"execute_command_async": {
		Type: "object",
		Properties: map[string]Property{
			"command":          {Type: "string", Description: "要执行的命令 / Command to execute"},
			"args":             {Type: "array", Description: "命令参数 / Command arguments", Items: &Items{Type: "string"}},
			"work_dir":         {Type: "string", Description: "工作目录 / Working directory"},
			"timeout":          {Type: "integer", Description: "超时时间(秒) / Timeout in seconds"},
			"permission_level": {Type: "integer", Description: "权限级别 / Permission level"},
			"user":             {Type: "string", Description: "执行用户 / Executing user"},
		},
		Required: []string{"command", "work_dir"},
	},
	"get_command_task_status": {
		Type: "object",
		Properties: map[string]Property{
			"task_id": {Type: "string", Description: "任务ID / Task ID"},
		},
		Required: []string{"task_id"},
	},
	"cancel_command_task": {
		Type: "object",
		Properties: map[string]Property{
			"task_id": {Type: "string", Description: "任务ID / Task ID"},
		},
		Required: []string{"task_id"},
	},
	"list_command_tasks": {
		Type:       "object",
		Properties: map[string]Property{},
		Required:   []string{},
	},
	"get_command_history": {
		Type: "object",
		Properties: map[string]Property{
			"limit": {Type: "integer", Description: "返回的最大记录数 / Maximum number of records to return"},
		},
		Required: []string{},
	},
	"set_permission_level": {
		Type: "object",
		Properties: map[string]Property{
			"level": {Type: "integer", Description: "权限级别(0-3) / Permission level (0-3)"},
		},
		Required: []string{"level"},
	},
	"get_permission_level": {
		Type:       "object",
		Properties: map[string]Property{},
		Required:   []string{},
	},
	"clear_command_history": {
		Type:       "object",
		Properties: map[string]Property{},
		Required:   []string{},
	},
	"get_command_task": {
		Type: "object",
		Properties: map[string]Property{
			"task_id": {Type: "string", Description: "任务ID / Task ID"},
		},
		Required: []string{"task_id"},
	},
}

// GetToolSchema 获取工具的输入模式 / Get tool input schema
func GetToolSchema(toolName string) JSONSchema {
	if schema, ok := ToolSchemas[toolName]; ok {
		return schema
	}
	// 返回空对象模式作为默认值 / Return empty object schema as default
	return JSONSchema{
		Type:       "object",
		Properties: map[string]Property{},
		Required:   []string{},
	}
}
