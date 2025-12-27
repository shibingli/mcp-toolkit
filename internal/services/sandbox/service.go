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

package sandbox

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"mcp-toolkit/pkg/types"

	"go.uber.org/zap"
)

// Service 文件系统服务 / Filesystem service
type Service struct {
	sandboxDir         string                        // 沙箱目录 / Sandbox directory
	logger             *zap.Logger                   // 日志记录器 / Logger
	currentWorkDir     string                        // 当前工作目录(相对于沙箱根目录) / Current working directory (relative to sandbox root)
	blacklistCommands  []string                      // 黑名单命令列表 / Blacklist commands
	blacklistDirs      []string                      // 黑名单目录列表 / Blacklist directories
	mu                 sync.RWMutex                  // 读写锁,保护黑名单和工作目录 / RWMutex to protect blacklist and working directory
	commandHistory     []*types.CommandHistoryEntry  // 命令执行历史 / Command execution history
	commandTasks       map[string]*types.CommandTask // 异步命令任务 / Async command tasks
	taskMu             sync.RWMutex                  // 任务锁 / Task mutex
	permissionLevel    types.CommandPermissionLevel  // 当前权限级别 / Current permission level
	defaultEnvironment map[string]string             // 默认环境变量 / Default environment variables
	auditLogger        *zap.Logger                   // 审计日志记录器 / Audit logger
}

// NewService 创建文件系统服务实例 / Create filesystem service instance
func NewService(sandboxDir string, logger *zap.Logger) (*Service, error) {
	// 参数验证 / Parameter validation
	if sandboxDir == "" {
		return nil, errors.New("sandbox directory cannot be empty")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	// 确保沙箱目录存在 / Ensure sandbox directory exists
	absPath, err := filepath.Abs(sandboxDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// 创建沙箱目录(如果不存在) / Create sandbox directory if not exists
	if err = os.MkdirAll(absPath, DefaultDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create sandbox directory: %w", err)
	}

	// 初始化黑名单 / Initialize blacklist
	blacklistCommands := make([]string, len(DefaultBlacklistCommands))
	copy(blacklistCommands, DefaultBlacklistCommands)

	blacklistDirs := make([]string, len(DefaultBlacklistDirectories))
	copy(blacklistDirs, DefaultBlacklistDirectories)

	// 创建审计日志记录器 / Create audit logger
	auditLogger := logger.Named("audit")

	return &Service{
		sandboxDir:         absPath,
		logger:             logger,
		currentWorkDir:     ".", // 默认工作目录为沙箱根目录 / Default working directory is sandbox root
		blacklistCommands:  blacklistCommands,
		blacklistDirs:      blacklistDirs,
		commandHistory:     make([]*types.CommandHistoryEntry, 0, 100),
		commandTasks:       make(map[string]*types.CommandTask),
		permissionLevel:    types.PermissionLevelStandard, // 默认标准权限 / Default standard permission
		defaultEnvironment: make(map[string]string),
		auditLogger:        auditLogger,
	}, nil
}

// validatePath 验证路径是否在沙箱目录内 / Validate if path is within sandbox directory
func (s *Service) validatePath(path string) (string, error) {
	path = strings.TrimSpace(path)

	// 参数验证 / Parameter validation
	if path == "" {
		return "", errors.New(types.ErrInvalidPath)
	}

	// 清理路径 / Clean path
	cleanedPath := filepath.Clean(path)

	// 检查路径遍历攻击 / Check for path traversal attacks
	if strings.Contains(cleanedPath, "..") {
		s.logger.Warn("path traversal attempt detected",
			zap.String("requested_path", path))
		return "", errors.New(types.ErrSandboxViolation)
	}

	// 获取绝对路径 / Get absolute path
	absPath := filepath.Join(s.sandboxDir, cleanedPath)
	finalPath := filepath.Clean(absPath)

	// 确保路径以沙箱目录开头 / Ensure path starts with sandbox directory
	// 使用filepath.Rel来确保路径在沙箱内 / Use filepath.Rel to ensure path is within sandbox
	relPath, err := filepath.Rel(s.sandboxDir, finalPath)
	if err != nil {
		s.logger.Warn("failed to get relative path",
			zap.String("requested_path", path),
			zap.Error(err))
		return "", errors.New(types.ErrSandboxViolation)
	}

	// 检查相对路径是否包含.. / Check if relative path contains ..
	if strings.HasPrefix(relPath, "..") || strings.Contains(relPath, string(filepath.Separator)+"..") {
		s.logger.Warn("sandbox violation attempt",
			zap.String("requested_path", path),
			zap.String("relative_path", relPath),
			zap.String("sandbox_dir", s.sandboxDir))
		return "", errors.New(types.ErrSandboxViolation)
	}

	return finalPath, nil
}

// CreateFile 创建文件 / Create file
func (s *Service) CreateFile(req *types.CreateFileRequest) (*types.OperationResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateCreateFileRequest(req); err != nil {
		return nil, err
	}

	validPath, err := s.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	// 确保父目录存在 / Ensure parent directory exists
	dir := filepath.Dir(validPath)
	if err = os.MkdirAll(dir, DefaultDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create parent directory: %w", err)
	}

	// 创建文件 / Create file
	file, err := os.Create(validPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 写入内容 / Write content
	if req.Content != "" {
		if _, err = file.WriteString(req.Content); err != nil {
			return nil, fmt.Errorf("failed to write content: %w", err)
		}
	}

	s.logger.Info("file created", zap.String("path", validPath))
	return &types.OperationResponse{
		Success: true,
		Message: types.MsgFileCreated,
	}, nil
}

// CreateDir 创建目录 / Create directory
func (s *Service) CreateDir(req *types.CreateDirRequest) (*types.OperationResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateCreateDirRequest(req); err != nil {
		return nil, err
	}

	validPath, err := s.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	if err = os.MkdirAll(validPath, DefaultDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	s.logger.Info("directory created", zap.String("path", validPath))
	return &types.OperationResponse{
		Success: true,
		Message: types.MsgDirCreated,
	}, nil
}

// ReadFile 读取文件 / Read file
func (s *Service) ReadFile(req *types.ReadFileRequest) (*types.ReadFileResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateReadFileRequest(req); err != nil {
		return nil, err
	}

	validPath, err := s.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(validPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.New(types.ErrFileNotFound)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return &types.ReadFileResponse{
		Content: string(content),
	}, nil
}

// WriteFile 写入文件 / Write file
func (s *Service) WriteFile(req *types.WriteFileRequest) (*types.OperationResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateWriteFileRequest(req); err != nil {
		return nil, err
	}

	validPath, err := s.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	// 确保父目录存在 / Ensure parent directory exists
	dir := filepath.Dir(validPath)
	if err = os.MkdirAll(dir, DefaultDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create parent directory: %w", err)
	}

	if err = os.WriteFile(validPath, []byte(req.Content), DefaultFilePerm); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	s.logger.Info("file written", zap.String("path", validPath))
	return &types.OperationResponse{
		Success: true,
		Message: types.MsgFileWritten,
	}, nil
}

// Delete 删除文件或目录 / Delete file or directory
func (s *Service) Delete(req *types.DeleteRequest) (*types.OperationResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateDeleteRequest(req); err != nil {
		return nil, err
	}

	validPath, err := s.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	if err = os.RemoveAll(validPath); err != nil {
		return nil, fmt.Errorf("failed to delete: %w", err)
	}

	s.logger.Info("deleted", zap.String("path", validPath))
	return &types.OperationResponse{
		Success: true,
		Message: types.MsgFileDeleted,
	}, nil
}

// DeleteFile 删除文件（仅限文件，不能删除目录）/ Delete file only (not directory)
func (s *Service) DeleteFile(req *types.DeleteFileRequest) (*types.OperationResponse, error) {
	// 参数验证 / Parameter validation
	if req.Path == "" {
		return nil, errors.New(types.ErrPathRequired)
	}

	validPath, err := s.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	// 检查是否为文件 / Check if it's a file
	info, err := os.Stat(validPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.New(types.ErrFileNotFound)
		}
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	if info.IsDir() {
		return nil, errors.New("path is a directory, use delete_directory instead")
	}

	if err = os.Remove(validPath); err != nil {
		return nil, fmt.Errorf("failed to delete file: %w", err)
	}

	s.logger.Info("file deleted", zap.String("path", validPath))
	return &types.OperationResponse{
		Success: true,
		Message: types.MsgFileDeleted,
	}, nil
}

// DeleteDirectory 删除目录（仅限目录，不能删除文件）/ Delete directory only (not file)
func (s *Service) DeleteDirectory(req *types.DeleteDirectoryRequest) (*types.OperationResponse, error) {
	// 参数验证 / Parameter validation
	if req.Path == "" {
		return nil, errors.New(types.ErrPathRequired)
	}

	validPath, err := s.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	// 检查是否为目录 / Check if it's a directory
	info, err := os.Stat(validPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.New(types.ErrFileNotFound)
		}
		return nil, fmt.Errorf("failed to stat directory: %w", err)
	}

	if !info.IsDir() {
		return nil, errors.New("path is a file, use delete_file instead")
	}

	// 默认递归删除 / Default to recursive delete
	if req.Recursive {
		if err = os.RemoveAll(validPath); err != nil {
			return nil, fmt.Errorf("failed to delete directory: %w", err)
		}
	} else {
		// 非递归删除，只能删除空目录 / Non-recursive delete, can only delete empty directory
		if err = os.Remove(validPath); err != nil {
			return nil, fmt.Errorf("failed to delete directory (directory may not be empty, use recursive=true): %w", err)
		}
	}

	s.logger.Info("directory deleted", zap.String("path", validPath), zap.Bool("recursive", req.Recursive))
	return &types.OperationResponse{
		Success: true,
		Message: "directory deleted successfully",
	}, nil
}

// Copy 复制文件或目录 / Copy file or directory
func (s *Service) Copy(req *types.CopyRequest) (*types.OperationResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateCopyRequest(req); err != nil {
		return nil, err
	}

	srcPath, err := s.validatePath(req.Source)
	if err != nil {
		return nil, err
	}

	dstPath, err := s.validatePath(req.Destination)
	if err != nil {
		return nil, err
	}

	// 获取源文件信息 / Get source file info
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.New(types.ErrFileNotFound)
		}
		return nil, fmt.Errorf("failed to stat source: %w", err)
	}

	if srcInfo.IsDir() {
		// 复制目录 / Copy directory
		if err = s.copyDir(srcPath, dstPath); err != nil {
			return nil, err
		}
	} else {
		// 复制文件 / Copy file
		if err = s.copyFile(srcPath, dstPath); err != nil {
			return nil, err
		}
	}

	s.logger.Info("copied", zap.String("source", srcPath), zap.String("destination", dstPath))
	return &types.OperationResponse{
		Success: true,
		Message: types.MsgFileCopied,
	}, nil
}

// CopyFile 复制文件（仅限文件）/ Copy file only
func (s *Service) CopyFile(req *types.CopyFileRequest) (*types.OperationResponse, error) {
	// 参数验证 / Parameter validation
	if req.Source == "" || req.Destination == "" {
		return nil, errors.New("source and destination are required")
	}

	srcPath, err := s.validatePath(req.Source)
	if err != nil {
		return nil, err
	}

	dstPath, err := s.validatePath(req.Destination)
	if err != nil {
		return nil, err
	}

	// 检查源是否为文件 / Check if source is a file
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.New(types.ErrFileNotFound)
		}
		return nil, fmt.Errorf("failed to stat source: %w", err)
	}

	if srcInfo.IsDir() {
		return nil, errors.New("source is a directory, use copy_directory instead")
	}

	if err = s.copyFile(srcPath, dstPath); err != nil {
		return nil, err
	}

	s.logger.Info("file copied", zap.String("source", srcPath), zap.String("destination", dstPath))
	return &types.OperationResponse{
		Success: true,
		Message: types.MsgFileCopied,
	}, nil
}

// CopyDirectory 复制目录（仅限目录）/ Copy directory only
func (s *Service) CopyDirectory(req *types.CopyDirectoryRequest) (*types.OperationResponse, error) {
	// 参数验证 / Parameter validation
	if req.Source == "" || req.Destination == "" {
		return nil, errors.New("source and destination are required")
	}

	srcPath, err := s.validatePath(req.Source)
	if err != nil {
		return nil, err
	}

	dstPath, err := s.validatePath(req.Destination)
	if err != nil {
		return nil, err
	}

	// 检查源是否为目录 / Check if source is a directory
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.New(types.ErrFileNotFound)
		}
		return nil, fmt.Errorf("failed to stat source: %w", err)
	}

	if !srcInfo.IsDir() {
		return nil, errors.New("source is a file, use copy_file instead")
	}

	if err = s.copyDir(srcPath, dstPath); err != nil {
		return nil, err
	}

	s.logger.Info("directory copied", zap.String("source", srcPath), zap.String("destination", dstPath))
	return &types.OperationResponse{
		Success: true,
		Message: "directory copied successfully",
	}, nil
}

// copyFile 复制单个文件 / Copy single file
func (s *Service) copyFile(src, dst string) error {
	// 确保目标目录存在 / Ensure destination directory exists
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, DefaultDirPerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() { _ = dstFile.Close() }()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// 复制文件权限 / Copy file permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	if err = os.Chmod(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}

// copyDir 递归复制目录 / Recursively copy directory
func (s *Service) copyDir(src, dst string) error {
	// 创建目标目录 / Create destination directory
	if err := os.MkdirAll(dst, DefaultDirPerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err = s.copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err = s.copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// Move 移动文件或目录 / Move file or directory
func (s *Service) Move(req *types.MoveRequest) (*types.OperationResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateMoveRequest(req); err != nil {
		return nil, err
	}

	srcPath, err := s.validatePath(req.Source)
	if err != nil {
		return nil, err
	}

	dstPath, err := s.validatePath(req.Destination)
	if err != nil {
		return nil, err
	}

	// 确保目标目录存在 / Ensure destination directory exists
	dstDir := filepath.Dir(dstPath)
	if err = os.MkdirAll(dstDir, DefaultDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	if err = os.Rename(srcPath, dstPath); err != nil {
		return nil, fmt.Errorf("failed to move: %w", err)
	}

	s.logger.Info("moved", zap.String("source", srcPath), zap.String("destination", dstPath))
	return &types.OperationResponse{
		Success: true,
		Message: types.MsgFileMoved,
	}, nil
}

// MoveFile 移动文件（仅限文件）/ Move file only
func (s *Service) MoveFile(req *types.MoveFileRequest) (*types.OperationResponse, error) {
	// 参数验证 / Parameter validation
	if req.Source == "" || req.Destination == "" {
		return nil, errors.New("source and destination are required")
	}

	srcPath, err := s.validatePath(req.Source)
	if err != nil {
		return nil, err
	}

	dstPath, err := s.validatePath(req.Destination)
	if err != nil {
		return nil, err
	}

	// 检查源是否为文件 / Check if source is a file
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.New(types.ErrFileNotFound)
		}
		return nil, fmt.Errorf("failed to stat source: %w", err)
	}

	if srcInfo.IsDir() {
		return nil, errors.New("source is a directory, use move_directory instead")
	}

	// 确保目标目录存在 / Ensure destination directory exists
	dstDir := filepath.Dir(dstPath)
	if err = os.MkdirAll(dstDir, DefaultDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	if err = os.Rename(srcPath, dstPath); err != nil {
		return nil, fmt.Errorf("failed to move file: %w", err)
	}

	s.logger.Info("file moved", zap.String("source", srcPath), zap.String("destination", dstPath))
	return &types.OperationResponse{
		Success: true,
		Message: types.MsgFileMoved,
	}, nil
}

// MoveDirectory 移动目录（仅限目录）/ Move directory only
func (s *Service) MoveDirectory(req *types.MoveDirectoryRequest) (*types.OperationResponse, error) {
	// 参数验证 / Parameter validation
	if req.Source == "" || req.Destination == "" {
		return nil, errors.New("source and destination are required")
	}

	srcPath, err := s.validatePath(req.Source)
	if err != nil {
		return nil, err
	}

	dstPath, err := s.validatePath(req.Destination)
	if err != nil {
		return nil, err
	}

	// 检查源是否为目录 / Check if source is a directory
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.New(types.ErrFileNotFound)
		}
		return nil, fmt.Errorf("failed to stat source: %w", err)
	}

	if !srcInfo.IsDir() {
		return nil, errors.New("source is a file, use move_file instead")
	}

	// 确保目标父目录存在 / Ensure destination parent directory exists
	dstDir := filepath.Dir(dstPath)
	if err = os.MkdirAll(dstDir, DefaultDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	if err = os.Rename(srcPath, dstPath); err != nil {
		return nil, fmt.Errorf("failed to move directory: %w", err)
	}

	s.logger.Info("directory moved", zap.String("source", srcPath), zap.String("destination", dstPath))
	return &types.OperationResponse{
		Success: true,
		Message: "directory moved successfully",
	}, nil
}

// ListDir 列出目录内容 / List directory contents
func (s *Service) ListDir(req *types.ListDirRequest) (*types.ListDirResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateListDirRequest(req); err != nil {
		return nil, err
	}

	validPath, err := s.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(validPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.New(types.ErrFileNotFound)
		}
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	fileInfos := make([]types.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			s.logger.Warn("failed to get file info", zap.String("name", entry.Name()), zap.Error(err))
			continue
		}

		fileInfos = append(fileInfos, types.FileInfo{
			Name:    entry.Name(),
			Path:    filepath.Join(req.Path, entry.Name()),
			Size:    info.Size(),
			IsDir:   entry.IsDir(),
			Mode:    info.Mode().String(),
			ModTime: info.ModTime(),
		})
	}

	return &types.ListDirResponse{
		Files: fileInfos,
	}, nil
}

// Search 搜索文件 / Search files
func (s *Service) Search(req *types.SearchRequest) (*types.SearchResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateSearchRequest(req); err != nil {
		return nil, err
	}

	validPath, err := s.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	var matchedFiles []types.FileInfo

	err = filepath.Walk(validPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过错误 / Skip errors
		}

		// 检查文件名是否匹配模式 / Check if filename matches pattern
		matched, err := filepath.Match(req.Pattern, info.Name())
		if err != nil {
			return nil // 跳过无效模式 / Skip invalid patterns
		}

		if matched {
			// 计算相对路径 / Calculate relative path
			relPath, err := filepath.Rel(s.sandboxDir, path)
			if err != nil {
				relPath = path
			}

			matchedFiles = append(matchedFiles, types.FileInfo{
				Name:    info.Name(),
				Path:    relPath,
				Size:    info.Size(),
				IsDir:   info.IsDir(),
				Mode:    info.Mode().String(),
				ModTime: info.ModTime(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	return &types.SearchResponse{
		Files: matchedFiles,
	}, nil
}

// BatchDelete 批量删除 / Batch delete
func (s *Service) BatchDelete(req *types.BatchDeleteRequest) (*types.OperationResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateBatchDeleteRequest(req); err != nil {
		return nil, err
	}

	var failedPaths []string

	for _, path := range req.Paths {
		validPath, err := s.validatePath(path)
		if err != nil {
			failedPaths = append(failedPaths, path)
			s.logger.Warn("failed to validate path", zap.String("path", path), zap.Error(err))
			continue
		}

		if err = os.RemoveAll(validPath); err != nil {
			failedPaths = append(failedPaths, path)
			s.logger.Warn("failed to delete", zap.String("path", validPath), zap.Error(err))
		}
	}

	if len(failedPaths) > 0 {
		return &types.OperationResponse{
			Success: false,
			Message: fmt.Sprintf("failed to delete %d paths: %v", len(failedPaths), failedPaths),
		}, nil
	}

	s.logger.Info("batch delete completed", zap.Int("count", len(req.Paths)))
	return &types.OperationResponse{
		Success: true,
		Message: types.MsgSuccess,
	}, nil
}

// FileStat 获取文件状态 / Get file status
func (s *Service) FileStat(req *types.FileStatRequest) (*types.FileInfo, error) {
	// 参数验证 / Parameter validation
	if err := validateFileStatRequest(req); err != nil {
		return nil, err
	}

	validPath, err := s.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(validPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.New(types.ErrFileNotFound)
		}
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// 计算相对路径 / Calculate relative path
	relPath, err := filepath.Rel(s.sandboxDir, validPath)
	if err != nil {
		relPath = req.Path
	}

	return &types.FileInfo{
		Name:    info.Name(),
		Path:    relPath,
		Size:    info.Size(),
		IsDir:   info.IsDir(),
		Mode:    info.Mode().String(),
		ModTime: info.ModTime(),
	}, nil
}

// FileExists 检查文件是否存在 / Check if file exists
func (s *Service) FileExists(req *types.FileExistsRequest) (*types.FileExistsResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateFileExistsRequest(req); err != nil {
		return nil, err
	}

	validPath, err := s.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(validPath)
	exists := !errors.Is(err, fs.ErrNotExist)

	return &types.FileExistsResponse{
		Exists: exists,
	}, nil
}

// GetCurrentTime 获取当前系统时间 / Get current system time
// 支持指定时区，如果未指定则使用系统本地时区
func (s *Service) GetCurrentTime(req *types.GetCurrentTimeRequest) (*types.GetTimeResponse, error) {
	var loc *time.Location
	var err error

	// 如果指定了时区，则加载该时区；否则使用本地时区
	if req != nil && req.TimeZone != "" {
		loc, err = time.LoadLocation(req.TimeZone)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone '%s': %w", req.TimeZone, err)
		}
	} else {
		loc = time.Local
	}

	now := time.Now().In(loc)
	zoneName, offset := now.Zone()

	// 计算时区偏移字符串 (如 +08:00)
	offsetHours := offset / 3600
	offsetMinutes := (offset % 3600) / 60
	if offsetMinutes < 0 {
		offsetMinutes = -offsetMinutes
	}
	var offsetStr string
	if offset >= 0 {
		offsetStr = fmt.Sprintf("+%02d:%02d", offsetHours, offsetMinutes)
	} else {
		offsetStr = fmt.Sprintf("-%02d:%02d", -offsetHours, offsetMinutes)
	}

	// 判断是否为夏令时（通过比较标准时区名称）
	_, winterOffset := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc).Zone()
	isDST := offset != winterOffset

	return &types.GetTimeResponse{
		DateTime:       now.Format("2006-01-02 15:04:05"),
		Date:           now.Format("2006-01-02"),
		Time:           now.Format("15:04:05"),
		TimeZone:       zoneName,
		TimeZoneOffset: offsetStr,
		Unix:           now.Unix(),
		UnixMilli:      now.UnixMilli(),
		Weekday:        now.Weekday().String(),
		IsDST:          isDST,
	}, nil
}
