package sandbox

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"mcp-toolkit/pkg/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterTools(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建MCP服务器 / Create MCP server
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    "Test Server",
			Version: "1.0.0",
		},
		&mcp.ServerOptions{
			Capabilities: &mcp.ServerCapabilities{
				Tools: &mcp.ToolCapabilities{},
			},
		},
	)

	// 注册工具 / Register tools
	service.RegisterTools(mcpServer)

	// 验证工具已注册 / Verify tools are registered
	assert.NotNil(t, mcpServer)
}

func TestHandleCreateFile(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	args := types.CreateFileRequest{
		Path:    "test.txt",
		Content: "Hello, World!",
	}

	result, resp, err := service.handleCreateFile(context.Background(), nil, args)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证文件已创建 / Verify file is created
	filePath := filepath.Join(tempDir, "test.txt")
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", string(content))
}

func TestHandleCreateDir(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	args := types.CreateDirRequest{
		Path: "testdir",
	}

	result, resp, err := service.handleCreateDir(context.Background(), nil, args)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证目录已创建 / Verify directory is created
	dirPath := filepath.Join(tempDir, "testdir")
	info, err := os.Stat(dirPath)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestHandleReadFile(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test file
	testContent := "Test content"
	testPath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testPath, []byte(testContent), 0644)
	require.NoError(t, err)

	args := types.ReadFileRequest{
		Path: "test.txt",
	}

	result, resp, err := service.handleReadFile(context.Background(), nil, args)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, testContent, resp.Content)
}

func TestHandleWriteFile(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	args := types.WriteFileRequest{
		Path:    "test.txt",
		Content: "New content",
	}

	result, resp, err := service.handleWriteFile(context.Background(), nil, args)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证文件内容 / Verify file content
	filePath := filepath.Join(tempDir, "test.txt")
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "New content", string(content))
}

func TestHandleDelete(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test file
	testPath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testPath, []byte("test"), 0644)
	require.NoError(t, err)

	args := types.DeleteRequest{
		Path: "test.txt",
	}

	result, resp, err := service.handleDelete(context.Background(), nil, args)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证文件已删除 / Verify file is deleted
	_, err = os.Stat(testPath)
	assert.True(t, os.IsNotExist(err))
}

func TestHandleCopy(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建源文件 / Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	err := os.WriteFile(srcPath, []byte("content"), 0644)
	require.NoError(t, err)

	args := types.CopyRequest{
		Source:      "source.txt",
		Destination: "dest.txt",
	}

	result, resp, err := service.handleCopy(context.Background(), nil, args)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证目标文件已创建 / Verify destination file is created
	dstPath := filepath.Join(tempDir, "dest.txt")
	content, err := os.ReadFile(dstPath)
	assert.NoError(t, err)
	assert.Equal(t, "content", string(content))
}

func TestHandleMove(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建源文件 / Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	err := os.WriteFile(srcPath, []byte("content"), 0644)
	require.NoError(t, err)

	args := types.MoveRequest{
		Source:      "source.txt",
		Destination: "moved.txt",
	}

	result, resp, err := service.handleMove(context.Background(), nil, args)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证文件已移动 / Verify file is moved
	dstPath := filepath.Join(tempDir, "moved.txt")
	_, err = os.Stat(dstPath)
	assert.NoError(t, err)

	_, err = os.Stat(srcPath)
	assert.True(t, os.IsNotExist(err))
}

func TestHandleListDir(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test files
	err := os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	args := types.ListDirRequest{
		Path: ".",
	}

	result, resp, err := service.handleListDir(context.Background(), nil, args)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Files, 1)
}

func TestHandleSearch(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test files
	err := os.WriteFile(filepath.Join(tempDir, "test1.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	args := types.SearchRequest{
		Path:    ".",
		Pattern: "test*.txt",
	}

	result, resp, err := service.handleSearch(context.Background(), nil, args)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Files, 1)
}

func TestHandleBatchDelete(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test files
	err := os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("content"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tempDir, "file2.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	args := types.BatchDeleteRequest{
		Paths: []string{"file1.txt", "file2.txt"},
	}

	result, resp, err := service.handleBatchDelete(context.Background(), nil, args)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
}

func TestHandleFileStat(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test file
	testPath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testPath, []byte("content"), 0644)
	require.NoError(t, err)

	args := types.FileStatRequest{
		Path: "test.txt",
	}

	result, resp, err := service.handleFileStat(context.Background(), nil, args)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, "test.txt", resp.Name)
}

func TestHandleFileExists(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test file
	testPath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testPath, []byte("content"), 0644)
	require.NoError(t, err)

	args := types.FileExistsRequest{
		Path: "test.txt",
	}

	result, resp, err := service.handleFileExists(context.Background(), nil, args)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.True(t, resp.Exists)
}

func TestHandleGetCurrentTime(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	args := types.GetCurrentTimeRequest{}

	result, resp, err := service.handleGetCurrentTime(context.Background(), nil, args)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.DateTime)
	assert.NotEmpty(t, resp.TimeZone)
	assert.NotZero(t, resp.Unix)
}
