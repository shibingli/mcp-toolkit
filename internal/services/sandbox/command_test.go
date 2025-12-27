package sandbox

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"mcp-toolkit/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteCommand(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	tests := []struct {
		name        string
		request     *types.ExecuteCommandRequest
		expectError bool
		checkOutput bool
	}{
		{
			name: "执行简单命令",
			request: &types.ExecuteCommandRequest{
				Command: getEchoCommand(),
				Args:    []string{"hello"},
				WorkDir: ".",
			},
			expectError: false,
			checkOutput: true,
		},
		{
			name: "黑名单命令应被拒绝",
			request: &types.ExecuteCommandRequest{
				Command: "shutdown",
				Args:    []string{},
				WorkDir: ".",
			},
			expectError: true,
		},
		{
			name: "无效工作目录",
			request: &types.ExecuteCommandRequest{
				Command: getEchoCommand(),
				Args:    []string{"test"},
				WorkDir: "../..",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.ExecuteCommand(tt.request)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
				if tt.checkOutput {
					assert.NotEmpty(t, resp.Stdout)
				}
			}
		})
	}
}

func TestGetCommandBlacklist(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	req := &types.GetCommandBlacklistRequest{}
	resp, err := service.GetCommandBlacklist(req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Commands)
	assert.NotEmpty(t, resp.Directories)
	assert.NotEmpty(t, resp.SystemDirectories)

	// 验证默认黑名单命令存在
	assert.Contains(t, resp.Commands, "shutdown")
	assert.Contains(t, resp.Commands, "format")
	// rm和rmdir已从黑名单移除，改为参数验证
	assert.NotContains(t, resp.Commands, "rm")
	assert.NotContains(t, resp.Commands, "rmdir")
}

func TestUpdateCommandBlacklist(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 添加自定义命令到黑名单
	req := &types.UpdateCommandBlacklistRequest{
		Commands:    []string{"custom_cmd"},
		Directories: []string{"/custom/dir"},
	}

	resp, err := service.UpdateCommandBlacklist(req)
	require.NoError(t, err)
	assert.True(t, resp.Success)

	// 验证黑名单已更新
	getReq := &types.GetCommandBlacklistRequest{}
	getResp, err := service.GetCommandBlacklist(getReq)
	require.NoError(t, err)
	assert.Contains(t, getResp.Commands, "custom_cmd")
	assert.Contains(t, getResp.Directories, "/custom/dir")
}

func TestGetWorkingDirectory(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	req := &types.GetWorkingDirectoryRequest{}
	resp, err := service.GetWorkingDirectory(req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, ".", resp.WorkDir)
}

func TestChangeDirectory(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试目录
	testDir := "test_dir"
	err := os.MkdirAll(filepath.Join(tempDir, testDir), 0755)
	require.NoError(t, err)

	// 切换到测试目录
	req := &types.ChangeDirectoryRequest{
		Path: testDir,
	}
	resp, err := service.ChangeDirectory(req)
	require.NoError(t, err)
	assert.True(t, resp.Success)

	// 验证工作目录已更改
	getReq := &types.GetWorkingDirectoryRequest{}
	getResp, err := service.GetWorkingDirectory(getReq)
	require.NoError(t, err)
	assert.Equal(t, testDir, getResp.WorkDir)
}

func TestIsCommandBlacklisted(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	tests := []struct {
		name        string
		command     string
		blacklisted bool
	}{
		{
			name:        "黑名单命令",
			command:     "shutdown",
			blacklisted: true,
		},
		{
			name:        "黑名单命令(大写)",
			command:     "SHUTDOWN",
			blacklisted: true,
		},
		{
			name:        "带路径的黑名单命令",
			command:     "/usr/bin/shutdown",
			blacklisted: true,
		},
		{
			name:        "非黑名单命令",
			command:     getEchoCommand(),
			blacklisted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.isCommandBlacklisted(tt.command)
			assert.Equal(t, tt.blacklisted, result)
		})
	}
}

func TestIsDirectoryBlacklisted(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	tests := []struct {
		name        string
		directory   string
		blacklisted bool
	}{
		{
			name:        "系统目录",
			directory:   getSystemDirectory(),
			blacklisted: true,
		},
		{
			name:        "沙箱目录",
			directory:   tempDir,
			blacklisted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.isDirectoryBlacklisted(tt.directory)
			assert.Equal(t, tt.blacklisted, result)
		})
	}
}

// getEchoCommand 获取跨平台的echo命令 / Get cross-platform echo command
func getEchoCommand() string {
	if runtime.GOOS == "windows" {
		return "cmd"
	}
	return "echo"
}

// getSystemDirectory 获取系统目录用于测试 / Get system directory for testing
func getSystemDirectory() string {
	switch runtime.GOOS {
	case "windows":
		return "C:\\Windows"
	case "darwin":
		return "/System"
	default:
		return "/bin"
	}
}

// TestValidateCommandPaths 测试命令路径验证 / Test command path validation
func TestValidateCommandPaths(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test file
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	// 获取系统文件路径 / Get system file path
	var systemFilePath string
	if runtime.GOOS == "windows" {
		systemFilePath = "C:\\Windows\\System32\\notepad.exe"
	} else {
		systemFilePath = "/etc/passwd"
	}

	tests := []struct {
		name        string
		command     string
		args        []string
		workDir     string
		expectError bool
	}{
		{
			name:        "rm命令删除沙箱内文件",
			command:     "rm",
			args:        []string{"test.txt"},
			workDir:     tempDir,
			expectError: false,
		},
		{
			name:        "rm命令尝试删除沙箱外文件",
			command:     "rm",
			args:        []string{"../../etc/passwd"},
			workDir:     tempDir,
			expectError: true,
		},
		{
			name:        "rmdir命令删除沙箱内目录",
			command:     "rmdir",
			args:        []string{"subdir"},
			workDir:     tempDir,
			expectError: false,
		},
		{
			name:        "rm命令使用绝对路径删除沙箱外文件",
			command:     "rm",
			args:        []string{systemFilePath},
			workDir:     tempDir,
			expectError: true,
		},
		{
			name:        "非路径敏感命令不需要验证",
			command:     "echo",
			args:        []string{"../../etc/passwd"},
			workDir:     tempDir,
			expectError: false,
		},
		{
			name:        "rm命令带选项参数",
			command:     "rm",
			args:        []string{"-rf", "test.txt"},
			workDir:     tempDir,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateCommandPaths(tt.command, tt.args, tt.workDir)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestExecuteRmCommand 测试rm命令执行 / Test rm command execution
func TestExecuteRmCommand(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建测试文件 / Create test file
	testFile := filepath.Join(tempDir, "to_delete.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	// 验证文件存在 / Verify file exists
	_, err = os.Stat(testFile)
	require.NoError(t, err)

	// 根据操作系统选择删除命令 / Choose delete command based on OS
	var command string
	var args []string
	if runtime.GOOS == "windows" {
		command = "cmd"
		args = []string{"/c", "del", "to_delete.txt"}
	} else {
		command = "rm"
		args = []string{"to_delete.txt"}
	}

	// 执行删除命令 / Execute delete command
	req := &types.ExecuteCommandRequest{
		Command: command,
		Args:    args,
		WorkDir: ".",
	}

	resp, err := service.ExecuteCommand(req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 0, resp.ExitCode)

	// 验证文件已被删除 / Verify file is deleted
	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))
}

// TestExecuteRmCommandOutsideSandbox 测试rm命令沙箱保护 / Test rm command sandbox protection
func TestExecuteRmCommandOutsideSandbox(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 获取系统文件路径 / Get system file path
	var systemFilePath string
	if runtime.GOOS == "windows" {
		systemFilePath = "C:\\Windows\\System32\\notepad.exe"
	} else {
		systemFilePath = "/etc/passwd"
	}

	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "尝试删除沙箱外文件(相对路径)",
			command: "rm",
			args:    []string{"../../etc/passwd"},
		},
		{
			name:    "尝试删除系统文件(绝对路径)",
			command: "rm",
			args:    []string{systemFilePath},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &types.ExecuteCommandRequest{
				Command: tt.command,
				Args:    tt.args,
				WorkDir: ".",
			}

			_, err := service.ExecuteCommand(req)
			assert.Error(t, err)
			if err != nil {
				assert.Contains(t, err.Error(), "outside sandbox")
			}
		})
	}
}
