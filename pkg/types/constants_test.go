// Copyright 2024 MCP Toolkit Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestGetDefaultSandboxDir 测试GetDefaultSandboxDir函数 / Test GetDefaultSandboxDir function
func TestGetDefaultSandboxDir(t *testing.T) {
	// 保存原始环境变量 / Save original environment variables
	originalTemp := os.Getenv("TEMP")
	originalTmp := os.Getenv("TMP")
	defer func() {
		// 恢复原始环境变量 / Restore original environment variables
		if originalTemp != "" {
			_ = os.Setenv("TEMP", originalTemp)
		} else {
			_ = os.Unsetenv("TEMP")
		}
		if originalTmp != "" {
			_ = os.Setenv("TMP", originalTmp)
		} else {
			_ = os.Unsetenv("TMP")
		}
	}()

	tests := []struct {
		name     string
		setup    func()
		validate func(t *testing.T, result string)
	}{
		{
			name: "Windows with TEMP environment variable",
			setup: func() {
				if runtime.GOOS == "windows" {
					_ = os.Setenv("TEMP", `C:\Users\Test\AppData\Local\Temp`)
				}
			},
			validate: func(t *testing.T, result string) {
				if runtime.GOOS == "windows" {
					expected := filepath.Join(`C:\Users\Test\AppData\Local\Temp`, "mcp_sandbox_toolkit")
					if result != expected {
						t.Errorf("Expected %s, got %s", expected, result)
					}
				}
			},
		},
		{
			name: "Windows with TMP environment variable",
			setup: func() {
				if runtime.GOOS == "windows" {
					_ = os.Unsetenv("TEMP")
					_ = os.Setenv("TMP", `C:\Temp`)
				}
			},
			validate: func(t *testing.T, result string) {
				if runtime.GOOS == "windows" {
					expected := filepath.Join(`C:\Temp`, "mcp_sandbox_toolkit")
					if result != expected {
						t.Errorf("Expected %s, got %s", expected, result)
					}
				}
			},
		},
		{
			name: "Windows without TEMP/TMP environment variables",
			setup: func() {
				if runtime.GOOS == "windows" {
					_ = os.Unsetenv("TEMP")
					_ = os.Unsetenv("TMP")
				}
			},
			validate: func(t *testing.T, result string) {
				if runtime.GOOS == "windows" {
					expected := `C:\Temp\mcp_sandbox_toolkit`
					if result != expected {
						t.Errorf("Expected %s, got %s", expected, result)
					}
				}
			},
		},
		{
			name:  "Linux/MacOS default path",
			setup: func() {},
			validate: func(t *testing.T, result string) {
				if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
					expected := "/tmp/mcp_sandbox_toolkit"
					if result != expected {
						t.Errorf("Expected %s, got %s", expected, result)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result := GetDefaultSandboxDir()
			tt.validate(t, result)

			// 验证路径不为空 / Verify path is not empty
			if result == "" {
				t.Error("GetDefaultSandboxDir returned empty string")
			}

			// 验证路径包含mcp_sandbox_toolkit / Verify path contains mcp_sandbox_toolkit
			if !strings.Contains(result, "mcp_sandbox_toolkit") && result != DefaultSandboxDir {
				t.Errorf("Expected path to contain 'mcp_sandbox_toolkit', got %s", result)
			}
		})
	}
}

// TestGetDefaultSandboxDir_CurrentOS 测试当前操作系统的默认沙箱路径 / Test default sandbox path for current OS
func TestGetDefaultSandboxDir_CurrentOS(t *testing.T) {
	result := GetDefaultSandboxDir()

	// 验证返回值不为空 / Verify result is not empty
	if result == "" {
		t.Fatal("GetDefaultSandboxDir returned empty string")
	}

	// 根据操作系统验证路径格式 / Verify path format based on OS
	switch runtime.GOOS {
	case "windows":
		// Windows路径应该包含反斜杠或使用filepath.Separator / Windows path should contain backslash or use filepath.Separator
		if !strings.Contains(result, "mcp_sandbox_toolkit") {
			t.Errorf("Windows path should contain 'mcp_sandbox_toolkit', got %s", result)
		}
	case "linux", "darwin":
		// Linux/MacOS路径应该以/tmp开头 / Linux/MacOS path should start with /tmp
		if !strings.HasPrefix(result, "/tmp/") {
			t.Errorf("Linux/MacOS path should start with '/tmp/', got %s", result)
		}
	default:
		// 其他系统应该返回DefaultSandboxDir / Other systems should return DefaultSandboxDir
		if result != DefaultSandboxDir {
			t.Errorf("Unknown OS should return DefaultSandboxDir, got %s", result)
		}
	}

	t.Logf("Current OS: %s, Default sandbox dir: %s", runtime.GOOS, result)
}
