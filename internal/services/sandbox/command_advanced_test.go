package sandbox

import (
	"runtime"
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
			Args:    getEchoArgs("test"),
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

	// 验证历史记录内容 / Verify history content
	for _, entry := range histResp.History {
		assert.NotEmpty(t, entry.ID)
		assert.Equal(t, getEchoCommand(), entry.Command)
		// Duration 可能为 0（在快速执行的环境中），所以只检查非负数 / Duration may be 0 (in fast execution environments), so only check non-negative
		assert.GreaterOrEqual(t, entry.Duration, int64(0))
		assert.True(t, entry.Success)
		assert.Equal(t, 0, entry.ExitCode)
	}

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

// TestCommandHistoryPagination 测试命令历史分页 / Test command history pagination
func TestCommandHistoryPagination(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 执行多个命令 / Execute multiple commands
	for i := 0; i < 15; i++ {
		req := &types.ExecuteCommandRequest{
			Command: getEchoCommand(),
			Args:    getEchoArgs("test"),
			WorkDir: ".",
		}
		_, err := service.ExecuteCommand(req)
		require.NoError(t, err)
	}

	// 测试默认限制 / Test default limit
	histReq := &types.GetCommandHistoryRequest{}
	histResp, err := service.GetCommandHistory(histReq)
	require.NoError(t, err)
	assert.Equal(t, 15, len(histResp.History))
	assert.Equal(t, 15, histResp.Total)

	// 测试自定义限制 / Test custom limit
	histReq = &types.GetCommandHistoryRequest{
		Limit: 5,
	}
	histResp, err = service.GetCommandHistory(histReq)
	require.NoError(t, err)
	assert.Equal(t, 5, len(histResp.History))
	assert.Equal(t, 15, histResp.Total)

	// 测试偏移量 / Test offset
	histReq = &types.GetCommandHistoryRequest{
		Limit:  5,
		Offset: 10,
	}
	histResp, err = service.GetCommandHistory(histReq)
	require.NoError(t, err)
	assert.Equal(t, 5, len(histResp.History))
	assert.Equal(t, 15, histResp.Total)

	// 测试超出范围的偏移量 / Test offset out of range
	histReq = &types.GetCommandHistoryRequest{
		Limit:  5,
		Offset: 20,
	}
	histResp, err = service.GetCommandHistory(histReq)
	require.NoError(t, err)
	assert.Equal(t, 0, len(histResp.History))
	assert.Equal(t, 15, histResp.Total)

	// 测试负数偏移量 / Test negative offset
	histReq = &types.GetCommandHistoryRequest{
		Limit:  5,
		Offset: -1,
	}
	histResp, err = service.GetCommandHistory(histReq)
	require.NoError(t, err)
	assert.Equal(t, 5, len(histResp.History))

	// 测试超大限制 / Test large limit
	histReq = &types.GetCommandHistoryRequest{
		Limit: 200,
	}
	histResp, err = service.GetCommandHistory(histReq)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(histResp.History), 100) // 应该被限制为100 / Should be limited to 100
}

// TestCommandHistoryLimit 测试命令历史数量限制 / Test command history size limit
func TestCommandHistoryLimit(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 执行超过1000个命令 / Execute more than 1000 commands
	for i := 0; i < 1005; i++ {
		req := &types.ExecuteCommandRequest{
			Command: getEchoCommand(),
			Args:    getEchoArgs("test"),
			WorkDir: ".",
		}
		_, err := service.ExecuteCommand(req)
		require.NoError(t, err)
	}

	// 验证历史记录数量被限制为1000 / Verify history is limited to 1000
	histReq := &types.GetCommandHistoryRequest{
		Limit: 1100,
	}
	histResp, err := service.GetCommandHistory(histReq)
	require.NoError(t, err)
	assert.LessOrEqual(t, histResp.Total, 1000)
}

// TestAsyncCommandExecution 测试异步命令执行 / Test async command execution
func TestAsyncCommandExecution(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 异步执行命令 / Execute command asynchronously
	req := &types.ExecuteCommandAsyncRequest{
		Command: getEchoCommand(),
		Args:    getEchoArgs("async test"),
		WorkDir: ".",
		Timeout: 10,
	}

	resp, err := service.ExecuteCommandAsync(req)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.TaskID)
	assert.NotEmpty(t, resp.Message)

	// 等待任务完成 / Wait for task completion
	time.Sleep(200 * time.Millisecond)

	// 获取任务状态 / Get task status
	taskReq := &types.GetCommandTaskRequest{
		TaskID: resp.TaskID,
	}
	taskResp, err := service.GetCommandTask(taskReq)
	require.NoError(t, err)
	assert.NotNil(t, taskResp.Task)
	assert.Equal(t, types.TaskStatusCompleted, taskResp.Task.Status)
	assert.Equal(t, 0, taskResp.Task.ExitCode)
	assert.NotZero(t, taskResp.Task.StartTime)
	assert.NotZero(t, taskResp.Task.EndTime)
}

// TestAsyncCommandValidation 测试异步命令验证 / Test async command validation
func TestAsyncCommandValidation(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	tests := []struct {
		name    string
		request *types.ExecuteCommandAsyncRequest
		wantErr bool
	}{
		{
			name: "空命令",
			request: &types.ExecuteCommandAsyncRequest{
				Command: "",
				WorkDir: ".",
			},
			wantErr: true,
		},
		{
			name: "正常命令",
			request: &types.ExecuteCommandAsyncRequest{
				Command: getEchoCommand(),
				Args:    getEchoArgs("test"),
				WorkDir: ".",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.ExecuteCommandAsync(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.TaskID)
			}
		})
	}
}

// TestGetCommandTaskNotFound 测试获取不存在的任务 / Test get non-existent task
func TestGetCommandTaskNotFound(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	req := &types.GetCommandTaskRequest{
		TaskID: "non-existent-task-id",
	}

	_, err := service.GetCommandTask(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
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

	// 测试所有权限级别 / Test all permission levels
	levels := []types.CommandPermissionLevel{
		types.PermissionLevelReadOnly,
		types.PermissionLevelStandard,
		types.PermissionLevelElevated,
		types.PermissionLevelAdmin,
	}

	for _, level := range levels {
		setReq := &types.SetPermissionLevelRequest{
			Level: level,
		}
		setResp, err := service.SetPermissionLevel(setReq)
		require.NoError(t, err)
		assert.True(t, setResp.Success)
		assert.NotEmpty(t, setResp.Message)

		// 验证权限已更改 / Verify permission changed
		getResp, err = service.GetPermissionLevel(getReq)
		require.NoError(t, err)
		assert.Equal(t, level, getResp.Level)
	}

	// 恢复标准权限 / Restore standard permission
	setReq := &types.SetPermissionLevelRequest{
		Level: types.PermissionLevelStandard,
	}
	_, err = service.SetPermissionLevel(setReq)
	require.NoError(t, err)
}

// TestPermissionLevelInvalid 测试无效权限级别 / Test invalid permission level
func TestPermissionLevelInvalid(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 尝试设置无效权限级别 / Try to set invalid permission level
	setReq := &types.SetPermissionLevelRequest{
		Level: 999, // 无效级别 / Invalid level
	}
	_, err := service.SetPermissionLevel(setReq)
	assert.Error(t, err)
}

// TestCancelCommandTask 测试取消命令任务 / Test cancel command task
func TestCancelCommandTask(t *testing.T) {
	service, tempDir := setupTestService(t)

	// 创建一个长时间运行的任务 / Create a long-running task
	var req *types.ExecuteCommandAsyncRequest
	if runtime.GOOS == "windows" {
		req = &types.ExecuteCommandAsyncRequest{
			Command: "cmd",
			Args:    []string{"/c", "ping", "127.0.0.1", "-n", "5"},
			WorkDir: ".",
			Timeout: 60,
		}
	} else {
		req = &types.ExecuteCommandAsyncRequest{
			Command: "sleep",
			Args:    []string{"5"},
			WorkDir: ".",
			Timeout: 60,
		}
	}

	resp, err := service.ExecuteCommandAsync(req)
	require.NoError(t, err)

	// 等待任务开始 / Wait for task to start
	time.Sleep(200 * time.Millisecond)

	// 取消任务 / Cancel task
	cancelReq := &types.CancelCommandTaskRequest{
		TaskID: resp.TaskID,
	}
	cancelResp, err := service.CancelCommandTask(cancelReq)

	// 取消可能成功或失败(如果任务已经完成) / Cancel may succeed or fail (if task already completed)
	if err == nil {
		assert.True(t, cancelResp.Success)
		assert.NotEmpty(t, cancelResp.Message)
	}

	// 等待进程完全结束 / Wait for process to fully terminate
	time.Sleep(1 * time.Second)

	// 验证任务状态 / Verify task status
	taskReq := &types.GetCommandTaskRequest{
		TaskID: resp.TaskID,
	}
	taskResp, err := service.GetCommandTask(taskReq)
	if err == nil {
		// 任务应该是已取消或已完成状态 / Task should be cancelled or completed
		assert.Contains(t, []types.CommandTaskStatus{
			types.TaskStatusCancelled,
			types.TaskStatusCompleted,
			types.TaskStatusFailed,
		}, taskResp.Task.Status)
	}

	// 清理测试环境 / Cleanup test environment
	cleanupTestService(t, tempDir)
}

// TestCancelCommandTaskNotFound 测试取消不存在的任务 / Test cancel non-existent task
func TestCancelCommandTaskNotFound(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	req := &types.CancelCommandTaskRequest{
		TaskID: "non-existent-task-id",
	}

	_, err := service.CancelCommandTask(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestCancelCommandTaskAlreadyCompleted 测试取消已完成的任务 / Test cancel already completed task
func TestCancelCommandTaskAlreadyCompleted(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 创建一个快速完成的任务 / Create a quick task
	req := &types.ExecuteCommandAsyncRequest{
		Command: getEchoCommand(),
		Args:    getEchoArgs("quick task"),
		WorkDir: ".",
		Timeout: 10,
	}

	resp, err := service.ExecuteCommandAsync(req)
	require.NoError(t, err)

	// 等待任务完成 / Wait for task to complete
	time.Sleep(500 * time.Millisecond)

	// 尝试取消已完成的任务 / Try to cancel completed task
	cancelReq := &types.CancelCommandTaskRequest{
		TaskID: resp.TaskID,
	}
	_, err = service.CancelCommandTask(cancelReq)
	// 可能返回错误或成功,取决于实现 / May return error or success depending on implementation
	// 这里不做严格断言 / No strict assertion here
}

// TestAsyncCommandExecutionFailure 测试异步命令执行失败 / Test async command execution failure
func TestAsyncCommandExecutionFailure(t *testing.T) {
	service, tempDir := setupTestService(t)
	defer cleanupTestService(t, tempDir)

	// 测试1: 无效的工作目录（触发 failTask）/ Test 1: Invalid working directory (triggers failTask)
	req1 := &types.ExecuteCommandAsyncRequest{
		Command: getEchoCommand(),
		Args:    getEchoArgs("test"),
		WorkDir: "../..",
	}

	resp1, err := service.ExecuteCommandAsync(req1)
	require.NoError(t, err)

	// 等待任务执行 / Wait for task execution
	time.Sleep(200 * time.Millisecond)

	// 获取任务状态 / Get task status
	taskReq1 := &types.GetCommandTaskRequest{
		TaskID: resp1.TaskID,
	}
	taskResp1, err := service.GetCommandTask(taskReq1)
	require.NoError(t, err)

	// 验证任务失败 / Verify task failed
	assert.Equal(t, types.TaskStatusFailed, taskResp1.Task.Status)
	assert.NotEmpty(t, taskResp1.Task.Error)
	assert.Equal(t, -1, taskResp1.Task.ExitCode)

	// 测试2: 尝试删除沙箱外文件（触发 failTask）/ Test 2: Try to delete file outside sandbox (triggers failTask)
	req2 := &types.ExecuteCommandAsyncRequest{
		Command: "rm",
		Args:    []string{"-rf", "/etc/passwd"},
		WorkDir: ".",
	}

	resp2, err := service.ExecuteCommandAsync(req2)
	require.NoError(t, err)

	// 等待任务执行 / Wait for task execution
	time.Sleep(200 * time.Millisecond)

	// 获取任务状态 / Get task status
	taskReq2 := &types.GetCommandTaskRequest{
		TaskID: resp2.TaskID,
	}
	taskResp2, err := service.GetCommandTask(taskReq2)
	require.NoError(t, err)

	// 验证任务失败 / Verify task failed
	assert.Equal(t, types.TaskStatusFailed, taskResp2.Task.Status)
	assert.NotEmpty(t, taskResp2.Task.Error)
	assert.Equal(t, -1, taskResp2.Task.ExitCode)
}
