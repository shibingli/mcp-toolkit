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

import "time"

// FileInfo 文件信息结构体 / File information structure
type FileInfo struct {
	Name    string    `json:"name"`     // 文件名 / File name
	Path    string    `json:"path"`     // 文件路径 / File path
	Size    int64     `json:"size"`     // 文件大小(字节) / File size in bytes
	IsDir   bool      `json:"is_dir"`   // 是否为目录 / Whether it is a directory
	Mode    string    `json:"mode"`     // 文件权限 / File permissions
	ModTime time.Time `json:"mod_time"` // 修改时间 / Modification time
}

// CreateFileRequest 创建文件请求 / Create file request
type CreateFileRequest struct {
	Path    string `json:"path"`    // 文件路径 / File path
	Content string `json:"content"` // 文件内容 / File content
}

// CreateDirRequest 创建目录请求 / Create directory request
type CreateDirRequest struct {
	Path string `json:"path"` // 目录路径 / Directory path
}

// ReadFileRequest 读取文件请求 / Read file request
type ReadFileRequest struct {
	Path string `json:"path"` // 文件路径 / File path
}

// WriteFileRequest 写入文件请求 / Write file request
type WriteFileRequest struct {
	Path    string `json:"path"`    // 文件路径 / File path
	Content string `json:"content"` // 文件内容 / File content
}

// DeleteRequest 删除请求 / Delete request
type DeleteRequest struct {
	Path string `json:"path"` // 文件或目录路径 / File or directory path
}

// CopyRequest 复制请求 / Copy request
type CopyRequest struct {
	Source      string `json:"source"`      // 源路径 / Source path
	Destination string `json:"destination"` // 目标路径 / Destination path
}

// MoveRequest 移动请求 / Move request
type MoveRequest struct {
	Source      string `json:"source"`      // 源路径 / Source path
	Destination string `json:"destination"` // 目标路径 / Destination path
}

// ListDirRequest 列出目录请求 / List directory request
type ListDirRequest struct {
	Path string `json:"path"` // 目录路径 / Directory path
}

// SearchRequest 搜索请求 / Search request
type SearchRequest struct {
	Path    string `json:"path"`    // 搜索路径 / Search path
	Pattern string `json:"pattern"` // 搜索模式 / Search pattern
}

// BatchDeleteRequest 批量删除请求 / Batch delete request
type BatchDeleteRequest struct {
	Paths []string `json:"paths"` // 文件或目录路径列表 / List of file or directory paths
}

// FileStatRequest 获取文件状态请求 / Get file status request
type FileStatRequest struct {
	Path string `json:"path"` // 文件路径 / File path
}

// FileExistsRequest 检查文件是否存在请求 / Check file exists request
type FileExistsRequest struct {
	Path string `json:"path"` // 文件路径 / File path
}

// CreateFileResponse 创建文件响应 / Create file response
type CreateFileResponse = OperationResponse

// CreateDirResponse 创建目录响应 / Create directory response
type CreateDirResponse = OperationResponse

// WriteFileResponse 写入文件响应 / Write file response
type WriteFileResponse = OperationResponse

// DeleteResponse 删除响应 / Delete response
type DeleteResponse = OperationResponse

// CopyResponse 复制响应 / Copy response
type CopyResponse = OperationResponse

// MoveResponse 移动响应 / Move response
type MoveResponse = OperationResponse

// BatchDeleteResponse 批量删除响应 / Batch delete response
type BatchDeleteResponse = OperationResponse

// FileStatResponse 文件状态响应 / File stat response
type FileStatResponse = FileInfo

// FileExistsResponse 检查文件是否存在响应 / Check file exists response
type FileExistsResponse struct {
	Exists bool `json:"exists"` // 是否存在 / Whether exists
}

// ReadFileResponse 读取文件响应 / Read file response
type ReadFileResponse struct {
	Content string `json:"content"` // 文件内容 / File content
}

// ListDirResponse 列出目录响应 / List directory response
type ListDirResponse struct {
	Files []FileInfo `json:"files"` // 文件列表 / File list
}

// SearchResponse 搜索响应 / Search response
type SearchResponse struct {
	Files []FileInfo `json:"files"` // 匹配的文件列表 / Matched file list
}

// OperationResponse 操作响应 / Operation response
type OperationResponse struct {
	Success bool   `json:"success"` // 是否成功 / Whether successful
	Message string `json:"message"` // 消息 / Message
}

// GetCurrentTimeResponse 获取当前时间响应 / Get current time response
type GetCurrentTimeResponse = GetTimeResponse

// GetTimeResponse 获取时间响应 / Get time response
type GetTimeResponse struct {
	Time     time.Time `json:"time"`      // 当前时间 / Current time
	TimeZone string    `json:"time_zone"` // 时区 / Time zone
	Unix     int64     `json:"unix"`      // Unix时间戳 / Unix timestamp
}

// ExecuteCommandRequest 执行命令请求 / Execute command request
type ExecuteCommandRequest struct {
	Command string   `json:"command"`           // 要执行的命令 / Command to execute
	Args    []string `json:"args,omitempty"`    // 命令参数 / Command arguments
	WorkDir string   `json:"work_dir"`          // 工作目录(相对于沙箱根目录) / Working directory (relative to sandbox root)
	Timeout int      `json:"timeout,omitempty"` // 超时时间(秒),0表示不限制 / Timeout in seconds, 0 means no limit
}

// ExecuteCommandResponse 执行命令响应 / Execute command response
type ExecuteCommandResponse struct {
	Success  bool   `json:"success"`   // 是否成功 / Whether successful
	ExitCode int    `json:"exit_code"` // 退出码 / Exit code
	Stdout   string `json:"stdout"`    // 标准输出 / Standard output
	Stderr   string `json:"stderr"`    // 标准错误 / Standard error
	Message  string `json:"message"`   // 消息 / Message
}

// GetCommandBlacklistRequest 获取命令黑名单请求 / Get command blacklist request
type GetCommandBlacklistRequest struct{}

// GetCommandBlacklistResponse 获取命令黑名单响应 / Get command blacklist response
type GetCommandBlacklistResponse struct {
	Commands          []string `json:"commands"`           // 黑名单命令列表 / Blacklist command list
	Directories       []string `json:"directories"`        // 黑名单目录列表 / Blacklist directory list
	SystemDirectories []string `json:"system_directories"` // 系统目录列表 / System directory list
}

// UpdateCommandBlacklistRequest 更新命令黑名单请求 / Update command blacklist request
type UpdateCommandBlacklistRequest struct {
	Commands    []string `json:"commands,omitempty"`    // 要添加的黑名单命令 / Commands to add to blacklist
	Directories []string `json:"directories,omitempty"` // 要添加的黑名单目录 / Directories to add to blacklist
}

// UpdateCommandBlacklistResponse 更新命令黑名单响应 / Update command blacklist response
type UpdateCommandBlacklistResponse = OperationResponse

// GetWorkingDirectoryRequest 获取当前工作目录请求 / Get working directory request
type GetWorkingDirectoryRequest struct{}

// GetWorkingDirectoryResponse 获取当前工作目录响应 / Get working directory response
type GetWorkingDirectoryResponse struct {
	WorkDir string `json:"work_dir"` // 当前工作目录(相对于沙箱根目录) / Current working directory (relative to sandbox root)
}

// ChangeDirectoryRequest 切换工作目录请求 / Change directory request
type ChangeDirectoryRequest struct {
	Path string `json:"path"` // 目标目录路径(相对于沙箱根目录) / Target directory path (relative to sandbox root)
}

// ChangeDirectoryResponse 切换工作目录响应 / Change directory response
type ChangeDirectoryResponse = OperationResponse

// CommandPermissionLevel 命令执行权限级别 / Command execution permission level
type CommandPermissionLevel int

const (
	// PermissionLevelReadOnly 只读权限 - 只能执行查询类命令 / Read-only permission - can only execute query commands
	PermissionLevelReadOnly CommandPermissionLevel = iota
	// PermissionLevelStandard 标准权限 - 可以执行大部分命令 / Standard permission - can execute most commands
	PermissionLevelStandard
	// PermissionLevelElevated 提升权限 - 可以执行所有非黑名单命令 / Elevated permission - can execute all non-blacklisted commands
	PermissionLevelElevated
	// PermissionLevelAdmin 管理员权限 - 可以执行所有命令(包括修改黑名单) / Admin permission - can execute all commands (including modifying blacklist)
	PermissionLevelAdmin
)

// CommandHistoryEntry 命令执行历史记录条目 / Command execution history entry
type CommandHistoryEntry struct {
	ID              string                 `json:"id"`               // 历史记录ID / History entry ID
	Command         string                 `json:"command"`          // 执行的命令 / Executed command
	Args            []string               `json:"args"`             // 命令参数 / Command arguments
	WorkDir         string                 `json:"work_dir"`         // 工作目录 / Working directory
	StartTime       time.Time              `json:"start_time"`       // 开始时间 / Start time
	EndTime         time.Time              `json:"end_time"`         // 结束时间 / End time
	Duration        int64                  `json:"duration"`         // 执行时长(毫秒) / Execution duration in milliseconds
	ExitCode        int                    `json:"exit_code"`        // 退出码 / Exit code
	Success         bool                   `json:"success"`          // 是否成功 / Whether successful
	User            string                 `json:"user,omitempty"`   // 执行用户 / Executing user
	PermissionLevel CommandPermissionLevel `json:"permission_level"` // 权限级别 / Permission level
	Environment     map[string]string      `json:"environment"`      // 环境变量 / Environment variables
}

// ExecuteCommandAsyncRequest 异步执行命令请求 / Execute command async request
type ExecuteCommandAsyncRequest struct {
	Command         string                 `json:"command"`                    // 要执行的命令 / Command to execute
	Args            []string               `json:"args,omitempty"`             // 命令参数 / Command arguments
	WorkDir         string                 `json:"work_dir"`                   // 工作目录 / Working directory
	Timeout         int                    `json:"timeout,omitempty"`          // 超时时间(秒) / Timeout in seconds
	Environment     map[string]string      `json:"environment,omitempty"`      // 环境变量 / Environment variables
	PermissionLevel CommandPermissionLevel `json:"permission_level,omitempty"` // 权限级别 / Permission level
	User            string                 `json:"user,omitempty"`             // 执行用户 / Executing user
}

// ExecuteCommandAsyncResponse 异步执行命令响应 / Execute command async response
type ExecuteCommandAsyncResponse struct {
	TaskID  string `json:"task_id"` // 任务ID / Task ID
	Message string `json:"message"` // 消息 / Message
}

// CommandTaskStatus 命令任务状态 / Command task status
type CommandTaskStatus string

const (
	// TaskStatusPending 等待执行 / Pending execution
	TaskStatusPending CommandTaskStatus = "pending"
	// TaskStatusRunning 正在执行 / Running
	TaskStatusRunning CommandTaskStatus = "running"
	// TaskStatusCompleted 执行完成 / Completed
	TaskStatusCompleted CommandTaskStatus = "completed"
	// TaskStatusFailed 执行失败 / Failed
	TaskStatusFailed CommandTaskStatus = "failed"
	// TaskStatusCancelled 已取消 / Cancelled
	TaskStatusCancelled CommandTaskStatus = "cancelled"
)

// CommandTask 命令执行任务 / Command execution task
type CommandTask struct {
	ID              string                 `json:"id"`               // 任务ID / Task ID
	Command         string                 `json:"command"`          // 命令 / Command
	Args            []string               `json:"args"`             // 参数 / Arguments
	WorkDir         string                 `json:"work_dir"`         // 工作目录 / Working directory
	Status          CommandTaskStatus      `json:"status"`           // 状态 / Status
	StartTime       time.Time              `json:"start_time"`       // 开始时间 / Start time
	EndTime         time.Time              `json:"end_time"`         // 结束时间 / End time
	ExitCode        int                    `json:"exit_code"`        // 退出码 / Exit code
	Stdout          string                 `json:"stdout"`           // 标准输出 / Standard output
	Stderr          string                 `json:"stderr"`           // 标准错误 / Standard error
	Error           string                 `json:"error,omitempty"`  // 错误信息 / Error message
	User            string                 `json:"user,omitempty"`   // 执行用户 / Executing user
	PermissionLevel CommandPermissionLevel `json:"permission_level"` // 权限级别 / Permission level
	Environment     map[string]string      `json:"environment"`      // 环境变量 / Environment variables
}

// GetCommandTaskRequest 获取命令任务请求 / Get command task request
type GetCommandTaskRequest struct {
	TaskID string `json:"task_id"` // 任务ID / Task ID
}

// GetCommandTaskResponse 获取命令任务响应 / Get command task response
type GetCommandTaskResponse struct {
	Task *CommandTask `json:"task"` // 任务信息 / Task information
}

// CancelCommandTaskRequest 取消命令任务请求 / Cancel command task request
type CancelCommandTaskRequest struct {
	TaskID string `json:"task_id"` // 任务ID / Task ID
}

// CancelCommandTaskResponse 取消命令任务响应 / Cancel command task response
type CancelCommandTaskResponse = OperationResponse

// GetCommandHistoryRequest 获取命令历史请求 / Get command history request
type GetCommandHistoryRequest struct {
	Limit  int    `json:"limit,omitempty"`  // 返回记录数量限制 / Limit of returned records
	Offset int    `json:"offset,omitempty"` // 偏移量 / Offset
	User   string `json:"user,omitempty"`   // 按用户过滤 / Filter by user
}

// GetCommandHistoryResponse 获取命令历史响应 / Get command history response
type GetCommandHistoryResponse struct {
	History []*CommandHistoryEntry `json:"history"` // 历史记录列表 / History entry list
	Total   int                    `json:"total"`   // 总记录数 / Total count
}

// ClearCommandHistoryRequest 清空命令历史请求 / Clear command history request
type ClearCommandHistoryRequest struct{}

// ClearCommandHistoryResponse 清空命令历史响应 / Clear command history response
type ClearCommandHistoryResponse = OperationResponse

// SetPermissionLevelRequest 设置权限级别请求 / Set permission level request
type SetPermissionLevelRequest struct {
	Level CommandPermissionLevel `json:"level"` // 权限级别 / Permission level
}

// SetPermissionLevelResponse 设置权限级别响应 / Set permission level response
type SetPermissionLevelResponse = OperationResponse

// GetPermissionLevelRequest 获取权限级别请求 / Get permission level request
type GetPermissionLevelRequest struct{}

// GetPermissionLevelResponse 获取权限级别响应 / Get permission level response
type GetPermissionLevelResponse struct {
	Level CommandPermissionLevel `json:"level"` // 当前权限级别 / Current permission level
}
