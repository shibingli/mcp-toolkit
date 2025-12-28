package sandbox

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"mcp-toolkit/pkg/types"

	"go.uber.org/zap"
)

// ExecuteCommand 执行命令 / Execute command
func (s *Service) ExecuteCommand(req *types.ExecuteCommandRequest) (*types.ExecuteCommandResponse, error) {
	// 记录开始时间 / Record start time
	startTime := time.Now()

	// 参数验证 / Parameter validation
	if err := validateExecuteCommandRequest(req); err != nil {
		return nil, err
	}

	// 权限检查 / Permission check
	if err := s.checkCommandPermission(req.Command, 0); err != nil {
		s.logger.Warn("command permission denied",
			zap.String("command", req.Command),
			zap.Error(err))
		return nil, err
	}

	// 获取必要的配置信息 / Get necessary configuration
	s.mu.RLock()
	// 检查命令是否在黑名单中 / Check if command is in blacklist
	if s.isCommandBlacklisted(req.Command) {
		s.mu.RUnlock()
		s.logger.Warn("blocked blacklisted command",
			zap.String("command", req.Command))
		return nil, errors.New(types.ErrCommandBlacklisted)
	}

	// 验证工作目录 / Validate working directory
	workDir := req.WorkDir
	if workDir == "" {
		workDir = s.currentWorkDir
	}

	validWorkDir, err := s.validatePath(workDir)
	if err != nil {
		s.mu.RUnlock()
		return nil, err
	}

	// 检查工作目录是否在黑名单中 / Check if working directory is in blacklist
	if s.isDirectoryBlacklisted(validWorkDir) {
		s.mu.RUnlock()
		s.logger.Warn("blocked blacklisted directory",
			zap.String("directory", validWorkDir))
		return nil, errors.New(types.ErrDirectoryBlacklisted)
	}

	// 验证命令参数中的路径 / Validate paths in command arguments
	if err := s.validateCommandPaths(req.Command, req.Args, validWorkDir); err != nil {
		s.mu.RUnlock()
		s.logger.Warn("command contains invalid paths",
			zap.String("command", req.Command),
			zap.Strings("args", req.Args),
			zap.Error(err))
		return nil, err
	}

	permLevel := s.permissionLevel
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

	// 捕获输出 / Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令 / Execute command
	s.logger.Info("executing command",
		zap.String("command", req.Command),
		zap.Strings("args", req.Args),
		zap.String("work_dir", validWorkDir))

	err = cmd.Run()

	// 获取退出码 / Get exit code
	exitCode := 0
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
	}

	// 记录结束时间 / Record end time
	endTime := time.Now()

	success := exitCode == 0
	message := types.MsgCommandExecuted
	if !success {
		message = fmt.Sprintf("命令执行完成,退出码: %d / Command completed with exit code: %d", exitCode, exitCode)
	}

	s.logger.Info("command executed",
		zap.String("command", req.Command),
		zap.Int("exit_code", exitCode),
		zap.Bool("success", success))

	// 添加到历史记录 / Add to history
	entry := createHistoryEntry(
		req.Command, req.Args, workDir,
		startTime, endTime,
		exitCode, success,
		"", // user - 可以从context中获取 / can be obtained from context
		permLevel,
		nil, // environment
	)
	s.addCommandHistory(entry)

	return &types.ExecuteCommandResponse{
		Success:  success,
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Message:  message,
	}, nil
}

// GetCommandBlacklist 获取命令黑名单 / Get command blacklist
func (s *Service) GetCommandBlacklist(_ *types.GetCommandBlacklistRequest) (*types.GetCommandBlacklistResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	commands := make([]string, len(s.blacklistCommands))
	copy(commands, s.blacklistCommands)

	directories := make([]string, len(s.blacklistDirs))
	copy(directories, s.blacklistDirs)

	systemDirs := make([]string, len(SystemDirectories))
	copy(systemDirs, SystemDirectories)

	return &types.GetCommandBlacklistResponse{
		Commands:          commands,
		Directories:       directories,
		SystemDirectories: systemDirs,
	}, nil
}

// UpdateCommandBlacklist 更新命令黑名单 / Update command blacklist
func (s *Service) UpdateCommandBlacklist(req *types.UpdateCommandBlacklistRequest) (*types.UpdateCommandBlacklistResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateUpdateCommandBlacklistRequest(req); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 添加命令到黑名单 / Add commands to blacklist
	if len(req.Commands) > 0 {
		for _, cmd := range req.Commands {
			cmd = strings.TrimSpace(cmd)
			if cmd == "" {
				continue
			}
			// 检查是否已存在 / Check if already exists
			exists := false
			for _, existing := range s.blacklistCommands {
				if existing == cmd {
					exists = true
					break
				}
			}
			if !exists {
				s.blacklistCommands = append(s.blacklistCommands, cmd)
			}
		}
	}

	// 添加目录到黑名单 / Add directories to blacklist
	if len(req.Directories) > 0 {
		for _, dir := range req.Directories {
			dir = strings.TrimSpace(dir)
			if dir == "" {
				continue
			}
			// 检查是否已存在 / Check if already exists
			exists := false
			for _, existing := range s.blacklistDirs {
				if existing == dir {
					exists = true
					break
				}
			}
			if !exists {
				s.blacklistDirs = append(s.blacklistDirs, dir)
			}
		}
	}

	s.logger.Info("command blacklist updated",
		zap.Int("total_commands", len(s.blacklistCommands)),
		zap.Int("total_directories", len(s.blacklistDirs)))

	return &types.UpdateCommandBlacklistResponse{
		Success: true,
		Message: types.MsgBlacklistUpdated,
	}, nil
}

// GetWorkingDirectory 获取当前工作目录 / Get working directory
func (s *Service) GetWorkingDirectory(_ *types.GetWorkingDirectoryRequest) (*types.GetWorkingDirectoryResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &types.GetWorkingDirectoryResponse{
		WorkDir: s.currentWorkDir,
	}, nil
}

// ChangeDirectory 切换工作目录 / Change directory
func (s *Service) ChangeDirectory(req *types.ChangeDirectoryRequest) (*types.ChangeDirectoryResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateChangeDirectoryRequest(req); err != nil {
		return nil, err
	}

	// 验证路径 / Validate path
	validPath, err := s.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	// 检查目录是否存在 / Check if directory exists
	info, err := os.Stat(validPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New(types.ErrPathNotFound)
		}
		return nil, fmt.Errorf("failed to stat directory: %w", err)
	}

	if !info.IsDir() {
		return nil, errors.New(types.ErrNotDirectory)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查目录是否在黑名单中 / Check if directory is in blacklist
	if s.isDirectoryBlacklisted(validPath) {
		s.logger.Warn("blocked blacklisted directory",
			zap.String("directory", validPath))
		return nil, errors.New(types.ErrDirectoryBlacklisted)
	}

	// 计算相对路径 / Calculate relative path
	relPath, err := filepath.Rel(s.sandboxDir, validPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path: %w", err)
	}

	s.currentWorkDir = relPath
	s.logger.Info("working directory changed",
		zap.String("new_work_dir", relPath))

	return &types.ChangeDirectoryResponse{
		Success: true,
		Message: types.MsgDirectoryChanged,
	}, nil
}

// isCommandBlacklisted 检查命令是否在黑名单中 / Check if command is in blacklist
func (s *Service) isCommandBlacklisted(command string) bool {
	command = strings.TrimSpace(strings.ToLower(command))

	// 提取命令名称(去除路径) / Extract command name (remove path)
	cmdName := filepath.Base(command)

	// 去除Windows的.exe后缀 / Remove .exe suffix on Windows
	if runtime.GOOS == "windows" {
		cmdName = strings.TrimSuffix(cmdName, ".exe")
	}

	for _, blacklisted := range s.blacklistCommands {
		if strings.ToLower(blacklisted) == cmdName {
			return true
		}
	}

	return false
}

// isDirectoryBlacklisted 检查目录是否在黑名单中 / Check if directory is in blacklist
func (s *Service) isDirectoryBlacklisted(dir string) bool {
	dir = filepath.Clean(dir)
	sandboxDir := filepath.Clean(s.sandboxDir)

	// 如果目录在沙箱内，则不检查黑名单（沙箱目录本身及其子目录都是安全的）
	// If directory is within sandbox, skip blacklist check (sandbox directory and its subdirectories are safe)
	if strings.HasPrefix(dir, sandboxDir) {
		return false
	}

	// 检查自定义黑名单 / Check custom blacklist
	for _, blacklisted := range s.blacklistDirs {
		blacklisted = filepath.Clean(blacklisted)
		if strings.HasPrefix(dir, blacklisted) {
			return true
		}
	}

	// 检查系统目录 / Check system directories
	for _, sysDir := range SystemDirectories {
		sysDir = filepath.Clean(sysDir)
		if strings.HasPrefix(dir, sysDir) {
			return true
		}
	}

	return false
}

// validateCommandPaths 验证命令参数中的路径是否在沙箱内 / Validate paths in command arguments are within sandbox
func (s *Service) validateCommandPaths(command string, args []string, workDir string) error {
	// 需要验证路径的命令列表 / Commands that need path validation
	pathSensitiveCommands := map[string]bool{
		"rm":     true,
		"rmdir":  true,
		"del":    true, // Windows
		"erase":  true, // Windows
		"rd":     true, // Windows
		"remove": true,
	}

	// 检查命令是否需要路径验证 / Check if command needs path validation
	cmdName := filepath.Base(command)
	cmdName = strings.ToLower(strings.TrimSuffix(cmdName, filepath.Ext(cmdName)))

	if !pathSensitiveCommands[cmdName] {
		return nil // 不需要验证 / No validation needed
	}

	// 解析并验证参数中的路径 / Parse and validate paths in arguments
	for _, arg := range args {
		// 跳过选项参数 / Skip option arguments
		if strings.HasPrefix(arg, "-") {
			continue
		}
		// Windows命令选项以/开头，但要排除路径(如C:/path)
		if runtime.GOOS == "windows" && strings.HasPrefix(arg, "/") && len(arg) > 1 && arg[1] != ':' && !strings.Contains(arg, "\\") && !strings.Contains(arg, "/") {
			continue
		}

		// 检查是否是路径参数 / Check if it's a path argument
		if arg == "" {
			continue
		}

		// 解析路径 / Parse path
		var targetPath string
		if filepath.IsAbs(arg) {
			targetPath = arg
		} else {
			targetPath = filepath.Join(workDir, arg)
		}

		// 清理路径 / Clean path
		targetPath = filepath.Clean(targetPath)

		// 验证路径是否在沙箱内 / Validate path is within sandbox
		sandboxDir := filepath.Clean(s.sandboxDir)
		if !strings.HasPrefix(targetPath, sandboxDir) {
			return fmt.Errorf("%s: path '%s' is outside sandbox directory", types.ErrSandboxViolation, arg)
		}

		// 检查路径是否在黑名单目录中 / Check if path is in blacklisted directory
		if s.isDirectoryBlacklisted(targetPath) {
			return fmt.Errorf("%s: path '%s' is in blacklisted directory", types.ErrDirectoryBlacklisted, arg)
		}
	}

	return nil
}
