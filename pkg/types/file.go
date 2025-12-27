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

// Package types 文件操作相关类型定义 / File operation related type definitions
package types

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

// DeleteRequest 删除请求（自动判断文件或目录）/ Delete request (auto-detect file or directory)
type DeleteRequest struct {
	Path string `json:"path"` // 文件或目录路径 / File or directory path
}

// DeleteFileRequest 删除文件请求 / Delete file request
type DeleteFileRequest struct {
	Path string `json:"path"` // 文件路径 / File path
}

// DeleteDirectoryRequest 删除目录请求 / Delete directory request
type DeleteDirectoryRequest struct {
	Path      string `json:"path"`      // 目录路径 / Directory path
	Recursive bool   `json:"recursive"` // 是否递归删除子目录和文件 / Whether to recursively delete subdirectories and files
}

// CopyRequest 复制请求（自动判断文件或目录）/ Copy request (auto-detect file or directory)
type CopyRequest struct {
	Source      string `json:"source"`      // 源路径 / Source path
	Destination string `json:"destination"` // 目标路径 / Destination path
}

// CopyFileRequest 复制文件请求 / Copy file request
type CopyFileRequest struct {
	Source      string `json:"source"`      // 源文件路径 / Source file path
	Destination string `json:"destination"` // 目标文件路径 / Destination file path
}

// CopyDirectoryRequest 复制目录请求 / Copy directory request
type CopyDirectoryRequest struct {
	Source      string `json:"source"`      // 源目录路径 / Source directory path
	Destination string `json:"destination"` // 目标目录路径 / Destination directory path
}

// MoveRequest 移动请求（自动判断文件或目录）/ Move request (auto-detect file or directory)
type MoveRequest struct {
	Source      string `json:"source"`      // 源路径 / Source path
	Destination string `json:"destination"` // 目标路径 / Destination path
}

// MoveFileRequest 移动文件请求 / Move file request
type MoveFileRequest struct {
	Source      string `json:"source"`      // 源文件路径 / Source file path
	Destination string `json:"destination"` // 目标文件路径 / Destination file path
}

// MoveDirectoryRequest 移动目录请求 / Move directory request
type MoveDirectoryRequest struct {
	Source      string `json:"source"`      // 源目录路径 / Source directory path
	Destination string `json:"destination"` // 目标目录路径 / Destination directory path
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

// DeleteFileResponse 删除文件响应 / Delete file response
type DeleteFileResponse = OperationResponse

// DeleteDirectoryResponse 删除目录响应 / Delete directory response
type DeleteDirectoryResponse = OperationResponse

// CopyResponse 复制响应 / Copy response
type CopyResponse = OperationResponse

// CopyFileResponse 复制文件响应 / Copy file response
type CopyFileResponse = OperationResponse

// CopyDirectoryResponse 复制目录响应 / Copy directory response
type CopyDirectoryResponse = OperationResponse

// MoveResponse 移动响应 / Move response
type MoveResponse = OperationResponse

// MoveFileResponse 移动文件响应 / Move file response
type MoveFileResponse = OperationResponse

// MoveDirectoryResponse 移动目录响应 / Move directory response
type MoveDirectoryResponse = OperationResponse

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
