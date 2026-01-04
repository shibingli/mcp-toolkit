package sandbox

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetSystemDirectories 测试获取系统目录列表 / Test get system directories
func TestGetSystemDirectories(t *testing.T) {
	dirs := getSystemDirectories()

	// 验证返回的目录列表不为空 / Verify the directory list is not empty
	assert.NotEmpty(t, dirs)

	// 根据操作系统验证特定目录 / Verify specific directories based on OS
	if runtime.GOOS == "windows" {
		// Windows 系统应该包含这些目录 / Windows should contain these directories
		found := false
		for _, dir := range dirs {
			if dir == "C:\\Windows" || dir == "C:\\Program Files" || dir == "C:\\System32" {
				found = true
				break
			}
		}
		assert.True(t, found, "Windows系统目录列表应包含系统关键目录")
	} else {
		// Unix/Linux 系统应该包含这些目录 / Unix/Linux should contain these directories
		found := false
		for _, dir := range dirs {
			if dir == "/bin" || dir == "/sbin" || dir == "/usr" || dir == "/etc" {
				found = true
				break
			}
		}
		assert.True(t, found, "Unix/Linux系统目录列表应包含系统关键目录")
	}
}
