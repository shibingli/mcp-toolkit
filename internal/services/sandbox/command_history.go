package sandbox

import (
	"time"

	"mcp-toolkit/pkg/types"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// addCommandHistory 添加命令执行历史记录 / Add command execution history
func (s *Service) addCommandHistory(entry *types.CommandHistoryEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 限制历史记录数量,保留最近1000条 / Limit history size, keep last 1000 entries
	if len(s.commandHistory) >= 1000 {
		s.commandHistory = s.commandHistory[1:]
	}

	s.commandHistory = append(s.commandHistory, entry)

	// 记录审计日志 / Log audit entry
	s.auditLogger.Info("command executed",
		zap.String("id", entry.ID),
		zap.String("command", entry.Command),
		zap.Strings("args", entry.Args),
		zap.String("work_dir", entry.WorkDir),
		zap.String("user", entry.User),
		zap.Int("permission_level", int(entry.PermissionLevel)),
		zap.Time("start_time", entry.StartTime),
		zap.Time("end_time", entry.EndTime),
		zap.Int64("duration_ms", entry.Duration),
		zap.Int("exit_code", entry.ExitCode),
		zap.Bool("success", entry.Success),
	)
}

// GetCommandHistory 获取命令执行历史 / Get command execution history
func (s *Service) GetCommandHistory(req *types.GetCommandHistoryRequest) (*types.GetCommandHistoryResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 设置默认值 / Set default values
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// 过滤历史记录 / Filter history
	var filtered []*types.CommandHistoryEntry
	for _, entry := range s.commandHistory {
		// 按用户过滤 / Filter by user
		if req.User != "" && entry.User != req.User {
			continue
		}
		filtered = append(filtered, entry)
	}

	total := len(filtered)

	// 应用分页 / Apply pagination
	start := offset
	if start > total {
		start = total
	}

	end := start + limit
	if end > total {
		end = total
	}

	result := filtered[start:end]

	return &types.GetCommandHistoryResponse{
		History: result,
		Total:   total,
	}, nil
}

// ClearCommandHistory 清空命令执行历史 / Clear command execution history
func (s *Service) ClearCommandHistory(req *types.ClearCommandHistoryRequest) (*types.ClearCommandHistoryResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := len(s.commandHistory)
	s.commandHistory = make([]*types.CommandHistoryEntry, 0, 100)

	s.logger.Info("command history cleared",
		zap.Int("cleared_count", count))

	s.auditLogger.Info("command history cleared",
		zap.Int("cleared_count", count))

	return &types.ClearCommandHistoryResponse{
		Success: true,
		Message: types.MsgSuccess,
	}, nil
}

// createHistoryEntry 创建历史记录条目 / Create history entry
func createHistoryEntry(command string, args []string, workDir string, startTime, endTime time.Time,
	exitCode int, success bool, user string, permissionLevel types.CommandPermissionLevel,
	environment map[string]string) *types.CommandHistoryEntry {

	return &types.CommandHistoryEntry{
		ID:              uuid.New().String(),
		Command:         command,
		Args:            args,
		WorkDir:         workDir,
		StartTime:       startTime,
		EndTime:         endTime,
		Duration:        endTime.Sub(startTime).Milliseconds(),
		ExitCode:        exitCode,
		Success:         success,
		User:            user,
		PermissionLevel: permissionLevel,
		Environment:     environment,
	}
}
