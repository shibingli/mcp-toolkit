package sandbox

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"mcp-toolkit/pkg/types"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ExecuteCommandAsync 异步执行命令 / Execute command asynchronously
func (s *Service) ExecuteCommandAsync(req *types.ExecuteCommandAsyncRequest) (*types.ExecuteCommandAsyncResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateExecuteCommandAsyncRequest(req); err != nil {
		return nil, err
	}

	// 权限检查 / Permission check
	if err := s.checkCommandPermission(req.Command, req.PermissionLevel); err != nil {
		return nil, err
	}

	// 创建任务 / Create task
	taskID := uuid.New().String()
	task := &types.CommandTask{
		ID:              taskID,
		Command:         req.Command,
		Args:            req.Args,
		WorkDir:         req.WorkDir,
		Status:          types.TaskStatusPending,
		User:            req.User,
		PermissionLevel: req.PermissionLevel,
		Environment:     req.Environment,
	}

	// 保存任务 / Save task
	s.taskMu.Lock()
	s.commandTasks[taskID] = task
	s.taskMu.Unlock()

	// 异步执行 / Execute asynchronously
	go s.executeTaskAsync(task, req)

	s.logger.Info("async command task created",
		zap.String("task_id", taskID),
		zap.String("command", req.Command),
		zap.Strings("args", req.Args))

	return &types.ExecuteCommandAsyncResponse{
		TaskID:  taskID,
		Message: "command task created successfully",
	}, nil
}

// executeTaskAsync 异步执行任务 / Execute task asynchronously
func (s *Service) executeTaskAsync(task *types.CommandTask, req *types.ExecuteCommandAsyncRequest) {
	// 更新任务状态为运行中 / Update task status to running
	s.updateTaskStatus(task.ID, types.TaskStatusRunning)
	task.StartTime = time.Now()

	// 获取必要的配置信息 / Get necessary configuration
	s.mu.RLock()
	// 验证工作目录 / Validate working directory
	workDir := req.WorkDir
	if workDir == "" {
		workDir = s.currentWorkDir
	}

	validWorkDir, err := s.validatePath(workDir)
	if err != nil {
		s.mu.RUnlock()
		s.failTask(task, err.Error())
		return
	}

	// 检查工作目录是否在黑名单中 / Check if working directory is blacklisted
	if s.isDirectoryBlacklisted(validWorkDir) {
		s.mu.RUnlock()
		s.failTask(task, types.ErrDirectoryBlacklisted)
		return
	}

	// 验证命令参数中的路径 / Validate paths in command arguments
	if err := s.validateCommandPaths(req.Command, req.Args, validWorkDir); err != nil {
		s.mu.RUnlock()
		s.failTask(task, err.Error())
		return
	}
	s.mu.RUnlock() // 释放锁,准备执行命令 / Release lock before executing command

	// 设置超时 / Set timeout
	timeout := time.Duration(req.Timeout) * time.Second
	if req.Timeout == 0 {
		timeout = DefaultCommandTimeout * time.Second
	}
	if timeout > MaxCommandTimeout*time.Second {
		timeout = MaxCommandTimeout * time.Second
	}

	// 创建上下文 / Create context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 构建命令 / Build command
	var cmd *exec.Cmd
	if len(req.Args) > 0 {
		cmd = exec.CommandContext(ctx, req.Command, req.Args...)
	} else {
		cmd = exec.CommandContext(ctx, req.Command)
	}

	// 设置工作目录 / Set working directory
	cmd.Dir = validWorkDir

	// 设置环境变量 / Set environment variables
	if len(req.Environment) > 0 {
		env := make([]string, 0, len(req.Environment))
		for k, v := range req.Environment {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = env
	}

	// 捕获输出 / Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令 / Execute command
	err = cmd.Run()

	task.EndTime = time.Now()
	task.Stdout = stdout.String()
	task.Stderr = stderr.String()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			task.ExitCode = exitErr.ExitCode()
		} else {
			task.ExitCode = -1
		}
		task.Error = err.Error()
		s.updateTaskStatus(task.ID, types.TaskStatusFailed)
	} else {
		task.ExitCode = 0
		s.updateTaskStatus(task.ID, types.TaskStatusCompleted)
	}

	// 添加到历史记录 / Add to history
	entry := createHistoryEntry(
		req.Command, req.Args, workDir,
		task.StartTime, task.EndTime,
		task.ExitCode, task.Status == types.TaskStatusCompleted,
		req.User, req.PermissionLevel, req.Environment,
	)
	s.addCommandHistory(entry)
}

// updateTaskStatus 更新任务状态 / Update task status
func (s *Service) updateTaskStatus(taskID string, status types.CommandTaskStatus) {
	s.taskMu.Lock()
	defer s.taskMu.Unlock()

	if task, exists := s.commandTasks[taskID]; exists {
		task.Status = status
	}
}

// failTask 标记任务失败 / Mark task as failed
func (s *Service) failTask(task *types.CommandTask, errorMsg string) {
	task.EndTime = time.Now()
	task.Error = errorMsg
	task.ExitCode = -1
	s.updateTaskStatus(task.ID, types.TaskStatusFailed)

	s.logger.Error("async command task failed",
		zap.String("task_id", task.ID),
		zap.String("command", task.Command),
		zap.String("error", errorMsg))
}

// GetCommandTask 获取命令任务信息 / Get command task information
func (s *Service) GetCommandTask(req *types.GetCommandTaskRequest) (*types.GetCommandTaskResponse, error) {
	if req.TaskID == "" {
		return nil, errors.New("task_id is required")
	}

	s.taskMu.RLock()
	defer s.taskMu.RUnlock()

	task, exists := s.commandTasks[req.TaskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", req.TaskID)
	}

	return &types.GetCommandTaskResponse{
		Task: task,
	}, nil
}

// CancelCommandTask 取消命令任务 / Cancel command task
func (s *Service) CancelCommandTask(req *types.CancelCommandTaskRequest) (*types.CancelCommandTaskResponse, error) {
	if req.TaskID == "" {
		return nil, errors.New("task_id is required")
	}

	s.taskMu.Lock()
	defer s.taskMu.Unlock()

	task, exists := s.commandTasks[req.TaskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", req.TaskID)
	}

	// 只能取消等待中或运行中的任务 / Can only cancel pending or running tasks
	if task.Status != types.TaskStatusPending && task.Status != types.TaskStatusRunning {
		return nil, fmt.Errorf("task cannot be cancelled, current status: %s", task.Status)
	}

	task.Status = types.TaskStatusCancelled
	task.EndTime = time.Now()

	s.logger.Info("command task cancelled",
		zap.String("task_id", req.TaskID))

	return &types.CancelCommandTaskResponse{
		Success: true,
		Message: "task cancelled successfully",
	}, nil
}

// validateExecuteCommandAsyncRequest 验证异步执行命令请求 / Validate execute command async request
func validateExecuteCommandAsyncRequest(req *types.ExecuteCommandAsyncRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.Command == "" {
		return errors.New(types.ErrInvalidCommand)
	}

	return nil
}
