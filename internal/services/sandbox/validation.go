package sandbox

import (
	"errors"
	"fmt"

	"mcp-toolkit/pkg/types"
)

// validateCreateFileRequest 验证创建文件请求 / Validate create file request
func validateCreateFileRequest(req *types.CreateFileRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Path == "" {
		return errors.New(types.ErrInvalidPath)
	}
	if len(req.Path) > MaxPathLength {
		return fmt.Errorf("path length exceeds maximum allowed length of %d", MaxPathLength)
	}
	if len(req.Content) > int(MaxFileSize) {
		return fmt.Errorf("content size exceeds maximum allowed size of %d bytes", MaxFileSize)
	}
	return nil
}

// validateCreateDirRequest 验证创建目录请求 / Validate create directory request
func validateCreateDirRequest(req *types.CreateDirRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Path == "" {
		return errors.New(types.ErrInvalidPath)
	}
	return nil
}

// validateReadFileRequest 验证读取文件请求 / Validate read file request
func validateReadFileRequest(req *types.ReadFileRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Path == "" {
		return errors.New(types.ErrInvalidPath)
	}
	return nil
}

// validateWriteFileRequest 验证写入文件请求 / Validate write file request
func validateWriteFileRequest(req *types.WriteFileRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Path == "" {
		return errors.New(types.ErrInvalidPath)
	}
	if len(req.Path) > MaxPathLength {
		return fmt.Errorf("path length exceeds maximum allowed length of %d", MaxPathLength)
	}
	if len(req.Content) > int(MaxFileSize) {
		return fmt.Errorf("content size exceeds maximum allowed size of %d bytes", MaxFileSize)
	}
	return nil
}

// validateDeleteRequest 验证删除请求 / Validate delete request
func validateDeleteRequest(req *types.DeleteRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Path == "" {
		return errors.New(types.ErrInvalidPath)
	}
	return nil
}

// validateCopyRequest 验证复制请求 / Validate copy request
func validateCopyRequest(req *types.CopyRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Source == "" {
		return errors.New("source path cannot be empty")
	}
	if req.Destination == "" {
		return errors.New("destination path cannot be empty")
	}
	if req.Source == req.Destination {
		return errors.New("source and destination cannot be the same")
	}
	return nil
}

// validateMoveRequest 验证移动请求 / Validate move request
func validateMoveRequest(req *types.MoveRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Source == "" {
		return errors.New("source path cannot be empty")
	}
	if req.Destination == "" {
		return errors.New("destination path cannot be empty")
	}
	if req.Source == req.Destination {
		return errors.New("source and destination cannot be the same")
	}
	return nil
}

// validateListDirRequest 验证列出目录请求 / Validate list directory request
func validateListDirRequest(req *types.ListDirRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Path == "" {
		return errors.New(types.ErrInvalidPath)
	}
	return nil
}

// validateSearchRequest 验证搜索请求 / Validate search request
func validateSearchRequest(req *types.SearchRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Path == "" {
		return errors.New(types.ErrInvalidPath)
	}
	if req.Pattern == "" {
		return errors.New("search pattern cannot be empty")
	}
	return nil
}

// validateBatchDeleteRequest 验证批量删除请求 / Validate batch delete request
func validateBatchDeleteRequest(req *types.BatchDeleteRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if len(req.Paths) == 0 {
		return errors.New("paths list cannot be empty")
	}
	if len(req.Paths) > MaxBatchDeleteCount {
		return fmt.Errorf("batch delete count exceeds maximum allowed count of %d", MaxBatchDeleteCount)
	}
	for i, path := range req.Paths {
		if path == "" {
			return fmt.Errorf("path at index %d cannot be empty", i)
		}
		if len(path) > MaxPathLength {
			return fmt.Errorf("path at index %d exceeds maximum allowed length of %d", i, MaxPathLength)
		}
	}
	return nil
}

// validateFileStatRequest 验证文件状态请求 / Validate file stat request
func validateFileStatRequest(req *types.FileStatRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Path == "" {
		return errors.New(types.ErrInvalidPath)
	}
	return nil
}

// validateFileExistsRequest 验证文件存在请求 / Validate file exists request
func validateFileExistsRequest(req *types.FileExistsRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Path == "" {
		return errors.New(types.ErrInvalidPath)
	}
	return nil
}

// validateExecuteCommandRequest 验证执行命令请求 / Validate execute command request
func validateExecuteCommandRequest(req *types.ExecuteCommandRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Command == "" {
		return errors.New(types.ErrInvalidCommand)
	}
	if req.Timeout < 0 {
		return errors.New("timeout cannot be negative")
	}
	if req.Timeout > MaxCommandTimeout {
		return fmt.Errorf("timeout exceeds maximum allowed timeout of %d seconds", MaxCommandTimeout)
	}
	return nil
}

// validateUpdateCommandBlacklistRequest 验证更新命令黑名单请求 / Validate update command blacklist request
func validateUpdateCommandBlacklistRequest(req *types.UpdateCommandBlacklistRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if len(req.Commands) == 0 && len(req.Directories) == 0 {
		return errors.New("at least one command or directory must be provided")
	}
	return nil
}

// validateChangeDirectoryRequest 验证切换目录请求 / Validate change directory request
func validateChangeDirectoryRequest(req *types.ChangeDirectoryRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Path == "" {
		return errors.New(types.ErrInvalidPath)
	}
	return nil
}
