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

package sandbox

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"mcp-toolkit/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestDownloadFile 测试下载文件功能 / Test download file functionality
func TestDownloadFile(t *testing.T) {
	// 创建临时沙箱目录 / Create temporary sandbox directory
	sandboxDir := t.TempDir()
	logger := zap.NewNop()

	service, err := NewService(sandboxDir, logger)
	require.NoError(t, err)

	// 创建测试HTTP服务器 / Create test HTTP server
	testContent := "Hello, World! This is test content."
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(testContent))
	}))
	defer server.Close()

	t.Run("成功下载文件", func(t *testing.T) {
		req := &types.DownloadFileRequest{
			URL:  server.URL,
			Path: "downloads/test.txt",
		}

		resp, err := service.DownloadFile(req)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "file downloaded successfully", resp.Message)
		assert.Equal(t, filepath.Join("downloads", "test.txt"), resp.SandboxPath)
		assert.Greater(t, resp.Size, int64(0))
		assert.Equal(t, "text/plain", resp.ContentType)

		// 验证文件内容 / Verify file content
		content, err := os.ReadFile(filepath.Join(sandboxDir, "downloads", "test.txt"))
		require.NoError(t, err)
		assert.Equal(t, testContent, string(content))
	})

	t.Run("使用POST方法下载", func(t *testing.T) {
		postServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		}))
		defer postServer.Close()

		req := &types.DownloadFileRequest{
			URL:    postServer.URL,
			Path:   "downloads/post_response.json",
			Method: "POST",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"test":"data"}`,
		}

		resp, err := service.DownloadFile(req)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "application/json", resp.ContentType)
	})

	t.Run("自定义超时时间", func(t *testing.T) {
		req := &types.DownloadFileRequest{
			URL:     server.URL,
			Path:    "downloads/timeout_test.txt",
			Timeout: 60,
		}

		resp, err := service.DownloadFile(req)
		require.NoError(t, err)
		assert.True(t, resp.Success)
	})

	t.Run("自定义请求头", func(t *testing.T) {
		headerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "CustomValue", r.Header.Get("X-Custom-Header"))
			_, _ = w.Write([]byte("OK"))
		}))
		defer headerServer.Close()

		req := &types.DownloadFileRequest{
			URL:  headerServer.URL,
			Path: "downloads/header_test.txt",
			Headers: map[string]string{
				"X-Custom-Header": "CustomValue",
			},
		}

		resp, err := service.DownloadFile(req)
		require.NoError(t, err)
		assert.True(t, resp.Success)
	})

	t.Run("HTTP错误状态码", func(t *testing.T) {
		errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer errorServer.Close()

		req := &types.DownloadFileRequest{
			URL:  errorServer.URL,
			Path: "downloads/error.txt",
		}

		_, err := service.DownloadFile(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("无效URL", func(t *testing.T) {
		req := &types.DownloadFileRequest{
			URL:  "invalid-url",
			Path: "downloads/invalid.txt",
		}

		_, err := service.DownloadFile(req)
		assert.Error(t, err)
	})

	t.Run("空URL", func(t *testing.T) {
		req := &types.DownloadFileRequest{
			URL:  "",
			Path: "downloads/empty.txt",
		}

		_, err := service.DownloadFile(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "URL cannot be empty")
	})

	t.Run("空路径", func(t *testing.T) {
		req := &types.DownloadFileRequest{
			URL:  server.URL,
			Path: "",
		}

		_, err := service.DownloadFile(req)
		assert.Error(t, err)
	})

	t.Run("路径遍历攻击防护", func(t *testing.T) {
		req := &types.DownloadFileRequest{
			URL:  server.URL,
			Path: "../../../etc/passwd",
		}

		_, err := service.DownloadFile(req)
		assert.Error(t, err)
	})

	t.Run("自动创建父目录", func(t *testing.T) {
		req := &types.DownloadFileRequest{
			URL:  server.URL,
			Path: "deep/nested/directory/file.txt",
		}

		resp, err := service.DownloadFile(req)
		require.NoError(t, err)
		assert.True(t, resp.Success)

		// 验证目录已创建 / Verify directory was created
		dirPath := filepath.Join(sandboxDir, "deep/nested/directory")
		info, err := os.Stat(dirPath)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("默认HTTP方法为GET", func(t *testing.T) {
		methodServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			_, _ = w.Write([]byte("OK"))
		}))
		defer methodServer.Close()

		req := &types.DownloadFileRequest{
			URL:  methodServer.URL,
			Path: "downloads/default_method.txt",
		}

		resp, err := service.DownloadFile(req)
		require.NoError(t, err)
		assert.True(t, resp.Success)
	})

	t.Run("超时限制", func(t *testing.T) {
		req := &types.DownloadFileRequest{
			URL:     server.URL,
			Path:    "downloads/max_timeout.txt",
			Timeout: 500, // 超过最大值 / Exceeds maximum
		}

		resp, err := service.DownloadFile(req)
		require.NoError(t, err)
		assert.True(t, resp.Success)
	})

	t.Run("nil请求", func(t *testing.T) {
		_, err := service.DownloadFile(nil)
		assert.Error(t, err)
	})
}

// TestValidateDownloadFileRequest 测试下载文件请求验证 / Test download file request validation
func TestValidateDownloadFileRequest(t *testing.T) {
	t.Run("有效请求", func(t *testing.T) {
		req := &types.DownloadFileRequest{
			URL:  "https://example.com/file.txt",
			Path: "downloads/file.txt",
		}
		err := validateDownloadFileRequest(req)
		assert.NoError(t, err)
	})

	t.Run("nil请求", func(t *testing.T) {
		err := validateDownloadFileRequest(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})

	t.Run("空URL", func(t *testing.T) {
		req := &types.DownloadFileRequest{
			URL:  "",
			Path: "downloads/file.txt",
		}
		err := validateDownloadFileRequest(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "URL cannot be empty")
	})

	t.Run("空路径", func(t *testing.T) {
		req := &types.DownloadFileRequest{
			URL:  "https://example.com/file.txt",
			Path: "",
		}
		err := validateDownloadFileRequest(req)
		assert.Error(t, err)
	})

	t.Run("无效HTTP方法", func(t *testing.T) {
		req := &types.DownloadFileRequest{
			URL:    "https://example.com/file.txt",
			Path:   "downloads/file.txt",
			Method: "INVALID",
		}
		err := validateDownloadFileRequest(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid HTTP method")
	})

	t.Run("有效HTTP方法", func(t *testing.T) {
		methods := []string{"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH", "OPTIONS"}
		for _, method := range methods {
			req := &types.DownloadFileRequest{
				URL:    "https://example.com/file.txt",
				Path:   "downloads/file.txt",
				Method: method,
			}
			err := validateDownloadFileRequest(req)
			assert.NoError(t, err, fmt.Sprintf("Method %s should be valid", method))
		}
	})

	t.Run("路径过长", func(t *testing.T) {
		longPath := string(make([]byte, MaxPathLength+1))
		req := &types.DownloadFileRequest{
			URL:  "https://example.com/file.txt",
			Path: longPath,
		}
		err := validateDownloadFileRequest(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds maximum")
	})

	t.Run("非HTTP/HTTPS URL", func(t *testing.T) {
		req := &types.DownloadFileRequest{
			URL:  "ftp://example.com/file.txt",
			Path: "downloads/file.txt",
		}
		err := validateDownloadFileRequest(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid URL format")
	})
}

// TestIsValidURL 测试URL验证函数 / Test URL validation function
func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"有效HTTP URL", "http://example.com", true},
		{"有效HTTPS URL", "https://example.com", true},
		{"有效HTTP URL带路径", "http://example.com/path/to/file", true},
		{"有效HTTPS URL带路径", "https://example.com/path/to/file", true},
		{"无效URL - 空字符串", "", false},
		{"无效URL - 无协议", "example.com", false},
		{"无效URL - FTP协议", "ftp://example.com", false},
		{"无效URL - 只有协议", "http://", false},
		{"无效URL - 格式错误", "not-a-url", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidURL(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDownloadFileSizeLimit 测试文件大小限制 / Test file size limit
func TestDownloadFileSizeLimit(t *testing.T) {
	service, sandboxDir := setupTestService(t)
	defer cleanupTestService(t, sandboxDir)

	t.Run("Content-Length超过限制", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 设置Content-Length为超过限制的大小
			w.Header().Set("Content-Length", fmt.Sprintf("%d", MaxDownloadFileSize+1))
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		req := &types.DownloadFileRequest{
			URL:  server.URL,
			Path: "downloads/large-file.bin",
		}

		_, err := service.DownloadFile(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds maximum allowed size")
	})

	t.Run("实际下载大小超过限制", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 不设置Content-Length，但写入超过限制的数据
			w.WriteHeader(http.StatusOK)
			// 写入超过限制的数据
			data := make([]byte, 1024*1024) // 1MB chunks
			for i := int64(0); i < MaxDownloadFileSize/int64(len(data))+2; i++ {
				_, _ = w.Write(data)
			}
		}))
		defer server.Close()

		req := &types.DownloadFileRequest{
			URL:  server.URL,
			Path: "downloads/large-file.bin",
		}

		_, err := service.DownloadFile(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds maximum allowed size")

		// 验证不完整的文件已被删除
		filePath := filepath.Join(sandboxDir, "downloads/large-file.bin")
		_, err = os.Stat(filePath)
		assert.True(t, os.IsNotExist(err))
	})
}

// TestDownloadFileRedirectLimit 测试重定向限制 / Test redirect limit
func TestDownloadFileRedirectLimit(t *testing.T) {
	service, sandboxDir := setupTestService(t)
	defer cleanupTestService(t, sandboxDir)

	t.Run("重定向次数过多", func(t *testing.T) {
		redirectCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			redirectCount++
			if redirectCount <= MaxHTTPRedirects+1 {
				// 继续重定向
				http.Redirect(w, r, r.URL.String()+"?redirect="+fmt.Sprint(redirectCount), http.StatusFound)
			} else {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("final content"))
			}
		}))
		defer server.Close()

		req := &types.DownloadFileRequest{
			URL:  server.URL,
			Path: "downloads/redirect.txt",
		}

		_, err := service.DownloadFile(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stopped after")
	})

	t.Run("正常重定向", func(t *testing.T) {
		redirectCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			redirectCount++
			if redirectCount < 3 {
				// 重定向3次（在限制内）
				http.Redirect(w, r, r.URL.String()+"?redirect="+fmt.Sprint(redirectCount), http.StatusFound)
			} else {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("final content"))
			}
		}))
		defer server.Close()

		req := &types.DownloadFileRequest{
			URL:  server.URL,
			Path: "downloads/redirect.txt",
		}

		resp, err := service.DownloadFile(req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Success)

		// 验证文件内容
		filePath := filepath.Join(sandboxDir, "downloads/redirect.txt")
		content, err := os.ReadFile(filePath)
		assert.NoError(t, err)
		assert.Equal(t, "final content", string(content))
	})
}

// TestDownloadFileSkipTLSVerify 测试跳过TLS证书验证 / Test skip TLS certificate verification
func TestDownloadFileSkipTLSVerify(t *testing.T) {
	service, sandboxDir := setupTestService(t)
	defer cleanupTestService(t, sandboxDir)

	t.Run("使用自签名证书的HTTPS服务器", func(t *testing.T) {
		// 创建一个使用自签名证书的HTTPS测试服务器
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("secure content"))
		}))
		defer server.Close()

		// 不跳过证书验证应该失败
		t.Run("不跳过证书验证-应该失败", func(t *testing.T) {
			req := &types.DownloadFileRequest{
				URL:           server.URL,
				Path:          "downloads/secure-fail.txt",
				SkipTLSVerify: false,
			}

			_, err := service.DownloadFile(req)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "certificate")
		})

		// 跳过证书验证应该成功
		t.Run("跳过证书验证-应该成功", func(t *testing.T) {
			req := &types.DownloadFileRequest{
				URL:           server.URL,
				Path:          "downloads/secure-success.txt",
				SkipTLSVerify: true,
			}

			resp, err := service.DownloadFile(req)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.True(t, resp.Success)

			// 验证文件内容
			filePath := filepath.Join(sandboxDir, "downloads/secure-success.txt")
			content, err := os.ReadFile(filePath)
			assert.NoError(t, err)
			assert.Equal(t, "secure content", string(content))
		})
	})

	t.Run("普通HTTP不受影响", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("http content"))
		}))
		defer server.Close()

		req := &types.DownloadFileRequest{
			URL:           server.URL,
			Path:          "downloads/http.txt",
			SkipTLSVerify: true, // HTTP不需要TLS，这个参数应该被忽略
		}

		resp, err := service.DownloadFile(req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Success)
	})
}
