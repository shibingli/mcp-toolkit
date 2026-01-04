package sandbox

import (
	"testing"

	"mcp-toolkit/pkg/types"

	"github.com/stretchr/testify/assert"
)

// TestGetPermissionLevelName 测试获取权限级别名称 / Test get permission level name
func TestGetPermissionLevelName(t *testing.T) {
	tests := []struct {
		name     string
		level    types.CommandPermissionLevel
		expected string
	}{
		{
			name:     "只读权限",
			level:    types.PermissionLevelReadOnly,
			expected: "read-only",
		},
		{
			name:     "标准权限",
			level:    types.PermissionLevelStandard,
			expected: "standard",
		},
		{
			name:     "提升权限",
			level:    types.PermissionLevelElevated,
			expected: "elevated",
		},
		{
			name:     "管理员权限",
			level:    types.PermissionLevelAdmin,
			expected: "admin",
		},
		{
			name:     "未知权限",
			level:    types.CommandPermissionLevel(999),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPermissionLevelName(tt.level)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestRequiresElevatedPermission 测试是否需要提升权限 / Test requires elevated permission
func TestRequiresElevatedPermission(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{
			name:     "chmod命令需要提升权限",
			command:  "chmod",
			expected: true,
		},
		{
			name:     "chown命令需要提升权限",
			command:  "chown",
			expected: true,
		},
		{
			name:     "chgrp命令需要提升权限",
			command:  "chgrp",
			expected: true,
		},
		{
			name:     "sudo命令需要提升权限",
			command:  "sudo",
			expected: true,
		},
		{
			name:     "su命令需要提升权限",
			command:  "su",
			expected: true,
		},
		{
			name:     "kill命令需要提升权限",
			command:  "kill",
			expected: true,
		},
		{
			name:     "killall命令需要提升权限",
			command:  "killall",
			expected: true,
		},
		{
			name:     "echo命令不需要提升权限",
			command:  "echo",
			expected: false,
		},
		{
			name:     "ls命令不需要提升权限",
			command:  "ls",
			expected: false,
		},
		{
			name:     "rm命令不需要提升权限",
			command:  "rm",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := requiresElevatedPermission(tt.command)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsReadOnlyCommand 测试是否为只读命令 / Test is read only command
func TestIsReadOnlyCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{
			name:     "ls是只读命令",
			command:  "ls",
			expected: true,
		},
		{
			name:     "dir是只读命令",
			command:  "dir",
			expected: true,
		},
		{
			name:     "cat是只读命令",
			command:  "cat",
			expected: true,
		},
		{
			name:     "type是只读命令",
			command:  "type",
			expected: true,
		},
		{
			name:     "head是只读命令",
			command:  "head",
			expected: true,
		},
		{
			name:     "tail是只读命令",
			command:  "tail",
			expected: true,
		},
		{
			name:     "grep是只读命令",
			command:  "grep",
			expected: true,
		},
		{
			name:     "find是只读命令",
			command:  "find",
			expected: true,
		},
		{
			name:     "echo是只读命令",
			command:  "echo",
			expected: true,
		},
		{
			name:     "pwd是只读命令",
			command:  "pwd",
			expected: true,
		},
		{
			name:     "whoami是只读命令",
			command:  "whoami",
			expected: true,
		},
		{
			name:     "rm不是只读命令",
			command:  "rm",
			expected: false,
		},
		{
			name:     "mkdir不是只读命令",
			command:  "mkdir",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isReadOnlyCommand(tt.command)
			assert.Equal(t, tt.expected, result)
		})
	}
}
