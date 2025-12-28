package sandbox

import (
	"fmt"
	"strings"

	"mcp-toolkit/pkg/types"

	"go.uber.org/zap"
)

// 只读权限允许的命令 / Read-only permission allowed commands
var readOnlyCommands = map[string]bool{
	"ls":     true,
	"dir":    true,
	"cat":    true,
	"type":   true,
	"echo":   true,
	"pwd":    true,
	"cd":     true,
	"find":   true,
	"grep":   true,
	"head":   true,
	"tail":   true,
	"wc":     true,
	"stat":   true,
	"file":   true,
	"which":  true,
	"where":  true,
	"whoami": true,
}

// 标准权限禁止的命令 / Standard permission prohibited commands
var standardProhibitedCommands = map[string]bool{
	"chmod":   true,
	"chown":   true,
	"chgrp":   true,
	"sudo":    true,
	"su":      true,
	"kill":    true,
	"killall": true,
}

// SetPermissionLevel 设置权限级别 / Set permission level
func (s *Service) SetPermissionLevel(req *types.SetPermissionLevelRequest) (*types.SetPermissionLevelResponse, error) {
	// 验证权限级别 / Validate permission level
	if req.Level < types.PermissionLevelReadOnly || req.Level > types.PermissionLevelAdmin {
		return nil, fmt.Errorf("invalid permission level: %d", req.Level)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	oldLevel := s.permissionLevel
	s.permissionLevel = req.Level

	s.logger.Info("permission level changed",
		zap.Int("old_level", int(oldLevel)),
		zap.Int("new_level", int(req.Level)))

	s.auditLogger.Info("permission level changed",
		zap.Int("old_level", int(oldLevel)),
		zap.Int("new_level", int(req.Level)))

	return &types.SetPermissionLevelResponse{
		Success: true,
		Message: "permission level updated successfully",
	}, nil
}

// GetPermissionLevel 获取当前权限级别 / Get current permission level
func (s *Service) GetPermissionLevel(req *types.GetPermissionLevelRequest) (*types.GetPermissionLevelResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &types.GetPermissionLevelResponse{
		Level: s.permissionLevel,
	}, nil
}

// checkCommandPermission 检查命令权限 / Check command permission
func (s *Service) checkCommandPermission(command string, requestLevel types.CommandPermissionLevel) error {
	s.mu.RLock()
	currentLevel := s.permissionLevel
	s.mu.RUnlock()

	// 使用请求中的权限级别,如果未指定则使用当前级别 / Use permission level from request, or current level if not specified
	level := requestLevel
	if level == 0 {
		level = currentLevel
	}

	// 提取命令基本名称 / Extract command base name
	cmdName := strings.ToLower(strings.TrimSpace(command))
	if idx := strings.LastIndexAny(cmdName, "/\\"); idx >= 0 {
		cmdName = cmdName[idx+1:]
	}

	// 只读权限检查 / Read-only permission check
	if level == types.PermissionLevelReadOnly {
		if !readOnlyCommands[cmdName] {
			return fmt.Errorf("command '%s' is not allowed with read-only permission", command)
		}
		return nil
	}

	// 标准权限检查 / Standard permission check
	if level == types.PermissionLevelStandard {
		if standardProhibitedCommands[cmdName] {
			return fmt.Errorf("command '%s' requires elevated permission", command)
		}
	}

	// 提升权限和管理员权限允许所有非黑名单命令 / Elevated and admin permissions allow all non-blacklisted commands
	// 黑名单检查在ExecuteCommand中进行 / Blacklist check is done in ExecuteCommand

	return nil
}

// getPermissionLevelName 获取权限级别名称 / Get permission level name
func getPermissionLevelName(level types.CommandPermissionLevel) string {
	switch level {
	case types.PermissionLevelReadOnly:
		return "read-only"
	case types.PermissionLevelStandard:
		return "standard"
	case types.PermissionLevelElevated:
		return "elevated"
	case types.PermissionLevelAdmin:
		return "admin"
	default:
		return "unknown"
	}
}

// requiresElevatedPermission 检查命令是否需要提升权限 / Check if command requires elevated permission
func requiresElevatedPermission(command string) bool {
	cmdName := strings.ToLower(strings.TrimSpace(command))
	if idx := strings.LastIndexAny(cmdName, "/\\"); idx >= 0 {
		cmdName = cmdName[idx+1:]
	}
	return standardProhibitedCommands[cmdName]
}

// isReadOnlyCommand 检查是否为只读命令 / Check if command is read-only
func isReadOnlyCommand(command string) bool {
	cmdName := strings.ToLower(strings.TrimSpace(command))
	if idx := strings.LastIndexAny(cmdName, "/\\"); idx >= 0 {
		cmdName = cmdName[idx+1:]
	}
	return readOnlyCommands[cmdName]
}
