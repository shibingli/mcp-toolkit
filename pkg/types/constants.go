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
	"os"
	"path/filepath"
	"runtime"
)

const (
	// ProtocolVersion MCP协议版本 / MCP protocol version
	ProtocolVersion = "2024-11-05"

	// ServerName MCP服务器名称 / MCP server name
	ServerName = "mcp-toolkit"

	// ServerVersion MCP服务器版本 / MCP server version
	ServerVersion = "1.0.1"

	// DefaultSandboxDir 默认沙箱目录（相对路径） / Default sandbox directory (relative path)
	DefaultSandboxDir = "./sandbox"
)

// GetDefaultSandboxDir 根据操作系统返回默认的沙箱目录路径 / Get default sandbox directory path based on OS
// Linux/MacOS: /tmp/mcp_sandbox_toolkit
// Windows: %TEMP%\mcp_sandbox_toolkit 或 C:\Temp\mcp_sandbox_toolkit
func GetDefaultSandboxDir() string {
	switch runtime.GOOS {
	case "windows":
		// Windows系统：优先使用TEMP环境变量，否则使用C:\Temp\mcp_sandbox_toolkit
		if tempDir := os.Getenv("TEMP"); tempDir != "" {
			return filepath.Join(tempDir, "mcp_sandbox_toolkit")
		}
		if tempDir := os.Getenv("TMP"); tempDir != "" {
			return filepath.Join(tempDir, "mcp_sandbox_toolkit")
		}
		return `C:\Temp\mcp_sandbox_toolkit`
	case "linux", "darwin":
		// Linux/MacOS系统：使用/tmp/mcp_sandbox_toolkit
		return "/tmp/mcp_sandbox_toolkit"
	default:
		// 其他系统：使用当前目录下的sandbox
		return DefaultSandboxDir
	}
}

const (
	// ErrSandboxViolation 沙箱违规错误 / Sandbox violation error
	ErrSandboxViolation = "path is outside sandbox directory"

	// ErrFileNotFound 文件未找到错误 / File not found error
	ErrFileNotFound = "file not found"

	// ErrInvalidPath 无效路径错误 / Invalid path error
	ErrInvalidPath = "invalid path"

	// ErrPermissionDenied 权限拒绝错误 / Permission denied error
	ErrPermissionDenied = "permission denied"

	// ErrIsDirectory 是目录错误 / Is directory error
	ErrIsDirectory = "is a directory"

	// ErrNotDirectory 不是目录错误 / Not a directory error
	ErrNotDirectory = "not a directory"

	// ErrDirectoryNotEmpty 目录非空错误 / Directory not empty error
	ErrDirectoryNotEmpty = "directory not empty"

	// ErrPathNotFound 路径未找到错误 / Path not found error
	ErrPathNotFound = "path not found"

	// ErrCommandBlacklisted 命令在黑名单中错误 / Command blacklisted error
	ErrCommandBlacklisted = "command is blacklisted"

	// ErrDirectoryBlacklisted 目录在黑名单中错误 / Directory blacklisted error
	ErrDirectoryBlacklisted = "directory is blacklisted"

	// ErrInvalidCommand 无效命令错误 / Invalid command error
	ErrInvalidCommand = "invalid command"

	// ErrPathRequired 路径必填错误 / Path required error
	ErrPathRequired = "path is required"

	// MsgSuccess 成功消息 / Success message
	MsgSuccess = "operation completed successfully"

	// MsgFileCreated 文件创建成功消息 / File created successfully message
	MsgFileCreated = "file created successfully"

	// MsgDirCreated 目录创建成功消息 / Directory created successfully message
	MsgDirCreated = "directory created successfully"

	// MsgFileDeleted 文件删除成功消息 / File deleted successfully message
	MsgFileDeleted = "file deleted successfully"

	// MsgFileCopied 文件复制成功消息 / File copied successfully message
	MsgFileCopied = "file copied successfully"

	// MsgFileMoved 文件移动成功消息 / File moved successfully message
	MsgFileMoved = "file moved successfully"

	// MsgFileWritten 文件写入成功消息 / File written successfully message
	MsgFileWritten = "file written successfully"

	// MsgCommandExecuted 命令执行成功消息 / Command executed successfully message
	MsgCommandExecuted = "command executed successfully"

	// MsgBlacklistUpdated 黑名单更新成功消息 / Blacklist updated successfully message
	MsgBlacklistUpdated = "blacklist updated successfully"

	// MsgDirectoryChanged 目录切换成功消息 / Directory changed successfully message
	MsgDirectoryChanged = "directory changed successfully"
)
