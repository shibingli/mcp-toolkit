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

// Package types 通用类型定义 / Common type definitions
// 本文件包含所有模块共用的基础类型
// 具体功能类型已拆分到以下文件：
//   - file.go: 文件操作相关类型
//   - time.go: 时间相关类型
//   - command.go: 命令执行相关类型
//   - sysinfo.go: 系统信息相关类型
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

// OperationResponse 操作响应 / Operation response
// 通用的操作结果响应结构体，被多个模块复用
type OperationResponse struct {
	Success bool   `json:"success"` // 是否成功 / Whether successful
	Message string `json:"message"` // 消息 / Message
}
