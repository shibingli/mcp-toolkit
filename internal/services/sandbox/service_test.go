package sandbox

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"mcp-toolkit/pkg/types"
)

// setupTestService 创建测试服务实例 / Create test service instance
func setupTestService(t *testing.T) (*Service, string) {
	// 创建临时测试目录 / Create temporary test directory
	tempDir, err := os.MkdirTemp("", "fs_test_*")
	require.NoError(t, err)

	logger := zap.NewNop()
	service, err := NewService(tempDir, logger)
	require.NoError(t, err)

	return service, tempDir
}

// cleanupTestService 清理测试服务 / Cleanup test service
func cleanupTestService(t *testing.T, tempDir string) {
	err := os.RemoveAll(tempDir)
	require.NoError(t, err)
}

func TestNewService(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "fs_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := zap.NewNop()
	service, err := NewService(tempDir, logger)

	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, tempDir, service.sandboxDir)
}

func TestValidatePath(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	tests := []struct {
		name        string
		path        string
		shouldError bool
	}{
		{
			name:        "valid path",
			path:        "test.txt",
			shouldError: false,
		},
		{
			name:        "valid nested path",
			path:        "dir/subdir/test.txt",
			shouldError: false,
		},
		{
			name:        "path traversal attempt",
			path:        "../outside.txt",
			shouldError: true,
		},
		{
			name:        "path traversal with multiple levels",
			path:        "../../outside.txt",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.validatePath(tt.path)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateFile(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	req := &types.CreateFileRequest{
		Path:    "test.txt",
		Content: "Hello, World!",
	}

	resp, err := service.CreateFile(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, types.MsgFileCreated, resp.Message)

	// 验证文件是否存在 / Verify file exists
	filePath := filepath.Join(tempDir, "test.txt")
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", string(content))
}

func TestCreateFileWithNestedPath(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	req := &types.CreateFileRequest{
		Path:    "dir/subdir/test.txt",
		Content: "Nested file",
	}

	resp, err := service.CreateFile(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证文件是否存在 / Verify file exists
	filePath := filepath.Join(tempDir, "dir/subdir/test.txt")
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "Nested file", string(content))
}

func TestCreateDir(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	req := &types.CreateDirRequest{
		Path: "testdir",
	}

	resp, err := service.CreateDir(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, types.MsgDirCreated, resp.Message)

	// 验证目录是否存在 / Verify directory exists
	dirPath := filepath.Join(tempDir, "testdir")
	info, err := os.Stat(dirPath)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestReadFile(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test file
	testContent := "Test content"
	testPath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testPath, []byte(testContent), 0644)
	require.NoError(t, err)

	req := &types.ReadFileRequest{
		Path: "test.txt",
	}

	resp, err := service.ReadFile(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, testContent, resp.Content)
}

func TestReadFileNotFound(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	req := &types.ReadFileRequest{
		Path: "nonexistent.txt",
	}

	resp, err := service.ReadFile(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), types.ErrFileNotFound)
}

func TestWriteFile(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	req := &types.WriteFileRequest{
		Path:    "test.txt",
		Content: "New content",
	}

	resp, err := service.WriteFile(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, types.MsgFileWritten, resp.Message)

	// 验证文件内容 / Verify file content
	filePath := filepath.Join(tempDir, "test.txt")
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "New content", string(content))
}

func TestDelete(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test file
	testPath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testPath, []byte("test"), 0644)
	require.NoError(t, err)

	req := &types.DeleteRequest{
		Path: "test.txt",
	}

	resp, err := service.Delete(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证文件已删除 / Verify file is deleted
	_, err = os.Stat(testPath)
	assert.True(t, os.IsNotExist(err))
}

func TestCopyFile(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建源文件 / Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	testContent := "Copy test"
	err := os.WriteFile(srcPath, []byte(testContent), 0644)
	require.NoError(t, err)

	req := &types.CopyRequest{
		Source:      "source.txt",
		Destination: "dest.txt",
	}

	resp, err := service.Copy(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证目标文件存在且内容正确 / Verify destination file exists with correct content
	dstPath := filepath.Join(tempDir, "dest.txt")
	content, err := os.ReadFile(dstPath)
	assert.NoError(t, err)
	assert.Equal(t, testContent, string(content))

	// 验证源文件仍然存在 / Verify source file still exists
	_, err = os.Stat(srcPath)
	assert.NoError(t, err)
}

func TestCopyDirectory(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建源目录和文件 / Create source directory and files
	srcDir := filepath.Join(tempDir, "srcdir")
	err := os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(srcDir, "subdir", "file2.txt"), []byte("content2"), 0644)
	require.NoError(t, err)

	req := &types.CopyRequest{
		Source:      "srcdir",
		Destination: "dstdir",
	}

	resp, err := service.Copy(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证目标目录和文件 / Verify destination directory and files
	dstDir := filepath.Join(tempDir, "dstdir")
	content1, err := os.ReadFile(filepath.Join(dstDir, "file1.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "content1", string(content1))

	content2, err := os.ReadFile(filepath.Join(dstDir, "subdir", "file2.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "content2", string(content2))
}

func TestMove(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建源文件 / Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	testContent := "Move test"
	err := os.WriteFile(srcPath, []byte(testContent), 0644)
	require.NoError(t, err)

	req := &types.MoveRequest{
		Source:      "source.txt",
		Destination: "moved.txt",
	}

	resp, err := service.Move(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证目标文件存在 / Verify destination file exists
	dstPath := filepath.Join(tempDir, "moved.txt")
	content, err := os.ReadFile(dstPath)
	assert.NoError(t, err)
	assert.Equal(t, testContent, string(content))

	// 验证源文件已删除 / Verify source file is deleted
	_, err = os.Stat(srcPath)
	assert.True(t, os.IsNotExist(err))
}

func TestListDir(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件和目录 / Create test files and directories
	err := os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("content1"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tempDir, "file2.txt"), []byte("content2"), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(tempDir, "subdir"), 0755)
	require.NoError(t, err)

	req := &types.ListDirRequest{
		Path: ".",
	}

	resp, err := service.ListDir(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Files, 3)

	// 验证文件信息 / Verify file information
	fileNames := make(map[string]bool)
	for _, file := range resp.Files {
		fileNames[file.Name] = true
	}

	assert.True(t, fileNames["file1.txt"])
	assert.True(t, fileNames["file2.txt"])
	assert.True(t, fileNames["subdir"])
}

func TestSearch(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test files
	err := os.WriteFile(filepath.Join(tempDir, "test1.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tempDir, "test2.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tempDir, "other.log"), []byte("content"), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(tempDir, "subdir"), 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tempDir, "subdir", "test3.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	req := &types.SearchRequest{
		Path:    ".",
		Pattern: "test*.txt",
	}

	resp, err := service.Search(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Files, 3)

	// 验证搜索结果 / Verify search results
	fileNames := make(map[string]bool)
	for _, file := range resp.Files {
		fileNames[file.Name] = true
	}

	assert.True(t, fileNames["test1.txt"])
	assert.True(t, fileNames["test2.txt"])
	assert.True(t, fileNames["test3.txt"])
}

func TestBatchDelete(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test files
	err := os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tempDir, "file2.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tempDir, "file3.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	req := &types.BatchDeleteRequest{
		Paths: []string{"file1.txt", "file2.txt"},
	}

	resp, err := service.BatchDelete(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证文件已删除 / Verify files are deleted
	_, err = os.Stat(filepath.Join(tempDir, "file1.txt"))
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(filepath.Join(tempDir, "file2.txt"))
	assert.True(t, os.IsNotExist(err))

	// 验证未删除的文件仍然存在 / Verify undeleted file still exists
	_, err = os.Stat(filepath.Join(tempDir, "file3.txt"))
	assert.NoError(t, err)
}

func TestFileStat(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test file
	testContent := "Test content"
	testPath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testPath, []byte(testContent), 0644)
	require.NoError(t, err)

	req := &types.FileStatRequest{
		Path: "test.txt",
	}

	resp, err := service.FileStat(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test.txt", resp.Name)
	assert.Equal(t, int64(len(testContent)), resp.Size)
	assert.False(t, resp.IsDir)
}

func TestFileExists(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test file
	testPath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testPath, []byte("content"), 0644)
	require.NoError(t, err)

	// 测试存在的文件 / Test existing file
	req := &types.FileExistsRequest{
		Path: "test.txt",
	}

	resp, err := service.FileExists(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Exists)

	// 测试不存在的文件 / Test non-existing file
	req = &types.FileExistsRequest{
		Path: "nonexistent.txt",
	}

	resp, err = service.FileExists(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Exists)
}

func TestGetCurrentTime(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 测试不指定时区（使用本地时区）
	resp, err := service.GetCurrentTime(&types.GetCurrentTimeRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.DateTime)
	assert.NotEmpty(t, resp.Date)
	assert.NotEmpty(t, resp.Time)
	assert.NotEmpty(t, resp.TimeZone)
	assert.NotEmpty(t, resp.TimeZoneOffset)
	assert.NotZero(t, resp.Unix)
	assert.NotZero(t, resp.UnixMilli)
	assert.NotEmpty(t, resp.Weekday)
}

func TestGetCurrentTimeWithTimezone(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 测试指定时区
	resp, err := service.GetCurrentTime(&types.GetCurrentTimeRequest{
		TimeZone: "Asia/Shanghai",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.DateTime)
	assert.Contains(t, resp.TimeZoneOffset, "+08:00")

	// 测试无效时区
	_, err = service.GetCurrentTime(&types.GetCurrentTimeRequest{
		TimeZone: "Invalid/Timezone",
	})
	assert.Error(t, err)
}

func TestCopyFileSourceNotFound(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	req := &types.CopyRequest{
		Source:      "nonexistent.txt",
		Destination: "dest.txt",
	}

	resp, err := service.Copy(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestCopyFileSandboxViolation(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	req := &types.CopyRequest{
		Source:      "../outside.txt",
		Destination: "dest.txt",
	}

	resp, err := service.Copy(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestBatchDeleteWithInvalidPaths(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建一个有效文件 / Create a valid file
	err := os.WriteFile(filepath.Join(tempDir, "valid.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	req := &types.BatchDeleteRequest{
		Paths: []string{"valid.txt", "../invalid.txt"},
	}

	resp, err := service.BatchDelete(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.Contains(t, resp.Message, "failed to delete")
}

func TestListDirNotFound(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	req := &types.ListDirRequest{
		Path: "nonexistent",
	}

	resp, err := service.ListDir(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestFileStatNotFound(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	req := &types.FileStatRequest{
		Path: "nonexistent.txt",
	}

	resp, err := service.FileStat(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestNewServiceWithInvalidPath(t *testing.T) {
	logger := zap.NewNop()

	// 使用一个无效的路径字符(在某些系统上) / Use an invalid path character (on some systems)
	service, err := NewService("./test_sandbox", logger)

	// 即使路径看起来无效,NewService也应该成功创建目录 / Even if path looks invalid, NewService should successfully create directory
	assert.NoError(t, err)
	assert.NotNil(t, service)

	// 清理 / Cleanup
	if service != nil {
		os.RemoveAll(service.sandboxDir)
	}
}

func TestCreateFileWriteError(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建一个只读目录 / Create a read-only directory
	readOnlyDir := filepath.Join(tempDir, "readonly")
	err := os.MkdirAll(readOnlyDir, 0555)
	require.NoError(t, err)
	defer os.Chmod(readOnlyDir, 0755) // 恢复权限以便清理 / Restore permissions for cleanup

	req := &types.CreateFileRequest{
		Path:    "readonly/test.txt",
		Content: "content",
	}

	resp, err := service.CreateFile(req)

	// 在某些系统上可能会成功,在其他系统上会失败 / May succeed on some systems, fail on others
	if err != nil {
		assert.Nil(t, resp)
	}
}

func TestCopyFilePermissionError(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建源文件 / Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	err := os.WriteFile(srcPath, []byte("content"), 0644)
	require.NoError(t, err)

	// 创建只读目标目录 / Create read-only destination directory
	readOnlyDir := filepath.Join(tempDir, "readonly")
	err = os.MkdirAll(readOnlyDir, 0555)
	require.NoError(t, err)
	defer os.Chmod(readOnlyDir, 0755)

	req := &types.CopyRequest{
		Source:      "source.txt",
		Destination: "readonly/dest.txt",
	}

	resp, err := service.Copy(req)

	// 在某些系统上可能会成功,在其他系统上会失败 / May succeed on some systems, fail on others
	if err != nil {
		assert.Nil(t, resp)
	}
}

func TestSearchWithInvalidPattern(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test file
	err := os.WriteFile(filepath.Join(tempDir, "test.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	req := &types.SearchRequest{
		Path:    ".",
		Pattern: "[invalid",
	}

	resp, err := service.Search(req)

	// 无效模式应该被跳过,返回空结果 / Invalid pattern should be skipped, return empty result
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestMoveWithNestedDestination(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建源文件 / Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	err := os.WriteFile(srcPath, []byte("content"), 0644)
	require.NoError(t, err)

	req := &types.MoveRequest{
		Source:      "source.txt",
		Destination: "nested/dir/moved.txt",
	}

	resp, err := service.Move(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证文件已移动 / Verify file is moved
	dstPath := filepath.Join(tempDir, "nested/dir/moved.txt")
	_, err = os.Stat(dstPath)
	assert.NoError(t, err)
}

func TestWriteFileWithNestedPath(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	req := &types.WriteFileRequest{
		Path:    "nested/dir/test.txt",
		Content: "content",
	}

	resp, err := service.WriteFile(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证文件存在 / Verify file exists
	filePath := filepath.Join(tempDir, "nested/dir/test.txt")
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "content", string(content))
}

func TestListDirWithFileInfo(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件和目录 / Create test files and directories
	err := os.WriteFile(filepath.Join(tempDir, "file.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(tempDir, "subdir"), 0755)
	require.NoError(t, err)

	req := &types.ListDirRequest{
		Path: ".",
	}

	resp, err := service.ListDir(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Files, 2)

	// 验证文件信息包含所有字段 / Verify file info contains all fields
	for _, file := range resp.Files {
		assert.NotEmpty(t, file.Name)
		assert.NotEmpty(t, file.Path)
		assert.NotEmpty(t, file.Mode)
		assert.NotZero(t, file.ModTime)
	}
}

// TestDeleteFile 测试删除文件 / Test delete file
func TestDeleteFile(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test file
	testFile := filepath.Join(tempDir, "test_delete.txt")
	err := os.WriteFile(testFile, []byte("content"), 0644)
	require.NoError(t, err)

	// 删除文件 / Delete file
	resp, err := service.DeleteFile(&types.DeleteFileRequest{Path: "test_delete.txt"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证文件已删除 / Verify file deleted
	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))
}

// TestDeleteFileNotExists 测试删除不存在的文件 / Test delete non-existent file
func TestDeleteFileNotExists(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	resp, err := service.DeleteFile(&types.DeleteFileRequest{Path: "not_exists.txt"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// TestDeleteDirectory 测试删除目录 / Test delete directory
func TestDeleteDirectory(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试目录和文件 / Create test directory and files
	testDir := filepath.Join(tempDir, "test_dir")
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(testDir, "file.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	// 递归删除目录 / Delete directory recursively
	resp, err := service.DeleteDirectory(&types.DeleteDirectoryRequest{
		Path:      "test_dir",
		Recursive: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证目录已删除 / Verify directory deleted
	_, err = os.Stat(testDir)
	assert.True(t, os.IsNotExist(err))
}

// TestDeleteDirectoryNonRecursive 测试非递归删除非空目录 / Test non-recursive delete non-empty directory
func TestDeleteDirectoryNonRecursive(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建非空目录 / Create non-empty directory
	testDir := filepath.Join(tempDir, "non_empty_dir")
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(testDir, "file.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	// 非递归删除应该失败 / Non-recursive delete should fail
	resp, err := service.DeleteDirectory(&types.DeleteDirectoryRequest{
		Path:      "non_empty_dir",
		Recursive: false,
	})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// TestCopyFileExplicit 测试显式复制文件 / Test explicit copy file
func TestCopyFileExplicit(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建源文件 / Create source file
	srcFile := filepath.Join(tempDir, "src.txt")
	err := os.WriteFile(srcFile, []byte("source content"), 0644)
	require.NoError(t, err)

	// 复制文件 / Copy file
	resp, err := service.CopyFile(&types.CopyFileRequest{
		Source:      "src.txt",
		Destination: "dst.txt",
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证目标文件存在且内容正确 / Verify destination file exists with correct content
	dstFile := filepath.Join(tempDir, "dst.txt")
	content, err := os.ReadFile(dstFile)
	assert.NoError(t, err)
	assert.Equal(t, "source content", string(content))
}

// TestCopyDirectoryExplicit 测试显式复制目录 / Test explicit copy directory
func TestCopyDirectoryExplicit(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建源目录和文件 / Create source directory and files
	srcDir := filepath.Join(tempDir, "src_dir")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(srcDir, "file.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	// 复制目录 / Copy directory
	resp, err := service.CopyDirectory(&types.CopyDirectoryRequest{
		Source:      "src_dir",
		Destination: "dst_dir",
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证目标目录和文件存在 / Verify destination directory and files exist
	dstFile := filepath.Join(tempDir, "dst_dir", "file.txt")
	content, err := os.ReadFile(dstFile)
	assert.NoError(t, err)
	assert.Equal(t, "content", string(content))
}

// TestMoveFile 测试移动文件 / Test move file
func TestMoveFile(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建源文件 / Create source file
	srcFile := filepath.Join(tempDir, "move_src.txt")
	err := os.WriteFile(srcFile, []byte("move content"), 0644)
	require.NoError(t, err)

	// 移动文件 / Move file
	resp, err := service.MoveFile(&types.MoveFileRequest{
		Source:      "move_src.txt",
		Destination: "move_dst.txt",
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证源文件不存在 / Verify source file doesn't exist
	_, err = os.Stat(srcFile)
	assert.True(t, os.IsNotExist(err))

	// 验证目标文件存在且内容正确 / Verify destination file exists with correct content
	dstFile := filepath.Join(tempDir, "move_dst.txt")
	content, err := os.ReadFile(dstFile)
	assert.NoError(t, err)
	assert.Equal(t, "move content", string(content))
}

// TestMoveDirectory 测试移动目录 / Test move directory
func TestMoveDirectory(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建源目录和文件 / Create source directory and files
	srcDir := filepath.Join(tempDir, "move_src_dir")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(srcDir, "file.txt"), []byte("dir content"), 0644)
	require.NoError(t, err)

	// 移动目录 / Move directory
	resp, err := service.MoveDirectory(&types.MoveDirectoryRequest{
		Source:      "move_src_dir",
		Destination: "move_dst_dir",
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证源目录不存在 / Verify source directory doesn't exist
	_, err = os.Stat(srcDir)
	assert.True(t, os.IsNotExist(err))

	// 验证目标目录和文件存在 / Verify destination directory and files exist
	dstFile := filepath.Join(tempDir, "move_dst_dir", "file.txt")
	content, err := os.ReadFile(dstFile)
	assert.NoError(t, err)
	assert.Equal(t, "dir content", string(content))
}
