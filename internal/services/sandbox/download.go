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
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mcp-toolkit/pkg/types"

	"go.uber.org/zap"
)

const (
	// DefaultDownloadTimeout 默认下载超时时间（秒）/ Default download timeout in seconds
	DefaultDownloadTimeout = 30

	// MaxDownloadTimeout 最大下载超时时间（秒）/ Maximum download timeout in seconds
	MaxDownloadTimeout = 300

	// MaxDownloadFileSize 最大下载文件大小（字节）/ Maximum download file size in bytes
	// 默认限制为 1GB，防止磁盘空间耗尽 / Default limit is 1GB to prevent disk space exhaustion
	MaxDownloadFileSize = 1024 * 1024 * 1024 // 1GB

	// MaxHTTPRedirects 最大HTTP重定向次数 / Maximum HTTP redirects
	MaxHTTPRedirects = 10

	// DefaultHTTPMethod 默认HTTP方法 / Default HTTP method
	DefaultHTTPMethod = "GET"
)

// DownloadFile 下载文件 / Download file
func (s *Service) DownloadFile(req *types.DownloadFileRequest) (*types.DownloadFileResponse, error) {
	// 参数验证 / Parameter validation
	if err := validateDownloadFileRequest(req); err != nil {
		return nil, err
	}

	// 验证并获取保存路径 / Validate and get save path
	validPath, err := s.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	// 确保父目录存在 / Ensure parent directory exists
	dir := filepath.Dir(validPath)
	if err = os.MkdirAll(dir, DefaultDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create parent directory: %w", err)
	}

	// 设置默认值 / Set default values
	method := req.Method
	if method == "" {
		method = DefaultHTTPMethod
	}
	method = strings.ToUpper(method)

	timeout := req.Timeout
	if timeout <= 0 {
		timeout = DefaultDownloadTimeout
	}
	if timeout > MaxDownloadTimeout {
		timeout = MaxDownloadTimeout
	}

	// 创建HTTP客户端，限制重定向次数 / Create HTTP client with redirect limit
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= MaxHTTPRedirects {
				return fmt.Errorf("stopped after %d redirects", MaxHTTPRedirects)
			}
			return nil
		},
	}

	// 如果需要跳过TLS证书验证 / If skip TLS verification is needed
	if req.SkipTLSVerify {
		s.logger.Warn("skipping TLS certificate verification - this is insecure and should only be used in development",
			zap.String("url", req.URL))
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // #nosec G402 - This is intentionally configurable for development environments
			},
		}
	}

	// 创建HTTP请求 / Create HTTP request
	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(method, req.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// 设置请求头 / Set request headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// 如果没有设置User-Agent，使用默认值 / Set default User-Agent if not provided
	if httpReq.Header.Get("User-Agent") == "" {
		httpReq.Header.Set("User-Agent", fmt.Sprintf("%s/%s", types.ServerName, types.ServerVersion))
	}

	s.logger.Info("downloading file",
		zap.String("url", req.URL),
		zap.String("method", method),
		zap.String("path", validPath),
		zap.Int("timeout", timeout))

	// 发送HTTP请求 / Send HTTP request
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// 检查HTTP状态码 / Check HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	// 检查Content-Length，如果提供了则验证文件大小 / Check Content-Length if provided
	contentLength := resp.ContentLength
	if contentLength > 0 && contentLength > MaxDownloadFileSize {
		return nil, fmt.Errorf("file size (%d bytes) exceeds maximum allowed size (%d bytes)", contentLength, MaxDownloadFileSize)
	}

	// 创建目标文件 / Create target file
	file, err := os.Create(validPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	// 使用LimitReader限制读取大小，防止超大文件 / Use LimitReader to limit read size
	limitedReader := io.LimitReader(resp.Body, MaxDownloadFileSize+1)

	// 下载文件内容 / Download file content
	size, err := io.Copy(file, limitedReader)

	// 关闭文件句柄 / Close file handle
	_ = file.Close()

	if err != nil {
		// 下载失败，删除不完整的文件 / Delete incomplete file on failure
		_ = os.Remove(validPath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// 检查是否超过最大文件大小 / Check if file size exceeds maximum
	if size > MaxDownloadFileSize {
		_ = os.Remove(validPath)
		return nil, fmt.Errorf("file size (%d bytes) exceeds maximum allowed size (%d bytes)", size, MaxDownloadFileSize)
	}

	// 获取相对路径 / Get relative path
	relPath, err := filepath.Rel(s.sandboxDir, validPath)
	if err != nil {
		relPath = req.Path
	}

	s.logger.Info("file downloaded successfully",
		zap.String("url", req.URL),
		zap.String("path", validPath),
		zap.Int64("size", size))

	return &types.DownloadFileResponse{
		Success:      true,
		Message:      "file downloaded successfully",
		SandboxPath:  relPath,
		AbsolutePath: validPath,
		Size:         size,
		ContentType:  resp.Header.Get("Content-Type"),
	}, nil
}
