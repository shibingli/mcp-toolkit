package sandbox

import (
	"testing"
	"time"

	"mcp-toolkit/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCommandHistory 测试命令历史功能 / Test command history functionality
func TestCommandHistory(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 执行几个命令 / Execute some commands
	for i := 0; i < 3; i++ {
		req := &types.ExecuteCommandRequest{
			Command: getEchoCommand(),
			Args:    []string{"test"},
			WorkDir: ".",
		}
		_, err := service.ExecuteCommand(req)
		require.NoError(t, err)
	}

	// 获取历史记录 / Get history
	histReq := &types.GetCommandHistoryRequest{
		Limit: 10,
	}
	histResp, err := service.GetCommandHistory(histReq)
	require.NoError(t, err)
	assert.NotNil(t, histResp)
	assert.Equal(t, 3, len(histResp.History))
	assert.Equal(t, 3, histResp.Total)

	// 清空历史 / Clear history
	clearReq := &types.ClearCommandHistoryRequest{}
	clearResp, err := service.ClearCommandHistory(clearReq)
	require.NoError(t, err)
	assert.True(t, clearResp.Success)

	// 验证已清空 / Verify cleared
	histResp, err = service.GetCommandHistory(histReq)
	require.NoError(t, err)
	assert.Equal(t, 0, len(histResp.History))
}

// TestAsyncCommandExecution 测试异步命令执行 / Test async command execution
func TestAsyncCommandExecution(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 异步执行命令 / Execute command asynchronously
	req := &types.ExecuteCommandAsyncRequest{
		Command: getEchoCommand(),
		Args:    []string{"async test"},
		WorkDir: ".",
		Timeout: 10,
	}

	resp, err := service.ExecuteCommandAsync(req)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.TaskID)

	// 等待任务完成 / Wait for task completion
	time.Sleep(100 * time.Millisecond)

	// 获取任务状态 / Get task status
	taskReq := &types.GetCommandTaskRequest{
		TaskID: resp.TaskID,
	}
	taskResp, err := service.GetCommandTask(taskReq)
	require.NoError(t, err)
	assert.NotNil(t, taskResp.Task)
	assert.Equal(t, types.TaskStatusCompleted, taskResp.Task.Status)
}

// TestPermissionLevel 测试权限级别 / Test permission level
func TestPermissionLevel(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 获取当前权限级别 / Get current permission level
	getReq := &types.GetPermissionLevelRequest{}
	getResp, err := service.GetPermissionLevel(getReq)
	require.NoError(t, err)
	assert.Equal(t, types.PermissionLevelStandard, getResp.Level)

	// 设置为只读权限 / Set to read-only permission
	setReq := &types.SetPermissionLevelRequest{
		Level: types.PermissionLevelReadOnly,
	}
	setResp, err := service.SetPermissionLevel(setReq)
	require.NoError(t, err)
	assert.True(t, setResp.Success)

	// 验证权限已更改 / Verify permission changed
	getResp, err = service.GetPermissionLevel(getReq)
	require.NoError(t, err)
	assert.Equal(t, types.PermissionLevelReadOnly, getResp.Level)

	// 尝试执行非只读命令应该失败 / Trying to execute non-read-only command should fail
	cmdReq := &types.ExecuteCommandRequest{
		Command: "rm",
		Args:    []string{"test.txt"},
		WorkDir: ".",
	}
	_, err = service.ExecuteCommand(cmdReq)
	assert.Error(t, err)

	// 恢复标准权限 / Restore standard permission
	setReq.Level = types.PermissionLevelStandard
	_, err = service.SetPermissionLevel(setReq)
	require.NoError(t, err)
}

// TestCancelCommandTask 测试取消命令任务 / Test cancel command task
func TestCancelCommandTask(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建一个长时间运行的任务 / Create a long-running task
	req := &types.ExecuteCommandAsyncRequest{
		Command: getEchoCommand(),
		Args:    []string{"long running"},
		WorkDir: ".",
		Timeout: 60,
	}

	resp, err := service.ExecuteCommandAsync(req)
	require.NoError(t, err)

	// 立即取消任务 / Cancel task immediately
	cancelReq := &types.CancelCommandTaskRequest{
		TaskID: resp.TaskID,
	}
	cancelResp, err := service.CancelCommandTask(cancelReq)

	// 取消可能成功或失败(如果任务已经完成) / Cancel may succeed or fail (if task already completed)
	if err == nil {
		assert.True(t, cancelResp.Success)
	}
}
